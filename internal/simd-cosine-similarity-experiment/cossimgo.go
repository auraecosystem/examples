package main

import (
	"math"
	"unsafe"
)

// cosSimGo is a pure Go implementation of cosine similarity.
func cosSimGo(a, b []float32) float32 {
	var dot, na, nb float32
	for i, x := range a {
		y := b[i]
		dot += x * y
		na += x * x
		nb += y * y
	}
	return dot / float32(math.Sqrt(float64(na*nb)))
}

// dotProductGo is a pure Go implementation of dot product.
func dotProductGo(a, b []float32) float32 {
	var dot float32
	for i, x := range a {
		dot += x * b[i]
	}
	return dot
}

func dotProductGoUnrolled(a, b []float32) float32 {
	var dot0, dot1, dot2, dot3 float32
	for i := 0; i < len(a); i += 8 {
		p := (*[8]float32)(unsafe.Pointer(&a[i]))
		q := (*[8]float32)(unsafe.Pointer(&b[i]))
		dot0 += p[0]*q[0] + p[1]*q[1]
		dot1 += p[2]*q[2] + p[3]*q[3]
		dot2 += p[4]*q[4] + p[5]*q[5]
		dot3 += p[6]*q[6] + p[7]*q[7]
	}
	return dot0 + dot1 + dot2 + dot3
}

// sumReduceSliceGo sum-reduces two slices into a scalar.
func sumReduceSliceGo(a, b []float32) float32 {
	var sum float32
	for i, x := range a {
		sum += x + b[i]
	}
	return sum
}
