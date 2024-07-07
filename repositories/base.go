package repositories

import (
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
