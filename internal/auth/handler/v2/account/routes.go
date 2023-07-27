package accountV2

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes register account version 2 routes with /api prefix
func RegisterRoutes(r fiber.Router) {
	r.Post("/register", Register)
	r.Post("/debug/register", RegisterDebug)
	r.Post("/debug/register/_batch", RegisterDebugInBatch)
	r.Post("/login", Login)
	r.Post("/logout", Logout)
	r.Put("/register", ResetPassword)
	r.Get("/verify/email", VerifyWithEmail)
	r.Get("/verify/email/:email", VerifyWithEmailOld)
	r.Get("/verify/apikey", VerifyWithApikey)
	r.Delete("/users/me", DeleteUserByMe)
}
