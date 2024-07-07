package repositories

import (
	"bigmind/xcheck-be/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	// GetByID(id int) (*models.User, error)
	GetAll() ([]models.User, error)
	Create(user *models.User) (models.User, error)
	GetByUsername(username string) (models.User, error)
	GetByID(uid int64) (models.User, error)
	// Update(user *models.User) error
	// Delete(id int) error

	Signin(username string, password string) (models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) GetAll() ([]models.User, error) {
	var users []models.User

	err := repo.db.Preload("Role").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *userRepository) Create(user *models.User) (models.User, error) {
	var err = repo.db.Omit("password").Create(user).Error
	if err != nil {
		return models.User{}, err
	}
	return *user, nil
}

func (repo *userRepository) GetByUsername(username string) (models.User, error) {
	user := models.User{
		Username: username,
	}
	err := repo.db.Preload("Role").First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (repo *userRepository) GetByID(id int64) (models.User, error) {
	user := models.User{
		ID: id,
	}
	err := repo.db.Preload("Role").First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (repo *userRepository) Signin(username string, password string) (models.User, error) {
	user := models.User{
		Username: username,
	}
	err := repo.db.Preload("Role").First(&user, "username = ?", username).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
