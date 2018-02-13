package p751toolbox

// SIDH implementation parameters correspond to p751

// The x-coordinate of P_A = [3^239](11, oddsqrt(11^3 + 11)) on E_0(F_p)
var Affine_xPA = PrimeFieldElement{A: Fp751Element{0xd56fe52627914862, 0x1fad60dc96b5baea, 0x1e137d0bf07ab91, 0x404d3e9252161964, 0x3c5385e4cd09a337, 0x4476426769e4af73, 0x9790c6db989dfe33, 0xe06e1c04d2aa8b5e, 0x38c08185edea73b9, 0xaa41f678a4396ca6, 0x92b9259b2229e9a0, 0x2f9326818be0}}

// The y-coordinate of P_A = [3^239](11, oddsqrt(11^3 + 11)) on E_0(F_p)
var Affine_yPA = PrimeFieldElement{A: Fp751Element{0x332bd16fbe3d7739, 0x7e5e20ff2319e3db, 0xea856234aefbd81b, 0xe016df7d6d071283, 0x8ae42796f73cd34f, 0x6364b408a4774575, 0xa71c97f17ce99497, 0xda03cdd9aa0cbe71, 0xe52b4fda195bd56f, 0xdac41f811fce0a46, 0x9333720f0ee84a61, 0x1399f006e578}}

// The x-coordinate of P_B = [2^372](6, oddsqrt(6^3 + 6)) on E_0(F_p)
var Affine_xPB = PrimeFieldElement{A: Fp751Element{0xf1a8c9ed7b96c4ab, 0x299429da5178486e, 0xef4926f20cd5c2f4, 0x683b2e2858b4716a, 0xdda2fbcc3cac3eeb, 0xec055f9f3a600460, 0xd5a5a17a58c3848b, 0x4652d836f42eaed5, 0x2f2e71ed78b3a3b3, 0xa771c057180add1d, 0xc780a5d2d835f512, 0x114ea3b55ac1}}

// The y-coordinate of P_B = [2^372](6, oddsqrt(6^3 + 6)) on E_0(F_p)
var Affine_yPB = PrimeFieldElement{A: Fp751Element{0xd1e1471273e3736b, 0xf9301ba94da241fe, 0xe14ab3c17fef0a85, 0xb4ddd26a037e9e62, 0x66142dfb2afeb69, 0xe297cb70649d6c9e, 0x214dfc6e8b1a0912, 0x9f5ba818b01cf859, 0x87d15b4907c12828, 0xa4da70c53a880dbf, 0xac5df62a72c8f253, 0x2e26a42ec617}}

// The value of (a+2)/4 for the starting curve E_0 with a=0: this is 1/2
var E0_aPlus2Over4 = PrimeFieldElement{A: Fp751Element{0x124d6, 0x0, 0x0, 0x0, 0x0, 0xb8e0000000000000, 0x9c8a2434c0aa7287, 0xa206996ca9a378a3, 0x6876280d41a41b52, 0xe903b49f175ce04f, 0xf8511860666d227, 0x4ea07cff6e7f}}

const (
	// The secret key size, in bytes.
	SecretKeySize = 48
	// The public key size, in bytes.
	PublicKeySize = 564
	// The shared secret size, in bytes.
	SharedSecretSize = 188
	// Fixed parameters for isogeny tree computations
	MaxAlice = 185

	MaxBob = 239
	// Alice's mask values
	MaskAliceByte1 = 0

	MaskAliceByte2 = 15

	MaskAliceByte3 = 254
	// Bob's mask value
	MaskBobByte = 3

	// Sample rate to obtain a value in [0,3^238]
	SampleRate = 102
)

var AliceIsogenyStrategy = [MaxAlice]int{0, 1, 1, 2, 2, 2, 3, 4, 4, 4, 4, 5, 5,
	6, 7, 8, 8, 9, 9, 9, 9, 9, 9, 9, 12, 11, 12, 12, 13, 14, 15, 16, 16, 16, 16,
	16, 16, 17, 17, 18, 18, 17, 21, 17, 18, 21, 20, 21, 21, 21, 21, 21, 22, 25, 25,
	25, 26, 27, 28, 28, 29, 30, 31, 32, 32, 32, 32, 32, 32, 32, 33, 33, 33, 35, 36,
	36, 33, 36, 35, 36, 36, 35, 36, 36, 37, 38, 38, 39, 40, 41, 42, 38, 39, 40, 41,
	42, 40, 46, 42, 43, 46, 46, 46, 46, 48, 48, 48, 48, 49, 49, 48, 53, 54, 51, 52,
	53, 54, 55, 56, 57, 58, 59, 59, 60, 62, 62, 63, 64, 64, 64, 64, 64, 64, 64, 64,
	65, 65, 65, 65, 65, 66, 67, 65, 66, 67, 66, 69, 70, 66, 67, 66, 69, 70, 69, 70,
	70, 71, 72, 71, 72, 72, 74, 74, 75, 72, 72, 74, 74, 75, 72, 72, 74, 75, 75, 72,
	72, 74, 75, 75, 77, 77, 79, 80, 80, 82}

var BobIsogenyStrategy = [MaxBob]int{0, 1, 1, 2, 2, 2, 3, 3, 4, 4, 4, 5, 5, 5, 6,
	7, 8, 8, 8, 8, 9, 9, 9, 9, 9, 10, 12, 12, 12, 12, 12, 12, 13, 14, 14, 15, 16,
	16, 16, 16, 16, 17, 16, 16, 17, 19, 19, 20, 21, 22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 24, 24, 25, 27, 27, 28, 28, 29, 28, 29, 28, 28, 28, 30, 28, 28, 28, 29,
	30, 33, 33, 33, 33, 34, 35, 37, 37, 37, 37, 38, 38, 37, 38, 38, 38, 38, 38, 39,
	43, 38, 38, 38, 38, 43, 40, 41, 42, 43, 48, 45, 46, 47, 47, 48, 49, 49, 49, 50,
	51, 50, 49, 49, 49, 49, 51, 49, 53, 50, 51, 50, 51, 51, 51, 52, 55, 55, 55, 56,
	56, 56, 56, 56, 58, 58, 61, 61, 61, 63, 63, 63, 64, 65, 65, 65, 65, 66, 66, 65,
	65, 66, 66, 66, 66, 66, 66, 66, 71, 66, 73, 66, 66, 71, 66, 73, 66, 66, 71, 66,
	73, 68, 68, 71, 71, 73, 73, 73, 75, 75, 78, 78, 78, 80, 80, 80, 81, 81, 82, 83,
	84, 85, 86, 86, 86, 86, 86, 87, 86, 88, 86, 86, 86, 86, 88, 86, 88, 86, 86, 86,
	88, 88, 86, 86, 86, 93, 90, 90, 92, 92, 92, 93, 93, 93, 93, 93, 97, 97, 97, 97,
	97, 97}

// SIKE implementation parameters correspond to p751
const (
	// MessageSize in bytes.
	MessageSize = 32
	// CiphertextSize in bytes = SIDHPublicKeySize + MessageSize
	CiphertextSize = 596
	// SIKESecretKeySize in bytes = MessageSize + SIDHPublicKeySize + SIDHSecretKeySize
	SIKESecretKeySize = 644
	// G is a custom vallue for cSHAKE256
	G = "0"
	// H is a custom vallue for cSHAKE256
	H = "1"
	// P is a custom vallue for cSHAKE256
	P = "2"
)
