package service

import (
	"context"

	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/schema"
)

type AccountService interface {
	Login(
		ctx context.Context,
		email, password string,
	) (
		response *schema.TokenResponse,
		err error,
	)
	Register(ctx context.Context, email, password, verificationCode string) (response *schema.TokenResponse, err error)
	ResetPassword(ctx context.Context, email, password, verificationCode string) (response *schema.TokenResponse, err error)
}

type accountService struct {
	*Service
	repository repository.AccountRepository
}

func NewAccountService(service *Service, repository repository.AccountRepository) AccountService {
	return &accountService{Service: service, repository: repository}
}

func (a *accountService) Login(
	ctx context.Context,
	email, password string,
) (
	response *schema.TokenResponse,
	err error,
) {
	//TODO implement me
	panic("implement me")
}

func (a *accountService) Register(ctx context.Context, email, password, verificationCode string) (response *schema.TokenResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *accountService) ResetPassword(ctx context.Context, email, password, verificationCode string) (response *schema.TokenResponse, err error) {
	//TODO implement me
	panic("implement me")
}
