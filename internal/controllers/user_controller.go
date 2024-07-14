package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"bigmind/xcheck-be/constant"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// GetAllUser retrieves all users based on the specified parameters.
//
// c *gin.Context: Gin context
// Return type(s): None
// @Summary      Get All users
// @Tags         users
// @ID			 user-get-all
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      404
// @Failure      500
// @Security	 BearerAuth
// @Router       /users [get]
func (u UserController) GetAllUser(c *gin.Context) {
	// params := MakeQueryParams(c.Request.URL.Query(), []string{"role_id"})
	pageParams, params := MakePaginationQueryParams(c.Request.URL.Query(), []string{"role_id"})
	// log.Println(pageParams)
	// allUsers, err := u.service.GetAllUser(params)
	allUsers, count, err := u.service.GetPaginateAllUser(pageParams, params)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	meta := utils.MetaResponse{
		Page:  pageParams.Page(count),
		Limit: pageParams.Limit(count),
		Total: int(count),
	}

	c.JSON(http.StatusOK, utils.BuildResponseWithPaginate(http.StatusOK, constant.Success, "", allUsers, &meta))
}

// CreateUser creates a new user.
//
// It takes a Gin context as a parameter.
// It binds the JSON request body to a User struct.
// It validates the User struct.
// It calls the CreateUser method of the UserService.
// It returns the created User as a JSON response.
// @Summary      Create user
// @Tags         users
// @ID			 user-create
// @Produce      json
// @Param		 user	body		models.UserRequest	true	"User"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      404
// @Failure      500
// @Security	 BearerAuth
// @Router       /users [post]
func (u UserController) CreateUser(c *gin.Context) {
	defer utils.ResponseHandler(c)

	var user *models.User

	c.Next()
	c.BindJSON(&user)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(user)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(constant.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	result, err := u.service.CreateUser(user)
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", result))
}

// GetUserByID retrieves a user by their ID and returns it as a JSON response.
//
// Parameters:
// - c: The gin Context for handling HTTP request and response.
// Returns: None
// @Summary      Get user by ID
// @ID			 user-get-by-id
// @Tags         users
// @Param id   path int true "User ID"
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      404
// @Failure      500
// @Security		BearerAuth
// @Router       /users/{id} [get]
func (u UserController) GetUserByID(c *gin.Context) {
	defer utils.ResponseHandler(c)
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	var user models.User
	user, err = u.service.GetUserByID(int64(uid))
	if err != nil {
		utils.PanicException(constant.DataNotFound, errors.New("data not found").Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", user))
}
