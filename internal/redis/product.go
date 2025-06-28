package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elangreza14/superindo/internal/domain"
	"github.com/elangreza14/superindo/internal/params"
	"github.com/redis/go-redis/v9"
)

type (
	ProductRepo struct {
		cache *redis.Client
	}
)

const (
	totalProductKeys = "total"
	prefixProduct    = "product:"
)

func NewProductRepo(cache *redis.Client) *ProductRepo {
	return &ProductRepo{
		cache: cache,
	}
}

func (pr *ProductRepo) SetProduct(ctx context.Context, req params.ListProductQueryParams, totalProducts int, listProducts []domain.Product) error {
	keyRaw := prefixProduct + req.GetParamsKey()

	products := make(map[string]any)

	str, err := json.Marshal(listProducts)
	if err != nil {
		return err
	}
	products[req.GetOrderingKey()] = str
	products[totalProductKeys] = totalProducts

	if err := pr.cache.HSet(ctx, keyRaw, products).Err(); err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepo) GetProductData(ctx context.Context, req params.ListProductQueryParams) ([]domain.Product, error) {
	keyRaw := prefixProduct + req.GetParamsKey()

	res, err := pr.cache.HGet(ctx, keyRaw, req.GetOrderingKey()).Result()
	if err != nil {
		return nil, err
	}

	var listProducts []domain.Product
	if err := json.Unmarshal([]byte(res), &listProducts); err != nil {
		return nil, err
	}
	return listProducts, nil
}

func (pr *ProductRepo) GetProductTotal(ctx context.Context, req params.ListProductQueryParams) (int, error) {
	keyRaw := prefixProduct + req.GetParamsKey()

	res, err := pr.cache.HGet(ctx, keyRaw, totalProductKeys).Result()
	if err != nil {
		return 0, err
	}

	total, err := strconv.Atoi(res)
	if err != nil {
		return 0, fmt.Errorf("failed to convert product total: %w", err)
	}
	return total, nil
}

func (pr *ProductRepo) FlushAll(ctx context.Context) error {
	var cursor uint64
	pattern := prefixProduct + "*"
	for {
		keys, nextCursor, err := pr.cache.ScanType(ctx, cursor, pattern, 100, "hash").Result()
		if err != nil {
			return err
		}
		for _, key := range keys {
			if err := pr.cache.Del(ctx, key).Err(); err != nil {
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}
