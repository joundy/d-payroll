package main

import (
	"d-payroll/config"
	"d-payroll/controller/http"
	repository "d-payroll/repository/db"
	authservice "d-payroll/service/auth"
	userservice "d-payroll/service/user"
)

func main() {
	config := config.NewConfig()
	db, err := repository.NewDBHelper(*config)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// repositories

	userDB := repository.NewUserDB(db.DB)

	// services

	userSvc := userservice.NewUserService(userDB)
	authSvc := authservice.NewAuthService(config, userSvc)

	// deliveries http

	httpApp := http.NewHttpApp(config)

	http.NewUserHttp(httpApp, userSvc)
	http.NewAuthHttp(httpApp, authSvc)

	httpApp.Listen()
}
