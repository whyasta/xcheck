package repositories

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	// GetByID(id int) (*models.User, error)
	Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.User, int64, error)
	//FindAll(params map[string]interface{}) ([]models.User, error)
	FindAll(paginate *utils.Paginate, filter []utils.Filter, sorts []utils.Sort) ([]models.User, int64, error)

	Save(user *models.User) (models.User, error)
	FindByUsername(username string) (models.User, error)
	FindByID(uid int64) (models.User, error)
	Update(id int64, event *map[string]interface{}) (models.User, error)

	Signin(username string, password string) (models.User, error)
	CreateAuth(id int64) (utils.AuthDetails, error)
	FindByAuth(uid int64, authUuid string) (models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) Paginate(paginate *utils.Paginate, params map[string]interface{}) ([]models.User, int64, error) {
	var users []models.User
	var count int64

	tx := repo.db.
		Scopes(paginate.PaginatedResult).
		Omit("Password", "AuthUuids").
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

func (repo *userRepository) FindAll(paginate *utils.Paginate, filters []utils.Filter, sorts []utils.Sort) ([]models.User, int64, error) {
	var records []models.User
	var count int64

	tx := repo.db.
		Scopes(paginate.PaginatedResult)

	tx = tx.Preload("Role").Omit("Password", "AuthUuids")

	if len(filters) > 0 {
		for _, filter := range filters {
			newFilter := utils.NewFilter(filter.Property, filter.Operation, filter.Collation, filter.Value, filter.Items)
			tx = newFilter.FilterResult("", tx)
		}
	}

	// log.Println("sorts: ", sorts)
	if len(sorts) > 0 {
		for _, sort := range sorts {
			newSort := utils.NewSort(sort.Property, sort.Direction)
			tx = newSort.SortResult(tx)
		}
	}

	tx = tx.Find(&records)

	if len(filters) <= 0 {
		tx.Limit(-1).Offset(-1)
	}
	tx.Count(&count)

	err := tx.Error

	return records, count, err
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
	err = repo.db.Omit("Password").Omit("AuthUuids").Preload("Role").First(&user, "id = ?", id).Error
	return
}

func (repo *userRepository) FindByAuth(id int64, authId string) (user models.User, err error) {
	err = repo.db.Omit("Password").
		Omit("AuthUuids").Preload("Role").
		First(&user, "id = ? AND auth_uuids LIKE ?", id, "%"+authId+"%").
		Error
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

func (repo *userRepository) Update(id int64, data *map[string]interface{}) (models.User, error) {
	return BaseUpdate[models.User](*repo.db, id, data)
}

func (repo *userRepository) CreateAuth(id int64) (utils.AuthDetails, error) {
	var user models.User
	err := repo.db.Omit("Password").Preload("Role").First(&user, "id = ?", id).Error
	if err != nil {
		return utils.AuthDetails{}, err
	}

	authUuid := uuid.New().String()

	/*
			// uuids := user.AuthUuids
			var uuids []interface{}
			if user.AuthUuids != "" {
				if err := json.Unmarshal([]byte(user.AuthUuids), &uuids); err != nil {
					fmt.Println(err)
					return utils.AuthDetails{}, err
				}
			}

		    uuids = append(uuids, authUuid)
			jsonStr, _ := json.Marshal(uuids)
			user.AuthUuids = string(jsonStr)

			if err := repo.db.Table("users").
				Where("id = ?", id).
				Updates(map[string]interface{}{"auth_uuids": jsonStr}).
				Error; err != nil {
				return utils.AuthDetails{}, err
			}*/

	var authD utils.AuthDetails
	authD.AuthUuid = authUuid
	authD.UserId = uint64(user.ID)
	return authD, nil
}
