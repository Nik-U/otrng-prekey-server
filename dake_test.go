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

func generateScalarFrom(data ...byte) ed448.Scalar {
	v := [privKeyLength]byte{}
	copy(v[:], data[:])
	return ed448.NewScalar(v[:])
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
	_, ok := d1.deserialize([]byte{
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

	c.Assert(ok, Equals, true)
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

func (s *GenericServerSuite) Test_dake2Message_shouldSerializeCorrectly(c *C) {
	d2 := &dake2Message{}
	d2.instanceTag = 0x4253112B
	d2.serverIdentity = []byte("prekey1.example.org")
	d2.serverFingerprint = fingerprint{0x13, 0x14, 0x15, 0xAB, 0xCC, 0xAC, 0xDC}
	d2.s = generatePointFrom([symKeyLength]byte{0x41, 0x12, 0xAC, 0xDF, 0xBD, 0xBF, 0xFE})
	d2.sigma = &ringSignature{
		c1: generateScalarFrom(0x01, 0x42, 0x12, 0xAB, 0xFC),
		r1: generateScalarFrom(0x01, 0x42, 0x12, 0xAB, 0xFD),
		c2: generateScalarFrom(0x02, 0x42, 0x15, 0xAB, 0xFC),
		r2: generateScalarFrom(0x03, 0x42, 0x16, 0xAC, 0xFD),
		c3: generateScalarFrom(0x04, 0x42, 0x17, 0xAD, 0xFC),
		r3: generateScalarFrom(0x05, 0x42, 0x18, 0xAE, 0xFD),
	}

	expected := []byte{
		// version
		0x00, 0x04,

		// message type
		0x02,

		// instance tag
		0x42, 0x53, 0x11, 0x2B,

		// prekey server composite identity
		// identity
		0x00, 0x00, 0x00, 0x13,
		0x70, 0x72, 0x65, 0x6b, 0x65, 0x79, 0x31, 0x2e,
		0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e,
		0x6f, 0x72, 0x67,

		// prekey server composite identity
		// fingerprint
		0x00, 0x00, 0x00, 0x38,
		0x13, 0x14, 0x15, 0xab, 0xcc, 0xac, 0xdc, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// s (point)
		0xd1, 0xad, 0xbe, 0x3a, 0xdd, 0x60, 0xc4, 0xbf,
		0xe0, 0xd8, 0x02, 0x85, 0x5b, 0x60, 0x6c, 0x3d,
		0xc3, 0x0a, 0x18, 0x6c, 0x77, 0xdc, 0xf8, 0x15,
		0xa0, 0x3b, 0x28, 0x20, 0x3c, 0xde, 0x5e, 0xa8,
		0x86, 0xf2, 0xa0, 0x93, 0xfa, 0xbd, 0x46, 0xd6,
		0x29, 0xef, 0x85, 0x4e, 0xfd, 0x5b, 0x3e, 0xa6,
		0x96, 0xeb, 0x17, 0x4f, 0x0a, 0xdb, 0x30, 0xf0,
		0x80,

		// sigma
		// sigma c1
		0x01, 0x42, 0x12, 0xab, 0xfc, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// sigma r1
		0x01, 0x42, 0x12, 0xab, 0xfd, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// sigma c2
		0x02, 0x42, 0x15, 0xab, 0xfc, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// sigma r2
		0x03, 0x42, 0x16, 0xac, 0xfd, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// sigma c3
		0x04, 0x42, 0x17, 0xad, 0xfc, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// sigma r3
		0x05, 0x42, 0x18, 0xae, 0xfd, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	c.Assert(d2.serialize(), DeepEquals, expected)
}

func (s *GenericServerSuite) Test_dake2Message_shouldDeserializeCorrectly(c *C) {
	d2 := &dake2Message{}
	_, ok := d2.deserialize([]byte{
		// version
		0x00, 0x04,

		// message type
		0x02,

		// instance tag
		0x42, 0x53, 0x11, 0x2B,

		// prekey server composite identity
		// identity
		0x00, 0x00, 0x00, 0x13,
		0x70, 0x72, 0x65, 0x6b, 0x65, 0x79, 0x31, 0x2e,
		0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e,
		0x6f, 0x72, 0x67,

		// prekey server composite identity
		// fingerprint
		0x00, 0x00, 0x00, 0x38,
		0x13, 0x14, 0x15, 0xab, 0xcc, 0xac, 0xdc, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// s (point)
		0xd1, 0xad, 0xbe, 0x3a, 0xdd, 0x60, 0xc4, 0xbf,
		0xe0, 0xd8, 0x02, 0x85, 0x5b, 0x60, 0x6c, 0x3d,
		0xc3, 0x0a, 0x18, 0x6c, 0x77, 0xdc, 0xf8, 0x15,
		0xa0, 0x3b, 0x28, 0x20, 0x3c, 0xde, 0x5e, 0xa8,
		0x86, 0xf2, 0xa0, 0x93, 0xfa, 0xbd, 0x46, 0xd6,
		0x29, 0xef, 0x85, 0x4e, 0xfd, 0x5b, 0x3e, 0xa6,
		0x96, 0xeb, 0x17, 0x4f, 0x0a, 0xdb, 0x30, 0xf0,
		0x80,

		// sigma
		// sigma c1
		0x01, 0x42, 0x12, 0xab, 0xfc, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// sigma r1
		0x01, 0x42, 0x12, 0xab, 0xfd, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// sigma c2
		0x02, 0x42, 0x15, 0xab, 0xfc, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// sigma r2
		0x03, 0x42, 0x16, 0xac, 0xfd, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// sigma c3
		0x04, 0x42, 0x17, 0xad, 0xfc, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// sigma r3
		0x05, 0x42, 0x18, 0xae, 0xfd, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})

	c.Assert(ok, Equals, true)
	c.Assert(d2.instanceTag, Equals, uint32(0x4253112B))
	c.Assert(d2.serverIdentity, DeepEquals, []byte("prekey1.example.org"))
	c.Assert(d2.serverFingerprint, DeepEquals, fingerprint{0x13, 0x14, 0x15, 0xAB, 0xCC, 0xAC, 0xDC})
	c.Assert(d2.s.Equals(generatePointFrom([symKeyLength]byte{0x41, 0x12, 0xAC, 0xDF, 0xBD, 0xBF, 0xFE})), Equals, true)
	c.Assert(d2.sigma.c1, DeepEquals, generateScalarFrom(0x01, 0x42, 0x12, 0xAB, 0xFC))
	c.Assert(d2.sigma.r1, DeepEquals, generateScalarFrom(0x01, 0x42, 0x12, 0xAB, 0xFD))
	c.Assert(d2.sigma.c2, DeepEquals, generateScalarFrom(0x02, 0x42, 0x15, 0xAB, 0xFC))
	c.Assert(d2.sigma.r2, DeepEquals, generateScalarFrom(0x03, 0x42, 0x16, 0xAC, 0xFD))
	c.Assert(d2.sigma.c3, DeepEquals, generateScalarFrom(0x04, 0x42, 0x17, 0xAD, 0xFC))
	c.Assert(d2.sigma.r3, DeepEquals, generateScalarFrom(0x05, 0x42, 0x18, 0xAE, 0xFD))
}
