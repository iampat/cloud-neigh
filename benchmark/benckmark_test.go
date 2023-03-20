package benchmark_test

import "testing"

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

var result int

func benchmarkFib(i int, b *testing.B) {
	var r int
	for n := 0; n < b.N; n++ {
		r = Fib(i)
	}
	result = r
}

func BenchmarkFib10(b *testing.B) { benchmarkFib(10, b) }

// func BenchmarkFib2(b *testing.B)  { benchmarkFib(2, b) }
// func BenchmarkFib3(b *testing.B)  { benchmarkFib(3, b) }
// func BenchmarkFib10(b *testing.B) { benchmarkFib(10, b) }
// func BenchmarkFib20(b *testing.B) { benchmarkFib(20, b) }
// func BenchmarkFib40(b *testing.B) { benchmarkFib(40, b) }

// var result int

// func BenchmarkFibComplete(b *testing.B) {
// 	var r int
// 	for n := 0; n < b.N; n++ {
// 		// always record the result of Fib to prevent
// 		// the compiler eliminating the function call.
// 		r = Fib(10)
// 	}
// 	// always store the result to a package level variable
// 	// so the compiler cannot eliminate the Benchmark itself.
// 	result = r
// }
