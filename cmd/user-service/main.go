package main

import (
	"github.com/andreas-bauer/simple-go-user-service/pkg/http/rest"
)

func main() {
	server := &rest.Instance{}
	server.Start()
}
