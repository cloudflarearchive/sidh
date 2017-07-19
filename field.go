package cln16sidh

//------------------------------------------------------------------------------
// Extension Field
//------------------------------------------------------------------------------

// Represents an element of the extension field F_{p^2}.
type ExtensionFieldElement struct {
	// This field element is in Montgomery form, so that the value `a` is
	// represented by `aR mod p`.
	a fp751Element
	// This field element is in Montgomery form, so that the value `b` is
	// represented by `bR mod p`.
	b fp751Element
}

// Set dest = lhs * rhs.
//
// Allowed to overlap lhs or rhs with dest.
func (dest *ExtensionFieldElement) Mul(lhs, rhs *ExtensionFieldElement) {
	// Let (a,b,c,d) = (lhs.a,lhs.b,rhs.a,rhs.b).
	a := &lhs.a
	b := &lhs.b
	c := &rhs.a
	d := &rhs.b

	// We want to compute
	//
	// (a + bi)*(c + di) = (a*c - b*d) + (a*d + b*c)i
	//
	// Use Karatsuba's trick: note that
	//
	// (b - a)*(c - d) = (b*c + a*d) - a*c - b*d
	//
	// so (a*d + b*c) = (b-a)*(c-d) + a*c + b*d.

	var ac, bd fp751X2
	fp751Mul(&ac, a, c)				// = a*c*R*R
	fp751Mul(&bd, b, d)				// = b*d*R*R

	var b_minus_a, c_minus_d fp751Element
	fp751SubReduced(&b_minus_a, b, a)		// = (b-a)*R
	fp751SubReduced(&c_minus_d, c, d)		// = (c-d)*R

	var ad_plus_bc fp751X2
	fp751Mul(&ad_plus_bc, &b_minus_a, &c_minus_d)	// = (b-a)*(c-d)*R*R
	fp751X2AddLazy(&ad_plus_bc, &ad_plus_bc, &ac)	// = ((b-a)*(c-d) + a*c)*R*R
	fp751X2AddLazy(&ad_plus_bc, &ad_plus_bc, &bd)	// = ((b-a)*(c-d) + a*c + b*d)*R*R

	fp751MontgomeryReduce(&dest.b, &ad_plus_bc)	// = (a*d + b*c)*R mod p

	var ac_minus_bd fp751X2
	fp751X2SubLazy(&ac_minus_bd, &ac, &bd)		// = (a*c - b*d)*R*R
	fp751MontgomeryReduce(&dest.a, &ac_minus_bd)	// = (a*c - b*d)*R mod p
}

// Set dest = x * x
//
// Allowed to overlap dest with x.
func (dest *ExtensionFieldElement) Sqr(x *ExtensionFieldElement) {
	a := &x.a
	b := &x.b

	// We want to compute
	//
	// (a + bi)*(a + bi) = (a^2 - b^2) + 2abi.

	var a2, a_plus_b, a_minus_b fp751Element
	fp751AddReduced(&a2, a, a)				// = a*R + a*R = 2*a*R
	fp751AddReduced(&a_plus_b, a, b)			// = a*R + b*R = (a+b)*R
	fp751SubReduced(&a_minus_b, a, b)			// = a*R - b*R = (a-b)*R

	var asq_minus_bsq, ab2 fp751X2
	fp751Mul(&asq_minus_bsq, &a_plus_b, &a_minus_b)		// = (a+b)*(a-b)*R*R = (a^2 - b^2)*R*R
	fp751Mul(&ab2, &a2, b)					// = 2*a*b*R*R

	fp751MontgomeryReduce(&dest.a, &asq_minus_bsq)		// = (a^2 - b^2)*R mod p
	fp751MontgomeryReduce(&dest.b, &ab2)			// = 2*a*b*R mod p
}

// Set dest = lhs + rhs.
//
// Allowed to overlap lhs or rhs with dest.
func (dest *ExtensionFieldElement) Add(lhs, rhs *ExtensionFieldElement) {
	fp751AddReduced(&dest.a, &lhs.a, &rhs.a)
	fp751AddReduced(&dest.b, &lhs.b, &rhs.b)
}

// Set dest = lhs - rhs.
//
// Allowed to overlap lhs or rhs with dest.
func (dest *ExtensionFieldElement) Sub(lhs, rhs *ExtensionFieldElement) {
	fp751SubReduced(&dest.a, &lhs.a, &rhs.a)
	fp751SubReduced(&dest.b, &lhs.b, &rhs.b)
}

// Returns true if lhs = rhs.  Takes variable time.
func (lhs *ExtensionFieldElement) VartimeEq(rhs *ExtensionFieldElement) bool {
	return lhs.a.vartimeEq(rhs.a) && lhs.b.vartimeEq(rhs.b)
}

//------------------------------------------------------------------------------
// Prime Field
//------------------------------------------------------------------------------

// Represents an element of the prime field F_p.
type PrimeFieldElement struct {
	// This field element is in Montgomery form, so that the value `a` is
	// represented by `aR mod p`.
	a fp751Element
}

// Set dest = lhs * rhs.
//
// Allowed to overlap lhs or rhs with dest.
func (dest *PrimeFieldElement) Mul(lhs, rhs *PrimeFieldElement) {
	a := &lhs.a				// = a*R
	b := &rhs.a				// = b*R

	var ab fp751X2
	fp751Mul(&ab, a, b)			// = a*b*R*R
	fp751MontgomeryReduce(&dest.a, &ab)	// = a*b*R mod p
}

// Set dest = x^(2^k), for k >= 1, by repeated squarings.
//
// Allowed to overlap x with dest.
func (dest *PrimeFieldElement) Pow2k(x *PrimeFieldElement, k uint8) {
	dest.Sqr(x)
	for i := uint8(1); i < k; i++ {
		dest.Sqr(dest)
	}
}

// Set dest = x^2
//
// Allowed to overlap x with dest.
func (dest *PrimeFieldElement) Sqr(x *PrimeFieldElement) {
	a := &x.a				// = a*R
	b := &x.a				// = b*R

	var ab fp751X2
	fp751Mul(&ab, a, b)			// = a*b*R*R
	fp751MontgomeryReduce(&dest.a, &ab)	// = a*b*R mod p
}

// Set dest = lhs + rhs.
//
// Allowed to overlap lhs or rhs with dest.
func (dest *PrimeFieldElement) Add(lhs, rhs *PrimeFieldElement) {
	fp751AddReduced(&dest.a, &lhs.a, &rhs.a)
}

// Set dest = lhs - rhs.
//
// Allowed to overlap lhs or rhs with dest.
func (dest *PrimeFieldElement) Sub(lhs, rhs *PrimeFieldElement) {
	fp751SubReduced(&dest.a, &lhs.a, &rhs.a)
}

// Returns true if lhs = rhs.  Takes variable time.
func (lhs *PrimeFieldElement) VartimeEq(rhs *PrimeFieldElement) bool {
	return lhs.a.vartimeEq(rhs.a)
}

// Set dest = x^((p-3)/4)
//
// Allowed to overlap x with dest.
func (dest *PrimeFieldElement) P34(x *PrimeFieldElement) {
	// Sliding-window strategy computed with Sage, awk, sed, and tr
	powStrategy := [137]uint8{5, 7, 6, 2, 10, 4, 6, 9, 8, 5, 9, 4, 7, 5, 5, 4, 8, 3, 9, 5, 5, 4, 10, 4, 6, 6, 6, 5, 8, 9, 3, 4, 9, 4, 5, 6, 6, 2, 9, 4, 5, 5, 5, 7, 7, 9, 4, 6, 4, 8, 5, 8, 6, 6, 2, 9, 7, 4, 8, 8, 8, 4, 6, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 2}
	mulStrategy := [137]uint8{31, 23, 21, 1, 31, 7, 7, 7, 9, 9, 19, 15, 23, 23, 11, 7, 25, 5, 21, 17, 11, 5, 17, 7, 11, 9, 23, 9, 1, 19, 5, 3, 25, 15, 11, 29, 31, 1, 29, 11, 13, 9, 11, 27, 13, 19, 15, 31, 3, 29, 23, 31, 25, 11, 1, 21, 19, 15, 15, 21, 29, 13, 23, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 31, 3}
	initialMul := uint8(27)

	// Build a lookup table of odd multiples of x.
	lookup := [16]PrimeFieldElement{}
	xx := &PrimeFieldElement{}
	xx.Sqr(x) // Set xx = x^2
	lookup[0] = *x
	for i := 1; i < 16; i++ {
		lookup[i].Mul(&lookup[i-1], xx)
	}
	// Now lookup = {x, x^3, x^5, ... }
	// so that lookup[i] = x^{2*i + 1}
	// so that lookup[k/2] = x^k, for odd k

	*dest = lookup[initialMul/2]
	for i := uint8(0); i < 137; i++ {
		dest.Pow2k(dest, powStrategy[i])
		dest.Mul(dest, &lookup[mulStrategy[i]/2])
	}
}

//------------------------------------------------------------------------------
// Internals
//------------------------------------------------------------------------------

const fp751NumWords = 12

// Internal representation of an element of the base field F_p.
//
// This type is distinct from PrimeFieldElement in that no particular meaning
// is assigned to the representation -- it could represent an element in
// Montgomery form, or not.  Tracking the meaning of the field element is left
// to higher types.
type fp751Element [fp751NumWords]uint64

// Represents an intermediate product of two elements of the base field F_p.
type fp751X2 [2 * fp751NumWords]uint64

// Compute z = x + y (mod p).
//go:noescape
func fp751AddReduced(z, x, y *fp751Element)

// Compute z = x - y (mod p).
//go:noescape
func fp751SubReduced(z, x, y *fp751Element)

// Compute z = x + y, without reducing mod p.
//go:noescape
func fp751AddLazy(z, x, y *fp751Element)

// Compute z = x + y, without reducing mod p.
//go:noescape
func fp751X2AddLazy(z, x, y *fp751X2)

// Compute z = x - y, without reducing mod p.
//go:noescape
func fp751X2SubLazy(z, x, y *fp751X2)

// Compute z = x * y.
//go:noescape
func fp751Mul(z *fp751X2, x, y *fp751Element)

// Perform Montgomery reduction: set z = x R^{-1} (mod p).
// Destroys the input value.
//go:noescape
func fp751MontgomeryReduce(z *fp751Element, x *fp751X2)

// Reduce a field element in [0, 2*p) to one in [0,p).
//go:noescape
func fp751StrongReduce(x *fp751Element)

func (x fp751Element) vartimeEq(y fp751Element) bool {
	fp751StrongReduce(&x)
	fp751StrongReduce(&y)
	eq := true
	for i := 0; i < fp751NumWords; i++ {
		eq = (x[i] == y[i]) && eq
	}

	return eq
}
