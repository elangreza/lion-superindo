package handler

//go:generate mockgen -source $GOFILE -destination ../../mock/handler/mock_$GOFILE -package mock$GOPACKAGE

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/elangreza14/superindo/internal/params"
)

type (
	ProductService interface {
		ListProduct(ctx context.Context, args params.ListProductQueryParams) (*params.ListProductResponses, error)
		CreateProduct(ctx context.Context, req params.CreateProductRequest) (*params.CreateProductResponse, error)
	}

	ProductHandler struct {
		svc ProductService
	}
)

func NewProductHandler(svc ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (ph *ProductHandler) ListProductHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	query := &params.ListProductQueryParams{}

	if r.URL.Query().Get("page") != "" {
		query.Page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			Error(w, http.StatusBadRequest, errors.New("not valid page"))
			return
		}
	}

	if r.URL.Query().Get("limit") != "" {
		query.Limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			Error(w, http.StatusBadRequest, errors.New("not valid limit"))
			return
		}
	}

	query.Search = r.URL.Query().Get("search")
	query.Types = r.URL.Query()["types"]
	query.Sorts = r.URL.Query()["sorts"]
	if err := query.Validate(); err != nil {
		Error(w, http.StatusBadRequest, err)
		return
	}

	res, err := ph.svc.ListProduct(r.Context(), *query)
	if err != nil {
		slog.Error("controller", "service", err.Error())
		Error(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	Success(w, http.StatusOK, res)
}

func (ph *ProductHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	body := params.CreateProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		Error(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		Error(w, http.StatusBadRequest, err)
		return
	}

	res, err := ph.svc.CreateProduct(r.Context(), body)
	if err != nil {
		slog.Error("controller", "service", err.Error())
		if err.Error() == "product already exist" {
			Error(w, http.StatusConflict, errors.New("product already exist"))
			return
		}
		Error(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	Success(w, http.StatusCreated, res)
}

func (ph *ProductHandler) ProductHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ph.ListProductHandler(w, r)
	case http.MethodPost:
		ph.CreateProductHandler(w, r)
	default:
		Error(w, http.StatusMethodNotAllowed, errors.New("invalid method"))
	}
}
