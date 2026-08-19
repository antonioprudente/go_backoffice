package services

import (
	"errors"
	"example/go_backoffice/dto/user"
	"example/go_backoffice/enums"
	"example/go_backoffice/mappers"
	"example/go_backoffice/models"
	"example/go_backoffice/policies"
	"example/go_backoffice/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	GetAllByRole(role string) ([]user.UserResponse, error)
	GetUserByIDAndRole(id uint, targetRole string, actor policies.AuthContext) (*user.UserResponse, error)
	CreateUser(request *user.UserRequest, actor policies.AuthContext) (*user.UserResponse, error)
	ChangeStatus(userID uint, status enums.Status) (*user.UserResponse, error)
	DeleteUser(id string) error
}

type userService struct {
	db           *gorm.DB
	repo         repositories.UserRepo
	scopeRepo    repositories.ScopeRepo
	agencyOpRepo repositories.AgencyOperatorRepo
	policy       policies.UserPolicy
}

func NewUserService(
	db *gorm.DB,
	repo repositories.UserRepo,
	scopeRepo repositories.ScopeRepo,
	agencyOpRepo repositories.AgencyOperatorRepo,
	policy policies.UserPolicy,
) UserService {
	return &userService{
		db:           db,
		repo:         repo,
		scopeRepo:    scopeRepo,
		agencyOpRepo: agencyOpRepo,
		policy:       policy,
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

	if err := s.policy.View(actor, target); err != nil {
		return nil, err
	}

	response := mappers.ToUserResponse(target)
	return &response, nil
}

func (s *userService) CreateUser(request *user.UserRequest, actor policies.AuthContext) (*user.UserResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	request.Password = string(hashed)
	newUser := mappers.ToUserModel(request)

	if err := s.policy.Create(actor, newUser); err != nil {
		return nil, err
	}
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		userRepo := s.repo.WithTx(tx)
		agencyOpRepo := s.agencyOpRepo.WithTx(tx)

		if err := userRepo.Create(newUser); err != nil {
			return err
		}
		// Se un OPERATOR crea una AGENCY, va creato anche il collegamento
		// nella tabella pivot AgencyOperator
		if actor.Role == enums.RoleOperator.String() && newUser.Role == enums.RoleAgency {
			agencyOp := &models.AgencyOperator{
				OperatorID: actor.UserID,
				AgencyID:   newUser.ID,
			}
			if _, err := agencyOpRepo.AssignAgency(agencyOp); err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return nil, txErr
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
