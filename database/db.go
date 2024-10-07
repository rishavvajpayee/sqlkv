package database

import (
	"database/sql"
)

func InitAppDB(engine string, path string) (*sql.DB, error) {
	db, e := sql.Open(engine, path)
	return db, e
}
