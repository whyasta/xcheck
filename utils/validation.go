package utils

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func DateValidation(fl validator.FieldLevel) bool {
	// log.Println(fl.Field().String())
	if fl.Field().String() == "" {
		return true
	}
	_, err := time.Parse("2006-01-02", fl.Field().String())
	return err == nil
}

func BarcodeValidation(fl validator.FieldLevel) bool {
	log.Println("BarcodeValidation", fl.Field().String())
	if fl.Field().String() == "" {
		return false
	}
	var alphanumeric = regexp.MustCompile("^[a-zA-Z0-9_-]*$")
	return alphanumeric.MatchString(fl.Field().String())
}

func InSessionValidator(fl validator.FieldLevel) bool {
	param := fl.Field().String()
	validValues := []string{"admin", "user", "guest"} // Define allowed values

	for _, value := range validValues {
		if param == value {
			return true
		}
	}
	return false
}

func InitValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
	return validate
}

// func FormatValidationError(err error) map[string]string {
func FormatValidationError(err error, obj interface{}) string {
	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		// jsonKey := GetJSONTag(obj, err.StructField())
		// if jsonKey == "" {
		// 	jsonKey = err.StructField()
		// }
		var message string
		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", err.Field())
		case "in":
			message = fmt.Sprintf("%s must be one of the valid values", err.Field())
		case "min":
			message = fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
		case "max":
			message = fmt.Sprintf("%s must be at most %s", err.Field(), err.Param())
		default:
			message = fmt.Sprintf("%s is not valid", err.Field())
		}
		errors[err.Field()] = message
	}

	var pairs []string

	for _, value := range errors {
		// pair := fmt.Sprintf("%s: %s", key, value)
		pair := fmt.Sprintf("%s", value)
		pairs = append(pairs, pair)
	}

	return strings.Join(pairs, ", ")
}

func GetJSONTag(obj interface{}, fieldName string) string {
	val := reflect.TypeOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	field, found := val.FieldByName(fieldName)
	if !found {
		return ""
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return ""
	}

	return strings.Split(jsonTag, ",")[0] // Get only the key, ignoring options like "omitempty"
}
