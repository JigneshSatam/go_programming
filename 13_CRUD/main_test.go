package main

import "testing"

func BenchmarkCalculate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// callInit()
		callOld()
	}
}
