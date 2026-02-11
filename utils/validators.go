package utils

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidateFutureDate(fl validator.FieldLevel) bool {
	dateTime, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return dateTime.After(time.Now())
}
