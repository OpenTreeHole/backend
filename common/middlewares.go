package common

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"strconv"
	"time"
)

func GetUserID(c *fiber.Ctx) (int, error) {
	// get user id from header: X-Consumer-Username if through Kong
	username := c.Get("X-Consumer-Username")
	if username != "" {
		id, err := strconv.Atoi(username)
		if err == nil {
			return id, nil
		}
	}

	// get user id from jwt
	// ID and UserID are both valid
	var user struct {
		ID     int `json:"id"`
		UserID int `json:"user_id"`
	}

	token := GetJWTToken(c)
	if token == "" {
		return 0, Unauthorized("Unauthorized")
	}

	err := ParseJWTToken(token, &user)
	if err != nil {
		return 0, Unauthorized("Unauthorized")
	}

	if user.ID != 0 {
		return user.ID, nil
	} else if user.UserID != 0 {
		return user.UserID, nil
	}

	return 0, Unauthorized("Unauthorized")
}

func MiddlewareGetUserID(c *fiber.Ctx) error {
	userID, err := GetUserID(c)
	if err == nil {
		c.Locals("user_id", userID)
	}

	return c.Next()
}

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
