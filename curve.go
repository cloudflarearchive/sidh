package cln16sidh

// A point on the projective line P^1(F_{p^2}).
//
// XXX understand and explain what's going on with this as a moduli space
type ProjectiveCurveParameters struct {
	A ExtensionFieldElement
	C ExtensionFieldElement
}

type CachedCurveParameters struct {
	Aplus2C ExtensionFieldElement
	C4      ExtensionFieldElement
}

// = 256
var const256 = ExtensionFieldElement{
	a: fp751Element{0x249ad67, 0x0, 0x0, 0x0, 0x0, 0x730000000000000, 0x738154969973da8b, 0x856657c146718c7f, 0x461860e4e363a697, 0xf9fd6510bba838cd, 0x4e1a3c3f06993c0c, 0x55abef5b75c7},
	b: fp751Element{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
}

func (curveParams *ProjectiveCurveParameters) jInvariant() *ExtensionFieldElement {
	var v0, v1, v2, v3 ExtensionFieldElement
	A := &curveParams.A
	C := &curveParams.C
	v0.Square(C)           // C^2
	v1.Square(A)           // A^2
	v2.Add(&v0, &v0)       // 2C^2
	v3.Add(&v2, &v0)       // 3C^2
	v2.Add(&v2, &v2)       // 4C^2
	v2.Sub(&v1, &v2)       // A^2 - 4C^2
	v1.Sub(&v1, &v3)       // A^2 - 3C^2
	v3.Square(&v1)         // (A^2 - 3C^2)^2
	v3.Mul(&v3, &v1)       // (A^2 - 3C^2)^3
	v0.Square(&v0)         // C^4
	v3.Mul(&v3, &const256) // 256(A^2 - 3C^2)^3
	v2.Mul(&v2, &v0)       // C^4(A^2 - 4C^2)
	v2.Inv(&v2)            // 1/C^4(A^2 - 4C^2)
	v0.Mul(&v3, &v2)       // 256(A^2 - 3C^2)^3 / C^4(A^2 - 4C^2)
	return &v0
}

// Compute cached parameters A + 2C, 4C.
func (curve *ProjectiveCurveParameters) cachedParams() CachedCurveParameters {
	var cached CachedCurveParameters
	cached.Aplus2C.Add(&curve.C, &curve.C)          // = 2*C
	cached.C4.Add(&cached.Aplus2C, &cached.Aplus2C) // = 4*C
	cached.Aplus2C.Add(&cached.Aplus2C, &curve.A)   // = 2*C + A
	return cached
}

// A point on the projective line P^1(F_{p^2}).
//
// This represents a point on the (Kummer line) of a Montgomery curve.  The
// curve is specified by a ProjectiveCurveParameters struct.
type ProjectivePoint struct {
	x ExtensionFieldElement // this is actually X, but can't be named that
	z ExtensionFieldElement // this is actually Z, but can't be named that
}

func (point *ProjectivePoint) toAffine() *ExtensionFieldElement {
	affine_x := new(ExtensionFieldElement)
	affine_x.Inv(&point.z).Mul(affine_x, &point.x)
	return affine_x
}

func (lhs *ProjectivePoint) VartimeEq(rhs *ProjectivePoint) bool {
	var t0, t1 ExtensionFieldElement
	t0.Mul(&lhs.x, &rhs.z)
	t1.Mul(&lhs.z, &rhs.x)
	return t0.VartimeEq(&t1)
}

func ProjectivePointConditionalSwap(xP, xQ *ProjectivePoint, choice uint8) {
	ExtensionFieldElementConditionalSwap(&xP.x, &xQ.x, choice)
	ExtensionFieldElementConditionalSwap(&xP.z, &xQ.z, choice)
}

// Given xP = x(P), xQ = x(Q), and xPmQ = x(P-Q), compute xR = x(P+Q).
//
// Returns xR to allow chaining.  Safe to overlap xP, xQ, xR.
func (xR *ProjectivePoint) Add(xP, xQ, xPmQ *ProjectivePoint) *ProjectivePoint {
	// Algorithm 1 of Costello-Smith.
	var v0, v1, v2, v3, v4 ExtensionFieldElement
	v0.Add(&xP.x, &xP.z)               // X_P + Z_P
	v1.Sub(&xQ.x, &xQ.z).Mul(&v1, &v0) // (X_Q - Z_Q)(X_P + Z_P)
	v0.Sub(&xP.x, &xP.z)               // X_P - Z_P
	v2.Add(&xQ.x, &xQ.z).Mul(&v2, &v0) // (X_Q + Z_Q)(X_P - Z_P)
	v3.Add(&v1, &v2).Square(&v3)       // 4(X_Q X_P - Z_Q Z_P)^2
	v4.Sub(&v1, &v2).Square(&v4)       // 4(X_Q Z_P - Z_Q X_P)^2
	v0.Mul(&xPmQ.z, &v3)               // 4X_{P-Q}(X_Q X_P - Z_Q Z_P)^2
	xR.z.Mul(&xPmQ.x, &v4)             // 4Z_{P-Q}(X_Q Z_P - Z_Q X_P)^2
	xR.x = v0
	return xR
}

// Given xP = x(P) and cached curve parameters Aplus2C = A + 2*C, C4 = 4*C, compute xQ = x([2]P).
//
// Returns xQ to allow chaining.  Safe to overlap xP, xQ.
func (xQ *ProjectivePoint) Double(xP *ProjectivePoint, curve *CachedCurveParameters) *ProjectivePoint {
	// Algorithm 2 of Costello-Smith, amended to work with projective curve coefficients.
	var v1, v2, v3, xz4 ExtensionFieldElement
	v1.Add(&xP.x, &xP.z).Square(&v1) // (X+Z)^2
	v2.Sub(&xP.x, &xP.z).Square(&v2) // (X-Z)^2
	xz4.Sub(&v1, &v2)                // 4XZ = (X+Z)^2 - (X-Z)^2
	v2.Mul(&v2, &curve.C4)           // 4C(X-Z)^2
	xQ.x.Mul(&v1, &v2)               // 4C(X+Z)^2(X-Z)^2
	v3.Mul(&xz4, &curve.Aplus2C)     // 4XZ(A + 2C)
	v3.Add(&v3, &v2)                 // 4XZ(A + 2C) + 4C(X-Z)^2
	xQ.z.Mul(&v3, &xz4)              // (4XZ(A + 2C) + 4C(X-Z)^2)4XZ
	// Now (xQ.x : xQ.z)
	//   = (4C(X+Z)^2(X-Z)^2 : (4XZ(A + 2C) + 4C(X-Z)^2)4XZ )
	//   = ((X+Z)^2(X-Z)^2 : (4XZ((A + 2C)/4C) + (X-Z)^2)4XZ )
	//   = ((X+Z)^2(X-Z)^2 : (4XZ((a + 2)/4) + (X-Z)^2)4XZ )
	return xQ
}

// Given the curve parameters, xP = x(P), and k >= 1, compute xQ = x([2^k]P).
//
// Returns xQ to allow chaining.  Safe to overlap xP, xQ.
func (xQ *ProjectivePoint) Pow2k(curve *ProjectiveCurveParameters, xP *ProjectivePoint, k uint32) *ProjectivePoint {
	if k == 0 {
		panic("Called Pow2k with k == 0")
	}

	cachedParams := curve.cachedParams()
	*xQ = *xP
	for i := uint32(0); i < k; i++ {
		xQ.Double(xQ, &cachedParams)
	}

	return xQ
}

// Given xP = x(P) and cached curve parameters Aplus2C = A + 2*C, C4 = 4*C, compute xQ = x([3]P).
//
// Returns xQ to allow chaining.  Safe to overlap xP, xQ.
func (xQ *ProjectivePoint) Triple(xP *ProjectivePoint, curve *CachedCurveParameters) *ProjectivePoint {
	// Uses the efficient Montgomery tripling formulas from Costello-Longa-Naehrig.
	var v0, v1, v2, v3, v4, v5 ExtensionFieldElement
	// Compute (X_2 : Z_2) = x([2]P)
	v2.Sub(&xP.x, &xP.z)           // X - Z
	v3.Add(&xP.x, &xP.z)           // X + Z
	v0.Square(&v2)                 // (X-Z)^2
	v1.Square(&v3)                 // (X+Z)^2
	v4.Mul(&v0, &curve.C4)         // 4C(X-Z)^2
	v5.Mul(&v4, &v1)               // 4C(X-Z)^2(X+Z)^2 = X_2
	v1.Sub(&v1, &v0)               // (X+Z)^2 - (X-Z)^2 = 4XZ
	v0.Mul(&v1, &curve.Aplus2C)    // 4XZ(A+2C)
	v4.Add(&v4, &v0).Mul(&v4, &v1) // (4C(X-Z)^2 + 4XZ(A+2C))4XZ = Z_2
	// Compute (X_3 : Z_3) = x(P + [2]P)
	v0.Add(&v5, &v4).Mul(&v0, &v2) // (X_2 + Z_2)(X-Z)
	v1.Sub(&v5, &v4).Mul(&v1, &v3) // (X_2 - Z_2)(X+Z)
	v4.Sub(&v0, &v1).Square(&v4)   // 4(XZ_2 - ZX_2)^2
	v5.Add(&v0, &v1).Square(&v5)   // 4(XX_2 - ZZ_2)^2
	v2.Mul(&xP.z, &v5)             // 4Z(XX_2 - ZZ_2)^2
	xQ.z.Mul(&xP.x, &v4)           // 4X(XZ_2 - ZX_2)^2
	xQ.x = v2
	return xQ
}

// Given the curve parameters, xP = x(P), and k >= 1, compute xQ = x([2^k]P).
//
// Returns xQ to allow chaining.  Safe to overlap xP, xQ.
func (xQ *ProjectivePoint) Pow3k(curve *ProjectiveCurveParameters, xP *ProjectivePoint, k uint32) *ProjectivePoint {
	if k == 0 {
		panic("Called Pow3k with k == 0")
	}

	cachedParams := curve.cachedParams()
	*xQ = *xP
	for i := uint32(0); i < k; i++ {
		xQ.Triple(xQ, &cachedParams)
	}

	return xQ
}

// Given x(P) and a scalar m in little-endian bytes, compute x([m]P) using the
// Montgomery ladder.  This is described in Algorithm 8 of Costello-Smith.
//
// This function's execution time is dependent only on the byte-length of the
// input scalar.  All scalars of the same input length execute in uniform time.
// The scalar can be padded with zero bytes to ensure a uniform length.
//
// Safe to overlap the source with the destination.
func (xQ *ProjectivePoint) ScalarMult(curve *ProjectiveCurveParameters, xP *ProjectivePoint, scalar []uint8) *ProjectivePoint {
	cachedParams := curve.cachedParams()
	var x0, x1, tmp ProjectivePoint

	x0.x.One()
	x0.z.Zero()
	x1 = *xP

	// Iterate over the bits of the scalar, top to bottom
	prevBit := uint8(0)
	for i := len(scalar) - 1; i >= 0; i-- {
		scalarByte := scalar[i]
		for j := 7; j >= 0; j-- {
			bit := (scalarByte >> uint(j)) & 0x1
			ProjectivePointConditionalSwap(&x0, &x1, (bit ^ prevBit))
			tmp.Double(&x0, &cachedParams)
			x1.Add(&x0, &x1, xP)
			x0 = tmp
			prevBit = bit
		}
	}
	// now prevBit is the lowest bit of the scalar
	ProjectivePointConditionalSwap(&x0, &x1, prevBit)
	*xQ = x0
	return xQ
}

// Given x(P), x(Q), x(P-Q), as well as a scalar m in little-endian bytes,
// compute x(P + [m]Q) using the "three-point ladder" of de Feo, Jao, and Plut.
//
// Safe to overlap the source with the destination.
//
// This function's execution time is dependent only on the byte-length of the
// input scalar.  All scalars of the same input length execute in uniform time.
// The scalar can be padded with zero bytes to ensure a uniform length.
//
// The algorithm, as described in de Feo-Jao-Plut, is as follows:
//
// (x0, x1, x2) <--- (x(O), x(Q), x(P))
//
// for i = |m| down to 0, indexing the bits of m:
//     Invariant: (x0, x1, x2) == (x( [t]Q ), x( [t+1]Q ), x( P + [t]Q ))
//          where t = m//2^i is the high bits of m, starting at i
//     if m_i == 0:
//         (x0, x1, x2) <--- (xDBL(x0), xADD(x1, x0, x(Q)), xADD(x2, x0, x(P)))
//         Invariant: (x0, x1, x2) == (x( [2t]Q ), x( [2t+1]Q ), x( P + [2t]Q ))
//                                 == (x( [t']Q ), x( [t'+1]Q ), x( P + [t']Q ))
//              where t' = m//2^{i-1} is the high bits of m, starting at i-1
//     if m_i == 1:
//         (x0, x1, x2) <--- (xADD(x1, x0, x(Q)), xDBL(x1), xADD(x2, x1, x(P-Q)))
//         Invariant: (x0, x1, x2) == (x( [2t+1]Q ), x( [2t+2]Q ), x( P + [2t+1]Q ))
//                                 == (x( [t']Q ),   x( [t'+1]Q ), x( P + [t']Q ))
//              where t' = m//2^{i-1} is the high bits of m, starting at i-1
// return x2
//
// Notice that the roles of (x0,x1) and (x(P), x(P-Q)) swap depending on the
// current bit of the scalar.  Instead of swapping which operations we do, we
// can swap variable names, producing the following uniform algorithm:
//
// (x0, x1, x2) <--- (x(O), x(Q), x(P))
// (y0, y1) <--- (x(P), x(P-Q))
//
// for i = |m| down to 0, indexing the bits of m:
//      (x0, x1) <--- SWAP( m_{i+1} xor m_i, (x0,x1) )
//      (y0, y1) <--- SWAP( m_{i+1} xor m_i, (y0,y1) )
//      (x0, x1, x2) <--- ( xDBL(x0), xADD(x1,x0,x(Q)), xADD(x2, x0, y0) )
//
// return x2
//
func (xR *ProjectivePoint) ThreePointLadder(curve *ProjectiveCurveParameters, xP, xQ, xPmQ *ProjectivePoint, scalar []uint8) *ProjectivePoint {
	cachedParams := curve.cachedParams()
	var x0, x1, x2, y0, y1, tmp ProjectivePoint

	// (x0, x1, x2) <--- (x(O), x(Q), x(P))
	x0.x.One()
	x0.z.Zero()
	x1 = *xQ
	x2 = *xP
	// (y0, y1) <--- (x(P), x(P-Q))
	y0 = *xP
	y1 = *xPmQ

	// Iterate over the bits of the scalar, top to bottom
	prevBit := uint8(0)
	for i := len(scalar) - 1; i >= 0; i-- {
		scalarByte := scalar[i]
		for j := 7; j >= 0; j-- {
			bit := (scalarByte >> uint(j)) & 0x1
			ProjectivePointConditionalSwap(&x0, &x1, (bit ^ prevBit))
			ProjectivePointConditionalSwap(&y0, &y1, (bit ^ prevBit))
			x2.Add(&x2, &x0, &y0) // = xADD(x2, x0, y0)
			tmp.Double(&x0, &cachedParams)
			x1.Add(&x1, &x0, xQ) // = xADD(x1, x0, x(Q))
			x0 = tmp             // = xDBL(x0)
			prevBit = bit
		}
	}

	*xR = x2
	return xR
}

// Represents a 3-isogeny phi, holding the data necessary to evaluate phi.
type Cached3Isogeny struct {
	x ExtensionFieldElement
	z ExtensionFieldElement
}

// Given a three-torsion point x3 = x(P_3) on the curve E_(A:C), construct the
// three-isogeny phi : E_(A:C) -> E_(A:C)/<P_3> = E_(A':C').
//
// Returns a tuple (codomain, isogeny) = (E_(A':C'), phi).
func Compute3Isogeny(x3 *ProjectivePoint) (ProjectiveCurveParameters, Cached3Isogeny) {
	var isogeny Cached3Isogeny
	isogeny.x = x3.x
	isogeny.z = x3.z
	// We want to compute
	// (A':C') = (Z^4 + 18X^2Z^2 - 27X^4 : 4XZ^3)
	// To do this, use the identity 18X^2Z^2 - 27X^4 = 9X^2(2Z^2 - 3X^2)
	var codomain ProjectiveCurveParameters
	var v0, v1, v2, v3 ExtensionFieldElement
	v1.Square(&x3.x)               // = X^2
	v0.Add(&v1, &v1).Add(&v1, &v0) // = 3X^2
	v1.Add(&v0, &v0).Add(&v1, &v0) // = 9X^2
	v2.Square(&x3.z)               // = Z^2
	v3.Square(&v2)                 // = Z^4
	v2.Add(&v2, &v2)               // = 2Z^2
	v0.Sub(&v2, &v0)               // = 2Z^2 - 3X^2
	v1.Mul(&v1, &v0)               // = 9X^2(2Z^2 - 3X^2)
	v0.Mul(&x3.x, &x3.z)           // = XZ
	v0.Add(&v0, &v0)               // = 2XZ
	codomain.A.Add(&v3, &v1)       // = Z^4 + 9X^2(2Z^2 - 3X^2)
	codomain.C.Mul(&v0, &v2)       // = 4XZ^3

	return codomain, isogeny
}

// Given a 3-isogeny phi and a point xP = x(P), compute x(Q), the x-coordinate
// of the image Q = phi(P) of P under phi : E_(A:C) -> E_(A':C').
//
// The output xQ = x(Q) is then a point on the curve E_(A':C'); the curve
// parameters are returned by the Compute3Isogeny function used to construct
// phi.
func (phi *Cached3Isogeny) Eval(xP *ProjectivePoint) ProjectivePoint {
	var xQ ProjectivePoint
	var t0, t1, t2 ExtensionFieldElement
	t0.Mul(&phi.x, &xP.x) // = X3*XP
	t1.Mul(&phi.z, &xP.z) // = Z3*XP
	t2.Sub(&t0, &t1)      // = X3*XP - Z3*ZP
	t0.Mul(&phi.z, &xP.x) // = Z3*XP
	t1.Mul(&phi.x, &xP.z) // = X3*ZP
	t0.Sub(&t0, &t1)      // = Z3*XP - X3*ZP
	t2.Square(&t2)        // = (X3*XP - Z3*ZP)^2
	t0.Square(&t0)        // = (Z3*XP - X3*ZP)^2
	xQ.x.Mul(&t2, &xP.x)  // = XP*(X3*XP - Z3*ZP)^2
	xQ.z.Mul(&t0, &xP.z)  // = ZP*(Z3*XP - X3*ZP)^2

	return xQ
}

// Represents a 4-isogeny phi, holding the data necessary to evaluate phi.
type Cached4Isogeny struct {
	Xsq_plus_Zsq  ExtensionFieldElement
	Xsq_minus_Zsq ExtensionFieldElement
	XZ2           ExtensionFieldElement
	Xpow4         ExtensionFieldElement
	Zpow4         ExtensionFieldElement
}

// Given a four-torsion point x4 = x(P_4) on the curve E_(A:C), compute the
// coefficients of the codomain E_(A':C') of the four-isogeny phi : E_(A:C) ->
// E_(A:C)/<P_4>.
//
// Returns a tuple (codomain, isogeny) = (E_(A':C') : phi).
func Compute4Isogeny(x4 *ProjectivePoint) (ProjectiveCurveParameters, Cached4Isogeny) {
	var codomain ProjectiveCurveParameters
	var isogeny Cached4Isogeny
	var v0, v1 ExtensionFieldElement
	v0.Square(&x4.x)                                     // = X4^2
	v1.Square(&x4.z)                                     // = Z4^2
	isogeny.Xsq_plus_Zsq.Add(&v0, &v1)                   // = X4^2 + Z4^2
	isogeny.Xsq_minus_Zsq.Add(&v0, &v1)                  // = X4^2 - Z4^2
	isogeny.XZ2.Add(&x4.x, &x4.z)                        // = X4 + Z4
	isogeny.XZ2.Square(&isogeny.XZ2)                     // = X4^2 + Z4^2 + 2X4Z4
	isogeny.XZ2.Sub(&isogeny.XZ2, &isogeny.Xsq_plus_Zsq) // = 2X4Z4
	isogeny.Xpow4.Square(&v0)                            // = X4^4
	isogeny.Zpow4.Square(&v1)                            // = Z4^4
	v0.Add(&isogeny.Xpow4, &isogeny.Xpow4)               // = 2X4^4
	v0.Sub(&v0, &isogeny.Zpow4)                          // = 2X4^4 - Z4^4
	codomain.A.Add(&v0, &v0)                             // = 2(2X4^4 - Z4^4)
	codomain.C = isogeny.Zpow4                           // = Z4^4

	return codomain, isogeny
}

// Given a 4-isogeny phi and a point xP = x(P), compute x(Q), the x-coordinate
// of the image Q = phi(P) of P under phi : E_(A:C) -> E_(A':C').
//
// The output xQ = x(Q) is then a point on the curve E_(A':C'); the curve
// parameters are returned by the Compute4Isogeny function used to construct
// phi.
func (phi *Cached4Isogeny) Eval(xP *ProjectivePoint) ProjectivePoint {
	var xQ ProjectivePoint
	var t0, t1, t2 ExtensionFieldElement
	// We want to compute formula (7) of Costello-Longa-Naehrig, namely
	//
	// Xprime = (2*X_4*Z*Z_4 - (X_4^2 + Z_4^2)*X)*(X*X_4 - Z*Z_4)^2*X
	// Zprime = (2*X*X_4*Z_4 - (X_4^2 + Z_4^2)*Z)*(X_4*Z - X*Z_4)^2*Z
	//
	// To do this we adapt the method in the MSR implementation, which computes
	//
	// X_Q = Xprime*( 16*(X_4 + Z_4)*(X_4 - Z_4)*X_4^2*Z_4^4 )
	// Z_Q = Zprime*( 16*(X_4 + Z_4)*(X_4 - Z_4)*X_4^2*Z_4^4 )
	//
	t0.Mul(&xP.x, &phi.XZ2)                      // = 2*X*X_4*Z_4
	t1.Mul(&xP.z, &phi.Xsq_plus_Zsq)             // = (X_4^2 + Z_4^2)*Z
	t0.Sub(&t0, &t1)                             // = -X_4^2*Z + 2*X*X_4*Z_4 - Z*Z_4^2
	t1.Mul(&xP.z, &phi.Xsq_minus_Zsq)            // = (X_4^2 - Z_4^2)*Z
	t2.Sub(&t0, &t1).Square(&t2)                 // = 4*(X_4*Z - X*Z_4)^2*X_4^2
	t0.Mul(&t0, &t1).Add(&t0, &t0).Add(&t0, &t0) // = 4*(2*X*X_4*Z_4 - (X_4^2 + Z_4^2)*Z)*(X_4^2 - Z_4^2)*Z
	t1.Add(&t0, &t2)                             // = 4*(X*X_4 - Z*Z_4)^2*Z_4^2
	t0.Mul(&t0, &t2)                             // = Zprime * 16*(X_4 + Z_4)*(X_4 - Z_4)*X_4^2
	xQ.z.Mul(&t0, &phi.Zpow4)                    // = Zprime * 16*(X_4 + Z_4)*(X_4 - Z_4)*X_4^2*Z_4^4
	t2.Mul(&t2, &phi.Zpow4)                      // = 4*(X_4*Z - X*Z_4)^2*X_4^2*Z_4^4
	t0.Mul(&t1, &phi.Xpow4)                      // = 4*(X*X_4 - Z*Z_4)^2*X_4^4*Z_4^2
	t0.Sub(&t2, &t0)                             // = -4*(X*X_4^2 - 2*X_4*Z*Z_4 + X*Z_4^2)*X*(X_4^2 - Z_4^2)*X_4^2*Z_4^2
	xQ.x.Mul(&t1, &t0)                           // = Xprime * 16*(X_4 + Z_4)*(X_4 - Z_4)*X_4^2*Z_4^4

	return xQ
}
