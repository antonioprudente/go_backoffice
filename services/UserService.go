// services/UserService.go
package services

import (
	"example/go_backoffice/dto/user"
	"example/go_backoffice/enums"
	"example/go_backoffice/mappers"
	"example/go_backoffice/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAllByRole(role string) ([]user.UserResponse, error)
	GetUserByIDAndRole(id uint, role string) (*user.UserResponse, error)
	CreateUser(request *user.UserRequest) (*user.UserResponse, error)
	ChangeStatus(userID uint, status enums.Status) (*user.UserResponse, error)
	DeleteUser(id string) error
}

type userService struct {
	repo repositories.UserRepo
}

func NewUserService(repo repositories.UserRepo) UserService {
	return &userService{repo: repo}
}

// LISTA UTENTI
func (s *userService) GetAllByRole(role string) ([]user.UserResponse, error) {
	list, err := s.repo.GetAllByRole(role)
	if err != nil {
		return nil, err
	}
	responses := mappers.ToUserResponses(list)
	return responses, nil
}

// TROVA UTANTE DA ID
func (s *userService) GetUserByIDAndRole(id uint, role string) (*user.UserResponse, error) {
	user, err := s.repo.GetByIDAndRole(id, role)
	if err != nil {
		return nil, err
	}

	response := mappers.ToUserResponse(user)
	return &response, nil
}

// CREAZIONE USER
func (s *userService) CreateUser(request *user.UserRequest) (*user.UserResponse, error) {
	// crittografia password
	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	request.Password = string(hashed)       // aggiornamento password nella request con l'hash
	newUser := mappers.ToUserModel(request) // conversione da request a model

	// salvataggio del model nel db tramite repository
	if err := s.repo.Create(newUser); err != nil {
		return nil, err
	}
	response := mappers.ToUserResponse(newUser) // conversione da model a response
	return &response, nil
}

// AGGIORNAMENTO STATO
func (s *userService) ChangeStatus(userID uint, status enums.Status) (*user.UserResponse, error) {
	updated, err := s.repo.UpdateStatus(userID, string(status))
	if err != nil {
		return nil, err
	}
	response := mappers.ToUserResponse(updated)
	return &response, nil
}

// CANCELLAZIONE UTENTE
func (s *userService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}
