package services

import (
	"bigmind/xcheck-be/models"
	"bigmind/xcheck-be/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// var once sync.Once

type UserService struct {
	repo repositories.UserRepository
}

// var instance *UserService

// NewUserService: construction function, injected by user repository
func NewUserService(r repositories.UserRepository) *UserService {
	// once.Do(func() {
	// 	instance = &UserService{
	// 		repo: r,
	// 	}
	// })
	// return instance
	return &UserService{
		repo: r,
	}
}

func (s *UserService) GetAll() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) Create(user *models.User) (models.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.PasswordHash = string(hash)
	return s.repo.Create(user)
}

func (s *UserService) Signin(username string, password string) (models.User, error) {
	user, err := s.repo.Signin(username, password)
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

func (s *UserService) GetByUsername(uname string) (models.User, error) {
	return s.repo.GetByUsername(uname)
}

func (s *UserService) GetByID(uid int64) (models.User, error) {
	return s.repo.GetByID(uid)
}
