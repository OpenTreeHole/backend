package repository

import (
	"context"

	"github.com/opentreehole/backend/internal/model"
)

type ReviewRepository interface {
	Repository

	FindReviewsByCourseID(ctx context.Context, courseID int) (reviews []*model.Review, err error)

	FindReviewsByCourseIDs(ctx context.Context, courseIDs []int) (reviews []*model.Review, err error)
}

type reviewRepository struct {
	Repository
}

func NewReviewRepository(repository Repository) ReviewRepository {
	return &reviewRepository{Repository: repository}
}

/* 接口实现 */

func (r *reviewRepository) FindReviewsByCourseID(ctx context.Context, courseID int) (reviews []*model.Review, err error) {
	reviews = make([]*model.Review, 5)
	err = r.GetDB(ctx).Where("course_id = ?", courseID).Find(&reviews).Error
	return
}

func (r *reviewRepository) FindReviewsByCourseIDs(ctx context.Context, courseIDs []int) (reviews []*model.Review, err error) {
	reviews = make([]*model.Review, 5)
	err = r.GetDB(ctx).Where("course_id IN ?", courseIDs).Find(&reviews).Error
	return
}
