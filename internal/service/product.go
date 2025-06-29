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
		ListProducts(ctx context.Context, req params.ListProductsQueryParams) ([]domain.Product, error)
		CountProducts(ctx context.Context, req params.ListProductsQueryParams) (int, error)
		CreateProduct(ctx context.Context, req params.CreateProductRequest) (int, error)
	}

	CacheRepo interface {
		FlushAllProducts(ctx context.Context) error
		CacheProducts(ctx context.Context, req params.ListProductsQueryParams, countProducts int, listProducts []domain.Product) error
		GetCachedProducts(ctx context.Context, req params.ListProductsQueryParams) (listProducts []domain.Product, err error)
		GetCachedProductCount(ctx context.Context, req params.ListProductsQueryParams) (countProducts int, err error)
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

func (ps *ProductService) ListProducts(ctx context.Context, req params.ListProductsQueryParams) (*params.ListProductsResponses, error) {
	products, err := ps.cache.GetCachedProducts(ctx, req)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("cache error: %w", err)
	}

	var isListProductsFromDB bool
	if len(products) == 0 && err == redis.Nil {
		products, err = ps.db.ListProducts(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("db error: %w", err)
		}
		isListProductsFromDB = true
	}

	countProducts, err := ps.cache.GetCachedProductCount(ctx, req)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("cache error (total products): %w", err)
	}

	var isCountProductsFromDB bool
	if countProducts == 0 && err == redis.Nil {
		countProducts, err = ps.db.CountProducts(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("db error (total products): %w", err)
		}
		isCountProductsFromDB = true
	}

	// If both products and countProducts are from the database, cache them
	// If either products or countProducts is from the cache, we do not cache them again
	if isListProductsFromDB || isCountProductsFromDB {
		err = ps.cache.CacheProducts(ctx, req, countProducts, products)
		if err != nil {
			return nil, err
		}
	}

	res := params.ListProductsResponses{}

	if countProducts == 0 {
		return &res, nil
	}

	res.TotalData = countProducts
	res.TotalPage = (countProducts + int(req.Limit) - 1) / int(req.Limit)

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
	products, err := ps.db.CountProducts(ctx, params.ListProductsQueryParams{Search: req.Name})
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

	if err := ps.cache.FlushAllProducts(ctx); err != nil {
		return nil, fmt.Errorf("failed to flush cache: %w", err)
	}

	return &params.CreateProductResponse{ID: id}, nil
}
