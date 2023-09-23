package repository

type CourseGroupRepository interface {
	Repository
}

type courseGroupRepository struct {
	Repository
}

func NewCourseGroupRepository(repository Repository) CourseGroupRepository {
	return &courseGroupRepository{Repository: repository}
}

/* 接口实现 */
