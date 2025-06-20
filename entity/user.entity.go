package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleEmployee UserRole = "EMPLOYEE"
)

type User struct {
	Id       *uint
	Username string
	Password string
	Role     UserRole

	UserInfo *UserInfo

	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type UserInfo struct {
	MonthlySalary *int
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
