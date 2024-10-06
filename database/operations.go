package database

import (
	"sqlkv/config"
	"strconv"
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

func DbSetKey(app *config.AppConfig, key string, value string) (string, error) {
	var exists bool
	err := app.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM kv WHERE key = ?)", key).Scan(&exists)
	if err != nil {
		return "", err
	}
	if exists {
		_, err := app.DB.Exec("UPDATE kv SET value=? WHERE key=?", value, key)
		if err != nil {
			return "", err
		}
		return value, nil
	}
	row, err := app.DB.Exec(
		"INSERT INTO kv (key, value) VALUES (? , ?)",
		key,
		value,
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
