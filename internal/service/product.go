package service

//go:generate mockgen -source $GOFILE -destination ../../mock/service/mock_$GOFILE -package mock$GOPACKAGE

import (
	"context"
	"errors"

	"github.com/elangreza14/superindo/internal/domain"
	"github.com/elangreza14/superindo/internal/params"
)

type (
	ProductRepo interface {
		ListProduct(ctx context.Context, req params.ListProductQueryParams) ([]domain.Product, error)
		TotalProduct(ctx context.Context, req params.ListProductQueryParams, withCache bool) (int, error)
		CreateProduct(ctx context.Context, req params.CreateProductRequest) (int, error)
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

func (ps *ProductService) ListProduct(ctx context.Context, req params.ListProductQueryParams) (*params.ListProductResponses, error) {
	totalProducts, err := ps.repo.TotalProduct(ctx, req, true)
	if err != nil {
		return nil, err
	}

	res := params.ListProductResponses{
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
				Price:    product.Price,
				Type:     product.ProductType.Name,
				UpdateAt: updatedAt,
			})
	}

	return &res, nil
}

func (ps *ProductService) CreateProduct(ctx context.Context, req params.CreateProductRequest) (*params.CreateProductResponse, error) {
	products, err := ps.repo.TotalProduct(ctx, params.ListProductQueryParams{Search: req.Name}, false)
	if err != nil {
		return nil, err
	}

	if products > 0 {
		return nil, errors.New("product already exist")
	}

	id, err := ps.repo.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return &params.CreateProductResponse{ID: id}, nil
}
