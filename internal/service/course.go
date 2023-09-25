package service

import (
	"context"

	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/schema"
)

type CourseService interface {
	Service

	ListCoursesV1(ctx context.Context) (response []*schema.CourseGroupV1Response, err error)
}

type courseService struct {
	Service
	courseGroupRepository repository.CourseGroupRepository
	courseRepository      repository.CourseRepository
}

func NewCourseService(
	service Service,
	courseRepository repository.CourseRepository,
	courseGroupRepository repository.CourseGroupRepository,
) CourseService {
	return &courseService{
		Service:               service,
		courseRepository:      courseRepository,
		courseGroupRepository: courseGroupRepository,
	}
}

func (s *courseService) ListCoursesV1(ctx context.Context) (response []*schema.CourseGroupV1Response, err error) {
	groups, err := s.courseGroupRepository.FindAllGroups(ctx, repository.WithGroupCourses())
	if err != nil {
		return nil, err
	}

	response = make([]*schema.CourseGroupV1Response, 0, len(groups))
	for _, group := range groups {
		response = append(response, new(schema.CourseGroupV1Response).FromModel(nil, group, nil))
	}

	return response, nil
}
