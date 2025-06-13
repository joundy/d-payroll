package dto

import (
	"d-payroll/entity"
	"time"
)

type CreateUserInfoBodyDto struct {
	MonthlySalary *int `json:"monthly_salary" validate:"required"`
}

type CreateUserBodyDto struct {
	Username string                 `json:"username" validate:"required"`
	Password string                 `json:"password" validate:"required"`
	Role     string                 `json:"role" validate:"required,oneof=ADMIN EMPLOYEE"`
	UserInfo *CreateUserInfoBodyDto `json:"user_info"`
}

func (c *CreateUserBodyDto) ToUserEntity() *entity.User {
	var userInfo *entity.UserInfo
	if c.UserInfo != nil {
		userInfo = &entity.UserInfo{
			MonthlySalary: c.UserInfo.MonthlySalary,
		}
	}
	return &entity.User{
		Username: c.Username,
		Password: c.Password,
		Role:     entity.UserRole(c.Role),
		UserInfo: userInfo,
	}
}

type userInfoDto struct {
	MonthlySalary *int `json:"monthly_salary"`
}

type userResponseDto struct {
	Id        *uint        `json:"id,omitempty"`
	Username  string       `json:"username"`
	Role      string       `json:"role"`
	UserInfo  *userInfoDto `json:"user_info,omitempty"`
	CreatedAt *time.Time   `json:"created_at,omitempty"`
	UpdatedAt *time.Time   `json:"updated_at,omitempty"`
}

func (r *userResponseDto) fromUserEntity(user *entity.User) {
	r.Id = user.Id
	r.Username = user.Username
	r.Role = string(user.Role)
	if user.UserInfo != nil {
		r.UserInfo = &userInfoDto{
			MonthlySalary: user.UserInfo.MonthlySalary,
		}
	}

	if user.CreatedAt != nil {
		r.CreatedAt = user.CreatedAt
	}

	if user.UpdatedAt != nil {
		r.UpdatedAt = user.UpdatedAt
	}
}

type CreateUserResponseDto userResponseDto
type GetUserByIdResponseDto userResponseDto

func (c *CreateUserResponseDto) FromUserEntity(user *entity.User) {
	(*userResponseDto)(c).fromUserEntity(user)
}

func (g *GetUserByIdResponseDto) FromUserEntity(user *entity.User) {
	(*userResponseDto)(g).fromUserEntity(user)
}
