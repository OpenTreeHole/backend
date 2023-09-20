//go:build wireinject
// +build wireinject

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

var HandlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewAccountHandler,
	handler.NewDocsHandler,
)

var ServiceSet = wire.NewSet(
	service.NewService,
	service.NewAccountService,
)

var RepositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRepository,
	repository.NewAccountRepository,
)

func NewApp() (*server.Server, func(), error) {
	wire.Build(
		config.NewConfig,
		log.NewLogger,
		cache.NewCache,
		RepositorySet,
		ServiceSet,
		HandlerSet,
		handler.NewValidater,
		server.NewServer,
	)
	return &server.Server{}, nil, nil
}
