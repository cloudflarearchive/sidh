package sidh

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"

	. "github.com/cloudflare/p751sidh/p751toolbox"
)

func newp751SIDH() SIDH {
	return SIDH{
		secretKeySize:        SecretKeySize,
		publicKeySize:        PublicKeySize,
		sharedSecretSize:     SharedSecretSize,
		maxAlice:             MaxAlice,
		maxBob:               MaxBob,
		maskAliceByte1:       MaskAliceByte1,
		maskAliceByte2:       MaskAliceByte2,
		maskAliceByte3:       MaskAliceByte3,
		maskBobByte:          MaskBobByte,
		sampleRate:           SampleRate,
		aliceIsogenyStrategy: AliceIsogenyStrategy[:],
		bobIsogenyStrategy:   BobIsogenyStrategy[:],
	}
}

func TestScheme(t *testing.T) {
	testSIDH(newp751SIDH(), t)
}

func testSIDH(sidhp751 SIDH, t testing.TB) {
	var aliceSecret = new(SIDHSecretKeyAlice)
	var bobSecret = new(SIDHSecretKeyBob)
	var alicePublic = new(SIDHPublicKeyAlice)
	var bobPublic = new(SIDHPublicKeyBob)
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
	if !bytes.Equal(secret1[:], secret2[:]) {
		t.Fatalf("The two shared keys: %d, %d do not match", secret1, secret2)
	}
}

func BenchmarkSIDH(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testSIDH(newp751SIDH(), b)
	}
}
