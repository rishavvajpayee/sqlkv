package main

import (
	"context"
	"database/sql"
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
	if err := InitialSeedDatabase(db); err != nil {
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

func InitialSeedDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS kv
		(id INTEGER PRIMARY KEY,
		key TEXT NOT NULL,
		value TEXT,
		value_type TEXT,
		expires_in TIMESTAMP DEFAULT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
	`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS
		history (
			seeded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			seeded BOOLEAN DEFAULT FALSE
		);`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		INSERT INTO history (seeded) VALUES (0)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS
		idx_kv_key ON kv(key);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS
		idx_kv_expires_in ON kv(expires_in);
	`)
	if err != nil {
		return err
	}

	return nil
}
