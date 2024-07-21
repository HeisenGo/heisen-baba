package presenter

import (
	"github.com/go-playground/validator/v10"
)

// This is the validator instance
// for more information see: https://github.com/go-playground/validator
var validate = validator.New()

type XValidator struct {
	validator *validator.Validate
}

func (v XValidator) Validate(data interface{}) []Response {
	var validationErrors []Response

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem Response

			elem.Success = false
			elem.Error = err.Error() // Export field value

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

var appValidator *XValidator

func GetValidator() *XValidator {
	if appValidator == nil {
		appValidator = &XValidator{
			validator: validate,
		}
	}
	return appValidator
}
