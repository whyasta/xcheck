package repositories

import (
	"bigmind/xcheck-be/models"

	"gorm.io/gorm"
)

type RoleRepository interface {
	CreateRole(role *models.UserRole) (models.UserRole, error)
	GetAllRole() ([]models.UserRole, error)
	// GetRoleById(uid int) (*models.UserRole, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{
		db: db,
	}
}

func (repo *roleRepository) CreateRole(role *models.UserRole) (models.UserRole, error) {
	var err = repo.db.Create(role).Error
	if err != nil {
		return models.UserRole{}, err
	}
	return *role, nil
}

func (repo *roleRepository) GetAllRole() ([]models.UserRole, error) {
	var roles []models.UserRole

	err := repo.db.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
