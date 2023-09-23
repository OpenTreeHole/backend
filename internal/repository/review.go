package repository

type ReviewRepository interface {
	Repository
}

type reviewRepository struct {
	Repository
}

func NewReviewRepository(repository Repository) ReviewRepository {
	return &reviewRepository{Repository: repository}
}

/* 接口实现 */
