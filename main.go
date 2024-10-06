package main

import (
	"log/slog"
	"net/http"
	"os"
	"sqlkv/config"
	"sqlkv/database"
	"sqlkv/handlers"

	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

func serverHealthCheck(context echo.Context) error {
	return context.JSON(http.StatusOK, map[string]string{"status": "success"})
}

func main() {
	// Load config
	if err := config.LoadConfig(); err != nil {
		slog.Any("error", err)
		os.Exit(1)
	}

	// Initialize DB
	db, err := database.InitAppDB("sqlite3", "./sqlkv.db")
	if err != nil {
		slog.Any("error", err)
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

	// Routes
	server.GET("/", serverHealthCheck)
	server.GET("/seed", handlers.Seed)
	server.GET("kv/get/:key", handlers.GetKey)
	server.POST("kv/set", handlers.SetKey)

	slog.Info("Starting HTTP server on :8000")
	server.Logger.Fatal(server.Start(":8000"))
}
