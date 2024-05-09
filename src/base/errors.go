package base

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	ory "github.com/ory/kratos-client-go"
	"github.com/sirupsen/logrus"
	"net/http"
)

type FieldError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type ServiceError struct {
	error
	Summary string
	Detail  any
	Status  int
}

func (sError *ServiceError) String() string {
	return sError.Summary
}

func logKratosError(message string, originalError error) {
	var kratosError *ory.GenericOpenAPIError
	converted := errors.As(originalError, &kratosError)
	var model ory.ErrorGeneric

	if converted {
		model, converted = kratosError.Model().(ory.ErrorGeneric)
	}
	if converted {
		Logger.WithFields(logrus.Fields{
			"message":        message,
			"error":          originalError.Error(),
			"kratos_message": model.Error.Message,
			"reason":         model.Error.Reason,
		}).Warn("Error interaction with Kratos server")
	} else {
		Logger.WithFields(logrus.Fields{
			"message":       message,
			"error_message": originalError.Error(),
		}).Warn("Error interaction with Kratos server")
	}
}

func NewKratosError(message string, originalError error) ServiceError {
	logKratosError(message, originalError)
	return ServiceError{Summary: message}
}

func NewQueryParamError(paramName string, err error) ServiceError {
	return ServiceError{
		Summary: fmt.Sprintf("Invalid format for query param '%s'", paramName),
		Detail:  err.Error(),
		Status:  http.StatusBadRequest,
	}
}

func NewPathParamRequiredError(paramName string) ServiceError {
	return ServiceError{
		Summary: fmt.Sprintf("Path parameter '%s' required", paramName),
		Status:  http.StatusBadRequest,
	}
}

func NewPathParamError(paramName string, err error) ServiceError {
	return ServiceError{
		Summary: fmt.Sprintf("Invalid format for path param '%s'", paramName),
		Detail:  err.Error(),
		Status:  http.StatusBadRequest,
	}
}

func WrapValidationErrors(err error) error {
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		errorDetails := make([]FieldError, len(validationErr))
		for i, fe := range validationErr {
			detail := getErrorMessageForTag(fe.Tag())
			if detail == "" {
				detail = "Unknown validation error for " + fe.Field()
			}
			errorDetails[i] = FieldError{
				Name:    fe.Field(),
				Message: detail,
			}
		}
		return ServiceError{
			Summary: "Data validation failed",
			Detail:  errorDetails,
			Status:  422,
		}
	} else {
		return ServiceError{
			Summary: err.Error(),
			Status:  422,
		}
	}
}
