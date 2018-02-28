// sidhp751_interface implements all the required interface APIs for performing
// SIDH key-exchange over p751 parameter. It imports the sidhp751 implementation
// and converts all the underlying data structures to generic interface structures.
// Author: Amir Jalali		ajalali2016@fau.edu
package sidh

import (
	"io"
)

import . "github.com/cloudflare/p751sidh/sidh/sidhp751"

type p751SIDH struct {
	SIDH
}

func NewP751SIDH() SIDH {
	return &p751SIDH{}
}

func (e *p751SIDH) GenerateAliceKeypair(rand io.Reader) (publicKey PublicKeyAlice, secretKey SecretKeyAlice, err error) {
	var sidhPublicKey = new(SIDHPublicKeyAlice)
	var sidhSecretKey = new(SIDHSecretKeyAlice)

	sidhPublicKey, sidhSecretKey, err = GenerateAliceKeypair(rand)
	return sidhPublicKey, sidhSecretKey, err
}

func (e *p751SIDH) GenerateBobKeypair(rand io.Reader) (publicKey PublicKeyBob, secretKey SecretKeyBob, err error) {
	var sidhPublicKey = new(SIDHPublicKeyBob)
	var sidhSecretKey = new(SIDHSecretKeyBob)

	sidhPublicKey, sidhSecretKey, err = GenerateBobKeypair(rand)
	return sidhPublicKey, sidhSecretKey, err
}

func (e *p751SIDH) PublicKeyAlice(secretKey SecretKeyAlice) (publicKey PublicKeyAlice) {
	var sidhPublicKey SIDHPublicKeyAlice
	var sidhSecretKey = new(SIDHSecretKeyAlice)

	sidhSecretKey = secretKey.(*SIDHSecretKeyAlice)
	sidhPublicKey = sidhSecretKey.PublicKey()
	return &sidhPublicKey
}

func (e *p751SIDH) PublicKeyBob(secretKey SecretKeyBob) (publicKey PublicKeyBob) {
	var sidhPublicKey SIDHPublicKeyBob
	var sidhSecretKey = new(SIDHSecretKeyBob)

	sidhSecretKey = secretKey.(*SIDHSecretKeyBob)
	sidhPublicKey = sidhSecretKey.PublicKey()
	return &sidhPublicKey
}

func (e *p751SIDH) SharedSecretAlice(aliceSecret SecretKeyAlice, bobPublic PublicKeyBob) []byte {
	var sharedSecret [SharedSecretSize]byte
	var sidhSecretKey = new(SIDHSecretKeyAlice)
	var sidhPublicKey = new(SIDHPublicKeyBob)

	sidhSecretKey = aliceSecret.(*SIDHSecretKeyAlice)
	sidhPublicKey = bobPublic.(*SIDHPublicKeyBob)
	sharedSecret = sidhSecretKey.SharedSecret(sidhPublicKey)
	return sharedSecret[:]
}

func (e *p751SIDH) SharedSecretBob(bobSecret SecretKeyBob, alicePublic PublicKeyAlice) []byte {
	var sharedSecret [SharedSecretSize]byte
	var sidhSecretKey = new(SIDHSecretKeyBob)
	var sidhPublicKey = new(SIDHPublicKeyAlice)

	sidhSecretKey = bobSecret.(*SIDHSecretKeyBob)
	sidhPublicKey = alicePublic.(*SIDHPublicKeyAlice)
	sharedSecret = sidhSecretKey.SharedSecret(sidhPublicKey)
	return sharedSecret[:]
}
