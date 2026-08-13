package services

import (
	"example/go_backoffice/dto/user"
	"example/go_backoffice/enums"
	"example/go_backoffice/mappers"
	"example/go_backoffice/models"
	"example/go_backoffice/policies"
	"example/go_backoffice/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAllByRole(role string) ([]user.UserResponse, error)
	GetAllByRoleScoped(role string, actor policies.AuthContext) ([]user.UserResponse, error)
	GetUserByIDAndRole(id uint, role string) (*user.UserResponse, error)
	CreateUser(request *user.UserRequest, actor policies.AuthContext) (*user.UserResponse, error)
	ChangeStatus(userID uint, status enums.Status) (*user.UserResponse, error)
	DeleteUser(id string) error
}

type userService struct {
	repo         repositories.UserRepo
	scopeRepo    repositories.ScopeRepo
	agencyPolicy policies.AgencyPolicy
	scopePolicy  policies.ScopePolicy
}

func NewUserService(
	repo repositories.UserRepo,
	scopeRepo repositories.ScopeRepo,
	agencyPolicy policies.AgencyPolicy,
	scopePolicy policies.ScopePolicy,
) UserService {
	return &userService{repo: repo, scopeRepo: scopeRepo, agencyPolicy: agencyPolicy, scopePolicy: scopePolicy}
}

func (s *userService) GetAllByRole(role string) ([]user.UserResponse, error) {
	list, err := s.repo.GetAllByRole(role)
	if err != nil {
		return nil, err
	}
	return mappers.ToUserResponses(list), nil
}

func (s *userService) GetAllByRoleScoped(role string, actor policies.AuthContext) ([]user.UserResponse, error) {
	var scope policies.Scope
	var err error

	if enums.Role(role) == enums.RoleAgent {
		scope, err = s.scopePolicy.AgentScope(actor)
	} else {
		scope, err = s.scopePolicy.AgencyScope(actor)
	}
	if err != nil {
		return nil, err
	}

	if scope.Unrestricted {
		return s.GetAllByRole(role)
	}

	var list []models.User
	if scope.FilterByForeignID {
		list, err = s.repo.GetAllByRoleAndForeignIDs(role, scope.IDs)
	} else {
		list, err = s.repo.GetAllByRoleAndIDs(role, scope.IDs)
	}
	if err != nil {
		return nil, err
	}
	return mappers.ToUserResponses(list), nil
}

func (s *userService) GetUserByIDAndRole(id uint, role string) (*user.UserResponse, error) {

	u, err := s.repo.GetByIDAndRole(id, role)
	if err != nil {
		return nil, err
	}
	response := mappers.ToUserResponse(u)
	return &response, nil
}

func (s *userService) CreateUser(request *user.UserRequest, actor policies.AuthContext) (*user.UserResponse, error) {
	if enums.Role(request.Role) == enums.RoleAgency {
		if err := s.agencyPolicy.CanCreateAgency(actor, request.ForeignId); err != nil {
			return nil, err
		}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	request.Password = string(hashed)
	newUser := mappers.ToUserModel(request)

	if err := s.repo.Create(newUser); err != nil {
		return nil, err
	}

	if enums.Role(request.Role) == enums.RoleAgency && actor.Role == enums.RoleOperator.String() {
		if err := s.scopeRepo.AssignAgencyToOperator(actor.UserID, newUser.ID); err != nil {
			return nil, err
		}
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
