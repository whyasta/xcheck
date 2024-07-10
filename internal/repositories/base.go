package repositories

import (
	"reflect"

	"gorm.io/gorm"
)

type Repository struct {
	User *userRepository
	Role *roleRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
		Role: NewRoleRepository(db),
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
