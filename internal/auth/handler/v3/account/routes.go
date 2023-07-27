package accountV3

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes register account version 3 routes with /api/v3 prefix
func RegisterRoutes(r fiber.Router) {
	r.Post("/register", Register)
	r.Post("/debug/register", RegisterDebug)
	r.Post("/login", Login)
	r.Post("/logout", Logout)
	r.Post("/refresh", Refresh)
	r.Put("/register", ResetPassword)
	r.Get("/verify/email", VerifyWithEmail)
	r.Delete("/users/me", DeleteUserByMe)
	r.Delete("/users/:id", DeleteUserByID)
}
