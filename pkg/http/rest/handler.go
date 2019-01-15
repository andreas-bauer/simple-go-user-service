/**
 * Copyright (c) 2019 Andreas Bauer
 *
 * SPDX-License-Identifier: MIT
 */
package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andreas-bauer/simple-go-user-service/pkg/user"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func Router(service *user.Service) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users/", handleGetAllUsers(service)).Methods("GET")
	router.HandleFunc("/users/{Email}", handleGetUser(service)).Methods("GET")
	router.HandleFunc("/users/", handleCreateUser(service)).Methods("POST")
	router.HandleFunc("/users/{Email}", handleDeleteUser(service)).Methods("DELETE")
	return router
}

func handleGetAllUsers(service *user.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Http Request: GET /users/")

		users, err := service.FindAll()
		if err != nil {
			logrus.WithError(err)
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})
}

func handleGetUser(service *user.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		logrus.Info("Http Request: GET /users/", params["Email"])

		user, err := service.FindByEmail(params["Email"])
		if err != nil {
			logrus.WithError(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})
}

func handleCreateUser(service *user.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Http Request: POST /users/")
		var user user.User
		_ = json.NewDecoder(r.Body).Decode(&user)

		err := service.Save(&user)
		if err != nil {
			msg := fmt.Sprintf("Unable to store user %v", user)
			logrus.WithError(err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})
}

func handleDeleteUser(service *user.Service) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Http Request: DELETE /users/")
		params := mux.Vars(r)
		err := service.Delete(params["Email"])

		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
