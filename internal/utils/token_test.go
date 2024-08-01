package utils

import "testing"

// BenchmarkGetToken function measures the performance of generating short tokens using the Get method of a ShortToken object.
// The NewShortToken function is used to create a ShortToken object with a specified length.
// The ShortToken interface provides two methods:
// - Get: returns a new random short token
// - Check: checks the length and alphabet of a given token
// The shortToken struct represents a ShortToken implementation. It has two fields:
// - length: the desired length of the short token
// - bufSize: the size of the byte buffer used to generate the token
// The Get method of the shortToken struct generates a random byte buffer using the crypto/rand package. It then returns a shortened BASE64 representation of the buffer as a string.
// The Check method of the shortToken struct validates the length and alphabet of a given token. It checks if the length of the token matches the desired length and if the token conforms
func BenchmarkGetToken(b *testing.B) {
	token := NewShortToken(20)
	b.Run("create 20", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			token.Get()
		}
	})
}
