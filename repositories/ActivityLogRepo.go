package repositories

import (
	"example/go_backoffice/models"

	"gorm.io/gorm"
)

type ActivityLogRepo interface {
	Create(log *models.ActivityLog) error
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
