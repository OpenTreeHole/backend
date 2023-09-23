package repository

type CourseRepository interface {
	Repository
}

type courseRepository struct {
	Repository
}

func NewCourseRepository(repository Repository) CourseRepository {
	return &courseRepository{Repository: repository}
}

/* 接口实现 */
