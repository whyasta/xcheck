package repositories

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"
	"log"

	"gorm.io/gorm"
)

type UserRepository interface {
	// GetByID(id int) (*models.User, error)
	Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.User, int64, error)

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
	// log.Println(params)
	err := repo.db.
		// Scopes(NewPaginate(params["limit"], params["page"]).PaginatedResult).
		// Preload("Role").
		Joins("Role").
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

func (repo *userRepository) Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.User, int64, error) {
	var users []models.User
	var count int64

	log.Println(paginate)

	tx := repo.db.
		Scopes(paginate.PaginatedResult).
		Preload("Role").
		Where(params).
		Find(&users)

	err := tx.Error

	if err != nil {
		return nil, 0, err
	}

	for i := range users {
		users[i].Password = ""
	}

	tx.Limit(-1).Offset(-1)
	tx.Count(&count)

	return users, count, nil
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
		// Omit("Password").
		// Preload("Role").
		Joins("Role").
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
		// Preload("Role").
		Where("username = ?", username).
		Joins("Role").
		First(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}
