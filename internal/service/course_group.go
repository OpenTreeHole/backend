package service

import (
	"context"

	"github.com/opentreehole/backend/internal/model"
	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/schema"
)

type CourseGroupService interface {
	Service

	GetGroupByIDV1(ctx context.Context, user *model.User, id int) (response *schema.CourseGroupV1Response, err error)
}

type courseGroupService struct {
	Service
	courseGroupRepository repository.CourseGroupRepository
}

func NewCourseGroupService(
	service Service,
	courseGroupRepository repository.CourseGroupRepository,
) CourseGroupService {
	return &courseGroupService{
		Service:               service,
		courseGroupRepository: courseGroupRepository,
	}
}

func (c *courseGroupService) GetGroupByIDV1(ctx context.Context, user *model.User, id int) (response *schema.CourseGroupV1Response, err error) {
	// 获取课程组，同时加载课程和评教，和评教对应的历史记录，和评教对应的用户成就
	group, err := c.courseGroupRepository.FindGroupByID(ctx, id, repository.WithReviewsAndHistory(), repository.WithReviewsUserAchievements())
	if err != nil {
		return nil, err
	}

	response = new(schema.CourseGroupV1Response).FromModel(user, group)
	return
}
