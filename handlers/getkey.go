package handlers

import (
	"net/http"
	"sqlkv/config"
	"sqlkv/database"

	"github.com/labstack/echo"
)

type SetKeyRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func GetKey(context echo.Context) error {
	var (
		app = context.Get("app").(*config.AppConfig)
	)
	key := context.Param("key")
	value, err := database.DbGetKey(app, key)
	if err != nil {
		return context.JSON(500, map[string]string{"message": "key not found"})
	}
	return context.JSON(http.StatusOK, map[string]string{key: value})
}

func SetKey(context echo.Context) error {
	var app = context.Get("app").(*config.AppConfig)
	request := new(SetKeyRequest)

	err := context.Bind(request)
	if err != nil {
		return context.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request body"})
	}
	value, err := database.DbSetKey(app, request.Key, request.Value)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, map[string]string{"message": "failed to set key"})
	}
	return context.JSON(http.StatusOK, map[string]string{"message": "key set successfully", "id": value})
}
