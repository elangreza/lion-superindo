package service

//go:generate mockgen -source $GOFILE -destination ../../mock/service/mock_$GOFILE -package mock$GOPACKAGE

import (
	"context"
	"fmt"

	"github.com/elangreza14/lion-superindo/internal/domain"
	"github.com/elangreza14/lion-superindo/internal/params"
	errs "github.com/elangreza14/lion-superindo/pkg/error"
	"github.com/redis/go-redis/v9"
)

type (
	DbRepo interface {
		ListProduct(ctx context.Context, req params.ListProductQueryParams) ([]domain.Product, error)
		TotalProduct(ctx context.Context, req params.ListProductQueryParams) (int, error)
		CreateProduct(ctx context.Context, req params.CreateProductRequest) (int, error)
	}

	CacheRepo interface {
		FlushAll(ctx context.Context) error
		SetProduct(ctx context.Context, req params.ListProductQueryParams, totalProducts int, listProducts []domain.Product) error
		GetProductData(ctx context.Context, req params.ListProductQueryParams) (listProducts []domain.Product, err error)
		GetProductTotal(ctx context.Context, req params.ListProductQueryParams) (totalProducts int, err error)
	}

	ProductService struct {
		db    DbRepo
		cache CacheRepo
	}
)

func NewProductService(repo DbRepo, cache CacheRepo) *ProductService {
	return &ProductService{
		db:    repo,
		cache: cache,
	}
}

func (ps *ProductService) ListProduct(ctx context.Context, req params.ListProductQueryParams) (*params.ListProductResponses, error) {
	products, err := ps.cache.GetProductData(ctx, req)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("cache error: %w", err)
	}

	if len(products) == 0 && err == redis.Nil {
		products, err = ps.db.ListProduct(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("db error: %w", err)
		}
	}

	totalProducts, err := ps.cache.GetProductTotal(ctx, req)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("cache error (total products): %w", err)
	}

	if totalProducts == 0 && err == redis.Nil {
		totalProducts, err = ps.db.TotalProduct(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("db error (total products): %w", err)
		}
	}

	err = ps.cache.SetProduct(ctx, req, totalProducts, products)
	if err != nil {
		return nil, err
	}

	res := params.ListProductResponses{}

	if totalProducts == 0 {
		return &res, nil
	}

	res.TotalData = totalProducts
	res.TotalPage = (totalProducts + int(req.Limit) - 1) / int(req.Limit)

	res.Products = make([]params.ProductResponse, 0, len(products))
	for _, product := range products {
		res.Products = append(res.Products, params.ProductResponse{
			ID:        product.ID,
			Name:      product.Name,
			Price:     product.Price,
			Type:      product.ProductType.Name,
			CreatedAt: product.CreatedAt,
		})
	}

	return &res, nil
}

func (ps *ProductService) CreateProduct(ctx context.Context, req params.CreateProductRequest) (*params.CreateProductResponse, error) {
	products, err := ps.db.TotalProduct(ctx, params.ListProductQueryParams{Search: req.Name})
	if err != nil {
		return nil, err
	}

	if products > 0 {
		return nil, errs.AlreadyExistError{
			Message: fmt.Sprintf("product %s", req.Name),
		}
	}

	id, err := ps.db.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	if err := ps.cache.FlushAll(ctx); err != nil {
		return nil, fmt.Errorf("failed to flush cache: %w", err)
	}

	return &params.CreateProductResponse{ID: id}, nil
}
