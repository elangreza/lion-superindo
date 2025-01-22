package handler

import (
	"context"
	"net/http"

	"github.com/elangreza14/superindo/internal/params"
	"github.com/gin-gonic/gin"
)

type ProductService interface {
	ListProduct(ctx context.Context, args params.ProductQueryParams) ([]params.ProductResponse, error)
}

// ProductHandler ...
type ProductHandler struct {
	svc ProductService
}

// NewProductHandler ...
func NewProductHandler(svc ProductService) *ProductHandler {
	return &ProductHandler{
		svc: svc,
	}
}

func (ph *ProductHandler) ListProductHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := params.ProductQueryParams{}
		err := c.BindQuery(&query)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		if err = query.Validate(); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		res, err := ph.svc.ListProduct(c, query)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
