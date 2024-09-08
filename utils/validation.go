package utils

import (
	"log"
	"regexp"
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
