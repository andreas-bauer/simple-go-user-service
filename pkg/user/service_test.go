/**
 * Copyright (c) 2019 Andreas Bauer
 *
 * SPDX-License-Identifier: MIT
 */
package user

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type FakeDB struct{}

func TestFindAll(t *testing.T) {
	initTestData()
	s := NewService(&FakeDB{})

	users, err := s.FindAll()

	if len(users) != 3 {
		t.Fatal("Expected size of all users to be 3 but is ", len(users))
	}
	if err != nil {
		t.Fatal("Expected error to be nil")
	}
}

func TestFindByEmail(t *testing.T) {
	initTestData()
	s := NewService(&FakeDB{})
	expectedEmail := "tu2@test.com"

	user, err := s.FindByEmail(expectedEmail)

	if err != nil {
		t.Fatal(err)
	}
	if user.Email != expectedEmail {
		t.Fatalf("Expected user's email to be %v but is %v", expectedEmail, user.Email)
	}
}

func TestFindByEmail_NotFound(t *testing.T) {
	initTestData()
	s := NewService(&FakeDB{})
	expectedEmail := "notexisting@test.com"

	_, err := s.FindByEmail(expectedEmail)

	if err == nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	initTestData()
	s := NewService(&FakeDB{})

	err := s.Delete("tu2@test.com")

	if err != nil {
		t.Fatal()
	}
	if len(users) != 2 {
		t.Fatal("Expected size of all users to be 2 but is ", len(users))
	}
}

func TestSave(t *testing.T) {
	initTestData()
	s := NewService(&FakeDB{})
	user := &User{Name: "New Test User", Email: "t@t.com", Password: "secret", Role: "admin"}

	s.Save(user)

	if len(users) != 4 {
		t.Fatal("Expected size of all users to be 4 but is ", len(users))
	}

	resultUser, err := userByEmail("t@t.com")

	if err != nil {
		t.Fatal()
	}

	if resultUser.Role != "ADMIN" {
		t.Fatalf("User's role should be ADMIN but is %v", resultUser.Role)
	}

	if resultUser.Password == "secret" {
		t.Fatalf("User's passord should be encrypted")
	}

	err = bcrypt.CompareHashAndPassword([]byte(resultUser.Password), []byte("secret"))
	if err != nil {
		t.Fatal("User's password is not hashed correct")
	}
}

func TestSave_InvalidRole(t *testing.T) {
	initTestData()
	s := NewService(&FakeDB{})
	user := &User{Name: "New Test User", Email: "t@t.com", Password: "secret", Role: "notExistingRole"}

	err := s.Save(user)

	if err == nil {
		t.Fatal("Should return error because of wrong role")
	}
}

/*
 *	Fake DB and Test Data
 */

var users []*User

func initTestData() {
	users = []*User{}
	users = append(users, &User{Name: "TestUser1", Email: "tu1@test.com", Password: "secret", Role: "ADMIN"})
	users = append(users, &User{Name: "TestUser2", Email: "tu2@test.com", Password: "secret", Role: "ADMIN"})
	users = append(users, &User{Name: "TestUser3", Email: "tu3@test.com", Password: "secret", Role: "ADMIN"})
}

func (db *FakeDB) FindAll() (result []*User, err error) {
	return users, nil
}

func (db *FakeDB) FindByEmail(email string) (*User, error) {
	user, err := userByEmail(email)
	return user, err
}

func (db *FakeDB) Delete(email string) (err error) {
	for index, item := range users {
		if item.Email == email {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	return nil
}

func (db *FakeDB) Save(user *User) (err error) {
	users = append(users, user)
	return nil
}

func userByEmail(email string) (*User, error) {
	for _, item := range users {
		if item.Email == email {
			return item, nil
		}
	}

	return nil, errors.New("No user found")
}
