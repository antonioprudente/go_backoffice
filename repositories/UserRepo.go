package repositories

import (
	"fmt"
	"strconv"

	"example/go_backoffice/models"

	"gorm.io/gorm"
)

type UserRepo interface {
	WithTx(tx *gorm.DB) UserRepo

	GetAllByRole(role string) ([]models.User, error)
	GetByIDAndRole(id uint, role string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	UpdateStatusByIdAndRole(id uint, role string, status string) (*models.User, error)
	Delete(id string) error

	GetAllByRoleAndIDs(role string, ids []uint) ([]models.User, error)
	GetAllByRoleAndForeignIDs(role string, foreignIDs []uint) ([]models.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) WithTx(tx *gorm.DB) UserRepo {
	return &userRepo{db: tx}
}

func (r *userRepo) GetAllByRole(role string) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("role = ?", role).Find(&users).Error
	return users, err
}

func (r *userRepo) GetByIDAndRole(id uint, role string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).Where("role = ?", role).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Implementazione
func (r *userRepo) UpdateStatusByIdAndRole(id uint, role string, status string) (*models.User, error) {
	result := r.db.Model(&models.User{}).
		Where("id = ? AND role = ?", id, role).
		Update("status", status)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound // oppure un errore custom tipo "utente non trovato per questo ruolo"
	}

	var updatedUser models.User
	if err := r.db.Where("id = ? AND role = ?", id, role).First(&updatedUser).Error; err != nil {
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

// GetAllByRoleAndIDs filtra per ruolo e per un set di ID (usato per OPERATOR->agenti
// o AGENT->propri sottoagenti)
func (r *userRepo) GetAllByRoleAndIDs(role string, ids []uint) ([]models.User, error) {
	var users []models.User
	if len(ids) == 0 {
		return users, nil
	}
	err := r.db.Where("role = ? AND id IN ?", role, ids).Find(&users).Error
	return users, err
}

// GetAllByRoleAndForeignIDs filtra per ruolo (tipicamente AGENCY) e per un set di
// foreign_id (gli AgentID a cui sono agganciate)
func (r *userRepo) GetAllByRoleAndForeignIDs(role string, foreignIDs []uint) ([]models.User, error) {
	var users []models.User
	if len(foreignIDs) == 0 {
		return users, nil
	}
	err := r.db.Where("role = ? AND foreign_id IN ?", role, foreignIDs).Find(&users).Error
	return users, err
}
