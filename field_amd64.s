#include "textflag.h"

// p751 + 1
#define P751P1_5   $0xEEB0000000000000
#define P751P1_6   $0xE3EC968549F878A8
#define P751P1_7   $0xDA959B1A13F7CC76
#define P751P1_8   $0x084E9867D6EBE876
#define P751P1_9   $0x8562B5045CB25748
#define P751P1_10  $0x0E12909F97BADC66
#define P751P1_11  $0x00006FE5D541F71C

#define P751_0     $0xFFFFFFFFFFFFFFFF
#define P751_5     $0xEEAFFFFFFFFFFFFF
#define P751_6     $0xE3EC968549F878A8
#define P751_7     $0xDA959B1A13F7CC76
#define P751_8     $0x084E9867D6EBE876
#define P751_9     $0x8562B5045CB25748
#define P751_10    $0x0E12909F97BADC66
#define P751_11    $0x00006FE5D541F71C

#define P751X2_0   $0xFFFFFFFFFFFFFFFE
#define P751X2_1   $0xFFFFFFFFFFFFFFFF
#define P751X2_5   $0xDD5FFFFFFFFFFFFF
#define P751X2_6   $0xC7D92D0A93F0F151
#define P751X2_7   $0xB52B363427EF98ED
#define P751X2_8   $0x109D30CFADD7D0ED
#define P751X2_9   $0x0AC56A08B964AE90
#define P751X2_10  $0x1C25213F2F75B8CD
#define P751X2_11  $0x0000DFCBAA83EE38

// The MSR code uses these registers for parameter passing.  Keep using
// them to avoid significant code changes.  This means that when the Go
// assembler does something strange, we can diff the machine code
// against a different assembler to find out what Go did.

#define REG_P1 DI
#define REG_P2 SI
#define REG_P3 DX

// We can't write MOVQ $0, AX because Go's assembler incorrectly
// optimizes this to XOR AX, AX, which clobbers the carry flags.
//
// This bug was defined to be "correct" behaviour (cf.
// https://github.com/golang/go/issues/12405 ) by declaring that the MOV
// pseudo-instruction clobbers flags, although this fact is mentioned
// nowhere in the documentation for the Go assembler.
//
// Defining MOVQ to clobber flags has the effect that it is never safe
// to interleave MOVQ with ADCQ and SBBQ instructions.  Since this is
// required to write a carry chain longer than registers' working set,
// all of the below code therefore relies on the unspecified and
// undocumented behaviour that MOV won't clobber flags, except in the
// case of the above-mentioned bug.
//
// However, there's also no specification of which instructions
// correspond to machine instructions, and which are
// pseudo-instructions (i.e., no specification of what the assembler
// actually does), so this doesn't seem much worse than usual.
//
// Avoid the bug by dropping the bytes for `mov eax, 0` in directly:

#define ZERO_AX_WITHOUT_CLOBBERING_FLAGS BYTE	$0xB8; BYTE $0; BYTE $0; BYTE $0; BYTE $0;

TEXT ·Fp751Add(SB), NOSPLIT, $0-24

	MOVQ	z+0(FP), REG_P3
	MOVQ	x+8(FP), REG_P1
	MOVQ	y+16(FP), REG_P2

	MOVQ	(REG_P1), R8
	MOVQ	(8)(REG_P1), R9
	MOVQ	(16)(REG_P1), R10
	MOVQ	(24)(REG_P1), R11
	MOVQ	(32)(REG_P1), R12
	MOVQ	(40)(REG_P1), R13
	MOVQ	(48)(REG_P1), R14
	MOVQ	(56)(REG_P1), R15
	MOVQ	(64)(REG_P1), CX
	ADDQ	(REG_P2), R8
	ADCQ	(8)(REG_P2), R9
	ADCQ	(16)(REG_P2), R10
	ADCQ	(24)(REG_P2), R11
	ADCQ	(32)(REG_P2), R12
	ADCQ	(40)(REG_P2), R13
	ADCQ	(48)(REG_P2), R14
	ADCQ	(56)(REG_P2), R15
	ADCQ	(64)(REG_P2), CX
	MOVQ	(72)(REG_P1), AX
	ADCQ	(72)(REG_P2), AX
	MOVQ	AX, (72)(REG_P3)
	MOVQ	(80)(REG_P1), AX
	ADCQ	(80)(REG_P2), AX
	MOVQ	AX, (80)(REG_P3)
	MOVQ	(88)(REG_P1), AX
	ADCQ	(88)(REG_P2), AX
	MOVQ	AX, (88)(REG_P3)

	MOVQ	P751X2_0, AX
	SUBQ	AX, R8
	MOVQ	P751X2_1, AX
	SBBQ	AX, R9
	SBBQ	AX, R10
	SBBQ	AX, R11
	SBBQ	AX, R12
	MOVQ	P751X2_5, AX
	SBBQ	AX, R13
	MOVQ	P751X2_6, AX
	SBBQ	AX, R14
	MOVQ	P751X2_7, AX
	SBBQ	AX, R15
	MOVQ	P751X2_8, AX
	SBBQ	AX, CX
	MOVQ	R8, (REG_P3)
	MOVQ	R9, (8)(REG_P3)
	MOVQ	R10, (16)(REG_P3)
	MOVQ	R11, (24)(REG_P3)
	MOVQ	R12, (32)(REG_P3)
	MOVQ	R13, (40)(REG_P3)
	MOVQ	R14, (48)(REG_P3)
	MOVQ	R15, (56)(REG_P3)
	MOVQ	CX, (64)(REG_P3)
	MOVQ	(72)(REG_P3), R8
	MOVQ	(80)(REG_P3), R9
	MOVQ	(88)(REG_P3), R10
	MOVQ	P751X2_9, AX
	SBBQ	AX, R8
	MOVQ	P751X2_10, AX
	SBBQ	AX, R9
	MOVQ	P751X2_11, AX
	SBBQ	AX, R10
	MOVQ	R8, (72)(REG_P3)
	MOVQ	R9, (80)(REG_P3)
	MOVQ	R10, (88)(REG_P3)
	ZERO_AX_WITHOUT_CLOBBERING_FLAGS
	SBBQ	$0, AX

	MOVQ	P751X2_0, SI
	ANDQ	AX, SI
	MOVQ	P751X2_1, R8
	ANDQ	AX, R8
	MOVQ	P751X2_5, R9
	ANDQ	AX, R9
	MOVQ	P751X2_6, R10
	ANDQ	AX, R10
	MOVQ	P751X2_7, R11
	ANDQ	AX, R11
	MOVQ	P751X2_8, R12
	ANDQ	AX, R12
	MOVQ	P751X2_9, R13
	ANDQ	AX, R13
	MOVQ	P751X2_10, R14
	ANDQ	AX, R14
	MOVQ	P751X2_11, R15
	ANDQ	AX, R15

	MOVQ	(REG_P3), AX
	ADDQ	SI, AX
	MOVQ	AX, (REG_P3)
	MOVQ	(8)(REG_P3), AX
	ADCQ	R8, AX
	MOVQ	AX, (8)(REG_P3)
	MOVQ	(16)(REG_P3), AX
	ADCQ	R8, AX
	MOVQ	AX, (16)(REG_P3)
	MOVQ	(24)(REG_P3), AX
	ADCQ	R8, AX
	MOVQ	AX, (24)(REG_P3)
	MOVQ	(32)(REG_P3), AX
	ADCQ	R8, AX
	MOVQ	AX, (32)(REG_P3)
	MOVQ	(40)(REG_P3), AX
	ADCQ	R9, AX
	MOVQ	AX, (40)(REG_P3)
	MOVQ	(48)(REG_P3), AX
	ADCQ	R10, AX
	MOVQ	AX, (48)(REG_P3)
	MOVQ	(56)(REG_P3), AX
	ADCQ	R11, AX
	MOVQ	AX, (56)(REG_P3)
	MOVQ	(64)(REG_P3), AX
	ADCQ	R12, AX
	MOVQ	AX, (64)(REG_P3)
	MOVQ	(72)(REG_P3), AX
	ADCQ	R13, AX
	MOVQ	AX, (72)(REG_P3)
	MOVQ	(80)(REG_P3), AX
	ADCQ	R14, AX
	MOVQ	AX, (80)(REG_P3)
	MOVQ	(88)(REG_P3), AX
	ADCQ	R15, AX
	MOVQ	AX, (88)(REG_P3)

	RET

TEXT ·Fp751Sub(SB), NOSPLIT, $0-24

	MOVQ	z+0(FP),  REG_P3
	MOVQ	x+8(FP),  REG_P1
	MOVQ	y+16(FP),  REG_P2

	MOVQ	(REG_P1), R8
	MOVQ	(8)(REG_P1), R9
	MOVQ	(16)(REG_P1), R10
	MOVQ	(24)(REG_P1), R11
	MOVQ	(32)(REG_P1), R12
	MOVQ	(40)(REG_P1), R13
	MOVQ	(48)(REG_P1), R14
	MOVQ	(56)(REG_P1), R15
	MOVQ	(64)(REG_P1), CX
	SUBQ	(REG_P2), R8
	SBBQ	(8)(REG_P2), R9
	SBBQ	(16)(REG_P2), R10
	SBBQ	(24)(REG_P2), R11
	SBBQ	(32)(REG_P2), R12
	SBBQ	(40)(REG_P2), R13
	SBBQ	(48)(REG_P2), R14
	SBBQ	(56)(REG_P2), R15
	SBBQ	(64)(REG_P2), CX
	MOVQ	R8, (REG_P3)
	MOVQ	R9, (8)(REG_P3)
	MOVQ	R10, (16)(REG_P3)
	MOVQ	R11, (24)(REG_P3)
	MOVQ	R12, (32)(REG_P3)
	MOVQ	R13, (40)(REG_P3)
	MOVQ	R14, (48)(REG_P3)
	MOVQ	R15, (56)(REG_P3)
	MOVQ	CX, (64)(REG_P3)
	MOVQ	(72)(REG_P1), AX
	SBBQ	(72)(REG_P2), AX
	MOVQ	AX, (72)(REG_P3)
	MOVQ	(80)(REG_P1), AX
	SBBQ	(80)(REG_P2), AX
	MOVQ	AX, (80)(REG_P3)
	MOVQ	(88)(REG_P1), AX
	SBBQ	(88)(REG_P2), AX
	MOVQ	AX, (88)(REG_P3)
	ZERO_AX_WITHOUT_CLOBBERING_FLAGS
	SBBQ	$0, AX

	MOVQ	P751X2_0, SI
	ANDQ	AX, SI
	MOVQ	P751X2_1, R8
	ANDQ	AX, R8
	MOVQ	P751X2_5, R9
	ANDQ	AX, R9
	MOVQ	P751X2_6, R10
	ANDQ	AX, R10
	MOVQ	P751X2_7, R11
	ANDQ	AX, R11
	MOVQ	P751X2_8, R12
	ANDQ	AX, R12
	MOVQ	P751X2_9, R13
	ANDQ	AX, R13
	MOVQ	P751X2_10, R14
	ANDQ	AX, R14
	MOVQ	P751X2_11, R15
	ANDQ	AX, R15

	MOVQ	(REG_P3), AX
	ADDQ	SI, AX
	MOVQ	AX, (REG_P3)
	MOVQ	(8)(REG_P3), AX
	ADCQ	R8, AX
	MOVQ	AX, (8)(REG_P3)
	MOVQ	(16)(REG_P3), AX
	ADCQ	R8, AX
	MOVQ	AX, (16)(REG_P3)
	MOVQ	(24)(REG_P3), AX
	ADCQ	R8, AX
	MOVQ	AX, (24)(REG_P3)
	MOVQ	(32)(REG_P3), AX
	ADCQ	R8, AX
	MOVQ	AX, (32)(REG_P3)
	MOVQ	(40)(REG_P3), AX
	ADCQ	R9, AX
	MOVQ	AX, (40)(REG_P3)
	MOVQ	(48)(REG_P3), AX
	ADCQ	R10, AX
	MOVQ	AX, (48)(REG_P3)
	MOVQ	(56)(REG_P3), AX
	ADCQ	R11, AX
	MOVQ	AX, (56)(REG_P3)
	MOVQ	(64)(REG_P3), AX
	ADCQ	R12, AX
	MOVQ	AX, (64)(REG_P3)
	MOVQ	(72)(REG_P3), AX
	ADCQ	R13, AX
	MOVQ	AX, (72)(REG_P3)
	MOVQ	(80)(REG_P3), AX
	ADCQ	R14, AX
	MOVQ	AX, (80)(REG_P3)
	MOVQ	(88)(REG_P3), AX
	ADCQ	R15, AX
	MOVQ	AX, (88)(REG_P3)

	RET

