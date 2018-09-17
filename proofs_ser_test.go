package prekeyserver

import (
	"math/big"

	"github.com/otrv4/ed448"
	. "gopkg.in/check.v1"
)

func (s *GenericServerSuite) Test_ecdhProof_shouldSerializeCorrectly(c *C) {
	p := &ecdhProof{}
	val := [64]byte{0x42, 0x01, 0x00, 0x02}
	p.c = val[:]
	p.v = ed448.NewScalar([]byte{0xFF, 0xFF, 0xF1, 0xD0, 0x01})

	expected := []byte{
		0x42, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		0xff, 0xff, 0xf1, 0xd0, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	c.Assert(p.serialize(), DeepEquals, expected)
}

func (s *GenericServerSuite) Test_ecdhProof_shouldDeserializeCorrectly(c *C) {
	p := &ecdhProof{}
	_, ok := p.deserialize([]byte{
		0x42, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		0xff, 0xff, 0xf1, 0xd0, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})

	c.Assert(ok, Equals, true)
	val := [64]byte{0x42, 0x01, 0x00, 0x02}
	c.Assert(p.c, DeepEquals, val[:])
	c.Assert(p.v.Equals(ed448.NewScalar([]byte{0xFF, 0xFF, 0xF1, 0xD0, 0x01})), Equals, true)
}

func (s *GenericServerSuite) Test_ecdhProof_shouldFailOnMissingC(c *C) {
	p := &ecdhProof{}
	vv, ok := p.deserialize([]byte{
		0x42, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})

	c.Assert(ok, Equals, false)
	c.Assert(vv, IsNil)
}

func (s *GenericServerSuite) Test_ecdhProof_shouldFailOnMissingV(c *C) {
	p := &ecdhProof{}
	vv, ok := p.deserialize([]byte{
		0x42, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		0xff, 0xff, 0xf1, 0xd0, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})

	c.Assert(ok, Equals, false)
	c.Assert(vv, IsNil)
}

func (s *GenericServerSuite) Test_dhProof_shouldSerializeCorrectly(c *C) {
	p := &dhProof{}
	val := [64]byte{0x43, 0x01, 0x00, 0x02}
	p.c = val[:]
	p.v = new(big.Int).SetBytes([]byte{0xF0, 0xFF, 0xF1, 0xD0, 0x01})

	expected := []byte{
		0x43, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		0x00, 0x00, 0x00, 0x05, 0xf0, 0xff, 0xf1, 0xd0,
		0x01,
	}

	c.Assert(p.serialize(), DeepEquals, expected)
}

func (s *GenericServerSuite) Test_dhProof_shouldDeserializeCorrectly(c *C) {
	p := &dhProof{}
	_, ok := p.deserialize([]byte{
		0x43, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		0x00, 0x00, 0x00, 0x05, 0xf0, 0xff, 0xf1, 0xd0,
		0x01,
	})

	c.Assert(ok, Equals, true)
	val := [64]byte{0x43, 0x01, 0x00, 0x02}
	c.Assert(p.c, DeepEquals, val[:])
	c.Assert(p.v.Cmp(new(big.Int).SetBytes([]byte{0xF0, 0xFF, 0xF1, 0xD0, 0x01})), Equals, 0)
}

func (s *GenericServerSuite) Test_dhProof_shouldFailOnMissingC(c *C) {
	p := &dhProof{}
	vv, ok := p.deserialize([]byte{
		0x42, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})

	c.Assert(ok, Equals, false)
	c.Assert(vv, IsNil)
}

func (s *GenericServerSuite) Test_dhProof_shouldFailOnMissingV(c *C) {
	p := &dhProof{}
	vv, ok := p.deserialize([]byte{
		0x42, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		0x00, 0x00, 0x00, 0x05, 0xf0, 0xff, 0xf1, 0xd0,
	})

	c.Assert(ok, Equals, false)
	c.Assert(vv, IsNil)
}