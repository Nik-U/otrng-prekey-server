package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"os"
	"time"

	pks "github.com/otrv4/otrng-prekey-server"
	. "gopkg.in/check.v1"
)

func (s *RawServerSuite) Test_flowTest_success(c *C) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	defer func() {
		os.Stdout = old
	}()
	os.Stdout = w
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	rs := &rawServer{}
	rs.load(pks.CreateFactory(fixtureRand()))
	*listenIP = "localhost"
	*listenPort = 0
	go rs.run()

	time.Sleep(time.Duration(100) * time.Millisecond)

	a := rs.l.Addr()
	con, _ := net.Dial(a.Network(), a.String())
	defer con.Close()

	ensembleRetrievalQueryMessage := "AAQQEkRVEQAAABBzaXRhQGV4YW1wbGUub3JnAAAAAQQ=."
	expectedResult := "AAQOEkRVEQAAAC5ObyBQcmVrZXkgTWVzc2FnZXMgYXZhaWxhYmxlIGZvciB0aGlzIGlkZW50aXR5."
	from := "rama@example.org"

	toSend := []byte{}
	toSend = appendShort(toSend, uint16(len(from)))
	toSend = append(toSend, []byte(from)...)
	toSend = appendShort(toSend, uint16(len(ensembleRetrievalQueryMessage)))
	toSend = append(toSend, []byte(ensembleRetrievalQueryMessage)...)

	n, e := con.Write(toSend)
	c.Assert(e, IsNil)
	c.Assert(n, Equals, 65)

	res, e := ioutil.ReadAll(con)
	c.Assert(e, IsNil)
	_, ss, _ := extractShort(res)
	c.Assert(ss, Equals, uint16(77))
	c.Assert(string(res[2:]), Equals, expectedResult)

	w.Close()
	sout := <-outC
	c.Assert(sout, Equals,
		"Starting server on localhost:0...\n"+
			"  [3B72D580C05DE282 3A14B02B682636BF 58F291A7E831D237 ECE8FC14DA50A187 A50ACF665442AB2D 140E140B813CFCCA 993BC02AA4A3D35C]\n")

	rs.l.Close()
}