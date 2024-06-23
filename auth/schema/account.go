package schema

type EmailModel struct {
	// email in email blacklist
	Email string `json:"email" query:"email" validate:"isValidEmail"`
}

type LoginRequest struct {
	EmailModel
	Password string `json:"password" minLength:"8" validate:"required,min=8"`
}
