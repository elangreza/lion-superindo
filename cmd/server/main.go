package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/elangreza14/lion-superindo/cmd/server/config"
	"github.com/elangreza14/lion-superindo/internal/handler"
	"github.com/elangreza14/lion-superindo/internal/postgresql"
	redisRepo "github.com/elangreza14/lion-superindo/internal/redis"
	"github.com/elangreza14/lion-superindo/internal/service"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

// TODO google wire
// TODO finish unit test
// TODO swagger or openapi
// TODO POSTMAN collection
// TODO custom error handling
// TODO move into newest git repository
// TODO add ratelimiter

func main() {
	cfg, err := config.LoadConfig()
	errChecker(err)

	dbPool, err := setupDB(cfg)
	errChecker(err)
	defer dbPool.Close()

	cacheClient, err := setupCache(cfg)
	errChecker(err)

	productDb := postgresql.NewProductRepo(dbPool)
	productCache := redisRepo.NewProductRepo(cacheClient)
	productService := service.NewProductService(productDb, productCache)
	productHandler := handler.NewProductHandler(productService)

	mux := http.NewServeMux()
	mux.HandleFunc("/product", productHandler.ProductHandler)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.HTTP_PORT),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	slog.Info("server started", "port", cfg.HTTP_PORT)

	<-gracefulShutdown(context.Background(), 5*time.Second,
		func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
		func(ctx context.Context) error {
			return dbPool.Close()
		},
		func(ctx context.Context) error {
			cacheClient.Close()
			return nil
		},
	)
}

func errChecker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func setupDB(cfg *config.Config) (*sql.DB, error) {
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

func setupCache(cfg *config.Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.REDIS_HOSTNAME, cfg.REDIS_PORT),
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return redisClient, nil
}

type operation func(ctx context.Context) error

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops ...operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		slog.Info("shutting down")

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		go func() {
			<-ctx.Done()
			slog.Info("force quit the app")
			wait <- struct{}{}
		}()

		var wg sync.WaitGroup

		for key, op := range ops {
			wg.Add(1)
			go func(key int, op operation) {
				defer wg.Done()
				processName := fmt.Sprintf("process %d", key)

				if err := op(ctx); err != nil {
					slog.Error(processName, "err", err.Error())
					return
				}

				slog.Info(processName, "message", "success")
			}(key, op)
		}

		wg.Wait()
	}()

	return wait
}
