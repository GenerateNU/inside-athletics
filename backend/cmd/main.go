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
	"gorm.io/gorm"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbUrl := os.Getenv("DB_CONNECTION_STRING")

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Successefully connected to Supabase 🚀")
	defer db.Close()

	app := server.CreateApp(db)
	app.Server.Listen("localhost:8080")

	// gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Server is shutting down...")

	if err := app.Server.ShutdownWithContext(ctx); err != nil {
		log.Fatal("Failed to shutdown server:", err)
	}

	slog.Info("Server shut down successfully")
}
