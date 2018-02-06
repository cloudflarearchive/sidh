// This file contains APIs for Supersingular Isogeny Key Encapsulation (SIKE) protocol
// which is recently submitted to NIST PQC standardization workshop.
// The underlying core functions are based on SIDH API developed by Cloudflare.
// SIKE requires NIST's approved hash functions such as sha3 and cSHAKE256 to encapsulate
// the key based on NIST specifications. 
//
//Author: Amir Jalali  				ajalali2016@fau.edu
//Date: Feb 2018

package p751sidh

import (
	"io"
	"golang.org/x/crypto/sha3"
	"bytes"
)

const (
	// The message size, in bytes.
	MessageSize = 32
	// The ciphertext size, in bytes = PublicKeySize + MessageSize
	CiphertextSize = 596
	// SIKE secret key size, in bytes = MessageSize + PublicKeySize + SecretKeySize
	SIKESecretKeySize = 644
	// Custom values for cSHAKE256
	G = "0"
	H = "1"
	P = "2"
)

type SIKESecretKey struct {
	Scalar [MessageSize]byte
	SecretKey *SIDHSecretKeyBob
	PublicKey *SIDHPublicKeyBob
}

type SIKECipherText struct {
	PublicKey *SIDHPublicKeyAlice
	Scalar [MessageSize]byte
}

type SIKESharedSecret struct {
	Scalar [SharedSecretSize]byte
}

// SIKE keypair generation generates SIKE secret-key and public-key.
// The secret-key contains a random message + secret-key + public-key.   
func GenerateKeyPair(rand io.Reader) (publicKey *SIDHPublicKeyBob, SIKEsecretKey *SIKESecretKey, err error) {
	
	publicKey = new(SIDHPublicKeyBob)
	SIKEsecretKey = new(SIKESecretKey)
	var secretKey = new(SIDHSecretKeyBob)

	// Randomly generate 32 byte s
	_, err = io.ReadFull(rand, SIKEsecretKey.Scalar[:])
	if err != nil {
		return nil, nil, err
	}

	// Generate Encryptor secret-key and public-key
	publicKey, secretKey, err = GenerateBobKeypair(rand)

	// Copy secretKey and publicKey into SIKEsecretKey
	SIKEsecretKey.SecretKey = secretKey
	SIKEsecretKey.PublicKey = publicKey

	return
}

// SIKE encapsulation gets the public-key as input and generates 
// the SIKE ciphertex and shared secret. The generated ciphertet contains 
// the public-key and a random message 
func Encapsulation(rand io.Reader, publicKey *SIDHPublicKeyBob) (cipherText *SIKECipherText, sharedSecret *SIKESharedSecret, err error) {
	cipherText = new(SIKECipherText)
	sharedSecret = new(SIKESharedSecret)
	var ephemeral_sk = new(SIDHSecretKeyAlice)
	var jinvariant [SharedSecretSize]byte
	var h_ [MessageSize]byte
	var tmp = make([]byte, (CiphertextSize + MessageSize))

	// Generate ephemeral secretKey G(m||pk) mod oA
	_, err = io.ReadFull(rand, tmp[:MessageSize])
	if err != nil {
		return nil, nil, err
	}

	// Append publicKey to message and hash it
	publicKey.ToBytes(tmp[MessageSize:])
	sha3.CShakeSum256(ephemeral_sk.Scalar[:], tmp[:CiphertextSize], []byte(G))
	
	// Perform mod oA
	ephemeral_sk.Scalar[47] = 0
	ephemeral_sk.Scalar[46] &= 15 // clear high bits, so scalar < 2^372
	ephemeral_sk.Scalar[0] &= 254 // clear low bit, so scalar is even

	// Encryption
	var tmp_pk = new(SIDHPublicKeyAlice)
	*tmp_pk = ephemeral_sk.PublicKey()
	cipherText.PublicKey =  tmp_pk
	jinvariant = ephemeral_sk.SharedSecret(publicKey)
	sha3.CShakeSum256(h_[:], jinvariant[:], []byte(P))

	for i := 0; i < MessageSize; i++ {
		cipherText.Scalar[i] = tmp[i] ^ h_[i]
	}

	// Generate shared secret: ss = H(m||ct)
	cipherText.PublicKey.ToBytes(tmp[MessageSize:])
	copy(tmp[CiphertextSize:], cipherText.Scalar[:])
	sha3.CShakeSum256(sharedSecret.Scalar[:], tmp[:], []byte(H))

	return
}

// SIKE decapsulation get the SIKE secret-key and ciphertext as inputs
// and computes the shared secret. 
func Decapsulation(sikeSecretKey *SIKESecretKey, cipherText *SIKECipherText) (sharedSecret *SIKESharedSecret) {
	sharedSecret = new(SIKESharedSecret)
	var ephemeral_sk = new(SIDHSecretKeyAlice)
	var jinvariant [SharedSecretSize]byte
	var h_ [MessageSize]byte
	var c0_ = new(SIDHPublicKeyAlice)
	var c0_bytes [PublicKeySize]byte
	var c1_bytes [PublicKeySize]byte
	var tmp = make([]byte, (CiphertextSize + MessageSize))
	
	// Decrypt
	jinvariant = sikeSecretKey.SecretKey.SharedSecret(cipherText.PublicKey)
	sha3.CShakeSum256(h_[:], jinvariant[:], []byte(P))
	for i:= 0; i < MessageSize; i++ {
		tmp[i] = cipherText.Scalar[i] ^ h_[i]
	}

	// Generate ephemeral secretKey G(m||pk) mod oA
	sikeSecretKey.PublicKey.ToBytes(tmp[MessageSize:])
	sha3.CShakeSum256(ephemeral_sk.Scalar[:], tmp[:CiphertextSize], []byte(G))
	ephemeral_sk.Scalar[47] = 0
	ephemeral_sk.Scalar[46] &= 15 // clear high bits, so scalar < 2^372
	ephemeral_sk.Scalar[0] &= 254 // clear low bit, so scalar is even

	// Generate shared secret ss = H(m||ct) or return ss = H(s||ct)
	*c0_ = ephemeral_sk.PublicKey()
	c0_.ToBytes(c0_bytes[:])
	cipherText.PublicKey.ToBytes(c1_bytes[:])
	if !bytes.Equal(c0_bytes[:], c1_bytes[:]){
		copy(tmp[:MessageSize], sikeSecretKey.Scalar[:])
	}
	cipherText.PublicKey.ToBytes(tmp[MessageSize:])
	copy(tmp[CiphertextSize:], cipherText.Scalar[:])
	sha3.CShakeSum256(sharedSecret.Scalar[:], tmp[:], []byte(H))

	return
}
