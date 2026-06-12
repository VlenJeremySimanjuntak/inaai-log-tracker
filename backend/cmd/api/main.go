// cmd/api/main.go
package main

import (
	"context"
	"log"
	"net/http"
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
	"github.com/labstack/echo/v4/middleware"
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

	eventBus := event.NewBus(5)
	defer eventBus.Stop()

	eventBus.Subscribe("log.status.updated", func(ctx context.Context, e event.Event) error {
		log.Printf("📢 Event: log status updated. Data: %v", e.Data)
		return nil
	})

	logRepo := repository.NewMysqlLogRepository(db)
	logUsecase := usecase.NewLogUsecase(logRepo, eventBus)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	// Handler untuk Incident Logs (WAJIB)
	delivery.NewLogHandler(e, logUsecase)

	// Handler untuk AI Summary
	summaryHandler := delivery.NewSummaryHandler(logUsecase)
	e.GET("/api/summary/latest", summaryHandler.GetLatestSummary)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8081"
	}

	go func() {
		log.Printf("🚀 Aplikasi Backend Mengetuk Port %s...", port)
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutdown server: %v", err)
		}
	}()

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