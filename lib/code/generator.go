package code

import (
	"golang.org/x/exp/rand"
	"time"
)

const (
	SixElementsCodeLength = 6
)

// GenerateCodeFromSalt generates 6 symbols length code
func GenerateCodeFromSalt(salt string) string {
	buffer := make([]byte, SixElementsCodeLength)
	rand.Seed(uint64(time.Now().UnixNano()))

	for i := range buffer {
		buffer[i] = salt[rand.Intn(len(salt))]
	}

	return string(buffer)
}
