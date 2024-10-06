package util

import (
	"log/slog"
	"os"
)

func SetupLogger() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return logger
}
