package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/opentreehole/backend/internal/model"
)

type ReviewRepository interface {
	Repository

	FindReviewsByCourseIDs(ctx context.Context, courseIDs []int, condition func(db *gorm.DB) *gorm.DB) (reviews []*model.Review, err error)

	FindReviews(ctx context.Context, condition func(db *gorm.DB) *gorm.DB) (reviews []*model.Review, err error)

	GetReviewByID(ctx context.Context, id int) (review *model.Review, err error)

	GetReview(ctx context.Context, condition func(tx *gorm.DB) *gorm.DB) (review *model.Review, err error)

	FindReviewVotes(ctx context.Context, reviewIDs []int, userIDs []int) (votes []*model.ReviewVote, err error)

	CreateReview(ctx context.Context, review *model.Review) (err error)

	UpdateReview(ctx context.Context, userID int, oldReview *model.Review, newReview *model.Review) (err error)
}

type reviewRepository struct {
	Repository
}

func NewReviewRepository(repository Repository) ReviewRepository {
	return &reviewRepository{Repository: repository}
}

/* 接口实现 */

func (r *reviewRepository) FindReviewsByCourseIDs(
	ctx context.Context,
	courseIDs []int,
	condition func(db *gorm.DB) *gorm.DB,
) (
	reviews []*model.Review,
	err error,
) {
	reviews = make([]*model.Review, 0, 5)
	err = condition(r.GetDB(ctx).Where("course_id IN ?", courseIDs)).Find(&reviews).Error
	return
}

func (r *reviewRepository) FindReviewVotes(ctx context.Context, reviewIDs []int, userIDs []int) (votes []*model.ReviewVote, err error) {
	votes = make([]*model.ReviewVote, 0, 5)
	if len(reviewIDs) == 0 && len(userIDs) == 0 {
		return
	}
	db := r.GetDB(ctx)
	if len(reviewIDs) > 0 {
		db = db.Where("review_id IN ?", reviewIDs)
	}
	if len(userIDs) > 0 {
		db = db.Where("user_id IN ?", userIDs)
	}
	err = db.Find(&votes).Error
	return
}

func (r *reviewRepository) FindReviews(ctx context.Context, condition func(db *gorm.DB) *gorm.DB) (reviews []*model.Review, err error) {
	reviews = make([]*model.Review, 0, 5)
	err = condition(r.GetDB(ctx)).Preload("History").
		Preload("UserAchievements.Achievement").Find(&reviews).Error
	return
}

func (r *reviewRepository) GetReview(ctx context.Context, condition func(tx *gorm.DB) *gorm.DB) (review *model.Review, err error) {
	review = new(model.Review)
	err = condition(r.GetDB(ctx)).Preload("History").
		Preload("UserAchievements.Achievement").First(review).Error
	return
}

func (r *reviewRepository) CreateReview(ctx context.Context, review *model.Review) (err error) {
	return r.GetDB(ctx).Create(review).Error
}

func (r *reviewRepository) GetReviewByID(ctx context.Context, id int) (review *model.Review, err error) {
	review = new(model.Review)
	err = r.GetDB(ctx).Preload("History").
		Preload("UserAchievements.Achievement").First(review, id).Error
	return
}

func (r *reviewRepository) UpdateReview(ctx context.Context, userID int, oldReview *model.Review, newReview *model.Review) (err error) {
	// 存储到 review_history 中
	err = r.GetDB(ctx).Create(&model.ReviewHistory{
		ReviewID: oldReview.ID,
		Title:    oldReview.Title,
		Content:  oldReview.Content,
		AlterBy:  userID,
	}).Error
	if err != nil {
		return
	}

	// 更新 review
	return r.GetDB(ctx).Model(oldReview).
		Select("Title", "Content", "Rank").Updates(newReview).Error
}
