package repositories

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"

	"gorm.io/gorm"
)

type GateRepository interface {
	Save(gate *models.Gate) (models.Gate, error)
	BulkSave(gates *[]models.Gate) ([]models.Gate, error)
	Update(id int64, data *map[string]interface{}) (models.Gate, error)
	Delete(uid int64) (models.Gate, error)
	FindAll(paginate *utils.Paginate, filter []utils.Filter, sorts []utils.Sort) ([]models.Gate, int64, error)
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

func (repo *gateRepository) Save(gate *models.Gate) (models.Gate, error) {
	return BaseInsert(*repo.base.GetDB(), *gate)
}

func (repo *gateRepository) BulkSave(gates *[]models.Gate) ([]models.Gate, error) {
	return BaseInsert(*repo.base.GetDB(), *gates)
}

func (repo *gateRepository) FindByID(id int64) (models.Gate, error) {
	return BaseFindByID[models.Gate](*repo.base.GetDB(), "gates", id, []string{})
}

func (repo *gateRepository) FindAll(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Gate, int64, error) {
	return BasePaginateWithFilter[[]models.Gate](*repo.base.GetDB(), []string{}, paginate, filters, sorts)
}

func (repo *gateRepository) Delete(id int64) (models.Gate, error) {
	return BaseSoftDelete[models.Gate](*repo.base.GetDB(), id)
}

func (repo *gateRepository) Update(id int64, data *map[string]interface{}) (models.Gate, error) {
	return BaseUpdate[models.Gate](*repo.base.GetDB(), id, data)
}
