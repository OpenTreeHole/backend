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
	reviewRepository      repository.ReviewRepository
}

func NewCourseGroupService(
	service Service,
	courseGroupRepository repository.CourseGroupRepository,
	reviewRepository repository.ReviewRepository,
) CourseGroupService {
	return &courseGroupService{
		Service:               service,
		courseGroupRepository: courseGroupRepository,
		reviewRepository:      reviewRepository,
	}
}

func (c *courseGroupService) GetGroupByIDV1(ctx context.Context, user *model.User, id int) (response *schema.CourseGroupV1Response, err error) {
	// 获取课程组，同时加载课程
	// 这里不预加载课程的评论，因为评论作为动态的数据，应该独立作缓存，提高缓存粒度和缓存更新频率
	group, err := c.courseGroupRepository.FindGroupByID(ctx, id, repository.WithGroupCourses())
	if err != nil {
		return nil, err
	}

	// 获取课程组的所有课程的所有评论，同时加载评论的历史记录和用户成就
	courseIDs := make([]int, len(group.Courses))
	for i, course := range group.Courses {
		courseIDs[i] = course.ID
	}
	reviews, err := c.reviewRepository.FindReviewsByCourseIDs(ctx, courseIDs, repository.ReviewWithHistory(), repository.ReviewWithUserAchievements())
	if err != nil {
		return nil, err
	}

	// 将评论按照课程分组
	reviewsMap := make(map[int][]*model.Review)
	for _, review := range reviews {
		if reviewsMap[review.CourseID] == nil {
			reviewsMap[review.CourseID] = make([]*model.Review, 0)
		}
		reviewsMap[review.CourseID] = append(reviewsMap[review.CourseID], review)
	}

	// 将评论分组放入课程中
	for _, course := range group.Courses {
		course.Reviews = reviewsMap[course.ID]
	}

	// 获取课程组的所有课程的所有评论的自己的投票
	reviewIDs := make([]int, 0)
	for _, review := range reviews {
		reviewIDs = append(reviewIDs, review.ID)
	}
	votes, err := c.reviewRepository.FindReviewVotes(ctx, reviewIDs, []int{user.ID})
	if err != nil {
		return nil, err
	}

	// 转换为 map, 便于查找
	votesMap := make(map[int]*model.ReviewVote)
	for _, vote := range votes {
		votesMap[vote.ReviewID] = vote
	}

	// 将课程组转换为响应
	response = new(schema.CourseGroupV1Response).FromModel(user, group, votesMap)
	return
}
