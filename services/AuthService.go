package services

import (
	"errors"
	"time"

	"example/go_backoffice/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("credenziali non valide")

type AuthService interface {
	Login(email string, password string) (string, error)
}

type authService struct {
	repo      repositories.UserRepo
	jwtSecret []byte
}

func NewAuthService(repo repositories.UserRepo, jwtSecret string) AuthService {
	return &authService{repo: repo, jwtSecret: []byte(jwtSecret)}
}

func (s *authService) Login(email string, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
