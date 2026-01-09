package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Database represents a database
type Database struct {
	ID        int
	Name      string
	State     string
	CreateDate time.Time
	FilePath  string
	IsSystem  bool
}

// Catalog represents database catalog
type Catalog struct {
	masterDB *sql.DB
	dataDir  string
}

// NewCatalog creates a new database catalog
func NewCatalog(dataDir string, masterDB *sql.DB) *Catalog {
	// Create data directory if it doesn't exist
	os.MkdirAll(dataDir, 0755)

	// Create master database if it doesn't exist
	masterPath := filepath.Join(dataDir, "master.db")
	if _, err := os.Stat(masterPath); os.IsNotExist(err) {
		db, err := sql.Open("sqlite", masterPath)
		if err != nil {
			fmt.Printf("Error creating master database: %v\n", err)
			return nil
		}
		defer db.Close()

		// Create system tables
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS sys_databases (
				database_id INTEGER PRIMARY KEY,
				name TEXT NOT NULL UNIQUE,
				state TEXT DEFAULT 'ONLINE',
				create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
				file_path TEXT,
				is_system BOOLEAN DEFAULT 0
			);

			CREATE TABLE IF NOT EXISTS sys_procedures (
				procedure_id INTEGER PRIMARY KEY,
				database_id INTEGER NOT NULL,
				name TEXT NOT NULL,
				definition TEXT NOT NULL,
				create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (database_id) REFERENCES sys_databases(database_id)
			);

			CREATE TABLE IF NOT EXISTS sys_functions (
				function_id INTEGER PRIMARY KEY,
				database_id INTEGER NOT NULL,
				name TEXT NOT NULL,
				definition TEXT NOT NULL,
				return_type TEXT,
				create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (database_id) REFERENCES sys_databases(database_id)
			);

			-- Insert system databases
			INSERT INTO sys_databases (database_id, name, state, is_system)
			VALUES
				(1, 'master', 'ONLINE', 1),
				(2, 'tempdb', 'ONLINE', 1),
				(3, 'model', 'ONLINE', 1),
				(4, 'msdb', 'ONLINE', 1)
			ON CONFLICT(database_id) DO NOTHING;
		`)
		if err != nil {
			fmt.Printf("Error creating system tables: %v\n", err)
			return nil
		}
	}

	return &Catalog{
		masterDB: masterDB,
		dataDir:  dataDir,
	}
}

// ListDatabases returns list of all databases
func (c *Catalog) ListDatabases() ([]Database, error) {
	query := `
		SELECT database_id, name, state, create_date, file_path, is_system
		FROM sys_databases
		ORDER BY is_system DESC, name ASC
	`

	rows, err := c.masterDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying databases: %w", err)
	}
	defer rows.Close()

	var databases []Database
	for rows.Next() {
		var db Database
		var createDate string
		err := rows.Scan(
			&db.ID,
			&db.Name,
			&db.State,
			&createDate,
			&db.FilePath,
			&db.IsSystem,
		)
		if err != nil {
			continue
		}

		// Parse create date
		db.CreateDate, _ = time.Parse("2006-01-02 15:04:05", createDate)

		databases = append(databases, db)
	}

	return databases, nil
}

// ListDatabasesFromFilesystem scans filesystem for .db files
func (c *Catalog) ListDatabasesFromFilesystem() ([]Database, error) {
	var databases []Database

	// Scan directory for .db files
	files, err := os.ReadDir(c.dataDir)
	if err != nil {
		return nil, fmt.Errorf("error reading data directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".db") {
			dbName := strings.TrimSuffix(file.Name(), ".db")

			// Get file info
			info, err := file.Info()
			if err != nil {
				continue
			}

			databases = append(databases, Database{
				Name:      dbName,
				State:     "ONLINE",
				CreateDate: info.ModTime(),
				FilePath:  filepath.Join(c.dataDir, file.Name()),
				IsSystem:  isSystemDatabase(dbName),
			})
		}
	}

	// Sort: system databases first, then by name
	sort.Slice(databases, func(i, j int) bool {
		if databases[i].IsSystem != databases[j].IsSystem {
			return databases[i].IsSystem
		}
		return databases[i].Name < databases[j].Name
	})

	// Assign database IDs
	for i := range databases {
		databases[i].ID = i + 1
	}

	return databases, nil
}

// CreateDatabase creates a new database
func (c *Catalog) CreateDatabase(dbName string) (*Database, error) {
	// Validate database name
	if !isValidDatabaseName(dbName) {
		return nil, fmt.Errorf("invalid database name '%s'", dbName)
	}

	// Check if system database
	if isSystemDatabase(dbName) {
		return nil, fmt.Errorf("cannot create system database '%s'", dbName)
	}

	// Check if already exists
	dbPath := filepath.Join(c.dataDir, dbName+".db")
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		return nil, fmt.Errorf("database '%s' already exists", dbName)
	}

	// Create new database file
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error creating database file: %w", err)
	}
	defer db.Close()

	// Create system tables in new database
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sys_objects (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			create_date DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS sys_columns (
			id INTEGER PRIMARY KEY,
			object_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			is_nullable BOOLEAN DEFAULT 1,
			FOREIGN KEY (object_id) REFERENCES sys_objects(id)
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("error creating system tables: %w", err)
	}

	// Add to catalog
	query := `
		INSERT INTO sys_databases (name, state, file_path, create_date)
		VALUES (?, 'ONLINE', ?, CURRENT_TIMESTAMP)
	`
	result, err := c.masterDB.Exec(query, dbName, dbPath)
	if err != nil {
		return nil, fmt.Errorf("error adding to catalog: %w", err)
	}

	// Get database ID
	id, _ := result.LastInsertId()

	database := &Database{
		ID:        int(id),
		Name:      dbName,
		State:     "ONLINE",
		CreateDate: time.Now(),
		FilePath:  dbPath,
		IsSystem:  false,
	}

	return database, nil
}

// DropDatabase drops a database
func (c *Catalog) DropDatabase(dbName string) error {
	// Check if system database
	if isSystemDatabase(dbName) {
		return fmt.Errorf("cannot drop system database '%s'", dbName)
	}

	// Get database path from catalog
	var dbPath string
	err := c.masterDB.QueryRow(
		"SELECT file_path FROM sys_databases WHERE name = ?",
		dbName,
	).Scan(&dbPath)
	if err != nil {
		return fmt.Errorf("database '%s' not found", dbName)
	}

	// Delete database file
	err = os.Remove(dbPath)
	if err != nil {
		return fmt.Errorf("error deleting database file: %w", err)
	}

	// Remove from catalog
	query := `DELETE FROM sys_databases WHERE name = ?`
	_, err = c.masterDB.Exec(query, dbName)
	if err != nil {
		return fmt.Errorf("error removing from catalog: %w", err)
	}

	// Drop procedures and functions
	c.masterDB.Exec("DELETE FROM sys_procedures WHERE database_id = (SELECT database_id FROM sys_databases WHERE name = ?)", dbName)
	c.masterDB.Exec("DELETE FROM sys_functions WHERE database_id = (SELECT database_id FROM sys_databases WHERE name = ?)", dbName)

	return nil
}

// GetDatabase returns a database by name
func (c *Catalog) GetDatabase(dbName string) (*Database, error) {
	query := `
		SELECT database_id, name, state, create_date, file_path, is_system
		FROM sys_databases
		WHERE name = ?
	`

	var db Database
	var createDate string
	err := c.masterDB.QueryRow(query, dbName).Scan(
		&db.ID,
		&db.Name,
		&db.State,
		&createDate,
		&db.FilePath,
		&db.IsSystem,
	)

	if err != nil {
		return nil, fmt.Errorf("database '%s' not found", dbName)
	}

	// Parse create date
	db.CreateDate, _ = time.Parse("2006-01-02 15:04:05", createDate)

	return &db, nil
}

// OpenDatabase opens a database connection
func (c *Catalog) OpenDatabase(dbName string) (*sql.DB, error) {
	// Get database from catalog
	db, err := c.GetDatabase(dbName)
	if err != nil {
		return nil, err
	}

	// Open database
	conn, err := sql.Open("sqlite", db.FilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	return conn, nil
}

// CreateProcedure creates a stored procedure in a database
func (c *Catalog) CreateProcedure(dbName, procName, definition string) error {
	// Get database ID
	var dbID int
	err := c.masterDB.QueryRow(
		"SELECT database_id FROM sys_databases WHERE name = ?",
		dbName,
	).Scan(&dbID)
	if err != nil {
		return fmt.Errorf("database '%s' not found", dbName)
	}

	// Insert procedure
	query := `
		INSERT INTO sys_procedures (database_id, name, definition, create_date)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(database_id, name) DO UPDATE SET definition = excluded.definition
	`
	_, err = c.masterDB.Exec(query, dbID, procName, definition)
	return err
}

// GetProcedures returns procedures for a database
func (c *Catalog) GetProcedures(dbName string) ([]map[string]interface{}, error) {
	query := `
		SELECT name, definition, create_date
		FROM sys_procedures
		WHERE database_id = (SELECT database_id FROM sys_databases WHERE name = ?)
		ORDER BY name ASC
	`

	rows, err := c.masterDB.Query(query, dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var procedures []map[string]interface{}
	for rows.Next() {
		var name, definition, createDate string
		err := rows.Scan(&name, &definition, &createDate)
		if err != nil {
			continue
		}

		procedures = append(procedures, map[string]interface{}{
			"name":        name,
			"definition":  definition,
			"create_date": createDate,
		})
	}

	return procedures, nil
}

// CreateFunction creates a function in a database
func (c *Catalog) CreateFunction(dbName, funcName, definition, returnType string) error {
	// Get database ID
	var dbID int
	err := c.masterDB.QueryRow(
		"SELECT database_id FROM sys_databases WHERE name = ?",
		dbName,
	).Scan(&dbID)
	if err != nil {
		return fmt.Errorf("database '%s' not found", dbName)
	}

	// Insert function
	query := `
		INSERT INTO sys_functions (database_id, name, definition, return_type, create_date)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(database_id, name) DO UPDATE SET definition = excluded.definition, return_type = excluded.return_type
	`
	_, err = c.masterDB.Exec(query, dbID, funcName, definition, returnType)
	return err
}

// GetFunctions returns functions for a database
func (c *Catalog) GetFunctions(dbName string) ([]map[string]interface{}, error) {
	query := `
		SELECT name, definition, return_type, create_date
		FROM sys_functions
		WHERE database_id = (SELECT database_id FROM sys_databases WHERE name = ?)
		ORDER BY name ASC
	`

	rows, err := c.masterDB.Query(query, dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var functions []map[string]interface{}
	for rows.Next() {
		var name, definition, returnType, createDate string
		err := rows.Scan(&name, &definition, &returnType, &createDate)
		if err != nil {
			continue
		}

		functions = append(functions, map[string]interface{}{
			"name":        name,
			"definition":  definition,
			"return_type": returnType,
			"create_date": createDate,
		})
	}

	return functions, nil
}

// isValidDatabaseName validates database name
func isValidDatabaseName(name string) bool {
	if name == "" {
		return false
	}

	// Check for invalid characters
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return false
		}
	}

	return true
}

// isSystemDatabase checks if database is a system database
func isSystemDatabase(name string) bool {
	systemDatabases := []string{"master", "tempdb", "model", "msdb"}
	for _, sys := range systemDatabases {
		if strings.EqualFold(name, sys) {
			return true
		}
	}
	return false
}
