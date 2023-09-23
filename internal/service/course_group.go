package service

import (
	"github.com/opentreehole/backend/internal/repository"
)

type CourseGroupService interface {
	Service
}

type courseGroupService struct {
	Service
	repository repository.CourseGroupRepository
}

func NewCourseGroupService(service Service, repository repository.CourseGroupRepository) CourseGroupService {
	return &courseGroupService{Service: service, repository: repository}
}
