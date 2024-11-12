package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	DirName = ".local"
)

func createDirIfNotExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(path, 0755); err != nil {
				return fmt.Errorf("Failed to create directory: %w", err)
			}
		} else {
			return fmt.Errorf("Failed to get stats of directory: %w", err)
		}
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}
	return nil
}

// 'dbName' is the name of the database file. Pass :memory: for in-memory database.
func NewSQLite(dbName string) (*sql.DB, error) {
	if dbName != ":memory:" {
		if err := createDirIfNotExists(DirName); err != nil {
			return nil, err
		}
		dbName = strings.Replace(dbName, ".db", "", 1)
		dbName = fmt.Sprintf("%s/%s.db", DirName, dbName)
	}

	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		return nil, fmt.Errorf("Failed to open sqlite database: %w", err)
	}

	// SQLite optimizations
	stmts := [...]string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA locking_mode = NORMAL;",
		"PRAGMA busy_timeout = 10000;",
		"PRAGMA cache_size = 10000;",
		"PRAGMA foreign_keys = ON;",
	}

	var errList []error

	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			errList = append(errList, err)
		}
	}

	if len(errList) > 0 {
		return nil, errors.Join(errList...)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping SQLite database: %w", err)
	}

	return db, nil
}
