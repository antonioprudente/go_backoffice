package services

import (
	"errors"

	"example/go_backoffice/config"
	"example/go_backoffice/dto/auth"
	"example/go_backoffice/mappers"
	"example/go_backoffice/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(req *auth.LoginRequest, plainPassword string) (*auth.LoginResponse, error)
}

type authService struct {
	repo repositories.UserRepo
}

func NewAuthService(repo repositories.UserRepo) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Login(req *auth.LoginRequest, plainPassword string) (*auth.LoginResponse, error) {
	u, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("credenziali non valide")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword)); err != nil {
		return nil, errors.New("credenziali non valide")
	}

	token, err := config.GenerateToken(u.ID, string(u.Role))
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Token: token,
		User:  mappers.ToUserResponse(u),
	}, nil
}
