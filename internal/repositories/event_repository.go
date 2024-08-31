package repositories

import (
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"

	"gorm.io/gorm"
)

type EventRepository interface {
	Save(event *dto.EventRequest) (models.Event, error)
	BulkSave(events *[]dto.EventRequest) ([]models.Event, error)
	Update(id int64, event *map[string]interface{}) (models.Event, error)
	Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.Event, int64, error)
	GetFiltered(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Event, int64, error)
	Delete(uid int64) (models.Event, error)
	FindByID(uid int64) (models.Event, error)
	Summary(uid int64) dto.EventSummary
	Report(uid int64) (dto.EventReportResponse, error)
	GateTicketTypes(uid int64, gateIds []int64) []dto.EventGateTicketTypeResponse
}

type eventRepository struct {
	base BaseRepository
}

func NewEventRepository(db *gorm.DB) *eventRepository {
	return &eventRepository{
		base: NewBaseRepository(db, models.Event{}),
	}
}

func (repo *eventRepository) Save(eventDto *dto.EventRequest) (models.Event, error) {
	if eventDto.ID != 0 {
		var result = models.Event{}
		var err = repo.base.GetDB().Table("events").Save(&eventDto).First(&result).Error
		return result, err
	}

	event := eventDto.ToEntity()
	return BaseInsert(*repo.base.GetDB(), *event)
}

func (repo *eventRepository) BulkSave(eventDto *[]dto.EventRequest) ([]models.Event, error) {
	event := make([]models.Event, 0)
	for _, v := range *eventDto {
		event = append(event, *v.ToEntity())
	}
	return BaseInsert(*repo.base.GetDB(), event)
}

func (repo *eventRepository) Update(id int64, event *map[string]interface{}) (models.Event, error) {
	// var result = models.Event{}
	// var result models.Event

	// var err = repo.base.GetDB().Model(&result).
	// 	Table("events").
	// 	Clauses(clause.Returning{}).
	// 	Where("id = ?", id).
	// 	Updates(event).
	// 	First(&result).
	// 	Error
	// return result, err
	return BaseUpdate[models.Event](*repo.base.GetDB(), id, event)
}

func (repo *eventRepository) FindByID(id int64) (models.Event, error) {
	return BaseFindByID[models.Event](*repo.base.GetDB(), "events", id, []string{"TicketTypes", "Gates", "Sessions"})
	// data, err := repo.base.CommonFindByID("events", id)
	// if err != nil {
	// 	return models.Event{}, err
	// }

	// jsonStr, err := json.Marshal(data)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return models.Event{}, err
	// }
	// var event models.Event
	// if err := json.Unmarshal(jsonStr, &event); err != nil {
	// 	fmt.Println(err)
	// 	return models.Event{}, err
	// }
	// return event, nil
}

func (repo *eventRepository) Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.Event, int64, error) {
	// var events []models.Event
	// var count int64

	// tx := repo.base.GetDB().
	// 	Scopes(paginate.PaginatedResult).
	// 	Where(params).
	// 	Find(&events)

	// err := tx.Error

	// if err != nil {
	// 	return nil, 0, err
	// }

	// tx.Limit(-1).Offset(-1)
	// tx.Count(&count)

	// return events, count, nil
	return BasePaginate[[]models.Event](*repo.base.GetDB(), paginate, params)
}

func (repo *eventRepository) Delete(id int64) (models.Event, error) {
	return BaseSoftDelete[models.Event](*repo.base.GetDB(), id)
}

func (repo *eventRepository) GetFiltered(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Event, int64, error) {
	/*var events []models.Event
	var count int64
	// fmt.Println(filters)

	tx := repo.base.GetDB().
		Scopes(paginate.PaginatedResult)

	if len(filters) > 0 {
		for _, filter := range filters {
			newFilter := utils.NewFilter(filter.Property, filter.Operation, filter.Collation, filter.Value, filter.Items)
			tx = tx.Where(newFilter.FilterResult("", repo.base.GetDB()))
		}
	}

	tx = tx.Find(&events)

	err := tx.Error

	if err != nil {
		return nil, 0, err
	}

	if len(filters) <= 0 {
		tx.Limit(-1).Offset(-1)
		tx.Count(&count)
	}

	return events, count, nil*/
	return BasePaginateWithFilter[[]models.Event](*repo.base.GetDB(), []string{"TicketTypes", "Gates", "Sessions"}, paginate, filters, sorts)
}

func (repo *eventRepository) Summary(id int64) dto.EventSummary {
	var result = dto.EventSummary{}
	// var totalCheckIn int64
	// var totalCheckOut int64

	// subQuery := repo.base.GetDB().Select("gateAllocation_id").Where("event_id = ?", id).Table("gateAllocations")
	err := repo.base.GetDB().Table("barcodes").
		Select("count(id) as total_barcode",
			"SUM(CASE WHEN current_status = 'IN' THEN 1 ELSE 0 END) as ongoing_check_in",
			"SUM(CASE WHEN current_status = 'OUT' THEN 1 ELSE 0 END) as ongoing_check_out").
		Where("event_id = ?", id).
		Omit("total_ticket_type").
		Scan(&result).
		Error
	if err != nil {
		return dto.EventSummary{}
	}

	type UniqueCheck struct {
		TotalCheckIn  int64 `json:"total_check_in" mapstructure:"total_check_in"`
		TotalCheckOut int64 `json:"total_check_out" mapstructure:"total_check_out"`
	}

	var uniqueCheck UniqueCheck
	err = repo.base.GetDB().Table("barcode_logs").
		Select("SUM(CASE WHEN action = 'IN' THEN 1 ELSE 0 END) as total_check_in",
			"SUM(CASE WHEN action = 'OUT' THEN 1 ELSE 0 END) as total_check_out").
		Where("event_id = ?", id).
		Scan(&uniqueCheck).
		Error
	if err != nil {
		return dto.EventSummary{}
	}

	result.TotalCheckIn = uniqueCheck.TotalCheckIn
	result.TotalCheckOut = uniqueCheck.TotalCheckOut

	/*
		var ticketType []struct {
			TicketTypeName string `json:"ticket_type_name" mapstructure:"ticket_type_name"`
			TotalBarcode   int    `json:"total_barcode" mapstructure:"total_barcode"`
			TotalCheckIn   int    `json:"total_check_in" mapstructure:"total_check_in"`
			TotalCheckOut  int    `json:"total_check_out" mapstructure:"total_check_out"`
		}

		var jsonTicketType []map[string]interface{}
		repo.base.GetDB().Table("ticket_types").
			Select("ticket_type_name",
				"count(barcode) as total_barcode",
				"SUM(CASE WHEN current_status = 'IN' THEN 1 ELSE 0 END) as total_check_in",
				"SUM(CASE WHEN current_status = 'OUT' THEN 1 ELSE 0 END) as total_check_out").
			Joins("join barcodes on barcodes.ticket_type_id = ticket_types.id").
			Where("barcodes.event_id = ?", id).
			Group("barcodes.ticket_type_id").
			Scan(&ticketType)

		mapstructure.Decode(ticketType, &jsonTicketType)
		result.TotalTicketType = jsonTicketType
	*/

	// repo.base.GetDB().Table("barcode_logs").
	// 	Select("count(barcode) as total_check_in").
	// 	Where("event_id = ?", id).
	// 	Where("action = ?", "IN").
	// 	Scan(&totalCheckIn)

	// repo.base.GetDB().Table("barcode_logs").
	// 	Select("count(barcode) as total_check_in").
	// 	Where("event_id = ?", id).
	// 	Where("action = ?", "OUT").
	// 	Scan(&totalCheckOut)

	// result.TotalCheckIn = totalCheckIn
	// result.TotalCheckOut = totalCheckOut
	return result
}

func (repo *eventRepository) Report(id int64) (dto.EventReportResponse, error) {
	var result = dto.EventReportResponse{}

	tx := repo.base.GetDB().Table("events")
	err := tx.First(&result, "id = ?", id).Error
	//data, err := BaseFindByID[map[string]interface{}](*repo.base.GetDB(), "events", id, []string{""})

	result.EventSummary = repo.Summary(id)

	return result, err
}

func (repo *eventRepository) GateTicketTypes(eventId int64, gateIds []int64) []dto.EventGateTicketTypeResponse {
	var result = []dto.EventGateTicketTypeResponse{}
	repo.base.GetDB().Raw("SELECT DISTINCT bg.gate_id, g.gate_name, barcodes.ticket_type_id, tt.ticket_type_name FROM barcodes "+
		"JOIN ticket_types tt ON tt.id = barcodes.ticket_type_id "+
		"JOIN barcode_gates bg ON bg.barcode_id = barcodes.id "+
		"JOIN gates g ON g.id = bg.gate_id "+
		"WHERE bg.gate_id IN (?) and barcodes.event_id = ? GROUP BY bg.gate_id, tt.id", gateIds, eventId).
		Scan(&result)
	return result
	//result := repo.base.GetDB().
}
