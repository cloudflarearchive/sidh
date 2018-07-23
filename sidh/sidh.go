package sidh

import (
	"errors"
	"io"

	// TODO: This is needed by ExtensionFieldElement struct, which itself
	// 		 depends on implementation of p751.
	. "github.com/cloudflare/p751sidh/p751toolbox"
)

// -----------------------------------------------------------------------------
// Functions for traversing isogeny trees acoording to strategy. Key type 'A' is
//

// Traverses isogeny tree in order to compute xR, xP, xQ and xQmP needed
// for public key generation.
func traverseTreePublicKeyA(curve *ProjectiveCurveParameters, xR, phiP, phiQ, phiR *ProjectivePoint, pub *PublicKey) {
	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var i, sidx int

	cparam := curve.CalcCurveParamsEquiv4()
	phi := NewIsogeny4()
	strat := pub.params.A.IsogenyStrategy
	stratSz := len(strat)

	for j := 1; j <= stratSz; j++ {
		for i <= stratSz-j {
			points = append(points, *xR)
			indices = append(indices, i)

			k := strat[sidx]
			sidx++
			xR.Pow2k(&cparam, xR, 2*k)
			i += int(k)
		}

		cparam = phi.GenerateCurve(xR)
		for k := 0; k < len(points); k++ {
			points[k] = phi.EvaluatePoint(&points[k])
		}

		*phiP = phi.EvaluatePoint(phiP)
		*phiQ = phi.EvaluatePoint(phiQ)
		*phiR = phi.EvaluatePoint(phiR)

		// pop xR from points
		*xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}
}

// Traverses isogeny tree in order to compute xR needed
// for public key generation.
func traverseTreeSharedKeyA(curve *ProjectiveCurveParameters, xR *ProjectivePoint, pub *PublicKey) {
	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var i, sidx int

	cparam := curve.CalcCurveParamsEquiv4()
	phi := NewIsogeny4()
	strat := pub.params.A.IsogenyStrategy
	stratSz := len(strat)

	for j := 1; j <= stratSz; j++ {
		for i <= stratSz-j {
			points = append(points, *xR)
			indices = append(indices, i)

			k := strat[sidx]
			sidx++
			xR.Pow2k(&cparam, xR, 2*k)
			i += int(k)
		}

		cparam = phi.GenerateCurve(xR)
		for k := 0; k < len(points); k++ {
			points[k] = phi.EvaluatePoint(&points[k])
		}

		// pop xR from points
		*xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}
}

// Traverses isogeny tree in order to compute xR, xP, xQ and xQmP needed
// for public key generation.
func traverseTreePublicKeyB(curve *ProjectiveCurveParameters, xR, phiP, phiQ, phiR *ProjectivePoint, pub *PublicKey) {
	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var i, sidx int

	cparam := curve.CalcCurveParamsEquiv3()
	phi := NewIsogeny3()
	strat := pub.params.B.IsogenyStrategy
	stratSz := len(strat)

	for j := 1; j <= stratSz; j++ {
		for i <= stratSz-j {
			points = append(points, *xR)
			indices = append(indices, i)

			k := strat[sidx]
			sidx++
			xR.Pow3k(&cparam, xR, k)
			i += int(k)
		}

		cparam = phi.GenerateCurve(xR)
		for k := 0; k < len(points); k++ {
			points[k] = phi.EvaluatePoint(&points[k])
		}

		*phiP = phi.EvaluatePoint(phiP)
		*phiQ = phi.EvaluatePoint(phiQ)
		*phiR = phi.EvaluatePoint(phiR)

		// pop xR from points
		*xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}
}

// Traverses isogeny tree in order to compute xR, xP, xQ and xQmP needed
// for public key generation.
func traverseTreeSharedKeyB(curve *ProjectiveCurveParameters, xR *ProjectivePoint, pub *PublicKey) {
	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var i, sidx int

	cparam := curve.CalcCurveParamsEquiv3()
	phi := NewIsogeny3()
	strat := pub.params.B.IsogenyStrategy
	stratSz := len(strat)

	for j := 1; j <= stratSz; j++ {
		for i <= stratSz-j {
			points = append(points, *xR)
			indices = append(indices, i)

			k := strat[sidx]
			sidx++
			xR.Pow3k(&cparam, xR, k)
			i += int(k)
		}

		cparam = phi.GenerateCurve(xR)
		for k := 0; k < len(points); k++ {
			points[k] = phi.EvaluatePoint(&points[k])
		}

		// pop xR from points
		*xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}
}

// -----------------------------------------------------------------------------
// Key generation functions
//

// Generate a private key for "Alice".  Note that because this library does not
// implement SIDH validation, each keypair must be used for at most one
// shared secret computation.
func (prv *PrivateKey) generatePrivateKeyA(rand io.Reader) error {
	_, err := io.ReadFull(rand, prv.Scalar)
	if err != nil {
		return err
	}

	// Bit-twiddle to ensure scalar is in 2*[0,2^371):
	prv.Scalar[prv.params.SecretKeySize-1] = prv.params.A.MaskBytes[0]
	prv.Scalar[prv.params.SecretKeySize-2] &= prv.params.A.MaskBytes[1] // clear high bits, so scalar < 2^372
	prv.Scalar[0] &= prv.params.A.MaskBytes[2]                          // clear low bit, so scalar is even

	// We actually want scalar in 2*(0,2^371), but the above procedure
	// generates 0 with probability 2^(-371), which isn't worth checking
	// for.
	return nil
}

// Generate a private key for "Bob".  Note that because this library does not
// implement SIDH validation, each keypair must be used for at most one
// shared secret computation.
func (prv *PrivateKey) generatePrivateKeyB(rand io.Reader) error {
	// Perform rejection sampling to obtain a random value in [0,3^238]:
	var ok uint8
	for i := uint(0); i < prv.params.SampleRate; i++ {
		_, err := io.ReadFull(rand, prv.Scalar)
		if err != nil {
			return err
		}
		// Mask the high bits to obtain a uniform value in [0,2^378):
		// TODO: simply run it in loop, if rand distribution is uniform you surelly get non 0
		//       if not - better die, keep looping, hang, whatever, but don't generate secure key
		prv.Scalar[prv.params.SecretKeySize-1] &= prv.params.B.MaskBytes[0]

		// Accept if scalar < 3^238 (this happens w/ prob ~0.5828)
		// TODO this is specific to P751
		ok = checkLessThanThree238(prv.Scalar)
		if ok == 0 {
			break
		}
	}
	// ok is nonzero if all sampleRate trials failed.
	// This happens with probability 0.41719...^102 < 2^(-128), i.e., never
	if ok != 0 {
		// In case this happens user should retry. In practice it is highly
		// improbable (< 2^-128).
		return errors.New("sidh: private key generation failed")
	}

	// Multiply by 3 to get a scalar in 3*[0,3^238):
	multiplyByThree(prv.Scalar)
	// We actually want scalar in 2*(0,2^371), but the above procedure
	// generates 0 with probability 2^(-371), which isn't worth checking
	// for.
	return nil
}

// Generate a public key in the 2-torsion group
func publicKeyGenA(prv *PrivateKey) (pub *PublicKey) {
	var xPA, xQA, xRA ProjectivePoint
	var xPB, xQB, xRB, xR ProjectivePoint
	var invZP, invZQ, invZR ExtensionFieldElement
	var tmp ProjectiveCurveParameters
	var phi = NewIsogeny4()
	pub = NewPublicKey(prv.params.Id, KeyVariant_SIDH_A)

	// Load points for A
	xPA.FromAffine(&prv.params.A.Affine_P)
	xPA.Z.One()
	xQA.FromAffine(&prv.params.A.Affine_Q)
	xQA.Z.One()
	xRA.FromAffine(&prv.params.A.Affine_R)
	xRA.Z.One()

	// Load points for B
	xRB.FromAffine(&prv.params.B.Affine_R)
	xRB.Z.One()
	xQB.FromAffine(&prv.params.B.Affine_Q)
	xQB.Z.One()
	xPB.FromAffine(&prv.params.B.Affine_P)
	xPB.Z.One()

	// Find isogeny kernel
	tmp.A.Zero()
	tmp.C.One()
	xR = RightToLeftLadder(&tmp, &xPA, &xQA, &xRA, prv.params.A.SecretBitLen, prv.Scalar)

	// Reset params object and travers isogeny tree
	tmp.A.Zero()
	tmp.C.One()
	traverseTreePublicKeyA(&tmp, &xR, &xPB, &xQB, &xRB, pub)

	// Secret isogeny
	phi.GenerateCurve(&xR)
	xPA = phi.EvaluatePoint(&xPB)
	xQA = phi.EvaluatePoint(&xQB)
	xRA = phi.EvaluatePoint(&xRB)
	ExtensionFieldBatch3Inv(&xPA.Z, &xQA.Z, &xRA.Z, &invZP, &invZQ, &invZR)

	pub.affine_xP.Mul(&xPA.X, &invZP)
	pub.affine_xQ.Mul(&xQA.X, &invZQ)
	pub.affine_xQmP.Mul(&xRA.X, &invZR)
	return
}

// Generate a public key in the 3-torsion group
func publicKeyGenB(prv *PrivateKey) (pub *PublicKey) {
	var xPB, xQB, xRB, xR ProjectivePoint
	var xPA, xQA, xRA ProjectivePoint
	var invZP, invZQ, invZR ExtensionFieldElement
	var tmp ProjectiveCurveParameters
	var phi = NewIsogeny3()
	pub = NewPublicKey(prv.params.Id, prv.keyVariant)

	// Load points for B
	xRB.FromAffine(&prv.params.B.Affine_R)
	xRB.Z.One()
	xQB.FromAffine(&prv.params.B.Affine_Q)
	xQB.Z.One()
	xPB.FromAffine(&prv.params.B.Affine_P)
	xPB.Z.One()

	// Load points for A
	xPA.FromAffine(&prv.params.A.Affine_P)
	xPA.Z.One()
	xQA.FromAffine(&prv.params.A.Affine_Q)
	xQA.Z.One()
	xRA.FromAffine(&prv.params.A.Affine_R)
	xRA.Z.One()

	tmp.A.Zero()
	tmp.C.One()
	xR = RightToLeftLadder(&tmp, &xPB, &xQB, &xRB, prv.params.B.SecretBitLen, prv.Scalar)

	tmp.A.Zero()
	tmp.C.One()
	traverseTreePublicKeyB(&tmp, &xR, &xPA, &xQA, &xRA, pub)

	phi.GenerateCurve(&xR)
	xPB = phi.EvaluatePoint(&xPA)
	xQB = phi.EvaluatePoint(&xQA)
	xRB = phi.EvaluatePoint(&xRA)
	ExtensionFieldBatch3Inv(&xPB.Z, &xQB.Z, &xRB.Z, &invZP, &invZQ, &invZR)

	pub.affine_xP.Mul(&xPB.X, &invZP)
	pub.affine_xQ.Mul(&xQB.X, &invZQ)
	pub.affine_xQmP.Mul(&xRB.X, &invZR)
	return
}

// -----------------------------------------------------------------------------
// Key agreement functions
//

// Establishing shared keys in in 2-torsion group
func deriveSecretA(prv *PrivateKey, pub *PublicKey) []byte {
	var sharedSecret = make([]byte, pub.params.SharedSecretSize)
	var cparam ProjectiveCurveParameters
	var xP, xQ, xQmP ProjectivePoint
	var xR ProjectivePoint
	var phi = NewIsogeny4()

	// Recover curve coefficients
	cparam.RecoverCoordinateA(&pub.affine_xP, &pub.affine_xQ, &pub.affine_xQmP)
	cparam.C.One()

	// Find kernel of the morphism
	xP.FromAffine(&pub.affine_xP)
	xQ.FromAffine(&pub.affine_xQ)
	xQmP.FromAffine(&pub.affine_xQmP)
	xR = RightToLeftLadder(&cparam, &xP, &xQ, &xQmP, pub.params.A.SecretBitLen, prv.Scalar)

	// Traverse isogeny tree
	traverseTreeSharedKeyA(&cparam, &xR, pub)

	// Calculate j-invariant on isogeneus curve
	c := phi.GenerateCurve(&xR)
	cparam.RecoverCurveCoefficients4(&c)
	cparam.Jinvariant(sharedSecret)
	return sharedSecret
}

// Establishing shared keys in in 3-torsion group
func deriveSecretB(prv *PrivateKey, pub *PublicKey) []byte {
	var sharedSecret = make([]byte, pub.params.SharedSecretSize)
	var xP, xQ, xQmP ProjectivePoint
	var xR ProjectivePoint
	var cparam ProjectiveCurveParameters
	var phi = NewIsogeny3()

	// Recover curve coefficients
	cparam.RecoverCoordinateA(&pub.affine_xP, &pub.affine_xQ, &pub.affine_xQmP)
	cparam.C.One()

	// Find kernel of the morphism
	xP.FromAffine(&pub.affine_xP)
	xQ.FromAffine(&pub.affine_xQ)
	xQmP.FromAffine(&pub.affine_xQmP)
	xR = RightToLeftLadder(&cparam, &xP, &xQ, &xQmP, pub.params.B.SecretBitLen, prv.Scalar)

	// Traverse isogeny tree
	traverseTreeSharedKeyB(&cparam, &xR, pub)

	// Calculate j-invariant on isogeneus curve
	c := phi.GenerateCurve(&xR)
	cparam.RecoverCurveCoefficients3(&c)
	cparam.Jinvariant(sharedSecret)
	return sharedSecret
}
