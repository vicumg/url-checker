package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"urlChecker/internal/application/service"
	httpInfra "urlChecker/internal/infrastructure/http"
	"urlChecker/internal/infrastructure/logger"
	"urlChecker/internal/infrastructure/repository"
	"urlChecker/internal/interface/api"
)

func main() {

	dbPath := "./data/monitors.db"

	os.MkdirAll("./data", 0755)

	// repo := repository.NewMemoryRepository()
	repo, err := repository.NewSQLiteRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to init repository: %v", err)
	}
	defer repo.Close()

	fileLogger, err := logger.NewFileLogger("./logs")
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer fileLogger.Close()

	monitorService := service.NewMonitorService(repo)
	checkerService := service.NewCheckerService(repo, fileLogger)
	handler := api.NewHandler(monitorService)
	router := httpInfra.NewRouter(handler)

	// HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go checkerService.Start(ctx)

	go func() {
		fmt.Println("Server started on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")

	// Graceful shutdown
	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server stopped")
}
