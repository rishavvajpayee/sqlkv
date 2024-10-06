package handlers

import (
	"net/http"
	"sqlkv/config"
	"sqlkv/database"

	"github.com/labstack/echo"
)

func Seed(context echo.Context) error {
	var (
		app = context.Get("app").(*config.AppConfig)
	)
	logger := app.Logger
	seeded, err := database.CheckSeeded(app.DB)
	if err != nil {
		logger.Error("Failed to check if database is seeded", "error", err)
		return context.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to check if database is seeded"})
	}
	if !seeded {
		if err := database.SeedDatabase(context); err != nil {
			logger.Error("Failed to seed database", "error", err)
			return context.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to seed database"})
		}
		logger.Info("Database seeded successfully")
		return context.JSON(http.StatusOK, map[string]string{"message": "Database seeded successfully"})
	}
	return context.JSON(http.StatusOK, map[string]string{"message": "Database already seeded"})
}
