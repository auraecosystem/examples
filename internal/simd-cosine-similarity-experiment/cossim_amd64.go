package main

// These functions calculate the cosine similarity between a and b
func cosSimScalar(a, b []float32) float32
func cosSimAVX2(a, b []float32) float32
func cosSimFMA(a, b []float32) float32

func dotProductAVX2(a, b []float32) float32
func dotProductScalar(a, b []float32) float32
