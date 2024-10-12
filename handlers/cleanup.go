package handlers

import (
	"context"
	"database/sql"
	"log/slog"
	"time"
)

func DbCleanUp(ctx context.Context, db *sql.DB) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("Cleanup context cancelled (Go Routine Killed)")
			return
		default:
			if _, err := db.Exec(`
				DELETE FROM kv WHERE expires_in <= CURRENT_TIMESTAMP;
			`); err != nil {
				slog.Error("Failed to cleanup database", "error", err)
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}
