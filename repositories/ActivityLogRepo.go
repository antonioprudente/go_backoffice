package repositories

import (
	"example/go_backoffice/models"

	"gorm.io/gorm"
)

type ActivityLogRepo interface {
	Create(log *models.ActivityLog) error
	GetAll() ([]*models.ActivityLog, error)
	GetByID(id uint) (*models.ActivityLog, error)
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
	err := r.db.Order("created_at desc").Find(&logs).Error
	return logs, err
}

func (r *activityLogRepo) GetByID(id uint) (*models.ActivityLog, error) {
	var log models.ActivityLog
	err := r.db.Preload("Actor").Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}
