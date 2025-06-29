package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	errs "github.com/elangreza/lion-superindo/pkg/error"
)

func NewRoutes(productHandler *ProductHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/product", productHandler.ProductHandler)
	return mux
}

func Success(w http.ResponseWriter, status int, res any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{"data": res})
}

type APIError struct {
	Message string `json:"error"`
}

func Error(w http.ResponseWriter, status int, err error) {
	var apiErr APIError
	switch {
	case errors.As(err, &errs.AlreadyExistError{}):
		slog.Error("controller", "service", err.Error())
		status = errs.AlreadyExistError{}.HttpStatusCode()
		apiErr.Message = err.Error()
	case errors.As(err, &errs.ValidationError{}):
		slog.Error("controller", "request", err.Error())
		status = errs.ValidationError{}.HttpStatusCode()
		apiErr.Message = err.Error()
	case errors.As(err, &errs.MethodNotAllowedError{}):
		slog.Error("controller", "request", err.Error())
		status = errs.MethodNotAllowedError{}.HttpStatusCode()
		apiErr.Message = err.Error()
	default:
		slog.Error("controller", "service", err.Error())
		status = http.StatusInternalServerError
		apiErr.Message = "server error"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(apiErr)
}
