package repositories

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"

	"gorm.io/gorm"
)

type EventRepository interface {
	Save(role *models.Event) (models.Event, error)
	Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.Event, int64, error)
	Delete(uid int64) (models.Event, error)
	FindByID(uid int64) (models.Event, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *eventRepository {
	return &eventRepository{
		db: db,
	}
}

func (repo *eventRepository) Save(role *models.Event) (models.Event, error) {
	return BaseInsert(*repo.db, *role)
}

func (repo *eventRepository) FindByID(id int64) (models.Event, error) {
	return BaseFindByID[models.Event](*repo.db, id)
}

func (repo *eventRepository) Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.Event, int64, error) {
	var events []models.Event
	var count int64

	tx := repo.db.
		Scopes(paginate.PaginatedResult).
		Where(params).
		Find(&events)

	err := tx.Error

	if err != nil {
		return nil, 0, err
	}

	tx.Limit(-1).Offset(-1)
	tx.Count(&count)

	return events, count, nil
}

func (repo *eventRepository) Delete(id int64) (models.Event, error) {
	return BaseSoftDelete[models.Event](*repo.db, id)
}
