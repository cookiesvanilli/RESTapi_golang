package service

import (
	"crypto/sha1"
	"fmt"
	todo "github.com/cookiesvanilli/go_app"
	"github.com/cookiesvanilli/go_app/pkg/repository"
)

const salt = "1q2w3e4r5t"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

//хеширование пароля
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
