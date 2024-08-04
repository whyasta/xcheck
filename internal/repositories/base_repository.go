package repositories

import (
	"bigmind/xcheck-be/utils"
	"log"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	User       *userRepository
	Role       *roleRepository
	Event      *eventRepository
	TicketType *ticketTypeRepository
	Gate       *gateRepository
	Session    *sessionRepository
	Base       *baseRepository
	Import     *importRepository
	Barcode    *barcodeRepository
	Schedule   *scheduleRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:       NewUserRepository(db),
		Role:       NewRoleRepository(db),
		Event:      NewEventRepository(db),
		TicketType: NewTicketTypeRepository(db),
		Gate:       NewGateRepository(db),
		Session:    NewSessionRepository(db),
		Schedule:   NewScheduleRepository(db),
		Import:     NewImportRepository(db),
		Barcode:    NewBarcodeRepository(db),
	}
}

// type BaseRepository[M BaseModel[E], E any] struct {
// 	db *gorm.DB
// }

// type BaseModel[E any] interface {
// 	ToEntity() E
// 	FromEntity(entity E) interface{}
// }

// func NewBaseRepository[M BaseModel[E], E any](db *gorm.DB) *BaseRepository[M, E] {
// 	return &BaseRepository[M, E]{
// 		db: db,
// 	}
// }

// func (r *BaseRepository[M, E]) Insert(ctx context.Context, entity *E) error {
// 	var start M
// 	model := start.FromEntity(*entity).(M)

// 	err := r.db.WithContext(ctx).Create(&model).Error
// 	if err != nil {
// 		return err
// 	}

// 	*entity = model.ToEntity()
// 	return nil
// }

type baseRepository struct {
	db    *gorm.DB
	model interface{}
}

type BaseRepository interface {
	GetDB() *gorm.DB
	BeginTx()
	CommitTx()
	RollbackTx()
	CommonInsert(table string, item interface{}) (interface{}, error)
	CommonUpdate(table string, params map[string]interface{}, data interface{}) (interface{}, error)
	CommonFindByID(table string, uid int64) (interface{}, error)
}

func NewBaseRepository(db *gorm.DB, model interface{}) BaseRepository {
	return &baseRepository{db, model}
}

func (br *baseRepository) GetDB() *gorm.DB {
	return br.db
}

func (br *baseRepository) BeginTx() {
	br.db = br.GetDB().Begin()
}

func (br *baseRepository) CommitTx() {
	br.GetDB().Commit()
}

func (br *baseRepository) RollbackTx() {
	br.GetDB().Rollback()
}

func (br *baseRepository) CommonInsert(table string, item interface{}) (interface{}, error) {
	var err = br.GetDB().Table(table).Create(&item).Error
	return item, err
}

func (br *baseRepository) CommonFindByID(table string, id int64) (interface{}, error) {
	result := make(map[string]interface{})
	err := br.GetDB().Model(br.model).Table(table).First(&result, "id = ?", id).Error
	return result, err
}

func (br *baseRepository) CommonUpdate(table string, params map[string]interface{}, data interface{}) (interface{}, error) {
	result := make(map[string]interface{})

	// fmt.Println(br.model)
	// var err = br.GetDB().
	// 	Table(table).
	// 	Clauses(clause.Returning{}).
	// 	Where(params).
	// 	Updates(data).
	// 	First(&result).
	// 	Error
	err := br.GetDB().Model(br.model).Table(table).First(&result, "id = ?", 1).Error
	return result, err

	// var err = br.GetDB().Table(table).Create(&item).Error
	// return item, err
}

func BaseInsert[M any](db gorm.DB, item M) (M, error) {
	var result = reflect.ValueOf(item).Interface().(M)
	var err = db.Create(&result).Error
	return result, err
}

func BaseFindByID[M any](db gorm.DB, tableName string, id int64, joins []string) (M, error) {
	var result M
	tx := db.Model(&result)
	if len(joins) > 0 {
		for _, join := range joins {
			tx = tx.Preload(join)
		}
	}

	err := tx.First(&result, "id = ?", id).Error
	// err := tx.Error
	return result, err
}

func BaseSoftDelete[M any](db gorm.DB, id int64) (M, error) {
	var result M
	err := db.Where("id = ?", id).Delete(result).Error
	return result, err
}

func BaseUpdate[M any](db gorm.DB, id int64, updated *map[string]interface{}) (M, error) {
	var result M

	var err = db.Model(&result).
		// Table("events").
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updated).
		First(&result).
		Error
	return result, err
}

func BasePaginate[M any](db gorm.DB, paginate *utils.Paginate, params map[string]interface{}) (M, int64, error) {
	var records M
	var count int64

	tx := db.
		Scopes(paginate.PaginatedResult).
		Where(params).
		Find(&records)

	tx.Limit(-1).Offset(-1)
	tx.Count(&count)

	err := tx.Error

	return records, count, err
}

func BasePaginateWithFilter[M any](db gorm.DB, joins []string, paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) (M, int64, error) {
	var records M
	var count int64

	tx := db.
		Scopes(paginate.PaginatedResult)

	for _, join := range joins {
		tx = tx.Preload(join)
	}

	if len(filters) > 0 {
		for _, filter := range filters {
			newFilter := utils.NewFilter(filter.Property, filter.Operation, filter.Collation, filter.Value, filter.Items)
			tx = newFilter.FilterResult("", tx)
		}
	}

	log.Println("sorts: ", sorts)
	if len(sorts) > 0 {
		for _, sort := range sorts {
			newSort := utils.NewSort(sort.Property, sort.Direction)
			tx = newSort.SortResult(tx)
		}
	}

	tx = tx.Find(&records)

	if len(filters) <= 0 {
		tx.Limit(-1).Offset(-1)
	}
	tx.Count(&count)

	err := tx.Error

	return records, count, err
}
