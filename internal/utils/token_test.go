package utils

import "testing"

func BenchmarkFibo(b *testing.B) {
	token := NewShortToken(20)
	b.Run("create 20", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			token.Get()
		}
	})
}
