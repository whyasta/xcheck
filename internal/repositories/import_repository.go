package repositories

import (
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"

	"gorm.io/gorm"
)

type ImportRepository interface {
	GetDB() *gorm.DB
	Save(role *models.Import) (models.Import, error)
	Update(id int64, data *map[string]interface{}) (models.Import, error)
	Delete(uid int64) (models.Import, error)
	FindAll(paginate *utils.Paginate, filter []utils.Filter, sorts []utils.Sort) ([]models.Import, int64, error)
	FindByID(uid int64) (models.Import, error)
	CheckValidImport(uid int64) (bool, error)
	CheckValidAssign(uid int64) (bool, error)
}

type importRepository struct {
	base BaseRepository
}

func NewImportRepository(db *gorm.DB) *importRepository {
	return &importRepository{
		base: NewBaseRepository(db, models.Import{}),
	}
}

func (repo *importRepository) GetDB() *gorm.DB {
	return repo.base.GetDB()
}

func (repo *importRepository) Save(role *models.Import) (models.Import, error) {
	return BaseInsert(*repo.base.GetDB(), *role)
}

func (repo *importRepository) FindByID(id int64) (models.Import, error) {
	return BaseFindByID[models.Import](*repo.base.GetDB(), "imports", id, []string{})
}

func (repo *importRepository) FindAll(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Import, int64, error) {
	return BasePaginateWithFilter[[]models.Import](*repo.base.GetDB(), []string{}, paginate, filters, sorts)
}

func (repo *importRepository) Delete(id int64) (models.Import, error) {
	return BaseSoftDelete[models.Import](*repo.base.GetDB(), id)
}

func (repo *importRepository) Update(id int64, data *map[string]interface{}) (models.Import, error) {
	return BaseUpdate[models.Import](*repo.base.GetDB(), id, data)
}

func (repo *importRepository) CheckValidImport(id int64) (bool, error) {
	var result models.Import
	err := repo.base.GetDB().
		Where("id = ?", id).
		Where("status = ?", constant.ImportStatusCompleted).
		First(&result).
		Error

	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *importRepository) CheckValidAssign(id int64) (bool, error) {
	var result models.GateAllocation
	err := repo.base.GetDB().
		Where("id = ?", id).
		First(&result).
		Error

	if err != nil {
		return false, err
	}
	return true, nil
}
