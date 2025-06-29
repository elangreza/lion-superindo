//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"
	"net/http"

	"github.com/elangreza14/lion-superindo/cmd/server/config"
	"github.com/elangreza14/lion-superindo/internal/handler"
	postgreRepo "github.com/elangreza14/lion-superindo/internal/postgresql"
	redisRepo "github.com/elangreza14/lion-superindo/internal/redis"
	"github.com/elangreza14/lion-superindo/internal/service"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

type ProductHandlerDeps struct {
	Mux         *http.ServeMux
	DB          *sql.DB
	RedisClient *redis.Client
}

var productSet = wire.NewSet(
	config.SetupDB,
	config.SetupCache,
	postgreRepo.NewRepo,
	wire.Bind(new(service.DbRepo), new(*postgreRepo.PostgresRepo)), // <-- Bind DbRepo interface
	redisRepo.NewRepo,
	wire.Bind(new(service.CacheRepo), new(*redisRepo.RedisRepo)), // <-- Bind CacheRepo interface
	service.NewProductService,
	wire.Bind(new(handler.ProductService), new(*service.ProductService)), // <-- This line binds interface to implementation
	handler.NewProductHandler,
	handler.NewRoutes,
)

func InitializeProductHandler(cfg *config.Config) (*ProductHandlerDeps, error) {
	wire.Build(
		productSet,
		wire.Struct(new(ProductHandlerDeps), "Mux", "DB", "RedisClient"),
	)
	return nil, nil
}
