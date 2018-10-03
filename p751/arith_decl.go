// +build amd64,!noasm

package p751

import (
	. "github.com/cloudflare/p751sidh/internal/isogeny"
	cpu "github.com/cloudflare/p751sidh/internal/utils"
)

// If choice = 0, leave x,y unchanged. If choice = 1, set x,y = y,x.
// If choice is neither 0 nor 1 then behaviour is undefined.
// This function executes in constant time.
//go:noescape
func fp751ConditionalSwap(x, y *FpElement, choice uint8)

// Compute z = x + y (mod p).
//go:noescape
func fp751AddReduced(z, x, y *FpElement)

// Compute z = x - y (mod p).
//go:noescape
func fp751SubReduced(z, x, y *FpElement)

// Compute z = x + y, without reducing mod p.
//go:noescape
func fp751AddLazy(z, x, y *FpElement)

// Compute z = x + y, without reducing mod p.
//go:noescape
func fp751X2AddLazy(z, x, y *FpElementX2)

// Compute z = x - y, without reducing mod p.
//go:noescape
func fp751X2SubLazy(z, x, y *FpElementX2)

// Compute z = x * y.
//go:noescape
func fp751Mul(z *FpElementX2, x, y *FpElement)

// Function pointer that should point to one of the
// fp751MontgomeryReduce implementations below.
// When set, it performs Montgomery reduction: set z = x R^{-1} (mod 2*p).
// It may destroy the input value.
var fp751MontgomeryReduce func(z *FpElement, x *FpElementX2)

//go:noescape
func fp751MontgomeryReduceBMI2ADX(z *FpElement, x *FpElementX2)

//go:noescape
func fp751MontgomeryReduceBMI2(z *FpElement, x *FpElementX2)

//go:noescape
func fp751MontgomeryReduceFallback(z *FpElement, x *FpElementX2)

// Reduce a field element in [0, 2*p) to one in [0,p).
//go:noescape
func fp751StrongReduce(x *FpElement)

// On initialization, set the fp751MontgomeryReduce function pointer to the
// fastest implementation depending on CPU capabilities.
func init() {
	if cpu.HasBMI2 {
		if cpu.HasADX {
			fp751MontgomeryReduce = fp751MontgomeryReduceBMI2ADX
		} else {
			fp751MontgomeryReduce = fp751MontgomeryReduceBMI2
		}
	} else {
		fp751MontgomeryReduce = fp751MontgomeryReduceFallback
	}
}
