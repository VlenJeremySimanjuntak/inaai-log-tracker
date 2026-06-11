package main

import (
	"log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	
	"backend-tracker/internal/config"
	"backend-tracker/internal/delivery"
	"backend-tracker/internal/repository"
	"backend-tracker/internal/usecase"
)

func main() {
	// Inisialisasi Database MySQL
	db := config.InitDB()
	defer db.Close()

	e := echo.New()

	// Middleware pengaman CORS agar frontend React bisa menembak API di port berbeda
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT},
	}))

	// Registrasi Dependency Injection (DI) Clean Architecture
	logRepo := repository.NewMysqlLogRepository(db)
	logUseCase := usecase.NewLogUsecase(logRepo)
	delivery.NewLogHandler(e, logUseCase)

	log.Println("🚀 Aplikasi Backend Mengetuk Port 8081...")
	e.Logger.Fatal(e.Start(":8081"))
}