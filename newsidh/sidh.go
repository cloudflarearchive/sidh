package sidh

import (
	"errors"
	"io"
)

import . "github.com/cloudflare/p751sidh/p751toolbox"

type SIDH struct {
	// The secret key size, in bytes.
	secretKeySize int
	// The public key size, in bytes.
	publicKeySize int
	// The shared secret size, in bytes.
	sharedSecretSize int
	// Fixed parameters for isogeny tree computations
	maxAlice int

	maxBob   int
	// Alice's mask values
	maskAliceByte1 byte

	maskAliceByte2 byte

	maskAliceByte3 byte
	// Bob's mask value
	maskBobByte byte
	// Sample rate to obtain a value in [0,3^238]
	sampleRate int
	// Alice's isogeny optimal strategy
	aliceIsogenyStrategy []int
	// Bob's isogeny optimal strategy
	bobIsogenyStrategy []int
}

// Bob's public key.
type SIDHPublicKeyBob struct {
	affine_xP   ExtensionFieldElement
	affine_xQ   ExtensionFieldElement
	affine_xQmP ExtensionFieldElement
}

// Alice's public key.
type SIDHPublicKeyAlice struct {
	affine_xP   ExtensionFieldElement
	affine_xQ   ExtensionFieldElement
	affine_xQmP ExtensionFieldElement
}

// Bob's secret key.
type SIDHSecretKeyBob struct {
	Scalar [SecretKeySize]byte
}

// Alice's secret key.
type SIDHSecretKeyAlice struct {
	Scalar [SecretKeySize]byte
}

// Read a public key from a byte slice.  The input must be at least 564 bytes long.
func (e *SIDH) FromBytesBob(input []byte) (pubKey *SIDHPublicKeyBob) {
	if len(input) < e.publicKeySize {
		panic("Too short input to SIDH pubkey FromBytes, expected 564 bytes")
	}
	pubKey.affine_xP.FromBytes(input[0:e.sharedSecretSize])
	pubKey.affine_xQ.FromBytes(input[e.sharedSecretSize : 2*e.sharedSecretSize])
	pubKey.affine_xQmP.FromBytes(input[2*e.sharedSecretSize : 3*e.sharedSecretSize])

	return pubKey
}

// Write a public key to a byte slice.  The output must be at least 564 bytes long.
func (e *SIDH) ToBytesBob(pubKey *SIDHPublicKeyBob) (output []byte) {
	if len(output) < e.publicKeySize {
		panic("Too short output for SIDH pubkey FromBytes, expected 564 bytes")
	}
	pubKey.affine_xP.ToBytes(output[0:e.sharedSecretSize])
	pubKey.affine_xQ.ToBytes(output[e.sharedSecretSize : 2*e.sharedSecretSize])
	pubKey.affine_xQmP.ToBytes(output[2*e.sharedSecretSize : 3*e.sharedSecretSize])

	return output
}

// Read a public key from a byte slice.  The input must be at least 564 bytes long.
func (e *SIDH) FromBytesAlice(input []byte) (pubKey *SIDHPublicKeyAlice) {
	if len(input) < e.publicKeySize {
		panic("Too short input to SIDH pubkey FromBytes, expected 564 bytes")
	}
	pubKey.affine_xP.FromBytes(input[0:e.sharedSecretSize])
	pubKey.affine_xQ.FromBytes(input[e.sharedSecretSize : 2*e.sharedSecretSize])
	pubKey.affine_xQmP.FromBytes(input[2*e.sharedSecretSize : 3*e.sharedSecretSize])

	return pubKey
}

// Write a public key to a byte slice.  The output must be at least 564 bytes long.
func (e *SIDH) ToBytesAlice(pubKey *SIDHPublicKeyAlice) (output []byte) {
	if len(output) < e.publicKeySize {
		panic("Too short output for SIDH pubkey FromBytes, expected 564 bytes")
	}
	pubKey.affine_xP.ToBytes(output[0:e.sharedSecretSize])
	pubKey.affine_xQ.ToBytes(output[e.sharedSecretSize : 2*e.sharedSecretSize])
	pubKey.affine_xQmP.ToBytes(output[2*e.sharedSecretSize : 3*e.sharedSecretSize])

	return output
}

// Generate a keypair for "Alice".  Note that because this library does not
// implement SIDH validation, each keypair should be used for at most one
// shared secret computation.
func (e *SIDH) GenerateAliceKeypair(rand io.Reader) (publicKey *SIDHPublicKeyAlice, secretKey *SIDHSecretKeyAlice, err error) {
	publicKey = new(SIDHPublicKeyAlice)
	secretKey = new(SIDHSecretKeyAlice)

	_, err = io.ReadFull(rand, secretKey.Scalar[:])
	if err != nil {
		return nil, nil, err
	}
	// Bit-twiddle to ensure scalar is in 2*[0,2^371):
	secretKey.Scalar[e.secretKeySize-1] = e.maskAliceByte1
	secretKey.Scalar[e.secretKeySize-2] &= e.maskAliceByte2 // clear high bits, so scalar < 2^372
	secretKey.Scalar[0] &= e.maskAliceByte3                 // clear low bit, so scalar is even

	// We actually want scalar in 2*(0,2^371), but the above procedure
	// generates 0 with probability 2^(-371), which isn't worth checking
	// for.

	*publicKey = e.PublicKeyAlice(secretKey)

	return
}

// Set result to zero if the input scalar is <= 3^238.
//go:noescape
func checkLessThanThree238(scalar *[SecretKeySize]byte, result *uint32)

// Set scalar = 3*scalar
//go:noescape
func multiplyByThree(scalar *[SecretKeySize]byte)

// Generate a keypair for "Bob".  Note that because this library does not
// implement SIDH validation, each keypair should be used for at most one
// shared secret computation.
func (e *SIDH) GenerateBobKeypair(rand io.Reader) (publicKey *SIDHPublicKeyBob, secretKey *SIDHSecretKeyBob, err error) {
	publicKey = new(SIDHPublicKeyBob)
	secretKey = new(SIDHSecretKeyBob)

	// Perform rejection sampling to obtain a random value in [0,3^238]:
	var ok uint32
	for i := 0; i < e.sampleRate; i++ {
		_, err = io.ReadFull(rand, secretKey.Scalar[:])
		if err != nil {
			return nil, nil, err
		}
		// Mask the high bits to obtain a uniform value in [0,2^378):
		secretKey.Scalar[e.secretKeySize-1] &= e.maskBobByte
		// Accept if scalar < 3^238 (this happens w/ prob ~0.5828)
		checkLessThanThree238(&secretKey.Scalar, &ok)
		if ok == 0 {
			break
		}
	}
	// ok is nonzero if all sampleRate trials failed.
	// This happens with probability 0.41719...^102 < 2^(-128), i.e., never
	if ok != 0 {
		return nil, nil, errors.New("WOW! An event with probability < 2^(-128) occurred!!")
	}

	// Multiply by 3 to get a scalar in 3*[0,3^238):
	multiplyByThree(&secretKey.Scalar)

	// We actually want scalar in 2*(0,2^371), but the above procedure
	// generates 0 with probability 3^(-238), which isn't worth checking
	// for.

	*publicKey = e.PublicKeyBob(secretKey)

	return
}

// Compute the corresponding public key for the given secret key.
func (e *SIDH) PublicKeyAlice(secretKey *SIDHSecretKeyAlice) SIDHPublicKeyAlice {
	var xP, xQ, xQmP, xR ProjectivePoint

	xP.FromAffinePrimeField(&Affine_xPB)     // = ( x_P : 1) = x(P_B)
	xQ.FromAffinePrimeField(&Affine_xPB)     //
	xQ.X.Neg(&xQ.X)                          // = (-x_P : 1) = x(Q_B)
	xQmP = DistortAndDifference(&Affine_xPB) // = x(Q_B - P_B)

	xR = SecretPoint(&Affine_xPA, &Affine_yPA, secretKey.Scalar[:])

	var currentCurve ProjectiveCurveParameters
	// Starting curve has a = 0, so (A:C) = (0,1)
	currentCurve.A.Zero()
	currentCurve.C.One()

	var firstPhi FirstFourIsogeny
	currentCurve, firstPhi = ComputeFirstFourIsogeny(&currentCurve)

	xP = firstPhi.Eval(&xP)
	xQ = firstPhi.Eval(&xQ)
	xQmP = firstPhi.Eval(&xQmP)
	xR = firstPhi.Eval(&xR)

	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var phi FourIsogeny

	var i = 0

	for j := 1; j < e.maxAlice; j++ {
		for i < e.maxAlice-j {
			points = append(points, xR)
			indices = append(indices, i)
			k := int(e.aliceIsogenyStrategy[e.maxAlice-i-j])
			xR.Pow2k(&currentCurve, &xR, uint32(2*k))
			i = i + k
		}
		currentCurve, phi = ComputeFourIsogeny(&xR)

		for k := 0; k < len(points); k++ {
			points[k] = phi.Eval(&points[k])
		}

		xP = phi.Eval(&xP)
		xQ = phi.Eval(&xQ)
		xQmP = phi.Eval(&xQmP)

		// pop xR from points
		xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}

	currentCurve, phi = ComputeFourIsogeny(&xR)
	xP = phi.Eval(&xP)
	xQ = phi.Eval(&xQ)
	xQmP = phi.Eval(&xQmP)

	var invZP, invZQ, invZQmP ExtensionFieldElement
	ExtensionFieldBatch3Inv(&xP.Z, &xQ.Z, &xQmP.Z, &invZP, &invZQ, &invZQmP)

	var publicKey SIDHPublicKeyAlice
	publicKey.affine_xP.Mul(&xP.X, &invZP)
	publicKey.affine_xQ.Mul(&xQ.X, &invZQ)
	publicKey.affine_xQmP.Mul(&xQmP.X, &invZQmP)

	return publicKey
}

// Compute the public key corresponding to the secret key.
func (e *SIDH) PublicKeyBob(secretKey *SIDHSecretKeyBob) SIDHPublicKeyBob {
	var xP, xQ, xQmP, xR ProjectivePoint

	xP.FromAffinePrimeField(&Affine_xPA)     // = ( x_P : 1) = x(P_A)
	xQ.FromAffinePrimeField(&Affine_xPA)     //
	xQ.X.Neg(&xQ.X)                          // = (-x_P : 1) = x(Q_A)
	xQmP = DistortAndDifference(&Affine_xPA) // = x(Q_B - P_B)

	xR = SecretPoint(&Affine_xPB, &Affine_yPB, secretKey.Scalar[:])

	var currentCurve ProjectiveCurveParameters
	// Starting curve has a = 0, so (A:C) = (0,1)
	currentCurve.A.Zero()
	currentCurve.C.One()

	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var phi ThreeIsogeny

	var i = 0

	for j := 1; j < e.maxBob; j++ {
		for i < e.maxBob-j {
			points = append(points, xR)
			indices = append(indices, i)
			k := int(e.bobIsogenyStrategy[e.maxBob-i-j])
			xR.Pow3k(&currentCurve, &xR, uint32(k))
			i = i + k
		}
		currentCurve, phi = ComputeThreeIsogeny(&xR)

		for k := 0; k < len(points); k++ {
			points[k] = phi.Eval(&points[k])
		}

		xP = phi.Eval(&xP)
		xQ = phi.Eval(&xQ)
		xQmP = phi.Eval(&xQmP)

		// pop xR from points
		xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}

	currentCurve, phi = ComputeThreeIsogeny(&xR)
	xP = phi.Eval(&xP)
	xQ = phi.Eval(&xQ)
	xQmP = phi.Eval(&xQmP)

	var invZP, invZQ, invZQmP ExtensionFieldElement
	ExtensionFieldBatch3Inv(&xP.Z, &xQ.Z, &xQmP.Z, &invZP, &invZQ, &invZQmP)

	var publicKey SIDHPublicKeyBob
	publicKey.affine_xP.Mul(&xP.X, &invZP)
	publicKey.affine_xQ.Mul(&xQ.X, &invZQ)
	publicKey.affine_xQmP.Mul(&xQmP.X, &invZQmP)

	return publicKey
}

// Compute (Alice's view of) a shared secret using Alice's secret key and Bob's public key.
func (e *SIDH) SharedSecretAlice(aliceSecret *SIDHSecretKeyAlice, bobPublic *SIDHPublicKeyBob) []byte {
	var currentCurve = RecoverCurveParameters(&bobPublic.affine_xP, &bobPublic.affine_xQ, &bobPublic.affine_xQmP)

	var xR, xP, xQ, xQmP ProjectivePoint

	xP.FromAffine(&bobPublic.affine_xP)
	xQ.FromAffine(&bobPublic.affine_xQ)
	xQmP.FromAffine(&bobPublic.affine_xQmP)

	xR.RightToLeftLadder(&currentCurve, &xP, &xQ, &xQmP, aliceSecret.Scalar[:])

	var firstPhi FirstFourIsogeny
	currentCurve, firstPhi = ComputeFirstFourIsogeny(&currentCurve)
	xR = firstPhi.Eval(&xR)

	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var phi FourIsogeny

	var i = 0

	for j := 1; j < e.maxAlice; j++ {
		for i < e.maxAlice-j {
			points = append(points, xR)
			indices = append(indices, i)
			k := int(e.aliceIsogenyStrategy[e.maxAlice-i-j])
			xR.Pow2k(&currentCurve, &xR, uint32(2*k))
			i = i + k
		}
		currentCurve, phi = ComputeFourIsogeny(&xR)

		for k := 0; k < len(points); k++ {
			points[k] = phi.Eval(&points[k])
		}

		// pop xR from points
		xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}

	currentCurve, _ = ComputeFourIsogeny(&xR)

	var sharedSecret = make([]byte, e.sharedSecretSize)
	var jInv = currentCurve.JInvariant()
	jInv.ToBytes(sharedSecret[:])
	return sharedSecret
}

// Compute (Bob's view of) a shared secret using Bob's secret key and Alice's public key.
func (e *SIDH) SharedSecretBob(bobSecret *SIDHSecretKeyBob, alicePublic *SIDHPublicKeyAlice) []byte {
	var currentCurve = RecoverCurveParameters(&alicePublic.affine_xP, &alicePublic.affine_xQ, &alicePublic.affine_xQmP)

	var xR, xP, xQ, xQmP ProjectivePoint

	xP.FromAffine(&alicePublic.affine_xP)
	xQ.FromAffine(&alicePublic.affine_xQ)
	xQmP.FromAffine(&alicePublic.affine_xQmP)

	xR.RightToLeftLadder(&currentCurve, &xP, &xQ, &xQmP, bobSecret.Scalar[:])

	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var phi ThreeIsogeny

	var i = 0

	for j := 1; j < e.maxBob; j++ {
		for i < e.maxBob-j {
			points = append(points, xR)
			indices = append(indices, i)
			k := int(e.bobIsogenyStrategy[e.maxBob-i-j])
			xR.Pow3k(&currentCurve, &xR, uint32(k))
			i = i + k
		}
		currentCurve, phi = ComputeThreeIsogeny(&xR)

		for k := 0; k < len(points); k++ {
			points[k] = phi.Eval(&points[k])
		}

		// pop xR from points
		xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}
	currentCurve, _ = ComputeThreeIsogeny(&xR)

	var sharedSecret = make([]byte, e.sharedSecretSize)
	var jInv = currentCurve.JInvariant()
	jInv.ToBytes(sharedSecret[:])
	return sharedSecret
}
