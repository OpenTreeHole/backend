package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/opentreehole/backend/internal/model"
	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/schema"
)

type ReviewService interface {
	Service

	CreateReview(
		ctx context.Context,
		user *model.User,
		courseID int,
		request *schema.CreateReviewV1Request,
	) (
		response *schema.ReviewV1Response,
		err error,
	)

	ModifyReview(
		ctx context.Context,
		user *model.User,
		courseID int,
		request *schema.ModifyReviewV1Request,
	) (
		response *schema.ReviewV1Response,
		err error,
	)
}

type reviewService struct {
	Service
	reviewRepository repository.ReviewRepository
	courseRepository repository.CourseRepository
}

func NewReviewService(
	service Service,
	reviewRepository repository.ReviewRepository,
	courseRepository repository.CourseRepository,
) ReviewService {
	return &reviewService{
		Service:          service,
		reviewRepository: reviewRepository,
		courseRepository: courseRepository,
	}
}

func (s *reviewService) CreateReview(
	ctx context.Context,
	user *model.User,
	courseID int,
	request *schema.CreateReviewV1Request,
) (
	response *schema.ReviewV1Response,
	err error,
) {
	// 查找 course
	_, err = s.courseRepository.FindCourseByID(ctx, courseID)
	if err != nil {
		return
	}

	// 防止重复创建
	reviews, err := s.reviewRepository.FindReviews(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("course_id = ? AND user_id = ?", courseID, user.ID)
	})
	if err != nil {
		return
	}
	if len(reviews) > 0 {
		return nil, schema.BadRequest("You cannot post more than one review.")
	}

	// 创建 review
	review := request.ToModel(user.ID, courseID)
	err = s.reviewRepository.CreateReview(ctx, review)

	// 重新加载 review
	review, err = s.reviewRepository.GetReviewByID(ctx, review.ID)

	// 创建 response
	return new(schema.ReviewV1Response).FromModel(user, review, nil), nil
}

func (s *reviewService) ModifyReview(
	ctx context.Context,
	user *model.User,
	reviewID int,
	request *schema.ModifyReviewV1Request,
) (
	response *schema.ReviewV1Response,
	err error,
) {
	// 查找 review
	review, err := s.reviewRepository.GetReviewByID(ctx, reviewID)
	if err != nil {
		return
	}

	// 验证权限
	if review.ReviewerID != user.ID && !user.IsAdmin {
		return nil, schema.Forbidden()
	}

	// 修改 review
	err = s.reviewRepository.UpdateReview(
		ctx,
		user.ID,
		review,
		request.ToModel(review.ReviewerID, review.CourseID),
	)
	if err != nil {
		return
	}

	// 查找 review
	review, err = s.reviewRepository.GetReviewByID(ctx, reviewID)
	if err != nil {
		return
	}

	// 创建 response
	return new(schema.ReviewV1Response).FromModel(user, review, nil), nil
}
