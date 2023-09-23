package repository

import (
	"context"

	"github.com/opentreehole/backend/internal/model"
)

type CourseRepository interface {
	Repository

	FindCoursesByGroupID(ctx context.Context, groupID int) (courses []*model.Course, err error)
}

type courseRepository struct {
	Repository
}

func NewCourseRepository(repository Repository) CourseRepository {
	return &courseRepository{Repository: repository}
}

/* 接口实现 */

func (r *courseRepository) FindCoursesByGroupID(ctx context.Context, groupID int) (courses []*model.Course, err error) {
	courses = make([]*model.Course, 5)
	err = r.GetDB(ctx).Where("course_group_id = ?", groupID).Find(&courses).Error
	return
}
