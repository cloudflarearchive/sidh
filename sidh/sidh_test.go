package sidh

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestP751(t *testing.T) {
	testSIDH(NewP751SIDH(), t)
}

func BenchmarkP751SIDH(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testSIDH(NewP751SIDH(), b)
	}
}

func testSIDH(sidhp751 SIDH, t testing.TB) {
	var aliceSecret SecretKeyAlice
	var bobSecret SecretKeyBob
	var alicePublic PublicKeyAlice
	var bobPublic PublicKeyBob
	var err error
	var secret1, secret2 []byte

	alicePublic, aliceSecret, err = sidhp751.GenerateAliceKeypair(rand.Reader)
	if err != nil {
		t.Error(err)
	}

	bobPublic, bobSecret, err = sidhp751.GenerateBobKeypair(rand.Reader)
	if err != nil {
		t.Error(err)
	}

	secret1 = sidhp751.SharedSecretAlice(aliceSecret, bobPublic)

	secret2 = sidhp751.SharedSecretBob(bobSecret, alicePublic)

	if !bytes.Equal(secret1, secret2) {
		t.Fatalf("The two shared keys: %d, %d do not match", secret1, secret2)
	}
}
