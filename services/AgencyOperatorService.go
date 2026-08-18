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
	AssignAgenciesToOperator(request *pivot.ArraysToOpRequest, actor policies.AuthContext) (*pivot.ArraysToOpResponse, error)
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
	if request == nil || request.AgencyId == nil {
		return nil, errors.New("agency_id mancante nella request")
	}

	_, err := s.userRepo.GetByIDAndRole(*request.AgencyId, enums.RoleAgency.String())
	if err != nil {
		return nil, errors.New("L'utente selezionato non è un'agenzia")
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

func (s *agencyOperatorService) AssignAgenciesToOperator(request *pivot.ArraysToOpRequest, actor policies.AuthContext) (*pivot.ArraysToOpResponse, error) {
	if request == nil {
		return nil, errors.New("request payload non valido")
	}
	model := mappers.ToArrAgencyOperatorModel(request)

	newPivots, err := s.repo.AssignAgenciesMassive(model)
	if err != nil {
		return nil, err
	}

	response := mappers.ToArrAgencyOperatorResponse(newPivots)
	return response, nil
}
