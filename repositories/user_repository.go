package repositories

import (
	"bigmind/xcheck-be/models"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type UserRepository interface {
	// GetByID(id int) (*models.User, error)
	GetAll() ([]*models.User, error)
	Create(user *models.User) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByID(uid int) (*models.User, error)
	// Update(user *models.User) error
	// Delete(id int) error

	Signin(username string, password string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) GetAll() ([]*models.User, error) {
	rows, err := repo.db.Query("SELECT id, username, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// func (repo *MySQLUserRepository) GetByID(id int) (*models.User, error) {
// 	row := repo.db.QueryRow("SELECT id, username, email FROM users WHERE id = ?", id)
// 	user := &models.User{}
// 	err := row.Scan(&user.ID, &user.Username, &user.Email)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

func (repo *userRepository) Create(user *models.User) (*models.User, error) {
	query := fmt.Sprintf("INSERT INTO users (username, password_hash, email, role) VALUES ('%s', '%s', '%s', '%s')", user.Username, user.PasswordHash, user.Email, user.Role)
	// fmt.Printf(query + "\n")

	insert, err := repo.db.Exec(query)
	if insert != nil {
		user_id, _ := insert.LastInsertId()
		content := &models.User{
			ID:       user_id,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		}
		return content, nil
	}

	return nil, err
}

func (repo *userRepository) Signin(username string, password string) (*models.User, error) {
	row := repo.db.QueryRow("SELECT id, username, email, password_hash, role FROM users WHERE username = ?", username)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) GetByUsername(username string) (*models.User, error) {
	row := repo.db.QueryRow("SELECT id, username, email, role FROM users WHERE username = ?", username)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) GetByID(id int) (*models.User, error) {
	row := repo.db.QueryRow("SELECT id, username, email, role FROM users WHERE id = ?", id)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}
