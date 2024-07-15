package repositories

import (
	"reflect"

	"gorm.io/gorm"
)

type Repository struct {
	User  *userRepository
	Role  *roleRepository
	Event *eventRepository
	Base  *baseRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:  NewUserRepository(db),
		Role:  NewRoleRepository(db),
		Event: NewEventRepository(db),
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

func BaseInsert[M any](db gorm.DB, item M) (M, error) {
	var result = reflect.ValueOf(item).Interface().(M)
	var err = db.Create(&result).Error
	return result, err
}

func BaseFindByID[M any](db gorm.DB, id int64) (M, error) {
	var result M
	err := db.First(&result, "id = ?", id).Error
	return result, err
}

func BaseSoftDelete[M any](db gorm.DB, id int64) (M, error) {
	var result M
	err := db.Where("id = ?", id).Delete(result).Error
	return result, err
}
