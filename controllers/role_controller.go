package controllers

import (
	"bigmind/xcheck-be/constant"
	"bigmind/xcheck-be/models"
	"bigmind/xcheck-be/services"
	"bigmind/xcheck-be/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RoleController struct {
	service *services.RoleService
}

func NewRoleController(service *services.RoleService) *RoleController {
	return &RoleController{
		service: service,
	}
}

// @Summary      Create role
// @Tags         roles
// @ID			 role-create
// @Produce      json
// @Param		 role	body		models.UserRoleRequest	true	"User Role"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      404
// @Failure      500
// @Security	 BearerAuth
// @Router       /roles [post]
func (r RoleController) CreateRole(c *gin.Context) {
	defer utils.ResponseHandler(c)

	var role *models.UserRole

	c.Next()
	c.BindJSON(&role)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(role)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(constant.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	result, err := r.service.CreateRole(role)
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", result))
}

// @Summary      Get All roles
// @Tags         roles
// @ID			 role-get-all
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      404
// @Failure      500
// @Security	 BearerAuth
// @Router       /roles [get]
func (r RoleController) GetAllRole(c *gin.Context) {
	rows, err := r.service.GetAllRole()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", rows))
}

// @Summary      Get by ID
// @ID			 role-get-by-id
// @Tags         roles
// @Param id     path int true "Role ID"
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      404
// @Failure      500
// @Security		BearerAuth
// @Router       /roles/{id} [get]
func (r RoleController) GetRoleByID(c *gin.Context) {
	defer utils.ResponseHandler(c)
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	var user models.UserRole
	user, err = r.service.GetRoleByID(int64(uid))
	if err != nil {
		utils.PanicException(constant.DataNotFound, errors.New("data not found").Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", user))
}
