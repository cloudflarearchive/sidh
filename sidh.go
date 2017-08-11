package cln16sidh

// The x-coordinate of P_A = [3^239](11, oddsqrt(11^3 + 11)) on E_0(F_p)
var torsionPointPAx = PrimeFieldElement{a: fp751Element{0xd56fe52627914862, 0x1fad60dc96b5baea, 0x1e137d0bf07ab91, 0x404d3e9252161964, 0x3c5385e4cd09a337, 0x4476426769e4af73, 0x9790c6db989dfe33, 0xe06e1c04d2aa8b5e, 0x38c08185edea73b9, 0xaa41f678a4396ca6, 0x92b9259b2229e9a0, 0x2f9326818be0}}

// The y-coordinate of P_A = [3^239](11, oddsqrt(11^3 + 11)) on E_0(F_p)
var torsionPointPAy = PrimeFieldElement{a: fp751Element{0x332bd16fbe3d7739, 0x7e5e20ff2319e3db, 0xea856234aefbd81b, 0xe016df7d6d071283, 0x8ae42796f73cd34f, 0x6364b408a4774575, 0xa71c97f17ce99497, 0xda03cdd9aa0cbe71, 0xe52b4fda195bd56f, 0xdac41f811fce0a46, 0x9333720f0ee84a61, 0x1399f006e578}}

// The x-coordinate of P_B = [2^372](6, oddsqrt(6^3 + 6)) on E_0(F_p)
var torsionPointPBx = PrimeFieldElement{a: fp751Element{0xf1a8c9ed7b96c4ab, 0x299429da5178486e, 0xef4926f20cd5c2f4, 0x683b2e2858b4716a, 0xdda2fbcc3cac3eeb, 0xec055f9f3a600460, 0xd5a5a17a58c3848b, 0x4652d836f42eaed5, 0x2f2e71ed78b3a3b3, 0xa771c057180add1d, 0xc780a5d2d835f512, 0x114ea3b55ac1}}

// The y-coordinate of P_B = [2^372](6, oddsqrt(6^3 + 6)) on E_0(F_p)
var torsionPointPBy = PrimeFieldElement{a: fp751Element{0xd1e1471273e3736b, 0xf9301ba94da241fe, 0xe14ab3c17fef0a85, 0xb4ddd26a037e9e62, 0x66142dfb2afeb69, 0xe297cb70649d6c9e, 0x214dfc6e8b1a0912, 0x9f5ba818b01cf859, 0x87d15b4907c12828, 0xa4da70c53a880dbf, 0xac5df62a72c8f253, 0x2e26a42ec617}}

// The value of (a+2)/4 for the starting curve E_0 with a=0: this is 1/2
var aPlus2Over4_E0 = PrimeFieldElement{a: fp751Element{0x124d6, 0x0, 0x0, 0x0, 0x0, 0xb8e0000000000000, 0x9c8a2434c0aa7287, 0xa206996ca9a378a3, 0x6876280d41a41b52, 0xe903b49f175ce04f, 0xf8511860666d227, 0x4ea07cff6e7f}}

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
