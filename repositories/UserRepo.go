package repositories

import (
	"fmt"
	"strconv"

	"example/go_backoffice/models"

	"gorm.io/gorm"
)

type UserRepo interface {
	GetAll() ([]models.User, error)
	GetByID(id string) (*models.User, error)
	Create(user *models.User) error
	UpdateStatus(id uint, status string) (*models.User, error)
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
	uid, err := parseID(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := r.db.First(&user, uid).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Implementazione
func (r *userRepo) UpdateStatus(id uint, status string) (*models.User, error) {
	if err := r.db.Model(&models.User{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return nil, err
	}
	var updatedUser models.User
	if err := r.db.First(&updatedUser, id).Error; err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (r *userRepo) Delete(id string) error {
	uid, err := parseID(id)
	if err != nil {
		return err
	}
	return r.db.Delete(&models.User{}, uid).Error
}

func parseID(id string) (uint, error) {
	parsed, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("id non valido: %w", err)
	}
	return uint(parsed), nil
}
