package main

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	. "github.com/opentreehole/backend/common"
	. "github.com/opentreehole/backend/image_hosting/api"
	. "github.com/opentreehole/backend/image_hosting/model"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler:          ErrorHandler,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
		BodyLimit:             128 * 1024 * 1024,
	})
	// will catch every panic
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// app.Use(common.MiddlewareGetUserID)

	// will catch every HTTP request
	// app.Use(MiddlewareCustomLogger)
	Init()
	router := app.Group("/api")
	router.Post("/uploadImage", UploadImage)
	router.Get("/i/:year/:month/:day/:identifier", GetImage) // get images based on the identifier(excluding the extension)

	go func() {
		slog.LogAttrs(context.Background(), slog.LevelInfo, "Service started.")
		err := app.Listen("0.0.0.0:8000")
		if err != nil {
			slog.LogAttrs(context.Background(), slog.LevelError, "Wrong hostname", slog.String("err", err.Error()))
		}

	}()

	interrupt := make(chan os.Signal, 1)

	// listen for interrupt signal (Ctrl+C)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt
	slog.LogAttrs(context.Background(), slog.LevelInfo, "Shutting down the server.")
	err := app.Shutdown()
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "shutdown failed", slog.String("err", err.Error()))
	}
}
