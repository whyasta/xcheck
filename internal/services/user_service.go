package services

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	u repositories.UserRepository
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
func NewUserService(u repositories.UserRepository) *UserService {
	// sync.Once.Do(func() {
	// 	instance = &UserService{
	// 		u: u,
	// 		r: r,
	// 	}
	// })
	// return instance
	return &UserService{
		u: u,
	}
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

func (s *UserService) GetUserByUsername(uname string) (models.User, error) {
	return s.u.FindByUsername(uname)
}

func (s *UserService) Signin(username string, password string) (models.User, error) {
	user, err := s.u.Signin(username, password)
	if err != nil {
		return models.User{}, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, errors.New("invalid username or password")
	}
	user.Password = ""
	return user, nil
}
