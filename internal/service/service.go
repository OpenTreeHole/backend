package service

import (
	"context"

	"github.com/opentreehole/backend/internal/config"
	"github.com/opentreehole/backend/pkg/log"
)

type Service interface {
	GetLogger(ctx context.Context) *log.Logger
	GetConfig(ctx context.Context) *config.AllConfig
}

type service struct {
	logger *log.Logger
	conf   *config.AtomicAllConfig
}

func NewService(logger *log.Logger, conf *config.AtomicAllConfig) Service {
	return &service{logger: logger, conf: conf}
}

func (s *service) GetLogger(_ context.Context) *log.Logger {
	return s.logger
}

func (s *service) GetConfig(_ context.Context) *config.AllConfig {
	return s.conf.Load()
}
