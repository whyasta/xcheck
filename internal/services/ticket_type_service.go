package services

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
	"errors"
	"strconv"
)

type TicketTypeService struct {
	r repositories.TicketTypeRepository
	b repositories.BaseRepository
}

func NewTicketTypeService(r repositories.TicketTypeRepository, b repositories.BaseRepository) *TicketTypeService {
	return &TicketTypeService{r, b}
}

func (s *TicketTypeService) CreateTicketType(data *models.TicketType) (models.TicketType, error) {
	return s.r.Save(data)
}

func (s *TicketTypeService) CreateBulkTicketType(types *[]models.TicketType) ([]models.TicketType, error) {
	return s.r.BulkSave(types)
}

func (s *TicketTypeService) UpdateTicketType(eventID int64, id int64, data *map[string]interface{}) (models.TicketType, error) {
	var filters = []utils.Filter{
		{
			Property:  "event_id",
			Operation: "=",
			Value:     strconv.Itoa(int(eventID)),
		},
		{
			Property:  "id",
			Operation: "=",
			Value:     strconv.Itoa(int(id)),
		},
	}
	rows, _, _ := s.r.FindAll(utils.NewPaginate(1, 0), filters, []utils.Sort{})

	if len(rows) == 0 {
		return models.TicketType{}, errors.New("record not found")
	}

	return s.r.Update(id, data)
	// fmt.Println("dsadas")
	// _, err := s.b.CommonUpdate("ticket_types", map[string]interface{}{"id": id}, data)
	// var items models.TicketType
	// mapstructure.Decode(result, &items)
	// return models.TicketType{}, err
}

func (s *TicketTypeService) GetAllTicketTypes(pageParams *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.TicketType, int64, error) {
	return s.r.FindAll(pageParams, filters, sorts)
}

func (s *TicketTypeService) GetTicketTypeByID(uid int64) (models.TicketType, error) {
	return s.r.FindByID(uid)
}

func (s *TicketTypeService) Delete(uid int64) (models.TicketType, error) {
	return s.r.Delete(uid)
}
