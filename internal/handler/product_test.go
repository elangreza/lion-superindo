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

	"github.com/elangreza14/superindo/internal/params"
	mockhandler "github.com/elangreza14/superindo/mock/handler"
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
	assert.Equal(t, resBody.Error, "invalid method")
}

func TestProductHandler_ListProductHandler_Error_When_Validate_Query_Params(t *testing.T) {
	mc := gomock.NewController(t)
	mockProductService := mockhandler.NewMockProductService(mc)
	ph := NewProductHandler(mockProductService)
	routes := NewRoutes(ph)

	r := httptest.NewRequest(http.MethodGet, "/product?sorts=test", nil)
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
	assert.Equal(t, resBody.Error, "not valid sort format")
}

func TestProductHandler_ListProductHandler_Error_When_Processing_ListProduct(t *testing.T) {
	mc := gomock.NewController(t)
	mockProductService := mockhandler.NewMockProductService(mc)
	ph := NewProductHandler(mockProductService)
	routes := NewRoutes(ph)
	mockProductService.EXPECT().ListProduct(gomock.Any(), gomock.Any()).Return(nil, errors.New("test"))

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

func TestProductHandler_ListProductHandler_Success(t *testing.T) {
	mc := gomock.NewController(t)
	mockProductService := mockhandler.NewMockProductService(mc)
	ph := NewProductHandler(mockProductService)
	routes := NewRoutes(ph)
	resMock := &params.ListProductResponses{
		TotalData: 1,
		TotalPage: 1,
		Products: []params.ProductResponse{
			{
				ID:       1,
				Name:     "semangka",
				Quantity: 1,
				Price:    1,
				Type:     "buah",
				UpdateAt: time.Now(),
			},
		},
	}
	mockProductService.EXPECT().ListProduct(gomock.Any(), gomock.Any()).Return(resMock, nil)

	r := httptest.NewRequest(http.MethodGet, "/product", nil)
	w := httptest.NewRecorder()
	routes.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	resBody := struct {
		Data params.ListProductResponses `json:"data"`
	}{}
	err = json.Unmarshal(body, &resBody)
	assert.NoError(t, err)
	assert.Equal(t, resBody.Data.TotalData, 1)
	assert.Equal(t, resBody.Data.TotalPage, 1)
	assert.Equal(t, len(resBody.Data.Products), 1)
}

func TestProductHandler_CreateOrUpdateProductHandler(t *testing.T) {
	mc := gomock.NewController(t)
	mockProductService := mockhandler.NewMockProductService(mc)
	ph := NewProductHandler(mockProductService)
	routes := NewRoutes(ph)

	reqBody := params.CreateOrUpdateProductRequest{
		Name:     "a",
		Quantity: -1,
		Price:    1,
		Type:     "a",
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
	assert.Equal(t, resBody.Error, "quantity cannot be negative")
}

func TestProductHandler_CreateOrUpdateProductHandler_Error_When_Processing_CreateOrUpdateProduct(t *testing.T) {
	mc := gomock.NewController(t)
	mockProductService := mockhandler.NewMockProductService(mc)
	ph := NewProductHandler(mockProductService)
	routes := NewRoutes(ph)

	reqBody := params.CreateOrUpdateProductRequest{
		Name:     "a",
		Quantity: 1,
		Price:    1,
		Type:     "a",
	}
	errPayload, _ := json.Marshal(reqBody)
	bodyReader := bytes.NewReader(errPayload)

	mockProductService.EXPECT().CreateOrUpdateProduct(gomock.Any(), gomock.Any()).Return(errors.New("test"))

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

func TestProductHandler_CreateOrUpdateProductHandler_Success(t *testing.T) {
	mc := gomock.NewController(t)
	mockProductService := mockhandler.NewMockProductService(mc)
	ph := NewProductHandler(mockProductService)
	routes := NewRoutes(ph)

	reqBody := params.CreateOrUpdateProductRequest{
		Name:     "a",
		Quantity: 1,
		Price:    1,
		Type:     "a",
	}
	errPayload, _ := json.Marshal(reqBody)
	bodyReader := bytes.NewReader(errPayload)

	mockProductService.EXPECT().CreateOrUpdateProduct(gomock.Any(), gomock.Any()).Return(nil)

	r := httptest.NewRequest(http.MethodPost, "/product", bodyReader)
	w := httptest.NewRecorder()
	routes.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

}
