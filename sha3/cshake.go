package sha3

// This allows you to customize SHAKE with a custom string,
// see shake.go for more comments on usage and security.

import (
	"io"
	"bytes"
	"encoding/binary"
)

type CShakeHash interface {
	io.Writer
	io.Reader
	Clone() ShakeHash
	Reset()
}

// NewcShake128 creates a new custom cSHAKE128 variable-output-length CShakeHash.
// Its generic security strength is 128 bits against all attacks if at
// least 32 bytes of its output are used.
func NewcShake128(custom []byte) CShakeHash { 
	cshake := &state{rate: 168, dsbyte: 0x04} 
	cshake.initcShake(custom)
	return cshake
}

// NewcShake256 creates a new custom cSHAKE128 variable-output-length CShakeHash.
// Its generic security strength is 256 bits against all attacks if
// at least 64 bytes of its output are used.
func NewcShake256(custom []byte) CShakeHash {
	cshake := &state{rate: 136, dsbyte: 0x04}
	cshake.initcShake(custom)
	return cshake
}

// CShakeSum128 writes an arbitrary-length digest of data into hash.
// custom string should be of length smaller than uint64
func CShakeSum128(hash, data, custom []byte) {
	h := NewcShake128(custom)
	h.Write(data)
	h.Read(hash)
}

// CShakeSum256 writes an arbitrary-length digest of data into hash.
// custom string should be of length smaller than uint64
func CShakeSum256(hash, data, custom []byte) {
	h := NewcShake256(custom)
	h.Write(data)
	h.Read(hash)
}

// The initialization of cShake
func (cshake *state) initcShake(custom []byte) {
	var cShakePad []byte
	cShakePad = append(cShakePad, left_encode(uint64(cshake.rate))...)
	cShakePad = append(cShakePad, []byte{1, 0}...) // left_encode(0)
	cShakePad = append(cShakePad, left_encode(uint64(len(custom) * 8))...)
	cShakePad = append(cShakePad, custom...)	
	padding_len := cshake.rate - (len(cShakePad) % cshake.rate)
	cShakePad = append(cShakePad, bytes.Repeat([]byte{0}, padding_len)...)
	cshake.Write(cShakePad)
}

// Helper function for the initialization of cShake
func left_encode(value uint64) []byte {
	var input [9]byte
	var offset uint
	if value == 0 {
		offset = 8
	} else {
		binary.BigEndian.PutUint64(input[1:], value)
		for offset = 0; offset < 9; offset ++ {
			if input[offset] != 0 {
				break
			}
		}
	}
	input[offset - 1] = byte(9 - offset)
	return input[offset - 1:]
}
