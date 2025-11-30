package error

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ValidationsResponse struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

var ErrValidator = map[string]string{}

func ErrValidationResponse(err error) (validationsResponse []ValidationsResponse) {
	var fieldErrors validator.ValidationErrors
	if errors.As(err, &fieldErrors) {
		for _, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				validationsResponse = append(validationsResponse, ValidationsResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is required", err.Field()),
				})
			case "email":
				validationsResponse = append(validationsResponse, ValidationsResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is not a valid email", err.Field()),
				})
			default:
				ErrValidator, ok := ErrValidator[err.Tag()]
				if ok {
					count := strings.Count(ErrValidator, "%s")
					if count == 1 {
						validationsResponse = append(validationsResponse, ValidationsResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(ErrValidator, err.Field()),
						})
					} else {
						validationsResponse = append(validationsResponse, ValidationsResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(ErrValidator, err.Field(), err.Param()),
						})
					}
				} else {
					validationsResponse = append(validationsResponse, ValidationsResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf("something went wrong on %s; %s", err.Field(), err.Tag()),
					})
				}
			}
		}
	}

	return validationsResponse
}

func WrapError(err error) error {
	logrus.Errorf("Error: %v", err)
	return err
}
