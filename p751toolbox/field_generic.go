// +build noasm arm64 arm

package p751toolbox

// helper used for uint128 representation
type uint128 struct {
	H, L uint64
}

// Adds 2 64bit digits in constant time.
// Returns result and carry (1 or 0)
func addc64(cin, a, b uint64) (ret, cout uint64) {
	t := a + cin
	ret = b + t
	cout = ((a & b) | ((a | b) & (^ret))) >> 63
	return
}

// Substracts 2 64bit digits in constant time.
// Returns result and borrow (1 or 0)
func subc64(bIn, a, b uint64) (ret, bOut uint64) {
	var tmp1 = a - b
	// Set bOut if bIn!=0 and tmp1==0 in constant time
	bOut = bIn & (1 ^ ((tmp1 | uint64(0-tmp1)) >> 63))
	// Constant time check if x<y
	bOut |= (a ^ ((a ^ b) | (uint64(a-b) ^ b))) >> 63
	ret = tmp1 - bIn
	return
}

// Multiplies 2 64bit digits in constant time
func mul64(a, b uint64) (res uint128) {
	var al, bl, ah, bh, albl, albh, ahbl, ahbh uint64
	var res1, res2, res3 uint64
	var carry, maskL, maskH, temp uint64

	maskL = (^maskL) >> 32
	maskH = ^maskL

	al = a & maskL
	ah = a >> 32
	bl = b & maskL
	bh = b >> 32

	albl = al * bl
	albh = al * bh
	ahbl = ah * bl
	ahbh = ah * bh
	res.L = albl & maskL

	res1 = albl >> 32
	res2 = ahbl & maskL
	res3 = albh & maskL
	temp = res1 + res2 + res3
	carry = temp >> 32
	res.L ^= temp << 32

	res1 = ahbl >> 32
	res2 = albh >> 32
	res3 = ahbh & maskL
	temp = res1 + res2 + res3 + carry
	res.H = temp & maskL
	carry = temp & maskH
	res.H ^= (ahbh & maskH) + carry
	return
}

// Compute z = x + y (mod p).
func fp751AddReduced(z, x, y *Fp751Element) {
	var carry uint64

	// z=x+y % p751
	for i := 0; i < fp751NumWords; i++ {
		z[i], carry = addc64(carry, x[i], y[i])
	}

	// z = z - p751x2
	carry = 0
	for i := 0; i < fp751NumWords; i++ {
		z[i], carry = subc64(carry, z[i], p751x2[i])
	}

	// z = z + p751x2
	mask := uint64(0 - carry)
	carry = 0
	for i := 0; i < fp751NumWords; i++ {
		z[i], carry = addc64(carry, z[i], p751x2[i]&mask)
	}
}

// Compute z = x - y (mod p).
func fp751SubReduced(z, x, y *Fp751Element) {
	var borrow uint64

	for i := 0; i < fp751NumWords; i++ {
		z[i], borrow = subc64(borrow, x[i], y[i])
	}

	mask := uint64(0 - borrow)
	borrow = 0

	for i := 0; i < fp751NumWords; i++ {
		z[i], borrow = addc64(borrow, z[i], p751x2[i]&mask)
	}
}

// Conditionally swaps bits in x and y in constant time.
// mask indicates bits to be swaped (set bits are swapped)
// For details see "Hackers Delight, 2.20"
//
// Implementation doesn't actually depend on a prime field.
func fp751ConditionalSwap(x, y *Fp751Element, mask uint8) {
	var tmp, mask64 uint64

	mask64 = 0 - uint64(mask)
	for i := 0; i < len(x); i++ {
		tmp = mask64 & (x[i] ^ y[i])
		x[i] = tmp ^ x[i]
		y[i] = tmp ^ y[i]
	}
}

// Perform Montgomery reduction: set z = x R^{-1} (mod 2*p)
// with R=2^768. Destroys the input value.
func fp751MontgomeryReduce(z *Fp751Element, x *fp751X2) {
	var carry, t, u, v uint64
	var uv uint128
	var count int

	count = 5 // number of 0 digits in the least significat part of p751 + 1

	for i := 0; i < fp751NumWords; i++ {
		for j := 0; j < i; j++ {
			if j < (i - count + 1) {
				uv = mul64(z[j], p751p1[i-j])
				v, carry = addc64(0, uv.L, v)
				u, carry = addc64(carry, uv.H, u)
				t += carry
			}
		}
		v, carry = addc64(0, v, x[i])
		u, carry = addc64(carry, u, 0)
		t += carry

		z[i] = v
		v = u
		u = t
		t = 0
	}

	for i := fp751NumWords; i < 2*fp751NumWords-1; i++ {
		if count > 0 {
			count--
		}
		for j := i - fp751NumWords + 1; j < fp751NumWords; j++ {
			if j < (fp751NumWords - count) {
				uv = mul64(z[j], p751p1[i-j])
				v, carry = addc64(0, uv.L, v)
				u, carry = addc64(carry, uv.H, u)
				t += carry
			}
		}
		v, carry = addc64(0, v, x[i])
		u, carry = addc64(carry, u, 0)

		t += carry
		z[i-fp751NumWords] = v
		v = u
		u = t
		t = 0
	}
	v, carry = addc64(0, v, x[2*fp751NumWords-1])
	z[fp751NumWords-1] = v
}

// Compute z = x * y.
func fp751Mul(z *fp751X2, x, y *Fp751Element) {
	var u, v, t uint64
	var carry uint64
	var uv uint128

	for i := uint64(0); i < fp751NumWords; i++ {
		for j := uint64(0); j <= i; j++ {
			uv = mul64(x[j], y[i-j])
			v, carry = addc64(0, uv.L, v)
			u, carry = addc64(carry, uv.H, u)
			t += carry
		}
		z[i] = v
		v = u
		u = t
		t = 0
	}

	for i := fp751NumWords; i < (2*fp751NumWords)-1; i++ {
		for j := i - fp751NumWords + 1; j < fp751NumWords; j++ {
			uv = mul64(x[j], y[i-j])
			v, carry = addc64(0, uv.L, v)
			u, carry = addc64(carry, uv.H, u)
			t += carry
		}
		z[i] = v
		v = u
		u = t
		t = 0
	}
	z[2*fp751NumWords-1] = v
}

// Compute z = x + y, without reducing mod p.
func fp751AddLazy(z, x, y *Fp751Element) {
	var carry uint64
	for i := 0; i < fp751NumWords; i++ {
		z[i], carry = addc64(carry, x[i], y[i])
	}
}

// Compute z = x + y, without reducing mod p.
func fp751X2AddLazy(z, x, y *fp751X2) {
	var carry uint64
	for i := 0; i < 2*fp751NumWords; i++ {
		z[i], carry = addc64(carry, x[i], y[i])
	}
}

// Reduce a field element in [0, 2*p) to one in [0,p).
func fp751StrongReduce(x *Fp751Element) {
	var borrow, mask uint64
	for i := 0; i < fp751NumWords; i++ {
		x[i], borrow = subc64(borrow, x[i], p751[i])
	}

	// Sets all bits if borrow = 1
	mask = 0 - borrow
	borrow = 0
	for i := 0; i < fp751NumWords; i++ {
		x[i], borrow = addc64(borrow, x[i], p751[i]&mask)
	}
}

// Compute z = x - y, without reducing mod p.
func fp751X2SubLazy(z, x, y *fp751X2) {
	var borrow, mask uint64
	for i := 0; i < len(z); i++ {
		z[i], borrow = subc64(borrow, x[i], y[i])
	}

	// Sets all bits if borrow = 1
	mask = 0 - borrow
	borrow = 0
	for i := fp751NumWords; i < len(z); i++ {
		z[i], borrow = addc64(borrow, z[i], p751[i-fp751NumWords]&mask)
	}
}
