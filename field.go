package cln16sidh

const Fp751NumWords = 12

type Fp751Element [Fp751NumWords]uint64

type Fp751UnreducedProduct [2 * Fp751NumWords]uint64

/// Compute z = x + y.
func Fp751Add(z, x, y *Fp751Element)

/// Compute z = x + y, without reducing mod p.
func Fp751LazyAdd(z, x, y *Fp751Element)

/// Compute z = x + y, without reducing mod p.
func Fp751UnreducedProductLazyAdd(z, x, y *Fp751UnreducedProduct)

/// Compute z = x - y.
func Fp751Sub(z, x, y *Fp751Element)

/// Compute z = x * y.
func Fp751Mul(z *Fp751UnreducedProduct, x, y *Fp751Element)

/// Reduce an UnreducedProduct to a field element: set z = x (mod p).
/// Destroys the input value.
func Fp751Reduce(z *Fp751Element, x *Fp751UnreducedProduct)
