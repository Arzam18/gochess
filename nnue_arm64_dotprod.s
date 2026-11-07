//go:build arm64 && gochess_dotprod

#include "textflag.h"

// ARMv8.2-A DotProd kernels. These functions are only linked into the explicit
// gochess_dotprod build. Keep the regular NNUE ARM64 assembly as the NEON path.

// nnueMatMul32x512_i8DotProd(input *byte, weightsT *int8, biases *int32, output *int32)
// input is CReLU-packed and therefore in [0,127], so it is valid as signed
// byte data for SDOT. Each SDOT accumulates 16 products into four int32 lanes.
TEXT ·nnueMatMul32x512_i8DotProd(SB), NOSPLIT, $0-32
	MOVD input+0(FP), R0
	MOVD weightsT+8(FP), R1
	MOVD biases+16(FP), R2
	MOVD output+24(FP), R3

	MOVD $8, R4                   // 8 groups of 4 output neurons
	MOVD $512, R5                 // row stride in bytes

i8dp_outer:
	MOVD R1, R6
	ADD R5, R1, R7
	ADD R5, R7, R8
	ADD R5, R8, R9

	VEOR V16.B16, V16.B16, V16.B16
	VEOR V17.B16, V17.B16, V17.B16
	VEOR V18.B16, V18.B16, V18.B16
	VEOR V19.B16, V19.B16, V19.B16

	MOVD R0, R10
	MOVD $32, R11                // 512 / 16 bytes

i8dp_inner:
	VLD1 (R10), [V24.B16]

	VLD1 (R6), [V25.B16]
	WORD $0x4E999710             // SDOT V16.4S, V24.16B, V25.16B
	VLD1 (R7), [V25.B16]
	WORD $0x4E999711             // SDOT V17.4S, V24.16B, V25.16B
	VLD1 (R8), [V25.B16]
	WORD $0x4E999712             // SDOT V18.4S, V24.16B, V25.16B
	VLD1 (R9), [V25.B16]
	WORD $0x4E999713             // SDOT V19.4S, V24.16B, V25.16B

	ADD $16, R10, R10
	ADD $16, R6, R6
	ADD $16, R7, R7
	ADD $16, R8, R8
	ADD $16, R9, R9
	SUBS $1, R11, R11
	BNE i8dp_inner

	// Sum the four int32 lanes and add the original bias.
	WORD $0x4EB0BE10             // ADDP V16.4S, V16.4S, V16.4S
	WORD $0x4EB0BE10             // ADDP V16.4S, V16.4S, V16.4S
	WORD $0x1E26020C             // FMOV W12, S16
	MOVW (R2), R13
	ADDW R13, R12, R12
	MOVW R12, (R3)

	WORD $0x4EB1BE31             // ADDP V17.4S, V17.4S, V17.4S
	WORD $0x4EB1BE31             // ADDP V17.4S, V17.4S, V17.4S
	WORD $0x1E26022C             // FMOV W12, S17
	MOVW 4(R2), R13
	ADDW R13, R12, R12
	MOVW R12, 4(R3)

	WORD $0x4EB2BE52             // ADDP V18.4S, V18.4S, V18.4S
	WORD $0x4EB2BE52             // ADDP V18.4S, V18.4S, V18.4S
	WORD $0x1E26024C             // FMOV W12, S18
	MOVW 8(R2), R13
	ADDW R13, R12, R12
	MOVW R12, 8(R3)

	WORD $0x4EB3BE73             // ADDP V19.4S, V19.4S, V19.4S
	WORD $0x4EB3BE73             // ADDP V19.4S, V19.4S, V19.4S
	WORD $0x1E26026C             // FMOV W12, S19
	MOVW 12(R2), R13
	ADDW R13, R12, R12
	MOVW R12, 12(R3)

	ADD $2048, R1, R1
	ADD $16, R2, R2
	ADD $16, R3, R3
	SUBS $1, R4, R4
	BNE i8dp_outer
	RET

// nnueV5L1Int8MatMulNDotProd(acc8 *byte, wT8 *int8, hidden *int32, accLen int, l1 int)
// The SCReLU pack is uint8 [0,255], while SDOT takes signed bytes. XOR with
// 0x80 maps x to signed (x-128). Add 128*sum(weights) after the dot product:
// sum((x-128)*w) + 128*sum(w) == sum(x*w), exactly.
TEXT ·nnueV5L1Int8MatMulNDotProd(SB), NOSPLIT, $0-40
	MOVD acc8+0(FP), R0
	MOVD wT8+8(FP), R1
	MOVD hidden+16(FP), R2
	MOVD accLen+24(FP), R3
	MOVD l1+32(FP), R4
	MOVD R3, R5                  // weight row stride
	LSR $4, R3, R3               // accLen / 16

	WORD $0x4F04E41A             // MOVI V26.16B, #0x80
l1i8dp_outer:
	VEOR V16.B16, V16.B16, V16.B16
	MOVD $0, R10                 // sum of signed weights for correction
	MOVD R0, R6
	MOVD R1, R7
	MOVD R3, R8

l1i8dp_inner:
	VLD1 (R6), [V2.B16]          // uint8 inputs
	VLD1 (R7), [V3.B16]          // int8 weights
	WORD $0x6E3A1C42             // EOR V2.16B, V2.16B, V26.16B (input - 128)
	WORD $0x4E839450             // SDOT V16.4S, V2.16B, V3.16B
	WORD $0x4E303871             // SADDLV H17, V3.16B
	WORD $0x4E022E2B             // SMOV X11, V17.H[0]
	ADD R11, R10, R10            // accumulate sum(weights)

	ADD $16, R6, R6
	ADD $16, R7, R7
	SUBS $1, R8, R8
	BNE l1i8dp_inner

	// Reduce four SDOT lanes, then restore the unsigned-input offset.
	WORD $0x4EB0BE10             // ADDP V16.4S, V16.4S, V16.4S
	WORD $0x4EB0BE10             // ADDP V16.4S, V16.4S, V16.4S
	WORD $0x1E26020C             // FMOV W12, S16
	LSL $7, R10, R11             // 128 * sum(weights)
	ADDW R11, R12, R12

	MOVW (R2), R9
	ADDW R12, R9, R9
	MOVW R9, (R2)

	ADD $4, R2, R2
	ADD R5, R1, R1
	SUBS $1, R4, R4
	BNE l1i8dp_outer
	RET
