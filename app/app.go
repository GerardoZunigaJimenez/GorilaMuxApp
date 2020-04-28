package app

import (
	"awesomeProject/config"
	"awesomeProject/model"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type App struct {
	Router *mux.Router
	DBInfo *config.DataBaseConnInfo
}

func (a *App) Initialize(info config.DataBaseConnInfo) {
	log.Println("Starting Awesome Project")
	a.setDBConnection(info)
	a.setApiRoutes()
}

// Run function exposes our app to the web at a given address
func (a *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, a.Router))
}

func (a *App) setDBConnection(info config.DataBaseConnInfo) {
	e := info.WaitUntilDataBaseIsUp()
	if e != nil {
		log.Fatal("The Database is not available to establish a connection with it")
	}

	//Assign DataBase Info to the App
	a.DBInfo = &info
	log.Println("DB connection set successfully!!")
}

func (a *App) setApiRoutes() {
	a.Router = mux.NewRouter().StrictSlash(true)

	api := a.Router.PathPrefix("/awesomeProject").Subrouter()
	api.HandleFunc(model.UserHandlerFuncUrl, a.createUser).Methods("POST")
	api.HandleFunc(model.UserHandlerFuncUrl, a.getAllUsers).Methods("Get")
	api.HandleFunc(model.UserHandlerCreateUserBulkUrl, a.createUserBulk).Methods("Post")
	api.HandleFunc(model.UserHandlerFetchByEmailUrl, a.getUserByEmail).Methods("GET")
	api.HandleFunc(model.UserHandlerFetchByIdUrl, a.getUserById).Methods("GET")

	log.Println("API routes set successfully!")
}

func logError(message error) {
	f, _ := os.OpenFile("/var/tmp/awesomeProject_logfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	log.SetOutput(f)
	log.Println(message.Error())
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
