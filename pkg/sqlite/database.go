package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DefaultDBPath = "./data/tds_server.db"
)

// Database represents the SQLite database connection
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	if dbPath == "" {
		dbPath = DefaultDBPath
	}

	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Connected to SQLite database at %s", dbPath)

	return &Database{db: db}, nil
}

// Initialize creates the necessary tables if they don't exist
func (d *Database) Initialize() error {
	// Create procedures table
	procedureTableSQL := `
	CREATE TABLE IF NOT EXISTS procedures (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		body TEXT NOT NULL,
		parameters TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	if _, err := d.db.Exec(procedureTableSQL); err != nil {
		return fmt.Errorf("failed to create procedures table: %w", err)
	}

	log.Println("Database tables initialized successfully")
	return nil
}

// GetDB returns the underlying database connection
func (d *Database) GetDB() *sql.DB {
	return d.db
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}
