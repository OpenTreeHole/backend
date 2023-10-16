package server

import (
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"

	"github.com/opentreehole/backend/internal/config"
	"github.com/opentreehole/backend/internal/schema"
	"github.com/opentreehole/backend/pkg/utils"
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

	token := utils.GetJWTToken(c)
	if token == "" {
		return 0, schema.Unauthorized("Unauthorized")
	}

	err := utils.ParseJWTToken(token, &user)
	if err != nil {
		return 0, schema.Unauthorized("Unauthorized")
	}

	if user.ID != 0 {
		return user.ID, nil
	} else if user.UserID != 0 {
		return user.UserID, nil
	}

	return 0, schema.Unauthorized("Unauthorized")
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

	output := log.Info().
		Int("status_code", c.Response().StatusCode()).
		Str("method", c.Method()).
		Str("origin_url", c.OriginalURL()).
		Str("remote_ip", c.Get("X-Real-IP")).
		Int64("latency", latency).
		Str("user_agent", c.Get("User-Agent"))
	if ok {
		output = output.Int("user_id", userID)
	}
	if chainErr != nil {
		output = output.Err(chainErr)
	}
	output.Msg("http log")
	return nil
}

func StackTraceHandler(_ *fiber.Ctx, e any) {
	log.Error().Any("panic", e).Bytes("stack", debug.Stack()).Msg("stacktrace")
}

func RegisterMiddlewares(conf *config.AtomicAllConfig) func(app *fiber.App) {
	return func(app *fiber.App) {
		app.Use(fiberrecover.New(fiberrecover.Config{EnableStackTrace: true, StackTraceHandler: StackTraceHandler}))
		app.Use(MiddlewareGetUserID)
		if conf.Load().Mode != "bench" {
			app.Use(MiddlewareCustomLogger)
		}
		if conf.Load().Mode == "dev" {
			app.Use(cors.New(cors.Config{AllowOrigins: "*"})) // for swag docs
		}
	}
}
