package repositories

import (
	"bigmind/xcheck-be/models"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type RoleRepository interface {
	Create(user *models.UserRole) (*models.UserRole, error)
	GetAll() ([]*models.UserRole, error)
	GetByID(uid int) (*models.UserRole, error)
}

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *roleRepository {
	return &roleRepository{
		db: db,
	}
}

func (repo *userRepository) CreateRole(role *models.UserRole) (*models.UserRole, error) {
	return nil, nil
}
