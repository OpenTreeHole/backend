package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateIdentifier() (string, error) {
	now := time.Now().UnixMicro()

	// Generate a random 6-byte (12-character) string
	randomBytes := make([]byte, 6)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	randomSuffix := hex.EncodeToString(randomBytes)

	// Combine the timestamp and random suffix
	return fmt.Sprintf("%x%s", now, randomSuffix), nil
}
