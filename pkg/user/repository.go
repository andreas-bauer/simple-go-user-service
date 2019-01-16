/**
 * Copyright (c) 2019 Andreas Bauer
 *
 * SPDX-License-Identifier: MIT
 */
package user

type Repository interface {
	FindAll() (result []*User, err error)
	FindByEmail(email string) (*User, error)
	Delete(email string) (err error)
	Save(*User) (err error)
}
