package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/opentreehole/backend/internal/model"
	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/schema"
)

type CourseGroupService interface {
	Service

	/* V1 */

	GetGroupByIDV1(
		ctx context.Context,
		user *model.User,
		id int,
	) (
		response *schema.CourseGroupV1Response,
		err error,
	)

	GetCourseGroupHash(
		ctx context.Context,
	) (
		response *schema.CourseGroupHashV1Response,
		err error,
	)

	RefreshCourseGroupHash(
		ctx context.Context,
	) (
		err error,
	)

	/* V3 */

	SearchCourseGroupV3(
		ctx context.Context,
		user *model.User,
		request *schema.CourseGroupSearchV3Request,
	) (
		response *schema.PagedResponse[schema.CourseGroupV3Response, any],
		err error,
	)

	GetGroupByIDV3(
		ctx context.Context,
		user *model.User,
		id int,
	) (
		response *schema.CourseGroupV3Response,
		err error,
	)
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

func (c *courseGroupService) GetGroupByIDV1(
	ctx context.Context,
	user *model.User,
	id int,
) (
	response *schema.CourseGroupV1Response,
	err error,
) {
	// 获取课程组，同时加载课程
	// 这里不预加载课程的评论，因为评论作为动态的数据，应该独立作缓存，提高缓存粒度和缓存更新频率
	group, err := c.courseGroupRepository.FindGroupByID(ctx, id, func(db *gorm.DB) *gorm.DB {
		return db.Preload("Courses")
	})
	if err != nil {
		return nil, err
	}

	// 获取课程组的所有课程的所有评论，同时加载评论的历史记录和用户成就
	courseIDs := make([]int, len(group.Courses))
	for i, course := range group.Courses {
		courseIDs[i] = course.ID
	}
	reviews, err := c.reviewRepository.FindReviews(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("course_id IN ?", courseIDs)
	})
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

func (c *courseGroupService) GetCourseGroupHash(ctx context.Context) (response *schema.CourseGroupHashV1Response, err error) {
	_, hash, err := c.courseGroupRepository.FindGroupsWithCourses(ctx, false)
	if err != nil {
		return nil, err
	}
	response = new(schema.CourseGroupHashV1Response).FromModel(hash)
	return
}

func (c *courseGroupService) RefreshCourseGroupHash(ctx context.Context) (err error) {
	_, _, err = c.courseGroupRepository.FindGroupsWithCourses(ctx, true)
	return
}

/* V3 */

func (c *courseGroupService) SearchCourseGroupV3(
	ctx context.Context,
	user *model.User,
	request *schema.CourseGroupSearchV3Request,
) (
	response *schema.PagedResponse[schema.CourseGroupV3Response, any],
	err error,
) {
	var (
		page     = request.Page
		pageSize = request.PageSize
		query    = request.Query
	)
	groups, err := c.courseGroupRepository.FindGroups(ctx, func(db *gorm.DB) *gorm.DB {
		if model.CourseCodeRegexp.MatchString(query) {
			db = db.Where("code LIKE ?", query+"%")
		} else {
			db = db.Where("name LIKE ?", "%"+query+"%")
		}
		if page > 0 {
			if pageSize == 0 {
				pageSize = 10
			}
			db = db.Limit(pageSize).Offset((page - 1) * pageSize)
		} else {
			if pageSize > 0 {
				db = db.Limit(pageSize)
			}
		}
		return db.Order("id")
	})
	if err != nil {
		return nil, err
	}

	items := make([]*schema.CourseGroupV3Response, 0, len(groups))
	for _, group := range groups {
		items = append(items, new(schema.CourseGroupV3Response).FromModel(user, group, nil))
	}
	response = &schema.PagedResponse[schema.CourseGroupV3Response, any]{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
	}
	return response, nil
}

func (c *courseGroupService) GetGroupByIDV3(
	ctx context.Context,
	user *model.User,
	id int,
) (
	response *schema.CourseGroupV3Response,
	err error,
) {
	group, err := c.courseGroupRepository.FindGroupByID(ctx, id, func(db *gorm.DB) *gorm.DB {
		return db.Preload("Courses")
	})
	if err != nil {
		return nil, err
	}

	// 获取课程组的所有课程的所有评论，同时加载评论的历史记录和用户成就
	courseIDs := make([]int, len(group.Courses))
	for i, course := range group.Courses {
		courseIDs[i] = course.ID
	}
	reviews, err := c.reviewRepository.FindReviews(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("course_id IN ?", courseIDs)
	})
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
	response = new(schema.CourseGroupV3Response).FromModel(user, group, votesMap)
	return
}
