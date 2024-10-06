package config

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DB     *sql.DB
	Logger *slog.Logger
}

type SrvConfig struct {
	SchemaSeedFilePath string
}

var ServerConfig SrvConfig

func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		slog.Any("error", err)
	}

	ServerConfig = SrvConfig{
		SchemaSeedFilePath: getEnvVar("SCHEMA_FILE_PATH", "./schema.sql"),
	}
	return nil
}

func getEnvVar(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
