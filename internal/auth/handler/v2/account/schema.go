package accountV2

import (
	"fmt"
	"strconv"
	"strings"
)

type EmailModel struct {
	// email in email blacklist
	Email string `json:"email" query:"email" validate:"isValidEmail"`
}

type LoginRequest struct {
	EmailModel
	Password string `json:"password" minLength:"8"`
}

type TokenResponse struct {
	Access  string `json:"access,omitempty"`
	Refresh string `json:"refresh,omitempty"`
	Message string `json:"message,omitempty"`
}

type RegisterRequest struct {
	LoginRequest
	Verification VerificationType `json:"verification" minLength:"6" maxLength:"6" swaggerType:"string"`
}

type RegisterInBatchRequest struct {
	Data []LoginRequest `json:"data"`
}

type VerificationType string

func (v *VerificationType) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	// Ignore null, like in the main JSON package.
	if s == "null" {
		return nil
	}

	number, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	*v = VerificationType(fmt.Sprintf("%06d", number))
	return nil
}

func (v *VerificationType) UnmarshalText(data []byte) error {
	s := strings.Trim(string(data), `"`)
	// Ignore null, like in the main JSON package.
	if s == "" {
		return nil
	}

	*v = VerificationType(s)
	return nil
}

type EmailVerifyResponse struct {
	Message string `json:"message"`
	Scope   string `json:"scope" enums:"register,reset"`
}

type ApikeyRequest struct {
	EmailModel
	Apikey        string `json:"apikey" query:"apikey"`
	CheckRegister bool   `json:"check_register" query:"check_register" default:"false"` // if true, return whether registered
}

type ApikeyResponse struct {
	EmailVerifyResponse
	Code string `json:"code"`
}
