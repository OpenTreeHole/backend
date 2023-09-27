package server

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/opentreehole/backend/internal/config"
	"github.com/opentreehole/backend/internal/handler"
	"github.com/opentreehole/backend/internal/schema"
	"github.com/opentreehole/backend/pkg/log"
)

type Server struct {
	config       *config.AtomicAllConfig
	logger       *log.Logger
	db           *gorm.DB
	app          *fiber.App
	rootRegister []handler.RouteRegister
	handlers     []handler.RouteRegister
}

func NewServer(
	// docs
	docsHandler handler.DocsHandler,

	// auth
	accountHandler handler.AccountHandler,

	// treehole
	divisionHandler handler.DivisionHandler,

	// curriculum board
	courseGroupHandler handler.CourseGroupHandler,
	courseHandler handler.CourseHandler,
	reviewHandler handler.ReviewHandler,

	// others
	logger *log.Logger,
	config *config.AtomicAllConfig,
	db *gorm.DB,
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
		handlers = append(handlers,
			courseGroupHandler,
			courseHandler,
			reviewHandler,
		)
	}
	return &Server{
		logger: logger,
		config: config,
		db:     db,
		rootRegister: []handler.RouteRegister{
			docsHandler,
		},
		handlers: handlers,
	}
}

func (s *Server) GetFiberApp() *fiber.App {
	if s.app != nil {
		return s.app
	}

	var disableStartupMessage = true
	if s.config.Load().Mode == "dev" {
		disableStartupMessage = false
	}
	s.app = fiber.New(fiber.Config{
		DisableStartupMessage: disableStartupMessage,
		ErrorHandler:          schema.ErrorHandler,
	})

	RegisterMiddlewares(s.config)(s.app)
	for _, h := range s.rootRegister {
		h.RegisterRoute(s.app)
	}
	for _, h := range s.handlers {
		h.RegisterRoute(s.app.Group("/api"))
	}
	s.app.Get("/api", func(c *fiber.Ctx) error {
		// TODO: add meta info
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})
	return s.app
}

func (s *Server) GetDB() *gorm.DB {
	return s.db
}

func (s *Server) Run() {
	app := s.GetFiberApp()

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
