package p751toolbox

// Interface for working with isogenies.
type Isogeny interface {
	// Given a torsion point on a curve computes isogenous curve.
	// Returns curve coefficients (A:C), so that E_(A/C) = E_(A/C)/<P>,
	// where P is a provided projective point. Sets also isogeny constants
	// that are needed for isogeny evaluation.
	GenerateCurve(*ProjectivePoint) CurveCoefficientsEquiv
	// Evaluates isogeny at caller provided point. Requires isogeny curve constants
	// to be earlier computed by GenerateCurve.
	EvaluatePoint(*ProjectivePoint) ProjectivePoint
}

// Stores Isogeny 4 curve constants
type isogeny4 struct {
	isogeny3
	K3 ExtensionFieldElement
}

// Stores Isogeny 3 curve constants
type isogeny3 struct {
	K1 ExtensionFieldElement
	K2 ExtensionFieldElement
}

// Constructs isogeny4 objects
func NewIsogeny4() Isogeny {
	return new(isogeny4)
}

// Constructs isogeny3 objects
func NewIsogeny3() Isogeny {
	return new(isogeny3)
}

// Given a three-torsion point p = x(PB) on the curve E_(A:C), construct the
// three-isogeny phi : E_(A:C) -> E_(A:C)/<P_3> = E_(A':C').
//
// Input: (XP_3: ZP_3), where P_3 has exact order 3 on E_A/C
// Output: * Curve coordinates (A' + 2C', A' - 2C') corresponding to E_A'/C' = A_E/C/<P3>
//		   * Isogeny phi with constants in F_p^2
func (phi *isogeny3) GenerateCurve(p *ProjectivePoint) CurveCoefficientsEquiv {
	var t0, t1, t2, t3, t4 ExtensionFieldElement
	var coefEq CurveCoefficientsEquiv
	var K1, K2 = &phi.K1, &phi.K2

	K1.Sub(&p.X, &p.Z)           // K1 = XP3 - ZP3
	t0.Square(K1)                // t0 = K1^2
	K2.Add(&p.X, &p.Z)           // K2 = XP3 + ZP3
	t1.Square(K2)                // t1 = K2^2
	t2.Add(&t0, &t1)             // t2 = t0 + t1
	t3.Add(K1, K2)               // t3 = K1 + K2
	t3.Square(&t3)               // t3 = t3^2
	t3.Sub(&t3, &t2)             // t3 = t3 - t2
	t2.Add(&t1, &t3)             // t2 = t1 + t3
	t3.Add(&t3, &t0)             // t3 = t3 + t0
	t4.Add(&t3, &t0)             // t4 = t3 + t0
	t4.Add(&t4, &t4)             // t4 = t4 + t4
	t4.Add(&t1, &t4)             // t4 = t1 + t4
	coefEq.C.Mul(&t2, &t4)       // A24m = t2 * t4
	t4.Add(&t1, &t2)             // t4 = t1 + t2
	t4.Add(&t4, &t4)             // t4 = t4 + t4
	t4.Add(&t0, &t4)             // t4 = t0 + t4
	t4.Mul(&t3, &t4)             // t4 = t3 * t4
	t0.Sub(&t4, &coefEq.C)       // t0 = t4 - A24m
	coefEq.A.Add(&coefEq.C, &t0) // A24p = A24m + t0
	return coefEq
}

// Given a 3-isogeny phi and a point pB = x(PB), compute x(QB), the x-coordinate
// of the image QB = phi(PB) of PB under phi : E_(A:C) -> E_(A':C').
//
// The output xQ = x(Q) is then a point on the curve E_(A':C'); the curve
// parameters are returned by the GenerateCurve function used to construct phi.
func (phi *isogeny3) EvaluatePoint(p *ProjectivePoint) ProjectivePoint {
	var t0, t1, t2 ExtensionFieldElement
	var q ProjectivePoint
	var K1, K2 = &phi.K1, &phi.K2
	var px, pz = &p.X, &p.Z

	t0.Add(px, pz)   // t0 = XQ + ZQ
	t1.Sub(px, pz)   // t1 = XQ - ZQ
	t0.Mul(K1, &t0)  // t2 = K1 * t0
	t1.Mul(K2, &t1)  // t1 = K2 * t1
	t2.Add(&t0, &t1) // t2 = t0 + t1
	t0.Sub(&t1, &t0) // t0 = t1 - t0
	t2.Square(&t2)   // t2 = t2 ^ 2
	t0.Square(&t0)   // t0 = t0 ^ 2
	q.X.Mul(px, &t2) // XQ'= XQ * t2
	q.Z.Mul(pz, &t0) // ZQ'= ZQ * t0
	return q
}

// Given a four-torsion point p = x(PB) on the curve E_(A:C), construct the
// four-isogeny phi : E_(A:C) -> E_(A:C)/<P_4> = E_(A':C').
//
// Input: (XP_4: ZP_4), where P_4 has exact order 4 on E_A/C
// Output: * Curve coordinates (A' + 2C', 4C') corresponding to E_A'/C' = A_E/C/<P4>
//		   * Isogeny phi with constants in F_p^2
func (phi *isogeny4) GenerateCurve(p *ProjectivePoint) CurveCoefficientsEquiv {
	var coefEq CurveCoefficientsEquiv
	var xp4, zp4 = &p.X, &p.Z
	var K1, K2, K3 = &phi.K1, &phi.K2, &phi.K3

	K2.Sub(xp4, zp4)
	K3.Add(xp4, zp4)
	K1.Square(zp4)
	K1.Add(K1, K1)
	coefEq.C.Square(K1)
	K1.Add(K1, K1)
	coefEq.A.Square(xp4)
	coefEq.A.Add(&coefEq.A, &coefEq.A)
	coefEq.A.Square(&coefEq.A)
	return coefEq
}

// Given a 4-isogeny phi and a point xP = x(P), compute x(Q), the x-coordinate
// of the image Q = phi(P) of P under phi : E_(A:C) -> E_(A':C').
//
// Input: Isogeny returned by GenerateCurve and point q=(Qx,Qz) from E0_A/C
// Output: Corresponding point q from E1_A'/C', where E1 is 4-isogenous to E0
func (phi *isogeny4) EvaluatePoint(p *ProjectivePoint) ProjectivePoint {
	var t0, t1 ExtensionFieldElement
	var q = *p
	var xq, zq = &q.X, &q.Z
	var K1, K2, K3 = &phi.K1, &phi.K2, &phi.K3

	t0.Add(xq, zq)
	t1.Sub(xq, zq)
	xq.Mul(&t0, K2)
	zq.Mul(&t1, K3)
	t0.Mul(&t0, &t1)
	t0.Mul(&t0, K1)
	t1.Add(xq, zq)
	zq.Sub(xq, zq)
	t1.Square(&t1)
	zq.Square(zq)
	xq.Add(&t0, &t1)
	t0.Sub(zq, &t0)
	xq.Mul(xq, &t1)
	zq.Mul(zq, &t0)
	return q
}
