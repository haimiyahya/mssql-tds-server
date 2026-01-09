package sqlexecutor

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/factory/mssql-tds-server/pkg/database"
	"github.com/factory/mssql-tds-server/pkg/sqlparser"
)

// ExecuteCreateDatabase executes a CREATE DATABASE statement
func (e *Executor) ExecuteCreateDatabase(stmt *sqlparser.CreateDatabaseStatement) error {
	// Create database using catalog
	var _ database.Catalog

	// Create database using catalog
	db, err := e.catalog.CreateDatabase(stmt.DatabaseName)
	if err != nil {
		return fmt.Errorf("error creating database '%s': %w", stmt.DatabaseName, err)
	}

	log.Printf("Created database: %s (ID: %d, Path: %s)", db.Name, db.ID, db.FilePath)

	// Open database connection and cache it
	conn, err := sql.Open("sqlite", db.FilePath)
	if err != nil {
		return fmt.Errorf("error opening database '%s': %w", stmt.DatabaseName, err)
	}

	e.connections[stmt.DatabaseName] = conn

	return nil
}

// ExecuteDropDatabase executes a DROP DATABASE statement
func (e *Executor) ExecuteDropDatabase(stmt *sqlparser.DropDatabaseStatement) error {
	// Check if database is currently in use
	if e.currentDBName == stmt.DatabaseName {
		return fmt.Errorf("cannot drop database '%s' because it is currently in use", stmt.DatabaseName)
	}

	// Drop database using catalog
	err := e.catalog.DropDatabase(stmt.DatabaseName)
	if err != nil {
		return fmt.Errorf("error dropping database '%s': %w", stmt.DatabaseName, err)
	}

	// Close and remove connection
	if conn, exists := e.connections[stmt.DatabaseName]; exists {
		conn.Close()
		delete(e.connections, stmt.DatabaseName)
	}

	log.Printf("Dropped database: %s", stmt.DatabaseName)

	return nil
}

// ExecuteUseDatabase executes a USE statement
func (e *Executor) ExecuteUseDatabase(stmt *sqlparser.UseDatabaseStatement) error {
	// Check if database exists
	db, err := e.catalog.GetDatabase(stmt.DatabaseName)
	if err != nil {
		return fmt.Errorf("database '%s' does not exist", stmt.DatabaseName)
	}

	// Close current connection
	if e.currentDB != nil {
		e.currentDB.Close()
	}

	// Open new database connection
	conn, err := sql.Open("sqlite", db.FilePath)
	if err != nil {
		return fmt.Errorf("error opening database '%s': %w", stmt.DatabaseName, err)
	}

	// Set as current database
	e.currentDB = conn
	e.currentDBName = stmt.DatabaseName

	// Cache connection
	e.connections[stmt.DatabaseName] = conn

	log.Printf("Using database: %s (Path: %s)", stmt.DatabaseName, db.FilePath)

	return nil
}

// ExecuteSysDatabasesQuery executes a sys.databases query
func (e *Executor) ExecuteSysDatabasesQuery() (*sql.Rows, error) {
	// Get list of databases from catalog
	databases, err := e.catalog.ListDatabases()
	if err != nil {
		return nil, fmt.Errorf("error listing databases: %w", err)
	}

	// Create in-memory result set (SQLite doesn't support returning rows from Go code directly)
	// We'll simulate this by creating a temporary table or returning formatted data
	// For now, let's create a more elegant solution

	// Store databases in a temporary in-memory SQLite table
	_, err = e.currentDB.Exec(`
		DROP TABLE IF EXISTS _sys_databases_temp;
		CREATE TEMPORARY TABLE _sys_databases_temp (
			name TEXT,
			database_id INTEGER,
			state TEXT,
			create_date TEXT
		);
	`)

	if err != nil {
		return nil, fmt.Errorf("error creating temp table: %w", err)
	}

	// Insert databases
	for _, db := range databases {
		createDate := db.CreateDate.Format("2006-01-02 15:04:05")
		_, err = e.currentDB.Exec(`
			INSERT INTO _sys_databases_temp (name, database_id, state, create_date)
			VALUES (?, ?, ?, ?)
		`, db.Name, db.ID, db.State, createDate)
		if err != nil {
			log.Printf("Error inserting database to temp table: %v", err)
			continue
		}
	}

	// Query the temporary table
	rows, err := e.currentDB.Query(`
		SELECT name, database_id, state, create_date
		FROM _sys_databases_temp
		ORDER BY name ASC
	`)

	if err != nil {
		return nil, fmt.Errorf("error querying sys.databases: %w", err)
	}

	return rows, nil
}

// ExecuteDatabaseCommands executes database-related commands
func (e *Executor) ExecuteDatabaseCommands(stmt *sqlparser.Statement) error {
	switch stmt.Type {
	case sqlparser.StatementTypeCreateDatabase:
		return e.ExecuteCreateDatabase(stmt.CreateDatabase)
	case sqlparser.StatementTypeDropDatabase:
		return e.ExecuteDropDatabase(stmt.DropDatabase)
	case sqlparser.StatementTypeUseDatabase:
		return e.ExecuteUseDatabase(stmt.UseDatabase)
	default:
		return fmt.Errorf("unsupported database command type: %s", stmt.Type)
	}
}

// CreateProcedure creates a stored procedure in the current database
func (e *Executor) CreateProcedure(procName, definition string) error {
	if e.currentDBName == "" {
		return fmt.Errorf("no database selected")
	}

	// Add to catalog
	err := e.catalog.CreateProcedure(e.currentDBName, procName, definition)
	if err != nil {
		return fmt.Errorf("error creating procedure '%s': %w", procName, err)
	}

	// Store procedure in database's sys_objects table
	_, err = e.currentDB.Exec(`
		INSERT OR REPLACE INTO sys_objects (name, type)
		VALUES (?, 'PROCEDURE')
	`, procName)

	if err != nil {
		return fmt.Errorf("error storing procedure: %w", err)
	}

	log.Printf("Created procedure: %s in database: %s", procName, e.currentDBName)

	return nil
}

// CreateFunction creates a function in the current database
func (e *Executor) CreateFunction(funcName, definition, returnType string) error {
	if e.currentDBName == "" {
		return fmt.Errorf("no database selected")
	}

	// Add to catalog
	err := e.catalog.CreateFunction(e.currentDBName, funcName, definition, returnType)
	if err != nil {
		return fmt.Errorf("error creating function '%s': %w", funcName, err)
	}

	// Store function in database's sys_objects table
	_, err = e.currentDB.Exec(`
		INSERT OR REPLACE INTO sys_objects (name, type)
		VALUES (?, 'FUNCTION')
	`, funcName)

	if err != nil {
		return fmt.Errorf("error storing function: %w", err)
	}

	log.Printf("Created function: %s in database: %s", funcName, e.currentDBName)

	return nil
}

// GetProcedures returns procedures for the current database
func (e *Executor) GetProcedures() ([]map[string]interface{}, error) {
	if e.currentDBName == "" {
		return nil, fmt.Errorf("no database selected")
	}

	return e.catalog.GetProcedures(e.currentDBName)
}

// GetFunctions returns functions for the current database
func (e *Executor) GetFunctions() ([]map[string]interface{}, error) {
	if e.currentDBName == "" {
		return nil, fmt.Errorf("no database selected")
	}

	return e.catalog.GetFunctions(e.currentDBName)
}

// TranslateDatabaseReferences translates SQL Server database references to SQLite
func (e *Executor) TranslateDatabaseReferences(query string) string {
	// SQL Server: DatabaseName.schema.Table
	// SQLite: DatabaseName.Table (after ATTACH)

	// Find database references and translate them
	// Pattern: database.schema.table -> database.table

	// This is a simple translation - for production, use proper SQL parsing
	// Remove schema (dbo) from three-part naming
	// database.dbo.table -> database.table

	// Match database.schema.table pattern
	re := regexp.MustCompile(`(\w+)\.dbo\.(\w+)`)
	query = re.ReplaceAllString(query, `$1.$2`)

	// Check for database references
	dbRefs := e.extractDatabaseReferences(query)

	// Attach databases that aren't already attached
	for _, dbName := range dbRefs {
		if !e.attachedDBs[dbName] {
			err := e.attachDatabase(dbName)
			if err != nil {
				log.Printf("Error attaching database %s: %v", dbName, err)
			}
		}
	}

	return query
}

// extractDatabaseReferences extracts database names from query
func (e *Executor) extractDatabaseReferences(query string) []string {
	var dbRefs []string
	seen := make(map[string]bool)

	// Simple pattern matching: find database.table references
	// Match: databaseName.tableName (but not schema.tableName)
	re := regexp.MustCompile(`(\w+)\.(\w+)`)
	matches := re.FindAllStringSubmatch(query, -1)

	for _, match := range matches {
		if len(match) > 1 {
			dbName := match[1]
			// Filter out schema names
			if !isSchemaName(dbName) && !seen[dbName] {
				dbRefs = append(dbRefs, dbName)
				seen[dbName] = true
			}
		}
	}

	return dbRefs
}

// isSchemaName checks if name is a schema name
func isSchemaName(name string) bool {
	schemas := []string{"dbo", "sys", "information_schema"}
	for _, schema := range schemas {
		if strings.EqualFold(name, schema) {
			return true
		}
	}
	return false
}

// attachDatabase attaches a database to the current connection
func (e *Executor) attachDatabase(dbName string) error {
	// Get database from catalog
	db, err := e.catalog.GetDatabase(dbName)
	if err != nil {
		return err
	}

	// Attach using SQLite ATTACH command
	query := fmt.Sprintf("ATTACH DATABASE '%s' AS %s", db.FilePath, dbName)
	_, err = e.currentDB.Exec(query)
	if err != nil {
		return err
	}

	e.attachedDBs[dbName] = true
	log.Printf("Attached database: %s", dbName)

	return nil
}

// GetCurrentDatabase returns the current database name
func (e *Executor) GetCurrentDatabase() string {
	return e.currentDBName
}

// HasCurrentDatabase checks if a database is currently selected
func (e *Executor) HasCurrentDatabase() bool {
	return e.currentDBName != ""
}
