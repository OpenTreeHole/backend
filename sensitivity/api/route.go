package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func RegisterRouts(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/api")
	})
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", swagger.HandlerDefault)

	api := app.Group("/api")
	registerRoutes(api)
}

func registerRoutes(r fiber.Router) {
	r.Get("", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to sensitivity checker API"})
	})

	r.Post("/sensitive/text", CheckSensitiveText)

}
