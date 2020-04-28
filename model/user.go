package model

import (
	"database/sql"
)

const (
	UserFetchByEmail = "email"
	UserFetchById    = "userId"

	UserHandlerFuncUrl         = "/user"
	UserHandlerFetchByEmailUrl = UserHandlerFuncUrl + "/{" + UserFetchByEmail + "}"
	UserHandlerFetchByIdUrl    = UserHandlerFuncUrl + "/id/{" + UserFetchById + "}"
	UserHandlerCreateUserBulkUrl= UserHandlerFuncUrl + "/bulk"
)

type User struct {
	Id              int64   `json:"id"`
	DbId            int64   `json:"-"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email"`
	Gender          string  `json:"gender"`
	PersonalAddress Address `json:`
}

type Address struct {
	UserId       int64  `json:"-"`
	Country      string `json:"first_name"`
	State        string `json:"state"`
	City         string `json:"city"`
	AddressLine1 string `json:"address"`
}

func (u *User) CreateUser(db *sql.DB) error {
	sqlStatement := `INSERT INTO awesome_user(user_id, first_name, last_name, email, gender) VALUES ($1,$2,$3,$4,$5)`
	_, err := db.Exec(sqlStatement, u.Id, u.FirstName, u.LastName, u.Email, u.Gender)

	u.FetchUserByEmail(db)
	return err
}

func (u *User) FetchUserById(db *sql.DB) error {
	sqlStatement := `select dbId, user_id, first_name, last_name, email, gender from awesome_user u where u.user_id = $1`
	return db.QueryRow(sqlStatement, u.Id).Scan(&u.DbId, &u.Id, &u.FirstName, &u.LastName, &u.Email, &u.Gender)
}

func (u *User) FetchUserByEmail(db *sql.DB) error {
	sqlStatement := `select dbId, user_id, first_name, last_name, email, gender from awesome_user u where u.email = $1`
	return db.QueryRow(sqlStatement, u.Email).Scan(&u.DbId, &u.Id, &u.FirstName, &u.LastName, &u.Email, &u.Gender)
}

func FetchAllUsers(db *sql.DB) ([]User, error) {
	sqlStatement := `select dbId, user_id, first_name, last_name, email, gender from awesome_user u`
	rows, err := db.Query(sqlStatement)

	if err != nil{
		return nil, err
	}
	users := make([]User,10,100)
	for rows.Next(){
		var u User
		if err := rows.Scan( &u.DbId, &u.Id, &u.FirstName, &u.LastName, &u.Email, &u.Gender ); err != nil {
			return nil, err
		}
		users = append( users, u )
	}

	return users, nil
}
