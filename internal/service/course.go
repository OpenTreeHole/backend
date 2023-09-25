package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/opentreehole/backend/internal/model"
	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/schema"
)

type CourseService interface {
	Service

	ListCoursesV1(ctx context.Context) (response []*schema.CourseGroupV1Response, err error)
	GetCourseV1(ctx context.Context, user *model.User, id int) (response *schema.CourseV1Response, err error)
	AddCourseV1(ctx context.Context, request *schema.CreateCourseV1Request) (response *schema.CourseV1Response, err error)
}

type courseService struct {
	Service
	courseGroupRepository repository.CourseGroupRepository
	courseRepository      repository.CourseRepository
	reviewRepository      repository.ReviewRepository
}

func NewCourseService(
	service Service,
	courseRepository repository.CourseRepository,
	courseGroupRepository repository.CourseGroupRepository,
	reviewRepository repository.ReviewRepository,
) CourseService {
	return &courseService{
		Service:               service,
		courseRepository:      courseRepository,
		courseGroupRepository: courseGroupRepository,
		reviewRepository:      reviewRepository,
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

func (s *courseService) GetCourseV1(ctx context.Context, user *model.User, id int) (response *schema.CourseV1Response, err error) {
	course, err := s.courseRepository.FindCourseByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 获取课程的评论，同时加载评论的历史记录和用户成就
	course.Reviews, err = s.reviewRepository.FindReviewsByCourseID(ctx, course.ID, repository.ReviewWithHistory(), repository.ReviewWithUserAchievements())

	// 获取所有评论的自己的投票
	reviewIDs := make([]int, 0)
	for _, review := range course.Reviews {
		reviewIDs = append(reviewIDs, review.ID)
	}
	votes, err := s.reviewRepository.FindReviewVotes(ctx, reviewIDs, []int{user.ID})
	if err != nil {
		return nil, err
	}

	// 转为 votesMap
	votesMap := make(map[int]*model.ReviewVote)
	for _, vote := range votes {
		votesMap[vote.ReviewID] = vote
	}

	return new(schema.CourseV1Response).FromModel(user, course, votesMap), nil
}

func (s *courseService) AddCourseV1(ctx context.Context, request *schema.CreateCourseV1Request) (response *schema.CourseV1Response, err error) {
	group, err := s.courseGroupRepository.FindGroupByCode(ctx, request.Code)
	if err != nil {
		// create group if not found
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		group = request.ToCourseGroupModel()
		err = s.courseGroupRepository.CreateGroup(ctx, group)
		if err != nil {
			return nil, err
		}
	}

	course := request.ToModel(group.ID)
	err = s.courseRepository.CreateCourse(ctx, course)
	if err != nil {
		return nil, err
	}

	return new(schema.CourseV1Response).FromModel(nil, course, nil), nil
}
