package rest

import "github.com/gorilla/mux"

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users/", GetUsers).Methods("GET")
	router.HandleFunc("/users/{Email}", GetUser).Methods("GET")
	router.HandleFunc("/users/", CreateUser).Methods("POST")
	router.HandleFunc("/users/{Email}", DeleteUser).Methods("DELETE")
	return router
}
