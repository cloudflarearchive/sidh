package cln16sidh

import (
	"testing"
)

func TestBobKeyGenFastVsSlow(t *testing.T) {
	// m_B = 3*randint(0,3^238)
	var m_B = [...]uint8{246, 217, 158, 190, 100, 227, 224, 181, 171, 32, 120, 72, 92, 115, 113, 62, 103, 57, 71, 252, 166, 121, 126, 201, 55, 99, 213, 234, 243, 228, 171, 68, 9, 239, 214, 37, 255, 242, 217, 180, 25, 54, 242, 61, 101, 245, 78}

	var bobSecretKey = SIDHSecretKey{scalar: m_B[:]}

	var fastPubKey = BobKeyGenFast(&torsionPointPAx, &torsionPointPBx, &torsionPointPBy, &bobSecretKey)
	var slowPubKey = BobKeyGenSlow(&torsionPointPAx, &torsionPointPBx, &torsionPointPBy, &bobSecretKey)

	if !fastPubKey.a.VartimeEq(&slowPubKey.a) {
		t.Error("Expected a = ", fastPubKey.a, "found", slowPubKey.a)
	}
	if !fastPubKey.affine_xP.VartimeEq(&slowPubKey.affine_xP) {
		t.Error("Expected affine_xP = ", fastPubKey.affine_xP, "found", slowPubKey.affine_xP)
	}
	if !fastPubKey.affine_xQ.VartimeEq(&slowPubKey.affine_xQ) {
		t.Error("Expected affine_xQ = ", fastPubKey.affine_xQ, "found", slowPubKey.affine_xQ)
	}
	if !fastPubKey.affine_xQmP.VartimeEq(&slowPubKey.affine_xQmP) {
		t.Error("Expected affine_xQmP = ", fastPubKey.affine_xQmP, "found", slowPubKey.affine_xQmP)
	}
}

func TestAliceKeyGenFastVsSlow(t *testing.T) {
	// m_A = 2*randint(0,2^371)
	var m_A = [...]uint8{248, 31, 9, 39, 165, 125, 79, 135, 70, 97, 87, 231, 221, 204, 245, 38, 150, 198, 187, 184, 199, 148, 156, 18, 137, 71, 248, 83, 111, 170, 138, 61, 112, 25, 188, 197, 132, 151, 1, 0, 207, 178, 24, 72, 171, 22, 11}

	var aliceSecretKey = SIDHSecretKey{scalar: m_A[:]}

	var fastPubKey = AliceKeyGenFast(&torsionPointPBx, &torsionPointPAx, &torsionPointPAy, &aliceSecretKey)
	var slowPubKey = AliceKeyGenSlow(&torsionPointPBx, &torsionPointPAx, &torsionPointPAy, &aliceSecretKey)

	if !fastPubKey.a.VartimeEq(&slowPubKey.a) {
		t.Error("Expected a = ", fastPubKey.a, "found", slowPubKey.a)
	}
	if !fastPubKey.affine_xP.VartimeEq(&slowPubKey.affine_xP) {
		t.Error("Expected affine_xP = ", fastPubKey.affine_xP, "found", slowPubKey.affine_xP)
	}
	if !fastPubKey.affine_xQ.VartimeEq(&slowPubKey.affine_xQ) {
		t.Error("Expected affine_xQ = ", fastPubKey.affine_xQ, "found", slowPubKey.affine_xQ)
	}
	if !fastPubKey.affine_xQmP.VartimeEq(&slowPubKey.affine_xQmP) {
		t.Error("Expected affine_xQmP = ", fastPubKey.affine_xQmP, "found", slowPubKey.affine_xQmP)
	}
}

func TestSecretPoint(t *testing.T) {
	// m_A = 2*randint(0,2^371)
	var m_A = [...]uint8{248, 31, 9, 39, 165, 125, 79, 135, 70, 97, 87, 231, 221, 204, 245, 38, 150, 198, 187, 184, 199, 148, 156, 18, 137, 71, 248, 83, 111, 170, 138, 61, 112, 25, 188, 197, 132, 151, 1, 0, 207, 178, 24, 72, 171, 22, 11}
	// m_B = 3*randint(0,3^238)
	var m_B = [...]uint8{246, 217, 158, 190, 100, 227, 224, 181, 171, 32, 120, 72, 92, 115, 113, 62, 103, 57, 71, 252, 166, 121, 126, 201, 55, 99, 213, 234, 243, 228, 171, 68, 9, 239, 214, 37, 255, 242, 217, 180, 25, 54, 242, 61, 101, 245, 78}

	var xR_A = SecretPoint(&torsionPointPAx, &torsionPointPAy, m_A[:])
	var xR_B = SecretPoint(&torsionPointPBx, &torsionPointPBy, m_B[:])

	var sageAffine_xR_A = ExtensionFieldElement{a: fp751Element{0x29f1dff12103d089, 0x7409b9bf955e0d87, 0xe812441c1cca7288, 0xc32b8b13efba55f9, 0xc3b76a80696d83da, 0x185dd4f93a3dc373, 0xfc07c1a9115b6717, 0x39bfcdd63b5c4254, 0xc4d097d51d41efd8, 0x4f893494389b21c7, 0x373433211d3d0446, 0x53c35ccc3d22}, b: fp751Element{0x722e718f33e40815, 0x8c5fc0fdf715667, 0x850fd292bbe8c74c, 0x212938a60fcbf5d3, 0xfdb2a099d58dc6e7, 0x232f83ab63c9c205, 0x23eda62fa5543f5e, 0x49b5758855d9d04f, 0x6b455e6642ef25d1, 0x9651162537470202, 0xfeced582f2e96ff0, 0x33a9e0c0dea8}}
	var sageAffine_xR_B = ExtensionFieldElement{a: fp751Element{0xdd4e66076e8499f5, 0xe7efddc6907519da, 0xe31f9955b337108c, 0x8e558c5479ffc5e1, 0xfee963ead776bfc2, 0x33aa04c35846bf15, 0xab77d91b23617a0d, 0xbdd70948746070e2, 0x66f71291c277e942, 0x187c39db2f901fce, 0x69262987d5d32aa2, 0xe1db40057dc}, b: fp751Element{0xd1b766abcfd5c167, 0x4591059dc8a382fa, 0x1ddf9490736c223d, 0xc96db091bdf2b3dd, 0x7b8b9c3dc292f502, 0xe5b18ad85e4d3e33, 0xc3f3479b6664b931, 0xa4f17865299e21e6, 0x3f7ef5b332fa1c6e, 0x875bedb5dab06119, 0x9b5a06ea2e23b93, 0x43d48296fb26}}

	var affine_xR_A = xR_A.toAffine()
	if !sageAffine_xR_A.VartimeEq(affine_xR_A) {
		t.Error("Expected \n", sageAffine_xR_A, "\nfound\n", affine_xR_A)
	}

	var affine_xR_B = xR_B.toAffine()
	if !sageAffine_xR_B.VartimeEq(affine_xR_B) {
		t.Error("Expected \n", sageAffine_xR_B, "\nfound\n", affine_xR_B)
	}
}

var keygenBenchPubKey SIDHPublicKey

func BenchmarkBobKeyGenFast(b *testing.B) {
	// m_B = 3*randint(0,3^238)
	var m_B = [...]uint8{246, 217, 158, 190, 100, 227, 224, 181, 171, 32, 120, 72, 92, 115, 113, 62, 103, 57, 71, 252, 166, 121, 126, 201, 55, 99, 213, 234, 243, 228, 171, 68, 9, 239, 214, 37, 255, 242, 217, 180, 25, 54, 242, 61, 101, 245, 78}

	var bobSecretKey = SIDHSecretKey{scalar: m_B[:]}

	for n := 0; n < b.N; n++ {
		keygenBenchPubKey = BobKeyGenFast(&torsionPointPAx, &torsionPointPBx, &torsionPointPBy, &bobSecretKey)
	}
}

func BenchmarkBobKeyGenSlow(b *testing.B) {
	// m_B = 3*randint(0,3^238)
	var m_B = [...]uint8{246, 217, 158, 190, 100, 227, 224, 181, 171, 32, 120, 72, 92, 115, 113, 62, 103, 57, 71, 252, 166, 121, 126, 201, 55, 99, 213, 234, 243, 228, 171, 68, 9, 239, 214, 37, 255, 242, 217, 180, 25, 54, 242, 61, 101, 245, 78}

	var bobSecretKey = SIDHSecretKey{scalar: m_B[:]}


	for n := 0; n < b.N; n++ {
		keygenBenchPubKey = BobKeyGenSlow(&torsionPointPAx, &torsionPointPBx, &torsionPointPBy, &bobSecretKey)
	}
}
