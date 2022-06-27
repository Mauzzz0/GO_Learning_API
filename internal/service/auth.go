package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"learning_api/internal/entity"
	"learning_api/internal/repository"
	"os"
	"strconv"
	"time"
)

type AuthService struct {
	repo repository.Authorization
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user entity.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(_user entity.UserLogin) (string, error) {
	user, err := s.repo.GetUser(_user.Username, generatePasswordHash(_user.Password))
	if err != nil {
		return "", err
	}

	ttl, err := strconv.ParseInt(os.Getenv("JWT_TOKEN_TTL_HOUR"), 10, 64)

	if err != nil {
		logrus.Fatal("JWT_TOKEN_TTL_HOUR env must be integer: %s", err.Error())
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(ttl) * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		}, user.Id})

	return token.SignedString([]byte(os.Getenv("JWT_TOKEN_KEY")))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing methods")
		}

		return []byte(os.Getenv("JWT_TOKEN_KEY")), nil
	})
	if err != nil {
		return 0, nil
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	// return fmt.Sprintf("%s", hash.Sum([]byte(os.Getenv("PWD_HASH_SALT"))))
	// Здесь проблемы с кодировкой!!

	return fmt.Sprintf("%s%s", password, "salt")
}
