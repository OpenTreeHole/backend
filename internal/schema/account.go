package schema

import (
	"github.com/opentreehole/backend/internal/pkg/types"
)

type EmailModel struct {
	// email in email blacklist
	// TODO: add email blacklist
	Email string `json:"email" query:"email" validate:"email"`
}

type LoginRequest struct {
	EmailModel
	Password string `json:"password" minLength:"8" maxLength:"32" validate:"min=8,max=32"`
}

// TokenResponse for Login / Register / ResetPassword
type TokenResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
	Message string `json:"message"`
}

type RegisterRequest struct {
	LoginRequest
	Verification types.VerificationCode `json:"verification" swaggertype:"string"`
}

type ResetPasswordRequest = RegisterRequest
