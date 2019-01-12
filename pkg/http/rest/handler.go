package rest

import (
	"encoding/json"
	"fmt"
	"github.com/andreas-bauer/simple-go-user-service/pkg/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (srv *Instance) GetUsers(writer http.ResponseWriter, req *http.Request) {
	logrus.Info("Http Request: GET /users/")

	users, err := srv.db.FindAll()
	if err != nil {
		logrus.WithError(err)
		writer.Header().Set("Content-Type", "text/plain")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(users)
}

func (srv *Instance) GetUser(writer http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	logrus.Info("Http Request: GET /users/", params["Email"])

	user, err := srv.db.FindByEmail(params["Email"])
	if err != nil {
		logrus.WithError(err)
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(user)
}

func (srv *Instance) CreateUser(writer http.ResponseWriter, req *http.Request) {
	logrus.Info("Http Request: POST /users/")
	var user model.User
	_ = json.NewDecoder(req.Body).Decode(&user)
	user.Role = strings.ToUpper(user.Role)

	existAlready := srv.db.ContainsUserWithEmail(user.Email)
	if existAlready {
		msg := fmt.Sprintf("User with email %v already exist!", user.Email)
		logrus.Error(msg)
		http.Error(writer, msg, http.StatusConflict)
		return
	}

	if !isValidRole(user.Role) {
		msg := fmt.Sprintf("User role '%v' is not a valid role. Available roles: %v", user.Role, model.Roles)
		logrus.Error(msg)
		http.Error(writer, msg, http.StatusInternalServerError)
		return
	}

	srv.db.Save(user)

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(user)
}

func (srv *Instance) DeleteUser(writer http.ResponseWriter, req *http.Request) {
	logrus.Info("Http Request: DELETE /users/")
	params := mux.Vars(req)
	err := srv.db.Delete(params["Email"])

	if err != nil {
		writer.Header().Set("Content-Type", "text/plain")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func isValidRole(role string) bool {
	for _, item := range model.Roles {
		if item == role {
			return true
		}
	}

	return false
}
