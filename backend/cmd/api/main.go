// cmd/api/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"backend-tracker/internal/config"
	"backend-tracker/internal/delivery"
	"backend-tracker/internal/event"
	"backend-tracker/internal/repository"
	"backend-tracker/internal/usecase"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "backend-tracker/docs"
)

// @title           Log Tracker API
// @version         1.0
// @description     API untuk pelaporan gangguan operasional dengan pessimistic locking.
// @host            localhost:8081
// @BasePath        /api

func main() {
	_ = godotenv.Load()
	db := config.InitDB()
	defer db.Close()

	// Inisialisasi event bus dengan 5 worker
	eventBus := event.NewBus(5)
	defer eventBus.Stop()

	// Subscribe contoh event
	eventBus.Subscribe("log.status.updated", func(ctx context.Context, e event.Event) error {
		log.Printf("📢 Event: log status updated. Data: %v", e.Data)
		// Di sini bisa ditambahkan notifikasi, audit, atau trigger AI summary
		return nil
	})

	logRepo := repository.NewMysqlLogRepository(db)
	logUsecase := usecase.NewLogUsecase(logRepo, eventBus) // kirim eventBus

	e := echo.New()
	delivery.NewLogHandler(e, logUsecase)

	// Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8081"
	}

	// Graceful shutdown
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutdown server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server exited gracefully")
}