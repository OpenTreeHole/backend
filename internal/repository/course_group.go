package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/opentreehole/backend/internal/model"
)

type CourseGroupRepository interface {
	Repository

	FindAllGroups(ctx context.Context, options ...FindGroupOption) (groups []*model.CourseGroup, err error)
	FindGroupByID(ctx context.Context, id int, options ...FindGroupOption) (group *model.CourseGroup, err error)
	FindGroupByCode(ctx context.Context, code string, options ...FindGroupOption) (group *model.CourseGroup, err error)
	CreateGroup(ctx context.Context, group *model.CourseGroup) (err error)
}

type courseGroupRepository struct {
	Repository
}

func NewCourseGroupRepository(repository Repository) CourseGroupRepository {
	return &courseGroupRepository{Repository: repository}
}

/* 接口选项类型 */
type findGroupOptions struct {
	PreloadFuncs []func(db *gorm.DB) *gorm.DB
}

type FindGroupOption func(*findGroupOptions)

func WithGroupCourses() FindGroupOption {
	return func(o *findGroupOptions) {
		o.PreloadFuncs = append(o.PreloadFuncs, func(db *gorm.DB) *gorm.DB {
			return db.Preload("Courses")
		})
	}
}

func WithGroupReviews() FindGroupOption {
	return func(o *findGroupOptions) {
		o.PreloadFuncs = append(o.PreloadFuncs, func(db *gorm.DB) *gorm.DB {
			return db.Preload("Courses.Reviews")
		})
	}
}

func WithReviewsAndHistory() FindGroupOption {
	return func(o *findGroupOptions) {
		o.PreloadFuncs = append(o.PreloadFuncs, func(db *gorm.DB) *gorm.DB {
			return db.Preload("Courses.Reviews.History")
		})
	}
}

func WithReviewsUserAchievements() FindGroupOption {
	return func(o *findGroupOptions) {
		o.PreloadFuncs = append(o.PreloadFuncs, func(db *gorm.DB) *gorm.DB {
			return db.Preload("Courses.Reviews.UserAchievements.Achievement")
		})
	}
}

/* 接口实现 */

func (r *courseGroupRepository) FindAllGroups(ctx context.Context, options ...FindGroupOption) (groups []*model.CourseGroup, err error) {
	var option findGroupOptions
	for _, opt := range options {
		opt(&option)
	}
	groups = make([]*model.CourseGroup, 5)
	db := r.GetDB(ctx)
	for _, f := range option.PreloadFuncs {
		db = f(db)
	}
	err = db.Find(&groups).Error
	return
}

func (r *courseGroupRepository) FindGroupByID(ctx context.Context, id int, options ...FindGroupOption) (group *model.CourseGroup, err error) {
	var option findGroupOptions
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

func (r *courseGroupRepository) FindGroupByCode(ctx context.Context, code string, options ...FindGroupOption) (group *model.CourseGroup, err error) {
	var option findGroupOptions
	for _, opt := range options {
		opt(&option)
	}
	group = &model.CourseGroup{}
	db := r.GetDB(ctx)
	for _, f := range option.PreloadFuncs {
		db = f(db)
	}
	err = db.Where("code = ?", code).First(group).Error
	return
}

func (r *courseGroupRepository) CreateGroup(ctx context.Context, group *model.CourseGroup) (err error) {
	return r.GetDB(ctx).Create(group).Error
}
