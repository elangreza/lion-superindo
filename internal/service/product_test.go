package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/elangreza14/superindo/internal/domain"
	"github.com/elangreza14/superindo/internal/params"
	mockservice "github.com/elangreza14/superindo/mock/service"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type TestProductServiceSuite struct {
	suite.Suite

	MockProductRepo *mockservice.MockProductRepo
	Ps              *ProductService
	Ctrl            *gomock.Controller
}

func (suite *TestProductServiceSuite) SetupSuite() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockProductRepo = mockservice.NewMockProductRepo(suite.Ctrl)
	suite.Ps = NewProductService(suite.MockProductRepo)
}

func (suite *TestProductServiceSuite) TearDownSuite() {
	suite.Ctrl.Finish()
}

func TestProductServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TestProductServiceSuite))
}

func (suite *TestProductServiceSuite) TestProductService_GetProducts() {
	suite.Run("error when getting total of product", func() {
		req := params.ListProductQueryParams{}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, req, true).Return(0, errors.New("test"))

		got, err := suite.Ps.ListProduct(ctx, req)
		suite.Error(err)
		suite.Nil(got)
	})

	suite.Run("got total products with 0 data", func() {
		req := params.ListProductQueryParams{}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, req, true).Return(0, nil)

		got, err := suite.Ps.ListProduct(ctx, req)
		suite.NoError(err)
		suite.NotNil(got)
		suite.Equal(int(0), got.TotalData)
		suite.Equal(int(0), got.TotalPage)
		suite.Equal(int(0), len(got.Products))
	})

	suite.Run("error when getting list of product", func() {
		req := params.ListProductQueryParams{}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, req, true).Return(1, nil)
		suite.MockProductRepo.EXPECT().ListProduct(ctx, req).Return(nil, errors.New("test"))

		got, err := suite.Ps.ListProduct(ctx, req)
		suite.Error(err)
		suite.Nil(got)
	})

	suite.Run("error when getting list of product", func() {
		req := params.ListProductQueryParams{}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, req, true).Return(1, nil)
		suite.MockProductRepo.EXPECT().ListProduct(ctx, req).Return(nil, errors.New("test"))

		got, err := suite.Ps.ListProduct(ctx, req)
		suite.Error(err)
		suite.Nil(got)
	})

	suite.Run("success with limit 2 and total products is 3", func() {
		req := params.ListProductQueryParams{Limit: 2}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, req, true).Return(3, nil)
		suite.MockProductRepo.EXPECT().ListProduct(ctx, req).Return([]domain.Product{
			{
				ID:       1,
				BaseDate: domain.BaseDate{CreatedAt: time.Now()},
			}, {
				ID:       2,
				BaseDate: domain.BaseDate{CreatedAt: time.Now()},
			}, {
				ID:       3,
				BaseDate: domain.BaseDate{CreatedAt: time.Now()},
			},
		}, nil)

		got, err := suite.Ps.ListProduct(ctx, req)
		suite.NoError(err)
		suite.NotNil(got)
		suite.Equal(got.TotalData, 3)
		suite.Equal(got.TotalPage, 2)
	})

	suite.Run("success with limit 2 and total products is 2", func() {
		updatedAt := time.Now()
		req := params.ListProductQueryParams{Limit: 2}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, req, true).Return(2, nil)
		suite.MockProductRepo.EXPECT().ListProduct(ctx, req).Return([]domain.Product{
			{
				ID:       1,
				BaseDate: domain.BaseDate{CreatedAt: time.Now()},
			}, {
				ID:       2,
				BaseDate: domain.BaseDate{CreatedAt: time.Now(), UpdatedAt: &updatedAt},
			},
		}, nil)

		got, err := suite.Ps.ListProduct(ctx, req)
		suite.NoError(err)
		suite.NotNil(got)
		suite.Equal(got.TotalData, 2)
		suite.Equal(got.TotalPage, 1)
	})
}

func (suite *TestProductServiceSuite) TestProductService_CreateProduct() {
	suite.Run("error when getting product", func() {
		req := params.CreateProductRequest{Name: "melon"}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, params.ListProductQueryParams{
			Search: "melon",
		}, false).Return(0, errors.New("test"))

		_, err := suite.Ps.CreateProduct(ctx, req)
		suite.Error(err)
	})
	suite.Run("product already exist", func() {
		req := params.CreateProductRequest{Name: "melon"}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, params.ListProductQueryParams{
			Search: "melon",
		}, false).Return(1, nil)

		_, err := suite.Ps.CreateProduct(ctx, req)
		suite.Error(err)
	})

	suite.Run("error when create", func() {
		req := params.CreateProductRequest{Name: "melon"}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, params.ListProductQueryParams{
			Search: "melon",
		}, false).Return(0, nil)
		suite.MockProductRepo.EXPECT().CreateProduct(ctx, req).Return(0, errors.New("test"))

		_, err := suite.Ps.CreateProduct(ctx, req)
		suite.Error(err)
	})

	suite.Run("create product", func() {
		req := params.CreateProductRequest{Name: "melon"}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, params.ListProductQueryParams{
			Search: "melon",
		}, false).Return(0, nil)
		suite.MockProductRepo.EXPECT().CreateProduct(ctx, req).Return(6, nil)

		res, err := suite.Ps.CreateProduct(ctx, req)
		suite.NoError(err)
		suite.NotNil(res)
		suite.Equal(res.ID, 6)
	})

}
