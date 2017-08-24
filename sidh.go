package cln16sidh

// The x-coordinate of P_A = [3^239](11, oddsqrt(11^3 + 11)) on E_0(F_p)
var affine_xPA = PrimeFieldElement{a: fp751Element{0xd56fe52627914862, 0x1fad60dc96b5baea, 0x1e137d0bf07ab91, 0x404d3e9252161964, 0x3c5385e4cd09a337, 0x4476426769e4af73, 0x9790c6db989dfe33, 0xe06e1c04d2aa8b5e, 0x38c08185edea73b9, 0xaa41f678a4396ca6, 0x92b9259b2229e9a0, 0x2f9326818be0}}

// The y-coordinate of P_A = [3^239](11, oddsqrt(11^3 + 11)) on E_0(F_p)
var affine_yPA = PrimeFieldElement{a: fp751Element{0x332bd16fbe3d7739, 0x7e5e20ff2319e3db, 0xea856234aefbd81b, 0xe016df7d6d071283, 0x8ae42796f73cd34f, 0x6364b408a4774575, 0xa71c97f17ce99497, 0xda03cdd9aa0cbe71, 0xe52b4fda195bd56f, 0xdac41f811fce0a46, 0x9333720f0ee84a61, 0x1399f006e578}}

// The x-coordinate of P_B = [2^372](6, oddsqrt(6^3 + 6)) on E_0(F_p)
var affine_xPB = PrimeFieldElement{a: fp751Element{0xf1a8c9ed7b96c4ab, 0x299429da5178486e, 0xef4926f20cd5c2f4, 0x683b2e2858b4716a, 0xdda2fbcc3cac3eeb, 0xec055f9f3a600460, 0xd5a5a17a58c3848b, 0x4652d836f42eaed5, 0x2f2e71ed78b3a3b3, 0xa771c057180add1d, 0xc780a5d2d835f512, 0x114ea3b55ac1}}

// The y-coordinate of P_B = [2^372](6, oddsqrt(6^3 + 6)) on E_0(F_p)
var affine_yPB = PrimeFieldElement{a: fp751Element{0xd1e1471273e3736b, 0xf9301ba94da241fe, 0xe14ab3c17fef0a85, 0xb4ddd26a037e9e62, 0x66142dfb2afeb69, 0xe297cb70649d6c9e, 0x214dfc6e8b1a0912, 0x9f5ba818b01cf859, 0x87d15b4907c12828, 0xa4da70c53a880dbf, 0xac5df62a72c8f253, 0x2e26a42ec617}}

// The value of (a+2)/4 for the starting curve E_0 with a=0: this is 1/2
var aPlus2Over4_E0 = PrimeFieldElement{a: fp751Element{0x124d6, 0x0, 0x0, 0x0, 0x0, 0xb8e0000000000000, 0x9c8a2434c0aa7287, 0xa206996ca9a378a3, 0x6876280d41a41b52, 0xe903b49f175ce04f, 0xf8511860666d227, 0x4ea07cff6e7f}}

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

type SIDHPublicKeyAlice struct {
	affine_xP   ExtensionFieldElement
	affine_xQ   ExtensionFieldElement
	affine_xQmP ExtensionFieldElement
}

type SIDHSecretKeyBob struct {
	scalar []uint8
}

type SIDHSecretKeyAlice struct {
	scalar []uint8
}

// Compute the corresponding public key for the given secret key, using the
// fast isogeny-tree strategy.
func (secretKey *SIDHSecretKeyAlice) PublicKey() SIDHPublicKeyAlice {
	var xP, xQ, xQmP, xR ProjectivePoint

	xP.fromAffinePrimeField(&affine_xPB)     // = ( x_P : 1) = x(P_B)
	xQ.fromAffinePrimeField(&affine_xPB)     //
	xQ.x.Neg(&xQ.x)                          // = (-x_P : 1) = x(Q_B)
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
	ExtensionFieldBatch3Inv(&xP.z, &xQ.z, &xQmP.z, &invZP, &invZQ, &invZQmP)

	var publicKey SIDHPublicKeyAlice
	publicKey.affine_xP.Mul(&xP.x, &invZP)
	publicKey.affine_xQ.Mul(&xQ.x, &invZQ)
	publicKey.affine_xQmP.Mul(&xQmP.x, &invZQmP)

	return publicKey
}

func (secretKey *SIDHSecretKeyBob) PublicKey() SIDHPublicKeyBob {
	var xP, xQ, xQmP, xR ProjectivePoint

	xP.fromAffinePrimeField(&affine_xPA)     // = ( x_P : 1) = x(P_A)
	xQ.fromAffinePrimeField(&affine_xPA)     //
	xQ.x.Neg(&xQ.x)                          // = (-x_P : 1) = x(Q_A)
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
	ExtensionFieldBatch3Inv(&xP.z, &xQ.z, &xQmP.z, &invZP, &invZQ, &invZQmP)

	var publicKey SIDHPublicKeyBob
	publicKey.affine_xP.Mul(&xP.x, &invZP)
	publicKey.affine_xQ.Mul(&xQ.x, &invZQ)
	publicKey.affine_xQmP.Mul(&xQmP.x, &invZQmP)

	return publicKey
}

func (aliceSecret *SIDHSecretKeyAlice) SharedSecret(bobPublic *SIDHPublicKeyBob) ExtensionFieldElement {
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

	return currentCurve.JInvariant()
}

func (bobSecret *SIDHSecretKeyBob) SharedSecret(alicePublic *SIDHPublicKeyAlice) ExtensionFieldElement {
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

	return currentCurve.JInvariant()
}

// Given the affine x-coordinate affine_xP of P, compute the x-coordinate
// x(\tau(P)-P) of \tau(P)-P.
func DistortAndDifference(affine_xP *PrimeFieldElement) ProjectivePoint {
	var xR ProjectivePoint
	var t0, t1 PrimeFieldElement
	t0.Square(affine_xP)         // = x_P^2
	t1.One().Add(&t1, &t0)       // = x_P^2 + 1
	xR.x.b = t1.a                // = 0 + (x_P^2 + 1)*i
	t0.Add(affine_xP, affine_xP) // = 2*x_P
	xR.z.a = t0.a                // = 2*x_P + 0*i

	return xR
}

// Given an affine point P = (x_P, y_P) in the prime-field subgroup of the
// starting curve E_0(F_p), together with a secret scalar m, compute x(P+[m]Q),
// where Q = \tau(P) is the image of P under the distortion map described
// below.
//
// The computation uses basically the same strategy as the
// Costello-Longa-Naehrig implementation:
//
// 1. Use the standard Montgomery ladder to compute x([m]Q), x([m+1]Q)
//
// 2. Use Okeya-Sakurai coordinate recovery to recover [m]Q from Q, x([m]Q),
// x([m+1]Q)
//
// 3. Use P and [m]Q to compute x(P + [m]Q)
//
// The distortion map \tau is defined as
//
// \tau : E_0(F_{p^2}) ---> E_0(F_{p^2})
//
// \tau : (x,y) |---> (-x, iy).
//
// The image of the distortion map is the _trace-zero_ subgroup of E_0(F_{p^2})
// defined by Tr(P) = P + \pi_p(P) = id, where \pi_p((x,y)) = (x^p, y^p) is the
// p-power Frobenius map.  To see this, take P = (x,y) \in E_0(F_{p^2}).  Then
// Tr(P) = id if and only if \pi_p(P) = -P, so that
//
// -P = (x, -y) = (x^p, y^p) = \pi_p(P);
//
// we have x^p = x if and only if x \in F_p, while y^p = -y if and only if y =
// i*y' for y' \in F_p.
//
// Thus (excepting the identity) every point in the trace-zero subgroup is of
// the form \tau((x,y)) = (-x,i*y) for (x,y) \in E_0(F_p).
//
// Since the Montgomery ladder only uses the x-coordinate, and the x-coordinate
// is always in the prime subfield, we can compute x([m]Q), x([m+1]Q) entirely
// in the prime subfield.
//
// The affine form of the relation for Okeya-Sakurai coordinate recovery is
// given on p. 13 of Costello-Smith:
//
// y_Q = ((x_P*x_Q + 1)*(x_P + x_Q + 2*a) - 2*a - x_R*(x_P - x_Q)^2)/(2*b*y_P),
//
// where R = Q + P and a,b are the Montgomery parameters.  In our setting
// (a,b)=(0,1) and our points are P=Q, Q=[m]Q, P+Q=[m+1]Q, so this becomes
//
// y_{mQ} = ((x_Q*x_{mQ} + 1)*(x_Q + x_{mQ}) - x_{m1Q}*(x_Q - x_{mQ})^2)/(2*y_Q)
//
// y_{mQ} = ((1 - x_P*x_{mQ})*(x_{mQ} - x_P) - x_{m1Q}*(x_P + x_{mQ})^2)/(2*y_P*i)
//
// y_{mQ} = i*((1 - x_P*x_{mQ})*(x_{mQ} - x_P) - x_{m1Q}*(x_P + x_{mQ})^2)/(-2*y_P)
//
// since (x_Q, y_Q) = (-x_P, y_P*i).  In projective coordinates this is
//
// Y_{mQ}' = ((Z_{mQ} - x_P*X_{mQ})*(X_{mQ} - x_P*Z_{mQ})*Z_{m1Q}
//          - X_{m1Q}*(X_{mQ} + x_P*Z_{mQ})^2)
//
// with denominator
//
// Z_{mQ}' = (-2*y_P*Z_{mQ}*Z_{m1Q})*Z_{mQ}.
//
// Setting
//
// X_{mQ}' = (-2*y_P*Z_{mQ}*Z_{m1Q})*X_{mQ}
//
// gives [m]Q = (X_{mQ}' : i*Y_{mQ}' : Z_{mQ}') with X,Y,Z all in F_p.  (Here
// the ' just denotes that we've added extra terms to the denominators during
// the computation of Y)
//
// To compute the x-coordinate x(P+[m]Q) from P and [m]Q, we use the affine
// addition formulas of section 2.2 of Costello-Smith.  We're only interested
// in the x-coordinate, giving
//
// X_R = Z_{mQ}*(i*Y_{mQ} - y_P*Z_{mQ})^2 - (x_P*Z_{mQ} + X_{mQ})*(X_{mQ} - x_P*Z_{mQ})^2
//
// Z_R = Z_{mQ}*(X_{mQ} - x_P*Z_{mQ})^2.
//
// Notice that although X_R \in F_{p^2}, we can split the computation into
// coordinates X_R = X_{R,a} + X_{R,b}*i as
//
// (i*Y_{mQ} - y_P*Z_{mQ})^2 = (y_P*Z_{mQ})^2 - Y_{mQ}^2 - 2*y_P*Z_{mQ}*Y_{mQ}*i,
//
// giving
//
// X_{R,a} = Z_{mQ}*((y_P*Z_{mQ})^2 - Y_{mQ}^2)
//         - (x_P*Z_{mQ} + X_{mQ})*(X_{mQ} - x_P*Z_{mQ})^2
//
// X_{R,b} = -2*y_P*Y_{mQ}*Z_{mQ}^2
//
// Z_R = Z_{mQ}*(X_{mQ} - x_P*Z_{mQ})^2.
//
// These formulas could probably be combined with the formulas for y-recover
// and computed more efficiently, but efficiency isn't the biggest concern
// here, since the bulk of the cost is already in the ladder.
func SecretPoint(affine_xP, affine_yP *PrimeFieldElement, scalar []uint8) ProjectivePoint {
	var xQ ProjectivePrimeFieldPoint
	xQ.fromAffine(affine_xP)
	xQ.x.Neg(&xQ.x)

	// Compute x([m]Q) = (X_{mQ} : Z_{mQ}), x([m+1]Q) = (X_{m1Q} : Z_{m1Q})
	var xmQ, xm1Q = ScalarMultPrimeField(&aPlus2Over4_E0, &xQ, scalar)

	// Now perform coordinate recovery:
	// [m]Q = (X_{mQ} : Y_{mQ}*i : Z_{mQ})
	var XmQ, YmQ, ZmQ PrimeFieldElement
	var t0, t1 PrimeFieldElement

	// Y_{mQ} = (Z_{mQ} - x_P*X_{mQ})*(X_{mQ} - x_P*Z_{mQ})*Z_{m1Q}
	//         - X_{m1Q}*(X_{mQ} + x_P*Z_{mQ})^2
	t0.Mul(affine_xP, &xmQ.x)       // = x_P*X_{mQ}
	YmQ.Sub(&xmQ.z, &t0)            // = Z_{mQ} - x_P*X_{mQ}
	t1.Mul(affine_xP, &xmQ.z)       // = x_P*Z_{mQ}
	t0.Sub(&xmQ.x, &t1)             // = X_{mQ} - x_P*Z_{mQ}
	YmQ.Mul(&YmQ, &t0)              // = (Z_{mQ} - x_P*X_{mQ})*(X_{mQ} - x_P*Z_{mQ})
	YmQ.Mul(&YmQ, &xm1Q.z)          // = (Z_{mQ} - x_P*X_{mQ})*(X_{mQ} - x_P*Z_{mQ})*Z_{m1Q}
	t1.Add(&t1, &xmQ.x).Square(&t1) // = (X_{mQ} + x_P*Z_{mQ})^2
	t1.Mul(&t1, &xm1Q.x)            // = X_{m1Q}*(X_{mQ} + x_P*Z_{mQ})^2
	YmQ.Sub(&YmQ, &t1)              // = Y_{mQ}

	// Z_{mQ} = -2*(Z_{mQ}^2 * Z_{m1Q} * y_P)
	t0.Mul(&xmQ.z, &xm1Q.z).Mul(&t0, affine_yP) // = Z_{mQ} * Z_{m1Q} * y_P
	t0.Neg(&t0)                                 // = -1*(Z_{mQ} * Z_{m1Q} * y_P)
	t0.Add(&t0, &t0)                            // = -2*(Z_{mQ} * Z_{m1Q} * y_P)
	ZmQ.Mul(&xmQ.z, &t0)                        // = -2*(Z_{mQ}^2 * Z_{m1Q} * y_P)

	// We added terms to the denominator Z_{mQ}, so multiply them to X_{mQ}
	// X_{mQ} = -2*X_{mQ}*Z_{mQ}*Z_{m1Q}*y_P
	XmQ.Mul(&xmQ.x, &t0)

	// Now compute x(P + [m]Q) = (X_Ra + i*X_Rb : Z_R)
	var XRa, XRb, ZR PrimeFieldElement

	XRb.Square(&ZmQ).Mul(&XRb, &YmQ) // = Y_{mQ} * Z_{mQ}^2
	XRb.Mul(&XRb, affine_yP)         // = Y_{mQ} * y_P * Z_{mQ}^2
	XRb.Add(&XRb, &XRb)              // = 2 * Y_{mQ} * y_P * Z_{mQ}^2
	XRb.Neg(&XRb)                    // = -2 * Y_{mQ} * y_P * Z_{mQ}^2

	t0.Mul(affine_yP, &ZmQ).Square(&t0) // = (y_P * Z_{mQ})^2
	t1.Square(&YmQ)                     // = Y_{mQ}^2
	XRa.Sub(&t0, &t1)                   // = (y_P * Z_{mQ})^2 - Y_{mQ}^2
	XRa.Mul(&XRa, &ZmQ)                 // = Z_{mQ}*((y_P * Z_{mQ})^2 - Y_{mQ}^2)
	t0.Mul(affine_xP, &ZmQ)             // = x_P * Z_{mQ}
	t1.Add(&XmQ, &t0)                   // = X_{mQ} + x_P*Z_{mQ}
	t0.Sub(&XmQ, &t0)                   // = X_{mQ} - x_P*Z_{mQ}
	t0.Square(&t0)                      // = (X_{mQ} - x_P*Z_{mQ})^2
	t1.Mul(&t1, &t0)                    // = (X_{mQ} + x_P*Z_{mQ})*(X_{mQ} - x_P*Z_{mQ})^2
	XRa.Sub(&XRa, &t1)                  // = Z_{mQ}*((y_P*Z_{mQ})^2 - Y_{mQ}^2) - (X_{mQ} + x_P*Z_{mQ})*(X_{mQ} - x_P*Z_{mQ})^2

	ZR.Mul(&ZmQ, &t0) // = Z_{mQ}*(X_{mQ} - x_P*Z_{mQ})^2

	var xR ProjectivePoint
	xR.x.a = XRa.a
	xR.x.b = XRb.a
	xR.z.a = ZR.a

	return xR
}
