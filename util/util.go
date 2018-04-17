package util

import (
	"math/rand"
)

// Urandom produces a random stream of bytes
func Urandom(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(rand.Int31())
	}

	return b
}
