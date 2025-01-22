package service

import (
	"context"

	"github.com/elangreza14/superindo/internal/domain"
	"github.com/elangreza14/superindo/internal/params"
)

type (
	ProductRepo interface {
		ListProduct(ctx context.Context, req params.ProductQueryParams) (int, []domain.Product, error)
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

func (ps *ProductService) ListProduct(ctx context.Context, req params.ProductQueryParams) (*params.ProductsResponse, error) {
	totalProducts, products, err := ps.repo.ListProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	totalPage := totalProducts / int(req.Limit)
	if totalProducts%int(req.Limit) != 0 {
		totalPage++
	}

	res := params.ProductsResponse{
		TotalData: totalProducts,
		TotalPage: totalPage,
		Products:  []params.ProductResponse{},
	}
	for _, product := range products {
		updatedAt := product.CreatedAt
		if product.UpdatedAt != nil {
			updatedAt = *product.UpdatedAt
		}
		res.Products = append(res.Products,
			params.ProductResponse{
				ID:       product.ID,
				Name:     product.Name,
				Quantity: product.Quantity,
				Price:    product.Price,
				Type:     product.ProductType.Name,
				UpdateAt: updatedAt,
			})
	}

	return &res, nil
}
