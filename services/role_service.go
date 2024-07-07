package services

import (
	"bigmind/xcheck-be/repositories"
)

// var once sync.Once

type RoleService struct {
	repo repositories.RoleRepository
}
