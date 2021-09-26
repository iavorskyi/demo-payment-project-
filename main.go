package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var (
	dbConnectionString = "185.226.42.23:5432"
	dbUser             = "admin"
	dbPassword         = "HRJyfkG1XoePiMumcD7a406Ct2Q9538g"
	dbName             = "fleetsy"
	ren                = render.New()
)

func main() {
	log.Println("Listening port: 8003!!!")
	log.Fatal(http.ListenAndServe(":8003", Handlers()))
}

func Handlers() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/transactions/create", createTransaction).Methods("POST")
	r.HandleFunc("/users/{user_id}/balance", getBalance).Methods("GET")

	return r
}
