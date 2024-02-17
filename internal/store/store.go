package store

type User struct {
	Id       int
	Email    string
	Name     string
	Password string
	IsAdmin  bool
}

type UserStore interface {
	AddUser(name, email, password string) (*User, error)
	GetUserByEmail(email string) (*User, error)
}
