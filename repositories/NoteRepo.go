package repositories

import (
	"example/go_backoffice/models"

	"gorm.io/gorm"
)

type NoteRepo interface {
	WithTx(tx *gorm.DB) NoteRepo
	NewNote(note models.Note) (*models.Note, error)
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
