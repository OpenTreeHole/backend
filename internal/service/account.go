package service

import (
	"context"

	"github.com/opentreehole/backend/internal/repository"
	"github.com/opentreehole/backend/internal/schema"
)

type AccountService interface {
	Service
	Login(
		ctx context.Context,
		email, password string,
	) (
		response *schema.TokenResponse,
		err error,
	)
	Register(
		ctx context.Context,
		email, password, verificationCode string,
	) (
		response *schema.TokenResponse,
		err error,
	)
	ResetPassword(
		ctx context.Context,
		email, password, verificationCode string,
	) (
		response *schema.TokenResponse,
		err error,
	)
}

type accountService struct {
	Service
	repository repository.AccountRepository
}

func NewAccountService(service Service, repository repository.AccountRepository) AccountService {
	return &accountService{Service: service, repository: repository}
}

func (a *accountService) Login(
	ctx context.Context,
	email, password string,
) (
	response *schema.TokenResponse,
	err error,
) {
	// get user by email
	user, err := a.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// check password
	err = a.repository.CheckPassword(ctx, password, user.Password)
	if err != nil {
		return nil, err
	}

	// create jwt token
	access, refresh, err := a.repository.CreateJWTToken(ctx, user)
	if err != nil {
		return nil, err
	}

	return &schema.TokenResponse{
		Access:  access,
		Refresh: refresh,
		Message: "login success", // TODO: i18n
	}, nil
}

func (a *accountService) Register(
	ctx context.Context,
	email, password, verificationCode string,
) (
	response *schema.TokenResponse,
	err error,
) {
	scope := "register"

	// check verification code
	if a.GetConfig(ctx).Features.EmailVerification {
		err = a.repository.CheckVerificationCode(ctx, scope, email, verificationCode)
		if err != nil {
			return nil, err
		}

		// delete verification code
		defer func() {
			if err == nil {
				err = a.repository.DeleteVerificationCode(ctx, scope, email)
			}
		}()
	}

	// check if user has registered
	registered, err := a.repository.CheckIfUserExists(ctx, email)
	if err != nil {
		return nil, err
	}

	// check if user has been deleted
	deleted, err := a.repository.CheckIfUserDeleted(ctx, email)
	if err != nil {
		return nil, err
	}

	if registered {
		if !deleted {
			return nil, schema.BadRequest("该用户已注册，如果忘记密码，请使用忘记密码功能找回") // TODO: 避免硬编码, i18n, ErrUserAlreadyExists
		} else {
			return nil, schema.BadRequest("注销账号后禁止注册") // TODO: 避免硬编码, i18n, ErrUserNotAllowedToRegister
		}
	}

	var access, refresh string
	// in transaction
	err = a.repository.Transaction(ctx, func(ctx context.Context) error {
		// create user
		user, err := a.repository.CreateUser(ctx, email, password)
		if err != nil {
			return err
		}

		// create shamir emails
		if a.repository.GetConfig(ctx).Features.Shamir {
			// TODO implement me
		}

		// create kong consumer
		if a.repository.GetConfig(ctx).Features.ExternalGateway {
			// TODO implement me
		}

		// create jwt token
		access, refresh, err = a.repository.CreateJWTToken(ctx, user)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// response
	return &schema.TokenResponse{
		Access:  access,
		Refresh: refresh,
		Message: "register success", // TODO: i18n
	}, nil
}

func (a *accountService) ResetPassword(
	ctx context.Context,
	email, password, verificationCode string,
) (
	response *schema.TokenResponse,
	err error,
) {
	//TODO implement me
	panic("implement me")
}
