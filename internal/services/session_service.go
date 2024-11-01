package services

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
	"errors"
	"strconv"
)

type SessionService struct {
	r repositories.SessionRepository
}

func NewSessionService(r repositories.SessionRepository) *SessionService {
	return &SessionService{r}
}

func (s *SessionService) CreateSession(data *models.Session) (models.Session, error) {
	return s.r.Save(data)
}

func (s *SessionService) CreateBulkSession(sessions *[]models.Session) ([]models.Session, error) {
	return s.r.BulkSave(sessions)
}

func (s *SessionService) UpdateSession(eventID int64, id int64, data *map[string]interface{}) (models.Session, error) {
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
		return models.Session{}, errors.New("record not found")
	}

	return s.r.Update(id, data)
}

func (s *SessionService) GetAllSessions(pageParams *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.Session, int64, error) {
	return s.r.FindAll(pageParams, filters, sorts)
}

func (s *SessionService) GetSessionByID(uid int64) (models.Session, error) {
	return s.r.FindByID(uid)
}

func (s *SessionService) Delete(uid int64) (models.Session, error) {
	return s.r.Delete(uid)
}
