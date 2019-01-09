package model

type Role struct {
	ADMIN  string
	PUBLIC string
}

var Enum = &Role{
	ADMIN:  "ADMIN",
	PUBLIC: "PUBLIC",
}

var Roles = [...]string{Enum.ADMIN, Enum.PUBLIC}
