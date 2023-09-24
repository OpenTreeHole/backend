package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/opentreehole/backend/internal/model"
)

type CourseGroupRepository interface {
	Repository

	FindGroupByID(ctx context.Context, id int, options ...FindGroupByIDOption) (group *model.CourseGroup, err error)
}

type courseGroupRepository struct {
	Repository
}

func NewCourseGroupRepository(repository Repository) CourseGroupRepository {
	return &courseGroupRepository{Repository: repository}
}

/* 接口选项类型 */
type findGroupByIdOptions struct {
	PreloadFuncs []func(db *gorm.DB) *gorm.DB
}

type FindGroupByIDOption func(*findGroupByIdOptions)

func WithGroupCourses() FindGroupByIDOption {
	return func(o *findGroupByIdOptions) {
		o.PreloadFuncs = append(o.PreloadFuncs, func(db *gorm.DB) *gorm.DB {
			return db.Preload("Courses")
		})
	}
}

func WithGroupReviews() FindGroupByIDOption {
	return func(o *findGroupByIdOptions) {
		o.PreloadFuncs = append(o.PreloadFuncs, func(db *gorm.DB) *gorm.DB {
			return db.Preload("Courses.Reviews")
		})
	}
}

func WithReviewsAndHistory() FindGroupByIDOption {
	return func(o *findGroupByIdOptions) {
		o.PreloadFuncs = append(o.PreloadFuncs, func(db *gorm.DB) *gorm.DB {
			return db.Preload("Courses.Reviews.History")
		})
	}
}

func WithReviewsUserAchievements() FindGroupByIDOption {
	return func(o *findGroupByIdOptions) {
		o.PreloadFuncs = append(o.PreloadFuncs, func(db *gorm.DB) *gorm.DB {
			return db.Preload("Courses.Reviews.UserAchievements.Achievement")
		})
	}
}

/* 接口实现 */

func (r *courseGroupRepository) FindGroupByID(ctx context.Context, id int, options ...FindGroupByIDOption) (group *model.CourseGroup, err error) {
	var option findGroupByIdOptions
	for _, opt := range options {
		opt(&option)
	}
	group = &model.CourseGroup{}
	db := r.GetDB(ctx)
	for _, f := range option.PreloadFuncs {
		db = f(db)
	}
	err = db.First(group, id).Error
	return
}
