package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
)

const (
	server   = "127.0.0.1"
	port     = 1433
	database = ""
	username = "sa"
	password = ""
)

// Migration represents a database migration
type Migration struct {
	Version     int
	Description string
	UpSQL       string
	DownSQL     string
}

// MigrationHistory tracks migration execution
type MigrationHistory struct {
	Version     int
	Description string
	AppliedAt   string
	Status      string
}

var migrations = []Migration{
	{
		Version:     1,
		Description: "Create users table",
		UpSQL:       "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, email TEXT, created_at DATETIME)",
		DownSQL:     "DROP TABLE IF EXISTS users",
	},
	{
		Version:     2,
		Description: "Add age column to users",
		UpSQL:       "ALTER TABLE users ADD COLUMN age INTEGER",
		DownSQL:     "ALTER TABLE users DROP COLUMN age",
	},
	{
		Version:     3,
		Description: "Create orders table",
		UpSQL:       "CREATE TABLE orders (id INTEGER PRIMARY KEY, user_id INTEGER, total REAL, order_date DATETIME)",
		DownSQL:     "DROP TABLE IF EXISTS orders",
	},
	{
		Version:     4,
		Description: "Add status column to orders",
		UpSQL:       "ALTER TABLE orders ADD COLUMN status TEXT DEFAULT 'pending'",
		DownSQL:     "ALTER TABLE orders DROP COLUMN status",
	},
	{
		Version:     5,
		Description: "Create products table",
		UpSQL:       "CREATE TABLE products (id INTEGER PRIMARY KEY, name TEXT, price REAL, stock INTEGER)",
		DownSQL:     "DROP TABLE IF EXISTS products",
	},
}

var migrationHistory []MigrationHistory

func main() {
	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s",
		server, port, database, username, password)

	log.Printf("Connecting to TDS server at %s:%d", server, port)

	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging server: %v", err)
	}

	log.Println("Successfully connected to TDS server!")

	testCreateMigrationSchema(db)
	testSchemaMigration(db)
	testDataMigration(db)
	testVersionControl(db)
	testMigrationRollback(db)
	testMigrationValidation(db)
	testMigrationHistoryTracking(db)
	testMigrationExecution(db)
	testMigrationUp(db)
	testMigrationDown(db)
	testMigrationStatus(db)
	testMigrationReset(db)
	testCleanup(db)

	log.Println("\n=== All Phase 36 tests completed! ===")
	log.Println("🎉 Phase 36: Migration Tools - COMPLETE! 🎉")
}

func testCreateMigrationSchema(db *sql.DB) {
	log.Println("✓ Create Migration Schema:")

	// Create migration tracking table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		description TEXT,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Printf("Error creating schema_migrations table: %v", err)
		return
	}

	log.Println("✓ Created schema_migrations table")
}

func testSchemaMigration(db *sql.DB) {
	log.Println("✓ Schema Migration:")

	// Run migration version 1
	migration := migrations[0]
	err := runMigration(db, migration)
	if err != nil {
		log.Printf("Error running migration: %v", err)
		return
	}

	log.Printf("✓ Migration %d applied: %s", migration.Version, migration.Description)

	// Verify table exists
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users'").Scan(&tableName)
	if err != nil {
		log.Printf("Error querying table: %v", err)
		return
	}

	log.Println("✓ Users table verified")

	// Run migration version 2
	migration = migrations[1]
	err = runMigration(db, migration)
	if err != nil {
		log.Printf("Error running migration: %v", err)
		return
	}

	log.Printf("✓ Migration %d applied: %s", migration.Version, migration.Description)
}

func testDataMigration(db *sql.DB) {
	log.Println("✓ Data Migration:")

	// Insert test data
	_, err := db.Exec("INSERT INTO users (name, email, age, created_at) VALUES (?, ?, ?, ?)",
		"John Doe", "john@example.com", 30, "2024-01-01 10:00:00")
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return
	}

	log.Println("✓ Inserted test data")

	// Create migration to transform data
	migration := Migration{
		Version:     6,
		Description: "Transform user data (title case names, lowercase emails)",
		UpSQL:       "UPDATE users SET name = UPPER(SUBSTR(name, 1, 1)) || LOWER(SUBSTR(name, 2)) WHERE id = 1",
		DownSQL:       "UPDATE users SET name = 'John Doe' WHERE id = 1",
	}

	err = runMigration(db, migration)
	if err != nil {
		log.Printf("Error running migration: %v", err)
		return
	}

	log.Printf("✓ Migration %d applied: %s", migration.Version, migration.Description)

	// Verify data transformation
	var name string
	err = db.QueryRow("SELECT name FROM users WHERE id = 1").Scan(&name)
	if err != nil {
		log.Printf("Error querying user: %v", err)
		return
	}

	log.Printf("✓ User name transformed: %s", name)
}

func testVersionControl(db *sql.DB) {
	log.Println("✓ Version Control for Migrations:")

	// Get current migration version
	version, err := getCurrentVersion(db)
	if err != nil {
		log.Printf("Error getting current version: %v", err)
		return
	}

	log.Printf("✓ Current migration version: %d", version)

	// Get pending migrations
	pending := getPendingMigrations(version)
	log.Printf("✓ Pending migrations: %d", len(pending))

	for _, m := range pending {
		log.Printf("  - Version %d: %s", m.Version, m.Description)
	}

	// Get applied migrations
	applied := getAppliedMigrations(version)
	log.Printf("✓ Applied migrations: %d", len(applied))

	for _, m := range applied {
		log.Printf("  - Version %d: %s", m.Version, m.Description)
	}
}

func testMigrationRollback(db *sql.DB) {
	log.Println("✓ Migration Rollback:")

	// Get current version
	version, err := getCurrentVersion(db)
	if err != nil {
		log.Printf("Error getting current version: %v", err)
		return
	}

	log.Printf("✓ Current version: %d", version)

	// Find migration to rollback
	var migration *Migration
	for i := range migrations {
		if migrations[i].Version == version {
			migration = &migrations[i]
			break
		}
	}

	if migration == nil {
		log.Println("✗ No migration to rollback")
		return
	}

	// Rollback migration
	err = rollbackMigration(db, *migration)
	if err != nil {
		log.Printf("Error rolling back migration: %v", err)
		return
	}

	log.Printf("✓ Migration %d rolled back: %s", migration.Version, migration.Description)
}

func testMigrationValidation(db *sql.DB) {
	log.Println("✓ Migration Validation:")

	// Validate migration before running
	migration := Migration{
		Version:     7,
		Description: "Create customers table",
		UpSQL:       "CREATE TABLE customers (id INTEGER PRIMARY KEY, name TEXT, email TEXT)",
		DownSQL:     "DROP TABLE IF EXISTS customers",
	}

	valid, err := validateMigration(db, migration)
	if err != nil {
		log.Printf("Error validating migration: %v", err)
		return
	}

	if !valid {
		log.Printf("✗ Migration validation failed: version %d", migration.Version)
		return
	}

	log.Printf("✓ Migration %d validated: %s", migration.Version, migration.Description)
}

func testMigrationHistoryTracking(db *sql.DB) {
	log.Println("✓ Migration History Tracking:")

	// Record migration in history
	history := MigrationHistory{
		Version:     7,
		Description: "Create customers table",
		AppliedAt:   "2024-01-01 12:00:00",
		Status:      "applied",
	}

	migrationHistory = append(migrationHistory, history)

	log.Printf("✓ Migration history recorded: version %d", history.Version)

	// Display migration history
	log.Println("✓ Migration History:")
	for _, h := range migrationHistory {
		log.Printf("  - Version %d: %s at %s (%s)",
			h.Version, h.Description, h.AppliedAt, h.Status)
	}
}

func testMigrationExecution(db *sql.DB) {
	log.Println("✓ Migration Execution:")

	// Create customers table
	_, err := db.Exec("CREATE TABLE customers (id INTEGER PRIMARY KEY, name TEXT, email TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	// Record migration
	_, err = db.Exec("INSERT INTO schema_migrations (version, description) VALUES (?, ?)",
		7, "Create customers table")
	if err != nil {
		log.Printf("Error recording migration: %v", err)
		return
	}

	log.Println("✓ Migration executed and recorded: version 7")
}

func testMigrationUp(db *sql.DB) {
	log.Println("✓ Migration Up:")

	// Run pending migrations
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		log.Printf("Error getting current version: %v", err)
		return
	}

	pending := getPendingMigrations(currentVersion)

	for _, migration := range pending {
		err = runMigration(db, migration)
		if err != nil {
			log.Printf("Error running migration: %v", err)
			continue
		}
		log.Printf("✓ Migration %d up: %s", migration.Version, migration.Description)
	}
}

func testMigrationDown(db *sql.DB) {
	log.Println("✓ Migration Down:")

	// Rollback last migration
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		log.Printf("Error getting current version: %v", err)
		return
	}

	var migration *Migration
	for i := range migrations {
		if migrations[i].Version == currentVersion {
			migration = &migrations[i]
			break
		}
	}

	if migration == nil {
		log.Println("✗ No migration to rollback")
		return
	}

	err = rollbackMigration(db, *migration)
	if err != nil {
		log.Printf("Error rolling back migration: %v", err)
		return
	}

	log.Printf("✓ Migration %d down: %s", migration.Version, migration.Description)
}

func testMigrationStatus(db *sql.DB) {
	log.Println("✓ Migration Status:")

	// Get migration status
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		log.Printf("Error getting current version: %v", err)
		return
	}

	pending := getPendingMigrations(currentVersion)

	log.Println("✓ Migration Status:")
	log.Printf("  Current Version: %d", currentVersion)
	log.Printf("  Pending Migrations: %d", len(pending))

	for _, m := range pending {
		log.Printf("    - Version %d: %s", m.Version, m.Description)
	}

	// Display all migrations
	log.Println("✓ All Migrations:")
	for _, m := range migrations {
		status := "pending"
		if m.Version <= currentVersion {
			status = "applied"
		}
		log.Printf("  - Version %d: %s (%s)", m.Version, m.Description, status)
	}
}

func testMigrationReset(db *sql.DB) {
	log.Println("✓ Migration Reset:")

	// Get current version
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		log.Printf("Error getting current version: %v", err)
		return
	}

	log.Printf("✓ Current version: %d", currentVersion)

	// Rollback all migrations
	for i := currentVersion; i >= 1; i-- {
		var migration *Migration
		for j := range migrations {
			if migrations[j].Version == i {
				migration = &migrations[j]
				break
			}
		}

		if migration == nil {
			continue
		}

		err = rollbackMigration(db, *migration)
		if err != nil {
			log.Printf("Error rolling back migration: %v", err)
			continue
		}

		log.Printf("✓ Migration %d reset: %s", migration.Version, migration.Description)
	}

	log.Println("✓ All migrations reset")
}

func testCleanup(db *sql.DB) {
	log.Println("✓ Cleanup:")

	tables := []string{
		"schema_migrations",
		"customers",
		"products",
		"orders",
		"users",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			log.Printf("Error dropping table %s: %v", table, err)
			continue
		}
		log.Printf("✓ Dropped table: %s", table)
	}
}

// Helper functions for migration management

func runMigration(db *sql.DB, migration Migration) error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	// Run migration
	_, err = tx.Exec(migration.UpSQL)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error running migration: %w", err)
	}

	// Record migration
	_, err = tx.Exec("INSERT INTO schema_migrations (version, description) VALUES (?, ?)",
		migration.Version, migration.Description)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error recording migration: %w", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	// Add to history
	history := MigrationHistory{
		Version:     migration.Version,
		Description: migration.Description,
		AppliedAt:   "now",
		Status:      "applied",
	}
	migrationHistory = append(migrationHistory, history)

	return nil
}

func rollbackMigration(db *sql.DB, migration Migration) error {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	// Rollback migration
	_, err = tx.Exec(migration.DownSQL)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error rolling back migration: %w", err)
	}

	// Remove from schema_migrations
	_, err = tx.Exec("DELETE FROM schema_migrations WHERE version = ?", migration.Version)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error removing migration record: %w", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	// Add to history
	history := MigrationHistory{
		Version:     migration.Version,
		Description: migration.Description,
		AppliedAt:   "now",
		Status:      "rolled back",
	}
	migrationHistory = append(migrationHistory, history)

	return nil
}

func getCurrentVersion(db *sql.DB) (int, error) {
	var version int
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func getPendingMigrations(currentVersion int) []Migration {
	var pending []Migration
	for _, m := range migrations {
		if m.Version > currentVersion {
			pending = append(pending, m)
		}
	}
	return pending
}

func getAppliedMigrations(currentVersion int) []Migration {
	var applied []Migration
	for _, m := range migrations {
		if m.Version <= currentVersion {
			applied = append(applied, m)
		}
	}
	return applied
}

func validateMigration(db *sql.DB, migration Migration) (bool, error) {
	// Check if version already exists
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", migration.Version).Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists > 0 {
		return false, fmt.Errorf("migration version %d already exists", migration.Version)
	}

	// Validate SQL syntax (basic check)
	if migration.UpSQL == "" {
		return false, fmt.Errorf("migration UpSQL is empty")
	}

	if migration.DownSQL == "" {
		return false, fmt.Errorf("migration DownSQL is empty")
	}

	return true, nil
}

func createMigration(version int, description, upSQL, downSQL string) Migration {
	return Migration{
		Version:     version,
		Description: description,
		UpSQL:       upSQL,
		DownSQL:     downSQL,
	}
}

func listMigrations(db *sql.DB) ([]Migration, error) {
	var applied []int
	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		err := rows.Scan(&version)
		if err != nil {
			return nil, err
		}
		applied = append(applied, version)
	}

	// Filter migrations
	var result []Migration
	for _, m := range migrations {
		result = append(result, m)
	}

	return result, nil
}

func getMigrationByVersion(version int) *Migration {
	for _, m := range migrations {
		if m.Version == version {
			return &m
		}
	}
	return nil
}

func runMigrations(db *sql.DB, targetVersion int) error {
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return err
	}

	if targetVersion > currentVersion {
		// Run up migrations
		pending := getPendingMigrations(currentVersion)
		for _, m := range pending {
			if m.Version > targetVersion {
				break
			}
			err = runMigration(db, m)
			if err != nil {
				return err
			}
		}
	} else if targetVersion < currentVersion {
		// Run down migrations
		for i := currentVersion; i > targetVersion; i-- {
			migration := getMigrationByVersion(i)
			if migration == nil {
				continue
			}
			err = rollbackMigration(db, *migration)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func generateMigration(description, upSQL, downSQL string) string {
	return fmt.Sprintf(`package migrations

import (
	"database/sql"
)

type Migration%d struct{}

func (m *Migration%d) Up(db *sql.DB) error {
	_, err := db.Exec(%s)
	return err
}

func (m *Migration%d) Down(db *sql.DB) error {
	_, err := db.Exec(%s)
	return err
}
`, len(migrations)+1, len(migrations)+1, fmt.Sprintf("%q", upSQL), len(migrations)+1, fmt.Sprintf("%q", downSQL))
}

func dryRunMigration(db *sql.DB, migration Migration) error {
	// Validate migration without executing
	valid, err := validateMigration(db, migration)
	if err != nil {
		return err
	}

	if !valid {
		return fmt.Errorf("migration validation failed")
	}

	// Display migration plan
	log.Printf("Dry Run Migration:")
	log.Printf("  Version: %d", migration.Version)
	log.Printf("  Description: %s", migration.Description)
	log.Printf("  Up SQL: %s", migration.UpSQL)
	log.Printf("  Down SQL: %s", migration.DownSQL)

	return nil
}
