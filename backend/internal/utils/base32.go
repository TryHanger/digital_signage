package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateShortToken создаёт короткий токен длиной 12 символов (~48 бит)
func GenerateShortToken() string {
	b := make([]byte, 6)
	rand.Read(b)
	return hex.EncodeToString(b) // например "a1b2c3d4e5f6"
}
