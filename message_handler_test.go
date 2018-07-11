package prekeyserver

import (
	. "gopkg.in/check.v1"
)

func (s *GenericServerSuite) Test_otrngMessageHandler_handleMessage_errorsOnMessageParsing(c *C) {
	_, e := (&otrngMessageHandler{}).handleMessage("", []byte{0x01, 0x02, 0x03, 0x04})
	c.Assert(e, ErrorMatches, "invalid protocol version")
}

func (s *GenericServerSuite) Test_otrngMessageHandler_handleMessage_errorsOnIncorrectMessageReturn(c *C) {
	msg := generateStorageInformationRequestMessage([]byte{0x1, 0x02})
	_, e := (&otrngMessageHandler{}).handleMessage("", msg.serialize())
	c.Assert(e, ErrorMatches, "invalid toplevel message")
}

func (s *GenericServerSuite) Test_otrngMessageHandler_handleMessage_errorsOnErrorsFromResponse(c *C) {
	stor := createInMemoryStorage()
	serverKey := deriveEDDSAKeypair([symKeyLength]byte{0x25, 0x25, 0x25, 0x25, 0x25, 0x25, 0x25, 0x25})
	gs := &GenericServer{
		identity:    "masterOfKeys.example.org",
		rand:        fixtureRand(),
		key:         serverKey,
		fingerprint: serverKey.pub.fingerprint(),
		storageImpl: stor,
	}
	d1 := generateDake1(sita.instanceTag, sita.clientProfile, gs.key.pub.k)
	_, e := (&otrngMessageHandler{s: gs}).handleMessage("someone@somewhere.org", d1.serialize())
	c.Assert(e, ErrorMatches, "invalid ring signature generation")
}