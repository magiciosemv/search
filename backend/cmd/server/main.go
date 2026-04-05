package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"solana-monitor/internal/config"
	"solana-monitor/internal/handlers"
	"solana-monitor/internal/middleware"
	"solana-monitor/internal/models"
	"solana-monitor/internal/services"
)

func main() {
	// Load .env file
	godotenv.Load() // loads from .env in current directory

	// Load configuration
	cfg := config.Load()

	// Ensure data directory exists
	dataDir := filepath.Dir(cfg.Database.Path)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Initialize database
	db, err := models.NewDB(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize schema
	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Initialize Solana service
	solana := services.NewSolanaService(cfg.Solana.RPCURL, cfg.Notification.ProxyURL)

	// Initialize notification service
	notifier := services.NewNotificationService(
		db,
		cfg.Notification.TelegramBotToken,
		services.SMTPConfig{
			Host:     cfg.Notification.SMTP.Host,
			Port:     cfg.Notification.SMTP.Port,
			Username: cfg.Notification.SMTP.Username,
			Password: cfg.Notification.SMTP.Password,
			From:     cfg.Notification.SMTP.From,
		},
		cfg.Notification.ProxyURL,
	)

	// Initialize monitor service
	monitor := services.NewMonitorService(db, solana, notifier)

	// Start monitor in background
	ctx, cancel := context.WithCancel(context.Background())
	go monitor.Start(ctx)

	// Initialize handlers
	h := handlers.NewHandler(db, solana, monitor, notifier)

	// Setup router
	router := setupRouter(h, cfg.Server.APIKey)

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(h *handlers.Handler, apiKey string) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check (no auth)
	router.GET("/api/health", h.Health)

	// All other API routes require auth
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(apiKey))
	{
		api.GET("/addresses", h.GetAddresses)
		api.POST("/addresses", h.CreateAddress)
		api.GET("/addresses/:id", h.GetAddress)
		api.PUT("/addresses/:id", h.UpdateAddress)
		api.DELETE("/addresses/:id", h.DeleteAddress)
		api.POST("/addresses/:id/refresh", h.RefreshAddressBalance)

		// Rules
		api.GET("/rules", h.GetRules)
		api.POST("/rules", h.CreateRule)
		api.PUT("/rules/:id", h.UpdateRule)
		api.DELETE("/rules/:id", h.DeleteRule)

		// Notifications
		api.GET("/notifications", h.GetNotifications)
		api.POST("/notifications", h.CreateNotification)
		api.PUT("/notifications/:id", h.UpdateNotification)
		api.DELETE("/notifications/:id", h.DeleteNotification)
		api.POST("/notifications/:id/test", h.TestNotification)

		// Alerts
		api.GET("/alerts", h.GetAlerts)
		api.GET("/alerts/stats", h.GetAlertStats)

		// Stats
		api.GET("/stats", h.GetStats)
	}

	return router
}