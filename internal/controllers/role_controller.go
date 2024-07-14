package controllers

import (
	"bigmind/xcheck-be/constant"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/services"
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

// CreateRole creates a new role based on the input User Role and returns the result.
//
// Parameters:
// - c: The gin Context for handling HTTP request and response.
// Returns: None
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

// GetAllRole retrieves all roles based on the specified parameters.
//
// Parameter(s):
//
//	c *gin.Context: Gin context
//
// Return type(s): None
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
	// params := MakeQueryParams(c.Request.URL.Query(), []string{"role_name"})
	// rows, err := r.service.GetAllRole(params)

	pageParams, params := MakePaginationQueryParams(c.Request.URL.Query(), []string{"role_id"})
	rows, count, err := r.service.GetPaginateAllRole(pageParams, params)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	meta := utils.MetaResponse{
		Page:  pageParams.GetPage(count),
		Limit: pageParams.GetLimit(count),
		Total: int(count),
	}

	c.JSON(http.StatusOK, utils.BuildResponseWithPaginate(http.StatusOK, constant.Success, "", rows, &meta))
}

// GetRoleByID retrieves a role by its ID from the database and returns it as a JSON response.
//
// Parameters:
// - c: The gin Context for handling HTTP request and response.
//
// Return: None.
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
