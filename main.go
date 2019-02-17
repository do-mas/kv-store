package main

import (
	"kv-store/main/api"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	// router.HandleFunc("/storage/{id}", api.GetStorage).Methods("GET")
	router.HandleFunc("/storage/{id}", api.GetFile).Methods("GET")
	router.HandleFunc("/storage", api.CreateStorage).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))

}
