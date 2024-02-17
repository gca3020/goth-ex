package services

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/gca3020/goth-ex/internal/store"
)

type AuthService struct {
	userStore store.UserStore
}

func NewAuthService(userStore store.UserStore) *AuthService {
	return &AuthService{
		userStore: userStore,
	}
}

func (s *AuthService) AddUser(name, email, password string) (*store.User, error) {
	encPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	slog.Info("Adding user", "name", name, "email", email, "password", password, "enc", encPassword)
	return s.userStore.AddUser(name, email, string(encPassword))
}

func (s *AuthService) LoginUser(email, password string) (*store.User, error) {
	user, err := s.userStore.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}
	return user, nil
}
