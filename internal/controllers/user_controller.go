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
	"github.com/mitchellh/mapstructure"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// swagger:route GET /users User getUserList
// Get User list
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
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
		PagingInfo: utils.PagingInfo{
			Page:  pageParams.GetPage(count),
			Limit: pageParams.GetLimit(count),
			Total: int(count),
		},
	}

	c.JSON(http.StatusOK, utils.BuildResponseWithPaginate(http.StatusOK, constant.Success, "", allUsers, &meta))
}

// swagger:route POST /users User createUser
// Create User
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
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

// swagger:route GET /users/{id} User getUser
// Get User by id
//
// security:
//   - Bearer: []
//
// responses:
//
// 200:
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

func (r UserController) UpdateUser(c *gin.Context) {
	defer utils.ResponseHandler(c)

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	var user *models.UserUpdateRequest
	var request = make(map[string]interface{})

	c.Next()
	c.BindJSON(&request)
	mapstructure.Decode(request, &user)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(user)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(constant.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	result, err := r.service.UpdateUser(int64(uid), &request)
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}
	result.Password = ""
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", result))
}