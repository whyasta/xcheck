package repositories

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type EventRepository interface {
	Save(event *models.Event) (models.Event, error)
	Update(id int64, event *map[string]interface{}) (models.Event, error)
	Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.Event, int64, error)
	GetFiltered(paginate *utils.Paginate, filters []utils.Filter) ([]models.Event, int64, error)
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

func (repo *eventRepository) Save(event *models.Event) (models.Event, error) {
	if event.ID != 0 {
		var result = models.Event{}
		var err = repo.base.GetDB().Save(&event).First(&result).Error
		return result, err
	}
	return BaseInsert(*repo.base.GetDB(), *event)
}

func (repo *eventRepository) Update(id int64, event *map[string]interface{}) (models.Event, error) {
	// var result = models.Event{}
	// var result models.Event

	// var err = repo.base.GetDB().Model(&result).
	// 	Table("events").
	// 	Clauses(clause.Returning{}).
	// 	Where("id = ?", id).
	// 	Updates(event).
	// 	First(&result).
	// 	Error
	// return result, err
	return BaseUpdate[models.Event](*repo.base.GetDB(), id, event)
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
	// var events []models.Event
	// var count int64

	// tx := repo.base.GetDB().
	// 	Scopes(paginate.PaginatedResult).
	// 	Where(params).
	// 	Find(&events)

	// err := tx.Error

	// if err != nil {
	// 	return nil, 0, err
	// }

	// tx.Limit(-1).Offset(-1)
	// tx.Count(&count)

	// return events, count, nil
	return BasePaginate[[]models.Event](*repo.base.GetDB(), paginate, params)
}

func (repo *eventRepository) Delete(id int64) (models.Event, error) {
	return BaseSoftDelete[models.Event](*repo.base.GetDB(), id)
}

func (repo *eventRepository) GetFiltered(paginate *utils.Paginate, filters []utils.Filter) ([]models.Event, int64, error) {
	/*var events []models.Event
	var count int64
	// log.Println(filters)

	tx := repo.base.GetDB().
		Scopes(paginate.PaginatedResult)

	if len(filters) > 0 {
		for _, filter := range filters {
			newFilter := utils.NewFilter(filter.Property, filter.Operation, filter.Collation, filter.Value, filter.Items)
			tx = tx.Where(newFilter.FilterResult("", repo.base.GetDB()))
		}
	}

	tx = tx.Find(&events)

	err := tx.Error

	if err != nil {
		return nil, 0, err
	}

	if len(filters) <= 0 {
		tx.Limit(-1).Offset(-1)
		tx.Count(&count)
	}

	return events, count, nil*/
	return BasePaginateWithFilter[[]models.Event](*repo.base.GetDB(), []string{"Gates", "Schedules"}, paginate, filters)
}
