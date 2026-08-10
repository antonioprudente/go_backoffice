// services/UserService.go
package services

import (
	"example/go_backoffice/dto/user"
	"example/go_backoffice/enums"
	"example/go_backoffice/mappers"
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id string) (*user.UserResponse, error)
	CreateUser(request *user.UserRequest) (*user.UserResponse, error)
	ChangeStatus(userID uint, status enums.Status) (*user.UserResponse, error)
	DeleteUser(id string) error
}

type userService struct {
	repo     repositories.UserRepo
	treeRepo repositories.AgentNodeRepo
}

func NewUserService(repo repositories.UserRepo, treeRepo repositories.AgentNodeRepo) UserService {
	return &userService{repo: repo, treeRepo: treeRepo}
}

// LISTA UTENTI
func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

// TROVA UTANTE DA ID
func (s *userService) GetUserByID(id string) (*user.UserResponse, error) {
	user, err := s.repo.GetByID(id)
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
