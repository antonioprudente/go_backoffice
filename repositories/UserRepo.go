package repositories

import (
	"fmt"
	"strconv"

	"example/go_backoffice/dto/user"
	"example/go_backoffice/models"

	"gorm.io/gorm"
)

type UserRepo interface {
	WithTx(tx *gorm.DB) UserRepo

	GetAllByRole(role string, filter user.UserFilter) ([]models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByIDAndRole(id uint, role string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	UpdateStatusByIdAndRole(id uint, role string, status string) (*models.User, error)
	UpdateForeignID(id uint, role string, foreignID uint) (*models.User, error)

	DeleteByIdAndRole(id uint, role string) error

	GetAllByRoleAndIDs(role string, ids []uint, filter user.UserFilter) ([]models.User, error)
	GetAllByRoleAndForeignIDs(role string, foreignIDs []uint, filter user.UserFilter) ([]models.User, error)
	GetAllByForeignID(foreignID uint, filter user.UserFilter) ([]models.User, error)
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

func (r *userRepo) GetAllByRole(role string, filter user.UserFilter) ([]models.User, error) {
	var users []models.User
	q := r.db.Where("role = ?", role)
	q = applyUserFilter(q, filter)
	err := q.Find(&users).Error
	return users, err
}

func (r *userRepo) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByIDAndRole(id uint, role string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Foreign").Where("id = ?", id).Where("role = ?", role).First(&user).Error; err != nil {
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

func (r *userRepo) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) UpdateStatusByIdAndRole(id uint, role string, status string) (*models.User, error) {
	result := r.db.Model(&models.User{}).
		Where("id = ? AND role = ?", id, role).
		Update("status", status)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var updatedUser models.User
	if err := r.db.Where("id = ? AND role = ?", id, role).First(&updatedUser).Error; err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (r *userRepo) DeleteByIdAndRole(id uint, role string) error {
	result := r.db.Where("id = ? AND role = ?", id, role).Delete(&models.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func parseID(id string) (uint, error) {
	parsed, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("id non valido: %w", err)
	}
	return uint(parsed), nil
}

func (r *userRepo) GetAllByRoleAndIDs(role string, ids []uint, filter user.UserFilter) ([]models.User, error) {
	var users []models.User
	if len(ids) == 0 {
		return users, nil
	}
	q := r.db.Where("role = ? AND id IN ?", role, ids)
	q = applyUserFilter(q, filter)
	err := q.Find(&users).Error
	return users, err
}

func (r *userRepo) GetAllByRoleAndForeignIDs(role string, foreignIDs []uint, filter user.UserFilter) ([]models.User, error) {
	var users []models.User
	if len(foreignIDs) == 0 {
		return users, nil
	}
	q := r.db.Where("role = ? AND foreign_id IN ?", role, foreignIDs)
	q = applyUserFilter(q, filter)
	err := q.Find(&users).Error
	return users, err
}

func (r *userRepo) GetAllByForeignID(foreignID uint, filter user.UserFilter) ([]models.User, error) {
	var users []models.User
	q := r.db.Where("foreign_id = ?", foreignID)
	q = applyUserFilter(q, filter)
	err := q.Find(&users).Error
	return users, err
}

func (r *userRepo) UpdateForeignID(id uint, role string, foreignID uint) (*models.User, error) {
	result := r.db.Model(&models.User{}).
		Where("id = ? AND role = ?", id, role).
		Update("foreign_id", foreignID)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return r.GetByIDAndRole(id, role)
}

func applyUserFilter(q *gorm.DB, filter user.UserFilter) *gorm.DB {
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		q = q.Where("first_name LIKE ? OR last_name LIKE ? OR username LIKE ? OR email LIKE ?", like, like, like, like)
	}
	return q
}
