package params

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	errs "github.com/elangreza14/lion-superindo/pkg/error"
)

type ProductResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type ListProductsResponses struct {
	TotalData int               `json:"total_data"`
	TotalPage int               `json:"total_page"`
	Products  []ProductResponse `json:"products"`
}

type ListProductsQueryParams struct {
	// can be searched by id or name
	Search string
	// can be filtered by product type
	Types []string

	// local var. used for caching key
	paramsKey string

	// Embedding PaginationParams for pagination and sorting
	PaginationParams
}

func (pqr *ListProductsQueryParams) Validate() error {

	if err := pqr.PaginationParams.Validate(); err != nil {
		return errs.ValidationError{Err: err}
	}

	validSortKeys := map[string]bool{
		"id": true, "created_at": true, "price": true, "name": true,
	}
	for sortKey := range pqr.GetSortMapping() {
		if !validSortKeys[sortKey] {
			return errs.ValidationError{Message: fmt.Sprint("%s not valid sort key", sortKey)}
		}
	}

	pqr.Search = strings.TrimSpace(pqr.Search)
	mapKey := map[string]any{
		"search": pqr.Search,
		"types":  pqr.Types,
	}

	key, err := json.Marshal(mapKey)
	if err != nil {
		return errs.ValidationError{Message: "failed to marshal query params for cache key"}
	}

	pqr.paramsKey = string(key)

	return nil
}

func (pqr *ListProductsQueryParams) GetParamsKey() string {
	return pqr.paramsKey
}

type CreateProductRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Type  string `json:"type"`
}

type CreateProductResponse struct {
	ID int `json:"id"`
}

func (pqr *CreateProductRequest) Validate() error {
	if len(pqr.Name) == 0 {
		return errs.ValidationError{Message: "name cannot be empty"}
	}
	if len(pqr.Type) == 0 {
		return errs.ValidationError{Message: "type cannot be empty"}
	}
	pqr.Type = strings.ToLower(pqr.Type)
	if pqr.Price < 0 {
		return errs.ValidationError{Message: "price cannot be negative"}
	}
	return nil
}
