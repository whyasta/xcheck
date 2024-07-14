package controllers

import (
	"bigmind/xcheck-be/constant"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{
		service: service,
	}
}

// swagger:route POST /auth/signin Auth authSignin
// Signin
//
// responses:
//
// 200:
func (u AuthController) Signin(c *gin.Context) {
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

	_, tokenPair, err := u.service.Signin(userLogin.Username, userLogin.Password)
	if err != nil {
		utils.PanicException(constant.Unauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseWithToken(http.StatusOK, constant.Success, tokenPair["token"], tokenPair["refresh_token"], "", utils.Null()))
}

// swagger:route POST /auth/token Auth authRefreshToken
// Refresh token
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (u AuthController) Refresh(c *gin.Context) {
	defer utils.ResponseHandler(c)
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	tokenReq := tokenReqBody{}
	c.Bind(&tokenReq)

	tokenPair, err := u.service.RefreshToken(tokenReq.RefreshToken)
	if err != nil {
		utils.PanicException(constant.Unauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseWithToken(http.StatusOK, constant.Success, tokenPair["token"], tokenPair["refresh_token"], "", utils.Null()))
}

// swagger:route POST /auth/signout Auth signout
// Signout
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (u AuthController) Signout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// swagger:route GET /auth/me Auth authMe
// Get current user
// security:
//   - Bearer: []
//
// responses:
//
// 200:
func (u AuthController) CurrentUser(c *gin.Context) {
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
