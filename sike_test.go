package p751sidh

import (
	"bytes"
	"crypto/rand"
	mathRand "math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

import . "github.com/cloudflare/p751sidh/p751toolbox"

func SIKEGenerateKeypair(quickCheckRand *mathRand.Rand, size int) reflect.Value {
	// use crypto/rand instead of the quickCheck-provided RNG
	_, SIKESecretKey, err := crypto_SIKE_keypair(rand.Reader)
	if err != nil {
		panic("error generating secret key")
	}
	return reflect.ValueOf(*SIKESecretKey)
}

func TestSIKESharedSecret(t *testing.T) {
	sharedSecretsMatch := func() bool {

		SIKEPublic, SIKESecret, err := crypto_SIKE_keypair(rand.Reader)
		if err != nil {
			panic("error generating key pair")
		}

		cipherText, sharedSecret_1, err := crypto_SIKE_encap(rand.Reader, SIKEPublic)
		if err != nil {
			panic("error generating key encapsulation")
		}

		sharedSecret_2 := crypto_SIKE_decap(SIKESecret, cipherText)

		return bytes.Equal(sharedSecret_1.Scalar[:], sharedSecret_2.Scalar[:])
	}

	if err := quick.Check(sharedSecretsMatch, nil); err != nil {
		t.Error(err)
	}
}

func BenchmarkSIKEKeypair(b *testing.B) {
	for n := 0; n < b.N; n++ {
		crypto_SIKE_keypair(rand.Reader)
	}
}

var benchSIKEKeyEncapPublicKey = SIDHPublicKeyBob{affine_xP: ExtensionFieldElement{A: Fp751Element{0x6e1b8b250595b5fb, 0x800787f5197d963b, 0x6f4a4e314162a8a4, 0xe75cba4d37c02128, 0x2212e7579817a216, 0xd8a5fdb0ab2f843c, 0x44230c9f998cfd6c, 0x311ff789b26aa292, 0x73d05c379ff53e40, 0xddd8f5a223bad56c, 0x94b611e6e931c8b5, 0x4d6b9bfe3555}, B: Fp751Element{0x1a3686cfc8381294, 0x57f089b14f639cc4, 0xdb6a1565f2f5cabe, 0x83d67e8f6a02f215, 0x1946272593815e87, 0x2d839631785ca74c, 0xf149dcb2dee2bee, 0x705acd79efe405bf, 0xae3769b67687fbed, 0xacd5e29f2c203cb0, 0xdd91f08fa3153e08, 0x5a9ad8cb7400}}, affine_xQ: ExtensionFieldElement{A: Fp751Element{0xd30ed48b8c0d0c4a, 0x949cad95959ec462, 0x188675581e9d1f2a, 0xf57ed3233d33031c, 0x564c6532f7283ce7, 0x80cbef8ee3b66ecb, 0x5c687359315f22ce, 0x1da950f8671fac50, 0x6fa6c045f513ef6, 0x25ffc65a8da12d4a, 0x8b0f4ac0f5244f23, 0xadcb0e07fd92}, B: Fp751Element{0x37a43cd933ebfec4, 0x2a2806ef28dacf84, 0xd671fe718611b71e, 0xef7d73f01a676326, 0x99db1524e5799cf2, 0x860271dfbf67ff62, 0xedc2a0a14114bcf, 0x6c7b9b14b1264e5a, 0xf52de61707dc38b4, 0xccddb13fcc691f5a, 0x80f37a1220163920, 0x6a9175b9d5a1}}, affine_xQmP: ExtensionFieldElement{A: Fp751Element{0xf08af9e695c626da, 0x7a4b4d52b54e1b38, 0x980272cd4c8b8c10, 0x1afcb6151d113176, 0xaef7dbd877c00f0c, 0xe8a5ea89078700c3, 0x520c1901aa8323fa, 0xfba049c947f3383a, 0x1c38abcab48be9af, 0x9f1212b923481ea, 0x1522da3457a7c293, 0xb746f78e3a61}, B: Fp751Element{0x48010d0b48491128, 0x6d1c5c509f99f450, 0xaa3522330e3a8a62, 0x872aaf46193b2bb2, 0xc89260a2d8508973, 0x98bbbebf5524be83, 0x35711d01d895c217, 0x5e44e09ec506ed7, 0xac653a760ef6fd58, 0x5837954e30ad688d, 0xcbd3e9a1b5661da8, 0x15547f5d091a}}}

func BenchmarkSIKEKeyEncap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		crypto_SIKE_encap(rand.Reader, &benchSIKEKeyEncapPublicKey)
	}
}

var msg = [32]uint8{31, 9, 39, 165, 125, 79, 135, 70, 97, 87, 231, 221, 204, 245, 38, 150, 198, 187, 184, 199, 148, 156, 18, 137, 71, 248, 83, 111, 170, 138, 61, 122}

var m_B = [48]uint8{246, 217, 158, 190, 100, 227, 224, 181, 171, 32, 120, 72, 92, 115, 113, 62, 103, 57, 71, 252, 166, 121, 126, 201, 55, 99, 213, 234, 243, 228, 171, 68, 9, 239, 214, 37, 255, 242, 217, 180, 25, 54, 242, 61, 101, 245, 78, 0}

var secretKeyBob = SIDHSecretKeyBob{Scalar: m_B}

var publicKeyBob = SIDHPublicKeyBob{affine_xP: ExtensionFieldElement{A: Fp751Element{0x6e1b8b250595b5fb, 0x800787f5197d963b, 0x6f4a4e314162a8a4, 0xe75cba4d37c02128, 0x2212e7579817a216, 0xd8a5fdb0ab2f843c, 0x44230c9f998cfd6c, 0x311ff789b26aa292, 0x73d05c379ff53e40, 0xddd8f5a223bad56c, 0x94b611e6e931c8b5, 0x4d6b9bfe3555}, B: Fp751Element{0x1a3686cfc8381294, 0x57f089b14f639cc4, 0xdb6a1565f2f5cabe, 0x83d67e8f6a02f215, 0x1946272593815e87, 0x2d839631785ca74c, 0xf149dcb2dee2bee, 0x705acd79efe405bf, 0xae3769b67687fbed, 0xacd5e29f2c203cb0, 0xdd91f08fa3153e08, 0x5a9ad8cb7400}}, affine_xQ: ExtensionFieldElement{A: Fp751Element{0xd30ed48b8c0d0c4a, 0x949cad95959ec462, 0x188675581e9d1f2a, 0xf57ed3233d33031c, 0x564c6532f7283ce7, 0x80cbef8ee3b66ecb, 0x5c687359315f22ce, 0x1da950f8671fac50, 0x6fa6c045f513ef6, 0x25ffc65a8da12d4a, 0x8b0f4ac0f5244f23, 0xadcb0e07fd92}, B: Fp751Element{0x37a43cd933ebfec4, 0x2a2806ef28dacf84, 0xd671fe718611b71e, 0xef7d73f01a676326, 0x99db1524e5799cf2, 0x860271dfbf67ff62, 0xedc2a0a14114bcf, 0x6c7b9b14b1264e5a, 0xf52de61707dc38b4, 0xccddb13fcc691f5a, 0x80f37a1220163920, 0x6a9175b9d5a1}}, affine_xQmP: ExtensionFieldElement{A: Fp751Element{0xf08af9e695c626da, 0x7a4b4d52b54e1b38, 0x980272cd4c8b8c10, 0x1afcb6151d113176, 0xaef7dbd877c00f0c, 0xe8a5ea89078700c3, 0x520c1901aa8323fa, 0xfba049c947f3383a, 0x1c38abcab48be9af, 0x9f1212b923481ea, 0x1522da3457a7c293, 0xb746f78e3a61}, B: Fp751Element{0x48010d0b48491128, 0x6d1c5c509f99f450, 0xaa3522330e3a8a62, 0x872aaf46193b2bb2, 0xc89260a2d8508973, 0x98bbbebf5524be83, 0x35711d01d895c217, 0x5e44e09ec506ed7, 0xac653a760ef6fd58, 0x5837954e30ad688d, 0xcbd3e9a1b5661da8, 0x15547f5d091a}}}

var benchSIKEKeyDecapSecretKey = SIKESecretKey{Scalar: msg, SecretKey: &secretKeyBob, PublicKey: &publicKeyBob}

var publicKeyAlice = SIDHPublicKeyAlice{affine_xP: ExtensionFieldElement{A: Fp751Element{0xea6b2d1e2aebb250, 0x35d0b205dc4f6386, 0xb198e93cb1830b8d, 0x3b5b456b496ddcc6, 0x5be3f0d41132c260, 0xce5f188807516a00, 0x54f3e7469ea8866d, 0x33809ef47f36286, 0x6fa45f83eabe1edb, 0x1b3391ae5d19fd86, 0x1e66daf48584af3f, 0xb430c14aaa87}, B: Fp751Element{0x97b41ebc61dcb2ad, 0x80ead31cb932f641, 0x40a940099948b642, 0x2a22fd16cdc7fe84, 0xaabf35b17579667f, 0x76c1d0139feb4032, 0x71467e1e7b1949be, 0x678ca8dadd0d6d81, 0x14445daea9064c66, 0x92d161eab4fa4691, 0x8dfbb01b6b238d36, 0x2e3718434e4e}}, affine_xQ: ExtensionFieldElement{A: Fp751Element{0xb055cf0ca1943439, 0xa9ff5de2fa6c69ed, 0x4f2761f934e5730a, 0x61a1dcaa1f94aa4b, 0xce3c8fadfd058543, 0xeac432aaa6701b8e, 0x8491d523093aea8b, 0xba273f9bd92b9b7f, 0xd8f59fd34439bb5a, 0xdc0350261c1fe600, 0x99375ab1eb151311, 0x14d175bbdbc5}, B: Fp751Element{0xffb0ef8c2111a107, 0x55ceca3825991829, 0xdbf8a1ccc075d34b, 0xb8e9187bd85d8494, 0x670aa2d5c34a03b0, 0xef9fe2ed2b064953, 0xc911f5311d645aee, 0xf4411f409e410507, 0x934a0a852d03e1a8, 0xe6274e67ae1ad544, 0x9f4bc563c69a87bc, 0x6f316019681e}}, affine_xQmP: ExtensionFieldElement{A: Fp751Element{0x6ffb44306a153779, 0xc0ffef21f2f918f3, 0x196c46d35d77f778, 0x4a73f80452edcfe6, 0x9b00836bce61c67f, 0x387879418d84219e, 0x20700cf9fc1ec5d1, 0x1dfe2356ec64155e, 0xf8b9e33038256b1c, 0xd2aaf2e14bada0f0, 0xb33b226e79a4e313, 0x6be576fad4e5}, B: Fp751Element{0x7db5dbc88e00de34, 0x75cc8cb9f8b6e11e, 0x8c8001c04ebc52ac, 0x67ef6c981a0b5a94, 0xc3654fbe73230738, 0xc6a46ee82983ceca, 0xed1aa61a27ef49f0, 0x17fe5a13b0858fe0, 0x9ae0ca945a4c6b3c, 0x234104a218ad8878, 0xa619627166104394, 0x556a01ff2e7e}}}

var benchSIKEKeyDecapCipherText = SIKECipherText{PublicKey: &publicKeyAlice, Scalar: msg}

func BenchmarkSIKEKeyDecap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		crypto_SIKE_decap(&benchSIKEKeyDecapSecretKey, &benchSIKEKeyDecapCipherText)
	}
}
