package utils

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var (
	ErrJWTTokenRequired = errors.New("jwt token required")
	ErrInvalidJWTToken  = errors.New("invalid jwt token")
)

// ParseJWTToken extracts and parse token, whatever start with "Bearer " or not
func ParseJWTToken(token string, user any) error {
	// remove "Bearer " prefix if exists
	if strings.HasPrefix(token, "Bearer ") {
		token = token[7:]
	}
	token = strings.TrimSpace(token)
	payloads := strings.SplitN(token, ".", 3) // extract "Bearer "
	if len(payloads) < 3 {
		return ErrJWTTokenRequired
	}

	payloadString := payloads[1]

	// jwt encoding ignores padding, so RawStdEncoding should be used instead of StdEncoding
	// jwt encoding uses url safe base64 encoding, so RawURLEncoding should be used instead of RawStdEncoding
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadString) // the middle one is payload
	if err != nil {
		log.Err(err).Str("payload_string", payloadString).Msg("jwt parse error")
		return ErrInvalidJWTToken
	}

	err = json.Unmarshal(payloadBytes, user)
	if err != nil {
		log.Err(err).Str("payload_string", payloadString).Msg("jwt parse error")
		return ErrInvalidJWTToken
	}

	return nil
}

// GetJWTToken extracts token from header or cookie
// return empty string if not found
func GetJWTToken(c *fiber.Ctx) string {
	tokenString := c.Get("Authorization") // token in header
	if tokenString == "" {
		tokenString = c.Cookies("access") // token in cookie
	}
	return tokenString
}
