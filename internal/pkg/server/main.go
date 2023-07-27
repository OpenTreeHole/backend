package server

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/opentreehole/backend/internal/pkg/config"
	"github.com/opentreehole/backend/internal/pkg/database"
	"github.com/opentreehole/backend/pkg/utils"
)

type Config struct {
	database.Config
	AppName             string
	RegisterMiddlewares func(app *fiber.App)
	RegisterRoutes      func(app *fiber.App)
}

func Init(serverConfig *Config) *fiber.App {
	// init config
	config.Init()

	// init database, including gorm, cache and search engine
	database.Init(serverConfig.Config)

	// init fiber
	app := fiber.New(fiber.Config{
		AppName:               serverConfig.AppName,
		DisableStartupMessage: true,
		JSONDecoder:           json.Unmarshal,
		JSONEncoder:           json.Marshal,
		ErrorHandler:          utils.ErrorHandler,
	})

	// register middlewares
	if serverConfig.RegisterMiddlewares != nil {
		serverConfig.RegisterMiddlewares(app)
	} else {
		// default middlewares
		RegisterMiddlewares(app)
	}

	// register routes
	if serverConfig.RegisterRoutes != nil {
		serverConfig.RegisterRoutes(app)
	}

	return app
}

func Run(serverConfig *Config) {
	app := Init(serverConfig)

	// start server
	go func() {
		err := app.Listen("0.0.0.0:" + strconv.Itoa(config.Config.Port))
		if err != nil {
			log.Fatal().Err(err).Msg("app listen failed")
		}
	}()

	interrupt := make(chan os.Signal, 1)

	// wait for CTRL-C interrupt
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	// close app
	err := app.Shutdown()
	if err != nil {
		log.Err(err).Msg("error shutdown app")
	}
}
