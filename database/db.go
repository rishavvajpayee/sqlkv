package database

import (
	"database/sql"
)

// func InitDB(engine string, path string, wg *sync.WaitGroup, dbChannel chan *sql.DB, err chan error) {
// 	defer wg.Done()
// 	db, e := sql.Open(engine, path)
// 	if e != nil {
// 		err <- e
// 		return
// 	}
// 	dbChannel <- db
// }

func InitAppDB(engine string, path string) (*sql.DB, error) {
	db, e := sql.Open(engine, path)
	return db, e
}
