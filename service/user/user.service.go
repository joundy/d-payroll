package userservice

import (
	"context"
	"d-payroll/entity"
	repository "d-payroll/repository/db"
	"d-payroll/repository/db/models"
)

type UserService interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	CreateUsers(ctx context.Context, users []*entity.User) ([]*entity.User, error)
	GetUserById(ctx context.Context, id uint) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetUserIds(ctx context.Context) ([]uint, error)
}

type userService struct {
	userDB repository.UserDB
}

func NewUserService(userDB repository.UserDB) UserService {
	return &userService{userDB: userDB}
}

func (s *userService) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := user.HashPassword()
	if err != nil {
		return nil, err
	}

	var userModel models.User
	userModel.FromUserEntity(user)

	err = s.userDB.CreateUser(ctx, &userModel)
	if err != nil {
		return nil, err
	}

	return userModel.ToUserEntity(), nil
}

func (s *userService) CreateUsers(ctx context.Context, users []*entity.User) ([]*entity.User, error) {
	userModels := make([]*models.User, len(users))
	for i, user := range users {
		if err := user.HashPassword(); err != nil {
			return nil, err
		}
		var model models.User
		model.FromUserEntity(user)
		userModels[i] = &model
	}
	if err := s.userDB.CreateUsers(ctx, userModels); err != nil {
		return nil, err
	}

	createdUsers := make([]*entity.User, len(userModels))
	for i, model := range userModels {
		createdUsers[i] = model.ToUserEntity()
	}
	return createdUsers, nil
}

func (s *userService) GetUserById(ctx context.Context, id uint) (*entity.User, error) {
	userModel, err := s.userDB.GetuserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return userModel.ToUserEntity(), nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	userModel, err := s.userDB.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return userModel.ToUserEntity(), nil
}

func (s *userService) GetUserIds(ctx context.Context) ([]uint, error) {
	return s.userDB.GetUserIds(ctx)
}
