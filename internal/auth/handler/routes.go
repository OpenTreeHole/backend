package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	accountV2 "github.com/opentreehole/backend/internal/auth/handler/v2/account"
	accountV3 "github.com/opentreehole/backend/internal/auth/handler/v3/account"
)

func RegisterRoutes(app *fiber.App) {
	// docs
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", swagger.HandlerDefault)

	// register v2 routes
	groupV2 := app.Group("/api")

	accountV2.RegisterRoutes(groupV2)

	// register v3 routes
	groupV3 := app.Group("/api/v3")

	accountV3.RegisterRoutes(groupV3)
}
