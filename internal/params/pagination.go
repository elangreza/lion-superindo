package params

import (
	"errors"
	"fmt"
	"strings"
)

type PaginationParams struct {
	Sorts []string
	Limit int
	Page  int

	// local var. Used for sorting in the DB
	sortMap map[string]string
	// local var. Used for cache data ordering
	orderingKey string
}

func (pqr *PaginationParams) Validate() error {
	if pqr.Page < 1 {
		pqr.Page = 1
	}

	if pqr.Limit < 1 {
		pqr.Limit = 5
	}

	var orderingBuilder strings.Builder
	if len(pqr.Sorts) > 0 {
		newSorts := []string{}
		for _, sort := range pqr.Sorts {
			if strings.Contains(sort, ",") {
				newSorts = append(newSorts, strings.Split(sort, ",")...)
			} else {
				newSorts = append(newSorts, sort)
			}
		}

		pqr.sortMap = make(map[string]string)
		for _, sortRaw := range newSorts {
			parts := strings.Split(sortRaw, ":")
			if len(parts) != 2 {
				return fmt.Errorf("%s is not valid sort format", sortRaw)
			}

			value := strings.ToLower(strings.TrimSpace(parts[0]))
			direction := strings.ToLower(strings.TrimSpace(parts[1]))

			if direction != "asc" && direction != "desc" {
				return errors.New("not valid sort direction")
			}

			pqr.sortMap[value] = direction
			orderingBuilder.WriteString(sortRaw)
		}
	}

	pqr.orderingKey = fmt.Sprintf("%d:%d:%s", pqr.Limit, pqr.Page, orderingBuilder.String())

	return nil
}

func (pqr *PaginationParams) GetSortMapping() map[string]string {
	return pqr.sortMap
}

func (pqr *PaginationParams) GetOrderingKey() string {
	return pqr.orderingKey
}
