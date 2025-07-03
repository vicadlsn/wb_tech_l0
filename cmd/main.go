package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"webtechl0/internal/cache"
	"webtechl0/internal/config"
	"webtechl0/internal/handler"
	"webtechl0/internal/kafka"
	"webtechl0/internal/models"
	"webtechl0/internal/postgres"
	"webtechl0/internal/repository"
	"webtechl0/internal/service"
)

const defaultConfigPath = "./configs/config.yml"

func main() {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	cfg, err := config.New(defaultConfigPath)
	if err != nil {
		log.Error("Failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := postgres.New(ctx, cfg.Database)
	if err != nil {
		log.Error("Failed to connect to postgres", slog.Any("error", err))
		os.Exit(1)
	}
	defer pool.Close()

	orderRepository := repository.NewOrderRepository(pool)
	orderCache := cache.NewLRUCache[string, *models.Order](10)
	orderService := service.NewOrderService(orderRepository, orderCache)

	orderService.FillCache(ctx)

	orderHandler := handler.NewOrderHandler(orderService, log)

	router := handler.NewRouter(orderHandler, log)
	addr := cfg.HTTP.Host + ":" + cfg.HTTP.Port
	server := http.Server{
		Addr:    addr,
		Handler: router,
	}

	consumerHandler := kafka.NewOrderHandler(orderService, log)
	consumer := kafka.NewConsumer(cfg.Kafka, consumerHandler, log)

	errChan := make(chan error, 1)
	log.Info("Starting kafka consumer")
	go func() {
		errChan <- consumer.Start(ctx)
	}()

	log.Info("Starting http server", slog.String("addr", addr))
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server error", slog.Any("error", err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutdown signal received")
	cancel()

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Error("Failed to shutdown http server", slog.Any("error", err))
	}

	if err := <-errChan; err != nil && err != context.Canceled {
		log.Error("Kafka consumer failed", slog.Any("error", err))
	}

}
