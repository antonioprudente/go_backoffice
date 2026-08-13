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
}

type agencyOperatorService struct {
	repo repositories.AgencyOperatorRepo
}

func NewAgencyOperatorService(repo repositories.AgencyOperatorRepo) AgencyOperatorService {
	return &agencyOperatorService{repo: repo}
}

func (s *agencyOperatorService) Create(request *pivot.AssignToOpRequest, actor policies.AuthContext) (*pivot.AssignToOpResponse, error) {
	if actor.Role != enums.RoleAdmin.String() {
		return nil, errors.New("unauthorized: solo gli amministratori possono eseguire questa operazione")
	}

	if request == nil {
		return nil, errors.New("request payload non valido")
	}

	model := mappers.ToAgencyOperatorModel(request)

	newPivot, err := s.repo.AssignAgency(model)
	if err != nil {
		return nil, err
	}

	response := mappers.ToAgencyOperatorResponse(newPivot)

	return response, nil

}
