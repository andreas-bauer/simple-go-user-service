package rest

import "github.com/gorilla/mux"

func Router(srv *Instance) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users/", GetUsers).Methods("GET")
	router.HandleFunc("/users/{Email}", srv.GetUser).Methods("GET")
	router.HandleFunc("/users/", CreateUser).Methods("POST")
	router.HandleFunc("/users/{Email}", DeleteUser).Methods("DELETE")
	return router
}
