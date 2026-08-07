package services

import (
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"

	"github.com/google/uuid"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id string) (*models.User, error)
	CreateUser(user *models.User) error
	DeleteUser(id string) error
}

type userService struct {
	repo repositories.UserRepo
}

func NewUserService(repo repositories.UserRepo) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) CreateUser(user *models.User) error {
	// Genera automaticamente un UUID per la stringa ID
	user.ID = uuid.New().String()

	// Qui è possibile inserire l'hashing della password prima del salvataggio
	return s.repo.Create(user)
}

func (s *userService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}
