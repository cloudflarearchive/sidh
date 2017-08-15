package cln16sidh

import (
	"testing"
)

// Test the first four-isogeny from the base curve E_0(F_{p^2})
func TestFirstFourIsogenyVersusSage(t *testing.T) {
	var xR, isogenized_xR, sageIsogenized_xR ProjectivePoint

	// sage: p = 2^372 * 3^239 - 1; Fp = GF(p)
	//   ***   Warning: increasing stack size to 2000000.
	//   ***   Warning: increasing stack size to 4000000.
	// sage: R.<x> = Fp[]
	// sage: Fp2 = Fp.extension(x^2 + 1, 'i')
	// sage: i = Fp2.gen()
	// sage: E0Fp = EllipticCurve(Fp, [0,0,0,1,0])
	// sage: E0Fp2 = EllipticCurve(Fp2, [0,0,0,1,0])
	// sage: x_PA = 11
	// sage: y_PA = -Fp(11^3 + 11).sqrt()
	// sage: x_PB = 6
	// sage: y_PB = -Fp(6^3 + 6).sqrt()
	// sage: P_A = 3^239 * E0Fp((x_PA,y_PA))
	// sage: P_B = 2^372 * E0Fp((x_PB,y_PB))
	// sage: def tau(P):
	// ....:     return E0Fp2( (-P.xy()[0], i*P.xy()[1]))
	// ....:
	// sage: m_B = 3*randint(0,3^238)
	// sage: m_A = 2*randint(0,2^371)
	// sage: R_A = E0Fp2(P_A) + m_A*tau(P_A)
	// sage: def y_recover(x, a):
	// ....:     return (x**3 + a*x**2 + x).sqrt()
	// ....:
	// sage: first_4_torsion_point = E0Fp2(1, y_recover(Fp2(1),0))
	// sage: sage_first_4_isogeny = E0Fp2.isogeny(first_4_torsion_point)
	// sage: a = Fp2(0)
	// sage: sage_isomorphism = sage_first_4_isogeny.codomain().isomorphism_to(EllipticCurve(Fp2, [0,(2*(a+6))/(a-2),0,1,0]))
	// sage: isogenized_R_A = sage_isomorphism(sage_first_4_isogeny(R_A))

	xR.fromAffine(&ExtensionFieldElement{a: fp751Element{0xa179cb7e2a95fce9, 0xbfd6a0f3a0a892c0, 0x8b2f0aa4250ab3f3, 0x2e7aa4dd4118732d, 0x627969e493acbc2a, 0x21a5b852c7b8cc83, 0x26084278586324f2, 0x383be1aa5aa947c0, 0xc6558ecbb5c0183e, 0xf1f192086a52b035, 0x4c58b755b865c1b, 0x67b4ceea2d2c}, b: fp751Element{0xfceb02a2797fecbf, 0x3fee9e1d21f95e99, 0xa1c4ce896024e166, 0xc09c024254517358, 0xf0255994b17b94e7, 0xa4834359b41ee894, 0x9487f7db7ebefbe, 0x3bbeeb34a0bf1f24, 0xfa7e5533514c6a05, 0x92b0328146450a9a, 0xfde71ca3fada4c06, 0x3610f995c2bd}})

	sageIsogenized_xR.fromAffine(&ExtensionFieldElement{a: fp751Element{0xff99e76f78da1e05, 0xdaa36bd2bb8d97c4, 0xb4328cee0a409daf, 0xc28b099980c5da3f, 0xf2d7cd15cfebb852, 0x1935103dded6cdef, 0xade81528de1429c3, 0x6775b0fa90a64319, 0x25f89817ee52485d, 0x706e2d00848e697, 0xc4958ec4216d65c0, 0xc519681417f}, b: fp751Element{0x742fe7dde60e1fb9, 0x801a3c78466a456b, 0xa9f945b786f48c35, 0x20ce89e1b144348f, 0xf633970b7776217e, 0x4c6077a9b38976e5, 0x34a513fc766c7825, 0xacccba359b9cd65, 0xd0ca8383f0fd0125, 0x77350437196287a, 0x9fe1ad7706d4ea21, 0x4d26129ee42d}})

	var a ExtensionFieldElement
	a.Zero()

	_, phi := ComputeFirstFourIsogeny(&a)

	isogenized_xR = phi.Eval(&xR)

	if !sageIsogenized_xR.VartimeEq(&isogenized_xR) {
		t.Error("\nExpected\n", sageIsogenized_xR.toAffine(), "\nfound\n", isogenized_xR.toAffine())
	}
}
