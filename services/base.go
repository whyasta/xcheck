package services

import "bigmind/xcheck-be/repositories"

type Service struct {
	UserService *UserService
}

func NewService(
	repositories *repositories.Repository,
) *Service {
	return &Service{
		UserService: NewUserService(repositories.User),
	}
}
