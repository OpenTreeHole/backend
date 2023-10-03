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

	VoteForReviewV1(
		ctx context.Context,
		user *model.User,
		reviewID int,
		request *schema.VoteForReviewV1Request,
	) (
		response *schema.ReviewV1Response,
		err error,
	)

	ListMyReviewsV1(
		ctx context.Context,
		user *model.User,
	) (
		response []*schema.MyReviewV1Response,
		err error,
	)

	GetRandomReviewV1(
		ctx context.Context,
		user *model.User,
	) (
		response *schema.RandomReviewV1Response,
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
	course, err := s.courseRepository.FindCourseByID(ctx, courseID)
	if err != nil {
		return
	}

	// 防止重复创建
	reviews, err := s.reviewRepository.FindReviews(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("course_id = ? AND reviewer_id = ?", courseID, user.ID)
	})
	if err != nil {
		return
	}
	if len(reviews) > 0 {
		return nil, schema.BadRequest("You cannot post more than one review.")
	}

	// 创建 review
	review := request.ToModel(user.ID, courseID)
	review.Course = course
	err = s.reviewRepository.CreateReview(ctx, review)
	if err != nil {
		return
	}

	// 重新加载 review
	review, err = s.reviewRepository.GetReviewByID(ctx, review.ID)
	if err != nil {
		return
	}

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

	// 加载 review_vote
	votes, err := s.reviewRepository.FindReviewVotes(ctx, []int{reviewID}, []int{user.ID})
	if err != nil {
		return
	}
	votesMap := map[int]*model.ReviewVote{review.ID: {Data: 0}}
	if len(votes) > 0 {
		votesMap[review.ID] = votes[0]
	}

	// 创建 response
	return new(schema.ReviewV1Response).FromModel(user, review, votesMap), nil
}

func (s *reviewService) VoteForReviewV1(
	ctx context.Context,
	user *model.User,
	reviewID int,
	request *schema.VoteForReviewV1Request,
) (
	response *schema.ReviewV1Response,
	err error,
) {
	// 查找 review
	review, err := s.reviewRepository.GetReviewByID(ctx, reviewID)
	if err != nil {
		return
	}

	// 获取用户对 review 的投票
	votes, err := s.reviewRepository.FindReviewVotes(ctx, []int{reviewID}, []int{user.ID})
	if err != nil {
		return
	}

	var newVote = 0
	if request.Upvote {
		newVote = 1
	} else {
		newVote = -1
	}

	if len(votes) > 0 {
		// 更新投票
		vote := votes[0]
		if (vote.Data == 1 && request.Upvote) || (vote.Data == -1 && !request.Upvote) {
			newVote = 0
		}
	}
	err = s.reviewRepository.UpdateReviewVote(ctx, user.ID, review, newVote)
	if err != nil {
		return
	}

	// 查找 review
	review, err = s.reviewRepository.GetReviewByID(ctx, reviewID)
	if err != nil {
		return
	}

	votesMap := map[int]*model.ReviewVote{review.ID: {Data: newVote}}

	// 创建 response
	return new(schema.ReviewV1Response).FromModel(user, review, votesMap), nil

}

func (s *reviewService) ListMyReviewsV1(
	ctx context.Context,
	user *model.User,
) (
	response []*schema.MyReviewV1Response,
	err error,
) {
	// 查找 review
	reviews, err := s.reviewRepository.FindReviews(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Preload("Course").Where("reviewer_id = ?", user.ID)
	})
	if err != nil {
		return
	}

	// 加载 review_vote
	reviewIDs := make([]int, 0, len(reviews))
	for _, review := range reviews {
		reviewIDs = append(reviewIDs, review.ID)
	}
	votes, err := s.reviewRepository.FindReviewVotes(ctx, reviewIDs, []int{user.ID})
	if err != nil {
		return
	}
	votesMap := make(map[int]*model.ReviewVote)
	for _, vote := range votes {
		votesMap[vote.ReviewID] = vote
	}

	// 创建 response
	response = make([]*schema.MyReviewV1Response, 0, len(reviews))
	for _, review := range reviews {
		response = append(response, new(schema.MyReviewV1Response).FromModel(review, votesMap))
	}

	return
}

func (s *reviewService) GetRandomReviewV1(
	ctx context.Context,
	user *model.User,
) (
	response *schema.RandomReviewV1Response,
	err error,
) {
	reviews, err := s.reviewRepository.FindReviews(ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Preload("Course")
		if db.Dialector.Name() == "mysql" {
			return db.Joins(`JOIN (SELECT ROUND(RAND() * ((SELECT MAX(id) FROM review) - (SELECT MIN(id) FROM review)) + (SELECT MIN(id) FROM review)) AS id) AS number_table`).
				Where("review.id >= number_table.id").Limit(1)
		} else {
			return db.Order("RANDOM()").Limit(1)
		}
	})
	if err != nil {
		return
	}

	if len(reviews) == 0 {
		return nil, schema.InternalServerError("Unable to fetch a random review. Retry later.")
	}

	review := reviews[0]

	// 加载 review_vote
	votes, err := s.reviewRepository.FindReviewVotes(ctx, []int{user.ID}, []int{review.ReviewerID})
	if err != nil {
		return
	}
	votesMap := map[int]*model.ReviewVote{review.ID: {Data: 0}}
	if len(votes) > 0 {
		votesMap[review.ID] = votes[0]
	}

	// 创建 response
	return new(schema.RandomReviewV1Response).FromModel(review, votesMap), nil
}
