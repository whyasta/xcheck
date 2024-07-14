package services

import "bigmind/xcheck-be/internal/repositories"

type Service struct {
	AuthService *AuthService
	UserService *UserService
	RoleService *RoleService
}

func NewService(
	repositories *repositories.Repository,
) *Service {
	return &Service{
		AuthService: NewAuthService(repositories.User),
		UserService: NewUserService(repositories.User),
		RoleService: NewRoleService(repositories.Role),
	}
}
