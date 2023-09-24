package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/opentreehole/backend/internal/model"
)

type ReviewRepository interface {
	Repository

	FindReviewsByCourseID(ctx context.Context, courseID int, options ...FindReviewsByCourseIDOption) (reviews []*model.Review, err error)

	FindReviewsByCourseIDs(ctx context.Context, courseIDs []int, options ...FindReviewsByCourseIDOption) (reviews []*model.Review, err error)

	FindReviewVotes(ctx context.Context, reviewIDs []int, userIDs []int) (votes []*model.ReviewVote, err error)
}

type reviewRepository struct {
	Repository
}

func NewReviewRepository(repository Repository) ReviewRepository {
	return &reviewRepository{Repository: repository}
}

/* 接口实现 */

type findReviewsByCourseIDOptions struct {
	PreloadFuncs []func(db *gorm.DB) *gorm.DB
}

type FindReviewsByCourseIDOption func(*findReviewsByCourseIDOptions)

func ReviewWithHistory() FindReviewsByCourseIDOption {
	return func(o *findReviewsByCourseIDOptions) {
		o.PreloadFuncs = append(o.PreloadFuncs, func(db *gorm.DB) *gorm.DB {
			return db.Preload("History")
		})
	}
}

func ReviewWithUserAchievements() FindReviewsByCourseIDOption {
	return func(o *findReviewsByCourseIDOptions) {
		o.PreloadFuncs = append(o.PreloadFuncs, func(db *gorm.DB) *gorm.DB {
			return db.Preload("UserAchievements.Achievement")
		})
	}
}

func (r *reviewRepository) FindReviewsByCourseID(ctx context.Context, courseID int, options ...FindReviewsByCourseIDOption) (reviews []*model.Review, err error) {
	var opts findReviewsByCourseIDOptions
	for _, option := range options {
		option(&opts)
	}
	reviews = make([]*model.Review, 0, 5)
	db := r.GetDB(ctx).Where("course_id = ?", courseID)
	for _, preloadFunc := range opts.PreloadFuncs {
		db = preloadFunc(db)
	}
	err = db.Find(&reviews).Error
	return
}

func (r *reviewRepository) FindReviewsByCourseIDs(ctx context.Context, courseIDs []int, options ...FindReviewsByCourseIDOption) (reviews []*model.Review, err error) {
	var opts findReviewsByCourseIDOptions
	for _, option := range options {
		option(&opts)
	}
	reviews = make([]*model.Review, 0, 5)
	db := r.GetDB(ctx).Where("course_id IN ?", courseIDs)
	for _, preloadFunc := range opts.PreloadFuncs {
		db = preloadFunc(db)
	}
	err = db.Find(&reviews).Error
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
