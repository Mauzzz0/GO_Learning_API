package service

import (
	"crypto/sha1"
	"fmt"
	"learning_api/internal/entity"
	"learning_api/internal/repository"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user entity.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	// return fmt.Sprintf("%s", hash.Sum([]byte(os.Getenv("PWD_HASH_SALT"))))
	// Здесь проблемы с кодировкой!!

	return fmt.Sprintf("%s%s", password, "salt")
}
