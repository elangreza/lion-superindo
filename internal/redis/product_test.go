package redis

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/elangreza/lion-superindo/internal/domain"
	"github.com/elangreza/lion-superindo/internal/params"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestProductRepo_CacheProducts(t *testing.T) {
	dbRedis, mockRedis := redismock.NewClientMock()
	pr := NewRepo(dbRedis)

	listProducts := []domain.Product{{ID: 1}}

	req := params.ListProductsQueryParams{}
	req.Validate()

	jsonListProduct, _ := json.Marshal(listProducts)

	products := make(map[string]any)
	products[req.GetOrderingKey()] = jsonListProduct
	products[countProductsKeys] = 1

	keyRaw := prefixProduct + req.GetParamsKey()
	mockRedis.ExpectHSet(keyRaw, products).SetVal(1)

	err := pr.CacheProducts(context.Background(), req, 1, listProducts)
	assert.NoError(t, err)
	if err := mockRedis.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductRepo_GetCachedProducts(t *testing.T) {
	dbRedis, mockRedis := redismock.NewClientMock()
	pr := NewRepo(dbRedis)

	listProducts := []domain.Product{{ID: 1}}

	req := params.ListProductsQueryParams{}
	req.Validate()

	jsonListProduct, _ := json.Marshal(listProducts)

	tableTest := []struct {
		name      string
		expectErr bool
		mock      func(m redismock.ClientMock)
	}{
		{
			name:      "success",
			expectErr: false,
			mock: func(m redismock.ClientMock) {
				keyRaw := prefixProduct + req.GetParamsKey()
				m.ExpectHGet(keyRaw, req.GetOrderingKey()).SetVal(string(jsonListProduct))
			},
		},
		{
			name:      "failed",
			expectErr: true,
			mock: func(m redismock.ClientMock) {
				keyRaw := prefixProduct + req.GetParamsKey()
				m.ExpectHGet(keyRaw, req.GetOrderingKey()).SetErr(errors.New("redis error"))
			},
		},
		{
			name:      "failed when parsing",
			expectErr: true,
			mock: func(m redismock.ClientMock) {
				keyRaw := prefixProduct + req.GetParamsKey()
				m.ExpectHGet(keyRaw, req.GetOrderingKey()).SetVal(string("1"))
			},
		},
	}

	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(mockRedis)
			products, err := pr.GetCachedProducts(context.Background(), req)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, products)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, listProducts, products)
			}
			if err := mockRedis.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestProductRepo_GetCachedProductCount(t *testing.T) {
	dbRedis, mockRedis := redismock.NewClientMock()
	pr := NewRepo(dbRedis)

	req := params.ListProductsQueryParams{}
	req.Validate()

	tableTest := []struct {
		name      string
		expectErr bool
		mock      func(m redismock.ClientMock)
		got       int
	}{
		{
			name:      "success",
			expectErr: false,
			mock: func(m redismock.ClientMock) {
				keyRaw := prefixProduct + req.GetParamsKey()
				m.ExpectHGet(keyRaw, countProductsKeys).SetVal(string("1"))
			},
			got: 1,
		},
		{
			name:      "failed",
			expectErr: true,
			mock: func(m redismock.ClientMock) {
				keyRaw := prefixProduct + req.GetParamsKey()
				m.ExpectHGet(keyRaw, countProductsKeys).SetErr(errors.New("redis error"))
			},
			got: 0,
		},
		{
			name:      "failed when parsing",
			expectErr: true,
			mock: func(m redismock.ClientMock) {
				keyRaw := prefixProduct + req.GetParamsKey()
				m.ExpectHGet(keyRaw, countProductsKeys).SetVal(string("a"))
			},
			got: 0,
		},
	}

	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(mockRedis)
			count, err := pr.GetCachedProductCount(context.Background(), req)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.got, count)
			if err := mockRedis.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestProductRepo_FlushAllProducts(t *testing.T) {
	dbRedis, mockRedis := redismock.NewClientMock()
	pr := NewRepo(dbRedis)
	pattern := prefixProduct + "*"

	tableTest := []struct {
		name      string
		expectErr bool
		mock      func(m redismock.ClientMock)
	}{
		{
			name: "success",
			mock: func(m redismock.ClientMock) {
				m.ExpectScanType(0, pattern, 100, "hash").SetVal([]string{"a"}, 0)
				m.ExpectDel("a").SetVal(1)
			},
			expectErr: false,
		},
		{
			name: "failed scan",
			mock: func(m redismock.ClientMock) {
				m.ExpectScanType(0, pattern, 100, "hash").SetErr(errors.New("scan error"))
			},
			expectErr: true,
		},
		{
			name: "failed delete",
			mock: func(m redismock.ClientMock) {
				m.ExpectScanType(0, pattern, 100, "hash").SetVal([]string{"a"}, 0)
				m.ExpectDel("a").SetErr(errors.New("delete error"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tableTest {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(mockRedis)
			err := pr.FlushAllProducts(context.Background())
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mockRedis.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
