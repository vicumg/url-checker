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
	"urlChecker/internal/infrastructure/repository"
	"urlChecker/internal/interface/api"
)

func main() {
	// Инициализация слоев
	repo := repository.NewMemoryRepository()
	monitorService := service.NewMonitorService(repo)
	checkerService := service.NewCheckerService(repo)
	handler := api.NewHandler(monitorService)
	router := httpInfra.NewRouter(handler)

	// HTTP сервер
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запуск checker service в отдельной goroutine
	go checkerService.Start(ctx)

	// Запуск HTTP сервера в отдельной goroutine
	go func() {
		fmt.Println("Server started on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Ожидание сигнала завершения
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
