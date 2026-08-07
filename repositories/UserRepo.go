package repositories

import (
	"example/go_backoffice/models"

	"gorm.io/gorm"
)

type UserRepo interface {
	GetAll() ([]models.User, error)
	GetByID(id string) (*models.User, error)
	Create(user *models.User) error
	Delete(id string) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepo) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) Delete(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}
