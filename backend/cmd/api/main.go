package main

import (
	"log"
	"os"

	"backend-tracker/internal/config"
	"backend-tracker/internal/delivery"
	"backend-tracker/internal/repository"
	"backend-tracker/internal/usecase"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load .env optional
	_ = godotenv.Load()

	// Init DB (dari config Anda)
	db := config.InitDB()
	defer db.Close()

	// Setup repository & usecase
	logRepo := repository.NewMysqlLogRepository(db)
	logUsecase := usecase.NewLogUsecase(logRepo)

	// Setup Echo
	e := echo.New()
	delivery.NewLogHandler(e, logUsecase)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Server berjalan di port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}