package rest

import "github.com/gorilla/mux"

func Router(srv *Instance) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users/", srv.GetUsers).Methods("GET")
	router.HandleFunc("/users/{Email}", srv.GetUser).Methods("GET")
	router.HandleFunc("/users/", srv.CreateUser).Methods("POST")
	router.HandleFunc("/users/{Email}", srv.DeleteUser).Methods("DELETE")
	return router
}
