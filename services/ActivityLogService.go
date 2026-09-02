package services

import (
	"example/go_backoffice/dto/log"
	"example/go_backoffice/mappers"
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"
)

type ActivityLogService interface {
	NewLog(log *models.ActivityLog) error
	GetActivity() ([]*log.LogResponse, error)
	GetLogDetailByID(id uint) (*log.LogResponse, error)
}

type activityLogService struct {
	repo repositories.ActivityLogRepo
}

func NewActivityLogService(repo repositories.ActivityLogRepo) ActivityLogService {
	return &activityLogService{repo: repo}
}

func (s *activityLogService) NewLog(log *models.ActivityLog) error {
	return s.repo.Create(log)
}

func (s *activityLogService) GetActivity() ([]*log.LogResponse, error) {
	logs, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return mappers.ToLogResponses(logs), nil
}

func (s *activityLogService) GetLogDetailByID(id uint) (*log.LogResponse, error) {
	log, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return mappers.ToLogResponse(log), nil
}
