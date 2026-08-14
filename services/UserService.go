package services

import (
	"errors"
	"example/go_backoffice/dto/user"
	"example/go_backoffice/enums"
	"example/go_backoffice/mappers"
	"example/go_backoffice/policies"
	"example/go_backoffice/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAllByRole(role string) ([]user.UserResponse, error)
	GetUserByIDAndRole(id uint, targetRole string, actor policies.AuthContext) (*user.UserResponse, error)
	CreateUser(request *user.UserRequest) (*user.UserResponse, error)
	ChangeStatus(userID uint, status enums.Status) (*user.UserResponse, error)
	DeleteUser(id string) error
}

type userService struct {
	repo      repositories.UserRepo
	scopeRepo repositories.ScopeRepo
}

func NewUserService(
	repo repositories.UserRepo,
	scopeRepo repositories.ScopeRepo,
) UserService {
	return &userService{
		repo:      repo,
		scopeRepo: scopeRepo,
	}
}

var ErrUnauthorized = errors.New("non hai i permessi per accedere a questa risorsa")

func (s *userService) GetAllByRole(role string) ([]user.UserResponse, error) {
	list, err := s.repo.GetAllByRole(role)
	if err != nil {
		return nil, err
	}
	return mappers.ToUserResponses(list), nil
}

func (s *userService) GetUserByIDAndRole(id uint, targetRole string, actor policies.AuthContext) (*user.UserResponse, error) {

	target, err := s.repo.GetByIDAndRole(id, targetRole)
	if err != nil {
		return nil, err
	}

	response := mappers.ToUserResponse(target)
	return &response, nil
}

func (s *userService) CreateUser(request *user.UserRequest) (*user.UserResponse, error) {

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	request.Password = string(hashed)
	newUser := mappers.ToUserModel(request)

	if err := s.repo.Create(newUser); err != nil {
		return nil, err
	}

	response := mappers.ToUserResponse(newUser)
	return &response, nil
}

func (s *userService) ChangeStatus(userID uint, status enums.Status) (*user.UserResponse, error) {
	updated, err := s.repo.UpdateStatus(userID, string(status))
	if err != nil {
		return nil, err
	}
	response := mappers.ToUserResponse(updated)
	return &response, nil
}

func (s *userService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}
