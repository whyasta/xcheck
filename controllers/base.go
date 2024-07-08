package controllers

import (
	"bigmind/xcheck-be/services"

	"gorm.io/gorm/utils"
)

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

func MakeQueryParams(params map[string][]string, queryParams []string) map[string]interface{} {
	newParams := make(map[string]interface{})
	for key, value := range params {
		if !utils.Contains(queryParams, key) {
			continue
		}
		newParams[key] = value[0]
	}
	return newParams
}
