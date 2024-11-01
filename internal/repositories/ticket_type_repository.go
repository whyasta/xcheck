package repositories

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketTypeRepository interface {
	Save(ticketType *models.TicketType) (models.TicketType, error)
	BulkSave(ticketTypes *[]models.TicketType) ([]models.TicketType, error)
	Update(id int64, data *map[string]interface{}) (models.TicketType, error)
	Delete(uid int64) (models.TicketType, error)
	FindAll(paginate *utils.Paginate, filter []utils.Filter, sorts []utils.Sort) ([]models.TicketType, int64, error)
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

func (repo *ticketTypeRepository) Save(ticketType *models.TicketType) (models.TicketType, error) {
	return BaseInsert(*repo.base.GetDB(), *ticketType)
}

func (repo *ticketTypeRepository) BulkSave(ticketTypes *[]models.TicketType) ([]models.TicketType, error) {
	//return BaseInsert(*repo.base.GetDB(), *ticketTypes)
	var err = repo.base.GetDB().Table("ticket_types").Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
	}).Create(&ticketTypes).Error
	return *ticketTypes, err
}

func (repo *ticketTypeRepository) FindByID(id int64) (models.TicketType, error) {
	return BaseFindByID[models.TicketType](*repo.base.GetDB(), "ticket_types", id, []string{})
}

func (repo *ticketTypeRepository) FindAll(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.TicketType, int64, error) {
	return BasePaginateWithFilter[[]models.TicketType](*repo.base.GetDB(), []string{}, paginate, filters, sorts)
}

func (repo *ticketTypeRepository) Delete(id int64) (models.TicketType, error) {
	return BaseSoftDelete[models.TicketType](*repo.base.GetDB(), id)
}

func (repo *ticketTypeRepository) Update(id int64, data *map[string]interface{}) (models.TicketType, error) {
	return BaseUpdate[models.TicketType](*repo.base.GetDB(), id, data)
}
