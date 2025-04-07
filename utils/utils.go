// utils.go
// Utility functions used across the application.

package utils

import (
	"crypto/rand"
	"fmt"
	"io"
)

// GenerateRandomFileName returns a random 8-character filename with a .txt extension.
// Used to name output files uniquely and safely.
func GenerateRandomFileName() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		panic("Failed to generate random filename: " + err.Error())
	}
	for i, b := range bytes {
		bytes[i] = letters[int(b)%len(letters)]
	}
	return fmt.Sprintf("%s.txt", string(bytes))
}
