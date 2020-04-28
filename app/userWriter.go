package app

import (
	"awesomeProject/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
	u := new(model.User)
	err := json.NewDecoder(r.Body).Decode(u)

	if err != nil {
		logError(err)
		respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if u.Id < 0 || len(u.FirstName) == 0 || len(u.LastName) == 0 || len(u.Email) == 0 || len(u.Gender) == 0 {
		respondWithError(w, http.StatusUnprocessableEntity, "Unprocessable Entity")
		return
	}

	sqlCon := a.DBInfo.CreateConnection()

	lookUByEmail := model.User{Email: u.Email}
	//The Users exist previously
	if lookUByEmail.FetchUserByEmail(sqlCon); lookUByEmail.DbId > 0  {
		m := map[string]model.User{"User Already Exists":lookUByEmail}
		respondWithJSON(w, http.StatusConflict, m)
		return
	}
	err = u.CreateUser(sqlCon)
	defer sqlCon.Close()

	if err != nil {
		logError(err)
		respondWithError(w, http.StatusServiceUnavailable, "There is a problem to register a new user")
	}

	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) getUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)[model.UserFetchByEmail]

	if len(email) == 0 {
		respondWithError(w, http.StatusBadRequest, "The email can not be null")
		return
	}

	sqlCon := a.DBInfo.CreateConnection()
	u := model.User{Email: email}
	u.FetchUserByEmail(sqlCon)
	defer sqlCon.Close()

	if u.DbId <= 0 {
		respondWithError(w, http.StatusNotFound, "There is not an user with the email "+u.Email)
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) getUserById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)[model.UserFetchById], 10, 64)

	if err != nil || id <= 0 {
		respondWithError(w, http.StatusBadRequest, "The user Id can not be null")
		return
	}

	sqlCon := a.DBInfo.CreateConnection()
	u := model.User{Id: id}
	u.FetchUserById(sqlCon)
	defer sqlCon.Close()

	if u.DbId <= 0 {
		respondWithError(w, http.StatusNotFound, fmt.Sprint("There is not an user with the ID ", u.Id))
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) getAllUsers(w http.ResponseWriter, r *http.Request) {

	sqlCon := a.DBInfo.CreateConnection()
	users, err := model.FetchAllUsers(sqlCon)
	defer sqlCon.Close()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (a *App) createUserBulk(w http.ResponseWriter, r *http.Request) {
	uSlice :=  make([]model.User,10,100 )
	err := json.NewDecoder( r.Body ).Decode( &uSlice )
	if err != nil {
		panic(err)
	}

	if err != nil {
		logError(err)
		respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if len(uSlice) == 0 {
		respondWithError(w, http.StatusUnprocessableEntity, "There are not elements to insert")
		return
	}

	sqlCon := a.DBInfo.CreateConnection()

	uError :=  map[string]model.User{}
	for i, u := range uSlice{

		//Adding to possible Data with Errors
		if u.Id < 0 || len(u.FirstName) == 0 || len(u.LastName) == 0 || len(u.Email) == 0 || len(u.Gender) == 0 {
			k := fmt.Sprint("Json ", (i+1), " has Invalid Data")
			uError[k] = u
			continue
		}

		lookU := model.User{Email: u.Email}

		if lookU.FetchUserByEmail(sqlCon); lookU.DbId > 0 {
			k := fmt.Sprint("Json ", (i+1), " already exists")
			uError[k] = u
			continue
		}

		err = u.CreateUser(sqlCon)
		if err != nil {
			k := fmt.Sprint("Json ", (i+1), " had problems during insertion")
			uError[k] = u
			continue
		}
	}

	defer sqlCon.Close()

	if len(uError) > 0 {
		respondWithJSON(w, http.StatusConflict, uError)
	} else {
		respondWithJSON(w, http.StatusCreated, uSlice)
	}
}
