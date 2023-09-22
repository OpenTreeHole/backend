package repository

type DivisionRepository interface {
	// TODO:
}

type divisionRepository struct {
	Repository
}

func NewDivisionRepository(repository Repository) DivisionRepository {
	return &divisionRepository{Repository: repository}
}
