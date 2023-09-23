package service

import (
	"github.com/opentreehole/backend/internal/repository"
)

type ReviewService interface {
	Service
}

type reviewService struct {
	Service
	repository repository.ReviewRepository
}

func NewReviewService(service Service, repository repository.ReviewRepository) ReviewService {
	return &reviewService{Service: service, repository: repository}
}
