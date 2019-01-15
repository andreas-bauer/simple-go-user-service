/**
 * Copyright (c) 2019 Andreas Bauer
 *
 * SPDX-License-Identifier: MIT
 */
package model

type User struct {
	Name     string `json: "name"`
	Email    string `json: "email"`
	Password string `json: "password"`
	Role     string `json: "role"`
}
