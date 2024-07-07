package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"bigmind/xcheck-be/constant"
	"bigmind/xcheck-be/models"
	"bigmind/xcheck-be/services"
	"bigmind/xcheck-be/token"
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
func (u UserController) GetAll(c *gin.Context) {
	allUsers, err := u.service.GetAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// fmt.Println("MySQL All Users:", allUsers)

	// c.JSON(http.StatusOK, gin.H{"message": "pong"})
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", allUsers))
}

func (u UserController) Create(c *gin.Context) {
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

	user, err = u.service.Create(user)
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", user))
}

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

	var user *models.User
	user, err = u.service.Signin(userLogin.Username, userLogin.Password)
	if err != nil {
		utils.PanicException(constant.Unauthorized, err.Error())
		return
	}

	tokenStr, err := token.GenerateToken(user.Username)
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseWithToken(http.StatusOK, constant.Success, tokenStr, "", utils.Null()))
}

func (u UserController) Signout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

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
func (u UserController) GetByID(c *gin.Context) {
	defer utils.ResponseHandler(c)
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	var user *models.User
	user, err = u.service.GetByID(uid)
	if err != nil {
		utils.PanicException(constant.DataNotFound, errors.New("user not found").Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", user))
}

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
	username, err := token.ExtractTokenID(c)

	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	var user *models.User
	user, err = u.service.GetByUsername(username)

	if err != nil {
		utils.PanicException(constant.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, constant.Success, "", user))
}
