# Simple Go (User) Micro-Service

The purpose of this project is to gain experience in Golang for micro-service development.
Also, a small comparison with Spring Boot should show how the resource footprint of both technologies differs.

## Build and Run

Run the application with go run:
```
cd cmd/user-service
go run main.go
```

Build a binary:
```
cd cmd/user-service
go build
```

To use the application within a Docker container you can use `docker-compose`, which will also start a container for MongoDB.
```
docker-compose up
```

## API Endpoints

| Operation | URI  | Description  | 
|---|---|---|
| GET | /users/  |  Get all users |
| GET | /users/{email}  | Get user by email |
| POST | /users/  | Create new user |
| DELETE | /users/{email} | Delete user by email |

Request body for user creation:
```JSON
{
  "name": "Rick Sanchez",
  "email": "rick@internet.com",
  "password": "secret pw",
  "role": "admin"
}
```

## Footprint

This section should give you an idea how the footprint of a small micro-service written in Go can be.
Docker is an essential technology in modern software development, so this micro-service should also run within a Docker container.   
To keep the footprint small, I used the [alpine](https://hub.docker.com/_/alpine) docker image as the base image for this micro-service.

The following numbers are gathered via the `docker stats` command.

* RAM usage in idle after startup: ~ 3 MB
* RAM usage after 1000 POST and 1000 GET requests: ~ 8 MB
* Binary size: 14 MB
* Docker image size (binary included): 19 MB
 
 
 **Comparison with Spring Boot 2.0**

A very popular framework for micro-service development is Spring Boot 2.0 for Java. 
The Spring Boot project for this small comparison has similar features and uses the same MongoDB with the same user model structure to store the user objects. 
For the Docker base image I used [adoptopenjdk/openjdk8-openj9:alpine](https://hub.docker.com/r/adoptopenjdk/openjdk8-openj9/), which contains the cloud-optimized  OpenJ9 Java virtual machine implementation, instead of Open JDK. 
Open J9 helps to reduce the memory footprint of Java applications.

* RAM usage in idle after startup: ~ 110 MB
* RAM usage after 1000 POST and 1000 GET requests: ~ 140 MB
* Jar size: 33 MB
* Docker image size (jar included): 250 MB


## License
MIT, see [LICENSE](https://github.com/andreas-bauer/simple-go-user-service/blob/master/LICENSE).
