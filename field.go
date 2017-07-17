package cln16sidh

const Fp751NumWords = 12

// Represents an element of the base field F_p, in Montgomery form.
type Fp751Element [Fp751NumWords]uint64

// Represents an intermediate product of two elements of the base field F_p.
type Fp751X2 [2 * Fp751NumWords]uint64

// Compute z = x + y (mod p).
//go:noescape
func Fp751AddReduced(z, x, y *Fp751Element)

// Compute z = x - y (mod p).
//go:noescape
func Fp751SubReduced(z, x, y *Fp751Element)

// Compute z = x + y, without reducing mod p.
//go:noescape
func Fp751AddLazy(z, x, y *Fp751Element)

// Compute z = x + y, without reducing mod p.
//go:noescape
func Fp751X2AddLazy(z, x, y *Fp751X2)

// Compute z = x * y.
//go:noescape
func Fp751Mul(z *Fp751X2, x, y *Fp751Element)

// Reduce an X2 to a field element: set z = x (mod p).
// Destroys the input value.
//go:noescape
func Fp751Reduce(z *Fp751Element, x *Fp751X2)
