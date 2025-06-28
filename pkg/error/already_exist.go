package errs

import (
	"fmt"
	"net/http"
)

type AlreadyExistError struct {
	Message string
}

func (a AlreadyExistError) Error() string {
	if a.Message == "" {
		return "already exist"
	}

	return fmt.Sprintf("%s already exist", a.Message)
}

func (a AlreadyExistError) HttpStatusCode() int {
	return http.StatusConflict
}
