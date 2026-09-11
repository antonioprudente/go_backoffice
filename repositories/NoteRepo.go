package repositories

import (
	"example/go_backoffice/dto/note"
	"example/go_backoffice/models"
	"time"

	"gorm.io/gorm"
)

type NoteRepo interface {
	WithTx(tx *gorm.DB) NoteRepo
	NewNote(note models.Note) (*models.Note, error)
	GetByID(id uint) (*models.Note, error)
	GetAll(filter note.NoteFilter) ([]*models.Note, error)
	GetAllByActorID(actorID uint, filter note.NoteFilter) ([]*models.Note, error)
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
	err := r.db.Preload("Actor").Preload("Target").Where("id = ?", id).First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *noteRepo) GetAll(filter note.NoteFilter) ([]*models.Note, error) {
	var notes []*models.Note

	q := r.db.Model(&models.Note{}).
		Order("notes.created_at desc").
		Preload("Actor").
		Preload("Target")

	q = applyNoteFilter(q, filter)

	if err := q.Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *noteRepo) GetAllByActorID(actorID uint, filter note.NoteFilter) ([]*models.Note, error) {
	var notes []*models.Note

	q := r.db.Model(&models.Note{}).
		Where("notes.actor_id = ?", actorID).
		Order("notes.created_at desc").
		Preload("Actor").
		Preload("Target")

	q = applyNoteFilter(q, filter)

	if err := q.Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func applyNoteFilter(q *gorm.DB, filter note.NoteFilter) *gorm.DB {

	if filter.CreateFrom != nil {
		from := time.Date(
			filter.CreateFrom.Year(), filter.CreateFrom.Month(), filter.CreateFrom.Day(),
			0, 0, 0, 0, filter.CreateFrom.Location(),
		)
		q = q.Where("notes.created_at >= ?", from)

		if filter.CreateTo != nil {
			to := time.Date(
				filter.CreateTo.Year(), filter.CreateTo.Month(), filter.CreateTo.Day(),
				0, 0, 0, 0, filter.CreateTo.Location(),
			).AddDate(0, 0, 1)
			q = q.Where("notes.created_at < ?", to)
		}
	}

	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		q = q.Joins("LEFT JOIN users AS actor ON actor.id = notes.actor_id").
			Joins("LEFT JOIN users AS target ON target.id = notes.target_id").
			Where("actor.username LIKE ? OR notes.content LIKE ? OR target.username LIKE ?", like, like, like)
	}

	return q
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
