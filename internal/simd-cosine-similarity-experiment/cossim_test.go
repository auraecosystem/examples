package main

import (
	"math/rand/v2"
	"testing"
)

// almostEqualF32 checks if a and b are "almost equal", within a reasonably
// permissive epsilon for comparison, because the vector versions add & multiply
// floats in a different order, which affects precision.
func almostEqualF32(a, b float32) bool {
	const float32Eps = 0.001
	abs := a - b
	if a-b < 0 {
		abs = -abs
	}
	return abs <= float32Eps
}

func makeRandSlice(sz int) []float32 {
	s := make([]float32, sz)
	for i := range sz {
		s[i] = rand.Float32()
	}
	return s
}

func TestDotProduct(t *testing.T) {
	sz := 8 * 128
	va := makeRandSlice(sz)
	vb := makeRandSlice(sz)

	dotWant := dotProductGo(va, vb)
	dotGotScalar := dotProductScalar(va, vb)
	dotGotAVX2 := dotProductAVX2(va, vb)

	if !almostEqualF32(dotWant, dotGotScalar) {
		t.Errorf("scalar got %v, want %v", dotGotScalar, dotWant)
	}
	if !almostEqualF32(dotWant, dotGotAVX2) {
		t.Errorf("AVX2 got %v, want %v", dotGotAVX2, dotWant)
	}
}

func TestDotProductSmall(t *testing.T) {
	aa := []float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0}
	aaa := append(aa, aa...)
	bb := []float32{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0}
	bbb := append(bb, bb...)

	dotWant := dotProductGo(aaa, bbb)
	dotGotScalar := dotProductScalar(aaa, bbb)
	dotGotAVX2 := dotProductAVX2(aaa, bbb)

	if !almostEqualF32(dotWant, dotGotScalar) {
		t.Errorf("scalar got %v, want %v", dotGotScalar, dotWant)
	}
	if !almostEqualF32(dotWant, dotGotAVX2) {
		t.Errorf("AVX2 got %v, want %v", dotGotAVX2, dotWant)
	}
}

func TestCosSimSmall(t *testing.T) {
	a := []float32{1.1, 2.2, 3.1, 1.5, 2.2, 3.3, 0.5, 3.5}
	b := []float32{0.7, -0.2, 0.9, -1.1, 2.6, -0.5, 0.8, 0.2}

	want := cosSimGo(a, b)
	gotScalar := cosSimScalar(a, b)
	gotAVX2 := cosSimAVX2(a, b)

	if !almostEqualF32(gotScalar, want) {
		t.Errorf("scalar got %v, want %v", gotScalar, want)
	}
	if !almostEqualF32(gotAVX2, want) {
		t.Errorf("AVX2 got %v, want %v", gotAVX2, want)
	}
}

func TestCosSim(t *testing.T) {
	sz := 8 * 128
	va := makeRandSlice(sz)
	vb := makeRandSlice(sz)

	cosWant := cosSimGo(va, vb)
	cosScalar := cosSimScalar(va, vb)
	cosAVX2 := cosSimAVX2(va, vb)
	cosFMA := cosSimFMA(va, vb)

	if !almostEqualF32(cosWant, cosScalar) {
		t.Errorf("scalar got %v, want %v", cosScalar, cosWant)
	}
	if !almostEqualF32(cosWant, cosAVX2) {
		t.Errorf("AVX2 got %v, want %v", cosAVX2, cosWant)
	}
	if !almostEqualF32(cosWant, cosFMA) {
		t.Errorf("FMA got %v, want %v", cosFMA, cosWant)
	}
}

const benchArrSize = 1 * 1024 * 1024

func BenchmarkCosSimGo(b *testing.B) {
	aa := makeRandSlice(benchArrSize)
	bb := makeRandSlice(benchArrSize)

	// Each cosine similarity calculation reads two slices of size benchArrSize
	// from memory. Each slice holds floats, so its total size in memory is
	// benchArrSize * 4 bytes.
	b.SetBytes(benchArrSize * 4 * 2)
	for b.Loop() {
		cosSimGo(aa, bb)
	}
}

func BenchmarkCosSimScalar(b *testing.B) {
	aa := makeRandSlice(benchArrSize)
	bb := makeRandSlice(benchArrSize)

	b.SetBytes(benchArrSize * 4 * 2)
	for b.Loop() {
		cosSimScalar(aa, bb)
	}
}

func BenchmarkCosSimAVX2(b *testing.B) {
	aa := makeRandSlice(benchArrSize)
	bb := makeRandSlice(benchArrSize)

	b.SetBytes(benchArrSize * 4 * 2)
	for b.Loop() {
		cosSimAVX2(aa, bb)
	}
}

func BenchmarkCosSimFMA(b *testing.B) {
	aa := makeRandSlice(benchArrSize)
	bb := makeRandSlice(benchArrSize)

	b.SetBytes(benchArrSize * 4 * 2)
	for b.Loop() {
		cosSimFMA(aa, bb)
	}
}

func BenchmarkDotProductGo(b *testing.B) {
	aa := makeRandSlice(benchArrSize)
	bb := makeRandSlice(benchArrSize)

	b.SetBytes(benchArrSize * 4 * 2)
	for b.Loop() {
		dotProductGo(aa, bb)
	}
}

func BenchmarkSumReduceSliceGo(b *testing.B) {
	aa := makeRandSlice(benchArrSize)
	bb := makeRandSlice(benchArrSize)

	b.SetBytes(benchArrSize * 4 * 2)
	for b.Loop() {
		sumReduceSliceGo(aa, bb)
	}
}
