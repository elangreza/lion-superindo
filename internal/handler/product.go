package handler

//go:generate mockgen -source $GOFILE -destination ../../mock/handler/mock_$GOFILE -package mock$GOPACKAGE

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/elangreza/lion-superindo/internal/params"
	errs "github.com/elangreza/lion-superindo/pkg/error"
)

type (
	ProductService interface {
		ListProducts(ctx context.Context, args params.ListProductsQueryParams) (*params.ListProductsResponses, error)
		CreateProduct(ctx context.Context, req params.CreateProductRequest) (*params.CreateProductResponse, error)
	}

	ProductHandler struct {
		svc ProductService
	}
)

func NewProductHandler(svc ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// ListProductsHandler godoc
//
//	@Summary		Get products
//	@Description	Get all products
//	@Tags			product
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	params.ListProductsResponses
//	@Failure		400	{object}	handler.APIError	"validation error"
//	@Failure		500	{object}	handler.APIError	"server error"
//	@Router			/product [get]
//	@Param			page	query	int			false	"Page number, default 1"
//	@Param			limit	query	int			false	"Limit number of products, default 10"
//	@Param			search	query	string		false	"Search by product name or id"
//	@Param			type	query	[]string	false	"Filter by product type. Repeat param for multiple values (e.g. type=buah&type=snack) or use comma-separated (type=buah,snack)."
//	@Param			sort	query	[]string	false	"Sort by field. Values can be created_at:asc, created_at:desc, price:asc, price:desc, name:asc, name:desc, id:asc, id:desc. Default: id:asc"
func (ph *ProductHandler) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	query := &params.ListProductsQueryParams{}

	if r.URL.Query().Get("page") != "" {
		query.Page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			Error(w, http.StatusBadRequest, errs.ValidationError{Message: "not valid page"})
			return
		}
	}

	if r.URL.Query().Get("limit") != "" {
		query.Limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			Error(w, http.StatusBadRequest, errs.ValidationError{Message: "not valid limit"})
			return
		}
	}

	query.Search = r.URL.Query().Get("search")
	query.Types = r.URL.Query()["type"]
	query.Sorts = r.URL.Query()["sort"]
	if err := query.Validate(); err != nil {
		Error(w, http.StatusBadRequest, err)
		return
	}

	res, err := ph.svc.ListProducts(r.Context(), *query)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	Success(w, http.StatusOK, res)
}

// CreateProductHandler godoc
//
//	@Summary		Create product
//	@Description	Create a new product
//	@Tags			product
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	params.CreateProductResponse
//	@Failure		400	{object}	handler.APIError	"validation error"
//	@Failure		409	{object}	handler.APIError	"conflict error, if product with same name already exists"
//	@Failure		500	{object}	handler.APIError	"server error"
//	@Router			/product [post]
//	@Param			body	body	params.CreateProductRequest	true	"Product data"
func (ph *ProductHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	body := params.CreateProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		Error(w, http.StatusBadRequest, errs.ValidationError{Message: err.Error()})
		return
	}

	if err := body.Validate(); err != nil {
		Error(w, http.StatusBadRequest, err)
		return
	}

	res, err := ph.svc.CreateProduct(r.Context(), body)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	Success(w, http.StatusCreated, res)
}

func (ph *ProductHandler) ProductHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ph.ListProductsHandler(w, r)
	case http.MethodPost:
		ph.CreateProductHandler(w, r)
	default:
		Error(w, http.StatusMethodNotAllowed, errs.MethodNotAllowedError{
			Method: r.Method,
		})
	}
}
