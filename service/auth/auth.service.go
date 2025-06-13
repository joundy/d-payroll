package authservice

import (
	"context"
	"d-payroll/config"
	"d-payroll/entity"
	internalerror "d-payroll/internal-error"
	userservice "d-payroll/service/user"
	"d-payroll/utils"
)

type AuthService interface {
	Login(ctx context.Context, login *entity.Login) (*entity.AuthToken, error)
}

type authService struct {
	config  *config.Config
	userSvc userservice.UserService
}

func NewAuthService(config *config.Config, userSvc userservice.UserService) AuthService {
	return &authService{config: config, userSvc: userSvc}
}

func (a *authService) Login(ctx context.Context, login *entity.Login) (*entity.AuthToken, error) {
	user, err := a.userSvc.GetUserByUsername(ctx, login.Username)
	if err != nil {
		return nil, err
	}

	if !user.VerifyPassword(login.Password) {
		return nil, &internalerror.InvalidCredentialsError{}
	}

	token, err := utils.GenerateToken(a.config.Auth.JwtSecret, &entity.AuthTokenPayload{
		Id:   *user.Id,
		Role: user.Role,
	})
	if err != nil {
		return nil, err
	}

	return &entity.AuthToken{Token: token}, nil
}
