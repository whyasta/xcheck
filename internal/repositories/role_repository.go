package repositories

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"
	"fmt"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Save(role *models.UserRole) (models.UserRole, error)
	FindAll(params map[string]interface{}) ([]models.UserRole, error)
	FindByID(uid int64) (models.UserRole, error)
	// BaseFindByID(id int64) (models.UserRole, error)

	Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.UserRole, int64, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{
		db: db,
	}
}

func (repo *roleRepository) Save(role *models.UserRole) (models.UserRole, error) {
	// var err = repo.db.Create(role).Error
	// if err != nil {
	// 	return models.UserRole{}, err
	// }
	// return *role, nil
	return BaseInsert(*repo.db, *role)
}

func (repo *roleRepository) FindAll(params map[string]interface{}) ([]models.UserRole, error) {
	var roles []models.UserRole
	fmt.Println(params)
	err := repo.db.Where(params).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (repo *roleRepository) FindByID(id int64) (models.UserRole, error) {
	// role := models.UserRole{
	// 	ID: id,
	// }
	// err := repo.db.First(&role).Error
	// if err != nil {
	// 	return models.UserRole{}, err
	// }
	// return role, nil
	return BaseFindByID[models.UserRole](*repo.db, id, []string{})
}

func (repo *roleRepository) Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.UserRole, int64, error) {
	var roles []models.UserRole
	var count int64

	// fmt.Println(paginate)

	tx := repo.db.
		Scopes(paginate.PaginatedResult).
		Where(params).
		Find(&roles)

	err := tx.Error

	if err != nil {
		return nil, 0, err
	}

	tx.Limit(-1).Offset(-1)
	tx.Count(&count)

	return roles, count, nil
}
