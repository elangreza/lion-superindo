package params

import (
	"errors"
	"strings"
	"time"
)

type ProductResponse struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Quantity int       `json:"quantity"`
	Price    int       `json:"price"`
	Type     string    `json:"type"`
	UpdateAt time.Time `json:"updated_at"`
}

type ProductQueryParams struct {
	// can be search with
	// id and name
	Search string `json:"search"`
	// can be filtered by product type
	Types []string `json:"types"`
	// can be used with
	// sort=updated_at:desc,price:asc,name:desc
	Sort []string `json:"sort"`
}

func (pqr *ProductQueryParams) Validate() error {
	pqr.Search = strings.TrimSpace(pqr.Search)

	for index, sortRaw := range pqr.Sort {
		sortStr := strings.Split(sortRaw, ":")
		if len(sortStr) != 2 {
			return errors.New("not valid sort format")
		}

		sortValue := sortStr[0]
		if sortValue != "updated_at" &&
			sortValue != "price" &&
			sortValue != "name" {
			return errors.New("not valid sort value")
		}

		sortDirection := sortStr[0]
		if sortDirection != "ASC" &&
			sortDirection != "DESC" &&
			sortDirection != "asc" &&
			sortDirection != "desc" {
			return errors.New("not valid sort direction")
		}

		pqr.Sort[index] = sortValue + " " + sortDirection
	}

	return nil
}
