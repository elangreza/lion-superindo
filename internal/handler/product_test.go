package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/elangreza/lion-superindo/internal/params"
	mockhandler "github.com/elangreza/lion-superindo/mock/handler"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var mockErrorResBody = struct {
	Error string `json:"error"`
}{}

func TestProductHandler_ProductHandler_Invalid_Method(t *testing.T) {
	r := httptest.NewRequest(http.MethodPatch, "/product", nil)
	w := httptest.NewRecorder()

	mc := gomock.NewController(t)
	mockProductService := mockhandler.NewMockProductService(mc)
	ph := NewProductHandler(mockProductService)
	routes := NewRoutes(ph)
	routes.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)

	resBody := mockErrorResBody
	err = json.Unmarshal(body, &resBody)
	assert.NoError(t, err)
	assert.Equal(t, resBody.Error, "method PATCH not allowed")
}

func TestProductHandler_ListProductsHandler_Error_When_Validate_Query_Params(t *testing.T) {
	mc := gomock.NewController(t)
	mockProductService := mockhandler.NewMockProductService(mc)
	ph := NewProductHandler(mockProductService)
	routes := NewRoutes(ph)

	r := httptest.NewRequest(http.MethodGet, "/product?sort=test", nil)
	w := httptest.NewRecorder()
	routes.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	resBody := mockErrorResBody
	err = json.Unmarshal(body, &resBody)
	assert.NoError(t, err)
	assert.Equal(t, resBody.Error, "validation error: test is not valid sort format")
}

func TestProductHandler_ListProductsHandler_Error_When_Processing_ListProducts(t *testing.T) {
	mc := gomock.NewController(t)
	mockProductService := mockhandler.NewMockProductService(mc)
	ph := NewProductHandler(mockProductService)
	routes := NewRoutes(ph)
	mockProductService.EXPECT().ListProducts(gomock.Any(), gomock.Any()).Return(nil, errors.New("test"))

	r := httptest.NewRequest(http.MethodGet, "/product", nil)
	w := httptest.NewRecorder()
	routes.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	resBody := mockErrorResBody
	err = json.Unmarshal(body, &resBody)
	assert.NoError(t, err)
	assert.Equal(t, resBody.Error, "server error")
}

func TestProductHandler_ListProductsHandler_Success(t *testing.T) {
	mc := gomock.NewController(t)
	mockProductService := mockhandler.NewMockProductService(mc)
	ph := NewProductHandler(mockProductService)
	routes := NewRoutes(ph)
	resMock := &params.ListProductsResponses{
		TotalData: 1,
		TotalPage: 1,
		Products: []params.ProductResponse{
			{
				ID:        1,
				Name:      "semangka",
				Price:     1,
				Type:      "buah",
				CreatedAt: time.Now(),
			},
		},
	}
	mockProductService.EXPECT().ListProducts(gomock.Any(), gomock.Any()).Return(resMock, nil)

	r := httptest.NewRequest(http.MethodGet, "/product", nil)
	w := httptest.NewRecorder()
	routes.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody := struct {
		Data params.ListProductsResponses `json:"data"`
	}{}
	err = json.Unmarshal(body, &resBody)
	assert.NoError(t, err)
	assert.Equal(t, resBody.Data.TotalData, 1)
	assert.Equal(t, resBody.Data.TotalPage, 1)
	assert.Equal(t, len(resBody.Data.Products), 1)
}

func TestProductHandler_CreateProductHandler_Error_When_Validate_Query(t *testing.T) {
	mc := gomock.NewController(t)
	mockProductService := mockhandler.NewMockProductService(mc)
	ph := NewProductHandler(mockProductService)
	routes := NewRoutes(ph)

	reqBody := params.CreateProductRequest{
		Name:  "a",
		Price: -1,
		Type:  "a",
	}
	errPayload, _ := json.Marshal(reqBody)
	bodyReader := bytes.NewReader(errPayload)

	r := httptest.NewRequest(http.MethodPost, "/product", bodyReader)
	w := httptest.NewRecorder()
	routes.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	resBody := mockErrorResBody
	err = json.Unmarshal(body, &resBody)
	assert.NoError(t, err)
	assert.Equal(t, resBody.Error, "validation error: price cannot be negative")
}

func TestProductHandler_CreateProductHandler_Error_When_Processing_CreateProduct(t *testing.T) {
	mc := gomock.NewController(t)
	mockProductService := mockhandler.NewMockProductService(mc)
	ph := NewProductHandler(mockProductService)
	routes := NewRoutes(ph)

	reqBody := params.CreateProductRequest{
		Name:  "a",
		Price: 1,
		Type:  "a",
	}
	errPayload, _ := json.Marshal(reqBody)
	bodyReader := bytes.NewReader(errPayload)

	mockProductService.EXPECT().CreateProduct(gomock.Any(), gomock.Any()).Return(nil, errors.New("test"))

	r := httptest.NewRequest(http.MethodPost, "/product", bodyReader)
	w := httptest.NewRecorder()
	routes.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	resBody := mockErrorResBody
	err = json.Unmarshal(body, &resBody)
	assert.NoError(t, err)
	assert.Equal(t, resBody.Error, "server error")
}

func TestProductHandler_CreateProductHandler_Success(t *testing.T) {
	mc := gomock.NewController(t)
	mockProductService := mockhandler.NewMockProductService(mc)
	ph := NewProductHandler(mockProductService)
	routes := NewRoutes(ph)

	reqBody := params.CreateProductRequest{
		Name:  "a",
		Price: 1,
		Type:  "a",
	}
	errPayload, _ := json.Marshal(reqBody)
	bodyReader := bytes.NewReader(errPayload)

	mockProductService.EXPECT().CreateProduct(gomock.Any(), gomock.Any()).Return(&params.CreateProductResponse{
		ID: 1,
	}, nil)

	r := httptest.NewRequest(http.MethodPost, "/product", bodyReader)
	w := httptest.NewRecorder()
	routes.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	resBodya := struct {
		Data params.CreateProductResponse `json:"data"`
	}{}
	err = json.Unmarshal(body, &resBodya)
	assert.NoError(t, err)
	assert.Equal(t, resBodya.Data.ID, int(1))

}
