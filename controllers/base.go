package controllers

import "bigmind/xcheck-be/services"

type Controller struct {
	UserController *UserController
	RoleController *RoleController
}

func NewController(
	services *services.Service,
) *Controller {
	return &Controller{
		UserController: NewUserController(services.UserService),
		RoleController: NewRoleController(services.RoleService),
	}
}
