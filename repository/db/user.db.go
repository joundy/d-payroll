package repository

import (
	"context"
	"d-payroll/repository/db/models"

	"gorm.io/gorm"
)

// TODO: optimize query, don't use preload use join instead
type UserDB interface {
	CreateUser(ctx context.Context, users *models.User) error
	CreateUsers(ctx context.Context, users []*models.User) error
	GetuserById(ctx context.Context, id int) (*models.User, error)
}

type userDB struct {
	DB *gorm.DB
}

func NewUserDB(db *gorm.DB) UserDB {
	return &userDB{DB: db}
}

func (e *userDB) CreateUser(ctx context.Context, user *models.User) error {
	return e.DB.WithContext(ctx).Create(user).Error
}

func (e *userDB) CreateUsers(ctx context.Context, users []*models.User) error {
	return e.DB.WithContext(ctx).Create(users).Error
}

func (e *userDB) GetuserById(ctx context.Context, id int) (*models.User, error) {
	var user models.User

	tx := e.DB.Begin()
	tx.WithContext(ctx).Preload("UserInfo").First(&user, id)
	tx.Commit()

	return &user, nil
}
