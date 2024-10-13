package controllers

import (
	"bigmind/xcheck-be/internal/constant/response"
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/services"
	"bigmind/xcheck-be/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type RedeemController struct {
	service *services.RedeemService
}

func NewRedeemController(service *services.RedeemService) *RedeemController {
	return &RedeemController{
		service: service,
	}
}

func (r RedeemController) GetAll(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	pageParams, filter, sort := MakePageFilterQueryParams(c.Request.URL.Query(), []string{"event_id"})
	filter = append(filter, utils.Filter{
		Property:  "event_id",
		Operation: "=",
		Value:     strconv.Itoa(eventID),
	})

	rows, count, err := r.service.GetFilteredRedeems(pageParams, filter, sort)

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

	c.JSON(http.StatusOK, utils.BuildResponseWithPaginate(http.StatusOK, response.Success, "", rows, &meta))
}

func (r RedeemController) Redeem(c *gin.Context) {
	defer utils.ResponseHandler(c)

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	var request *dto.RedeemRequest

	c.Next()
	c.BindJSON(&request)

	validate := validator.New(validator.WithRequiredStructEnabled())
	en := en.New()
	UniversalTranslator = ut.New(en, en)
	trans, _ := UniversalTranslator.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans)

	err = validate.Struct(request)
	if err != nil {
		fmt.Println(err)
		// errors := err.(validator.ValidationErrors)

		// errors := utils.TranslateError(err, trans)
		validatorErrs := err.(validator.ValidationErrors)
		var errors []error
		for _, e := range validatorErrs {
			translatedErr := fmt.Errorf(e.Translate(trans))
			errors = append(errors, translatedErr)
		}
		// fmt.Println(errsEn)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	data, err := r.service.Redeem(int64(eventID), request.OrderID)

	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", data))
}

func (r RedeemController) Check(c *gin.Context) {
	defer utils.ResponseHandler(c)

	// eventID, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	utils.PanicException(response.InvalidRequest, err.Error())
	// 	return
	// }

}

/*
func (r RedeemController) Import(c *gin.Context) {
	defer utils.ResponseHandler(c)

	var role *models.UserRole

	c.Next()
	c.BindJSON(&role)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(role)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.PanicException(response.InvalidRequest, fmt.Sprintf("Validation error: %s", errors))
		return
	}

	result, err := r.service.CreateRole(role)
	if err != nil {
		utils.PanicException(response.InvalidRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.BuildResponse(http.StatusOK, response.Success, "", result))
}
*/
