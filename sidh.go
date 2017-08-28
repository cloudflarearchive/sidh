package cln16sidh

import (
	"errors"
	"io"
)

const (
	PublicKeySize    = 564
	SharedSecretSize = 188
)

// The x-coordinate of P_A = [3^239](11, oddsqrt(11^3 + 11)) on E_0(F_p)
var affine_xPA = PrimeFieldElement{A: fp751Element{0xd56fe52627914862, 0x1fad60dc96b5baea, 0x1e137d0bf07ab91, 0x404d3e9252161964, 0x3c5385e4cd09a337, 0x4476426769e4af73, 0x9790c6db989dfe33, 0xe06e1c04d2aa8b5e, 0x38c08185edea73b9, 0xaa41f678a4396ca6, 0x92b9259b2229e9a0, 0x2f9326818be0}}

// The y-coordinate of P_A = [3^239](11, oddsqrt(11^3 + 11)) on E_0(F_p)
var affine_yPA = PrimeFieldElement{A: fp751Element{0x332bd16fbe3d7739, 0x7e5e20ff2319e3db, 0xea856234aefbd81b, 0xe016df7d6d071283, 0x8ae42796f73cd34f, 0x6364b408a4774575, 0xa71c97f17ce99497, 0xda03cdd9aa0cbe71, 0xe52b4fda195bd56f, 0xdac41f811fce0a46, 0x9333720f0ee84a61, 0x1399f006e578}}

// The x-coordinate of P_B = [2^372](6, oddsqrt(6^3 + 6)) on E_0(F_p)
var affine_xPB = PrimeFieldElement{A: fp751Element{0xf1a8c9ed7b96c4ab, 0x299429da5178486e, 0xef4926f20cd5c2f4, 0x683b2e2858b4716a, 0xdda2fbcc3cac3eeb, 0xec055f9f3a600460, 0xd5a5a17a58c3848b, 0x4652d836f42eaed5, 0x2f2e71ed78b3a3b3, 0xa771c057180add1d, 0xc780a5d2d835f512, 0x114ea3b55ac1}}

// The y-coordinate of P_B = [2^372](6, oddsqrt(6^3 + 6)) on E_0(F_p)
var affine_yPB = PrimeFieldElement{A: fp751Element{0xd1e1471273e3736b, 0xf9301ba94da241fe, 0xe14ab3c17fef0a85, 0xb4ddd26a037e9e62, 0x66142dfb2afeb69, 0xe297cb70649d6c9e, 0x214dfc6e8b1a0912, 0x9f5ba818b01cf859, 0x87d15b4907c12828, 0xa4da70c53a880dbf, 0xac5df62a72c8f253, 0x2e26a42ec617}}

// The value of (a+2)/4 for the starting curve E_0 with a=0: this is 1/2
var aPlus2Over4_E0 = PrimeFieldElement{A: fp751Element{0x124d6, 0x0, 0x0, 0x0, 0x0, 0xb8e0000000000000, 0x9c8a2434c0aa7287, 0xa206996ca9a378a3, 0x6876280d41a41b52, 0xe903b49f175ce04f, 0xf8511860666d227, 0x4ea07cff6e7f}}

const maxAlice = 185

var aliceIsogenyStrategy = [maxAlice]int{0, 1, 1, 2, 2, 2, 3, 4, 4, 4, 4, 5, 5,
	6, 7, 8, 8, 9, 9, 9, 9, 9, 9, 9, 12, 11, 12, 12, 13, 14, 15, 16, 16, 16, 16,
	16, 16, 17, 17, 18, 18, 17, 21, 17, 18, 21, 20, 21, 21, 21, 21, 21, 22, 25, 25,
	25, 26, 27, 28, 28, 29, 30, 31, 32, 32, 32, 32, 32, 32, 32, 33, 33, 33, 35, 36,
	36, 33, 36, 35, 36, 36, 35, 36, 36, 37, 38, 38, 39, 40, 41, 42, 38, 39, 40, 41,
	42, 40, 46, 42, 43, 46, 46, 46, 46, 48, 48, 48, 48, 49, 49, 48, 53, 54, 51, 52,
	53, 54, 55, 56, 57, 58, 59, 59, 60, 62, 62, 63, 64, 64, 64, 64, 64, 64, 64, 64,
	65, 65, 65, 65, 65, 66, 67, 65, 66, 67, 66, 69, 70, 66, 67, 66, 69, 70, 69, 70,
	70, 71, 72, 71, 72, 72, 74, 74, 75, 72, 72, 74, 74, 75, 72, 72, 74, 75, 75, 72,
	72, 74, 75, 75, 77, 77, 79, 80, 80, 82}

const maxBob = 239

var bobIsogenyStrategy = [maxBob]int{0, 1, 1, 2, 2, 2, 3, 3, 4, 4, 4, 5, 5, 5, 6,
	7, 8, 8, 8, 8, 9, 9, 9, 9, 9, 10, 12, 12, 12, 12, 12, 12, 13, 14, 14, 15, 16,
	16, 16, 16, 16, 17, 16, 16, 17, 19, 19, 20, 21, 22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 24, 24, 25, 27, 27, 28, 28, 29, 28, 29, 28, 28, 28, 30, 28, 28, 28, 29,
	30, 33, 33, 33, 33, 34, 35, 37, 37, 37, 37, 38, 38, 37, 38, 38, 38, 38, 38, 39,
	43, 38, 38, 38, 38, 43, 40, 41, 42, 43, 48, 45, 46, 47, 47, 48, 49, 49, 49, 50,
	51, 50, 49, 49, 49, 49, 51, 49, 53, 50, 51, 50, 51, 51, 51, 52, 55, 55, 55, 56,
	56, 56, 56, 56, 58, 58, 61, 61, 61, 63, 63, 63, 64, 65, 65, 65, 65, 66, 66, 65,
	65, 66, 66, 66, 66, 66, 66, 66, 71, 66, 73, 66, 66, 71, 66, 73, 66, 66, 71, 66,
	73, 68, 68, 71, 71, 73, 73, 73, 75, 75, 78, 78, 78, 80, 80, 80, 81, 81, 82, 83,
	84, 85, 86, 86, 86, 86, 86, 87, 86, 88, 86, 86, 86, 86, 88, 86, 88, 86, 86, 86,
	88, 88, 86, 86, 86, 93, 90, 90, 92, 92, 92, 93, 93, 93, 93, 93, 97, 97, 97, 97,
	97, 97}

type SIDHPublicKeyBob struct {
	affine_xP   ExtensionFieldElement
	affine_xQ   ExtensionFieldElement
	affine_xQmP ExtensionFieldElement
}

// Read a public key from a byte slice.  The input must be at least 564 bytes long.
func (pubKey *SIDHPublicKeyBob) FromBytes(input []byte) {
	if len(input) < 564 {
		panic("Too short input to SIDH pubkey FromBytes, expected 564 bytes")
	}
	pubKey.affine_xP.FromBytes(input[0:188])
	pubKey.affine_xQ.FromBytes(input[188:376])
	pubKey.affine_xQmP.FromBytes(input[376:564])
}

// Write a public key to a byte slice.  The output must be at least 564 bytes long.
func (pubKey *SIDHPublicKeyBob) ToBytes(output []byte) {
	if len(output) < 564 {
		panic("Too short output for SIDH pubkey FromBytes, expected 564 bytes")
	}
	pubKey.affine_xP.ToBytes(output[0:188])
	pubKey.affine_xQ.ToBytes(output[188:376])
	pubKey.affine_xQmP.ToBytes(output[376:564])
}

type SIDHPublicKeyAlice struct {
	affine_xP   ExtensionFieldElement
	affine_xQ   ExtensionFieldElement
	affine_xQmP ExtensionFieldElement
}

// Read a public key from a byte slice.  The input must be at least 564 bytes long.
func (pubKey *SIDHPublicKeyAlice) FromBytes(input []byte) {
	if len(input) < 564 {
		panic("Too short input to SIDH pubkey FromBytes, expected 564 bytes")
	}
	pubKey.affine_xP.FromBytes(input[0:188])
	pubKey.affine_xQ.FromBytes(input[188:376])
	pubKey.affine_xQmP.FromBytes(input[376:564])
}

// Write a public key to a byte slice.  The output must be at least 564 bytes long.
func (pubKey *SIDHPublicKeyAlice) ToBytes(output []byte) {
	if len(output) < 564 {
		panic("Too short output for SIDH pubkey FromBytes, expected 564 bytes")
	}
	pubKey.affine_xP.ToBytes(output[0:188])
	pubKey.affine_xQ.ToBytes(output[188:376])
	pubKey.affine_xQmP.ToBytes(output[376:564])
}

type SIDHSecretKeyBob struct {
	scalar []uint8
}

type SIDHSecretKeyAlice struct {
	scalar []uint8
}

func GenerateAliceKeypair(rand io.Reader) (publicKey *SIDHPublicKeyAlice, secretKey *SIDHSecretKeyAlice, err error) {
	publicKey = new(SIDHPublicKeyAlice)
	secretKey = new(SIDHSecretKeyAlice)

	scalar := new([47]byte)

	_, err = io.ReadFull(rand, scalar[:])
	if err != nil {
		return nil, nil, err
	}

	// Bit-twiddle to ensure scalar is in 2*[0,2^371):
	scalar[46] &= 15 // clear high bits, so scalar < 2^372
	scalar[0] &= 254 // clear low bit, so scalar is even

	// We actually want scalar in 2*(0,2^371), but the above procedure
	// generates 0 with probability 2^(-371), which isn't worth checking
	// for.

	secretKey.scalar = scalar[:]

	*publicKey = secretKey.PublicKey()

	return
}

// Set result to zero if the input scalar is <= 3^238.
//go:noescape
func checkLessThanThree238(scalar *[48]byte, result *uint32)

// Set scalar = 3*scalar
//go:noescape
func multiplyByThree(scalar *[48]byte)

func GenerateBobKeypair(rand io.Reader) (publicKey *SIDHPublicKeyBob, secretKey *SIDHSecretKeyBob, err error) {
	publicKey = new(SIDHPublicKeyBob)
	secretKey = new(SIDHSecretKeyBob)

	scalar := new([48]byte)

	// Perform rejection sampling to obtain a random value in [0,3^238]:
	var ok uint32
	for i := 0; i < 102; i++ {
		_, err = io.ReadFull(rand, scalar[:])
		if err != nil {
			return nil, nil, err
		}
		// Mask the high bits to obtain a uniform value in [0,2^378):
		scalar[47] &= 3
		// Accept if scalar < 3^238 (this happens w/ prob ~0.5828)
		checkLessThanThree238(scalar, &ok)
		if ok == 0 {
			break
		}
	}
	// ok is nonzero if all 102 trials failed.
	// This happens with probability 0.41719...^102 < 2^(-128), i.e., never
	if ok != 0 {
		return nil, nil, errors.New("WOW! An event with probability < 2^(-128) occurred!!")
	}

	// Multiply by 3 to get a scalar in 3*[0,3^238):
	multiplyByThree(scalar)
	// We actually want scalar in 2*(0,2^371), but the above procedure
	// generates 0 with probability 3^(-238), which isn't worth checking
	// for.

	secretKey.scalar = scalar[:]

	*publicKey = secretKey.PublicKey()

	return
}

// Compute the corresponding public key for the given secret key, using the
// fast isogeny-tree strategy.
func (secretKey *SIDHSecretKeyAlice) PublicKey() SIDHPublicKeyAlice {
	var xP, xQ, xQmP, xR ProjectivePoint

	xP.fromAffinePrimeField(&affine_xPB)     // = ( x_P : 1) = x(P_B)
	xQ.fromAffinePrimeField(&affine_xPB)     //
	xQ.X.Neg(&xQ.X)                          // = (-x_P : 1) = x(Q_B)
	xQmP = DistortAndDifference(&affine_xPB) // = x(Q_B - P_B)

	xR = SecretPoint(&affine_xPA, &affine_yPA, secretKey.scalar)

	var currentCurve ProjectiveCurveParameters
	// Starting curve has a = 0, so (A:C) = (0,1)
	currentCurve.A.Zero()
	currentCurve.C.One()

	var firstPhi FirstFourIsogeny
	currentCurve, firstPhi = ComputeFirstFourIsogeny(&currentCurve)

	xP = firstPhi.Eval(&xP)
	xQ = firstPhi.Eval(&xQ)
	xQmP = firstPhi.Eval(&xQmP)
	xR = firstPhi.Eval(&xR)

	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var phi FourIsogeny

	var i = 0

	for j := 1; j < 185; j++ {
		for i < 185-j {
			points = append(points, xR)
			indices = append(indices, i)
			k := int(aliceIsogenyStrategy[185-i-j])
			xR.Pow2k(&currentCurve, &xR, uint32(2*k))
			i = i + k
		}
		currentCurve, phi = ComputeFourIsogeny(&xR)

		for k := 0; k < len(points); k++ {
			points[k] = phi.Eval(&points[k])
		}

		xP = phi.Eval(&xP)
		xQ = phi.Eval(&xQ)
		xQmP = phi.Eval(&xQmP)

		// pop xR from points
		xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}

	currentCurve, phi = ComputeFourIsogeny(&xR)
	xP = phi.Eval(&xP)
	xQ = phi.Eval(&xQ)
	xQmP = phi.Eval(&xQmP)

	var invZP, invZQ, invZQmP ExtensionFieldElement
	ExtensionFieldBatch3Inv(&xP.Z, &xQ.Z, &xQmP.Z, &invZP, &invZQ, &invZQmP)

	var publicKey SIDHPublicKeyAlice
	publicKey.affine_xP.Mul(&xP.X, &invZP)
	publicKey.affine_xQ.Mul(&xQ.X, &invZQ)
	publicKey.affine_xQmP.Mul(&xQmP.X, &invZQmP)

	return publicKey
}

func (secretKey *SIDHSecretKeyBob) PublicKey() SIDHPublicKeyBob {
	var xP, xQ, xQmP, xR ProjectivePoint

	xP.fromAffinePrimeField(&affine_xPA)     // = ( x_P : 1) = x(P_A)
	xQ.fromAffinePrimeField(&affine_xPA)     //
	xQ.X.Neg(&xQ.X)                          // = (-x_P : 1) = x(Q_A)
	xQmP = DistortAndDifference(&affine_xPA) // = x(Q_B - P_B)

	xR = SecretPoint(&affine_xPB, &affine_yPB, secretKey.scalar)

	var currentCurve ProjectiveCurveParameters
	// Starting curve has a = 0, so (A:C) = (0,1)
	currentCurve.A.Zero()
	currentCurve.C.One()

	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var phi ThreeIsogeny

	var i = 0

	for j := 1; j < 239; j++ {
		for i < 239-j {
			points = append(points, xR)
			indices = append(indices, i)
			k := int(bobIsogenyStrategy[239-i-j])
			xR.Pow3k(&currentCurve, &xR, uint32(k))
			i = i + k
		}
		currentCurve, phi = ComputeThreeIsogeny(&xR)

		for k := 0; k < len(points); k++ {
			points[k] = phi.Eval(&points[k])
		}

		xP = phi.Eval(&xP)
		xQ = phi.Eval(&xQ)
		xQmP = phi.Eval(&xQmP)

		// pop xR from points
		xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}

	currentCurve, phi = ComputeThreeIsogeny(&xR)
	xP = phi.Eval(&xP)
	xQ = phi.Eval(&xQ)
	xQmP = phi.Eval(&xQmP)

	var invZP, invZQ, invZQmP ExtensionFieldElement
	ExtensionFieldBatch3Inv(&xP.Z, &xQ.Z, &xQmP.Z, &invZP, &invZQ, &invZQmP)

	var publicKey SIDHPublicKeyBob
	publicKey.affine_xP.Mul(&xP.X, &invZP)
	publicKey.affine_xQ.Mul(&xQ.X, &invZQ)
	publicKey.affine_xQmP.Mul(&xQmP.X, &invZQmP)

	return publicKey
}

func (aliceSecret *SIDHSecretKeyAlice) SharedSecret(bobPublic *SIDHPublicKeyBob) [SharedSecretSize]byte {
	var currentCurve = RecoverCurveParameters(&bobPublic.affine_xP, &bobPublic.affine_xQ, &bobPublic.affine_xQmP)

	var xR, xP, xQ, xQmP ProjectivePoint

	xP.fromAffine(&bobPublic.affine_xP)
	xQ.fromAffine(&bobPublic.affine_xQ)
	xQmP.fromAffine(&bobPublic.affine_xQmP)

	xR.ThreePointLadder(&currentCurve, &xP, &xQ, &xQmP, aliceSecret.scalar)

	var firstPhi FirstFourIsogeny
	currentCurve, firstPhi = ComputeFirstFourIsogeny(&currentCurve)
	xR = firstPhi.Eval(&xR)

	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var phi FourIsogeny

	var i = 0

	for j := 1; j < 185; j++ {
		for i < 185-j {
			points = append(points, xR)
			indices = append(indices, i)
			k := int(aliceIsogenyStrategy[185-i-j])
			xR.Pow2k(&currentCurve, &xR, uint32(2*k))
			i = i + k
		}
		currentCurve, phi = ComputeFourIsogeny(&xR)

		for k := 0; k < len(points); k++ {
			points[k] = phi.Eval(&points[k])
		}

		// pop xR from points
		xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}

	currentCurve, _ = ComputeFourIsogeny(&xR)

	var sharedSecret [SharedSecretSize]byte
	var jInv = currentCurve.JInvariant()
	jInv.ToBytes(sharedSecret[:])
	return sharedSecret
}

func (bobSecret *SIDHSecretKeyBob) SharedSecret(alicePublic *SIDHPublicKeyAlice) [SharedSecretSize]byte {
	var currentCurve = RecoverCurveParameters(&alicePublic.affine_xP, &alicePublic.affine_xQ, &alicePublic.affine_xQmP)

	var xR, xP, xQ, xQmP ProjectivePoint

	xP.fromAffine(&alicePublic.affine_xP)
	xQ.fromAffine(&alicePublic.affine_xQ)
	xQmP.fromAffine(&alicePublic.affine_xQmP)

	xR.ThreePointLadder(&currentCurve, &xP, &xQ, &xQmP, bobSecret.scalar)

	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var phi ThreeIsogeny

	var i = 0

	for j := 1; j < 239; j++ {
		for i < 239-j {
			points = append(points, xR)
			indices = append(indices, i)
			k := int(bobIsogenyStrategy[239-i-j])
			xR.Pow3k(&currentCurve, &xR, uint32(k))
			i = i + k
		}
		currentCurve, phi = ComputeThreeIsogeny(&xR)

		for k := 0; k < len(points); k++ {
			points[k] = phi.Eval(&points[k])
		}

		// pop xR from points
		xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}
	currentCurve, _ = ComputeThreeIsogeny(&xR)

	var sharedSecret [SharedSecretSize]byte
	var jInv = currentCurve.JInvariant()
	jInv.ToBytes(sharedSecret[:])
	return sharedSecret
}
