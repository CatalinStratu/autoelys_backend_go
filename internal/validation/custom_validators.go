package validation

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidators(v *validator.Validate) error {
	if err := v.RegisterValidation("phone_e164", validatePhoneE164); err != nil {
		return err
	}
	if err := v.RegisterValidation("strong_password", validateStrongPassword); err != nil {
		return err
	}
	return nil
}

func validatePhoneE164(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true
	}
	phoneRegex := regexp.MustCompile(`^\+\d{1,3}\d{9,12}$`)
	return phoneRegex.MatchString(phone)
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasLetter := false
	hasDigit := false

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}
