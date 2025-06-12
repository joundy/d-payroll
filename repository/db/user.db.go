package repository

import (
	"d-payroll/repository/db/models"

	"gorm.io/gorm"
)

type UserDB interface {
	CreateUser(users *models.User) error
	CreateUsers(users []*models.User) error
}

type userDB struct {
	DB *gorm.DB
}

func NewUserDB(db *gorm.DB) UserDB {
	return &userDB{DB: db}
}

func (e *userDB) CreateUser(user *models.User) error {
	return e.DB.Create(user).Error
}

func (e *userDB) CreateUsers(users []*models.User) error {
	return e.DB.Create(users).Error
}
