package repositories

import (
	"bigmind/xcheck-be/internal/constant"
	"bigmind/xcheck-be/internal/constant/response"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BarcodeRepository interface {
	Save(role *models.Barcode) (models.Barcode, error)
	Update(id int64, data *map[string]interface{}) (models.Barcode, error)
	Delete(uid int64) (models.Barcode, error)
	FindAll(joins []string, paginate *utils.Paginate, filter []utils.Filter, sorts []utils.Sort) ([]models.Barcode, int64, error)
	FindAllWithRelations(eventID int64, paginate *utils.Paginate, filter []utils.Filter, sorts []utils.Sort) ([]models.Barcode, int64, error)
	FindByID(uid int64) (models.Barcode, error)
	AssignBarcodes(importID int64, assignID int64, ticketTypeID int64) (int64, error)
	Scan(eventID int64, barcode string) (models.Barcode, response.ResponseStatus, error)
	CreateLog(eventID int64, userID int64, gateID int64, ticketTypeID int64, sessionID int64, barcode string, currentStatus constant.BarcodeStatus, action constant.BarcodeStatus, device string, reason string) (models.BarcodeLog, bool, error)
	CreateBulkLog(barcodes *[]models.BarcodeLog) error
	CreateBulk(barcodes *[]models.Barcode) error
	GetLatestScan(eventID int64, barcode string) (models.BarcodeLog, error)

	UpdateSingleBarcode(eventID int64, barcodeID int64, sessions []int64, gates []int64) error
	UpdateBulkBarcode(eventID int64, ticketTypeID int64, sessions []int64, gates []int64) error
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

func (repo *barcodeRepository) FindAllWithRelations(eventID int64, paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Barcode, int64, error) {
	var records []models.Barcode
	var count int64

	tx := repo.base.GetDB().Debug().
		Table("barcodes").
		Scopes(paginate.PaginatedResult)

	if len(filters) > 0 {
		for _, filter := range filters {
			newFilter := utils.NewFilter(filter.Property, filter.Operation, filter.Collation, filter.Value, filter.Items)
			tx = newFilter.FilterResult("", tx)
		}
	}

	tx = tx.Preload("TicketType", func(tx2 *gorm.DB) *gorm.DB {
		return tx2.Omit("EventID")
	}).Preload("Gates", func(tx2 *gorm.DB) *gorm.DB {
		return tx2.Omit("EventID")
	}).Preload("Sessions", func(tx2 *gorm.DB) *gorm.DB {
		return tx2.Omit("EventID")
	}).Preload("LatestScan", func(tx2 *gorm.DB) *gorm.DB {
		return tx2.Where("barcode_logs.event_id = ?", eventID).
			Joins("JOIN users ON users.id = barcode_logs.scanned_by").
			Joins("JOIN gates ON gates.id = barcode_logs.gate_id").
			Select("users.username as scanned_by_name", "barcode_logs.*", "gates.gate_name")
	}).Preload("History", func(tx2 *gorm.DB) *gorm.DB {
		return tx2.Where("barcode_logs.event_id = ?", eventID).
			Joins("JOIN users ON users.id = barcode_logs.scanned_by").
			Joins("JOIN gates ON gates.id = barcode_logs.gate_id").
			Select("users.username as scanned_by_name", "barcode_logs.*", "gates.gate_name")
	})

	if len(sorts) > 0 {
		for _, sort := range sorts {
			newSort := utils.NewSort(sort.Property, sort.Direction)
			tx = newSort.SortResult(tx)
		}
	}

	tx = tx.Find(&records)

	// fmt.Println(filters)
	// if len(filters) <= 0 {
	// }
	tx.Limit(-1).Offset(-1)
	tx.Count(&count)

	err := tx.Error

	return records, count, err
}

func (repo *barcodeRepository) Delete(id int64) (models.Barcode, error) {
	return BaseSoftDelete[models.Barcode](*repo.base.GetDB(), id)
}

func (repo *barcodeRepository) Update(id int64, data *map[string]interface{}) (models.Barcode, error) {
	return BaseUpdate[models.Barcode](*repo.base.GetDB(), id, data)
}

func (repo *barcodeRepository) AssignBarcodes(importID int64, assignID int64, ticketTypeID int64) (int64, error) {
	var importBarcodes []models.ImportBarcode

	var err error
	var count int64

	// Begin transaction
	repo.base.GetDB().Transaction(func(tx *gorm.DB) error {
		result := repo.base.GetDB().
			Table("raw_barcodes").
			Where("import_id = ?", importID).
			Updates(map[string]interface{}{"assign_status": 1})

		err = result.Error
		count = result.RowsAffected

		repo.base.GetDB().
			Table("raw_barcodes").
			Where("import_id = ?", importID).
			Find(&importBarcodes)

		fmt.Println(importBarcodes)

		// for each barcode
		barcodes := []models.Barcode{}
		for _, item := range importBarcodes {
			barcodes = append(barcodes, models.Barcode{
				Barcode: item.Barcode,
				// GateAllocationID: assignID,
				TicketTypeID:  ticketTypeID,
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
			Where("id = ?", importID).
			Updates(map[string]interface{}{"status": constant.ImportStatusAssigned})

		err = result.Error
		return err
	})

	return count, err
}

func (repo *barcodeRepository) AssignBarcodesWithEvent(importID int64, eventID int64, ticketTypeID int64, sessions []int64, gates []int64) (int64, int64, int64, []string, []string, error) {
	var importBarcodes []models.ImportBarcode

	var err error
	var count int64
	var failedCount int64
	var duplicateCount int64

	var failedValues []string
	var duplicateValues []string

	// Begin transaction
	repo.base.GetDB().Transaction(func(tx *gorm.DB) error {
		result := repo.base.GetDB().
			Table("raw_barcodes").
			Where("import_id = ?", importID).
			Updates(map[string]interface{}{"assign_status": 1})

		err = result.Error

		repo.base.GetDB().
			Table("raw_barcodes").
			Where("import_id = ?", importID).
			Find(&importBarcodes)

		// fmt.Println(importBarcodes)

		// for each barcode
		barcodes := []models.Barcode{}
		for _, item := range importBarcodes {
			var exists bool
			err = repo.base.GetDB().Table("barcodes").Select("count(*) > 0").Where("event_id = ? AND barcode = ?", eventID, item.Barcode).
				Find(&exists).
				Error
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("exists: %s %v\n", item.Barcode, exists)
			if exists {
				duplicateCount++
				duplicateValues = append(duplicateValues, item.Barcode)
				continue
			}

			rowBarcode := models.Barcode{
				Barcode:       item.Barcode,
				EventID:       eventID,
				TicketTypeID:  ticketTypeID,
				Flag:          constant.BarcodeFlagValid,
				CurrentStatus: constant.BarcodeStatusNull,
			}

			validate := validator.New(validator.WithRequiredStructEnabled())
			validate.RegisterValidation("barcode", utils.BarcodeValidation)
			err := validate.Struct(&rowBarcode)
			if err != nil {
				fmt.Println(err)
				failedCount++
				failedValues = append(failedValues, item.Barcode)
				continue
			}

			barcodes = append(barcodes, rowBarcode)
		}

		if len(barcodes) > 0 {

			result = repo.base.GetDB().Omit("Sessions", "Gates").
				Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "event_id"}, {Name: "barcode"}},
					DoNothing: true,
				}).
				Create(&barcodes)

			err = result.Error
			if err != nil {
				fmt.Println(err)
			}

			for _, item := range barcodes {
				var barcode = models.Barcode{
					ID: item.ID,
				}

				gateIds := []models.Gate{}
				for _, gateItem := range gates {
					var gate = models.Gate{
						ID: gateItem,
					}
					// gateIds = append(gateIds, gate)

					repo.base.GetDB().
						Table("gates").
						Where("id = ?", gateItem).
						First(&gate)

					gateIds = append(gateIds, gate)
				}
				err = repo.base.GetDB().Debug().
					Model(&barcode).
					Association("Gates").
					Replace(gateIds)
				if err != nil {
					fmt.Println(err)
				}

				sessionIds := []models.Session{}
				for _, sessItem := range sessions {
					var session = models.Session{
						ID: sessItem,
					}

					repo.base.GetDB().
						Table("sessions").
						Where("id = ?", sessItem).
						First(&session)

					sessionIds = append(sessionIds, session)
				}

				err = repo.base.GetDB().Debug().
					Model(&barcode).
					Omit("Sessions.session_start").
					Omit("Sessions.session_end").
					Association("Sessions").
					Replace(sessionIds)
				if err != nil {
					fmt.Println(err)
				}
			}

			count = result.RowsAffected
		}

		fmt.Println("failed", strings.Join(failedValues[:], ","))
		fmt.Println("duplicate", strings.Join(duplicateValues[:], ","))

		result = repo.base.GetDB().
			Table("imports").
			Where("id = ?", importID).
			Updates(map[string]interface{}{
				"status":           constant.ImportStatusAssigned,
				"success_count":    count,
				"failed_count":     failedCount,
				"duplicate_count":  duplicateCount,
				"failed_values":    strings.Join(failedValues[:], ","),
				"duplicate_values": strings.Join(duplicateValues[:], ","),
			})

		err = result.Error
		return err
	})

	// fmt.Println(count, failedCount, duplicateCount)

	return count, failedCount, duplicateCount, failedValues, duplicateValues, err
}

func (repo *barcodeRepository) Scan(eventID int64, barcode string) (models.Barcode, response.ResponseStatus, error) {
	var result models.Barcode

	tx := repo.base.GetDB().
		Table("barcodes")

	tx = tx.Preload("TicketType", func(tx2 *gorm.DB) *gorm.DB {
		return tx2.Omit("EventID")
	}).Preload("Gates", func(tx2 *gorm.DB) *gorm.DB {
		return tx2.Omit("EventID")
	}).Preload("Sessions", func(tx2 *gorm.DB) *gorm.DB {
		return tx2.Omit("EventID")
	})

	err := tx.Where("barcode = ?", barcode).
		Where("event_id = ?", eventID).
		First(&result).Error
	if err != nil {
		return result, response.EC01, errors.New("Barcode " + barcode + " not found")
	}

	return result, response.Checkin, err
}

func (repo *barcodeRepository) CreateLog(eventID int64, userID int64, gateID int64, ticketTypeID int64,
	sessionID int64, barcode string, currentStatus constant.BarcodeStatus,
	action constant.BarcodeStatus, device string, reason string) (models.BarcodeLog, bool, error) {
	// action := constant.BarcodeStatusIn
	firstCheckin := false
	if currentStatus == constant.BarcodeStatusNull {
		firstCheckin = true
	}

	log := models.BarcodeLog{
		Barcode:      barcode,
		Action:       action,
		ScannedAt:    time.Now(),
		ScannedBy:    userID,
		EventID:      eventID,
		GateID:       gateID,
		Device:       device,
		TicketTypeID: ticketTypeID,
		SessionID:    sessionID,
		Reason:       reason,
	}

	var err = repo.base.GetDB().Table("barcode_logs").Create(&log).Error
	if err == nil && reason == "" {
		err = repo.base.GetDB().
			Table("barcodes").
			Where("barcode = ?", barcode).
			Updates(map[string]interface{}{"current_status": action, "flag": constant.BarcodeFlagUsed}).
			Error
		// Update("current_status", action).
		// Update("flag", constant.BarcodeFlagUsed).

	}
	return log, firstCheckin, err
}

func (repo *barcodeRepository) GetLatestScan(eventID int64, barcode string) (models.BarcodeLog, error) {
	var result models.BarcodeLog
	err := repo.base.GetDB().
		Table("barcode_logs").
		Where("event_id = ?", eventID).
		Where("barcode = ?", barcode).
		Where("reason = ?", "").
		Order("scanned_at DESC").
		First(&result).Error
	return result, err
}

func (repo *barcodeRepository) CreateBulkLog(logs *[]models.BarcodeLog) error {
	var barcodeList []string
	uniqueMap := make(map[string]bool)

	var err error
	if (len(*logs)) > 0 {
		err := repo.base.GetDB().Table("barcode_logs").Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "event_id"}, {Name: "barcode"}, {Name: "scanned_at"}, {Name: "action"}},
		}).Omit("ID").Create(&logs).Error

		if err != nil {
			return err
		}

		for _, item := range *logs {
			if !uniqueMap[item.Barcode] && (item.Reason == "" || item.Reason == "null") {
				uniqueMap[item.Barcode] = true
				barcodeList = append(barcodeList, item.Barcode)
			}
		}
	}

	var barcodes []models.Barcode
	subQuery1 := repo.base.GetDB().Model(&models.BarcodeLog{}).Select("max(scanned_at)").Where("barcode IN (?)", barcodeList).Group("barcode")
	repo.base.GetDB().Table("barcode_logs").
		Select("barcode", "action as current_status", "scanned_at").
		Where("scanned_at IN (?)", subQuery1).
		Where("barcode IN (?)", barcodeList).
		Find(&barcodes)

	for _, item := range barcodes {
		err := repo.base.GetDB().Table("barcodes").Where("barcode = ?", item.Barcode).
			Updates(map[string]interface{}{"current_status": item.CurrentStatus, "flag": constant.BarcodeFlagUsed}).Error
		// Update("current_status", item.CurrentStatus).
		// Update("flag", constant.BarcodeFlagUsed).Error
		if err != nil {
			return err
		}
	}

	// err = repo.base.GetDB().Table("barcodes").Where("barcode IN (?)", barcodeList).
	// 	Update("flag", constant.BarcodeFlagUsed).Error
	// if err != nil {
	// 	return err
	// }

	return err
}

func (repo *barcodeRepository) CreateBulk(barcodes *[]models.Barcode) error {
	err := repo.base.GetDB().Table("barcodes").
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "event_id"}, {Name: "barcode"}},
			DoUpdates: clause.AssignmentColumns([]string{"current_status", "flag"}),
		}, clause.Returning{}).Omit("ID", "CreatedAt", "Gates", "Sessions").Create(&barcodes).Error

	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, item := range *barcodes {
		// err := repo.base.GetDB().Table("barcodes").
		// 	Clauses(clause.OnConflict{
		// 		Columns:   []clause.Column{{Name: "event_id"}, {Name: "barcode"}},
		// 		DoUpdates: clause.AssignmentColumns([]string{"current_status", "flag"}),
		// 	}, clause.Returning{}).Omit("ID", "CreatedAt", "Gates", "Sessions").Create(&item).Error

		// if err != nil {
		// 	fmt.Println(err)
		// 	return err
		// }

		fmt.Println(item.ID)
		if item.ID == 0 {
			var temp models.Barcode
			repo.base.GetDB().Table("barcodes").Where("event_id = ? AND barcode = ?", item.EventID, item.Barcode).
				Find(&temp)
			item.ID = temp.ID
		}

		// updated / created
		if item.ID > 0 {
			err = repo.base.GetDB().Debug().
				Model(&item).
				Association("Gates").
				Replace(item.Gates)
			if err != nil {
				fmt.Println(err)
			}

			err = repo.base.GetDB().Debug().
				Model(&item).
				Association("Sessions").
				Replace(item.Sessions)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func (repo *barcodeRepository) UpdateSingleBarcode(eventID int64, barcodeID int64, sessions []int64, gates []int64) error {
	tx := repo.base.GetDB().Begin()

	var barcodeModel models.Barcode
	err := tx.Debug().Where("event_id = ? AND id = ?", eventID, barcodeID).First(&barcodeModel).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if barcodeModel.Flag == constant.BarcodeFlagUsed {
		tx.Rollback()
		return errors.New("Barcode " + barcodeModel.Barcode + " already used")
	}

	gateIds := []models.Gate{}
	err = tx.Debug().Table("gates").
		Where("id IN (?)", gates).
		Find(&gateIds).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	sessionIds := []models.Session{}
	err = tx.Debug().Table("sessions").
		Where("id IN (?)", sessions).
		Find(&sessionIds).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var paramGates string
	for _, gate := range gateIds {
		if paramGates == "" {
			paramGates = fmt.Sprintf("('%v', '%v')", barcodeID, gate.ID)
		} else {
			paramGates = fmt.Sprintf("%s, ('%v', '%v')", paramGates, barcodeID, gate.ID)
		}
	}
	err = tx.Debug().
		Model(&barcodeModel).
		Association("Gates").
		Clear()
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Debug().Exec("INSERT INTO barcode_gates (barcode_id, gate_id) VALUES " + paramGates).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var paramSessions string
	for _, session := range sessionIds {
		if paramSessions == "" {
			paramSessions = fmt.Sprintf("('%v', '%v')", barcodeID, session.ID)
		} else {
			paramSessions = fmt.Sprintf("%s, ('%v', '%v')", paramSessions, barcodeID, session.ID)
		}
	}
	err = tx.Debug().
		Model(&barcodeModel).
		Association("Sessions").
		Clear()
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Debug().Exec("INSERT INTO barcode_sessions (barcode_id, session_id) VALUES " + paramSessions).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (repo *barcodeRepository) UpdateBulkBarcode(eventID int64, ticketTypeID int64, sessions []int64, gates []int64) error {
	tx := repo.base.GetDB().Begin()

	var ticketType models.TicketType
	err := tx.Debug().Where("event_id = ? AND id = ?", eventID, ticketTypeID).First(&ticketType).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	gateIds := []models.Gate{}
	err = tx.Debug().Table("gates").
		Where("id IN (?)", gates).
		Find(&gateIds).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	sessionIds := []models.Session{}
	err = tx.Debug().Table("sessions").
		Where("id IN (?)", sessions).
		Find(&sessionIds).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var barcodes []models.Barcode
	err = tx.Debug().Table("barcodes").Where("event_id = ? AND ticket_type_id = ? AND flag = 'VALID'", eventID, ticketTypeID).Find(&barcodes).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var paramGates string
	for _, item := range barcodes {
		for _, gate := range gateIds {
			if paramGates == "" {
				paramGates = fmt.Sprintf("('%v', '%v')", item.ID, gate.ID)
			} else {
				paramGates = fmt.Sprintf("%s, ('%v', '%v')", paramGates, item.ID, gate.ID)
			}
		}
	}
	err = tx.Debug().
		Model(&barcodes).
		Association("Gates").
		Clear()
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Debug().Exec("INSERT INTO barcode_gates (barcode_id, gate_id) VALUES " + paramGates).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	var paramSessions string
	for _, item := range barcodes {
		for _, session := range sessionIds {
			if paramSessions == "" {
				paramSessions = fmt.Sprintf("('%v', '%v')", item.ID, session.ID)
			} else {
				paramSessions = fmt.Sprintf("%s, ('%v', '%v')", paramSessions, item.ID, session.ID)
			}
		}
	}
	err = tx.Debug().
		Model(&barcodes).
		Association("Sessions").
		Clear()
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Debug().Exec("INSERT INTO barcode_sessions (barcode_id, session_id) VALUES " + paramSessions).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
