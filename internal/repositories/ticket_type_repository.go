package repositories

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"

	"gorm.io/gorm"
)

type TicketTypeRepository interface {
	Save(role *models.TicketType) (models.TicketType, error)
	Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.TicketType, int64, error)
	Delete(uid int64) (models.TicketType, error)
	FindByID(uid int64) (models.TicketType, error)
}

type ticketTypeRepository struct {
	base BaseRepository
}

func NewTicketTypeRepository(db *gorm.DB) *ticketTypeRepository {
	return &ticketTypeRepository{
		base: NewBaseRepository(db, models.TicketType{}),
	}
}

func (repo *ticketTypeRepository) Save(role *models.TicketType) (models.TicketType, error) {
	return BaseInsert(*repo.base.GetDB(), *role)
}

func (repo *ticketTypeRepository) FindByID(id int64) (models.TicketType, error) {
	return BaseFindByID[models.TicketType](*repo.base.GetDB(), id)
}

func (repo *ticketTypeRepository) Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.TicketType, int64, error) {
	var ticketTypes []models.TicketType
	var count int64

	tx := repo.base.GetDB().
		Scopes(paginate.PaginatedResult).
		Where(params).
		Find(&ticketTypes)

	err := tx.Error

	if err != nil {
		return nil, 0, err
	}

	tx.Limit(-1).Offset(-1)
	tx.Count(&count)

	return ticketTypes, count, nil
}

func (repo *ticketTypeRepository) Delete(id int64) (models.TicketType, error) {
	return BaseSoftDelete[models.TicketType](*repo.base.GetDB(), id)
}
