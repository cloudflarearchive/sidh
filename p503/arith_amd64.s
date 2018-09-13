// +build amd64,!noasm

#include "textflag.h"

// p503
#define P503_0     $0xFFFFFFFFFFFFFFFF
#define P503_1     $0xFFFFFFFFFFFFFFFF
#define P503_2     $0xFFFFFFFFFFFFFFFF
#define P503_3     $0xABFFFFFFFFFFFFFF
#define P503_4     $0x13085BDA2211E7A0
#define P503_5     $0x1B9BF6C87B7E7DAF
#define P503_6     $0x6045C6BDDA77A4D0
#define P503_7     $0x004066F541811E1E

// p503+1
#define P503P1_3   $0xAC00000000000000
#define P503P1_4   $0x13085BDA2211E7A0
#define P503P1_5   $0x1B9BF6C87B7E7DAF
#define P503P1_6   $0x6045C6BDDA77A4D0
#define P503P1_7   $0x004066F541811E1E

// p503x2
#define P503X2_0   $0xFFFFFFFFFFFFFFFE
#define P503X2_1   $0xFFFFFFFFFFFFFFFF
#define P503X2_2   $0xFFFFFFFFFFFFFFFF
#define P503X2_3   $0x57FFFFFFFFFFFFFF
#define P503X2_4   $0x2610B7B44423CF41
#define P503X2_5   $0x3737ED90F6FCFB5E
#define P503X2_6   $0xC08B8D7BB4EF49A0
#define P503X2_7   $0x0080CDEA83023C3C

// The MSR code uses these registers for parameter passing.  Keep using
// them to avoid significant code changes.  This means that when the Go
// assembler does something strange, we can diff the machine code
// against a different assembler to find out what Go did.

#define REG_P1 DI
#define REG_P2 SI
#define REG_P3 DX

TEXT ·fp503StrongReduce(SB), NOSPLIT, $0-8
	MOVQ	x+0(FP), REG_P1

	// Zero AX for later use:
	XORQ	AX, AX

	// Load p into registers:
	MOVQ	P503_0, R8
	// P503_{1,2} = P503_0, so reuse R8
	MOVQ	P503_3, R9
	MOVQ	P503_4, R10
	MOVQ	P503_5, R11
	MOVQ	P503_6, R12
	MOVQ	P503_7, R13

	// Set x <- x - p
	SUBQ	R8,  ( 0)(REG_P1)
	SBBQ	R8,  ( 8)(REG_P1)
	SBBQ	R8,  (16)(REG_P1)
	SBBQ	R9,  (24)(REG_P1)
	SBBQ	R10, (32)(REG_P1)
	SBBQ	R11, (40)(REG_P1)
	SBBQ	R12, (48)(REG_P1)
	SBBQ	R13, (56)(REG_P1)

	// Save carry flag indicating x-p < 0 as a mask
	SBBQ	$0, AX

	// Conditionally add p to x if x-p < 0
	ANDQ	AX, R8
	ANDQ	AX, R9
	ANDQ	AX, R10
	ANDQ	AX, R11
	ANDQ	AX, R12
	ANDQ	AX, R13

	ADDQ	R8, ( 0)(REG_P1)
	ADCQ	R8, ( 8)(REG_P1)
	ADCQ	R8, (16)(REG_P1)
	ADCQ	R9, (24)(REG_P1)
	ADCQ	R10,(32)(REG_P1)
	ADCQ	R11,(40)(REG_P1)
	ADCQ	R12,(48)(REG_P1)
	ADCQ	R13,(56)(REG_P1)

	RET

// TODO this could be factorized to one single function
TEXT ·fp503ConditionalSwap(SB),NOSPLIT,$0-17

	MOVQ	x+0(FP), REG_P1
	MOVQ	y+8(FP), REG_P2
	MOVB	choice+16(FP), AL	// AL = 0 or 1
	MOVBLZX	AL, AX				// AX = 0 or 1
	NEGQ	AX					// RAX = 0x00..00 or 0xff..ff

#ifndef CSWAP_BLOCK
#define CSWAP_BLOCK(idx) 	\
	MOVQ	(idx*8)(REG_P1), BX	\ // BX = x[idx]
	MOVQ 	(idx*8)(REG_P2), CX	\ // CX = y[idx]
	MOVQ	CX, DX				\ // DX = y[idx]
	XORQ	BX, DX				\ // DX = y[idx] ^ x[idx]
	ANDQ	AX, DX				\ // DX = (y[idx] ^ x[idx]) & mask
	XORQ	DX, BX				\ // BX = (y[idx] ^ x[idx]) & mask) ^ x[idx] = x[idx] or y[idx]
	XORQ	DX, CX				\ // CX = (y[idx] ^ x[idx]) & mask) ^ y[idx] = y[idx] or x[idx]
	MOVQ	BX, (idx*8)(REG_P1)	\
	MOVQ	CX, (idx*8)(REG_P2)
#endif

	CSWAP_BLOCK(0)
	CSWAP_BLOCK(1)
	CSWAP_BLOCK(2)
	CSWAP_BLOCK(3)
	CSWAP_BLOCK(4)
	CSWAP_BLOCK(5)
	CSWAP_BLOCK(6)
	CSWAP_BLOCK(7)

#ifdef CSWAP_BLOCK
#undef CSWAP_BLOCK
#endif

	RET

TEXT ·fp503AddReduced(SB),NOSPLIT,$0-24

	MOVQ	z+0(FP), REG_P3
	MOVQ	x+8(FP), REG_P1
	MOVQ	y+16(FP), REG_P2

    // Used later to calculate a mask
    XORQ    CX, CX

    // [R8-R15]: z = x + y
	MOVQ	( 0)(REG_P1), R8
	MOVQ	( 8)(REG_P1), R9
	MOVQ	(16)(REG_P1), R10
	MOVQ	(24)(REG_P1), R11
	MOVQ	(32)(REG_P1), R12
	MOVQ	(40)(REG_P1), R13
	MOVQ	(48)(REG_P1), R14
	MOVQ	(56)(REG_P1), R15
	ADDQ	( 0)(REG_P2), R8
	ADCQ	( 8)(REG_P2), R9
	ADCQ	(16)(REG_P2), R10
	ADCQ	(24)(REG_P2), R11
	ADCQ	(32)(REG_P2), R12
	ADCQ	(40)(REG_P2), R13
	ADCQ	(48)(REG_P2), R14
	ADCQ	(56)(REG_P2), R15

    MOVQ    P503X2_0, AX
    SUBQ    AX, R8
    MOVQ    P503X2_1, AX
    SBBQ    AX, R9
    SBBQ    AX, R10
    MOVQ    P503X2_3, AX
    SBBQ    AX, R11
    MOVQ    P503X2_4, AX
    SBBQ    AX, R12
    MOVQ    P503X2_5, AX
    SBBQ    AX, R13
    MOVQ    P503X2_6, AX
    SBBQ    AX, R14
    MOVQ    P503X2_7, AX
    SBBQ    AX, R15

    SBBQ    $0, CX             // mask

                               // move z to REG_P3
    MOVQ    R8,  ( 0)(REG_P3)
    MOVQ    R9,  ( 8)(REG_P3)
    MOVQ    R10, (16)(REG_P3)
    MOVQ    R11, (24)(REG_P3)
    MOVQ    R12, (32)(REG_P3)
    MOVQ    R13, (40)(REG_P3)
    MOVQ    R14, (48)(REG_P3)
    MOVQ    R15, (56)(REG_P3)

                                   // if z<0 add p503x2 back
    MOVQ    P503X2_0,   R8
    MOVQ    P503X2_1,   R9
    MOVQ    P503X2_3,   R10
    MOVQ    P503X2_4,   R11
    MOVQ    P503X2_5,   R12
    MOVQ    P503X2_6,   R13
    MOVQ    P503X2_7,   R14
    ANDQ    CX, R8
    ANDQ    CX, R9
    ANDQ    CX, R10
    ANDQ    CX, R11
    ANDQ    CX, R12
    ANDQ    CX, R13
    ANDQ    CX, R14
    ADDQ    R8, ( 0)(REG_P3)
    ADCQ    R9, ( 8)(REG_P3)
    ADCQ    R9, (16)(REG_P3)
    ADCQ    R10,(24)(REG_P3)
    ADCQ    R11,(32)(REG_P3)
    ADCQ    R12,(40)(REG_P3)
    ADCQ    R13,(48)(REG_P3)
    ADCQ    R14,(56)(REG_P3)

	RET

TEXT ·fp503SubReduced(SB), NOSPLIT, $0-24

    MOVQ    z+0(FP), REG_P3
    MOVQ    x+8(FP), REG_P1
    MOVQ    y+16(FP), REG_P2

    // Used later to calculate a mask
    XORQ    CX, CX

    MOVQ    ( 0)(REG_P1), R8
    MOVQ    ( 8)(REG_P1), R9
    MOVQ    (16)(REG_P1), R10
    MOVQ    (24)(REG_P1), R11
    MOVQ    (32)(REG_P1), R12
    MOVQ    (40)(REG_P1), R13
    MOVQ    (48)(REG_P1), R14
    MOVQ    (56)(REG_P1), R15

    SUBQ    ( 0)(REG_P2), R8
    SBBQ    ( 8)(REG_P2), R9
    SBBQ    (16)(REG_P2), R10
    SBBQ    (24)(REG_P2), R11
    SBBQ    (32)(REG_P2), R12
    SBBQ    (40)(REG_P2), R13
    SBBQ    (48)(REG_P2), R14
    SBBQ    (56)(REG_P2), R15

    // mask
    SBBQ    $0, CX

    // store x-y in REG_P3
    MOVQ    R8,  ( 0)(REG_P3)
    MOVQ    R9,  ( 8)(REG_P3)
    MOVQ    R10, (16)(REG_P3)
    MOVQ    R11, (24)(REG_P3)
    MOVQ    R12, (32)(REG_P3)
    MOVQ    R13, (40)(REG_P3)
    MOVQ    R14, (48)(REG_P3)
    MOVQ    R15, (56)(REG_P3)

    // if z<0 add p503x2 back
    MOVQ    P503X2_0,   R8
    MOVQ    P503X2_1,   R9
    MOVQ    P503X2_3,   R10
    MOVQ    P503X2_4,   R11
    MOVQ    P503X2_5,   R12
    MOVQ    P503X2_6,   R13
    MOVQ    P503X2_7,   R14
    ANDQ    CX, R8
    ANDQ    CX, R9
    ANDQ    CX, R10
    ANDQ    CX, R11
    ANDQ    CX, R12
    ANDQ    CX, R13
    ANDQ    CX, R14

    ADDQ    R8, ( 0)(REG_P3)
    ADCQ    R9, ( 8)(REG_P3)
    ADCQ    R9, (16)(REG_P3)
    ADCQ    R10,(24)(REG_P3)
    ADCQ    R11,(32)(REG_P3)
    ADCQ    R12,(40)(REG_P3)
    ADCQ    R13,(48)(REG_P3)
    ADCQ    R14,(56)(REG_P3)

	RET

TEXT ·fp503Mul(SB), $96-24
	// Uses variant of Karatsuba method.
	//
	// Here we store the destination in CX instead of in REG_P3 because the
	// multiplication instructions use DX as an implicit destination
	// operand: MULQ $REG sets DX:AX <-- AX * $REG.

	// Actual implementation
	MOVQ	z+0(FP), CX
	MOVQ	x+8(FP), REG_P1
	MOVQ	y+16(FP), REG_P2

	// RAX and RDX will be used for a mask (0-borrow)
	XORQ	AX, AX

	// RCX[0-3]: U1+U0
	MOVQ	(32)(REG_P1), R8
	MOVQ	(40)(REG_P1), R9
	MOVQ	(48)(REG_P1), R10
	MOVQ	(56)(REG_P1), R11
	ADDQ	( 0)(REG_P1), R8
	ADCQ	( 8)(REG_P1), R9
	ADCQ	(16)(REG_P1), R10
	ADCQ	(24)(REG_P1), R11
	MOVQ	R8,  ( 0)(CX)
	MOVQ	R9,  ( 8)(CX)
	MOVQ	R10, (16)(CX)
	MOVQ	R11, (24)(CX)

	SBBQ	$0, AX

	// R12-R15: V1+V0
	XORQ	DX, DX
	MOVQ	(32)(REG_P2), R12
	MOVQ	(40)(REG_P2), R13
	MOVQ	(48)(REG_P2), R14
	MOVQ	(56)(REG_P2), R15
	ADDQ	( 0)(REG_P2), R12
	ADCQ	( 8)(REG_P2), R13
	ADCQ	(16)(REG_P2), R14
	ADCQ	(24)(REG_P2), R15

	SBBQ	$0, DX

	// Store carries on stack
	MOVQ	AX, (64)(SP)
	MOVQ	DX, (72)(SP)

	// (SP[0-3],R8,R9,R10,R11) <- (AH+AL)*(BH+BL).
	// MUL using comba; In comments below U=AH+AL V=BH+BL

    // U0*V0
    MOVQ    (CX), AX
    MULQ    R12
    MOVQ    AX, (SP)        // C0
    MOVQ    DX, R8

    // U0*V1
    XORQ    R9, R9
    MOVQ    (CX), AX
    MULQ    R13
    ADDQ    AX, R8
    ADCQ    DX, R9

    // U1*V0
    XORQ    R10, R10
    MOVQ    (8)(CX), AX
    MULQ    R12
    ADDQ    AX, R8
    MOVQ    R8, (8)(SP)     // C1
    ADCQ    DX, R9
    ADCQ    $0, R10

    // U0*V2
    XORQ    R8, R8
    MOVQ    (CX), AX
    MULQ    R14
    ADDQ    AX, R9
    ADCQ    DX, R10
    ADCQ    $0, R8

    // U2*V0
    MOVQ    (16)(CX), AX
    MULQ    R12
    ADDQ    AX, R9
    ADCQ    DX, R10
    ADCQ    $0, R8

    // U1*V1
    MOVQ    (8)(CX), AX
    MULQ    R13
    ADDQ    AX, R9
    MOVQ    R9, (16)(SP)        // C2
    ADCQ    DX, R10
    ADCQ    $0, R8

    // U0*V3
    XORQ    R9, R9
    MOVQ    (CX), AX
    MULQ    R15
    ADDQ    AX, R10
    ADCQ    DX, R8
    ADCQ    $0, R9

    // U3*V0
    MOVQ    (24)(CX), AX
    MULQ    R12
    ADDQ    AX, R10
    ADCQ    DX, R8
    ADCQ    $0, R9

    // U1*V2
    MOVQ    (8)(CX), AX
    MULQ    R14
    ADDQ    AX, R10
    ADCQ    DX, R8
    ADCQ    $0, R9

    // U2*V1
    MOVQ    (16)(CX), AX
    MULQ    R13
    ADDQ    AX, R10
    MOVQ    R10, (24)(SP)       // C3
    ADCQ    DX, R8
    ADCQ    $0, R9

    // U1*V3
    XORQ    R10, R10
    MOVQ    (8)(CX), AX
    MULQ    R15
    ADDQ    AX, R8
    ADCQ    DX, R9
    ADCQ    $0, R10

    // U3*V1
    MOVQ    (24)(CX), AX
    MULQ    R13
    ADDQ    AX, R8
    ADCQ    DX, R9
    ADCQ    $0, R10

    // U2*V2
    MOVQ    (16)(CX), AX
    MULQ    R14
    ADDQ    AX, R8
    MOVQ    R8, (32)(SP)        // C4
    ADCQ    DX, R9
    ADCQ    $0, R10

    // U2*V3
    XORQ    R11, R11
    MOVQ    (16)(CX), AX
    MULQ    R15
    ADDQ    AX, R9
    ADCQ    DX, R10
    ADCQ    $0, R11

    // U3*V2
    MOVQ    (24)(CX), AX
    MULQ    R14
    ADDQ    AX, R9              // C5
    ADCQ    DX, R10
    ADCQ    $0, R11

    // U3*V3
    MOVQ    (24)(CX), AX
    MULQ    R15
    ADDQ    AX, R10             // C6
    ADCQ    DX, R11             // C7


	MOVQ    (64)(SP), AX
	ANDQ    AX, R12
	ANDQ    AX, R13
	ANDQ    AX, R14
	ANDQ    AX, R15
	ADDQ    R8, R12
	ADCQ    R9, R13
	ADCQ    R10, R14
	ADCQ    R11, R15


	MOVQ    (72)(SP), AX
	MOVQ    (CX), R8
	MOVQ    (8)(CX), R9
	MOVQ    (16)(CX), R10
	MOVQ    (24)(CX), R11
	ANDQ    AX, R8
	ANDQ    AX, R9
	ANDQ    AX, R10
	ANDQ    AX, R11
	ADDQ    R12, R8
	ADCQ    R13, R9
	ADCQ    R14, R10
	ADCQ    R15, R11
	MOVQ    R8, (32)(SP)
	MOVQ    R9, (40)(SP)
	MOVQ    R10, (48)(SP)
	MOVQ    R11, (56)(SP)

	// CX[0-7] <- AL*BL

    // U0*V0
    MOVQ    (REG_P1), R11
    MOVQ    (REG_P2), AX
    MULQ    R11
    XORQ    R9, R9
    MOVQ    AX, (CX)            // C0
    MOVQ    DX, R8

    // U0*V1
    MOVQ    (16)(REG_P1), R14
    MOVQ    (8)(REG_P2), AX
    MULQ    R11
    XORQ    R10, R10
    ADDQ    AX, R8
    ADCQ    DX, R9

    // U1*V0
    MOVQ    (8)(REG_P1), R12
    MOVQ    (REG_P2), AX
    MULQ    R12
    ADDQ    AX, R8
    MOVQ    R8, (8)(CX)         // C1
    ADCQ    DX, R9
    ADCQ    $0, R10

    // U0*V2
    XORQ    R8, R8
    MOVQ    (16)(REG_P2), AX
    MULQ    R11
    ADDQ    AX, R9
    ADCQ    DX, R10
    ADCQ    $0, R8

    // U2*V0
    MOVQ    (REG_P2), R13
    MOVQ    R14, AX
    MULQ    R13
    ADDQ    AX, R9
    ADCQ    DX, R10
    ADCQ    $0, R8

    // U1*V1
    MOVQ    (8)(REG_P2), AX
    MULQ    R12
    ADDQ    AX, R9
    MOVQ    R9, (16)(CX)        // C2
    ADCQ    DX, R10
    ADCQ    $0, R8

    // U0*V3
    XORQ    R9, R9
    MOVQ    (24)(REG_P2), AX
    MULQ    R11
    MOVQ    (24)(REG_P1), R15
    ADDQ    AX, R10
    ADCQ    DX, R8
    ADCQ    $0, R9

    // U3*V1
    MOVQ    R15, AX
    MULQ    R13
    ADDQ    AX, R10
    ADCQ    DX, R8
    ADCQ    $0, R9

    // U2*V2
    MOVQ    (16)(REG_P2), AX
    MULQ    R12
    ADDQ    AX, R10
    ADCQ    DX, R8
    ADCQ    $0, R9

    // U2*V3
    MOVQ    (8)(REG_P2), AX
    MULQ    R14
    ADDQ    AX, R10
    MOVQ    R10, (24)(CX)       // C3
    ADCQ    DX, R8
    ADCQ    $0, R9

    // U3*V2
    XORQ    R10, R10
    MOVQ    (24)(REG_P2), AX
    MULQ    R12
    ADDQ    AX, R8
    ADCQ    DX, R9
    ADCQ    $0, R10

    // U3*V1
    MOVQ    (8)(REG_P2), AX
    MULQ    R15
    ADDQ    AX, R8
    ADCQ    DX, R9
    ADCQ    $0, R10

    // U2*V2
    MOVQ    (16)(REG_P2), AX
    MULQ    R14
    ADDQ    AX, R8
    MOVQ    R8, (32)(CX)		// C4
    ADCQ    DX, R9
    ADCQ    $0, R10

    // U2*V3
    XORQ    R8, R8
    MOVQ    (24)(REG_P2), AX
    MULQ    R14
    ADDQ    AX, R9
    ADCQ    DX, R10
    ADCQ    $0, R8

    // U3*V2
    MOVQ    (16)(REG_P2), AX
    MULQ    R15
    ADDQ    AX, R9
    MOVQ    R9, (40)(CX)		// C5
    ADCQ    DX, R10
    ADCQ    $0, R8

    // U3*V3
    MOVQ    (24)(REG_P2), AX
    MULQ    R15
    ADDQ    AX, R10
    MOVQ    R10, (48)(CX)		// C6
    ADCQ    DX, R8
    MOVQ    R8, (56)(CX)		// C7

	// CX[8-15] <- AH*BH
    MOVQ    (32)(REG_P1), R11
    MOVQ    (32)(REG_P2), AX
    MULQ    R11
    XORQ    R9, R9
    MOVQ    AX, (64)(CX)        // C0
    MOVQ    DX, R8

    MOVQ    (48)(REG_P1), R14
    MOVQ    (40)(REG_P2), AX
    MULQ    R11
    XORQ    R10, R10
    ADDQ    AX, R8
    ADCQ    DX, R9

    MOVQ    (40)(REG_P1), R12
    MOVQ    (32)(REG_P2), AX
    MULQ    R12
    ADDQ    AX, R8
    MOVQ    R8, (72)(CX)        // C1
    ADCQ    DX, R9
    ADCQ    $0, R10

    XORQ    R8, R8
    MOVQ    (48)(REG_P2), AX
    MULQ    R11
    ADDQ    AX, R9
    ADCQ    DX, R10
    ADCQ    $0, R8

    MOVQ    (32)(REG_P2), R13
    MOVQ    R14, AX
    MULQ    R13
    ADDQ    AX, R9
    ADCQ    DX, R10
    ADCQ    $0, R8

    MOVQ    (40)(REG_P2), AX
    MULQ    R12
    ADDQ    AX, R9
    MOVQ    R9, (80)(CX)        // C2
    ADCQ    DX, R10
    ADCQ    $0, R8

    XORQ    R9, R9
    MOVQ    (56)(REG_P2), AX
    MULQ    R11
    MOVQ    (56)(REG_P1), R15
    ADDQ    AX, R10
    ADCQ    DX, R8
    ADCQ    $0, R9

    MOVQ    R15, AX
    MULQ    R13
    ADDQ    AX, R10
    ADCQ    DX, R8
    ADCQ    $0, R9

    MOVQ    (48)(REG_P2), AX
    MULQ    R12
    ADDQ    AX, R10
    ADCQ    DX, R8
    ADCQ    $0, R9

    MOVQ    (40)(REG_P2), AX
    MULQ    R14
    ADDQ    AX, R10
    MOVQ    R10, (88)(CX)       // C3
    ADCQ    DX, R8
    ADCQ    $0, R9

    XORQ    R10, R10
    MOVQ    (56)(REG_P2), AX
    MULQ    R12
    ADDQ    AX, R8
    ADCQ    DX, R9
    ADCQ    $0, R10

    MOVQ    (40)(REG_P2), AX
    MULQ    R15
    ADDQ    AX, R8
    ADCQ    DX, R9
    ADCQ    $0, R10

    MOVQ    (48)(REG_P2), AX
    MULQ    R14
    ADDQ    AX, R8
    MOVQ    R8, (96)(CX)        // C4
    ADCQ    DX, R9
    ADCQ    $0, R10

    XORQ    R8, R8
    MOVQ    (56)(REG_P2), AX
    MULQ    R14
    ADDQ    AX, R9
    ADCQ    DX, R10
    ADCQ    $0, R8

    MOVQ    (48)(REG_P2), AX
    MULQ    R15
    ADDQ    AX, R9
    MOVQ    R9, (104)(CX)       // C5
    ADCQ    DX, R10
    ADCQ    $0, R8

    MOVQ    (56)(REG_P2), AX
    MULQ    R15
    ADDQ    AX, R10
    MOVQ    R10, (112)(CX)      // C6
    ADCQ    DX, R8
    MOVQ    R8, (120)(CX)       // C7

	// [R8-R15] <- (AH+AL)*(BH+BL) - AL*BL
	MOVQ    (SP), R8
	SUBQ    (CX), R8
	MOVQ    (8)(SP), R9
	SBBQ    (8)(CX), R9
	MOVQ    (16)(SP), R10
	SBBQ    (16)(CX), R10
	MOVQ    (24)(SP), R11
	SBBQ    (24)(CX), R11
	MOVQ    (32)(SP), R12
	SBBQ    (32)(CX), R12
	MOVQ    (40)(SP), R13
	SBBQ    (40)(CX), R13
	MOVQ    (48)(SP), R14
	SBBQ    (48)(CX), R14
	MOVQ    (56)(SP), R15
	SBBQ    (56)(CX), R15

	// [R8-R15] <- (AH+AL)*(BH+BL) - AL*BL - AH*BH
	MOVQ    (64)(CX), AX
	SUBQ    AX, R8
	MOVQ    (72)(CX), AX
	SBBQ    AX, R9
	MOVQ    (80)(CX), AX
	SBBQ    AX, R10
	MOVQ    (88)(CX), AX
	SBBQ    AX, R11
	MOVQ    (96)(CX), AX
	SBBQ    AX, R12
	MOVQ    (104)(CX), DX
	SBBQ    DX, R13
	MOVQ    (112)(CX), DI
	SBBQ    DI, R14
	MOVQ    (120)(CX), SI
	SBBQ    SI, R15

	// Final result
	ADDQ    (32)(CX), R8
	MOVQ    R8, (32)(CX)
	ADCQ    (40)(CX), R9
	MOVQ    R9, (40)(CX)
	ADCQ    (48)(CX), R10
	MOVQ    R10, (48)(CX)
	ADCQ    (56)(CX), R11
	MOVQ    R11, (56)(CX)
	ADCQ    (64)(CX), R12
	MOVQ    R12, (64)(CX)
	ADCQ    (72)(CX), R13
	MOVQ    R13, (72)(CX)
	ADCQ    (80)(CX), R14
	MOVQ    R14, (80)(CX)
	ADCQ    (88)(CX), R15
	MOVQ    R15, (88)(CX)
	ADCQ    $0, AX
	MOVQ    AX, (96)(CX)
	ADCQ    $0, DX
	MOVQ    DX, (104)(CX)
	ADCQ    $0, DI
	MOVQ    DI, (112)(CX)
	ADCQ    $0, SI
	MOVQ    SI, (120)(CX)

	RET

TEXT ·fp503MontgomeryReduce(SB), $0-16

	MOVQ	z+0(FP), REG_P2
	MOVQ	x+8(FP), REG_P1

	MOVQ    (REG_P1), R11
	MOVQ    P503P1_3, AX
	MULQ    R11
	XORQ    R8, R8
	ADDQ    (24)(REG_P1), AX
	MOVQ    AX, (24)(REG_P2)
	ADCQ    DX, R8


	XORQ    R9, R9
	MOVQ    P503P1_4, AX
	MULQ    R11
	XORQ    R10, R10
	ADDQ    AX, R8
	ADCQ    DX, R9


	MOVQ    (8)(REG_P1), R12
	MOVQ    P503P1_3, AX
	MULQ    R12
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10
	ADDQ    (32)(REG_P1), R8
	MOVQ    R8, (32)(REG_P2)       // Z4
	ADCQ    $0, R9
	ADCQ    $0, R10


	XORQ    R8, R8
	MOVQ    P503P1_5, AX
	MULQ    R11
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8


	MOVQ    P503P1_4, AX
	MULQ    R12
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8


	MOVQ    (16)(REG_P1), R13
	MOVQ    P503P1_3, AX
	MULQ    R13
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8
	ADDQ    (40)(REG_P1), R9
	MOVQ    R9, (40)(REG_P2)       // Z5
	ADCQ    $0, R10
	ADCQ    $0, R8


	XORQ    R9, R9
	MOVQ    P503P1_6, AX
	MULQ    R11
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9


	MOVQ    P503P1_5, AX
	MULQ    R12
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9


	MOVQ    P503P1_4, AX
	MULQ    R13
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9


	MOVQ    (24)(REG_P2), R14
	MOVQ    P503P1_3, AX
	MULQ    R14
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9
	ADDQ    (48)(REG_P1), R10
	MOVQ    R10, (48)(REG_P2)      // Z6
	ADCQ    $0, R8
	ADCQ    $0, R9


	XORQ    R10, R10
	MOVQ    P503P1_7, AX
	MULQ    R11
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10


	MOVQ    P503P1_6, AX
	MULQ    R12
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10


	MOVQ    P503P1_5, AX
	MULQ    R13
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10


	MOVQ    P503P1_4, AX
	MULQ    R14
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10


	MOVQ    (32)(REG_P2), R15
	MOVQ    P503P1_3, AX
	MULQ    R15
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10
	ADDQ    (56)(REG_P1), R8
	MOVQ    R8, (56)(REG_P2)       // Z7
	ADCQ    $0, R9
	ADCQ    $0, R10


	XORQ    R8, R8
	MOVQ    P503P1_7, AX
	MULQ    R12
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8


	MOVQ    P503P1_6, AX
	MULQ    R13
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8


	MOVQ    P503P1_5, AX
	MULQ    R14
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8


	MOVQ    P503P1_4, AX
	MULQ    R15
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8


	MOVQ    (40)(REG_P2), CX
	MOVQ    P503P1_3, AX
	MULQ    CX
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8
	ADDQ    (64)(REG_P1), R9
	MOVQ    R9, (REG_P2)           // Z0
	ADCQ    $0, R10
	ADCQ    $0, R8


	XORQ    R9, R9
	MOVQ    P503P1_7, AX
	MULQ    R13
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9


	MOVQ    P503P1_6, AX
	MULQ    R14
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9


	MOVQ    P503P1_5, AX
	MULQ    R15
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9


	MOVQ    P503P1_4, AX
	MULQ    CX
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9


	MOVQ    (48)(REG_P2), R13
	MOVQ    P503P1_3, AX
	MULQ    R13
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9
	ADDQ    (72)(REG_P1), R10
	MOVQ    R10, (8)(REG_P2)           // Z1
	ADCQ    $0, R8
	ADCQ    $0, R9


	XORQ    R10, R10
	MOVQ    P503P1_7, AX
	MULQ    R14
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10


	MOVQ    P503P1_6, AX
	MULQ    R15
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10


	MOVQ    P503P1_5, AX
	MULQ    CX
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10


	MOVQ    P503P1_4, AX
	MULQ    R13
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10


	MOVQ    (56)(REG_P2), R14
	MOVQ    P503P1_3, AX
	MULQ    R14
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10
	ADDQ    (80)(REG_P1), R8
	MOVQ    R8, (16)(REG_P2)           // Z2
	ADCQ    $0, R9
	ADCQ    $0, R10


	XORQ    R8, R8
	MOVQ    P503P1_7, AX
	MULQ    R15
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8


	MOVQ    P503P1_6, AX
	MULQ    CX
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8


	MOVQ    P503P1_5, AX
	MULQ    R13
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8


	MOVQ    P503P1_4, AX
	MULQ    R14
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADCQ    $0, R8
	ADDQ    (88)(REG_P1), R9
	MOVQ    R9, (24)(REG_P2)       // Z3
	ADCQ    $0, R10
	ADCQ    $0, R8


	XORQ    R9, R9
	MOVQ    P503P1_7, AX
	MULQ    CX
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9


	MOVQ    P503P1_6, AX
	MULQ    R13
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9


	MOVQ    P503P1_5, AX
	MULQ    R14
	ADDQ    AX, R10
	ADCQ    DX, R8
	ADCQ    $0, R9
	ADDQ    (96)(REG_P1), R10
	MOVQ    R10, (32)(REG_P2)          // Z4
	ADCQ    $0, R8
	ADCQ    $0, R9


	XORQ    R10, R10
	MOVQ    P503P1_7, AX
	MULQ    R13
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10


	MOVQ    P503P1_6, AX
	MULQ    R14
	ADDQ    AX, R8
	ADCQ    DX, R9
	ADCQ    $0, R10
	ADDQ    (104)(REG_P1), R8          // Z5
	MOVQ    R8, (40)(REG_P2)           // Z5
	ADCQ    $0, R9
	ADCQ    $0, R10


	MOVQ    P503P1_7, AX
	MULQ    R14
	ADDQ    AX, R9
	ADCQ    DX, R10
	ADDQ    (112)(REG_P1), R9      // Z6
	MOVQ    R9, (48)(REG_P2)       // Z6
	ADCQ    $0, R10
	ADDQ    (120)(REG_P1), R10     // Z7
	MOVQ    R10, (56)(REG_P2)      // Z7

	RET

TEXT ·fp503AddLazy(SB), NOSPLIT, $0-24

	MOVQ z+0(FP), REG_P3
	MOVQ x+8(FP), REG_P1
	MOVQ y+16(FP), REG_P2

	MOVQ	(REG_P1), R8
	MOVQ	(8)(REG_P1), R9
	MOVQ	(16)(REG_P1), R10
	MOVQ	(24)(REG_P1), R11
	MOVQ	(32)(REG_P1), R12
	MOVQ	(40)(REG_P1), R13
	MOVQ	(48)(REG_P1), R14
	MOVQ	(56)(REG_P1), R15

	ADDQ	(REG_P2), R8
	ADCQ	(8)(REG_P2), R9
	ADCQ	(16)(REG_P2), R10
	ADCQ	(24)(REG_P2), R11
	ADCQ	(32)(REG_P2), R12
	ADCQ	(40)(REG_P2), R13
	ADCQ	(48)(REG_P2), R14
	ADCQ	(56)(REG_P2), R15

	MOVQ	R8, (REG_P3)
	MOVQ	R9, (8)(REG_P3)
	MOVQ	R10, (16)(REG_P3)
	MOVQ	R11, (24)(REG_P3)
	MOVQ	R12, (32)(REG_P3)
	MOVQ	R13, (40)(REG_P3)
	MOVQ	R14, (48)(REG_P3)
	MOVQ	R15, (56)(REG_P3)

	RET

TEXT ·fp503X2AddLazy(SB), NOSPLIT, $0-24

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
	MOVQ	(64)(REG_P1), AX
	MOVQ	(72)(REG_P1), BX
	MOVQ	(80)(REG_P1), CX

	ADDQ	(REG_P2), R8
	ADCQ	(8)(REG_P2), R9
	ADCQ	(16)(REG_P2), R10
	ADCQ	(24)(REG_P2), R11
	ADCQ	(32)(REG_P2), R12
	ADCQ	(40)(REG_P2), R13
	ADCQ	(48)(REG_P2), R14
	ADCQ	(56)(REG_P2), R15
	ADCQ	(64)(REG_P2), AX
	ADCQ	(72)(REG_P2), BX
	ADCQ	(80)(REG_P2), CX

	MOVQ	R8, (REG_P3)
	MOVQ	R9, (8)(REG_P3)
	MOVQ	R10, (16)(REG_P3)
	MOVQ	R11, (24)(REG_P3)
	MOVQ	R12, (32)(REG_P3)
	MOVQ	R13, (40)(REG_P3)
	MOVQ	R14, (48)(REG_P3)
	MOVQ	R15, (56)(REG_P3)
	MOVQ	AX, (64)(REG_P3)
	MOVQ	BX, (72)(REG_P3)
	MOVQ	CX, (80)(REG_P3)

	MOVQ	(88)(REG_P1), R8
	MOVQ	(96)(REG_P1), R9
	MOVQ	(104)(REG_P1), R10
	MOVQ	(112)(REG_P1), R11
	MOVQ	(120)(REG_P1), R12

	ADCQ	(88)(REG_P2), R8
	ADCQ	(96)(REG_P2), R9
	ADCQ	(104)(REG_P2), R10
	ADCQ	(112)(REG_P2), R11
	ADCQ	(120)(REG_P2), R12

	MOVQ	R8, (88)(REG_P3)
	MOVQ	R9, (96)(REG_P3)
	MOVQ	R10, (104)(REG_P3)
	MOVQ	R11, (112)(REG_P3)
	MOVQ	R12, (120)(REG_P3)

	RET

TEXT ·fp503X2SubLazy(SB), NOSPLIT, $0-24

	MOVQ z+0(FP), REG_P3
	MOVQ x+8(FP), REG_P1
	MOVQ y+16(FP), REG_P2
	// Used later to store result of 0-borrow
	XORQ CX, CX

	// SUBC for first 11 limbs
	MOVQ	(REG_P1), R8
	MOVQ	(8)(REG_P1), R9
	MOVQ	(16)(REG_P1), R10
	MOVQ	(24)(REG_P1), R11
	MOVQ	(32)(REG_P1), R12
	MOVQ	(40)(REG_P1), R13
	MOVQ	(48)(REG_P1), R14
	MOVQ	(56)(REG_P1), R15
	MOVQ	(64)(REG_P1), AX
	MOVQ	(72)(REG_P1), BX

	SUBQ	(REG_P2), R8
	SBBQ	(8)(REG_P2), R9
	SBBQ	(16)(REG_P2), R10
	SBBQ	(24)(REG_P2), R11
	SBBQ	(32)(REG_P2), R12
	SBBQ	(40)(REG_P2), R13
	SBBQ	(48)(REG_P2), R14
	SBBQ	(56)(REG_P2), R15
	SBBQ	(64)(REG_P2), AX
	SBBQ	(72)(REG_P2), BX

	MOVQ	R8, (REG_P3)
	MOVQ	R9, (8)(REG_P3)
	MOVQ	R10, (16)(REG_P3)
	MOVQ	R11, (24)(REG_P3)
	MOVQ	R12, (32)(REG_P3)
	MOVQ	R13, (40)(REG_P3)
	MOVQ	R14, (48)(REG_P3)
	MOVQ	R15, (56)(REG_P3)
	MOVQ	AX, (64)(REG_P3)
	MOVQ	BX, (72)(REG_P3)

	// SUBC for last 5 limbs
	MOVQ	(80)(REG_P1), 	R8
	MOVQ	(88)(REG_P1), 	R9
	MOVQ	(96)(REG_P1), 	R10
	MOVQ	(104)(REG_P1), 	R11
	MOVQ	(112)(REG_P1), 	R12
	MOVQ	(120)(REG_P1), 	R13

	SBBQ	(80)(REG_P2), R8
	SBBQ	(88)(REG_P2), R9
	SBBQ	(96)(REG_P2), R10
	SBBQ	(104)(REG_P2), R11
	SBBQ	(112)(REG_P2), R12
	SBBQ	(120)(REG_P2), R13

	MOVQ	R8, (80)(REG_P3)
	MOVQ	R9, (88)(REG_P3)
	MOVQ	R10, (96)(REG_P3)
	MOVQ	R11, (104)(REG_P3)
	MOVQ	R12, (112)(REG_P3)
	MOVQ	R13, (120)(REG_P3)

	// Now the carry flag is 1 if x-y < 0.  If so, add p*2^768.
	SBBQ	$0, CX

	// Load p into registers:
	MOVQ	P503_0, R8
	// P503_{1,2} = P503_0, so reuse R8
	MOVQ	P503_3, R9
	MOVQ	P503_4, R10
	MOVQ	P503_5, R11
	MOVQ	P503_6, R12
	MOVQ	P503_7, R13

	ANDQ	CX, R8
	ANDQ	CX, R9
	ANDQ	CX, R10
	ANDQ	CX, R11
	ANDQ	CX, R12
	ANDQ	CX, R13

	ADDQ	R8,  (64   )(REG_P3)
	ADCQ	R8,  (64+ 8)(REG_P3)
	ADCQ	R8,  (64+16)(REG_P3)
	ADCQ	R9,  (64+24)(REG_P3)
	ADCQ	R10, (64+32)(REG_P3)
	ADCQ	R11, (64+40)(REG_P3)
	ADCQ	R12, (64+48)(REG_P3)
	ADCQ	R13, (64+56)(REG_P3)

	RET

