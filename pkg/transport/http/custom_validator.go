package inthttp

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	CustomErrorMessages = map[string]string{
		"password_pattern": "Password Must have number, uppercase, lower case and special character",
	}
)

func RegisterCustomValigator(validator *validator.Validate) {
	validator.RegisterValidation("password_pattern", PasswordPattern)
}

func PasswordPattern(fl validator.FieldLevel) bool {
	pw := fl.Field().String()

	// Why not just comiple in one regex? smh golang doesn't support "?=" looks ahead syntax so... i made this
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(pw)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pw)
	hasDigit := regexp.MustCompile(`\d`).MatchString(pw)
	hasSpecial := regexp.MustCompile(`\W`).MatchString(pw)
	hasSpace := regexp.MustCompile(`\s`).MatchString(pw)

	return hasLower && hasUpper && hasDigit && hasSpecial && !hasSpace
}
