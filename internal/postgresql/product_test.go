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

	got, err := pr.ListProduct(context.Background(), params.ProductQueryParams{Limit: 1})
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

	got, err := pr.TotalProduct(context.Background(), params.ProductQueryParams{Limit: 1})
	assert.NotNil(t, pr)
	assert.NoError(t, err)
	assert.Equal(t, got, int(1))
	if err := mockSql.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
