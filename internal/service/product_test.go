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
	QueryParam      params.ProductQueryParams
	Ps              *ProductService
	Ctrl            *gomock.Controller
}

func (suite *TestProductServiceSuite) SetupSuite() {
	suite.Ctrl = gomock.NewController(suite.T())
	suite.MockProductRepo = mockservice.NewMockProductRepo(suite.Ctrl)
	suite.QueryParam = params.ProductQueryParams{}
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
		req := params.ProductQueryParams{}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, req).Return(0, errors.New("test"))

		got, err := suite.Ps.ListProduct(ctx, req)
		suite.Error(err)
		suite.Nil(got)
	})

	suite.Run("got total products with 0 data", func() {
		req := params.ProductQueryParams{}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, req).Return(0, nil)

		got, err := suite.Ps.ListProduct(ctx, req)
		suite.NoError(err)
		suite.NotNil(got)
		suite.Equal(int(0), got.TotalData)
		suite.Equal(int(0), got.TotalPage)
		suite.Equal(int(0), len(got.Products))
	})

	suite.Run("error when getting list of product", func() {
		req := params.ProductQueryParams{}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, req).Return(1, nil)
		suite.MockProductRepo.EXPECT().ListProduct(ctx, req).Return(nil, errors.New("test"))

		got, err := suite.Ps.ListProduct(ctx, req)
		suite.Error(err)
		suite.Nil(got)
	})

	suite.Run("error when getting list of product", func() {
		req := params.ProductQueryParams{}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, req).Return(1, nil)
		suite.MockProductRepo.EXPECT().ListProduct(ctx, req).Return(nil, errors.New("test"))

		got, err := suite.Ps.ListProduct(ctx, req)
		suite.Error(err)
		suite.Nil(got)
	})

	suite.Run("success with limit 2 and total products is 3", func() {
		req := params.ProductQueryParams{Limit: 2}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, req).Return(3, nil)
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
		req := params.ProductQueryParams{Limit: 2}
		ctx := context.Background()
		suite.MockProductRepo.EXPECT().TotalProduct(ctx, req).Return(2, nil)
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
