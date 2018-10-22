package storage

import (
	"testing"
	"time"
)

// BenchmarkSetKey benchmarks setting keys.
func BenchmarkSetKey(b *testing.B) {
	s, _ := New(30 * time.Second)
	for i := 0; i < b.N; i++ {
		s.Set(string(i), i)
	}
}

// BenchmarkGetKey benchmarks getting keys.
func BenchmarkGetKey(b *testing.B) {
	s, _ := New(30 * time.Second)
	for i := 0; i < b.N; i++ {
		s.Get(string(i))
	}
}

// BenchmarkSetGetKey benchmarks both setting and getting keys.
func BenchmarkSetGetKey(b *testing.B) {
	s, _ := New(30 * time.Second)
	for i := 0; i < b.N; i++ {
		s.Set(string(i), i)
		s.Get(string(i))
	}
}
