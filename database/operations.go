package database

import (
	"encoding/json"
	"fmt"
	"sqlkv/config"
	"strconv"
	"strings"
	"time"
)

func DbGetKey(app *config.AppConfig, key string) (interface{}, error) {
	var value string
	var valueType string
	row := app.DB.QueryRow("SELECT value, value_type FROM kv WHERE key = ?", key)
	err := row.Scan(&value, &valueType)
	if err != nil {
		return "", err
	}
	switch valueType {
	case "string":
		return value, nil
	case "number":
		if strings.Contains(value, ".") {
			return strconv.ParseFloat(value, 64)
		}
		return strconv.ParseInt(value, 10, 64)
	case "json":
		var jsonValue interface{}
		err := json.Unmarshal([]byte(value), &jsonValue)
		if err != nil {
			return nil, err
		}
		return jsonValue, nil
	default:
		return nil, fmt.Errorf("unknown value type %s", valueType)
	}
}

func DbSetKey(app *config.AppConfig, key string, value interface{}, exp int64) (string, error) {
	valueStr, valueType := convertValue(value)
	currentTime := time.Now().UTC()
	var expTime *time.Time

	if exp != 0 {
		t := currentTime.Add(time.Duration(exp) * time.Second)
		expTime = &t
	}

	// a single UPSERT operation instead of separate INSERT or UPDATE
	query := `
        INSERT INTO kv (key, value, value_type, expires_in)
        VALUES (?, ?, ?, ?)
        ON CONFLICT(key) DO UPDATE SET
        value = excluded.value,
        value_type = excluded.value_type,
        expires_in = excluded.expires_in
    `

	_, err := app.DB.Exec(query, key, valueStr, valueType, expTime)
	if err != nil {
		return "", err
	}

	return valueStr, nil
}

func convertValue(value interface{}) (string, string) {
	switch v := value.(type) {
	case string:
		return v, "string"
	case int, int32, int64, float32, float64:
		return fmt.Sprintf("%v", v), "number"
	case bool:
		return strconv.FormatBool(v), "boolean"
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return "", "unknown"
		}
		return string(jsonBytes), "json"
	}
}
