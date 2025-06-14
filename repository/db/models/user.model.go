package models

import (
	"d-payroll/entity"
	"d-payroll/utils"

	"gorm.io/gorm"
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleEmployee UserRole = "EMPLOYEE"
)

type User struct {
	gorm.Model

	Username string
	Password string
	Role     UserRole `gorm:"type:user_role"`

	UserInfo *UserInfo `gorm:"foreignKey:UserId;references:ID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = utils.TimeNow()
	u.UpdatedAt = utils.TimeNow()

	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = utils.TimeNow()
	return
}

type UserInfo struct {
	ID uint `gorm:"primarykey"`

	UserId        uint
	MonthlySalary *int
}

func (u *User) ToUserEntity() *entity.User {
	var userInfo *entity.UserInfo
	if u.UserInfo != nil {
		userInfo = &entity.UserInfo{
			MonthlySalary: u.UserInfo.MonthlySalary,
		}
	}
	return &entity.User{
		Id:        &u.ID,
		Username:  u.Username,
		Password:  u.Password,
		Role:      entity.UserRole(u.Role),
		UserInfo:  userInfo,
		CreatedAt: &u.CreatedAt,
		UpdatedAt: &u.UpdatedAt,
	}
}

func (u *User) FromUserEntity(user *entity.User) {
	u.Username = user.Username
	u.Password = user.Password
	u.Role = UserRole(user.Role)

	if user.UserInfo != nil {
		u.UserInfo = &UserInfo{
			MonthlySalary: user.UserInfo.MonthlySalary,
		}
	}

	if user.CreatedAt != nil {
		u.CreatedAt = *user.CreatedAt
	}

	if user.UpdatedAt != nil {
		u.UpdatedAt = *user.UpdatedAt
	}
}
