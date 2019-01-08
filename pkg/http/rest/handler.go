package rest

import (
	"encoding/json"
	"github.com/andreas-bauer/simple-go-user-service/pkg/model"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	for _, item := range users {
		if item.Email == params["Email"] {
			writer.Header().Set("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(item)
		}
	}
}

func CreateUser(writer http.ResponseWriter, req *http.Request) {
	var user model.User
	_ = json.NewDecoder(req.Body).Decode(&user)
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
