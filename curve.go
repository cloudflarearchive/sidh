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

func (xQ *ProjectivePoint) ScalarMult(curve *ProjectiveCurveParameters, xP *ProjectivePoint, scalar []uint8) *ProjectivePoint {
	cachedParams := curve.cachedParams()
	var x0, x1, tmp ProjectivePoint
	x0.x.One()
	x0.z.Zero()
	x1 = *xP
	prevBit := uint8(0)
	// Iterate over the bits of the scalar, top to bottom
	for i := len(scalar) - 1; i >= 0; i-- {
		scalarByte := scalar[i]
		for j := 7; j >= 0; j-- {
			bit := (scalarByte >> uint(j)) & 0x1
			ProjectivePointConditionalSwap(&x0, &x1, (bit ^ prevBit))
			// could avoid use of tmp by having unified double/add
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
