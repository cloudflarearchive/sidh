/*
Abstract:
This file contains APIs for Supersingular Isogeny Key Encapsulation (SIKE) protocol which is recently submitted to NIST PQC candidates.
The underlying core functions are based on SIDH keygeneration functions developed by Cloudflare.
SIKE requires NIST's approved hash functions such as sha3 and cSHAKE256 to encapsulate the key based on NIST specifications. The custome 
SHAKE256 is adopted from  https://gist.github.com/mimoo/7e815318e54d5c07c3330149ddf439c5
and placed inside "golang.org/x/crypto/sha3" path

Author: Amir Jalali  				ajalali2016@fau.edu
Date: Feb 2018
*/

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
)

// SIKE secret key.
type SIKESecretKey struct {
	//Scalar Message
	Scalar [MessageSize]byte
	//SecretKey Bob
	SecretKey *SIDHSecretKeyBob
	//PublicKey Bob
	PublicKey *SIDHPublicKeyBob
}

type SIKECipherText struct {
	//PublicKey Alice
	PublicKey *SIDHPublicKeyAlice
	//Scalar Message
	Scalar [MessageSize]byte
}

type SIKESharedSecret struct{

	Scalar [SharedSecretSize]byte
}

func crypto_SIKE_keypair(rand io.Reader) (publicKey *SIDHPublicKeyBob, SIKEsecretKey *SIKESecretKey, err error){
	// SIKE keypair generation
	// Outputs: SIKE secret key (message + secretkey + publickey)
	//			SIKE public key (public key)

	publicKey = new(SIDHPublicKeyBob)
	SIKEsecretKey = new(SIKESecretKey)
	var secretKey = new(SIDHSecretKeyBob)

	//first randomly generate 32 byte s
	_, err = io.ReadFull(rand, SIKEsecretKey.Scalar[:])
	if err != nil {
		return nil, nil, err
	}

	//next, generate Encryptor secret-key and public-key
	publicKey, secretKey, err = GenerateBobKeypair(rand)

	//Copy secretKey and publicKey into SIKEsecretKey
	SIKEsecretKey.SecretKey = secretKey
	SIKEsecretKey.PublicKey = publicKey

	return
}

func crypto_SIKE_encap(rand io.Reader, publicKey *SIDHPublicKeyBob) (cipherText *SIKECipherText, sharedSecret *SIKESharedSecret, err error) {
	// SIKE encapsulation
	// Input: publicKey
	// Output: shared secret ss
	//		   ciphertext message ct (publicKey + message)
	cipherText = new(SIKECipherText)
	sharedSecret = new(SIKESharedSecret)

	const G = "0"
	const H = "1"
	const P = "2"
	var ephemeral_sk = new(SIDHSecretKeyAlice)
	var jinvariant [SharedSecretSize]byte
	var h_ [MessageSize]byte
	//var tmp [CiphertextSize + MessageSize]byte
	var tmp = make([]byte, (CiphertextSize + MessageSize))

	//Generate ephemeral secretKey G(m||pk) mod oA
	_, err = io.ReadFull(rand, tmp[:MessageSize])
	if err != nil {
		return nil, nil, err
	}

	//Append publicKey to message and hash it
	publicKey.ToBytes(tmp[MessageSize:])
	sha3.CShakeSum256(ephemeral_sk.Scalar[:], tmp[:CiphertextSize], []byte(G))
	//perform mod oA
	ephemeral_sk.Scalar[47] = 0
	ephemeral_sk.Scalar[46] &= 15 // clear high bits, so scalar < 2^372
	ephemeral_sk.Scalar[0] &= 254 // clear low bit, so scalar is even

	//Encryption
	var tmp_pk = new(SIDHPublicKeyAlice)
	//tmp_pk, ephemeral_sk, err = GenerateAliceKeypair(rand)

	*tmp_pk = ephemeral_sk.PublicKey()
	cipherText.PublicKey =  tmp_pk
	jinvariant = ephemeral_sk.SharedSecret(publicKey)
	sha3.CShakeSum256(h_[:], jinvariant[:], []byte(P))

	for i := 0; i < MessageSize; i++ {
		cipherText.Scalar[i] = tmp[i] ^ h_[i]
	}

	//Generate shared secret: ss = H(m||ct)
	cipherText.PublicKey.ToBytes(tmp[MessageSize:])
	copy(tmp[CiphertextSize:], cipherText.Scalar[:])
	sha3.CShakeSum256(sharedSecret.Scalar[:], tmp[:], []byte(H))

	return
}

func crypto_SIKE_decap(SIKEsecretKey *SIKESecretKey, cipherText *SIKECipherText) (sharedSecret *SIKESharedSecret){
	// SIKE decapsulation
	// Input: secretKey
	//		  cipherText message
	// Output: sharedSecret
	sharedSecret = new(SIKESharedSecret)

	const G = "0"
	const H = "1"
	const P = "2"
	var ephemeral_sk = new(SIDHSecretKeyAlice)
	var jinvariant [SharedSecretSize]byte
	var h_ [MessageSize]byte
	var c0_ = new(SIDHPublicKeyAlice)
	var c0_bytes [PublicKeySize]byte
	var c1_bytes [PublicKeySize]byte
	//var tmp [CiphertextSize+MessageSize]byte
	var tmp = make([]byte, (CiphertextSize + MessageSize))
	//Decrypt
	jinvariant = SIKEsecretKey.SecretKey.SharedSecret(cipherText.PublicKey)
	sha3.CShakeSum256(h_[:], jinvariant[:], []byte(P))
	for i:= 0; i < MessageSize; i++ {
		tmp[i] = cipherText.Scalar[i] ^ h_[i]
	}

	//Generate ephemeral secretKey G(m||pk) mod oA
	SIKEsecretKey.PublicKey.ToBytes(tmp[MessageSize:])
	sha3.CShakeSum256(ephemeral_sk.Scalar[:], tmp[:CiphertextSize], []byte(G))
	ephemeral_sk.Scalar[47] = 0
	ephemeral_sk.Scalar[46] &= 15 // clear high bits, so scalar < 2^372
	ephemeral_sk.Scalar[0] &= 254 // clear low bit, so scalar is even


	//Generate shared secret ss = H(m||ct) or return ss = H(s||ct)
	*c0_ = ephemeral_sk.PublicKey()
	c0_.ToBytes(c0_bytes[:])
	cipherText.PublicKey.ToBytes(c1_bytes[:])
	if !bytes.Equal(c0_bytes[:], c1_bytes[:]){
		copy(tmp[:MessageSize], SIKEsecretKey.Scalar[:])
	}
	cipherText.PublicKey.ToBytes(tmp[MessageSize:])
	copy(tmp[CiphertextSize:], cipherText.Scalar[:])
	sha3.CShakeSum256(sharedSecret.Scalar[:], tmp[:], []byte(H))

	return
}
