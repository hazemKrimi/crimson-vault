package api

import (
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (validator *CustomValidator) Validate(i any) error {
	if err := validator.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func PasswordValidator(fieldLevel validator.FieldLevel) bool {
	password := fieldLevel.Field().String()

	var (
		upper    = regexp.MustCompile(`[A-Z]`)
		lower    = regexp.MustCompile(`[a-z]`)
		number   = regexp.MustCompile(`[0-9]`)
		special  = regexp.MustCompile(`[!@#~$%^&*()+|_{}:<>?,./;'\[\]\\-]`)
		minChars = 8
	)

	return len(password) >= minChars &&
		upper.MatchString(password) &&
		lower.MatchString(password) &&
		number.MatchString(password) &&
		special.MatchString(password)
}
