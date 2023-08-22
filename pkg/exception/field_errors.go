package exception

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/kondohiroki/go-grpc-boilerplate/pkg/validation"
)

// NewValidationFailedErrors reads ValidationErrors and converts to format of ExceptionErrors.
func NewValidationFailedErrors(validationErrs validator.ValidationErrors) *ExceptionErrors {
	errItems := make([]*ExceptionError, 0, len(validationErrs))
	for _, validationErr := range validationErrs {
		errItems = append(errItems, &ExceptionError{
			Message:      validation.Translate(validationErr),
			Type:         ERROR_TYPE_VALIDATION_ERROR,
			ErrorSubcode: SUBCODE_VALIDATION_FAILED,
		})
	}
	return &ExceptionErrors{
		GlobalMessage:  "validation failed",
		HttpStatusCode: http.StatusUnprocessableEntity,
		ErrItems:       errItems,
	}
}

// IsEmpty reports whether there is at least one errItems in an instance.
func (cErrs *ExceptionErrors) IsEmpty() bool {
	return len(cErrs.ErrItems) == 0
}

// Append manually appends error item to current errors.
func (cErrs *ExceptionErrors) Append(cErr *ExceptionError) *ExceptionErrors {
	cErrs.ErrItems = append(cErrs.ErrItems, cErr)
	return cErrs
}
