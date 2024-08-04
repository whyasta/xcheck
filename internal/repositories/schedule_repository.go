package repositories

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"

	"gorm.io/gorm"
)

type ScheduleRepository interface {
	Save(data *models.Schedule) (models.Schedule, error)
	Update(id int64, data *map[string]interface{}) (models.Schedule, error)
	Delete(uid int64) (models.Schedule, error)
	FindAll(paginate *utils.Paginate, filter []utils.Filter, sorts []utils.Sort) ([]models.Schedule, int64, error)
	FindByID(uid int64) (models.Schedule, error)
}

type scheduleRepository struct {
	base BaseRepository
}

func NewScheduleRepository(db *gorm.DB) *scheduleRepository {
	return &scheduleRepository{
		base: NewBaseRepository(db, models.Schedule{}),
	}
}

func (repo *scheduleRepository) Save(data *models.Schedule) (models.Schedule, error) {
	return BaseInsert(*repo.base.GetDB(), *data)
}

func (repo *scheduleRepository) FindByID(id int64) (models.Schedule, error) {
	return BaseFindByID[models.Schedule](*repo.base.GetDB(), id, []string{})
}

func (repo *scheduleRepository) FindAll(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Schedule, int64, error) {
	return BasePaginateWithFilter[[]models.Schedule](*repo.base.GetDB(), []string{"Session", "Gate", "Event"}, paginate, filters, sorts)
}

func (repo *scheduleRepository) Delete(id int64) (models.Schedule, error) {
	return BaseSoftDelete[models.Schedule](*repo.base.GetDB(), id)
}

func (repo *scheduleRepository) Update(id int64, data *map[string]interface{}) (models.Schedule, error) {
	return BaseUpdate[models.Schedule](*repo.base.GetDB(), id, data)
}
