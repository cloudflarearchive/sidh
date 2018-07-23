// +build amd64,!noasm

package sidh

// Set result to zero if the input scalar is <= 3^238, otherwise result is 1.
// Scalar must be array of 48 bytes. This function is specific to P751.
//go:noescape
func checkLessThanThree238(scalar []byte) uint8

// Multiply 48-byte scalar by 3 to get a scalar in 3*[0,3^238). This
// function is specific to P751.
//go:noescape
func multiplyByThree(scalar []byte)
