package sidh

import (
	p751 "github.com/cloudflare/p751sidh/p751"
	. "github.com/cloudflare/p751sidh/internal/isogeny"
)

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
		SecretKeySize:    p751.P751_SecretKeySize,
		PublicKeySize:    p751.P751_PublicKeySize,
		SharedSecretSize: p751.P751_SharedSecretSize,
		A: DomainParams{
			Affine_P:        p751.P751_affine_PA,
			Affine_Q:        p751.P751_affine_QA,
			Affine_R:        p751.P751_affine_RA,
			SecretBitLen:    p751.P751_SecretBitLenA,
			MaskBytes:       []byte{p751.P751_MaskAliceByte1, p751.P751_MaskAliceByte2, p751.P751_MaskAliceByte3},
			IsogenyStrategy: p751.P751_AliceIsogenyStrategy[:],
		},
		B: DomainParams{
			Affine_P:        p751.P751_affine_PB,
			Affine_Q:        p751.P751_affine_QB,
			Affine_R:        p751.P751_affine_RB,
			SecretBitLen:    p751.P751_SecretBitLenB,
			MaskBytes:       []byte{p751.P751_MaskBobByte},
			IsogenyStrategy: p751.P751_BobIsogenyStrategy[:],
		},
		OneFp2:  p751.P751_OneFp2,
		HalfFp2: p751.P751_HalfFp2,
		MsgLen: 32,
		// SIKEp751 provides 192 bit of classical security ([SIKE], 5.1)
		KemSize:    24,
		SampleRate: p751.P751_SampleRate,
		Bytelen: p751.P751_Bytelen,
		Op: p751.FieldOperations(),
	}

	sidhParams[FP_751] = p751
}
