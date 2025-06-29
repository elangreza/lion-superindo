package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elangreza/lion-superindo/internal/domain"
	"github.com/elangreza/lion-superindo/internal/params"
)

const (
	countProductsKeys = "count"
	prefixProduct     = "product:"
)

func (pr *RedisRepo) CacheProducts(ctx context.Context, req params.ListProductsQueryParams, countProducts int, listProducts []domain.Product) error {
	keyRaw := prefixProduct + req.GetParamsKey()

	str, err := json.Marshal(listProducts)
	if err != nil {
		return err
	}

	products := make(map[string]any)
	products[req.GetOrderingKey()] = str
	products[countProductsKeys] = countProducts

	return pr.cache.HSet(ctx, keyRaw, products).Err()
}

func (pr *RedisRepo) GetCachedProducts(ctx context.Context, req params.ListProductsQueryParams) ([]domain.Product, error) {
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

func (pr *RedisRepo) GetCachedProductCount(ctx context.Context, req params.ListProductsQueryParams) (int, error) {
	keyRaw := prefixProduct + req.GetParamsKey()

	res, err := pr.cache.HGet(ctx, keyRaw, countProductsKeys).Result()
	if err != nil {
		return 0, err
	}

	total, err := strconv.Atoi(res)
	if err != nil {
		return 0, fmt.Errorf("failed to convert product total: %w", err)
	}
	return total, nil
}

func (pr *RedisRepo) FlushAllProducts(ctx context.Context) error {
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
