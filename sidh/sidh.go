package sidh

import (
	"errors"
	"io"

	. "github.com/cloudflare/p751sidh/internal/isogeny"
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
	var op = CurveOperations{Params: pub.params}

	cparam := op.CalcCurveParamsEquiv4(curve)
	phi := Newisogeny4(op.Params.Op)
	strat := pub.params.A.IsogenyStrategy
	stratSz := len(strat)

	for j := 1; j <= stratSz; j++ {
		for i <= stratSz-j {
			points = append(points, *xR)
			indices = append(indices, i)

			k := strat[sidx]
			sidx++
			op.Pow2k(xR, &cparam, 2*k)
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
	var op = CurveOperations{Params: pub.params}

	cparam := op.CalcCurveParamsEquiv4(curve)
	phi := Newisogeny4(op.Params.Op)
	strat := pub.params.A.IsogenyStrategy
	stratSz := len(strat)

	for j := 1; j <= stratSz; j++ {
		for i <= stratSz-j {
			points = append(points, *xR)
			indices = append(indices, i)

			k := strat[sidx]
			sidx++
			op.Pow2k(xR, &cparam, 2*k)
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
	var op = CurveOperations{Params: pub.params}

	cparam := op.CalcCurveParamsEquiv3(curve)
	phi := Newisogeny3(op.Params.Op)
	strat := pub.params.B.IsogenyStrategy
	stratSz := len(strat)

	for j := 1; j <= stratSz; j++ {
		for i <= stratSz-j {
			points = append(points, *xR)
			indices = append(indices, i)

			k := strat[sidx]
			sidx++
			op.Pow3k(xR, &cparam, k)
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
	var op = CurveOperations{Params: pub.params}

	cparam := op.CalcCurveParamsEquiv3(curve)
	phi := Newisogeny3(op.Params.Op)
	strat := pub.params.B.IsogenyStrategy
	stratSz := len(strat)

	for j := 1; j <= stratSz; j++ {
		for i <= stratSz-j {
			points = append(points, *xR)
			indices = append(indices, i)

			k := strat[sidx]
			sidx++
			op.Pow3k(xR, &cparam, k)
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
	var invZP, invZQ, invZR Fp2Element
	var tmp ProjectiveCurveParameters

	pub = NewPublicKey(prv.params.Id, KeyVariant_SIDH_A)
	var op = CurveOperations{Params: pub.params}
	var phi = Newisogeny4(op.Params.Op)

	// Load points for A
	xPA = ProjectivePoint{X: prv.params.A.Affine_P, Z: prv.params.OneFp2}
	xQA = ProjectivePoint{X: prv.params.A.Affine_Q, Z: prv.params.OneFp2}
	xRA = ProjectivePoint{X: prv.params.A.Affine_R, Z: prv.params.OneFp2}

	// Load points for B
	xRB = ProjectivePoint{X: prv.params.B.Affine_R, Z: prv.params.OneFp2}
	xQB = ProjectivePoint{X: prv.params.B.Affine_Q, Z: prv.params.OneFp2}
	xPB = ProjectivePoint{X: prv.params.B.Affine_P, Z: prv.params.OneFp2}

	// Find isogeny kernel
	tmp.C = pub.params.OneFp2
	xR = op.ScalarMul3Pt(&tmp, &xPA, &xQA, &xRA, prv.params.A.SecretBitLen, prv.Scalar)

	// Reset params object and travers isogeny tree
	tmp.C = pub.params.OneFp2
	tmp.A.Zeroize()
	traverseTreePublicKeyA(&tmp, &xR, &xPB, &xQB, &xRB, pub)

	// Secret isogeny
	phi.GenerateCurve(&xR)
	xPA = phi.EvaluatePoint(&xPB)
	xQA = phi.EvaluatePoint(&xQB)
	xRA = phi.EvaluatePoint(&xRB)
	op.Fp2Batch3Inv(&xPA.Z, &xQA.Z, &xRA.Z, &invZP, &invZQ, &invZR)

	op.Params.Op.Mul(&pub.affine_xP, &xPA.X, &invZP)
	op.Params.Op.Mul(&pub.affine_xQ, &xQA.X, &invZQ)
	op.Params.Op.Mul(&pub.affine_xQmP, &xRA.X, &invZR)
	return
}

// Generate a public key in the 3-torsion group
func publicKeyGenB(prv *PrivateKey) (pub *PublicKey) {
	var xPB, xQB, xRB, xR ProjectivePoint
	var xPA, xQA, xRA ProjectivePoint
	var invZP, invZQ, invZR Fp2Element
	var tmp ProjectiveCurveParameters

	pub = NewPublicKey(prv.params.Id, prv.keyVariant)
	var op = CurveOperations{Params: pub.params}
	var phi = Newisogeny3(op.Params.Op)

	// Load points for B
	xRB = ProjectivePoint{X: prv.params.B.Affine_R, Z: prv.params.OneFp2}
	xQB = ProjectivePoint{X: prv.params.B.Affine_Q, Z: prv.params.OneFp2}
	xPB = ProjectivePoint{X: prv.params.B.Affine_P, Z: prv.params.OneFp2}

	// Load points for A
	xPA = ProjectivePoint{X: prv.params.A.Affine_P, Z: prv.params.OneFp2}
	xQA = ProjectivePoint{X: prv.params.A.Affine_Q, Z: prv.params.OneFp2}
	xRA = ProjectivePoint{X: prv.params.A.Affine_R, Z: prv.params.OneFp2}

	tmp.C = pub.params.OneFp2
	xR = op.ScalarMul3Pt(&tmp, &xPB, &xQB, &xRB, prv.params.B.SecretBitLen, prv.Scalar)

	tmp.C = pub.params.OneFp2
	tmp.A.Zeroize()
	traverseTreePublicKeyB(&tmp, &xR, &xPA, &xQA, &xRA, pub)

	phi.GenerateCurve(&xR)
	xPB = phi.EvaluatePoint(&xPA)
	xQB = phi.EvaluatePoint(&xQA)
	xRB = phi.EvaluatePoint(&xRA)
	op.Fp2Batch3Inv(&xPB.Z, &xQB.Z, &xRB.Z, &invZP, &invZQ, &invZR)

	op.Params.Op.Mul(&pub.affine_xP, &xPB.X, &invZP)
	op.Params.Op.Mul(&pub.affine_xQ, &xQB.X, &invZQ)
	op.Params.Op.Mul(&pub.affine_xQmP, &xRB.X, &invZR)
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
	var op = CurveOperations{Params: prv.params}
	var phi = Newisogeny4(op.Params.Op)

	// Recover curve coefficients
	cparam.C = pub.params.OneFp2
	op.RecoverCoordinateA(&cparam, &pub.affine_xP, &pub.affine_xQ, &pub.affine_xQmP)

	// Find kernel of the morphism
	xP = ProjectivePoint{X: pub.affine_xP, Z: pub.params.OneFp2}
	xQ = ProjectivePoint{X: pub.affine_xQ, Z: pub.params.OneFp2}
	xQmP = ProjectivePoint{X: pub.affine_xQmP, Z: pub.params.OneFp2}
	xR = op.ScalarMul3Pt(&cparam, &xP, &xQ, &xQmP, pub.params.A.SecretBitLen, prv.Scalar)

	// Traverse isogeny tree
	traverseTreeSharedKeyA(&cparam, &xR, pub)

	// Calculate j-invariant on isogeneus curve
	c := phi.GenerateCurve(&xR)
	op.RecoverCurveCoefficients4(&cparam, &c)
	op.Jinvariant(&cparam, sharedSecret)
	return sharedSecret
}

// Establishing shared keys in in 3-torsion group
func deriveSecretB(prv *PrivateKey, pub *PublicKey) []byte {
	var sharedSecret = make([]byte, pub.params.SharedSecretSize)
	var xP, xQ, xQmP ProjectivePoint
	var xR ProjectivePoint
	var cparam ProjectiveCurveParameters
	var op = CurveOperations{Params: prv.params}
	var phi = Newisogeny3(op.Params.Op)

	// Recover curve coefficients
	cparam.C = pub.params.OneFp2
	op.RecoverCoordinateA(&cparam, &pub.affine_xP, &pub.affine_xQ, &pub.affine_xQmP)

	// Find kernel of the morphism
	xP = ProjectivePoint{X: pub.affine_xP, Z: pub.params.OneFp2}
	xQ = ProjectivePoint{X: pub.affine_xQ, Z: pub.params.OneFp2}
	xQmP = ProjectivePoint{X: pub.affine_xQmP, Z: pub.params.OneFp2}
	xR = op.ScalarMul3Pt(&cparam, &xP, &xQ, &xQmP, pub.params.B.SecretBitLen, prv.Scalar)

	// Traverse isogeny tree
	traverseTreeSharedKeyB(&cparam, &xR, pub)

	// Calculate j-invariant on isogeneus curve
	c := phi.GenerateCurve(&xR)
	op.RecoverCurveCoefficients3(&cparam, &c)
	op.Jinvariant(&cparam, sharedSecret)
	return sharedSecret
}
