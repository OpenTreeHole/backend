package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateIdentifier() (string, error) {
	now := time.Now().UnixMicro()
	nowStr := fmt.Sprintf("%x", now)[:8] // first 8 characters of the timestamp
	// Generate a random 6 character string, 14 characters in total (different from the original image proxy (13 characters) )
	randomBytes := make([]byte, 3)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	randomSuffix := hex.EncodeToString(randomBytes)

	// Combine the timestamp and random suffix
	return fmt.Sprintf("%s%s", nowStr, randomSuffix), nil
}

func IsAllowedExtension(ext string) bool {
	allowedExtensions := []string{"jpg", "jpeg", "png", "gif", "webp", "bmp"}
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}
	return false
}
