package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type DataBaseConnInfo struct {
	Host       string
	Port       string
	User       string
	Password   string
	DataBase   string
	SslEnabled string
}

func (i *DataBaseConnInfo) CreateConnection() *sql.DB {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", i.Host, i.Port, i.User, i.Password, i.DataBase, i.SslEnabled)

	var err error

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Printf("Problem to open %v database connection:\n %v \n", i.DataBase, err)
		return nil
	}
	return db
}

func (i *DataBaseConnInfo) WaitUntilDataBaseIsUp() error {
	var err error
	c := 0

	for {
		time.Sleep(time.Microsecond * time.Duration(60000))
		db := i.CreateConnection()

		err = db.Ping()

		if err == nil || c < 5 {
			db.Close()
			break
		}

		log.Println("We are failing to ping the Database, attempt #", c)
		c++
	}

	if err != nil {
		log.Println("problem to ping the database server", err)
		return err
	}
	return nil
}
