package services

import (
	"errors"
	"fmt"
	"reflect"

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
	GetAllByRole(role string, actor policies.AuthContext, filter user.UserFilter) ([]user.UserResponse, error)
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

func (s *userService) GetAllByRole(role string, actor policies.AuthContext, filter user.UserFilter) ([]user.UserResponse, error) {
	list, err := s.scopedList(role, actor, filter)
	if err != nil {
		return nil, err
	}
	return mappers.ToUserResponses(list), nil
}

func (s *userService) scopedList(role string, actor policies.AuthContext, filter user.UserFilter) ([]models.User, error) {
	switch actor.Role {
	case enums.RoleAdmin.String():
		return s.repo.GetAllByRole(role, filter)
	case enums.RoleOperator.String():
		return s.scopedListForOperator(role, actor.UserID, filter)
	case enums.RoleAgent.String():
		return s.scopedListForAgent(role, actor.UserID, filter)
	case enums.RoleAgency.String():
		if role != enums.RoleUser.String() {
			return nil, policies.ErrForbidden
		}
		return s.repo.GetAllByRoleAndForeignIDs(role, []uint{actor.UserID}, filter)
	}
	return nil, policies.ErrUnknownRole
}

func (s *userService) scopedListForOperator(role string, operatorID uint, filter user.UserFilter) ([]models.User, error) {
	switch role {
	case enums.RoleAgent.String():
		ids, err := s.scopeRepo.AssignedAgentIDs(operatorID)
		if err != nil {
			return nil, err
		}
		return s.repo.GetAllByRoleAndIDs(role, ids, filter)
	case enums.RoleAgency.String():
		ids, err := s.scopeRepo.AssignedAgencyIDs(operatorID)
		if err != nil {
			return nil, err
		}
		return s.repo.GetAllByRoleAndIDs(role, ids, filter)
	case enums.RoleUser.String():
		agencyIDs, err := s.scopeRepo.AssignedAgencyIDs(operatorID)
		if err != nil {
			return nil, err
		}
		return s.repo.GetAllByRoleAndForeignIDs(role, agencyIDs, filter)
	}
	return nil, policies.ErrForbidden
}

func (s *userService) scopedListForAgent(role string, agentID uint, filter user.UserFilter) ([]models.User, error) {
	scopeIDs, err := s.scopeRepo.NodeChildrenAndSelfAgentIds(agentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.User{}, nil
		}
		return nil, err
	}

	switch role {
	case enums.RoleAgent.String():
		return s.repo.GetAllByRoleAndIDs(role, scopeIDs, filter)
	case enums.RoleAgency.String():
		return s.repo.GetAllByRoleAndForeignIDs(role, scopeIDs, filter)
	case enums.RoleUser.String():
		agencies, err := s.repo.GetAllByRoleAndForeignIDs(enums.RoleAgency.String(), scopeIDs, user.UserFilter{})
		if err != nil {
			return nil, err
		}
		agencyIDs := make([]uint, len(agencies))
		for i, a := range agencies {
			agencyIDs[i] = a.ID
		}
		return s.repo.GetAllByRoleAndForeignIDs(role, agencyIDs, filter)
	}
	return nil, policies.ErrForbidden
}

func (s *userService) GetUserByIDAndRole(id uint, targetRole string, actor policies.AuthContext) (*user.UserResponse, error) {
	target, err := s.repo.GetByIDAndRole(id, targetRole)
	if err != nil {
		return nil, err
	}

	if err := s.policy.View(actor, target); err != nil {
		return nil, err
	}

	// Passaggio del filtro vuoto user.UserFilter{} in quanto per i collegati non si applica la ricerca/paginazione
	linkedUsers, err := s.repo.GetAllByForeignID(target.ID, user.UserFilter{})
	if err != nil {
		return nil, err
	}

	linkedUsers, err = s.filterLinkedUsersForActor(linkedUsers, actor)
	if err != nil {
		return nil, err
	}

	response := mappers.ToUserResponse(target)
	response.LinkedUsers = mappers.ToUserResponses(linkedUsers)
	return &response, nil
}

// filterLinkedUsersForActor applica lo scoping ai collegati (agenti/agenzie/utenti
// con foreign_id = target.ID) in base al ruolo dell'attore che ha effettuato la
// richiesta. Per l'OPERATOR lo scope è dato dalle pivot agent_operator / agency_operator.
func (s *userService) filterLinkedUsersForActor(linkedUsers []models.User, actor policies.AuthContext) ([]models.User, error) {
	switch actor.Role {
	case enums.RoleAdmin.String():
		// L'ADMIN vede tutto, nessun filtro
		return linkedUsers, nil

	case enums.RoleOperator.String():
		agentIDs, err := s.scopeRepo.AssignedAgentIDs(actor.UserID)
		if err != nil {
			return nil, err
		}
		agencyIDs, err := s.scopeRepo.AssignedAgencyIDs(actor.UserID)
		if err != nil {
			return nil, err
		}
		agentSet := toSet(agentIDs)
		agencySet := toSet(agencyIDs)

		filtered := make([]models.User, 0, len(linkedUsers))
		for _, u := range linkedUsers {
			switch u.Role {
			case enums.RoleAgent:
				if agentSet[u.ID] {
					filtered = append(filtered, u)
				}
			case enums.RoleAgency:
				if agencySet[u.ID] {
					filtered = append(filtered, u)
				}
			case enums.RoleUser:
				// Uno USER è visibile solo se la sua agenzia (ForeignID) è
				// assegnata all'operatore
				if u.ForeignID != nil && agencySet[*u.ForeignID] {
					filtered = append(filtered, u)
				}
			}
		}
		return filtered, nil

	case enums.RoleAgent.String():
		scopeIDs, err := s.scopeRepo.NodeChildrenAndSelfAgentIds(actor.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return []models.User{}, nil
			}
			return nil, err
		}
		scopeSet := toSet(scopeIDs)

		filtered := make([]models.User, 0, len(linkedUsers))
		for _, u := range linkedUsers {
			switch u.Role {
			case enums.RoleAgent:
				if scopeSet[u.ID] {
					filtered = append(filtered, u)
				}
			case enums.RoleAgency:
				if u.ForeignID != nil && scopeSet[*u.ForeignID] {
					filtered = append(filtered, u)
				}
			case enums.RoleUser:
				filtered = append(filtered, u) // eventualmente da restringere ulteriormente
			}
		}
		return filtered, nil

	default:
		// AGENCY / USER: già scoperti dalla propria policy.View sul target,
		// nessun filtro aggiuntivo necessario
		return linkedUsers, nil
	}
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
		ActorID:     &actor.UserID,
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
		ActorID:     &actor.UserID,
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
		ActorID:     &actor.UserID,
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
		ActorID:     &actor.UserID,
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
		ActorID:     &actor.UserID,
		Action:      enums.Delete,
		TargetType:  reflect.TypeOf(existing).Elem().Name(),
		TargetID:    &id,
		Description: fmt.Sprintf("Eliminato l'utente '%s' (Ruolo: %s)", existing.Username, targetRole),
	})

	if err != nil {
		return err
	}

	return nil
}
