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

	"github.com/elangreza14/superindo/cmd/rest/config"
	"github.com/elangreza14/superindo/internal/handler"
	"github.com/elangreza14/superindo/internal/postgresql"
	"github.com/elangreza14/superindo/internal/service"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	errChecker(err)

	dbPool, err := setupDB(cfg)
	errChecker(err)
	defer dbPool.Close()

	redisPool := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	productRepo := postgresql.NewProductRepo(dbPool, redisPool)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	mux := http.NewServeMux()
	mux.HandleFunc("/product", productHandler.ProductHandler)

	srv := &http.Server{
		Addr:           cfg.HTTP_PORT,
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

	<-gracefulShutdown(context.Background(), 5*time.Second,
		func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
		func(ctx context.Context) error {
			return dbPool.Close()
		},
		func(ctx context.Context) error {
			return redisPool.Shutdown(ctx).Err()
		},
	)
}

func errChecker(err error) {
	if err != nil {
		panic(err)
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

	return sql.Open("postgres", connString)
}

type operation func(ctx context.Context) error

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops ...operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		fmt.Println("a")
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
		cancel()
	}()

	return wait
}
