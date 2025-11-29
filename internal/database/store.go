package database

import (
	"database/sql"
	"os"
	"path/filepath"
	
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	Db *sql.DB
}

var s_db *Store

func GetDatabase() (*Store, error) {
	if s_db != nil {
		return s_db, nil
	}

	dbPath := "/config/clonis.db"

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	s_db := &Store{Db: db}
	if err := s_db.migrate(); err != nil {
		return nil, err
	}

	return s_db, nil
}