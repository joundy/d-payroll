package main

import (
	"context"
	"d-payroll/cmd/seed/utils"
	"d-payroll/config"
	"d-payroll/entity"
	repository "d-payroll/repository/db"
	userservice "d-payroll/service/user"

	"github.com/go-faker/faker/v4"
)

func main() {

	config := config.NewConfig()
	db, err := repository.NewDBHelper(*config)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userDB := repository.NewUserDB(db.DB)
	userService := userservice.NewUserService(userDB)

	seedUserAdmin(config, userService)
	seedUserEmployees(userService)
}

func seedUserAdmin(config *config.Config, userSvc userservice.UserService) {
	ctx := context.Background()
	userSvc.CreateUser(ctx, &entity.User{
		Username: config.AdminUser.Username,
		Password: config.AdminUser.Password,
		Role:     entity.UserRoleAdmin,
	})
}

func seedUserEmployees(userSvc userservice.UserService) {
	ctx := context.Background()
	userEmployees := []*entity.User{}

	for i := 0; i < 100; i++ {
		monthlySalary := utils.GenerateNumberBetween(1000000, 10000000)

		userEmployees = append(userEmployees, &entity.User{
			Username: faker.Username(),
			Password: faker.Password(),
			Role:     entity.UserRoleEmployee,
			UserInfo: &entity.UserInfo{
				MonthlySalary: &monthlySalary,
			},
		})
	}

	userSvc.CreateUsers(ctx, userEmployees)
}
