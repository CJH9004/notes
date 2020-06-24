package main

import "testing"

func BenchmarkArray(b *testing.B) {
	d := [1000]int{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d[i%1000]++
	}
}

func BenchmarkSlice(b *testing.B) {
	d := make([]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d[i%1000]++
	}
}

func BenchmarkSliceWith100(b *testing.B) {
	d := make([]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d[i%100]++
	}
}

func BenchmarkSliceWith10(b *testing.B) {
	d := make([]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d[i%10]++
	}
}

func BenchmarkMap(b *testing.B) {
	d := make(map[int]int, 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d[i%1000]++
	}
}

func BenchmarkDynMap(b *testing.B) {
	d := make(map[int]int)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d[i%1000]++
	}
}

func BenchmarkDynMapWith100(b *testing.B) {
	d := make(map[int]int)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d[i%100]++
	}
}

func BenchmarkDynMapWith10(b *testing.B) {
	d := make(map[int]int)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d[i%10]++
	}
}
