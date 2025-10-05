package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"sqlkv/database"

	"github.com/labstack/echo"
)

type SetKeyRequest struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	ExpiresIn int64       `json:"expires_in,omitempty"`
}

func GetKey(db *sql.DB, slogger *slog.Logger) echo.HandlerFunc {
	// Get app config
	return func(context echo.Context) error {
		slogger.Info("Getting the Key")
		key := context.Param("key")
		value, err := database.DbGetKey(slogger, db, key)
		if err != nil {
			return context.JSON(http.StatusNotFound, map[string]string{"message": "key not found"})
		}
		return context.JSON(http.StatusOK, map[string]interface{}{key: value})
	}
}

func SetKey(db *sql.DB, slogger *slog.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(SetKeyRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
		}
		if _, err := database.DbSetKey(slogger, db, request.Key, request.Value, request.ExpiresIn); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to set key", "err": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Ok"})
	}
}
