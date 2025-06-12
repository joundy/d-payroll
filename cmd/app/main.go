package main

import (
	"d-payroll/config"
	repository "d-payroll/repository/db"
)

func main() {
	config := config.NewConfig()
	db, err := repository.NewDBHelper(*config)
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
