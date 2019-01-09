package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/andreas-bauer/simple-go-user-service/pkg/model"
	"github.com/gorilla/mux"
)

var users []model.User

func init() {
	users = append(users, model.User{Name: "Andi", Email: "andi@andi.de", Password: "1234", Role: "ADMIN"})
	users = append(users, model.User{Name: "Sepp", Email: "sepp@internet.de", Password: "545646", Role: "ADMIN"})
	users = append(users, model.User{Name: "Hans", Email: "hans@web.de", Password: "231234", Role: "ADMIN"})
}

func GetUsers(writer http.ResponseWriter, req *http.Request) {
	log.Println("GetUsers")

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(users)
}

func GetUser(writer http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	log.Println("GetUser")

	user, err := doGetUser(params["Email"])
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(user)
}

func CreateUser(writer http.ResponseWriter, req *http.Request) {
	var user model.User
	_ = json.NewDecoder(req.Body).Decode(&user)
	user.Role = strings.ToUpper(user.Role)

	_, err := doGetUser(user.Email)
	if err == nil {
		msg := fmt.Sprintf("User with email %v already exist!", user.Email)
		log.Println(msg)
		http.Error(writer, msg, http.StatusConflict)
		return
	}

	if !isValidRole(user.Role) {
		msg := fmt.Sprintf("User role '%v' is not a valid role. Available roles: %v", user.Role, model.Roles)
		log.Println(msg)
		http.Error(writer, msg, http.StatusInternalServerError)
		return
	}

	users = append(users, user)

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(user)
}

func DeleteUser(writer http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range users {
		if item.Email == params["Email"] {
			users = append(users[:index], users[index+1:]...)
			break
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(users)
	}
}

func doGetUser(email string) (model.User, error) {
	for _, item := range users {
		if item.Email == email {
			return item, nil
		}
	}

	return model.User{}, errors.New("User with email " + email + " does not exist")
}

func isValidRole(role string) bool {
	for _, item := range model.Roles {
		if item == role {
			return true
		}
	}

	return false
}
