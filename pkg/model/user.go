package model

type User struct {
	Name     string `json: "Name"`
	Email    string `json: "Email"`
	Password string `json: "Password"`
	Role     string `json: "Role"`
}
