package store

import (
	"errors"
	"fmt"
	"log/slog"
)

type MemUserStore struct {
	users []User
}

func NewMemUserStore() *MemUserStore {
	return &MemUserStore{users: []User{
		{Id: 1, Name: "Admin", Email: "admin@example.com", Password: "$2a$10$UF/l5lpDnTGInZ/qWYg1ZepN44nJP1EmznhP6fZmjymftr2V4WPB2", IsAdmin: true},
		{Id: 2, Name: "User", Email: "user@example.com", Password: "$2a$10$i3OpRdnz90uJFNW9O7CVveUJfwYxO7vws.VjaWR5AFH.vhn0sh16K", IsAdmin: false},
	}}
}

func (u *MemUserStore) AddUser(username, email, password string) (*User, error) {
	if _, err := u.GetUserByEmail(email); err == nil {
		return nil, errors.New("User already exists")
	}

	userId := len(u.users)
	user := &User{Id: userId, Name: username, Email: email, Password: password, IsAdmin: false}
	u.users = append(u.users, *user)
	return user, nil
}

func (u *MemUserStore) GetUserByEmail(email string) (*User, error) {
	slog.Info("Get Users", "email", email, "users", u.users)
	for _, user := range u.users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("invalid user email")
}
