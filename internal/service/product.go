package service

import (
	"github.com/elangreza14/superindo/internal/domain"
	"github.com/elangreza14/superindo/internal/params"
)

type (
	ProductRepo interface {
		ListProduct(req params.ProductQueryParams) ([]domain.Product, error)
	}

	ProductService struct {
		repo ProductRepo
	}
)

func NewProductService(repo ProductRepo) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (ps *ProductService) ListProduct(req params.ProductQueryParams) ([]params.ProductResponse, error) {
	products, err := ps.repo.ListProduct(req)
	if err != nil {
		return nil, err
	}

	res := []params.ProductResponse{}
	for _, product := range products {
		updatedAt := product.CreatedAt
		if product.UpdatedAt != nil {
			updatedAt = *product.UpdatedAt
		}
		res = append(res,
			params.ProductResponse{
				ID:       product.ID,
				Name:     product.Name,
				Quantity: product.Quantity,
				Price:    product.Price,
				Type:     product.ProductType.Name,
				UpdateAt: updatedAt,
			})
	}

	return res, nil
}
