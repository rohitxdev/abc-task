package repo

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ClassNotFoundError    = errors.New("Class not found")
	ClassFullError        = errors.New("Class is full")
	InvalidDateRangeError = errors.New("No class is available on the given date")
)

type Repo struct {
	// Must be an active SQLite database
	db *sql.DB
}

func New(db *sql.DB) (*Repo, error) {
	if err := MigrateUp(db); err != nil {
		return nil, err
	}
	return &Repo{db}, nil
}

func MigrateUp(db *sql.DB) error {
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS classes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		start_date INTEGER NOT NULL,
		end_date INTEGER NOT NULL,
		capacity INTEGER NOT NULL
	);`); err != nil {
		return err
	}
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS bookings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		class_id INTEGER NOT NULL,
		member_name TEXT NOT NULL,
		date INTEGER NOT NULL,
		FOREIGN KEY (class_id) REFERENCES classes(id)
	);`); err != nil {
		return err
	}
	return nil
}

// 'startDate' and 'endDate' are in UNIX timestamp format
func (r *Repo) CreateClass(ctx context.Context, name string, startDate int64, endDate int64, capacity uint) error {
	query := "INSERT INTO classes (name, start_date, end_date, capacity) VALUES (?, ?, ?, ?);"
	_, err := r.db.ExecContext(ctx, query, name, startDate, endDate, capacity)
	return err
}

// 'date' is in UNIX timestamp format
func (r *Repo) CreateBooking(ctx context.Context, classID uint64, memberName string, date int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, "SELECT start_date, end_date, capacity FROM classes WHERE id = ?;", classID)
	var startDate int64
	var endDate int64
	var capacity uint
	if err = row.Scan(&startDate, &endDate, &capacity); err != nil {
		if err == sql.ErrNoRows {
			return ClassNotFoundError
		}
		return err
	}
	if startDate > date || endDate < date {
		return InvalidDateRangeError
	}

	row = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM bookings WHERE class_id = ? AND date = ?;", classID, date)
	var occupancy uint
	if err = row.Scan(&occupancy); err != nil {
		return err
	}

	if occupancy >= capacity {
		return ClassFullError
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO bookings (class_id, member_name, date) VALUES (?, ?, ?);", classID, memberName, date); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
