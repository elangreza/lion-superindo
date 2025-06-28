//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/elangreza14/lion-superindo/cmd/server/config"
	"github.com/elangreza14/lion-superindo/internal/handler"
	postgreRepo "github.com/elangreza14/lion-superindo/internal/postgresql"
	redisRepo "github.com/elangreza14/lion-superindo/internal/redis"
	"github.com/elangreza14/lion-superindo/internal/service"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

type ProductHandlerDeps struct {
	Handler     *handler.ProductHandler
	DB          *sql.DB
	RedisClient *redis.Client
}

var productSet = wire.NewSet(
	SetupDB,
	SetupCache,
	postgreRepo.NewProductRepo,
	wire.Bind(new(service.DbRepo), new(*postgreRepo.ProductRepo)), // <-- Bind DbRepo interface
	redisRepo.NewProductRepo,
	wire.Bind(new(service.CacheRepo), new(*redisRepo.ProductRepo)), // <-- Bind CacheRepo interface
	service.NewProductService,
	wire.Bind(new(handler.ProductService), new(*service.ProductService)), // <-- This line binds interface to implementation
	handler.NewProductHandler,
)

func InitializeProductHandler(cfg *config.Config) (*ProductHandlerDeps, error) {
	wire.Build(
		productSet,
		wire.Struct(new(ProductHandlerDeps), "Handler", "DB", "RedisClient"),
	)
	return nil, nil
}

func SetupDB(cfg *config.Config) (*sql.DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.POSTGRES_USER,
		cfg.POSTGRES_PASSWORD,
		cfg.POSTGRES_HOSTNAME,
		cfg.POSTGRES_PORT,
		cfg.POSTGRES_DB,
		cfg.POSTGRES_SSL,
	)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SetupCache(cfg *config.Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.REDIS_HOSTNAME, cfg.REDIS_PORT),
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return redisClient, nil
}
