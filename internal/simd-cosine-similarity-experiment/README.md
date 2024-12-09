Implementation of Cosine Similarity in assembly (using Go syntax) using AVX2
instructions (with a variant using FMA).

The benchmarks use `testing.B.Loop`, which is new in Go 1.24; to run the
tests before 1.24 is released, use an updated gotip:

    gotip test -v -bench=.
