package controllers

import (
	"bigmind/xcheck-be/internal/constant/response"
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

func (u AuthController) Signin(c *gin.Context) {
	defer utils.ResponseHandler(c)

	var userLogin *models.UserLogin

	c.Next()
	c.BindJSON(&userLogin)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(userLogin)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	data, tokenPair, err := u.service.Signin(userLogin.Username, userLogin.Password)
	if err != nil {
		utils.PanicException(response.Unauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseWithToken(http.StatusOK, response.Success, tokenPair["token"], tokenPair["refresh_token"], "", data))
}

func (u AuthController) Refresh(c *gin.Context) {
	defer utils.ResponseHandler(c)
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	tokenReq := tokenReqBody{}
	c.Bind(&tokenReq)

	tokenPair, err := u.service.RefreshToken(tokenReq.RefreshToken)
	if err != nil {
		utils.PanicException(response.Unauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponseWithToken(http.StatusOK, response.Success, tokenPair["token"], tokenPair["refresh_token"], "", utils.Null()))
}

func (u AuthController) Signout(c *gin.Context) {
	defer utils.ResponseHandler(c)

	err := utils.BlacklistToken(utils.ExtractToken(c))
	if err != nil {
		utils.PanicException(response.Unauthorized, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok", "status": "SUCCESS"})
}

func (u AuthController) CurrentUser(c *gin.Context) {
	defer utils.ResponseHandler(c)
	uid, authID, err := utils.ExtractTokenID(c)

	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var user models.User
	user, err = u.service.GetUserByAuth(uid, authID)

	if err != nil {
		utils.PanicException(response.Unauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", user))
}

func (u AuthController) CheckAuthID(c *gin.Context) bool {
	defer utils.ResponseHandler(c)
	uid, authID, err := utils.ExtractTokenID(c)

	if err != nil {
		return false
	}

	_, err = u.service.GetUserByAuth(uid, authID)
	return err == nil
}

func (u AuthController) ResetPassword(c *gin.Context) {
	defer utils.ResponseHandler(c)

	var userLogin *models.UserLogin

	c.Next()
	c.BindJSON(&userLogin)

	validate := utils.InitValidator()
	err := validate.Struct(userLogin)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	err = u.service.ResetPassword(userLogin.Username, userLogin.Password)
	if err != nil {
		utils.PanicException(response.Unauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", utils.Null()))
}
