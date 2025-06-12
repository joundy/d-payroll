package models

import (
	"gorm.io/gorm"
)

type UserRole string

// TODO: proper enum
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

type UserInfo struct {
	ID uint `gorm:"primarykey"`

	UserId        uint
	MonthlySalary *int
}
