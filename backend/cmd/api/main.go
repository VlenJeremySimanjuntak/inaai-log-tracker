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

// @title 		Log Tracker API
// @version 	1.0
// @description API untuk pelaporan gangguan operasional.
// @host 		backend-service-production.up.railway.app
// @BasePath 	/api

func main() {
	// 1. Load .env (opsional di Railway, tapi aman untuk lokal)
	_ = godotenv.Load()

	// 2. Inisialisasi Database
	// Memastikan DB terkoneksi sebelum aplikasi lanjut
	db := config.InitDB()
	defer db.Close()

	// 3. Setup Event Bus
	eventBus := event.NewBus(5)
	defer eventBus.Stop()

	// Subscribe event dasar
	eventBus.Subscribe("log.status.updated", func(ctx context.Context, e event.Event) error {
		log.Printf("📢 Event: log status updated. Data: %v", e.Data)
		return nil
	})

	// 4. Inisialisasi Layer (Repository -> Usecase -> Delivery)
	logRepo := repository.NewMysqlLogRepository(db)
	logUsecase := usecase.NewLogUsecase(logRepo, eventBus)

	// 5. Setup Echo
	e := echo.New()
	e.Use(middleware.Logger()) // Tambahkan Logger agar bisa melihat akses di log
	e.Use(middleware.Recover()) // Tambahkan Recover agar server tidak mati jika ada panic
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, 
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}, // Tambahkan OPTIONS
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// 6. Routes
	delivery.NewLogHandler(e, logUsecase)
	summaryHandler := delivery.NewSummaryHandler(logUsecase)
	e.GET("/api/summary/latest", summaryHandler.GetLatestSummary)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// 7. Port Handling
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Standar port cloud sering menggunakan 8080
	}

	// 8. Graceful Shutdown Implementation
	go func() {
		log.Printf("🚀 Aplikasi Backend Mengetuk Port %s...", port)
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Gagal menjalankan server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Mematikan server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}
	log.Println("Server exited gracefully")
}