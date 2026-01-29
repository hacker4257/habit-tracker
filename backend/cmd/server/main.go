package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"habit-tracker/internal/config"
	"habit-tracker/internal/handler"
	"habit-tracker/internal/middleware"
	"habit-tracker/internal/repository"
	"habit-tracker/internal/service"
	"habit-tracker/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize repository
	repo, err := repository.NewRecordRepository(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to initialize repository: %v", err)
	}

	// Initialize service
	svc := service.NewRecordService(repo)

	// Initialize handler
	h := handler.NewRecordHandler(svc)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/records", h.HandleRecords)
	mux.HandleFunc("/api/records/", h.HandleRecord)
	mux.HandleFunc("/api/stats", h.HandleStats)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply middleware
	var handler http.Handler = mux
	handler = middleware.CORS(cfg.Server.AllowOrigins)(handler)
	handler = middleware.Logging(handler)
	handler = middleware.Recovery(handler)

	// Start server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: handler,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		logger.Info("Shutting down server...")
		server.Close()
	}()

	logger.Info("Server starting on http://localhost:%s", cfg.Server.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("Server error: %v", err)
	}
}
