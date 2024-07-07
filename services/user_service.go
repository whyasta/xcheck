package services

import (
	"bigmind/xcheck-be/models"
	"bigmind/xcheck-be/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// var once sync.Once

type UserService struct {
	u repositories.UserRepository
	r repositories.RoleRepository
}

// var instance *UserService

// NewUserService: construction function, injected by user repository
func NewUserService(u repositories.UserRepository, r repositories.RoleRepository) *UserService {
	// once.Do(func() {
	// 	instance = &UserService{
	// 		r: r,
	// 	}
	// })
	// return instance
	return &UserService{
		u: u,
		r: r,
	}
}

func (s *UserService) GetAll() ([]models.User, error) {
	return s.u.GetAll()
}

func (s *UserService) CreateUser(user *models.User) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = ""
	user.PasswordHash = string(hash)
	return s.u.Create(user)
}

func (s *UserService) Signin(username string, password string) (models.User, error) {
	user, err := s.u.Signin(username, password)
	if err != nil {
		return models.User{}, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return models.User{}, errors.New("invalid username or password")
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *UserService) GetUserByUsername(uname string) (models.User, error) {
	return s.u.GetByUsername(uname)
}

func (s *UserService) GetUserByID(uid int64) (models.User, error) {
	return s.u.GetByID(uid)
}
