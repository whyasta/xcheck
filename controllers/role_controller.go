package controllers

import (
	"bigmind/xcheck-be/services"
)

type RoleController struct {
	service *services.UserService
}

func NewRoleController(service *services.UserService) *UserController {
	return &UserController{
		service: service,
	}
}
