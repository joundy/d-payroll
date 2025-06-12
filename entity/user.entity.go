package entity

import "golang.org/x/crypto/bcrypt"

type UserRole string

const (
	UserRoleAdmin    UserRole = "ADMIN"
	UserRoleEmployee UserRole = "EMPLOYEE"
)

type User struct {
	Username string `json:"username"`
	Password string
	Role     UserRole `json:"role"`

	UserInfo *UserInfo `json:"user_info"`
}

type UserInfo struct {
	MonthlySalary *int `json:"monthly_salary"`
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
