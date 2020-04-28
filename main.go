package main

import (
	"awesomeProject/app"
	"awesomeProject/config"
)

const (
	APIPort = ":8080"
)

func main() {
	dbCon := config.DataBaseConnInfo{
		Host:       "localhost",
		Port:       "5432",
		User:       "golang_user",
		Password:   "golang_user",
		DataBase:   "awesomeProject",
		SslEnabled: "disable",
	}

	a := new(app.App)
	a.Initialize(dbCon)
	a.Run(APIPort)
}
