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
	params := MakeQueryParams(c.Request.URL.Query(), []string{"role_id"})
	allUsers, err := u.service.GetAllUser(params)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", allUsers))
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

// Signin signs in a user based on the provided credentials.
// It takes a Gin context as a parameter.
// @Summary		Signin
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			account	body		models.UserLogin	true	"User Login"
// @Success		200
// @Failure		400
// @Failure		401
// @Failure		404
// @Failure		500
// @Security		BearerAuth
// @Router			/auth/signin [post]
func (u UserController) Signin(c *gin.Context) {
	defer utils.ResponseHandler(c)

	var userLogin *models.UserLogin

	c.Next()
	c.BindJSON(&userLogin)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(userLogin)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(constant.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	var user models.User
	user, err = u.service.Signin(userLogin.Username, userLogin.Password)
	if err != nil {
		utils.PanicException(constant.Unauthorized, err.Error())
		return
	}

	tokenStr, err := utils.GenerateToken(user.Username)
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseWithToken(http.StatusOK, constant.Success, tokenStr, "", utils.Null()))
}

func (u UserController) Signout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
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

// CurrentUser retrieves the current user based on the token provided in the context.
//
// Parameters:
// c *gin.Context: Gin context containing the token information.
// Returns:
// None
// @Summary      Get current user
// @ID			 user-current
// @Tags         auth
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      404
// @Failure      500
// @Security		BearerAuth
// @Router       /auth/me [get]
func (u UserController) CurrentUser(c *gin.Context) {
	defer utils.ResponseHandler(c)
	username, err := utils.ExtractTokenID(c)

	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	var user models.User
	user, err = u.service.GetUserByUsername(username)

	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", user))
}
