package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sqlkv/config"
	"sqlkv/database"
	"sqlkv/handlers"
	"time"

	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load config
	if err := config.LoadConfig(); err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize DB
	db, err := database.InitAppDB("sqlite3", "./sqlkv.db")
	if err != nil {
		slog.Error("Failed to initialize DB", "error", err)
		os.Exit(1)
	}

	server := echo.New()

	app := &config.AppConfig{
		DB:     db,
		Logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	// Middleware to attach app config
	server.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			context.Set("app", app)
			return next(context)
		}
	})

	// Seed database
	if err := handlers.InitialSeedDatabase(db); err != nil {
		slog.Error("Failed to seed database", "error", err)
		os.Exit(1)
	}
	slog.Info("Initial seed database executed successfully")

	// Start cleanup routine
	ctx, dbCancel := context.WithCancel(context.Background())
	go handlers.DbCleanUp(ctx, db)

	// Routes
	server.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{"status": "success"})
	})
	server.GET("kv/get/:key", handlers.GetKey)
	server.POST("kv/set", handlers.SetKey)

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		dbCancel()
		ctx, routineCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer routineCancel()

		// Shutdown server gracefully
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Error shutting down server", "error", err)
		}
	}()

	slog.Info("Starting HTTP server", "address", ":8000")
	server.Logger.Fatal(server.Start(":8000"))
}
