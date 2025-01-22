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

type ProductsResponse struct {
	TotalData int               `json:"total_data"`
	TotalPage int               `json:"total_page"`
	Products  []ProductResponse `json:"products"`
}

type ProductQueryParams struct {
	// can be search with
	// id and name
	Search string `form:"search"`
	// can be filtered by product type
	Types []string `form:"types"`
	// can be used with
	// sort=updated_at:desc,price:asc,name:desc
	Sorts   []string `form:"sorts"`
	sortMap map[string]string
	Limit   uint64 `form:"limit"`
	Page    uint64 `form:"page"`
}

func (pqr *ProductQueryParams) Validate() error {
	if pqr.Limit == 0 {
		pqr.Limit = 5
	}

	pqr.Search = strings.TrimSpace(pqr.Search)

	if len(pqr.Sorts) > 0 {
		pqr.sortMap = make(map[string]string)
		for _, sortRaw := range pqr.Sorts {
			sortStr := strings.Split(sortRaw, ":")
			if len(sortStr) != 2 {
				return errors.New("not valid sort format")
			}

			sortValue := sortStr[0]
			switch sortValue {
			case "updated_at", "price", "name":
			default:
				return errors.New("not valid sort value")
			}

			sortDirection := strings.ToLower(sortStr[1])
			switch sortDirection {
			case "asc", "desc":
			default:
				return errors.New("not valid sort direction")
			}

			pqr.sortMap[sortValue] = sortDirection
		}
	}

	return nil
}

func (pqr *ProductQueryParams) GetSortMapping() map[string]string {
	return pqr.sortMap
}
