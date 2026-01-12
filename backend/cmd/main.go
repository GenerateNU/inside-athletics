package main

import (
	"context"
	"fmt"
	"inside-athletics/internal/server"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUrl := os.Getenv("IPV4")

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Successefully connected to Supabase 🚀")

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get database instance: %v\n", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	app := server.CreateApp(db)
	app.Server.Listen("localhost:8080")

	// gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Server is shutting down...")

	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Server.ShutdownWithContext(shutdownCtx); err != nil {
		log.Fatal("Failed to shutdown server:", err)
	}

	slog.Info("Server shut down successfully")
}
