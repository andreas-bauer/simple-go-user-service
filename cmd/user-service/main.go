package main

import (
	"github.com/andreas-bauer/simple-go-user-service/pkg/http/rest"
	"github.com/andreas-bauer/simple-go-user-service/pkg/model"
	"log"
	"net/http"
)

var users []model.User

const (
	port = ":8080"
)

func main() {
	router := rest.Router()

	log.Println("Starting server on port " + port)
	log.Fatal(http.ListenAndServe(port, router))
}
