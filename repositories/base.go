package repositories

import "database/sql"

type Repository struct {
	User *userRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
