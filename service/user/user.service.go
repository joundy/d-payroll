package userservice

import (
	"d-payroll/entity"
	repository "d-payroll/repository/db"
	"d-payroll/repository/db/models"
)

type UserService interface {
	CreateUser(user *entity.User) error
	CreateUsers(users []*entity.User) error
}

type userService struct {
	userDB repository.UserDB
}

func NewUserService(userDB repository.UserDB) UserService {
	return &userService{userDB: userDB}
}

// TODO: code duplications
func (s *userService) CreateUser(user *entity.User) error {
	err := user.HashPassword()
	if err != nil {
		return err
	}

	var userInfoModel models.UserInfo
	if user.UserInfo != nil {
		userInfoModel = models.UserInfo{
			MonthlySalary: user.UserInfo.MonthlySalary,
		}
	}

	userModel := models.User{
		Username: user.Username,
		Password: user.Password,
		Role:     models.UserRole(user.Role),
		UserInfo: &userInfoModel,
	}

	err = s.userDB.CreateUser(&userModel)

	return err
}

func (s *userService) CreateUsers(users []*entity.User) error {
	userModels := make([]*models.User, len(users))

	for i, user := range users {
		err := user.HashPassword()
		if err != nil {
			return err
		}

		var userInfoModel models.UserInfo
		if user.UserInfo != nil {
			userInfoModel = models.UserInfo{
				MonthlySalary: user.UserInfo.MonthlySalary,
			}
		}

		userModel := models.User{
			Username: user.Username,
			Password: user.Password,
			Role:     models.UserRole(user.Role),
			UserInfo: &userInfoModel,
		}

		userModels[i] = &userModel
	}

	return s.userDB.CreateUsers(userModels)
}
