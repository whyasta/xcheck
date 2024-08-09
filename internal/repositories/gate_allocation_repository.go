package repositories

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"

	"gorm.io/gorm"
)

type GateAllocationRepository interface {
	Save(data *models.GateAllocation) (models.GateAllocation, error)
	Update(id int64, data *map[string]interface{}) (models.GateAllocation, error)
	Delete(uid int64) (models.GateAllocation, error)
	FindAll(paginate *utils.Paginate, filter []utils.Filter, sorts []utils.Sort) ([]models.GateAllocation, int64, error)
	FindByID(uid int64) (models.GateAllocation, error)
}

type gateAllocationRepository struct {
	base BaseRepository
}

func NewGateAllocationRepository(db *gorm.DB) *gateAllocationRepository {
	return &gateAllocationRepository{
		base: NewBaseRepository(db, models.GateAllocation{}),
	}
}

func (repo *gateAllocationRepository) Save(data *models.GateAllocation) (models.GateAllocation, error) {
	return BaseInsert(*repo.base.GetDB(), *data)
}

func (repo *gateAllocationRepository) FindByID(id int64) (models.GateAllocation, error) {
	return BaseFindByID[models.GateAllocation](*repo.base.GetDB(), "gateAllocations", id, []string{})
}

func (repo *gateAllocationRepository) FindAll(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.GateAllocation, int64, error) {
	//return BasePaginateWithFilter[[]models.GateAllocation](*repo.base.GetDB(), []string{"Session", "Gate", "Event"}, paginate, filters, sorts)
	return BasePaginateWithFilter[[]models.GateAllocation](*repo.base.GetDB(), []string{}, paginate, filters, sorts)
}

func (repo *gateAllocationRepository) Delete(id int64) (models.GateAllocation, error) {
	return BaseSoftDelete[models.GateAllocation](*repo.base.GetDB(), id)
}

func (repo *gateAllocationRepository) Update(id int64, data *map[string]interface{}) (models.GateAllocation, error) {
	return BaseUpdate[models.GateAllocation](*repo.base.GetDB(), id, data)
}
