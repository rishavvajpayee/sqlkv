package database

import (
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"os"
	"sqlkv/config"
	"sync"

	"github.com/labstack/echo"
)

func ReadSchema(path string, wg *sync.WaitGroup, contentCh chan []byte, errCh chan error) {
	defer wg.Done()
	file, err := os.Open(path)
	if err != nil {
		errCh <- err
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		errCh <- err
		return
	}

	contentCh <- content
}

func SeedDatabase(context echo.Context) error {
	var app = context.Get("app").(*config.AppConfig)
	db := app.DB

	var wg sync.WaitGroup
	contentCh := make(chan []byte)
	errCh := make(chan error)

	wg.Add(1)
	go ReadSchema(config.ServerConfig.SchemaSeedFilePath, &wg, contentCh, errCh)

	go func() {
		wg.Wait()
		close(contentCh)
		close(errCh)
	}()

	var content []byte
	var err error

	for {
		select {
		case content = <-contentCh:
			if content == nil {
				slog.Error("Failed to read schema")
				return errors.New("failed to read schema")
			}
			slog.Info("Schema read successfully")
		case err = <-errCh:
			if err != nil {
				slog.Error("Failed to read schema", "error", err)
				return errors.New(err.Error())
			}
			slog.Info("a error occured while reading schema")
		}
		if len(content) > 0 {
			break
		}
	}
	_, execErr := db.Exec(string(content))
	if execErr != nil {
		return execErr
	}

	return nil
}

func CheckSeeded(db *sql.DB) (bool, error) {
	var seeded bool
	row := db.QueryRow("SELECT seeded FROM history")
	err := row.Scan(&seeded)
	if err != nil {
		return false, nil
	}
	return seeded, nil
}
