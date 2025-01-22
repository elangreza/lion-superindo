package postgresql

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elangreza14/superindo/internal/params"
	"github.com/stretchr/testify/assert"
)

func TestProductRepo_ListProduct(t *testing.T) {
	// mc := gomock.NewController(t)
	// mockCache := mockpostgresql.NewMockCache(mc)
	db, mockSql, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	pr := &ProductRepo{db, nil}

	// mockCache.EXPECT().Set("a", 1).Return(errors.New("cek"))

	rows := sqlmock.
		NewRows([]string{"id", "name", "quantity", "price", "product_type_name", "created_at", "updated_at"}).
		AddRow(1, "test", 1, 1, "test", time.Now(), nil)

	mockSql.
		ExpectQuery("SELECT (.+) FROM products").
		WillReturnRows(rows)

	got, err := pr.ListProduct(context.Background(), params.ListProductQueryParams{Limit: 1})
	assert.NotNil(t, pr)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	if err := mockSql.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductRepo_TotalProduct(t *testing.T) {
	// mc := gomock.NewController(t)
	// mockCache := mockpostgresql.NewMockCache(mc)
	db, mockSql, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	pr := &ProductRepo{db, nil}

	// mockCache.EXPECT().Set("a", 1).Return(errors.New("cek"))

	rows := sqlmock.
		NewRows([]string{"count(id)"}).
		AddRow(1)

	mockSql.
		ExpectQuery("SELECT (.+) FROM products").
		WillReturnRows(rows)

	got, err := pr.TotalProduct(context.Background(), params.ListProductQueryParams{Limit: 1})
	assert.NotNil(t, pr)
	assert.NoError(t, err)
	assert.Equal(t, got, int(1))
	if err := mockSql.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductRepo_CreateOrUpdateProduct(t *testing.T) {
	db, mockSql, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
	pr := &ProductRepo{db, nil}

	mockSql.ExpectBegin()
	mockSql.ExpectExec("INSERT INTO product_types").WithArgs("buah").WillReturnResult(sqlmock.NewResult(1, 1))
	mockSql.ExpectExec("INSERT INTO products").WithArgs("melon", 1, 1000, "buah").WillReturnResult(sqlmock.NewResult(1, 1))
	mockSql.ExpectCommit()

	err = pr.CreateOrUpdateProduct(context.Background(), params.CreateOrUpdateProductRequest{
		Name:     "melon",
		Quantity: 1,
		Price:    1000,
		Type:     "buah",
	})
	assert.NoError(t, err)
	if err := mockSql.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
