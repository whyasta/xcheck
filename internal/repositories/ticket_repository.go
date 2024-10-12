package repositories

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type TicketRepository interface {
	GetFiltered(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Ticket, int64, error)
	FindByOrderID(eventID int64, orderID string) (models.Ticket, error)
	Exist(eventID int64, orderBarcode string) (bool, error)
	GetImport(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Import, int64, error)
	GetImportDetail(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Ticket, int64, error)
	ValidateImport(importID int64, eventID int64) error
	ValidateRecord(eventID int64, row []string) (bool, error)
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *ticketRepository {
	return &ticketRepository{
		db: db,
	}
}

func (repo *ticketRepository) GetImport(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Import, int64, error) {
	result, count, err := BasePaginateWithFilter[[]models.Import](*repo.db, []string{}, paginate, filters, sorts)

	for i := 0; i < len(result); i++ {
		result[i].FileName = "https://" + config.GetAppConfig().MinioEndpoint + "/" + config.GetAppConfig().MinioBucket + "/" + result[i].FileName
	}

	return result, count, err
}

func (repo *ticketRepository) GetImportDetail(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Ticket, int64, error) {
	return BasePaginateWithFilter[[]models.Ticket](*repo.db, []string{"TicketType"}, paginate, filters, sorts)
}

func (repo *ticketRepository) GetFiltered(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Ticket, int64, error) {
	return BasePaginateWithFilter[[]models.Ticket](*repo.db, []string{"TicketType"}, paginate, filters, sorts)
}

func (repo *ticketRepository) FindByOrderID(eventID int64, orderID string) (models.Ticket, error) {
	var result models.Ticket
	tx := repo.db.Model(&result)
	err := tx.First(&result, "event_id = ? AND order_id = ?", eventID, orderID).Error
	return result, err
}

func (repo *ticketRepository) Exist(eventID int64, orderBarcode string) (bool, error) {
	var exists bool
	err := repo.db.Model(&models.Ticket{}).Select("1").Where("event_id = ? AND order_barcode = ?", eventID, orderBarcode).Limit(1).Find(&exists).Error
	return exists, err
}

func (repo *ticketRepository) ValidateRecord(eventID int64, row []string) (bool, error) {
	item := models.Ticket{
		ImportID:       "9999",
		EventID:        eventID,
		OrderBarcode:   row[0],
		OrderID:        row[1],
		TicketTypeName: row[2],
		Name:           row[3],
		Email:          row[4],
		PhoneNumber:    row[5],
	}

	var result models.TicketType
	err := repo.db.Table("ticket_types").
		Where("event_id = ? AND ticket_type_name = ?", item.EventID, item.TicketTypeName).
		First(&result).Error
	if err != nil {
		return false, errors.New(fmt.Sprintf("Invalid ticket type %s", item.TicketTypeName))
	}
	item.TicketTypeID = &result.ID

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(&item)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, nil
}

func (repo *ticketRepository) ValidateImport(importID int64, eventID int64) error {
	result := repo.db.
		Table("tickets").
		Where("import_id = ?", importID).
		Updates(map[string]interface{}{"event_id": eventID})
	err := result.Error

	var importTickets []models.Ticket

	repo.db.Table("tickets").
		Where("import_id = ?", importID).
		Find(&importTickets)

	var failedValues []string

	count := 0
	failedCount := 0
	for _, item := range importTickets {
		// check valid ticket types
		if item.TicketTypeID == nil {
			var result models.TicketType
			err = repo.db.Table("ticket_types").
				Where("event_id = ? AND ticket_type_name = ?", item.EventID, item.TicketTypeName).
				First(&result).Error
			if err != nil {
				item.Note = fmt.Sprintf("Ticket type not found: %s", item.TicketTypeName)

				repo.db.
					Table("tickets").
					Where("id = ?", item.ID).
					Updates(item)
				failedCount++

				failedValues = append(failedValues, item.OrderBarcode+"|"+item.OrderID+"|"+item.TicketTypeName)

				continue
			}
			item.TicketTypeID = &result.ID
		}

		item.AssignStatus = 1
		item.Quantity = 1

		validate := validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(&item)
		if err != nil {
			fmt.Println(err)
			failedCount++
			continue
		}

		repo.db.
			Table("tickets").
			Where("id = ?", item.ID).
			Updates(item)

		count++
	}

	err = repo.db.Table("imports").
		Where("id = ?", importID).
		Updates(map[string]interface{}{
			"status":           constant.ImportStatusAssigned,
			"success_count":    count,
			"failed_count":     failedCount,
			"duplicate_count":  0,
			"failed_values":    strings.Join(failedValues[:], ","),
			"duplicate_values": "",
		}).Error
	return err
}
