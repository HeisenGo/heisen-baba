package presenter

import (
	"encoding/json"
	"math"
	"time"

	"github.com/go-playground/validator/v10"
)

type Timestamp time.Time

const timestampFormat = "2006-01-02 15:04:05"

// MarshalJSON converts the Timestamp to JSON format.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	stamp := time.Time(t).Format(timestampFormat)
	return json.Marshal(stamp)
}

// UnmarshalJSON converts JSON data to a Timestamp.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var err error
	var parsedTime time.Time

	// Remove quotes from the JSON string
	trimmedData := string(data)
	trimmedData = trimmedData[1 : len(trimmedData)-1]

	parsedTime, err = time.Parse(timestampFormat, trimmedData)
	if err != nil {
		return err
	}

	*t = Timestamp(parsedTime)
	return nil
}

// String converts the Timestamp to string.
func (t Timestamp) String() string {
	return time.Time(t).Format(timestampFormat)
}

type PaginationResponse[T any] struct {
	Page       uint `json:"page"`
	PageSize   uint `json:"page_size"`
	TotalPages uint `json:"total_pages"`
	Data       []T  `json:"data"`
}

func NewPagination[T any](data []T, page, pageSize, total uint) *PaginationResponse[T] {
	totalPages := uint(0)
	if pageSize > 0 && total > 0 {
		totalPages = uint(math.Ceil(float64(total) / float64(pageSize)))
	}
	return &PaginationResponse[T]{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		Data:       data,
	}
}

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
