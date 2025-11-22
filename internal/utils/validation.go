package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// FormatValidationErrorsArray converts validator errors to user-friendly messages (multiple errors per field)
func FormatValidationErrorsArray(errs validator.ValidationErrors) map[string][]string {
	errors := make(map[string][]string)

	for _, err := range errs {
		fieldName := ToSnakeCase(err.Field())

		var message string
		switch err.Tag() {
		case "required":
			message = "This field is required."
		case "min":
			message = "Must be at least " + err.Param() + " characters."
		case "email":
			message = "Must be a valid email address."
		case "phone_e164":
			message = "Must be a valid phone number in E.164 format (e.g., +407xxxxxxxx)."
		case "strong_password":
			message = "Must be at least 8 characters long and contain both letters and digits."
		case "eqfield":
			message = "Passwords do not match."
		case "eq":
			if err.Field() == "AcceptedTerms" {
				message = "You must accept the terms and conditions."
			} else {
				message = "Invalid value."
			}
		default:
			message = "Invalid value for " + fieldName + "."
		}

		errors[fieldName] = append(errors[fieldName], message)
	}

	return errors
}

// FormatValidationErrorsSimple converts validator errors to simple string messages (one error per field)
func FormatValidationErrorsSimple(errors validator.ValidationErrors) map[string]string {
	errorMessages := make(map[string]string)
	for _, err := range errors {
		field := err.Field()
		switch err.Tag() {
		case "required":
			errorMessages[field] = field + " is required"
		case "email":
			errorMessages[field] = "Invalid email format"
		case "min":
			errorMessages[field] = field + " must be at least " + err.Param() + " characters"
		case "max":
			errorMessages[field] = field + " must not exceed " + err.Param() + " characters"
		case "gt":
			errorMessages[field] = field + " must be greater than " + err.Param()
		case "oneof":
			errorMessages[field] = field + " must be one of: " + err.Param()
		default:
			errorMessages[field] = field + " is invalid"
		}
	}
	return errorMessages
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
