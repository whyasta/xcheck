package services

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	u repositories.UserRepository
	b repositories.BaseRepository
}

// var roleInstance *RoleService

// NewUserService constructs a new UserService using the provided UserRepository.
//
// Parameter:
//
//	u - UserRepository: the user repository to be injected into the UserService.
//
// Return type:
//
//	*UserService: the newly created UserService instance.
func NewUserService(u repositories.UserRepository, b repositories.BaseRepository) *UserService {
	// sync.Once.Do(func() {
	// 	instance = &UserService{
	// 		u: u,
	// 		r: r,
	// 	}
	// })
	// return instance
	return &UserService{u, b}
}

func (s *UserService) GetAllUser(params map[string]interface{}) ([]models.User, error) {
	result, err := s.u.FindAll(params)
	return result, err
}

func (s *UserService) GetPaginateAllUser(pageParams *utils.Paginate, params map[string]interface{}) ([]models.User, int64, error) {
	result, count, err := s.u.Paginate(pageParams, params)
	return result, count, err
}

func (s *UserService) CreateUser(user *models.User) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(hash)
	return s.u.Save(user)
}

func (s *UserService) GetUserByID(uid int64) (models.User, error) {
	return s.u.FindByID(uid)
}

// func (s *UserService) GetUserByUsername(uname string) (models.User, error) {
// 	return s.u.FindByUsername(uname)
// }

func (s *UserService) UpdateUser(id int64, data *map[string]interface{}) (models.User, error) {
	return s.u.Update(id, data)
}
