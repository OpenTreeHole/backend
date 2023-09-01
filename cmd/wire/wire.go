//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/spf13/viper"

	"github.com/opentreehole/backend/internal/handler"
	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/server"
	"github.com/opentreehole/backend/internal/service"
	"github.com/opentreehole/backend/pkg/log"
)

var HandlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewAccountHandler,
)

var ServiceSet = wire.NewSet(
	service.NewService,
	service.NewAccountService,
)

var RepositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewCacher,
	repository.NewRepository,
	repository.NewAccountRepository,
)

func NewApp(*viper.Viper) (*server.Server, func(), error) {
	wire.Build(
		RepositorySet,
		ServiceSet,
		HandlerSet,
		server.NewServer,
		log.NewLogger,
		handler.NewValidater,
	)
	return &server.Server{}, nil, nil
}
