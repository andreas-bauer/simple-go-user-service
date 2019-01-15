package user

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

//Service service interface
type Service struct {
	repo Repository
}

//NewService create new service
func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

//Find all user.
func (s *Service) FindAll() ([]*User, error) {
	return s.repo.FindAll()
}

//Find user by its email address.
func (s *Service) FindByEmail(email string) (*User, error) {
	return s.repo.FindByEmail(email)
}

//Delete user with given email address.
func (s *Service) Delete(email string) error {
	return s.repo.Delete(email)
}

//Store user and hash its password.
func (s *Service) Save(user *User) error {
	user.Role = strings.ToUpper(user.Role)
	user.Password = hashAndSalt(user.Password)

	if !isValidRole(user.Role) {
		msg := fmt.Sprintf("User role '%v' is not a valid role. Available roles: %v", user.Role, Roles)
		logrus.Error(msg)
		return errors.New(msg)
	}

	return s.repo.Save(user)
}

func hashAndSalt(rawPW string) string {
	pw := []byte(rawPW)
	hash, err := bcrypt.GenerateFromPassword(pw, bcrypt.MinCost)

	if err != nil {
		logrus.WithError(err)
	}

	return string(hash)
}

func isValidRole(role string) bool {
	for _, item := range Roles {
		if item == role {
			return true
		}
	}

	return false
}
