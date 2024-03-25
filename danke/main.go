package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/opentreehole/backend/common"
	"github.com/opentreehole/backend/danke/api"
	"github.com/opentreehole/backend/danke/model"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	common.InitCache()
	model.Init()

	var disableStartupMessage = false
	if viper.GetString(common.EnvMode) == "prod" {
		disableStartupMessage = true
	}
	app := fiber.New(fiber.Config{
		ErrorHandler:          common.ErrorHandler,
		DisableStartupMessage: disableStartupMessage,
	})
	registerMiddlewares(app)
	api.RegisterRoutes(app.Group("/api"))
	go func() {
		err := app.Listen("0.0.0.0:8000")
		if err != nil {
			slog.LogAttrs(context.Background(), slog.LevelError, "app listen failed", slog.String("err", err.Error()))
		}
	}()

	interrupt := make(chan os.Signal, 1)

	// wait for CTRL-C interrupt
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	// close app
	err := app.Shutdown()
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "shutdown failed", slog.String("err", err.Error()))
	}
}

func registerMiddlewares(app *fiber.App) {
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	// app.Use(common.MiddlewareGetUserID)
	if viper.GetString(common.EnvMode) != "bench" {
		app.Use(common.MiddlewareCustomLogger)
	}
	app.Use(pprof.New())
}
