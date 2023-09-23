package server

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/opentreehole/backend/internal/config"
	"github.com/opentreehole/backend/internal/handler"
	"github.com/opentreehole/backend/internal/schema"
	"github.com/opentreehole/backend/pkg/log"
)

type Server struct {
	config       *config.AtomicAllConfig
	logger       *log.Logger
	rootRegister []handler.RouteRegister
	handlers     []handler.RouteRegister
}

func NewServer(
	accountHandler handler.AccountHandler,
	divisionHandler handler.DivisionHandler,
	docsHandler handler.DocsHandler,
	logger *log.Logger,
	config *config.AtomicAllConfig,
) *Server {
	var handlers []handler.RouteRegister

	if config.Load().Modules.Auth {
		handlers = append(handlers, accountHandler)
	}

	if config.Load().Modules.Notification {
		// TODO
	}

	if config.Load().Modules.Treehole {
		handlers = append(handlers, divisionHandler)
	}

	if config.Load().Modules.CurriculumBoard {
		// TODO
	}
	return &Server{
		logger: logger,
		config: config,
		rootRegister: []handler.RouteRegister{
			docsHandler,
		},
		handlers: handlers,
	}
}

func (s *Server) Run() {
	var app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          schema.ErrorHandler,
	})

	RegisterMiddlewares(s.config)(app)
	for _, h := range s.rootRegister {
		h.RegisterRoute(app)
	}
	for _, h := range s.handlers {
		h.RegisterRoute(app.Group("/api"))
	}

	// start server
	go func() {
		err := app.Listen("0.0.0.0:" + strconv.Itoa(s.config.Load().Port))
		if err != nil {
			s.logger.Fatal("error start server", zap.Error(err))
		}
	}()

	interrupt := make(chan os.Signal, 1)

	// wait for CTRL-C interrupt
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	// close app
	err := app.Shutdown()
	if err != nil {
		s.logger.Fatal("error shutdown server", zap.Error(err))
	}
}
