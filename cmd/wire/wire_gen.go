// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/opentreehole/backend/internal/config"
	"github.com/opentreehole/backend/internal/handler"
	"github.com/opentreehole/backend/internal/pkg/cache"
	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/server"
	"github.com/opentreehole/backend/internal/service"
	"github.com/opentreehole/backend/pkg/log"
)

// Injectors from wire.go:

func NewApp() (*server.Server, func(), error) {
	pointer := config.NewConfig()
	logger, cleanup := log.NewLogger(pointer)
	validate := handler.NewValidater()
	handlerHandler := handler.NewHandler(logger, validate)
	serviceService := service.NewService(logger)
	db := repository.NewDB(pointer, logger)
	cacheCache := cache.NewCache(pointer, logger)
	repositoryRepository := repository.NewRepository(db, cacheCache, logger, pointer)
	accountRepository := repository.NewAccountRepository(repositoryRepository)
	accountService := service.NewAccountService(serviceService, accountRepository)
	accountHandler := handler.NewAccountHandler(handlerHandler, accountService)
	serverServer := server.NewServer(accountHandler, logger, pointer)
	return serverServer, func() {
		cleanup()
	}, nil
}

// wire.go:

var HandlerSet = wire.NewSet(handler.NewHandler, handler.NewAccountHandler)

var ServiceSet = wire.NewSet(service.NewService, service.NewAccountService)

var RepositorySet = wire.NewSet(repository.NewDB, repository.NewRepository, repository.NewAccountRepository)
