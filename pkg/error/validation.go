package errs

import (
	"fmt"
	"net/http"
)

type ValidationError struct {
	Message string
	Err     error
}

func (v ValidationError) Error() string {
	if v.Message == "" {
		return "validation error"
	}

	if v.Err != nil {
		return fmt.Sprintf("validation error: %s", v.Err.Error())
	}

	return fmt.Sprintf("validation error: %s", v.Message)
}

func (a ValidationError) HttpStatusCode() int {
	return http.StatusBadRequest
}
