package repositories

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"

	"gorm.io/gorm"
)

type SessionRepository interface {
	Save(role *models.Session) (models.Session, error)
	Update(id int64, data *map[string]interface{}) (models.Session, error)
	Delete(uid int64) (models.Session, error)
	FindAll(paginate *utils.Paginate, filter []utils.Filter) ([]models.Session, int64, error)
	FindByID(uid int64) (models.Session, error)
}

type sessionRepository struct {
	base BaseRepository
}

func NewSessionRepository(db *gorm.DB) *sessionRepository {
	return &sessionRepository{
		base: NewBaseRepository(db, models.Session{}),
	}
}

func (repo *sessionRepository) Save(role *models.Session) (models.Session, error) {
	return BaseInsert(*repo.base.GetDB(), *role)
}

func (repo *sessionRepository) FindByID(id int64) (models.Session, error) {
	return BaseFindByID[models.Session](*repo.base.GetDB(), id, []string{})
}

func (repo *sessionRepository) FindAll(paginate *utils.Paginate, filters []utils.Filter) ([]models.Session, int64, error) {
	return BasePaginateWithFilter[[]models.Session](*repo.base.GetDB(), []string{}, paginate, filters)
}

func (repo *sessionRepository) Delete(id int64) (models.Session, error) {
	return BaseSoftDelete[models.Session](*repo.base.GetDB(), id)
}

func (repo *sessionRepository) Update(id int64, data *map[string]interface{}) (models.Session, error) {
	return BaseUpdate[models.Session](*repo.base.GetDB(), id, data)
}