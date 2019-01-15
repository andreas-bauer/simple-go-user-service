package user

type Repository interface {
	FindAll() (result []*User, err error)
	FindByEmail(email string) (*User, error)
	Delete(email string) (err error)
	Save(*User) (err error)
}
