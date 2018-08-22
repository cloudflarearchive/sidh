// +build amd64,!noasm

package p751toolbox

import (
	"golang.org/x/sys/cpu"
)

// If choice = 0, leave x,y unchanged. If choice = 1, set x,y = y,x.
// If choice is neither 0 nor 1 then behaviour is undefined.
// This function executes in constant time.
//go:noescape
func fp751ConditionalSwap(x, y *Fp751Element, choice uint8)

// Compute z = x + y (mod p).
//go:noescape
func fp751AddReduced(z, x, y *Fp751Element)

// Compute z = x - y (mod p).
//go:noescape
func fp751SubReduced(z, x, y *Fp751Element)

// Compute z = x + y, without reducing mod p.
//go:noescape
func fp751AddLazy(z, x, y *Fp751Element)

// Compute z = x + y, without reducing mod p.
//go:noescape
func fp751X2AddLazy(z, x, y *fp751X2)

// Compute z = x - y, without reducing mod p.
//go:noescape
func fp751X2SubLazy(z, x, y *fp751X2)

// Compute z = x * y.
//go:noescape
func fp751Mul(z *fp751X2, x, y *Fp751Element)

// Perform Montgomery reduction: set z = x R^{-1} (mod 2*p).
// Destroys the input value.
var fp751MontgomeryReduce func(z *Fp751Element, x *fp751X2)

//go:noescape
func fp751MontgomeryReduceBMI2ADX(z *Fp751Element, x *fp751X2)

//go:noescape
func fp751MontgomeryReduceBMI2(z *Fp751Element, x *fp751X2)

//go:noescape
func fp751MontgomeryReduceFallback(z *Fp751Element, x *fp751X2)

// Reduce a field element in [0, 2*p) to one in [0,p).
//go:noescape
func fp751StrongReduce(x *Fp751Element)

// Choose the fastest variant depending on CPU capabilities.
func init() {
	if cpu.X86.HasBMI2 {
		if cpu.X86.HasADX {
			fp751MontgomeryReduce = fp751MontgomeryReduceBMI2ADX
		} else {
			fp751MontgomeryReduce = fp751MontgomeryReduceBMI2
		}
	} else {
		fp751MontgomeryReduce = fp751MontgomeryReduceFallback
	}
}
