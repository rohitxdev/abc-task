package repo

import (
	"context"
	"database/sql"
)

type RepoError struct {
	Message string
}

func (e RepoError) Error() string {
	return e.Message
}

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

type Class struct {
	Name      string
	Capacity  uint
	StartDate int64
	EndDate   int64
	ID        int64
}

// 'startDate' and 'endDate' are in UNIX timestamp format
func (r *Repo) CreateClass(name string, startDate int64, endDate int64, capacity uint) error {
	query := "INSERT INTO classes (name, start_date, end_date, capacity) VALUES (?, ?, ?, ?);"
	_, err := r.db.ExecContext(context.Background(), query, name, startDate, endDate, capacity)
	return err
}

// 'date' is in UNIX timestamp format
func (r *Repo) CreateBooking(ctx context.Context, classID uint64, memberName string, date int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	row := tx.QueryRowContext(ctx, "SELECT start_date, end_date, capacity FROM classes WHERE id = ?;", classID)
	var startDate int64
	var endDate int64
	var capacity uint
	if err = row.Scan(&startDate, &endDate, &capacity); err != nil {
		if err == sql.ErrNoRows {
			return RepoError{Message: "Class not found"}
		}
		return err
	}
	if startDate > date || endDate < date {
		return RepoError{Message: "No class is available on the given date"}
	}

	row = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM bookings WHERE class_id = ? AND date = ?;", classID, date)
	var occupancy uint
	err = row.Scan(&occupancy)
	if err != nil {
		return err
	}

	if occupancy >= capacity {
		return RepoError{Message: "Class is full"}
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO bookings (class_id, member_name, date) VALUES (?, ?, ?);", classID, memberName, date); err != nil {
		return err
	}
	return tx.Commit()
}
