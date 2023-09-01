package repository

import (
	"github.com/opentreehole/backend/internal/model"
)

type AccountRepository interface {
	GetUserByEmail(email string) (user *model.User, err error)
}

type accountRepository struct {
	*Repository
}

func (a *accountRepository) GetUserByEmail(email string) (user *model.User, err error) {
	//TODO implement me
	panic("implement me")
}

func NewAccountRepository(repository *Repository) AccountRepository {
	return &accountRepository{Repository: repository}
}
