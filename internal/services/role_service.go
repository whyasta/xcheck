package services

import (
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
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
	return s.r.Save(role)
}

func (s *RoleService) GetAllRole(params map[string]interface{}) ([]models.UserRole, error) {
	return s.r.FindAll(params)
}

func (s *RoleService) GetPaginateAllRole(pageParams *utils.Paginate, params map[string]interface{}) ([]models.UserRole, int64, error) {
	result, count, err := s.r.Paginate(pageParams, params)
	return result, count, err
}

func (s *RoleService) GetRoleByID(uid int64) (models.UserRole, error) {
	return s.r.FindByID(uid)
}
