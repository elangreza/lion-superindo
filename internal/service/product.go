package service

//go:generate mockgen -source $GOFILE -destination ../../mock/service/mock_$GOFILE -package mock$GOPACKAGE

import (
	"context"

	"github.com/elangreza14/superindo/internal/domain"
	"github.com/elangreza14/superindo/internal/params"
)

type (
	ProductRepo interface {
		ListProduct(ctx context.Context, req params.ProductQueryParams) ([]domain.Product, error)
		TotalProduct(ctx context.Context, req params.ProductQueryParams) (int, error)
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
	totalProducts, err := ps.repo.TotalProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	res := params.ProductsResponse{
		TotalData: totalProducts,
		Products:  []params.ProductResponse{},
	}

	if totalProducts == 0 {
		return &res, nil
	}

	products, err := ps.repo.ListProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	res.TotalPage = totalProducts / int(req.Limit)
	if totalProducts%int(req.Limit) != 0 {
		res.TotalPage++
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
