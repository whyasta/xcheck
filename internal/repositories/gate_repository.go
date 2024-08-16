package repositories

import (
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GateRepository interface {
	Save(gate *dto.GateRequestDto) (models.Gate, error)
	BulkSave(gates *[]dto.GateRequestDto) ([]models.Gate, error)
	Update(id int64, data *dto.GateRequestDto) (models.Gate, error)
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

func (repo *gateRepository) Save(gate *dto.GateRequestDto) (models.Gate, error) {
	var entity = gate.ToEntity()
	result, err := BaseInsert(*repo.base.GetDB(), *entity)
	return result, err
}

func (repo *gateRepository) BulkSave(gates *[]dto.GateRequestDto) ([]models.Gate, error) {
	// result, err := BaseInsert(*repo.base.GetDB(), *gates)
	var result []models.Gate
	for _, item := range *gates {
		result = append(result, *item.ToEntity())
	}

	var err = repo.base.GetDB().Table("gates").Omit("users").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"gate_name"}),
	}).Create(&result).Error

	return result, err
}

func (repo *gateRepository) FindByID(id int64) (models.Gate, error) {
	var result models.Gate
	tx := repo.base.GetDB().Model(&result)
	tx = tx.Preload("Users", func(tx2 *gorm.DB) *gorm.DB {
		return tx2.Omit("Password", "AuthUUIDs")
	})
	err := tx.First(&result, "id = ?", id).Error
	return result, err

	//return BaseFindByID[models.Gate](*repo.base.GetDB().Omit("users.auth_uuids", "users.password"), "gates", id, []string{"Users"})
}

func (repo *gateRepository) FindAll(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Gate, int64, error) {
	return BasePaginateWithFilter[[]models.Gate](*repo.base.GetDB(), []string{"Users"}, paginate, filters, sorts)
}

func (repo *gateRepository) Delete(id int64) (models.Gate, error) {
	return BaseSoftDelete[models.Gate](*repo.base.GetDB(), id)
}

func (repo *gateRepository) Update(id int64, data *dto.GateRequestDto) (models.Gate, error) {
	var result = models.Gate{
		ID: id,
	}

	var users []models.User
	for _, userID := range data.Users {
		users = append(users, models.User{
			ID: *userID,
		})
	}

	repo.base.GetDB().
		Model(&result).
		Association("Users").
		Replace(users)

	err := repo.base.GetDB().Model(&result).
		Omit(clause.Associations).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(data).
		Error
	return result, err

	//return BaseUpdate[models.Gate](*repo.base.GetDB(), id, data)
}
