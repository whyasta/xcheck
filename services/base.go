package services

import "bigmind/xcheck-be/repositories"

type Service struct {
	UserService *UserService
	RoleService *RoleService
}

func NewService(
	repositories *repositories.Repository,
) *Service {
	return &Service{
		UserService: NewUserService(repositories.User),
		RoleService: NewRoleService(repositories.Role),
	}
}
