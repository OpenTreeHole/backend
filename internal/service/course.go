package service

import (
	"github.com/opentreehole/backend/internal/repository"
)

type CourseService interface {
	Service
}

type courseService struct {
	Service
	repository repository.CourseRepository
}

func NewCourseService(service Service, repository repository.CourseRepository) CourseService {
	return &courseService{Service: service, repository: repository}
}
