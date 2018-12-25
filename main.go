package main

import (
	"log"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

type User struct {
	Name     string `json: "Name"`
	Email    string `json: "Email"`
	Password string `json: "Password"`
	Role     string `json: "Role"`
}

var users []User

const (
	port = ":8080"
)

func main() {
	users = append(users, User{Name: "Andi", Email:"andi@andi.de", Password:"1234", Role:"ADMIN"})
	users = append(users, User{Name: "Sepp", Email:"sepp@internet.de", Password:"545646", Role:"ADMIN"})
	users = append(users, User{Name: "Hans", Email:"hans@web.de", Password:"231234", Role:"ADMIN"})

	router := GetRouter();

	log.Println("Starting server on port " + port)
	log.Fatal(http.ListenAndServe(port, router))
}

func GetUsers(writer http.ResponseWriter, req *http.Request) {
	log.Println("GetUsers")
	json.NewEncoder(writer).Encode(users)
}

func GetUser(writer http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	log.Println("GetUser")
	for _, item := range users {
		if item.Email == params["Email"] {
			json.NewEncoder(writer).Encode(item)
		}
	}
}

func CreateUser(writer http.ResponseWriter, req *http.Request) {
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)
	users = append(users, user)
	json.NewEncoder(writer).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range users {
		if item.Email == params["Email"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(users)
	}
}