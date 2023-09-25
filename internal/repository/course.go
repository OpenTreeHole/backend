package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/opentreehole/backend/internal/model"
)

type CourseRepository interface {
	Repository

	FindCourseByID(ctx context.Context, id int, options ...FindCourseOption) (course *model.Course, err error)
	FindCoursesByGroupID(ctx context.Context, groupID int) (courses []*model.Course, err error)
	CreateCourse(ctx context.Context, course *model.Course) (err error)
}

type courseRepository struct {
	Repository
}

func NewCourseRepository(repository Repository) CourseRepository {
	return &courseRepository{Repository: repository}
}

/* 接口选项 */

type findCourseOptions struct {
	PreloadFuncs []func(db *gorm.DB) *gorm.DB
}

type FindCourseOption func(*findCourseOptions)

func WithCourseReviews() FindCourseOption {
	return func(o *findCourseOptions) {
		o.PreloadFuncs = append(o.PreloadFuncs, func(db *gorm.DB) *gorm.DB {
			return db.Preload("Reviews")
		})
	}
}

/* 接口实现 */

func (r *courseRepository) FindCourseByID(ctx context.Context, id int, options ...FindCourseOption) (course *model.Course, err error) {
	var opts findCourseOptions
	for _, option := range options {
		option(&opts)
	}
	course = new(model.Course)
	db := r.GetDB(ctx).Where("id = ?", id)
	for _, option := range opts.PreloadFuncs {
		db = option(db)
	}
	err = db.First(course).Error
	return
}

func (r *courseRepository) FindCoursesByGroupID(ctx context.Context, groupID int) (courses []*model.Course, err error) {
	courses = make([]*model.Course, 5)
	err = r.GetDB(ctx).Where("course_group_id = ?", groupID).Find(&courses).Error
	return
}

func (r *courseRepository) CreateCourse(ctx context.Context, course *model.Course) (err error) {
	return r.GetDB(ctx).Create(course).Error
}
