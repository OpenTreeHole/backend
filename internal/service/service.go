package service

import (
	"github.com/opentreehole/backend/internal/config"
	"github.com/opentreehole/backend/pkg/log"
)

type Service struct {
	logger *log.Logger
	conf   *config.AtomicAllConfig
}

func NewService(logger *log.Logger, conf *config.AtomicAllConfig) *Service {
	return &Service{logger: logger, conf: conf}
}
