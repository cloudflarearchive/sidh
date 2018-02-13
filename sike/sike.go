// This file contains APIs for Supersingular Isogeny Key Encapsulation (SIKE) protocol
// which is recently submitted to NIST PQC standardization workshop.
// The underlying core functions are based on SIDH API developed by Cloudflare.
// SIKE requires NIST's approved hash functions such as sha3 and cSHAKE256 to encapsulate
// the key based on NIST specifications. The sha3 package along with cSHAKE256 implementation
// is placed inside the p751sidh package. Note that the standard sha3 package from
// "golang.org/x/crypto/sha3" does not include cSHAKE256 functions.
//
// Author: Amir Jalali  				ajalali2016@fau.edu
// Date: Feb 2018

package sike

import (
	"bytes"
	"io"
)

import . "github.com/cloudflare/p751sidh/sha3"

import . "github.com/cloudflare/p751sidh/sidh"

import . "github.com/cloudflare/p751sidh/p751toolbox"

// SIKESecretKey contains SIDHSecretKeyBob, SIDHPublicKeyBob, and Scalar
type SIKESecretKey struct {
	Scalar    [MessageSize]byte
	SecretKey *SIDHSecretKeyBob
	PublicKey *SIDHPublicKeyBob
}

// SIKEPublicKey is just a wrapper around SIDHPublicKeyBob
type SIKEPublicKey struct {
	PublicKey *SIDHPublicKeyBob
}

// SIKECipherText contains SIDHSecretKeyBob, SIDHPublicKeyBob, and Scalar
type SIKECipherText struct {
	PublicKey *SIDHPublicKeyAlice
	Scalar    [MessageSize]byte
}

// SIKESharedSecret contains SIDHSecretKeyBob, SIDHPublicKeyBob, and Scalar
type SIKESharedSecret struct {
	Scalar [SharedSecretSize]byte
}

// GenerateKeyPair generation generates SIKE secret-key and public-key.
// The secret-key contains a random message + secret-key + public-key.
func GenerateKeyPair(rand io.Reader) (publicKey *SIKEPublicKey, sikeSecretKey *SIKESecretKey, err error) {

	publicKey = new(SIKEPublicKey)
	sikeSecretKey = new(SIKESecretKey)
	var secretKey = new(SIDHSecretKeyBob)

	// Randomly generate 32 byte s
	_, err = io.ReadFull(rand, sikeSecretKey.Scalar[:])
	if err != nil {
		return nil, nil, err
	}

	// Generate Encryptor secret-key and public-key
	publicKey.PublicKey, secretKey, err = GenerateBobKeypair(rand)

	// Copy secretKey and publicKey into SIKEsecretKey
	sikeSecretKey.SecretKey = secretKey
	sikeSecretKey.PublicKey = publicKey.PublicKey

	return
}

// Encapsulation gets the public-key as input and generates
// the SIKE ciphertex and shared secret. The generated ciphertet contains
// the public-key and a random message
func Encapsulation(rand io.Reader, publicKey *SIKEPublicKey) (cipherText *SIKECipherText, sharedSecret *SIKESharedSecret, err error) {
	cipherText = new(SIKECipherText)
	sharedSecret = new(SIKESharedSecret)
	var ephemeralSk = new(SIDHSecretKeyAlice)
	var jinvariant [SharedSecretSize]byte
	var h [MessageSize]byte
	var tmp = make([]byte, (CiphertextSize + MessageSize))

	// Generate ephemeral secretKey G(m||pk) mod oA
	_, err = io.ReadFull(rand, tmp[:MessageSize])
	if err != nil {
		return nil, nil, err
	}

	// Append publicKey to message and hash it
	publicKey.PublicKey.ToBytes(tmp[MessageSize:])
	CShakeSum256(ephemeralSk.Scalar[:], tmp[:CiphertextSize], []byte(G))

	// Perform mod oA
	ephemeralSk.Scalar[SecretKeySize-1] = MaskAliceByte1
	ephemeralSk.Scalar[SecretKeySize-2] &= MaskAliceByte2 // clear high bits, so scalar < 2^372
	ephemeralSk.Scalar[0] &= MaskAliceByte3               // clear low bit, so scalar is even

	// Encryption
	var tmpPk = new(SIDHPublicKeyAlice)
	*tmpPk = ephemeralSk.PublicKey()
	cipherText.PublicKey = tmpPk
	jinvariant = ephemeralSk.SharedSecret(publicKey.PublicKey)
	CShakeSum256(h[:], jinvariant[:], []byte(P))

	for i := 0; i < MessageSize; i++ {
		cipherText.Scalar[i] = tmp[i] ^ h[i]
	}

	// Generate shared secret: ss = H(m||ct)
	cipherText.PublicKey.ToBytes(tmp[MessageSize:])
	copy(tmp[CiphertextSize:], cipherText.Scalar[:])
	CShakeSum256(sharedSecret.Scalar[:], tmp[:], []byte(H))

	return
}

// Decapsulation gets the SIKE secret-key and ciphertext as inputs
// and computes the shared secret.
func Decapsulation(sikeSecretKey *SIKESecretKey, cipherText *SIKECipherText) (sharedSecret *SIKESharedSecret) {
	sharedSecret = new(SIKESharedSecret)
	var ephemeralSk = new(SIDHSecretKeyAlice)
	var jinvariant [SharedSecretSize]byte
	var h [MessageSize]byte
	var c0 = new(SIDHPublicKeyAlice)
	var c0Bytes [PublicKeySize]byte
	var c1Bytes [PublicKeySize]byte
	var tmp = make([]byte, (CiphertextSize + MessageSize))

	// Decrypt
	jinvariant = sikeSecretKey.SecretKey.SharedSecret(cipherText.PublicKey)
	CShakeSum256(h[:], jinvariant[:], []byte(P))
	for i := 0; i < MessageSize; i++ {
		tmp[i] = cipherText.Scalar[i] ^ h[i]
	}

	// Generate ephemeral secretKey G(m||pk) mod oA
	sikeSecretKey.PublicKey.ToBytes(tmp[MessageSize:])
	CShakeSum256(ephemeralSk.Scalar[:], tmp[:CiphertextSize], []byte(G))
	ephemeralSk.Scalar[SecretKeySize-1] = MaskAliceByte1
	ephemeralSk.Scalar[SecretKeySize-2] &= MaskAliceByte2 // clear high bits, so scalar < 2^372
	ephemeralSk.Scalar[0] &= MaskAliceByte3               // clear low bit, so scalar is even

	// Generate shared secret ss = H(m||ct) or return ss = H(s||ct)
	*c0 = ephemeralSk.PublicKey()
	c0.ToBytes(c0Bytes[:])
	cipherText.PublicKey.ToBytes(c1Bytes[:])
	if !bytes.Equal(c0Bytes[:], c1Bytes[:]) {
		copy(tmp[:MessageSize], sikeSecretKey.Scalar[:])
	}
	cipherText.PublicKey.ToBytes(tmp[MessageSize:])
	copy(tmp[CiphertextSize:], cipherText.Scalar[:])
	CShakeSum256(sharedSecret.Scalar[:], tmp[:], []byte(H))

	return
}
