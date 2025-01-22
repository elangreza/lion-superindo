package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/elangreza14/superindo/internal/params"
	"github.com/gin-gonic/gin"
)

type ProductService interface {
	ListProduct(ctx context.Context, args params.ListProductQueryParams) (*params.ListProductResponses, error)
	CreateOrUpdateProduct(ctx context.Context, req params.CreateOrUpdateProductRequest) error
}

// ProductHandler ...
type ProductHandler struct {
	svc ProductService
}

// NewProductHandler ...
func NewProductHandler(svc ProductService) *ProductHandler {
	return &ProductHandler{
		svc: svc,
	}
}

func (ph *ProductHandler) ListProductHandler(c *gin.Context) {
	query := params.ListProductQueryParams{}
	err := c.BindQuery(&query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = query.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	res, err := ph.svc.ListProduct(c, query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func (ph *ProductHandler) CreateOrUpdateProductHandler(c *gin.Context) {
	body := params.CreateOrUpdateProductRequest{}
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = body.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := ph.svc.CreateOrUpdateProduct(c, body); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (ph *ProductHandler) ProductHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPost:
			ph.CreateOrUpdateProductHandler(c)
		case http.MethodGet:
			ph.ListProductHandler(c)
		default:
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, nil)
		}
	}
}

// func newHandler(w http.ResponseWriter, r *http.Request) {
//   fmt.Println("GET params were:", r.URL.Query())

//   // if only one expected
//   param1 := r.URL.Query().Get("param1")
//   if param1 != "" {
//     // ... process it, will be the first (only) if multiple were given
//     // note: if they pass in like ?param1=&param2= param1 will also be "" :|
//   }

//   // if multiples possible, or to process empty values like param1 in
//   // ?param1=&param2=something
//   param1s := r.URL.Query()["param1"]
//   if len(param1s) > 0 {
//     // ... process them ... or you could just iterate over them without a check
//     // this way you can also tell if they passed in the parameter as the empty string
//     // it will be an element of the array that is the empty string
//   }
// }

func (ph *ProductHandler) ProductHandler2(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		query := params.ListProductQueryParams{}
		query.Search = r.URL.Query().Get("search")
		query.Types = r.URL.Query()["types"]
		query.Sorts = r.URL.Query()["sorts"]
		if err := query.Validate(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		res, err := ph.svc.ListProduct(r.Context(), query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	case http.MethodPost:
		body := params.CreateOrUpdateProductRequest{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := body.Validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		if err := ph.svc.CreateOrUpdateProduct(r.Context(), body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(nil)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
	}
}

// func TestController_CreateShortUrl(t *testing.T) {
// 	domain.NewSqId()
// 	r := httptest.NewRequest(http.MethodPost, "/l", nil)
// 	w := httptest.NewRecorder()

// 	mc := gomock.NewController(t)
// 	mockSortUrlService := mockrest.NewMockshortUrlService(mc)
// 	routes := NewRoute(RouteParams{
// 		shortUrlService: mockSortUrlService,
// 	})
// 	routes.ServeHTTP(w, r)

// 	res := w.Result()
// 	defer res.Body.Close()
// 	bytes, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		t.Errorf("expected error to be nil got %v", err)
// 	}

// 	fmt.Println(string(bytes))
// 	fmt.Println(res.StatusCode)
// }
