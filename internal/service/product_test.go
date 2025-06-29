package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/elangreza/lion-superindo/internal/domain"
	"github.com/elangreza/lion-superindo/internal/params"
	mockservice "github.com/elangreza/lion-superindo/mock/service"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type TestProductServiceSuite struct {
	suite.Suite

	MockDbRepo    *mockservice.MockDbRepo
	MockCacheRepo *mockservice.MockCacheRepo
	Ps            *ProductService
	Ctrl          *gomock.Controller
}

func (suite *TestProductServiceSuite) SetupSuite() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockDbRepo = mockservice.NewMockDbRepo(suite.Ctrl)
	suite.MockCacheRepo = mockservice.NewMockCacheRepo(suite.Ctrl)
	suite.Ps = NewProductService(suite.MockDbRepo, suite.MockCacheRepo)
}

func (suite *TestProductServiceSuite) TearDownSuite() {
	suite.Ctrl.Finish()
}

func TestProductServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TestProductServiceSuite))
}

func (suite *TestProductServiceSuite) TestProductService_GetProducts() {
	suite.Run("error GetCachedProducts", func() {
		req := params.ListProductsQueryParams{}
		ctx := context.Background()
		suite.MockCacheRepo.EXPECT().GetCachedProducts(ctx, req).Return(nil, errors.New("test"))

		got, err := suite.Ps.ListProducts(ctx, req)
		suite.Error(err)
		suite.Nil(got)
	})

	suite.Run("err ListProducts", func() {
		req := params.ListProductsQueryParams{}
		ctx := context.Background()
		suite.MockCacheRepo.EXPECT().GetCachedProducts(ctx, req).Return([]domain.Product{}, redis.Nil)
		suite.MockDbRepo.EXPECT().ListProducts(ctx, req).Return(nil, errors.New("test"))

		got, err := suite.Ps.ListProducts(ctx, req)
		suite.Error(err)
		suite.Nil(got)
	})

	suite.Run("err GetCachedProductCount", func() {
		req := params.ListProductsQueryParams{}
		ctx := context.Background()
		suite.MockCacheRepo.EXPECT().GetCachedProducts(ctx, req).Return([]domain.Product{}, redis.Nil)
		suite.MockDbRepo.EXPECT().ListProducts(ctx, req).Return([]domain.Product{}, nil)
		suite.MockCacheRepo.EXPECT().GetCachedProductCount(ctx, req).Return(0, errors.New("test"))

		got, err := suite.Ps.ListProducts(ctx, req)
		suite.Error(err)
		suite.Nil(got)
	})

	suite.Run("err CountProducts", func() {
		req := params.ListProductsQueryParams{}
		ctx := context.Background()
		suite.MockCacheRepo.EXPECT().GetCachedProducts(ctx, req).Return([]domain.Product{}, redis.Nil)
		suite.MockDbRepo.EXPECT().ListProducts(ctx, req).Return([]domain.Product{}, nil)
		suite.MockCacheRepo.EXPECT().GetCachedProductCount(ctx, req).Return(0, redis.Nil)
		suite.MockDbRepo.EXPECT().CountProducts(ctx, req).Return(0, errors.New("test"))

		got, err := suite.Ps.ListProducts(ctx, req)
		suite.Error(err)
		suite.Nil(got)
	})

	suite.Run("err CacheProducts", func() {
		req := params.ListProductsQueryParams{}
		ctx := context.Background()
		suite.MockCacheRepo.EXPECT().GetCachedProducts(ctx, req).Return([]domain.Product{}, redis.Nil)
		suite.MockDbRepo.EXPECT().ListProducts(ctx, req).Return([]domain.Product{}, nil)
		suite.MockCacheRepo.EXPECT().GetCachedProductCount(ctx, req).Return(0, redis.Nil)
		suite.MockDbRepo.EXPECT().CountProducts(ctx, req).Return(0, nil)
		suite.MockCacheRepo.EXPECT().CacheProducts(ctx, req, 0, []domain.Product{}).Return(errors.New("test"))

		got, err := suite.Ps.ListProducts(ctx, req)
		suite.Error(err)
		suite.Nil(got)
	})

	suite.Run("success with using cached data", func() {
		req := params.ListProductsQueryParams{
			Search: "",
			Types:  []string{},
			PaginationParams: params.PaginationParams{
				Sorts: []string{},
				Limit: 2,
				Page:  1,
			},
		}
		ctx := context.Background()

		ListProducts := []domain.Product{
			{
				ID:    1,
				Name:  "milk",
				Price: 20000,
				ProductType: domain.ProductType{
					Name:      "dairy",
					CreatedAt: time.Now(),
				},
				CreatedAt: time.Now(),
			},
		}

		suite.MockCacheRepo.EXPECT().GetCachedProducts(ctx, req).Return(ListProducts, nil)
		suite.MockCacheRepo.EXPECT().GetCachedProductCount(ctx, req).Return(1, nil)

		got, err := suite.Ps.ListProducts(ctx, req)
		suite.NoError(err)
		suite.NotNil(got)
		suite.Equal(got.TotalData, 1)
		suite.Equal(got.TotalPage, 1)
	})

	suite.Run("success without using cached data", func() {
		req := params.ListProductsQueryParams{
			Search: "",
			Types:  []string{},
			PaginationParams: params.PaginationParams{
				Sorts: []string{},
				Limit: 2,
				Page:  1,
			},
		}
		ctx := context.Background()

		ListProducts := []domain.Product{
			{
				ID:    1,
				Name:  "milk",
				Price: 20000,
				ProductType: domain.ProductType{
					Name:      "dairy",
					CreatedAt: time.Now(),
				},
				CreatedAt: time.Now(),
			},
		}

		suite.MockCacheRepo.EXPECT().GetCachedProducts(ctx, req).Return([]domain.Product{}, redis.Nil)
		suite.MockDbRepo.EXPECT().ListProducts(ctx, req).Return(ListProducts, nil)
		suite.MockCacheRepo.EXPECT().GetCachedProductCount(ctx, req).Return(0, redis.Nil)
		suite.MockDbRepo.EXPECT().CountProducts(ctx, req).Return(1, nil)
		suite.MockCacheRepo.EXPECT().CacheProducts(ctx, req, 1, ListProducts).Return(nil)

		got, err := suite.Ps.ListProducts(ctx, req)
		suite.NoError(err)
		suite.NotNil(got)
		suite.Equal(got.TotalData, 1)
		suite.Equal(got.TotalPage, 1)
	})

	suite.Run("success with using cached data", func() {
		req := params.ListProductsQueryParams{
			Search: "",
			Types:  []string{},
			PaginationParams: params.PaginationParams{
				Sorts: []string{},
				Limit: 2,
				Page:  1,
			},
		}
		ctx := context.Background()

		ListProducts := []domain.Product{
			{
				ID:    1,
				Name:  "milk",
				Price: 20000,
				ProductType: domain.ProductType{
					Name:      "dairy",
					CreatedAt: time.Now(),
				},
				CreatedAt: time.Now(),
			},
		}

		suite.MockCacheRepo.EXPECT().GetCachedProducts(ctx, req).Return(ListProducts, nil)
		suite.MockCacheRepo.EXPECT().GetCachedProductCount(ctx, req).Return(0, nil)

		got, err := suite.Ps.ListProducts(ctx, req)
		suite.NoError(err)
		suite.NotNil(got)
		suite.Equal(got.TotalData, 0)
		suite.Equal(got.TotalPage, 0)
	})
}

func (suite *TestProductServiceSuite) TestProductService_CreateProduct() {
	suite.Run("error when getting product", func() {
		req := params.CreateProductRequest{Name: "melon"}
		ctx := context.Background()
		suite.MockDbRepo.EXPECT().CountProducts(ctx, params.ListProductsQueryParams{
			Search: "melon",
		}).Return(0, errors.New("test"))

		_, err := suite.Ps.CreateProduct(ctx, req)
		suite.Error(err)
	})
	suite.Run("product already exist", func() {
		req := params.CreateProductRequest{Name: "melon"}
		ctx := context.Background()
		suite.MockDbRepo.EXPECT().CountProducts(ctx, params.ListProductsQueryParams{
			Search: "melon",
		}).Return(1, nil)

		_, err := suite.Ps.CreateProduct(ctx, req)
		suite.Error(err)
	})

	suite.Run("error when create", func() {
		req := params.CreateProductRequest{Name: "melon"}
		ctx := context.Background()
		suite.MockDbRepo.EXPECT().CountProducts(ctx, params.ListProductsQueryParams{
			Search: "melon",
		}).Return(0, nil)
		suite.MockDbRepo.EXPECT().CreateProduct(ctx, req).Return(0, errors.New("test"))

		_, err := suite.Ps.CreateProduct(ctx, req)
		suite.Error(err)
	})

	suite.Run("error when flush all product", func() {
		req := params.CreateProductRequest{Name: "melon"}
		ctx := context.Background()
		suite.MockDbRepo.EXPECT().CountProducts(ctx, params.ListProductsQueryParams{
			Search: "melon",
		}).Return(0, nil)
		suite.MockDbRepo.EXPECT().CreateProduct(ctx, req).Return(6, nil)
		suite.MockCacheRepo.EXPECT().FlushAllProducts(ctx).Return(errors.New("test"))

		_, err := suite.Ps.CreateProduct(ctx, req)
		suite.Error(err)
	})

	suite.Run("success", func() {
		req := params.CreateProductRequest{Name: "melon"}
		ctx := context.Background()
		suite.MockDbRepo.EXPECT().CountProducts(ctx, params.ListProductsQueryParams{
			Search: "melon",
		}).Return(0, nil)
		suite.MockDbRepo.EXPECT().CreateProduct(ctx, req).Return(6, nil)
		suite.MockCacheRepo.EXPECT().FlushAllProducts(ctx).Return(nil)

		res, err := suite.Ps.CreateProduct(ctx, req)
		suite.NoError(err)
		suite.NotNil(res)
		suite.Equal(res.ID, 6)
	})
}
