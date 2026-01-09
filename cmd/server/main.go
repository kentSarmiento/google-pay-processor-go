package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/kentSarmiento/google-pay-processor-go/internal/application"
	"github.com/kentSarmiento/google-pay-processor-go/internal/infrastructure/config"
	"github.com/kentSarmiento/google-pay-processor-go/internal/infrastructure/crypto"
	"github.com/kentSarmiento/google-pay-processor-go/internal/infrastructure/logger"
	"github.com/kentSarmiento/google-pay-processor-go/internal/infrastructure/metrics"
	"github.com/kentSarmiento/google-pay-processor-go/internal/infrastructure/payment"
	"github.com/kentSarmiento/google-pay-processor-go/pkg/googlepay"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.NewZerologLogger(cfg.Logging.Level, cfg.Logging.Format)
	log.Info("Starting Google Pay Processor", map[string]interface{}{
		"version":     "1.0.0",
		"environment": cfg.GooglePay.Environment,
	})

	// Initialize metrics collector
	metricsCollector := metrics.NewPrometheusCollector()

	// Initialize token decryptor
	tokenDecryptor := crypto.NewTinkDecryptor(
		cfg.GooglePay.MerchantID,
		cfg.GooglePay.ProtocolVersion,
		log,
	)

	// Initialize payment processor
	// In production, replace MockCardProcessor with actual processor client
	paymentProcessor := payment.NewMockCardProcessor(log)

	// Initialize Google Pay service
	googlePayService := application.NewGooglePayService(
		tokenDecryptor,
		paymentProcessor,
		log,
		metricsCollector,
	)

	// Initialize HTTP handler
	handler := googlepay.NewHandler(googlePayService, log, metricsCollector)

	// Setup HTTP router
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/v1/googlepay/process", handler.ProcessPayment)
	mux.HandleFunc("/health", handler.HealthCheck)
	mux.HandleFunc("/ready", handler.ReadinessCheck)

	// Create main server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start metrics server if enabled
	var metricsServer *http.Server
	if cfg.Metrics.Enabled {
		metricsMux := http.NewServeMux()
		metricsMux.Handle(cfg.Metrics.Path, promhttp.Handler())

		metricsServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Metrics.Port),
			Handler: metricsMux,
		}

		go func() {
			log.Info("Starting metrics server", map[string]interface{}{
				"port": cfg.Metrics.Port,
				"path": cfg.Metrics.Path,
			})
			if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Error("Metrics server error", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}()
	}

	// Start server in a goroutine
	go func() {
		log.Info("Starting HTTP server", map[string]interface{}{
			"port": cfg.Server.Port,
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...", nil)

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", map[string]interface{}{
			"error": err.Error(),
		})
	}

	if metricsServer != nil {
		if err := metricsServer.Shutdown(ctx); err != nil {
			log.Error("Metrics server forced to shutdown", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	log.Info("Server exited", nil)
}
