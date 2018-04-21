package sike

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestP751(t *testing.T) {
	testSIKE(NewP751SIKE(), t)
}

func BenchmarkP751SIKE(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testSIKE(NewP751SIKE(), b)
	}
}

func testSIKE(sikep751 SIKE, t testing.TB) {
	var secretKey SecretKey
	var publicKey PublicKey
	var ciphertext Ciphertext
	var secret1, secret2 []byte 
	var err error

	publicKey, secretKey, err = sikep751.GenerateKeyPair(rand.Reader)
	if err != nil {
		t.Error(err)
	}

	ciphertext, secret1, err = sikep751.Encapsulation(rand.Reader, publicKey)
	if err != nil {
		t.Error(err)
	}

	secret2 = sikep751.Decapsulation(secretKey, ciphertext)

	if !bytes.Equal(secret1, secret2) {
		t.Fatalf("The two shared keys: %d, %d do not match", secret1, secret2)
	}
}