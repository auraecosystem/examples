package main

import "math"

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

// sumReduceSliceGo sum-reduces two slices into a scalar.
func sumReduceSliceGo(a, b []float32) float32 {
	var sum float32
	for i, x := range a {
		sum += x + b[i]
	}
	return sum
}
