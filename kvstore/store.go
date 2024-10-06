package kvstore

import (
	"database/sql"
	"sync"
)

type KVStore struct {
	DB *sql.DB
}

// NewKVStore creates a new KV store instance
func NewKVStore(db *sql.DB) *KVStore {
	return &KVStore{DB: db}
}

func (kv *KVStore) Get(key string) (string, error) {
	var value string
	err := kv.DB.QueryRow("SELECT value FROM kv WHERE key = ?", key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (kv *KVStore) Set(key string, value string) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	var exists bool
	err := kv.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM kv WHERE key = ?)", key).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		_, err := kv.DB.Exec("UPDATE kv SET value=? WHERE key=?", value, key)
		return err
	}
	_, err = kv.DB.Exec("INSERT INTO kv (key, value) VALUES (?, ?)", key, value)
	return err
}
