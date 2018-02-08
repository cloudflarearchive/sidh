package sike

import (
	"bytes"
	"crypto/rand"
	mathRand "math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

import . "github.com/cloudflare/p751sidh/sidh"

func SIKEGenerateKeypair(quickCheckRand *mathRand.Rand, size int) reflect.Value {
	// use crypto/rand instead of the quickCheck-provided RNG
	_, sikeSecretKey, err := GenerateKeyPair(rand.Reader)
	if err != nil {
		panic("error generating secret key")
	}
	return reflect.ValueOf(*sikeSecretKey)
}

func TestSIKESharedSecret(t *testing.T) {
	sharedSecretsMatch := func() bool {

		sikePublic, sikeSecret, err := GenerateKeyPair(rand.Reader)
		if err != nil {
			panic("error generating key pair")
		}

		cipherText, sharedSecret1, err := Encapsulation(rand.Reader, sikePublic)
		if err != nil {
			panic("error generating key encapsulation")
		}

		sharedSecret2 := Decapsulation(sikeSecret, cipherText)

		return bytes.Equal(sharedSecret1.Scalar[:], sharedSecret2.Scalar[:])
	}

	if err := quick.Check(sharedSecretsMatch, nil); err != nil {
		t.Error(err)
	}
}

func BenchmarkSIKEKeypair(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateKeyPair(rand.Reader)
	}
}

var publicKeyBob, secretKeyBob, err = GenerateBobKeypair(rand.Reader)

var benchSIKEKeyEncapPublicKey = SIKEPublicKey{publicKeyBob}

func BenchmarkSIKEKeyEncap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Encapsulation(rand.Reader, &benchSIKEKeyEncapPublicKey)
	}
}

var msg = [32]uint8{31, 9, 39, 165, 125, 79, 135, 70, 97, 87, 231, 221, 204, 245, 38, 150, 198, 187, 184, 199, 148, 156, 18, 137, 71, 248, 83, 111, 170, 138, 61, 122}

var publicKeyAlice, secretKeyAlice, _ = GenerateAliceKeypair(rand.Reader)

var benchSIKEKeyDecapSecretKey = SIKESecretKey{Scalar: msg, SecretKey: secretKeyBob, PublicKey: publicKeyBob}

var benchSIKEKeyDecapCipherText = SIKECipherText{PublicKey: publicKeyAlice, Scalar: msg}

func BenchmarkSIKEKeyDecap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Decapsulation(&benchSIKEKeyDecapSecretKey, &benchSIKEKeyDecapCipherText)
	}
}
