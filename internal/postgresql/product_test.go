package postgresql

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elangreza14/lion-superindo/internal/params"
	"github.com/stretchr/testify/assert"
)

func TestProductRepo(t *testing.T) {
	db, mockSql, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	pr := NewRepo(db)
	now := time.Now()

	testTable := []struct {
		name        string
		expectedErr bool
		mock        func(sqlmock.Sqlmock)
		reqParams   params.ListProductsQueryParams
	}{
		{
			name:        "success with full params",
			expectedErr: false,
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows([]string{"id", "name", "price", "product_type_name", "created_at"}).
					AddRow(1, "test", 1, "test", now)
				m.ExpectQuery("SELECT (.+) FROM products").WillReturnRows(rows)
			},
			reqParams: params.ListProductsQueryParams{
				PaginationParams: params.PaginationParams{
					Limit: 1,
					Page:  2,
					Sorts: []string{"id:desc"},
				},
				Search: "milk",
				Types:  []string{"dairy"},
			},
		},
		{
			name:        "error scanning from db",
			expectedErr: true,
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows([]string{"id", "name", "price", "product_type_name", "created_at"}).
					AddRow("a", "test", 1, "test", now)
				m.ExpectQuery("SELECT (.+) FROM products").WillReturnRows(rows)
			},
			reqParams: params.ListProductsQueryParams{
				PaginationParams: params.PaginationParams{
					Limit: 1,
					Page:  2,
					Sorts: []string{"id:desc"},
				},
				Search: "milk",
				Types:  []string{"dairy"},
			},
		},
		{
			name:        "failed with id and got err",
			expectedErr: true,
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT (.+) FROM products").WillReturnError(errors.New("test"))
			},
			reqParams: params.ListProductsQueryParams{
				PaginationParams: params.PaginationParams{
					Limit: 1,
					Page:  2,
				},
				Search: "1",
				Types:  []string{"dairy"},
			},
		},
		{
			name:        "failed rows",
			expectedErr: true,
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT (.+) FROM products").WillReturnRows(
					sqlmock.NewRows(
						[]string{},
					).CloseError(errors.New("test")))
			},
			reqParams: params.ListProductsQueryParams{
				PaginationParams: params.PaginationParams{
					Limit: 1,
					Page:  2,
				},
				Search: "1",
				Types:  []string{"dairy"},
			},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {

			test.mock(mockSql)
			test.reqParams.Validate()
			got, err := pr.ListProducts(context.Background(), test.reqParams)

			if test.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}

			if err := mockSql.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}

}

func TestProductRepo_CountProducts(t *testing.T) {
	db, mockSql, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	pr := NewRepo(db)

	testTable := []struct {
		name        string
		expectedErr bool
		mock        func(sqlmock.Sqlmock)
		reqParams   params.ListProductsQueryParams
		got         int
	}{
		{
			name:        "success",
			expectedErr: false,
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				m.ExpectQuery(regexp.QuoteMeta("SELECT")).
					WillReturnRows(rows)
			},
			reqParams: params.ListProductsQueryParams{
				Search: "a",
				PaginationParams: params.PaginationParams{
					Limit: 1,
				},
			},
			got: 1,
		},
		{
			name:        "failed",
			expectedErr: true,
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta("SELECT")).
					WillReturnError(errors.New("test"))
			},
			reqParams: params.ListProductsQueryParams{
				Search: "a",
				PaginationParams: params.PaginationParams{
					Limit: 1,
				},
			},
			got: 0,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			test.mock(mockSql)
			test.reqParams.Validate()
			got, err := pr.CountProducts(context.Background(), test.reqParams)
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.got, got)

			if err := mockSql.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestProductRepo_CreateProduct(t *testing.T) {
	dbSql, mockSql, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer dbSql.Close()
	pr := &PostgresRepo{dbSql}

	testTable := []struct {
		name        string
		expectedErr bool
		mock        func(sqlmock.Sqlmock)
		reqParams   params.CreateProductRequest
		got         int
	}{
		{
			name:        "success",
			expectedErr: false,
			mock: func(m sqlmock.Sqlmock) {
				mockSql.ExpectBegin()
				mockSql.ExpectExec("INSERT INTO product_types").WithArgs("buah").WillReturnResult(sqlmock.NewResult(1, 1))
				mockSql.ExpectQuery("INSERT INTO products").WithArgs("melon", 1000, "buah").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))
				mockSql.ExpectCommit()
			},
			reqParams: params.CreateProductRequest{
				Name:  "melon",
				Price: 1000,
				Type:  "buah",
			},
			got: 5,
		},
		{
			name:        "error INSERT INTO products",
			expectedErr: true,
			mock: func(m sqlmock.Sqlmock) {
				mockSql.ExpectBegin()
				mockSql.ExpectExec("INSERT INTO product_types").WithArgs("buah").WillReturnResult(sqlmock.NewResult(1, 1))
				mockSql.ExpectQuery("INSERT INTO products").WithArgs("melon", 1000, "buah").WillReturnError(errors.New("test"))
				mockSql.ExpectRollback()
			},
			reqParams: params.CreateProductRequest{
				Name:  "melon",
				Price: 1000,
				Type:  "buah",
			},
			got: 0,
		},
		{
			name:        "error INSERT INTO product_types",
			expectedErr: true,
			mock: func(m sqlmock.Sqlmock) {
				mockSql.ExpectBegin()
				mockSql.ExpectExec("INSERT INTO product_types").WithArgs("buah").WillReturnError(errors.New("test"))
				mockSql.ExpectRollback()
			},
			reqParams: params.CreateProductRequest{
				Name:  "melon",
				Price: 1000,
				Type:  "buah",
			},
			got: 0,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			test.mock(mockSql)
			id, err := pr.CreateProduct(context.Background(), test.reqParams)
			if test.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.got, id)

			if err := mockSql.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
