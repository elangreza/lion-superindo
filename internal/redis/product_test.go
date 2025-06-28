package redis

// import (
// 	"context"
// 	"encoding/json"
// 	"testing"
// 	"time"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/elangreza14/lion-superindo/internal/domain"
// 	"github.com/elangreza14/lion-superindo/internal/params"
// 	"github.com/go-redis/redismock/v9"
// 	"github.com/stretchr/testify/assert"
// )

// // func TestProductRepo_ListProduct(t *testing.T) {
// // 	db, mockSql, err := sqlmock.New()
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	defer db.Close()
// // 	dbRedis, mockRedis := redismock.NewClientMock()
// // 	pr := &ProductRepo{db, dbRedis}
// // 	now := time.Now()
// // 	products := []domain.Product{
// // 		{
// // 			ID:    1,
// // 			Name:  "test",
// // 			Price: 1,
// // 			ProductType: domain.ProductTypes{
// // 				Name:     "test",
// // 				BaseDate: domain.BaseDate{},
// // 			},
// // 			BaseDate: domain.BaseDate{
// // 				CreatedAt: now,
// // 				UpdatedAt: nil,
// // 			},
// // 		},
// // 	}
// // 	productJson, _ := json.Marshal(products)

// // 	mockRedis.ExpectGet("listProduct:").SetVal(string(productJson))

// // 	got, err := pr.ListProduct(context.Background(), params.ListProductQueryParams{Limit: 1})
// // 	assert.NotNil(t, pr)
// // 	assert.NoError(t, err)
// // 	assert.NotNil(t, got)
// // 	if err := mockSql.ExpectationsWereMet(); err != nil {
// // 		t.Errorf("there were unfulfilled expectations: %s", err)
// // 	}
// // 	if err := mockRedis.ExpectationsWereMet(); err != nil {
// // 		t.Errorf("there were unfulfilled expectations: %s", err)
// // 	}
// // }

// // func TestProductRepo_ListProduct_With_Cache_Is_Exist(t *testing.T) {
// // 	db, mockSql, err := sqlmock.New()
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	defer db.Close()
// // 	dbRedis, mockRedis := redismock.NewClientMock()
// // 	pr := &ProductRepo{db, dbRedis}
// // 	now := time.Now()

// // 	products := []domain.Product{
// // 		{
// // 			ID:    1,
// // 			Name:  "test",
// // 			Price: 1,
// // 			ProductType: domain.ProductTypes{
// // 				Name:     "test",
// // 				BaseDate: domain.BaseDate{},
// // 			},
// // 			BaseDate: domain.BaseDate{
// // 				CreatedAt: now,
// // 				UpdatedAt: nil,
// // 			},
// // 		},
// // 	}
// // 	productJson, _ := json.Marshal(products)

// // 	mockRedis.ExpectGet("listProduct:").RedisNil()
// // 	rows := sqlmock.
// // 		NewRows([]string{"id", "name", "price", "product_type_name", "created_at", "created_at"}).
// // 		AddRow(1, "test", 1, "test", now, nil)
// // 	mockSql.ExpectQuery("SELECT (.+) FROM products").WillReturnRows(rows)
// // 	mockRedis.ExpectSet("listProduct:", string(productJson), time.Second*60).SetVal(string(productJson))

// // 	got, err := pr.ListProduct(context.Background(), params.ListProductQueryParams{Limit: 1})
// // 	assert.NotNil(t, pr)
// // 	assert.NoError(t, err)
// // 	assert.NotNil(t, got)
// // 	if err := mockSql.ExpectationsWereMet(); err != nil {
// // 		t.Errorf("there were unfulfilled expectations: %s", err)
// // 	}
// // 	if err := mockRedis.ExpectationsWereMet(); err != nil {
// // 		t.Errorf("there were unfulfilled expectations: %s", err)
// // 	}
// // }

// // func TestProductRepo_TotalProduct_With_Cache_Is_Exist(t *testing.T) {
// // 	db, mockSql, err := sqlmock.New()
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	defer db.Close()
// // 	dbRedis, mockRedis := redismock.NewClientMock()
// // 	pr := &ProductRepo{db, dbRedis}

// // 	mockRedis.ExpectGet("totalProduct:").SetVal("1")

// // 	got, err := pr.TotalProduct(context.Background(), params.ListProductQueryParams{Limit: 1}, true)
// // 	assert.NotNil(t, pr)
// // 	assert.NoError(t, err)
// // 	assert.Equal(t, got, int(1))
// // 	if err := mockSql.ExpectationsWereMet(); err != nil {
// // 		t.Errorf("there were unfulfilled expectations: %s", err)
// // 	}
// // 	if err := mockRedis.ExpectationsWereMet(); err != nil {
// // 		t.Errorf("there were unfulfilled expectations: %s", err)
// // 	}
// // }

// // func TestProductRepo_TotalProduct(t *testing.T) {
// // 	db, mockSql, err := sqlmock.New()
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	defer db.Close()
// // 	dbRedis, mockRedis := redismock.NewClientMock()
// // 	pr := &ProductRepo{db, dbRedis}

// // 	mockRedis.ExpectGet("totalProduct:").RedisNil()
// // 	rows := sqlmock.NewRows([]string{"count(id)"}).AddRow(1)
// // 	mockSql.ExpectQuery("SELECT (.+) FROM products").WillReturnRows(rows)
// // 	mockRedis.ExpectSet("totalProduct:", 1, time.Second*60).SetVal("1")

// // 	got, err := pr.TotalProduct(context.Background(), params.ListProductQueryParams{Limit: 1}, true)
// // 	assert.NotNil(t, pr)
// // 	assert.NoError(t, err)
// // 	assert.Equal(t, got, int(1))
// // 	if err := mockSql.ExpectationsWereMet(); err != nil {
// // 		t.Errorf("there were unfulfilled expectations: %s", err)
// // 	}
// // 	if err := mockRedis.ExpectationsWereMet(); err != nil {
// // 		t.Errorf("there were unfulfilled expectations: %s", err)
// // 	}
// // }

// // func TestProductRepo_CreateProduct(t *testing.T) {
// // 	dbSql, mockSql, err := sqlmock.New()
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	defer dbSql.Close()
// // 	dbRedis, mockRedis := redismock.NewClientMock()
// // 	pr := &ProductRepo{dbSql, dbRedis}

// // 	mockSql.ExpectBegin()
// // 	mockSql.ExpectExec("INSERT INTO product_types").WithArgs("buah").WillReturnResult(sqlmock.NewResult(1, 1))
// // 	mockSql.ExpectQuery("INSERT INTO products").WithArgs("melon", 1000, "buah").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))
// // 	mockSql.ExpectCommit()
// // 	mockRedis.ExpectFlushAll().SetVal("")

// // 	id, err := pr.CreateProduct(context.Background(), params.CreateProductRequest{
// // 		Name:  "melon",
// // 		Price: 1000,
// // 		Type:  "buah",
// // 	})
// // 	assert.NoError(t, err)
// // 	if err := mockSql.ExpectationsWereMet(); err != nil {
// // 		t.Errorf("there were unfulfilled expectations: %s", err)
// // 	}
// // 	assert.Equal(t, id, 5)
// // }
