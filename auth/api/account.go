package api

import "github.com/gofiber/fiber/v2"

func Login(c *fiber.Ctx) (err error) {
	return c.JSON(nil)
}
