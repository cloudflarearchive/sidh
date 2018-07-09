package p751toolbox

// A point on the projective line P^1(F_{p^2}).
//
// This is used to work projectively with the curve coefficients.
type ProjectiveCurveParameters struct {
	A ExtensionFieldElement
	C ExtensionFieldElement
}

// Stores curve projective parameters equivalent to A/C. Meaning of the
// values depends on the context. When working with isogenies over
// subgroup that are powers of:
// * three then  (A:C) ~ (A+2C:A-2C)
// * four then   (A:C) ~ (A+2C:  4C)
// See Appendix A of SIKE for more details
type CurveCoefficientsEquiv struct {
	A ExtensionFieldElement
	C ExtensionFieldElement
}

// A point on the projective line P^1(F_{p^2}).
//
// This represents a point on the Kummer line of a Montgomery curve.  The
// curve is specified by a ProjectiveCurveParameters struct.
type ProjectivePoint struct {
	X ExtensionFieldElement
	Z ExtensionFieldElement
}

// A point on the projective line P^1(F_p).
//
// This represents a point on the (Kummer line) of the prime-field subgroup of
// the base curve E_0(F_p), defined by E_0 : y^2 = x^3 + x.
type ProjectivePrimeFieldPoint struct {
	X PrimeFieldElement
	Z PrimeFieldElement
}

func (params *ProjectiveCurveParameters) FromAffine(a *ExtensionFieldElement) {
	params.A = *a
	params.C.One()
}

// Computes j-invariant for a curve y2=x3+A/Cx+x with A,C in F_(p^2). Result
// is returned in jBytes buffer, encoded in little-endian format. Caller
// provided jBytes buffer has to be big enough to j-invariant value. In case
// of SIDH, buffer size must be at least size of shared secret.
// Implementation corresponds to Algorithm 9 from SIKE.
func (cparams *ProjectiveCurveParameters) Jinvariant(jBytes []byte) {
	var j, t0, t1 ExtensionFieldElement

	j.Square(&cparams.A)  // j  = A^2
	t1.Square(&cparams.C) // t1 = C^2
	t0.Add(&t1, &t1)      // t0 = t1 + t1
	t0.Sub(&j, &t0)       // t0 = j - t0
	t0.Sub(&t0, &t1)      // t0 = t0 - t1
	j.Sub(&t0, &t1)       // t0 = t0 - t1
	t1.Square(&t1)        // t1 = t1^2
	j.Mul(&j, &t1)        // t0 = t0 * t1
	t0.Add(&t0, &t0)      // t0 = t0 + t0
	t0.Add(&t0, &t0)      // t0 = t0 + t0
	t1.Square(&t0)        // t1 = t0^2
	t0.Mul(&t0, &t1)      // t0 = t0 * t1
	t0.Add(&t0, &t0)      // t0 = t0 + t0
	t0.Add(&t0, &t0)      // t0 = t0 + t0
	j.Inv(&j)             // j  = 1/j
	j.Mul(&t0, &j)        // j  = t0 * j

	j.ToBytes(jBytes)
}

// Given affine points x(P), x(Q) and x(Q-P) in a extension field F_{p^2}, function
// recorvers projective coordinate A of a curve. This is Algorithm 10 from SIKE.
func (curve *ProjectiveCurveParameters) RecoverCoordinateA(xp, xq, xr *ExtensionFieldElement) {
	var t0, t1 ExtensionFieldElement

	t1.Add(xp, xq)                            // t1 = Xp + Xq
	t0.Mul(xp, xq)                            // t0 = Xp * Xq
	curve.A.Mul(xr, &t1)                      // A  = X(q-p) * t1
	curve.A.Add(&curve.A, &t0)                // A  = A + t0
	t0.Mul(&t0, xr)                           // t0 = t0 * X(q-p)
	curve.A.Sub(&curve.A, &oneExtensionField) // A  = A - 1
	t0.Add(&t0, &t0)                          // t0 = t0 + t0
	t1.Add(&t1, xr)                           // t1 = t1 + X(q-p)
	t0.Add(&t0, &t0)                          // t0 = t0 + t0
	curve.A.Square(&curve.A)                  // A  = A^2
	t0.Inv(&t0)                               // t0 = 1/t0
	curve.A.Mul(&curve.A, &t0)                // A  = A * t0
	curve.A.Sub(&curve.A, &t1)                // A  = A - t1
}

// Computes equivalence (A:C) ~ (A+2C : A-2C)
func (curve *ProjectiveCurveParameters) CalcCurveParamsEquiv3() CurveCoefficientsEquiv {
	var coef CurveCoefficientsEquiv
	var c2 ExtensionFieldElement

	c2.Add(&curve.C, &curve.C)
	// A24p = A+2*C
	coef.A.Add(&curve.A, &c2)
	// A24m = A-2*C
	coef.C.Sub(&curve.A, &c2)
	return coef
}

// Computes equivalence (A:C) ~ (A+2C : 4C)
func (cparams *ProjectiveCurveParameters) CalcCurveParamsEquiv4() CurveCoefficientsEquiv {
	var coefEq CurveCoefficientsEquiv

	coefEq.C.Add(&cparams.C, &cparams.C)
	// A24p = A+2C
	coefEq.A.Add(&cparams.A, &coefEq.C)
	// C24 = 4*C
	coefEq.C.Add(&coefEq.C, &coefEq.C)
	return coefEq
}

// Helper function for RightToLeftLadder(). Returns A+2C / 4.
func (cparams *ProjectiveCurveParameters) calcAplus2Over4() (ret ExtensionFieldElement) {
	var tmp ExtensionFieldElement
	// 2C
	tmp.Add(&cparams.C, &cparams.C)
	// A+2C
	ret.Add(&cparams.A, &tmp)
	// 1/4C
	tmp.Add(&tmp, &tmp).Inv(&tmp)
	// A+2C/4C
	ret.Mul(&ret, &tmp)
	return
}

// Recovers (A:C) curve parameters from projectively equivalent (A+2C:A-2C).
func (cparams *ProjectiveCurveParameters) RecoverCurveCoefficients3(coefEq *CurveCoefficientsEquiv) {
	cparams.A.Add(&coefEq.A, &coefEq.C)
	// cparams.A = 2*(A+2C+A-2C) = 4A
	cparams.A.Add(&cparams.A, &cparams.A)
	// cparams.C = (A+2C-A+2C) = 4C
	cparams.C.Sub(&coefEq.A, &coefEq.C)
	return
}

// Recovers (A:C) curve parameters from projectively equivalent (A+2C:4C).
func (cparams *ProjectiveCurveParameters) RecoverCurveCoefficients4(coefEq *CurveCoefficientsEquiv) {
	var half = ExtensionFieldElement{
		A: Fp751Element{
			0x00000000000124D6, 0x0000000000000000, 0x0000000000000000,
			0x0000000000000000, 0x0000000000000000, 0xB8E0000000000000,
			0x9C8A2434C0AA7287, 0xA206996CA9A378A3, 0x6876280D41A41B52,
			0xE903B49F175CE04F, 0x0F8511860666D227, 0x00004EA07CFF6E7F},
		B: Fp751Element{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	// cparams.C = (4C)*1/2=2C
	cparams.C.Mul(&coefEq.C, &half)
	// cparams.A = A+2C - 2C = A
	cparams.A.Sub(&coefEq.A, &cparams.C)
	// cparams.C = 2C * 1/2 = C
	cparams.C.Mul(&cparams.C, &half)
	return
}

func (point *ProjectivePoint) FromAffinePrimeField(x *PrimeFieldElement) {
	point.X.A = x.A
	point.X.B = zeroExtensionField.B
	point.Z = oneExtensionField
}

func (point *ProjectivePoint) FromAffine(x *ExtensionFieldElement) {
	point.X = *x
	point.Z = oneExtensionField
}

func (point *ProjectivePrimeFieldPoint) FromAffine(x *PrimeFieldElement) {
	point.X = *x
	point.Z = onePrimeField
}

func (point *ProjectivePoint) ToAffine() *ExtensionFieldElement {
	affine_x := new(ExtensionFieldElement)
	affine_x.Inv(&point.Z).Mul(affine_x, &point.X)
	return affine_x
}

func (point *ProjectivePrimeFieldPoint) ToAffine() *PrimeFieldElement {
	affine_x := new(PrimeFieldElement)
	affine_x.Inv(&point.Z).Mul(affine_x, &point.X)
	return affine_x
}

func (lhs *ProjectivePoint) VartimeEq(rhs *ProjectivePoint) bool {
	var t0, t1 ExtensionFieldElement
	t0.Mul(&lhs.X, &rhs.Z)
	t1.Mul(&lhs.Z, &rhs.X)
	return t0.VartimeEq(&t1)
}

func (lhs *ProjectivePrimeFieldPoint) VartimeEq(rhs *ProjectivePrimeFieldPoint) bool {
	var t0, t1 PrimeFieldElement
	t0.Mul(&lhs.X, &rhs.Z)
	t1.Mul(&lhs.Z, &rhs.X)
	return t0.VartimeEq(&t1)
}

func ProjectivePointConditionalSwap(xP, xQ *ProjectivePoint, choice uint8) {
	ExtensionFieldConditionalSwap(&xP.X, &xQ.X, choice)
	ExtensionFieldConditionalSwap(&xP.Z, &xQ.Z, choice)
}

func ProjectivePrimeFieldPointConditionalSwap(xP, xQ *ProjectivePrimeFieldPoint, choice uint8) {
	PrimeFieldConditionalSwap(&xP.X, &xQ.X, choice)
	PrimeFieldConditionalSwap(&xP.Z, &xQ.Z, choice)
}

// Combined coordinate doubling and differential addition. Takes projective points
// P,Q,Q-P and (A+2C)/4C curve E coefficient. Returns 2*P and P+Q calculated on E.
// Function is used only by RightToLeftLadder. Corresponds to Algorithm 5 of SIKE
func xDblAdd(P, Q, QmP *ProjectivePoint, a24 *ExtensionFieldElement) (dblP, PaQ ProjectivePoint) {
	var t0, t1, t2 ExtensionFieldElement
	xQmP, zQmP := &QmP.X, &QmP.Z
	xPaQ, zPaQ := &PaQ.X, &PaQ.Z
	x2P, z2P := &dblP.X, &dblP.Z
	xP, zP := &P.X, &P.Z
	xQ, zQ := &Q.X, &Q.Z

	t0.Add(xP, zP)       // t0   = Xp+Zp
	t1.Sub(xP, zP)       // t1   = Xp-Zp
	x2P.Square(&t0)      // 2P.X = t0^2
	t2.Sub(xQ, zQ)       // t2   = Xq-Zq
	xPaQ.Add(xQ, zQ)     // Xp+q = Xq+Zq
	t0.Mul(&t0, &t2)     // t0   = t0 * t2
	z2P.Mul(&t1, &t1)    // 2P.Z = t1 * t1
	t1.Mul(&t1, xPaQ)    // t1   = t1 * Xp+q
	t2.Sub(x2P, z2P)     // t2   = 2P.X - 2P.Z
	x2P.Mul(x2P, z2P)    // 2P.X = 2P.X * 2P.Z
	xPaQ.Mul(a24, &t2)   // Xp+q = A24 * t2
	zPaQ.Sub(&t0, &t1)   // Zp+q = t0 - t1
	z2P.Add(xPaQ, z2P)   // 2P.Z = Xp+q + 2P.Z
	xPaQ.Add(&t0, &t1)   // Xp+q = t0 + t1
	z2P.Mul(z2P, &t2)    // 2P.Z = 2P.Z * t2
	zPaQ.Square(zPaQ)    // Zp+q = Zp+q ^ 2
	xPaQ.Square(xPaQ)    // Xp+q = Xp+q ^ 2
	zPaQ.Mul(xQmP, zPaQ) // Zp+q = Xq-p * Zp+q
	xPaQ.Mul(zQmP, xPaQ) // Xp+q = Zq-p * Xp+q
	return
}

// Given the curve parameters, xP = x(P), and k >= 0, compute x2P = x([2^k]P).
//
// Returns x2P to allow chaining.  Safe to overlap xP, x2P.
func (x2P *ProjectivePoint) Pow2k(params *CurveCoefficientsEquiv, xP *ProjectivePoint, k uint32) *ProjectivePoint {
	var t0, t1 ExtensionFieldElement

	*x2P = *xP
	x, z := &x2P.X, &x2P.Z

	for i := uint32(0); i < k; i++ {
		t0.Sub(x, z)           // t0  = Xp - Zp
		t1.Add(x, z)           // t1  = Xp + Zp
		t0.Square(&t0)         // t0  = t0 ^ 2
		t1.Square(&t1)         // t1  = t1 ^ 2
		z.Mul(&params.C, &t0)  // Z2p = C24 * t0
		x.Mul(z, &t1)          // X2p = Z2p * t1
		t1.Sub(&t1, &t0)       // t1  = t1 - t0
		t0.Mul(&params.A, &t1) // t0  = A24+ * t1
		z.Add(z, &t0)          // Z2p = Z2p + t0
		z.Mul(z, &t1)          // Zp  = Z2p * t1
	}

	return x2P
}

// Given the curve parameters, xP = x(P), and k >= 0, compute x3P = x([3^k]P).
//
// Returns x3P to allow chaining.  Safe to overlap xP, xR.
func (x3P *ProjectivePoint) Pow3k(params *CurveCoefficientsEquiv, xP *ProjectivePoint, k uint32) *ProjectivePoint {
	var t0, t1, t2, t3, t4, t5, t6 ExtensionFieldElement

	*x3P = *xP
	x, z := &x3P.X, &x3P.Z

	for i := uint32(0); i < k; i++ {
		t0.Sub(x, z)           // t0  = Xp - Zp
		t2.Square(&t0)         // t2  = t0^2
		t1.Add(x, z)           // t1  = Xp + Zp
		t3.Square(&t1)         // t3  = t1^2
		t4.Add(&t1, &t0)       // t4  = t1 + t0
		t0.Sub(&t1, &t0)       // t0  = t1 - t0
		t1.Square(&t4)         // t1  = t4^2
		t1.Sub(&t1, &t3)       // t1  = t1 - t3
		t1.Sub(&t1, &t2)       // t1  = t1 - t2
		t5.Mul(&t3, &params.A) // t5  = t3 * A24+
		t3.Mul(&t3, &t5)       // t3  = t5 * t3
		t6.Mul(&t2, &params.C) // t6  = t2 * A24-
		t2.Mul(&t2, &t6)       // t2  = t2 * t6
		t3.Sub(&t2, &t3)       // t3  = t2 - t3
		t2.Sub(&t5, &t6)       // t2  = t5 - t6
		t1.Mul(&t2, &t1)       // t1  = t2 * t1
		t2.Add(&t3, &t1)       // t2  = t3 + t1
		t2.Square(&t2)         // t2  = t2^2
		x.Mul(&t2, &t4)        // X3p = t2 * t4
		t1.Sub(&t3, &t1)       // t1  = t3 - t1
		t1.Square(&t1)         // t1  = t1^2
		z.Mul(&t1, &t0)        // Z3p = t1 * t0
	}
	return x3P
}

// RightToLeftLadder is a right-to-left point multiplication that given the
// x-coordinate of P, Q and P-Q calculates the x-coordinate of R=Q+[scalar]P.
// nbits must be smaller or equal to len(scalar).
func RightToLeftLadder(c *ProjectiveCurveParameters, P, Q, PmQ *ProjectivePoint,
	nbits uint, scalar []uint8) ProjectivePoint {
	var R0, R2, R1 ProjectivePoint

	aPlus2Over4 := c.calcAplus2Over4()
	R1 = *P
	R2 = *PmQ
	R0 = *Q

	// Iterate over the bits of the scalar, bottom to top
	prevBit := uint8(0)
	for i := uint(0); i < nbits; i++ {
		bit := (scalar[i>>3] >> (i & 7) & 1)
		swap := prevBit ^ bit
		prevBit = bit
		ProjectivePointConditionalSwap(&R1, &R2, swap)
		R0, R2 = xDblAdd(&R0, &R2, &R1, &aPlus2Over4)
	}

	ProjectivePointConditionalSwap(&R1, &R2, prevBit)
	return R1
}
