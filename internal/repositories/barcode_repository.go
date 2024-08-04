package repositories

import (
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BarcodeRepository interface {
	Save(role *models.Barcode) (models.Barcode, error)
	Update(id int64, data *map[string]interface{}) (models.Barcode, error)
	Delete(uid int64) (models.Barcode, error)
	FindAll(joins []string, paginate *utils.Paginate, filter []utils.Filter, sorts []utils.Sort) ([]models.Barcode, int64, error)
	FindByID(uid int64) (models.Barcode, error)
	AssignBarcodes(importId int64, assignId int64, ticketTypeId int64) (int64, error)
	Scan(barcode string) (models.Barcode, error)
	CreateLog(eventId int64, userId int64, barcode string, currentStatus constant.BarcodeStatus, action constant.BarcodeStatus) (bool, error)
	CreateBulkLog(barcodes *[]models.BarcodeLog) error
}

type barcodeRepository struct {
	base BaseRepository
}

func NewBarcodeRepository(db *gorm.DB) *barcodeRepository {
	return &barcodeRepository{
		base: NewBaseRepository(db, models.Barcode{}),
	}
}

func (repo *barcodeRepository) Save(role *models.Barcode) (models.Barcode, error) {
	return BaseInsert(*repo.base.GetDB(), *role)
}

func (repo *barcodeRepository) FindByID(id int64) (models.Barcode, error) {
	return BaseFindByID[models.Barcode](*repo.base.GetDB(), "barcodes", id, []string{})
}

func (repo *barcodeRepository) FindAll(joins []string, paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Barcode, int64, error) {
	return BasePaginateWithFilter[[]models.Barcode](*repo.base.GetDB(), joins, paginate, filters, sorts)
}

func (repo *barcodeRepository) Delete(id int64) (models.Barcode, error) {
	return BaseSoftDelete[models.Barcode](*repo.base.GetDB(), id)
}

func (repo *barcodeRepository) Update(id int64, data *map[string]interface{}) (models.Barcode, error) {
	return BaseUpdate[models.Barcode](*repo.base.GetDB(), id, data)
}

func (repo *barcodeRepository) AssignBarcodes(importId int64, assignId int64, ticketTypeId int64) (int64, error) {
	var importBarcodes []models.ImportBarcode

	var err error
	var count int64

	// Begin transaction
	repo.base.GetDB().Transaction(func(tx *gorm.DB) error {
		result := repo.base.GetDB().
			Table("raw_barcodes").
			Where("import_id = ?", importId).
			Updates(map[string]interface{}{"assign_status": 1})

		err = result.Error
		count = result.RowsAffected

		repo.base.GetDB().
			Table("raw_barcodes").
			Where("import_id = ?", importId).
			Find(&importBarcodes)

		fmt.Println(importBarcodes)

		// for each barcode
		barcodes := []models.Barcode{}
		for _, item := range importBarcodes {
			barcodes = append(barcodes, models.Barcode{
				Barcode:       item.Barcode,
				ScheduleID:    assignId,
				TicketTypeID:  ticketTypeId,
				Flag:          constant.BarcodeFlagValid,
				CurrentStatus: constant.BarcodeStatusNull,
			})
		}

		result = repo.base.GetDB().
			Table("barcodes").
			Create(barcodes)

		err = result.Error
		count = result.RowsAffected

		result = repo.base.GetDB().
			Table("imports").
			Where("id = ?", importId).
			Updates(map[string]interface{}{"status": constant.ImportStatusPaired})

		err = result.Error
		return err
	})

	return count, err
}

func (repo *barcodeRepository) Scan(barcode string) (models.Barcode, error) {
	var result models.Barcode
	err := repo.base.GetDB().
		Joins("Schedule").
		Joins("Schedule.Session").
		Where("barcode = ?", barcode).
		First(&result).
		Error
	if err != nil {
		return result, errors.New("barcode not found")
	}

	return result, err
}

func (repo *barcodeRepository) CreateLog(eventId int64, userId int64, barcode string, currentStatus constant.BarcodeStatus, action constant.BarcodeStatus) (bool, error) {
	// action := constant.BarcodeStatusIn
	firstCheckin := false
	if currentStatus == constant.BarcodeStatusNull {
		action = constant.BarcodeStatusIn
		firstCheckin = true
	} else if currentStatus == constant.BarcodeStatusIn {
		// action = constant.BarcodeStatusOut //tidak ada checkout
		action = constant.BarcodeStatusIn
	}

	log := models.BarcodeLog{
		Barcode:   barcode,
		Action:    action,
		ScannedAt: time.Now(),
		ScannedBy: userId,
		EventID:   eventId,
	}

	var err = repo.base.GetDB().Table("barcode_logs").Create(&log).Error
	if err == nil {
		err = repo.base.GetDB().
			Table("barcodes").
			Where("barcode = ?", barcode).
			Update("current_status", action).
			Update("flag", constant.BarcodeFlagUsed).
			Error
	}
	return firstCheckin, err
}

func (repo *barcodeRepository) CreateBulkLog(barcodes *[]models.BarcodeLog) error {
	var err = repo.base.GetDB().Table("barcode_logs").Create(&barcodes).Error
	return err
}
