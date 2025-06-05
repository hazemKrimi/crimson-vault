package api

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (v *CustomValidator) Validate(i any) error {
	if err := v.validator.Struct(i); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)

			for _, ve := range validationErrors {
				field := ve.Field()
				tag := ve.Tag()

				var msg string

				switch tag {
				case "required":
					msg = fmt.Sprintf("%s is required!", field)
				case "email":
					msg = fmt.Sprintf("%s must be a valid email!", field)
				case "alpha":
					msg = fmt.Sprintf("%s must only contain alphabetic characters!", field)
				case "e164":
					msg = fmt.Sprintf("%s must be a valid phone number in e164 format!", field)
				case "password":
					msg = fmt.Sprintf("%s must have at lease one uppercase, one lowercase, one number and one special character!", field)
				case "eqcsfield":
					msg = fmt.Sprintf("%s must be the same as %s!", field, ve.Param())
				default:
					msg = fmt.Sprintf("%s is not valid!", field)
				}

				errors[field] = msg
			}
			return echo.NewHTTPError(http.StatusBadRequest, errors)
		}
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
