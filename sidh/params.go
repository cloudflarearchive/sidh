package sidh

import . "github.com/cloudflare/p751sidh/p751toolbox"

type DomainParams struct {
	// P, Q and R=P-Q base points
	Affine_P, Affine_Q, Affine_R ExtensionFieldElement
	// Max size of secret key for x-torsion group
	SecretBitLen uint
	// MaskBytes
	MaskBytes []byte
	// Size of a compuatation strategy for x-torsion group
	IsogenyStrategy []uint32
}

type SidhParams struct {
	Id PrimeFieldId
	// The secret key size, in bytes.
	SecretKeySize int
	// The public key size, in bytes.
	PublicKeySize int
	// The shared secret size, in bytes.
	SharedSecretSize uint
	// 2- and 3-torsion group parameter definitions
	A, B DomainParams
	// Sample rate to obtain a value in [0,3^238]
	SampleRate uint
}

// Keeps mapping: SIDH prime field ID to domain parameters
var sidhParams = make(map[PrimeFieldId]SidhParams)

// Params returns domain parameters corresponding to finite field and identified by
// `id` provieded by the caller. Function panics in case `id` wasn't registered earlier.
func Params(id PrimeFieldId) *SidhParams {
	if val, ok := sidhParams[id]; ok {
		return &val
	}
	panic("sidh: SIDH Params ID unregistered")
}

func init() {
	p751 := SidhParams{
		Id:               FP_751,
		SecretKeySize:    P751_SecretKeySize,
		PublicKeySize:    P751_PublicKeySize,
		SharedSecretSize: P751_SharedSecretSize,
		A: DomainParams{
			Affine_P:        P751_affine_PA,
			Affine_Q:        P751_affine_QA,
			Affine_R:        P751_affine_RA,
			SecretBitLen:    P751_SecretBitLenA,
			MaskBytes:       []byte{P751_MaskAliceByte1, P751_MaskAliceByte2, P751_MaskAliceByte3},
			IsogenyStrategy: P751_AliceIsogenyStrategy[:],
		},
		B: DomainParams{
			Affine_P:        P751_affine_PB,
			Affine_Q:        P751_affine_QB,
			Affine_R:        P751_affine_RB,
			SecretBitLen:    P751_SecretBitLenB,
			MaskBytes:       []byte{P751_MaskBobByte},
			IsogenyStrategy: P751_BobIsogenyStrategy[:],
		},
		SampleRate: P751_SampleRate,
	}

	sidhParams[FP_751] = p751
}
