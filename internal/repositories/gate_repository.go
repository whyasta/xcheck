package repositories

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"

	"gorm.io/gorm"
)

type GateRepository interface {
	Save(role *models.Gate) (models.Gate, error)
	Delete(uid int64) (models.Gate, error)
	FindAll(paginate *utils.Paginate, filter []utils.Filter) ([]models.Gate, int64, error)
	FindByID(uid int64) (models.Gate, error)
}

type gateRepository struct {
	base BaseRepository
}

func NewGateRepository(db *gorm.DB) *gateRepository {
	return &gateRepository{
		base: NewBaseRepository(db, models.Gate{}),
	}
}

func (repo *gateRepository) Save(role *models.Gate) (models.Gate, error) {
	return BaseInsert(*repo.base.GetDB(), *role)
}

func (repo *gateRepository) FindByID(id int64) (models.Gate, error) {
	return BaseFindByID[models.Gate](*repo.base.GetDB(), id)
}

func (repo *gateRepository) FindAll(paginate *utils.Paginate, filters []utils.Filter) ([]models.Gate, int64, error) {
	return BasePaginateWithFilter[[]models.Gate](*repo.base.GetDB(), []string{}, paginate, filters)
}

func (repo *gateRepository) Delete(id int64) (models.Gate, error) {
	return BaseSoftDelete[models.Gate](*repo.base.GetDB(), id)
}
