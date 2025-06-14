package main

import (
	"d-payroll/config"
	"d-payroll/controller/http"
	repository "d-payroll/repository/db"
	attendanceservice "d-payroll/service/attendance"
	authservice "d-payroll/service/auth"
	overtimeservice "d-payroll/service/overtime"
	reimbursementservice "d-payroll/service/reimbursement"
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
	attendanceDB := repository.NewAttendanceDB(db.DB)
	reimbursementDB := repository.NewReimbursementDB(db.DB)
	overtimeDB := repository.NewOvertimeDB(db.DB)

	// services

	userSvc := userservice.NewUserService(userDB)
	authSvc := authservice.NewAuthService(config, userSvc)
	attendanceSvc := attendanceservice.NewAttendanceService(attendanceDB)
	reimbursementSvc := reimbursementservice.NewReimbursementService(reimbursementDB)
	overtimeSvc := overtimeservice.NewOvertimeService(config, overtimeDB, attendanceSvc)

	// deliveries http

	httpApp := http.NewHttpApp(config)

	http.NewUserHttp(httpApp, userSvc)
	http.NewAuthHttp(httpApp, authSvc)
	http.NewAttendanceHttp(httpApp, attendanceSvc)
	http.NewReimbursementHttp(httpApp, reimbursementSvc)
	http.NewOvertimeHttp(httpApp, overtimeSvc)

	httpApp.Listen()
}
