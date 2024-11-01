package repositories

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/dto"
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
	Exist(eventID int64, orderBarcode string) (bool, error)
	GetImport(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Import, int64, error)
	GetImportDetail(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Ticket, int64, error)
	ValidateImport(importID int64, eventID int64) error
	ValidateRecord(eventID int64, row []string) (bool, error)
	// redemption
	FindByOrderID(eventID int64, orderID string) (models.Ticket, error)
	FindByBarcode(eventID int64, orderBarcode string) (models.Ticket, error)
	Redeem(eventID int64, generateBarcode bool, photo string, note *string, data []dto.TicketRedeemDataRequest) ([]models.Ticket, error)
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

func (repo *ticketRepository) FindByBarcode(eventID int64, orderBarcode string) (models.Ticket, error) {
	var result models.Ticket
	tx := repo.db.Model(&result)
	err := tx.First(&result, "event_id = ? AND order_barcode = ?", eventID, orderBarcode).Error

	if result.Status == constant.TicketStatusRedeemed {
		return result, errors.New("Ticket already redeemed")
	}

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

func (repo *ticketRepository) Redeem(eventID int64, generateBarcode bool, photo string, note *string, data []dto.TicketRedeemDataRequest) ([]models.Ticket, error) {
	var result []models.Ticket
	var errorList []string

	tx := repo.db.Begin()
	for _, item := range data {
		var ticket models.Ticket
		err := tx.Model(&models.Ticket{}).First(&ticket, "event_id = ? AND id = ?", eventID, item.ID).Error
		if err != nil {
			errorList = append(errorList, fmt.Sprintf("Ticket %v not found", item.ID))
			continue
		}

		// check if ticket already redeemed
		if ticket.Status == constant.TicketStatusRedeemed {
			errorList = append(errorList, fmt.Sprintf("Ticket %s already redeemed", ticket.OrderBarcode))
			continue
		} else if ticket.Status == constant.TicketStatusExpired {
			errorList = append(errorList, fmt.Sprintf("Ticket %s already expired", ticket.OrderBarcode))
			continue
		} else if ticket.Status == constant.TicketStatusCanceled {
			errorList = append(errorList, fmt.Sprintf("Ticket %s already canceled", ticket.OrderBarcode))
			continue
		}

		// check associated barcode
		if generateBarcode == false {
			// var barcode models.Barcode
			// err := repo.db.Model(&models.Barcode{}).First(&barcode, "event_id = ? AND ticket_type_id = ? AND barcode = ?", eventID, ticket.TicketTypeID, item.AssociateBarcode).Error
			// if err != nil {
			// 	errorList = append(errorList, fmt.Sprintf("Barcode %s not found", item.AssociateBarcode))
			// }
		} else {
			// var barcode models.Barcode
			// barcode, err = BarcodeRepo.Generate(eventID, ticket.TicketTypeID, item.AssociateBarcode)
			// if err != nil {
			//     errorList = append(errorList, fmt.Sprintf("Barcode %s not found", item.AssociateBarcode))
			// }
		}

		ticket.AssociateBarcode = &item.AssociateBarcode
		ticket.Status = constant.TicketStatusRedeemed
		ticket.PhotoUrl = &photo

		result = append(result, ticket)
	}
	if len(errorList) > 0 {
		tx.Rollback()
		return nil, errors.New(strings.Join(errorList, ", "))
	}

	for _, item := range result {
		err := tx.Model(&item).
			Updates(map[string]interface{}{"note": note, "status": constant.TicketStatusRedeemed, "associate_barcode": item.AssociateBarcode, "photo_url": photo}).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return result, nil
}
