// +build amd64,!noasm

package p503

import (
	. "github.com/cloudflare/p751sidh/internal/isogeny"
	cpu "github.com/cloudflare/p751sidh/internal/utils"
)

// If choice = 0, leave x,y unchanged. If choice = 1, set x,y = y,x.
// If choice is neither 0 nor 1 then behaviour is undefined.
// This function executes in constant time.
//go:noescape
func fp503ConditionalSwap(x, y *FpElement, choice uint8)

// Compute z = x + y (mod p).
//go:noescape
func fp503AddReduced(z, x, y *FpElement)

// Compute z = x - y (mod p).
//go:noescape
func fp503SubReduced(z, x, y *FpElement)

// Compute z = x + y, without reducing mod p.
//go:noescape
func fp503AddLazy(z, x, y *FpElement)

// Compute z = x + y, without reducing mod p.
//go:noescape
func fp503X2AddLazy(z, x, y *FpElementX2)

// Compute z = x - y, without reducing mod p.
//go:noescape
func fp503X2SubLazy(z, x, y *FpElementX2)

// Reduce a field element in [0, 2*p) to one in [0,p).
//go:noescape
func fp503StrongReduce(x *FpElement)

// Function pointer to function computing z = x * y.
// Concrete implementation depends on capabilities of the CPU which
// are resolved at runtime. CPUs with ADCX, ADOX and MULX support
// run most optimized implementation
var fp503Mul func(z *FpElementX2, x, y *FpElement)

// Mul implementattion for legacy CPUs
//go:noescape
func mul(z *FpElementX2, x, y *FpElement)

// Mul implementation for CPUs supporting carry-less MULX multiplier.
//go:noescape
func mulWithMULX(z *FpElementX2, x, y *FpElement)

// Mul implementation for CPUs supporting two independent carry chain
// (ADOX/ADCX) instructions and carry-less MULX multiplier
//go:noescape
func mulWithMULXADX(z *FpElementX2, x, y *FpElement)

// Computes the Montgomery reduction z = x R^{-1} (mod 2*p). On return value
// of x may be changed. z=x not allowed.
var fp503MontgomeryReduce func(z *FpElement, x *FpElementX2)

func redc(z *FpElement, x *FpElementX2)

// Mul implementation for CPUs supporting carry-less MULX multiplier.
//go:noescape
func redcWithMULX(z *FpElement, x *FpElementX2)

// Mul implementation for CPUs supporting two independent carry chain
// (ADOX/ADCX) instructions and carry-less MULX multiplier
//go:noescape
func redcWithMULXADX(z *FpElement, x *FpElementX2)

// On initialization, set the fp503Mul function pointer to the
// fastest implementation depending on CPU capabilities.
func init() {
	if cpu.HasBMI2 {
		if cpu.HasADX {
			fp503Mul = mulWithMULXADX
			fp503MontgomeryReduce = redcWithMULXADX
		} else {
			fp503Mul = mulWithMULX
			fp503MontgomeryReduce = redcWithMULX
		}
	} else {
		fp503Mul = mul
		fp503MontgomeryReduce = redc
	}
}
