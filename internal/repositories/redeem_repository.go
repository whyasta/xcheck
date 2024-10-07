package repositories

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"

	"gorm.io/gorm"
)

type RedeemRepository interface {
	GetFiltered(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Redeem, int64, error)
	FindByOrderID(eventID int64, orderID string) (models.Redeem, error)
}

type redeemRepository struct {
	db *gorm.DB
}

func NewRedeemRepository(db *gorm.DB) *redeemRepository {
	return &redeemRepository{
		db: db,
	}
}

func (repo *redeemRepository) GetFiltered(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Redeem, int64, error) {
	return BasePaginateWithFilter[[]models.Redeem](*repo.db, []string{"TicketType"}, paginate, filters, sorts)
}

func (repo *redeemRepository) FindByOrderID(eventID int64, orderID string) (models.Redeem, error) {
	var result models.Redeem
	tx := repo.db.Model(&result)
	err := tx.First(&result, "event_id = ? AND order_id = ?", eventID, orderID).Error
	return result, err
}
