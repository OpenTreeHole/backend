package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type DocsHandler interface {
	RouteRegister
}

type docsHandler struct{}

func NewDocsHandler() DocsHandler {
	return &docsHandler{}
}

func (h *docsHandler) RegisterRoute(router fiber.Router) {
	router.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	router.Get("/docs/*", swagger.HandlerDefault)
}
