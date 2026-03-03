package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/vo1dFl0w/loyalty-card-service/internal/adapters/storage/postgres"
	"github.com/vo1dFl0w/loyalty-card-service/internal/config"
	ht "github.com/vo1dFl0w/loyalty-card-service/internal/transport/http"
	"github.com/vo1dFl0w/loyalty-card-service/internal/transport/http/httpgen"
	"github.com/vo1dFl0w/loyalty-card-service/internal/usecase"
	"github.com/vo1dFl0w/loyalty-card-service/pkg/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		log.Println(ctx, "startup", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	loggerCfg := logger.NewLoggerConfig(cfg.Server.Env, cfg.Server.LoggerTimeFormat)
	logger := logger.LoadLogger(loggerCfg)

	databaseDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password, cfg.DB.DBName, cfg.DB.Sslmode,
	)

	var db *sql.DB
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", databaseDSN)
		if err == nil {
			break
		}
		logger.Info("waiting for database...")
		time.Sleep(time.Second * 1)
	}
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	logger.Info("connection to the database established")

	storage := postgres.NewStorage(db)
	loyaltyCardRepo := storage.LoyaltyCardRepo()

	loyaltyCardSrv := usecase.NewLoyaltyCardService(loyaltyCardRepo)

	handler := ht.NewHandler(cfg, logger, loyaltyCardSrv)

	server, err := httpgen.NewServer(handler)
	if err != nil {
		return fmt.Errorf("new server: %w", err)
	}

	middlewares := handler.CORSMiddleware(
		handler.RequestIDMiddleware(
			handler.LoggerMiddleware(
				handler.RequestTimeoutMiddleware(server),
			),
		),
	)

	srv := http.Server{
		Addr:    cfg.Server.HTTPAddr,
		Handler: middlewares,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("server started", "host", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		} else {
			serverErr <- nil
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case s := <-sig:
		logger.Info("initialization gracefull shutdown", "signal", s)
		shutdownCtx, cancel := context.WithTimeout(ctx, cfg.Server.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}
		logger.Info("server gracefully stopped")
		return nil
	}
}
