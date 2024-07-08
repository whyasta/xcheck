package repositories

import (
	"bigmind/xcheck-be/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	// GetByID(id int) (*models.User, error)
	FindAllUser() ([]models.User, error)
	SaveUser(user *models.User) (models.User, error)
	FindUserByUsername(username string) (models.User, error)
	FindUserByID(uid int64) (models.User, error)

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

func (repo *userRepository) FindAllUser() ([]models.User, error) {
	var users []models.User

	err := repo.db.Omit("Password").Preload("Role").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *userRepository) SaveUser(user *models.User) (models.User, error) {
	var err = repo.db.Create(user).Error
	if err != nil {
		return models.User{}, err
	}
	user.Password = ""
	return *user, nil
}

func (repo *userRepository) FindUserByUsername(username string) (models.User, error) {
	user := models.User{
		Username: username,
	}
	err := repo.db.Omit("Password").Preload("Role").First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	user.Password = ""
	return user, nil
}

func (repo *userRepository) FindUserByID(id int64) (models.User, error) {
	user := models.User{
		ID: id,
	}
	err := repo.db.Omit("Password").Preload("Role").First(&user).Error
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
