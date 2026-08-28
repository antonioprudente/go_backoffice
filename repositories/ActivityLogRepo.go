package repositories

import (
	"example/go_backoffice/models"

	"gorm.io/gorm"
)

type ActivityLogRepo interface {
	Create(log *models.ActivityLog) error
	GetAll() ([]*models.ActivityLog, error)
}

type activityLogRepo struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) ActivityLogRepo {
	return &activityLogRepo{db: db}
}

func (r *activityLogRepo) Create(log *models.ActivityLog) error {
	return r.db.Create(log).Error
}

func (r *activityLogRepo) GetAll() ([]*models.ActivityLog, error) {
	var logs []*models.ActivityLog
	err := r.db.Find(&logs).Error
	return logs, err
}
