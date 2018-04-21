// sikep751_interface implements all the required interface APIs for performing
// Supersingular Isogeny Key Encapsulation (SIKE) mechanism over p751 parameter set.
// It imports the sikep751 implementation and converts all the underlying data structures
// to generic interface structures.
// Author: Amir Jalali		ajalali2016@fau.edu
package sike

import (
	"io"
)

import . "github.com/cloudflare/p751sidh/sike/sikep751"

type p751SIKE struct {
	SIKE
}

func NewP751SIKE() SIKE {
	return &p751SIKE{}
}

func (e *p751SIKE) GenerateKeyPair(rand io.Reader) (publicKey PublicKey, secretKey SecretKey, err error) {
	var sikePublicKey = new(SIKEPublicKey)
	var sikeSecretKey = new(SIKESecretKey)

	sikePublicKey, sikeSecretKey, err = GenerateKeyPair(rand)
	return sikePublicKey, sikeSecretKey, err
}

func (e *p751SIKE) Encapsulation(rand io.Reader, publicKey PublicKey) (ciphertext Ciphertext, sharedSecret []byte, err error) {
	var sikePublicKey = new(SIKEPublicKey)
	var sikeCiphertext = new(SIKECipherText)
	var sikeSharedSecret = new(SIKESharedSecret)

	sikePublicKey = publicKey.(*SIKEPublicKey)
	sikeCiphertext, sikeSharedSecret, err = Encapsulation(rand, sikePublicKey)
	return sikeCiphertext, sikeSharedSecret.Scalar[:], err
}

func (e *p751SIKE) Decapsulation(secretKey SecretKey, ciphertext Ciphertext) ([]byte) {
	var sikeSecretKey = new(SIKESecretKey)
	var sikeCiphertext = new(SIKECipherText)
	var sikeSharedSecret = new(SIKESharedSecret)

	sikeSecretKey = secretKey.(*SIKESecretKey)
	sikeCiphertext = ciphertext.(*SIKECipherText)
	sikeSharedSecret = Decapsulation(sikeSecretKey, sikeCiphertext)

	return sikeSharedSecret.Scalar[:]
}