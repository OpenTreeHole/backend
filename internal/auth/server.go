package auth

import (
	"github.com/opentreehole/backend/internal/auth/handler"
	"github.com/opentreehole/backend/internal/auth/model"
	"github.com/opentreehole/backend/internal/pkg/database"
	"github.com/opentreehole/backend/internal/pkg/server"
)

var Config = &server.Config{
	AppName:             "auth",
	RegisterMiddlewares: server.RegisterMiddlewares,
	RegisterRoutes:      handler.RegisterRoutes,
	Config: database.Config{
		GormModels:         []any{model.User{}},
		EnableGorm:         true,
		EnableCache:        true,
		EnableSearchEngine: false,
	},
}
