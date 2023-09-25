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
	docsHandler := handler.NewDocsHandler()
	pointer := config.NewConfig()
	logger, cleanup := log.NewLogger(pointer)
	validate := handler.NewValidater()
	handlerHandler := handler.NewHandler(logger, validate)
	serviceService := service.NewService(logger, pointer)
	db := repository.NewDB(pointer, logger)
	cacheCache := cache.NewCache(pointer, logger)
	repositoryRepository := repository.NewRepository(db, cacheCache, logger, pointer)
	accountRepository := repository.NewAccountRepository(repositoryRepository)
	accountService := service.NewAccountService(serviceService, accountRepository)
	accountHandler := handler.NewAccountHandler(handlerHandler, accountService)
	divisionHandler := handler.NewDivisionHandler(handlerHandler)
	courseGroupRepository := repository.NewCourseGroupRepository(repositoryRepository)
	reviewRepository := repository.NewReviewRepository(repositoryRepository)
	courseGroupService := service.NewCourseGroupService(serviceService, courseGroupRepository, reviewRepository)
	courseGroupHandler := handler.NewCourseGroupHandler(handlerHandler, courseGroupService, accountRepository)
	courseRepository := repository.NewCourseRepository(repositoryRepository)
	courseService := service.NewCourseService(serviceService, courseRepository, courseGroupRepository, reviewRepository)
	courseHandler := handler.NewCourseHandler(handlerHandler, courseService, courseGroupService, accountRepository)
	reviewService := service.NewReviewService(serviceService, reviewRepository, courseRepository)
	reviewHandler := handler.NewReviewHandler(handlerHandler, reviewService, accountRepository)
	serverServer := server.NewServer(docsHandler, accountHandler, divisionHandler, courseGroupHandler, courseHandler, reviewHandler, logger, pointer)
	return serverServer, func() {
		cleanup()
	}, nil
}

// wire.go:

var HandlerSet = wire.NewSet(handler.NewHandler, handler.NewAccountHandler, handler.NewDocsHandler, handler.NewDivisionHandler, handler.NewCourseGroupHandler, handler.NewCourseHandler, handler.NewReviewHandler)

var ServiceSet = wire.NewSet(service.NewService, service.NewAccountService, service.NewDivisionService, service.NewCourseGroupService, service.NewCourseService, service.NewReviewService)

var RepositorySet = wire.NewSet(repository.NewDB, repository.NewRepository, repository.NewAccountRepository, repository.NewDivisionRepository, repository.NewCourseGroupRepository, repository.NewCourseRepository, repository.NewReviewRepository, repository.NewAchievementRepository)
