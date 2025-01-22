package main

import (
	"database/sql"
	"fmt"

	"github.com/elangreza14/superindo/cmd/rest/config"
	"github.com/elangreza14/superindo/internal/handler"
	"github.com/elangreza14/superindo/internal/postgresql"
	"github.com/elangreza14/superindo/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	errChecker(err)

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.POSTGRES_USER,
		cfg.POSTGRES_PASSWORD,
		cfg.POSTGRES_HOSTNAME,
		cfg.POSTGRES_PORT,
		cfg.POSTGRES_DB,
		cfg.POSTGRES_SSL,
	)

	pool, err := sql.Open("postgres", connString)
	errChecker(err)
	defer pool.Close()

	productRepo := postgresql.NewProductRepo(pool, nil)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	srv := gin.New()
	srv.GET("/products", productHandler.ListProductHandler())
	srv.POST("/products", productHandler.CreateProductHandler())
	srv.Run(cfg.HTTP_PORT)
}

func errChecker(err error) {
	if err != nil {
		panic(err)
	}
}
