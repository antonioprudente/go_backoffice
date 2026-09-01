package repositories

import (
	"example/go_backoffice/models"

	"gorm.io/gorm"
)

type NoteRepo interface {
	WithTx(tx *gorm.DB) NoteRepo
	NewNote(note models.Note) (*models.Note, error)
	GetByID(id uint) (*models.Note, error)
	GetAll() ([]*models.Note, error)
	GetAllByActorID(actorID uint) ([]*models.Note, error)
	Update(note *models.Note) error
	Delete(id uint) (bool, error)
}

type noteRepo struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) NoteRepo {
	return &noteRepo{db: db}
}

func (r *noteRepo) WithTx(tx *gorm.DB) NoteRepo {
	return &noteRepo{db: tx}
}

func (r *noteRepo) NewNote(note models.Note) (*models.Note, error) {
	err := r.db.Create(&note).Error
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *noteRepo) GetByID(id uint) (*models.Note, error) {
	var note models.Note
	err := r.db.Where("id = ?", id).First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *noteRepo) GetAll() ([]*models.Note, error) {
	var notes []*models.Note
	err := r.db.Find(&notes).Error
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *noteRepo) GetAllByActorID(actorID uint) ([]*models.Note, error) {
	var notes []*models.Note
	err := r.db.Where("actor_id = ?", actorID).Find(&notes).Error
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *noteRepo) Update(note *models.Note) error {
	return r.db.Save(note).Error
}

func (r *noteRepo) Delete(id uint) (bool, error) {
	result := r.db.Where("id = ?", id).Delete(&models.Note{})
	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
