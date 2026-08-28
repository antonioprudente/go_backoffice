package services

import (
	"example/go_backoffice/models"
	"example/go_backoffice/repositories"
)

type ActivityLogService interface {
	NewLog(log *models.ActivityLog) error
	GetActivity() ([]*models.ActivityLog, error)
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

func (s *activityLogService) GetActivity() ([]*models.ActivityLog, error) {
	activity, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return activity, err
}
