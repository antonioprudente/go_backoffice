package services

import (
	"errors"
	"example/go_backoffice/dto/user"
	"example/go_backoffice/enums"
	"example/go_backoffice/mappers"
	"example/go_backoffice/models"
	"example/go_backoffice/policies"
	"example/go_backoffice/repositories"
	"fmt"
	"reflect"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	GetAllByRole(role string) ([]user.UserResponse, error)
	GetUserByIDAndRole(id uint, targetRole string, actor policies.AuthContext) (*user.UserResponse, error)
	CreateUser(request *user.UserRequest, actor policies.AuthContext) (*user.UserResponse, error)
	UpdateUser(id uint, request *user.UserRequest, actor policies.AuthContext) (*user.UserResponse, error)
	ChangeStatus(userID uint, targetRole string, status enums.Status, actor policies.AuthContext) (*user.UserResponse, error)
	ChangeForeignID(request user.ChangeForeignRequest, targetRole string, actor policies.AuthContext) (*user.UserResponse, error)
	DeleteUserByIdAndRole(id uint, targetRole string, actor policies.AuthContext) error
}

type userService struct {
	db           *gorm.DB
	repo         repositories.UserRepo
	scopeRepo    repositories.ScopeRepo
	agencyOpRepo repositories.AgencyOperatorRepo
	logService   ActivityLogService
	policy       policies.UserPolicy
}

func NewUserService(
	db *gorm.DB,
	repo repositories.UserRepo,
	scopeRepo repositories.ScopeRepo,
	agencyOpRepo repositories.AgencyOperatorRepo,
	logService ActivityLogService,
	policy policies.UserPolicy,
) UserService {
	return &userService{
		db:           db,
		repo:         repo,
		scopeRepo:    scopeRepo,
		agencyOpRepo: agencyOpRepo,
		logService:   logService,
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

	// Activity Log - Creazione Utente
	err = s.logService.NewLog(&models.ActivityLog{
		ActorID: &actor.UserID,
		//ActorRole:   enums.Role(actor.Role),
		Action:      enums.Create,
		TargetType:  reflect.TypeOf(newUser).Elem().Name(),
		TargetID:    &newUser.ID,
		Description: fmt.Sprintf("Creato nuovo utente '%s' con ruolo %s", newUser.Username, newUser.Role),
	})

	if err != nil {
		return nil, err
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
	if request.ForeignID != nil || existing.ForeignID != request.ForeignID {
		existing.ForeignID = request.ForeignID
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

	// Activity Log - Aggiornamento Dati
	err = s.logService.NewLog(&models.ActivityLog{
		ActorID: &actor.UserID,
		//ActorRole:   enums.Role(actor.Role),
		Action:      enums.Update,
		TargetType:  reflect.TypeOf(existing).Elem().Name(),
		TargetID:    &existing.ID,
		Description: fmt.Sprintf("Aggiornati i dati dell'utente '%s'", existing.Username),
	})

	if err != nil {
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

	var action enums.Action
	switch updated.Status {
	case enums.StatusActive:
		action = enums.Active

	case enums.StatusSuspended:
		action = enums.Suspend

	case enums.StatusBlocked:
		action = enums.Block
	}

	// Activity Log - Cambio Stato
	err = s.logService.NewLog(&models.ActivityLog{
		ActorID: &actor.UserID,
		//ActorRole:   enums.Role(actor.Role),
		Action:      action,
		TargetType:  reflect.TypeOf(updated).Elem().Name(),
		TargetID:    &updated.ID,
		Description: fmt.Sprintf("Stato dell'utente '%s' impostato a %s", updated.Username, status),
	})

	if err != nil {
		return nil, err
	}

	response := mappers.ToUserResponse(updated)
	return &response, nil
}

func (s *userService) ChangeForeignID(request user.ChangeForeignRequest, targetRole string, actor policies.AuthContext) (*user.UserResponse, error) {
	target, err := s.repo.GetByIDAndRole(request.UserID, targetRole)
	if err != nil {
		return nil, err
	}

	parentRole := enums.RoleUser.String()
	if targetRole == enums.RoleAgency.String() {
		parentRole = enums.RoleAgent.String()
	}

	if targetRole == enums.RoleUser.String() {
		parentRole = enums.RoleAgency.String()
	}

	newParent, err := s.repo.GetByIDAndRole(*request.TargetID, parentRole)
	if err != nil {
		return nil, err
	}

	if err := s.policy.Move(actor, target, newParent); err != nil {
		return nil, err
	}

	updated, err := s.repo.UpdateForeignID(request.UserID, targetRole, *request.TargetID)
	if err != nil {
		return nil, err
	}

	// Activity Log - Spostamento Relazionale
	err = s.logService.NewLog(&models.ActivityLog{
		ActorID: &actor.UserID,
		//ActorRole:   enums.Role(actor.Role),
		Action:      enums.Move,
		TargetType:  reflect.TypeOf(updated).Elem().Name(),
		TargetID:    &updated.ID,
		Description: fmt.Sprintf("Utente '%s' collegato alla nuova entità genitore #%d", updated.Username, *request.TargetID),
	})

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

	if err := s.repo.DeleteByIdAndRole(id, targetRole); err != nil {
		return err
	}

	// Activity Log - Eliminazione
	err = s.logService.NewLog(&models.ActivityLog{
		ActorID: &actor.UserID,
		//ActorRole:   enums.Role(actor.Role),
		Action:      enums.Delete,
		TargetType:  reflect.TypeOf(models.User{}).Elem().Name(),
		TargetID:    &id,
		Description: fmt.Sprintf("Eliminato l'utente '%s' (Ruolo: %s)", existing.Username, targetRole),
	})

	if err != nil {
		return err
	}

	return nil
}
