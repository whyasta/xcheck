package repositories

import (
	"bigmind/xcheck-be/models"

	"gorm.io/gorm"
)

type RoleRepository interface {
	SaveRole(role *models.UserRole) (models.UserRole, error)
	FindAllRole() ([]models.UserRole, error)
	FindRoleByID(uid int64) (models.UserRole, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{
		db: db,
	}
}

func (repo *roleRepository) SaveRole(role *models.UserRole) (models.UserRole, error) {
	var err = repo.db.Create(role).Error
	if err != nil {
		return models.UserRole{}, err
	}
	return *role, nil
}

func (repo *roleRepository) FindAllRole() ([]models.UserRole, error) {
	var roles []models.UserRole

	err := repo.db.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (repo *roleRepository) FindRoleByID(id int64) (models.UserRole, error) {
	role := models.UserRole{
		ID: id,
	}
	err := repo.db.First(&role).Error
	if err != nil {
		return models.UserRole{}, err
	}
	return role, nil
}
