package common

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"time"
)

func MiddlewareCustomLogger(c *fiber.Ctx) error {
	startTime := time.Now()
	chainErr := c.Next()

	if chainErr != nil {
		if err := c.App().ErrorHandler(c, chainErr); err != nil {
			_ = c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	latency := time.Since(startTime).Milliseconds()
	userID, ok := c.Locals("user_id").(int)

	attrs := []slog.Attr{
		slog.Int("status_code", c.Response().StatusCode()),
		slog.String("method", c.Method()),
		slog.String("origin_url", c.OriginalURL()),
		slog.String("remote_ip", c.Get("X-Real-IP")),
		slog.Int64("latency", latency),
	}
	if ok {
		attrs = append(attrs, slog.Int("user_id", userID))
	}
	var logLevel = slog.LevelInfo
	if chainErr != nil {
		attrs = append(attrs, slog.String("err", chainErr.Error()))
		logLevel = slog.LevelError
	}
	Logger.LogAttrs(context.Background(), logLevel, "http log", attrs...)
	return nil
}
