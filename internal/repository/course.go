package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/opentreehole/backend/internal/model"
)

type CourseRepository interface {
	Repository

	FindCourseByID(ctx context.Context, id int, conditions ...func(db *gorm.DB) *gorm.DB) (course *model.Course, err error)
	FindCoursesByGroupID(ctx context.Context, groupID int) (courses []*model.Course, err error)
	CreateCourse(ctx context.Context, course *model.Course) (err error)
}

type courseRepository struct {
	Repository
}

func NewCourseRepository(repository Repository) CourseRepository {
	return &courseRepository{Repository: repository}
}

/* 接口实现 */

func (r *courseRepository) FindCourseByID(ctx context.Context, id int, conditions ...func(db *gorm.DB) *gorm.DB) (course *model.Course, err error) {
	course = new(model.Course)
	db := r.GetDB(ctx)
	for _, condition := range conditions {
		condition(db)
	}
	err = db.First(course, id).Error
	return
}

func (r *courseRepository) FindCoursesByGroupID(ctx context.Context, groupID int) (courses []*model.Course, err error) {
	courses = make([]*model.Course, 5)
	err = r.GetDB(ctx).Where("course_group_id = ?", groupID).Find(&courses).Error
	return
}

func (r *courseRepository) CreateCourse(ctx context.Context, course *model.Course) (err error) {
	err = r.GetDB(ctx).Create(course).Error

	// clear cache
	return r.GetCache(ctx).Delete(ctx, "danke:course_group")
}
