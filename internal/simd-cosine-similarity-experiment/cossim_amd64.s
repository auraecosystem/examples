//go:build amd64

#include "textflag.h"

// cosSimScalar is a scalar (single float at a time) implementation of
// cosine similarity, for testing & comparison.
TEXT ·cosSimScalar(SB), NOSPLIT, $0
	// Load slice headers;
	// SI: address of a's data
	// DI: address of b's data
	// CX: length of a (total number of elements)
	MOVQ a+0(FP), SI
	MOVQ b+24(FP), DI
	MOVQ a+8(FP), CX

	// X0 accumulates a*b (dot product between a and b), X6 accumulates a*a,
	// X7 accumulates b*b
	XORPS X0, X0
	XORPS X6, X6
	XORPS X7, X7

	// for r8 := 0; r8 < len(a); r8++
	MOVQ $0, R8

Loop:
	CMPQ R8, CX
	JAE Done

	// r8 counts the floats, so to load bytes we scale r8*4
	MOVSS (SI)(R8*4), X1
	MOVSS (DI)(R8*4), X2

	// Efficiently sequence the instructions so they can be pipelined.
	MOVUPS X1, X5
	MULSS X2, X1
	MULSS X2, X2
	MULSS X5, X5
	ADDSS X1, X0
	ADDSS X2, X6
	ADDSS X5, X7

	INCQ R8
	JMP Loop

Done:
	// a*b / sqrt(a*a * b*b)
	MULSS X7, X6
	RSQRTSS X6, X6
	MULSS X6, X0

	MOVSS X0, ret+48(FP)
	RET

// cosSimAVX2 uses AVX2 instructions to calculate the cosine similarity between
// two float32 slices.
// For simplicity, this function assumes that len(a)==len(b) and len(a)%8==0
TEXT ·cosSimAVX2(SB), NOSPLIT, $0
	// Load slice headers;
	// SI: address of a's data
	// DI: address of b's data
	// CX: length of a (total number of elements)
	MOVQ a+0(FP), SI
	MOVQ b+24(FP), DI
	MOVQ a+8(FP), CX
	
	// Y0 accumulates a*b, Y6 accumulates a*a, Y7 accumulates b*b
	VXORPS Y0, Y0, Y0
	VXORPS Y6, Y6, Y6
	VXORPS Y7, Y7, Y7

	// for r8 := 0; r8 < len(a); r8 += 8
	MOVQ $0, R8

Loop:
	CMPQ R8, CX
	JAE Done

	// r8 counts the floats, so to load bytes we scale r8*4
	VMOVUPS 0(SI)(R8*4), Y1
	VMOVUPS 0(DI)(R8*4), Y2

	VMULPS Y1, Y2, Y3
	VMULPS Y1, Y1, Y4
	VMULPS Y2, Y2, Y5

	VADDPS Y3, Y0, Y0
	VADDPS Y4, Y6, Y6
	VADDPS Y5, Y7, Y7

	ADDQ $8, R8
	JMP Loop

Done:
	// Horizontally reduce each of a*b, a*a, a*b separately into its register
	VPERM2F128 $0b00000001, Y0, Y0, Y1
	VADDPS Y1, Y0, Y0
	VHADDPS Y0, Y0, Y0
	VHADDPS Y0, Y0, Y0

	VPERM2F128 $0b00000001, Y6, Y6, Y8
	VADDPS Y8, Y6, Y6
	VHADDPS Y6, Y6, Y6
	VHADDPS Y6, Y6, Y6

	VPERM2F128 $0b00000001, Y7, Y7, Y8
	VADDPS Y8, Y7, Y7
	VHADDPS Y7, Y7, Y7
	VHADDPS Y7, Y7, Y7

	// a*b / sqrt(a*a * b*b)
	MULSS X7, X6
	RSQRTSS X6, X6
	MULSS X6, X0

	VMOVD X0, ret+48(FP)
	VZEROALL
	RET

TEXT ·cosSimFMA(SB), NOSPLIT, $0
	// Load slice headers;
	// SI: address of a's data
	// DI: address of b's data
	// CX: length of a (total number of elements)
	MOVQ a+0(FP), SI
	MOVQ b+24(FP), DI
	MOVQ a+8(FP), CX
	
	// Y0 accumulates a*b, Y6 accumulates a*a, Y7 accumulates b*b
	VXORPS Y0, Y0, Y0
	VXORPS Y6, Y6, Y6
	VXORPS Y7, Y7, Y7

	// for r8 := 0; r8 < len(a); r8 += 8
	MOVQ $0, R8

Loop:
	CMPQ R8, CX
	JAE Done

	// r8 counts the floats, so to load bytes we scale r8*4
	VMOVUPS 0(SI)(R8*4), Y1
	VMOVUPS 0(DI)(R8*4), Y2

	VFMADD231PS Y1, Y2, Y0
	VFMADD231PS Y1, Y1, Y6
	VFMADD231PS Y2, Y2, Y7

	ADDQ $8, R8
	JMP Loop

Done:
	// Horizontally reduce each of a*b, a*a, a*b separately into its register
	VPERM2F128 $0b00000001, Y0, Y0, Y1
	VADDPS Y1, Y0, Y0
	VHADDPS Y0, Y0, Y0
	VHADDPS Y0, Y0, Y0

	VPERM2F128 $0b00000001, Y6, Y6, Y8
	VADDPS Y8, Y6, Y6
	VHADDPS Y6, Y6, Y6
	VHADDPS Y6, Y6, Y6

	VPERM2F128 $0b00000001, Y7, Y7, Y8
	VADDPS Y8, Y7, Y7
	VHADDPS Y7, Y7, Y7
	VHADDPS Y7, Y7, Y7

	// a*b / sqrt(a*a * b*b)
	MULSS X7, X6
	RSQRTSS X6, X6
	MULSS X6, X0

	VMOVD X0, ret+48(FP)
	VZEROALL
	RET

// dotProductAVX2 calculates the dot product of two []float32
// For simplicity, this function assumes that len(a)==len(b) and len(a)%8==0
TEXT ·dotProductAVX2(SB), NOSPLIT, $0
	// Load slice headers;
	// SI: address of a's data
	// DI: address of b's data
	// CX: length of a (total number of elements)
	MOVQ a+0(FP), SI
	MOVQ b+24(FP), DI
	MOVQ a+8(FP), CX
	
	// Y0 will hold the result; init to 0
	VXORPS Y0, Y0, Y0

	// for r8 := 0; r8 < len(a); r8 += 8
	MOVQ $0, R8

Loop:
	CMPQ R8, CX
	JAE Done

	// r8 counts the floats, so to load bytes we scale r8*4
	VMOVUPS 0(SI)(R8*4), Y1
	VMOVUPS 0(DI)(R8*4), Y2
	VMULPS Y1, Y2, Y3
	VADDPS Y3, Y0, Y0

	ADDQ $8, R8
	JMP Loop

Done:
	// Horizonal addition of a full 8-float YMM register
	VPERM2F128 $0b00000001, Y0, Y0, Y1
	VADDPS Y1, Y0, Y0
	VHADDPS Y0, Y0, Y0
	VHADDPS Y0, Y0, Y0
	VMOVD X0, ret+48(FP)
	VZEROALL
	RET

// dotProductScalar is a scalar (single float at a time) implementation of
// dot product, for testing & comparison.
TEXT ·dotProductScalar(SB), NOSPLIT, $0
	// Load slice headers;
	// SI: address of a's data
	// DI: address of b's data
	// CX: length of a (total number of elements)
	MOVQ a+0(FP), SI
	MOVQ b+24(FP), DI
	MOVQ a+8(FP), CX

	// Accumulator
	XORPS X0, X0

	// for r8 := 0; r8 < len(a); r8++
	MOVQ $0, R8

Loop:
	CMPQ R8, CX
	JAE Done

	MOVSS (SI)(R8*4), X1
	MOVSS (DI)(R8*4), X2
	MULSS X1, X2
	ADDSS X2, X0

	INCQ R8
	JMP Loop

Done:
	MOVSS X0, ret+48(FP)
	RET

