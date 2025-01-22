package handler

import (
	"context"
	"net/http"

	"github.com/elangreza14/superindo/internal/params"
	"github.com/gin-gonic/gin"
)

type ProductService interface {
	ListProduct(ctx context.Context, args params.ListProductQueryParams) (*params.ListProductResponses, error)
	CreateProduct(ctx context.Context, req params.CreateProductRequest) error
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
		query := params.ListProductQueryParams{}
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
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func (ph *ProductHandler) CreateProductHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		body := params.CreateProductRequest{}
		err := c.BindJSON(&body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		if err = body.Validate(); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}

		if err := ph.svc.CreateProduct(c, body); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusCreated, nil)
	}
}
