package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetHashSHA256(value string, salt string) string {
	encoder := sha256.New()
	encoder.Write([]byte(value + salt))
	return hex.EncodeToString(encoder.Sum(nil))
}
