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
	UpdateUser(id uint, request *user.UserRequest, actor policies.AuthContext) (*user.UserResponse, error)
	ChangeStatus(userID uint, targetRole string, status enums.Status, actor policies.AuthContext) (*user.UserResponse, error)
	DeleteUserByIdAndRole(id uint, targetRole string, actor policies.AuthContext) error
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

func (s *userService) UpdateUser(id uint, request *user.UserRequest, actor policies.AuthContext) (*user.UserResponse, error) {
	existing, err := s.repo.GetByIDAndRole(id, request.Role)
	if err != nil {
		return nil, err
	}

	if err := s.policy.Update(actor, existing); err != nil {
		return nil, err
	}

	if request.FirstName != "" || existing.FirstName != request.FirstName {
		existing.FirstName = request.FirstName
	}
	if request.LastName != "" || existing.LastName != request.LastName {
		existing.LastName = request.LastName
	}
	if request.Username != "" || existing.Username != request.Username {
		existing.Username = request.Username
	}
	if request.Email != "" || existing.Email != request.Email {
		existing.Email = request.Email
	}
	if request.ForeignId != nil || existing.ForeignId != request.ForeignId {
		existing.ForeignId = request.ForeignId
	}

	if request.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		existing.Password = string(hashed)
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	response := mappers.ToUserResponse(existing)
	return &response, nil
}

func (s *userService) ChangeStatus(userID uint, targetRole string, status enums.Status, actor policies.AuthContext) (*user.UserResponse, error) {

	existing, err := s.repo.GetByIDAndRole(userID, targetRole)
	if err != nil {
		return nil, err
	}

	if err := s.policy.UpdateStatus(actor, existing); err != nil {
		return nil, err
	}

	updated, err := s.repo.UpdateStatusByIdAndRole(userID, targetRole, string(status))
	if err != nil {
		return nil, err
	}

	response := mappers.ToUserResponse(updated)
	return &response, nil
}

func (s *userService) DeleteUserByIdAndRole(id uint, targetRole string, actor policies.AuthContext) error {
	existing, err := s.repo.GetByIDAndRole(id, targetRole)
	if err != nil {
		return err
	}

	if err := s.policy.Delete(actor, existing); err != nil {
		return err
	}
	return s.repo.DeleteByIdAndRole(id, targetRole)
}
