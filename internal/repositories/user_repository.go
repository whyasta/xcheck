package repositories

import (
	"bigmind/xcheck-be/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	// GetByID(id int) (*models.User, error)
	FindAll(params map[string]interface{}) ([]models.User, error)
	Save(user *models.User) (models.User, error)
	FindByUsername(username string) (models.User, error)
	FindByID(uid int64) (models.User, error)

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

func (repo *userRepository) FindAll(params map[string]interface{}) ([]models.User, error) {
	var users []models.User

	err := repo.db.
		Preload("Role").
		Where(params).
		Find(&users).
		Error

	if err != nil {
		return nil, err
	}

	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}

func (repo *userRepository) Save(user *models.User) (models.User, error) {
	if err := repo.db.Create(user).Error; err != nil {
		return models.User{}, err
	}
	user.Password = ""
	return *user, nil
}

func (repo *userRepository) FindByUsername(username string) (user models.User, err error) {
	err = repo.db.
		Omit("Password").
		Preload("Role").
		Where("username = ?", username).
		First(&user).
		Error

	return user, err
}

func (repo *userRepository) FindByID(id int64) (user models.User, err error) {
	err = repo.db.Omit("Password").Preload("Role").First(&user, "id = ?", id).Error
	return
}

func (repo *userRepository) Signin(username, password string) (models.User, error) {
	user := models.User{}
	if err := repo.db.
		Preload("Role").
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}
