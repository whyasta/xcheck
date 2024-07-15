package repositories

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type EventRepository interface {
	Save(role *models.Event) (models.Event, error)
	Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.Event, int64, error)
	Delete(uid int64) (models.Event, error)
	FindByID(uid int64) (models.Event, error)
}

type eventRepository struct {
	base BaseRepository
}

func NewEventRepository(db *gorm.DB) *eventRepository {
	return &eventRepository{
		base: NewBaseRepository(db, models.Event{}),
	}
}

func (repo *eventRepository) Save(role *models.Event) (models.Event, error) {
	return BaseInsert(*repo.base.GetDB(), *role)
}

func (repo *eventRepository) FindByID(id int64) (models.Event, error) {
	//return BaseFindByID[models.Event](*repo.base.GetDB(), id)
	data, err := repo.base.CommonFindByID("events", id)
	if err != nil {
		return models.Event{}, err
	}

	jsonStr, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return models.Event{}, err
	}
	var event models.Event
	if err := json.Unmarshal(jsonStr, &event); err != nil {
		fmt.Println(err)
		return models.Event{}, err
	}
	return event, nil
}

func (repo *eventRepository) Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.Event, int64, error) {
	var events []models.Event
	var count int64

	tx := repo.base.GetDB().
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
	return BaseSoftDelete[models.Event](*repo.base.GetDB(), id)
}
