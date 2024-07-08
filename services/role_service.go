package services

import (
	"bigmind/xcheck-be/models"
	"bigmind/xcheck-be/repositories"
)

type RoleService struct {
	r repositories.RoleRepository
}

func NewRoleService(r repositories.RoleRepository) *RoleService {
	// once.Do(func() {
	// 	instance = &RoleService{
	// 		r: r,
	// 	}
	// })
	// return instance
	return &RoleService{
		r: r,
	}
}

func (s *RoleService) CreateRole(role *models.UserRole) (models.UserRole, error) {
	return s.r.SaveRole(role)
}

func (s *RoleService) GetAllRole(params map[string]interface{}) ([]models.UserRole, error) {
	return s.r.FindAllRoles(params)
}

func (s *RoleService) GetRoleByID(uid int64) (models.UserRole, error) {
	return s.r.FindRoleByID(uid)
}
