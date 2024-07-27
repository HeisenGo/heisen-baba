package db_helper

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	return nil
}

// ValidateCityCountry checks if the given city and country combination exists in the database
func ValidateCityCountry(city, country string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM worldcities WHERE city = ? AND country = ?", city, country).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("database query failed: %w", err)
	}
	return count > 0, nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}