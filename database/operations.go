package database

import (
	"sqlkv/config"
	"strconv"
	"time"
)

func DbGetKey(app *config.AppConfig, key string) (string, error) {
	var value string
	row := app.DB.QueryRow("SELECT value FROM kv WHERE key = ?", key)
	err := row.Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func DbSetKey(app *config.AppConfig, key string, value string, exp int64) (string, error) {
	var exists bool
	err := app.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM kv WHERE key = ?)", key).Scan(&exists)
	if err != nil {
		return "", err
	}
	currentTime := time.Now().UTC()
	if exists {
		query := "UPDATE kv SET value=?"
		args := []interface{}{value, key}
		if exp != 0 {
			expTime := currentTime.Add(time.Duration(exp) * time.Second)
			query += ", expires_in=?"
			args = []interface{}{value, expTime, key}
		}
		_, err := app.DB.Exec(query+" WHERE key=?", args...)
		if err != nil {
			return "", err
		}
		return value, nil
	}

	var expTime time.Time
	if exp != 0 {
		expTime = currentTime.Add(time.Duration(exp) * time.Second)
	}

	row, err := app.DB.Exec(
		"INSERT INTO kv (key, value, expires_in) VALUES (? , ?, ?)",
		key,
		value,
		expTime,
	)
	if err != nil {
		return "", err
	}
	id, err := row.LastInsertId()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(id)), nil
}
