package prekeyserver

import (
	"time"

	"github.com/twstrike/ed448"
	. "gopkg.in/check.v1"
)

func generatePointFrom(data [symKeyLength]byte) ed448.Point {
	return generateECDHPublicKeyFrom(data).k
}

func generateECDHPublicKeyFrom(data [symKeyLength]byte) *publicKey {
	return deriveECDHKeypair(data).pub
}

func generateEDDSAPublicKeyFrom(data [symKeyLength]byte) *publicKey {
	return deriveEDDSAKeypair(data).pub
}

func (s *GenericServerSuite) Test_dake1Message_shouldSerializeCorrectly(c *C) {
	d1 := &dake1Message{}
	d1.instanceTag = 0x4253112A
	d1.i = generatePointFrom([symKeyLength]byte{0x42, 0x11, 0xAA, 0xDE, 0xAD, 0xBE, 0xEF})
	d1.clientProfile = &clientProfile{}
	d1.clientProfile.identifier = 0xABCDEF11
	d1.clientProfile.instanceTag = 0x4253112A
	d1.clientProfile.publicKey = generateEDDSAPublicKeyFrom([symKeyLength]byte{0xAB, 0x42})
	d1.clientProfile.versions = []byte{0x04}
	d1.clientProfile.expiration = time.Date(2034, 11, 5, 13, 46, 00, 12, time.UTC)
	d1.clientProfile.dsaKey = nil
	d1.clientProfile.transitionalSignature = nil
	d1.clientProfile.sig = &eddsaSignature{
		s: [114]byte{0x15, 0x00, 0x00, 0x00, 0x12},
	}

	expected := []byte{
		// version
		0x00, 0x04,

		// message type
		0x01,

		// instance tag
		0x42, 0x53, 0x11, 0x2A,

		// client profile:
		0x0, 0x0, 0x0, 0x5,

		// identifier
		0x0, 0x1, 0xab, 0xcd, 0xef, 0x11,

		// instance tag
		0x0, 0x2, 0x42, 0x53, 0x11, 0x2a,

		// public key
		0x00, 0x03, 0x85, 0x9f, 0x37, 0x1f, 0xf3, 0x4f,
		0x36, 0x44, 0x5a, 0x99, 0xca, 0x8a, 0x11, 0x17,
		0x6b, 0xb8, 0x1e, 0xe0, 0x60, 0x39, 0x32, 0x76,
		0x71, 0xf4, 0xc6, 0x83, 0x77, 0x01, 0x45, 0x27,
		0x35, 0x3c, 0x75, 0xae, 0xee, 0xaa, 0xf9, 0x79,
		0x69, 0xa0, 0xd8, 0x9a, 0x3a, 0xb1, 0x48, 0xf6,
		0x44, 0x41, 0x83, 0x30, 0x9f, 0x41, 0x38, 0x1b,
		0xf3, 0x29, 0x00,

		// versions
		0x00, 0x05, 0x00, 0x00, 0x00, 0x01, 0x04,

		// expiry
		0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x79, 0xf8,
		0xc7, 0x98,

		// signature
		0x15, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00,

		// i:
		0x62, 0x38, 0x7d, 0xcd, 0x13, 0x84, 0x21, 0x0e,
		0x62, 0xcf, 0xaf, 0x06, 0x7f, 0x49, 0x02, 0x8c,
		0xdd, 0xfe, 0x99, 0xb9, 0x01, 0x59, 0x66, 0x7d,
		0x57, 0x0d, 0xc0, 0xb7, 0x89, 0x2c, 0xfc, 0x5c,
		0xac, 0xb8, 0x24, 0x17, 0xe9, 0x4d, 0x36, 0x29,
		0x04, 0x0e, 0x6a, 0xd1, 0xb4, 0x2d, 0x1a, 0x55,
		0xb9, 0x24, 0x29, 0x23, 0x7e, 0x5b, 0xc9, 0xe6,
		0x00,
	}

	c.Assert(d1.serialize(), DeepEquals, expected)
}

func (s *GenericServerSuite) Test_dake1Message_shouldDeserializeCorrectly(c *C) {
	d1 := &dake1Message{}
	e := d1.deserialize([]byte{
		// version
		0x00, 0x04,

		// message type
		0x01,

		// instance tag
		0x42, 0x53, 0x11, 0x2A,

		// client profile:
		0x0, 0x0, 0x0, 0x5,

		// identifier
		0x0, 0x1, 0xab, 0xcd, 0xef, 0x11,

		// instance tag
		0x0, 0x2, 0x42, 0x53, 0x11, 0x2a,

		// public key
		0x00, 0x03, 0x85, 0x9f, 0x37, 0x1f, 0xf3, 0x4f,
		0x36, 0x44, 0x5a, 0x99, 0xca, 0x8a, 0x11, 0x17,
		0x6b, 0xb8, 0x1e, 0xe0, 0x60, 0x39, 0x32, 0x76,
		0x71, 0xf4, 0xc6, 0x83, 0x77, 0x01, 0x45, 0x27,
		0x35, 0x3c, 0x75, 0xae, 0xee, 0xaa, 0xf9, 0x79,
		0x69, 0xa0, 0xd8, 0x9a, 0x3a, 0xb1, 0x48, 0xf6,
		0x44, 0x41, 0x83, 0x30, 0x9f, 0x41, 0x38, 0x1b,
		0xf3, 0x29, 0x00,

		// versions
		0x00, 0x05, 0x00, 0x00, 0x00, 0x01, 0x04,

		// expiry
		0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x79, 0xf8,
		0xc7, 0x98,

		// signature
		0x15, 0x00, 0x00, 0x00, 0x12, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0,
		0x00, 0x00,

		// i:
		0x62, 0x38, 0x7d, 0xcd, 0x13, 0x84, 0x21, 0x0e,
		0x62, 0xcf, 0xaf, 0x06, 0x7f, 0x49, 0x02, 0x8c,
		0xdd, 0xfe, 0x99, 0xb9, 0x01, 0x59, 0x66, 0x7d,
		0x57, 0x0d, 0xc0, 0xb7, 0x89, 0x2c, 0xfc, 0x5c,
		0xac, 0xb8, 0x24, 0x17, 0xe9, 0x4d, 0x36, 0x29,
		0x04, 0x0e, 0x6a, 0xd1, 0xb4, 0x2d, 0x1a, 0x55,
		0xb9, 0x24, 0x29, 0x23, 0x7e, 0x5b, 0xc9, 0xe6,
		0x00,
	})

	c.Assert(e, IsNil)
	c.Assert(d1.instanceTag, Equals, uint32(0x4253112A))
	c.Assert(d1.clientProfile.identifier, Equals, uint32(0xABCDEF11))
	c.Assert(d1.clientProfile.instanceTag, Equals, uint32(0x4253112A))
	c.Assert(d1.clientProfile.publicKey.k.Equals(generateEDDSAPublicKeyFrom([symKeyLength]byte{0xAB, 0x42}).k), Equals, true)
	c.Assert(d1.clientProfile.versions, DeepEquals, []byte{0x04})
	c.Assert(d1.clientProfile.expiration, Equals, time.Date(2034, 11, 5, 13, 46, 00, 00, time.UTC))
	c.Assert(d1.clientProfile.dsaKey, IsNil)
	c.Assert(d1.clientProfile.transitionalSignature, IsNil)
	c.Assert(d1.clientProfile.sig, DeepEquals, &eddsaSignature{
		s: [114]byte{0x15, 0x00, 0x00, 0x00, 0x12},
	})
	c.Assert(d1.i.Equals(generatePointFrom([symKeyLength]byte{0x42, 0x11, 0xAA, 0xDE, 0xAD, 0xBE, 0xEF})), Equals, true)
}