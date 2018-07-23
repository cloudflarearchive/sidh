// +build !amd64 noasm

package sidh

var three238m1 = []uint8{
	0xf8, 0x84, 0x83, 0x82, 0x8a, 0x71, 0xcd, 0xed,
	0x14, 0x7a, 0x42, 0xd4, 0xbf, 0x35, 0x3b, 0x73,
	0x38, 0xcf, 0xd7, 0x94, 0xcf, 0x29, 0x82, 0xf8,
	0xd6, 0x2a, 0x7c, 0x0c, 0x99, 0x6c, 0xc5, 0x63,
	0xc7, 0x22, 0x42, 0x8f, 0x7e, 0xa8, 0x58, 0xb8,
	0xf5, 0xea, 0x25, 0xb5, 0xc6, 0xc9, 0x54, 0x02}

func addc8(cin, a, b uint8) (ret, cout uint8) {
	t := a + cin
	ret = b + t
	cout = ((a & b) | ((a | b) & (^ret))) >> 7
	return
}

func subc8(bIn, a, b uint8) (ret, bOut uint8) {
	var tmp1 = a - b
	ret = tmp1 - bIn
	// Set bOut if bIn!=0 and tmp1==0 in constant time
	bOut = bIn & (1 ^ ((tmp1 | uint8(0-tmp1)) >> 7))
	// Constant time check if a<b
	bOut |= (a ^ ((a ^ b) | (uint8(a-b) ^ b))) >> 7
	return
}

// Set result to zero if the input scalar is <= 3^238, otherwise result is 1.
// Scalar must be array of 48 bytes. This function is specific to P751.
func checkLessThanThree238(scalar []byte) uint8 {
	var borrow uint8
	for i := 0; i < len(three238m1); i++ {
		_, borrow = subc8(borrow, three238m1[i], scalar[i])
	}
	return borrow
}

// Multiply 48-byte scalar by 3 to get a scalar in 3*[0,3^238). This
// function is specific to P751.
func multiplyByThree(scalar []byte) {
	var carry uint8
	var dbl [48]uint8

	for i := 0; i < len(scalar); i++ {
		dbl[i], carry = addc8(carry, scalar[i], scalar[i])
	}
	for i := 0; i < len(scalar); i++ {
		scalar[i], carry = addc8(carry, dbl[i], scalar[i])
	}
}
