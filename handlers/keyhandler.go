package handlers

import (
	"database/sql"
	"net/http"
	"sqlkv/config"
	"sqlkv/database"

	"github.com/labstack/echo"
)

type SetKeyRequest struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	ExpiresIn int64       `json:"expires_in,omitempty"`
}

func InitialSeedDatabase(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS kv (id INTEGER PRIMARY KEY, key TEXT NOT NULL, value TEXT, value_type TEXT, expires_in TIMESTAMP DEFAULT NULL, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS history ( seeded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, seeded BOOLEAN DEFAULT FALSE );")
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO history (seeded) VALUES (0)")
	if err != nil {
		return err
	}
	return nil
}

func GetKey(context echo.Context) error {
	var (
		app = context.Get("app").(*config.AppConfig)
	)
	key := context.Param("key")
	value, err := database.DbGetKey(app, key)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, map[string]string{"message": "key not found"})
	}
	return context.JSON(http.StatusOK, map[string]interface{}{key: value})
}

func SetKey(context echo.Context) error {
	var app = context.Get("app").(*config.AppConfig)
	request := new(SetKeyRequest)

	err := context.Bind(request)
	if err != nil {
		return context.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}
	value, err := database.DbSetKey(app, request.Key, request.Value, request.ExpiresIn)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to set key"})
	}
	return context.JSON(http.StatusOK, map[string]string{"message": "key set successfully", "id": value})
}
