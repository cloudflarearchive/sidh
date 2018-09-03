package sidh

import (
	"errors"
	. "github.com/cloudflare/p751sidh/internal/isogeny"
	"io"
)

// Id's correspond to bitlength of the prime field characteristic
// Currently FP_751 is the only one supported by this implementation
const (
	FP_503 PrimeFieldId = iota
	FP_751
	FP_964
	maxPrimeFieldId
)

const (
	// First 2 bits identify SIDH variant third bit indicates
	// wether key is a SIKE variant (set) or SIDH (not set)

	// 001 - SIDH: corresponds to 2-torsion group
	KeyVariant_SIDH_A KeyVariant = 1 << 0
	// 010 - SIDH: corresponds to 3-torsion group
	KeyVariant_SIDH_B = 1 << 1
	// 110 - SIKE
	KeyVariant_SIKE = 1<<2 | KeyVariant_SIDH_B
)

// Base type for public and private key. Used mainly to carry domain
// parameters.
type key struct {
	// Domain parameters of the algorithm to be used with a key
	params *SidhParams
	// Flag indicates wether corresponds to 2-, 3-torsion group or SIKE
	keyVariant KeyVariant
}

// Defines operations on public key
type PublicKey struct {
	key
	affine_xP   Fp2Element
	affine_xQ   Fp2Element
	affine_xQmP Fp2Element
}

// Defines operations on private key
type PrivateKey struct {
	key
	// Secret key
	Scalar []byte
	// Used only by KEM
	S []byte
}

// Accessor to the domain parameters
func (key *key) Params() *SidhParams {
	return key.params
}

// Accessor to key variant
func (key *key) Variant() KeyVariant {
	return key.keyVariant
}

// NewPrivateKey initializes private key.
// Usage of this function guarantees that the object is correctly initialized.
func NewPrivateKey(id PrimeFieldId, v KeyVariant) *PrivateKey {
	prv := &PrivateKey{key: key{params: Params(id), keyVariant: v}}
	prv.Scalar = make([]byte, prv.params.SecretKeySize)
	if v == KeyVariant_SIKE {
		prv.S = make([]byte, prv.params.MsgLen)
	}
	return prv
}

// NewPublicKey initializes public key.
// Usage of this function guarantees that the object is correctly initialized.
func NewPublicKey(id PrimeFieldId, v KeyVariant) *PublicKey {
	return &PublicKey{key: key{params: Params(id), keyVariant: v}}
}

// Import clears content of the public key currently stored in the structure
// and imports key stored in the byte string. Returns error in case byte string
// size is wrong. Doesn't perform any validation.
func (pub *PublicKey) Import(input []byte) error {
	if len(input) != pub.Size() {
		return errors.New("sidh: input to short")
	}
	op := CurveOperations{Params: pub.params}
	ssSz := pub.params.SharedSecretSize
	op.Fp2FromBytes(&pub.affine_xP, input[0:ssSz])
	op.Fp2FromBytes(&pub.affine_xQ, input[ssSz:2*ssSz])
	op.Fp2FromBytes(&pub.affine_xQmP, input[2*ssSz:3*ssSz])
	return nil
}

// Exports currently stored key. In case structure hasn't been filled with key data
// returned byte string is filled with zeros.
func (pub *PublicKey) Export() []byte {
	output := make([]byte, pub.params.PublicKeySize)
	op := CurveOperations{Params: pub.params}
	ssSz := pub.params.SharedSecretSize
	op.Fp2ToBytes(output[0:ssSz], &pub.affine_xP)
	op.Fp2ToBytes(output[ssSz:2*ssSz], &pub.affine_xQ)
	op.Fp2ToBytes(output[2*ssSz:3*ssSz], &pub.affine_xQmP)
	return output
}

// Size returns size of the public key in bytes
func (pub *PublicKey) Size() int {
	return pub.params.PublicKeySize
}

// Exports currently stored key. In case structure hasn't been filled with key data
// returned byte string is filled with zeros.
func (prv *PrivateKey) Export() []byte {
	ret := make([]byte, len(prv.Scalar)+len(prv.S))
	copy(ret, prv.S)
	copy(ret[len(prv.S):], prv.Scalar)
	return ret
}

// Size returns size of the private key in bytes
func (prv *PrivateKey) Size() int {
	tmp := prv.params.SecretKeySize
	if prv.Variant() == KeyVariant_SIKE {
		tmp += int(prv.params.MsgLen)
	}
	return tmp
}

// Import clears content of the private key currently stored in the structure
// and imports key from octet string. In case of SIKE, the random value 'S'
// must be prepended to the value of actual private key (see SIKE spec for details).
// Function doesn't import public key value to PrivateKey object.
func (prv *PrivateKey) Import(input []byte) error {
	if len(input) != prv.Size() {
		return errors.New("sidh: input to short")
	}
	if len(prv.Scalar) != prv.params.SecretKeySize {
		return errors.New("sidh: object wrongly initialized")
	}
	copy(prv.S, input[:len(prv.S)])
	copy(prv.Scalar, input[len(prv.S):])
	return nil
}

// Generates random private key for SIDH or SIKE. Returns error
// in case user provided RNG or memory initialization fails.
func (prv *PrivateKey) Generate(rand io.Reader) error {
	var err error

	if (prv.keyVariant & KeyVariant_SIDH_A) == KeyVariant_SIDH_A {
		err = prv.generatePrivateKeyA(rand)
	} else {
		err = prv.generatePrivateKeyB(rand)
	}

	if prv.keyVariant == KeyVariant_SIKE && err == nil {
		_, err = io.ReadFull(rand, prv.S)
	}

	return err
}

// Generates public key.
//
// Constant time.
func (prv *PrivateKey) GeneratePublicKey() (*PublicKey) {
	if (prv.keyVariant & KeyVariant_SIDH_A) == KeyVariant_SIDH_A {
		return publicKeyGenA(prv)
	}
	return publicKeyGenB(prv)
}

// Computes a shared secret which is a j-invariant. Function requires that pub has
// different KeyVariant than prv. Length of returned output is 2*ceil(log_2 P)/8),
// where P is a prime defining finite field.
//
// It's important to notice that each keypair must not be used more than once
// to calculate shared secret.
//
// Function may return error. This happens only in case provided input is invalid.
// Constant time for properly initialized private and public key.
func DeriveSecret(prv *PrivateKey, pub *PublicKey) ([]byte, error) {

	if (pub == nil) || (prv == nil) {
		return nil, errors.New("sidh: invalid arguments")
	}

	if (pub.keyVariant == prv.keyVariant) || (pub.params.Id != prv.params.Id) {
		return nil, errors.New("sidh: public and private are incompatbile")
	}

	if (prv.keyVariant & KeyVariant_SIDH_A) == KeyVariant_SIDH_A {
		return deriveSecretA(prv, pub), nil
	} else {
		return deriveSecretB(prv, pub), nil
	}
}
