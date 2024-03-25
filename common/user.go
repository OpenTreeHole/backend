package common

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"strings"
	"time"
)

type UserClaims struct {
	ID        int        `json:"id,omitempty"`
	UserID    int        `json:"user_id,omitempty"`
	UID       int        `json:"uid,omitempty"`
	IsAdmin   bool       `json:"is_admin"`
	ExpiresAt *time.Time `json:"exp,omitempty"`
}

type User struct {
	ID      int  `json:"id"`
	IsAdmin bool `json:"is_admin"`
}

var (
	ErrJWTTokenRequired = Unauthorized("jwt token required")
	ErrInvalidJWTToken  = Unauthorized("invalid jwt token")
)

// GetJWTToken extracts token from header or cookie
// return empty string if not found
func GetJWTToken(c *fiber.Ctx) string {
	tokenString := c.Get("Authorization") // token in header
	if tokenString == "" {
		tokenString = c.Cookies("access") // token in cookie
	}
	return tokenString
}

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
		slog.Error("jwt parse error", "err", err, "payload_string", payloadString)
		return ErrInvalidJWTToken
	}

	err = json.Unmarshal(payloadBytes, user)
	if err != nil {
		slog.Error("jwt parse error", "err", err, "payload_string", payloadString)
		return ErrInvalidJWTToken
	}

	return nil
}

func GetCurrentUser(c *fiber.Ctx) (user *User, err error) {
	token := GetJWTToken(c)
	if token == "" {
		return nil, Unauthorized("Unauthorized")
	}

	var userClaims UserClaims
	err = ParseJWTToken(token, &userClaims)
	if err != nil {
		return nil, err
	}

	user = &User{}
	if userClaims.ID == 0 && userClaims.UserID == 0 && userClaims.UID == 0 {
		return nil, Unauthorized("Unauthorized")
	} else {
		if userClaims.ID != 0 {
			user.ID = userClaims.ID
		} else if userClaims.UserID != 0 {
			user.ID = userClaims.UserID
		} else {
			user.ID = userClaims.UID
		}
	}

	if userClaims.ExpiresAt != nil && userClaims.ExpiresAt.Before(time.Now()) {
		return nil, Unauthorized("token expired")
	}

	return
}
