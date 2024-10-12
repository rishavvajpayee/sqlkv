package handlers

import (
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

func getAppConfig(context echo.Context) *config.AppConfig {
	return context.Get("app").(*config.AppConfig)
}

func GetKey(context echo.Context) error {
	// Get app config
	app := getAppConfig(context)

	// Get key from path
	key := context.Param("key")

	// Get value from database
	value, err := database.DbGetKey(app, key)
	if err != nil {
		return context.JSON(http.StatusNotFound, map[string]string{"message": "key not found"})
	}
	return context.JSON(http.StatusOK, map[string]interface{}{key: value})
}

func SetKey(context echo.Context) error {
	// Get app config
	app := getAppConfig(context)

	// Bind request
	request := new(SetKeyRequest)

	if err := context.Bind(request); err != nil {
		return context.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}

	if _, err := database.DbSetKey(app, request.Key, request.Value, request.ExpiresIn); err != nil {
		return context.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to set key", "err": err.Error()})
	}
	return context.JSON(http.StatusOK, map[string]string{"message": "Ok"})
}
