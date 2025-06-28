package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	errs "github.com/elangreza14/lion-superindo/pkg/error"
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

func Error(w http.ResponseWriter, status int, err error) {
	if errors.As(err, &errs.AlreadyExistError{}) {
		status = errs.AlreadyExistError{}.HttpStatusCode()
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
