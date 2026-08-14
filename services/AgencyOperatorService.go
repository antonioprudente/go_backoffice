package services

import (
	"errors"
	"example/go_backoffice/dto/pivot"
	"example/go_backoffice/enums"
	"example/go_backoffice/mappers"
	"example/go_backoffice/policies"
	"example/go_backoffice/repositories"
)

type AgencyOperatorService interface {
	AssignAgencyToOperator(request *pivot.AssignToOpRequest, actor policies.AuthContext) (*pivot.AssignToOpResponse, error)
}

type agencyOperatorService struct {
	repo     repositories.AgencyOperatorRepo
	userRepo repositories.UserRepo
}

func NewAgencyOperatorService(
	repo repositories.AgencyOperatorRepo,
	userRepo repositories.UserRepo,
) AgencyOperatorService {
	return &agencyOperatorService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *agencyOperatorService) AssignAgencyToOperator(request *pivot.AssignToOpRequest, actor policies.AuthContext) (*pivot.AssignToOpResponse, error) {
	if actor.Role != enums.RoleAdmin.String() {
		return nil, errors.New("unauthorized: solo gli amministratori possono eseguire questa operazione")
	}

	if request == nil {
		return nil, errors.New("request payload non valido")
	}
	_, err := s.userRepo.GetByIDAndRole(*request.AgentId, enums.RoleAgency.String())
	if err != nil {
		return nil, errors.New("L'utente selezionato non è un agenzia")
	}

	_, err = s.userRepo.GetByIDAndRole(request.OperatorId, enums.RoleOperator.String())
	if err != nil {
		return nil, errors.New("L'utente selezionato non è un operatore")
	}
	model := mappers.ToAgencyOperatorModel(request)

	newPivot, err := s.repo.AssignAgency(model)
	if err != nil {
		return nil, err
	}

	response := mappers.ToAgencyOperatorResponse(newPivot)
	return response, nil
}
