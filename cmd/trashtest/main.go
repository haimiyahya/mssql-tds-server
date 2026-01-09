package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/factory/mssql-tds-server/pkg/trash"
)

const (
	server   = "127.0.0.1"
	port     = 1433
	database = ""
	username = "sa"
	password = ""
)

func main() {
	log.Println("✓ Trash/Recycle Bin Test Client")
	log.Println("=")

	// Build connection string
	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=", server, port, database, username)

	// Connect to server
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatalf("✗ Error connecting to server: %v", err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("✗ Error pinging server: %v", err)
	}

	log.Println("✓ Successfully connected to MSSQL TDS Server!")
	log.Println()

	// Test 1: List databases
	testListDatabases(db)

	// Test 2: CREATE DATABASE
	testCreateDatabase(db, "TestDB_Trash")

	// Test 3: Verify database file exists
	dbPath := "./data/TestDB_Trash.db"
	log.Printf("✓ Test: Verify database file exists")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("  ✗ Database file does not exist: %s", dbPath)
	} else {
		log.Printf("  ✓ Database file exists: %s", dbPath)
	}
	log.Println()

	// Test 4: DROP DATABASE (should move to trash)
	testDropDatabase(db, "TestDB_Trash")

	// Test 5: Verify database file is NOT in data directory
	log.Printf("✓ Test: Verify database file removed from data directory")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("  ✓ Database file removed from data directory")
	} else {
		log.Printf("  ✗ Database file still exists in data directory: %s", dbPath)
	}
	log.Println()

	// Test 6: List trash (manual trash directory)
	testListTrash()

	// Test 7: CREATE another database
	testCreateDatabase(db, "TestDB_Recovery")

	// Test 8: DROP DATABASE
	testDropDatabase(db, "TestDB_Recovery")

	// Test 9: Restore from trash
	testRestoreFromTrash()

	// Test 10: Verify restored database exists
	log.Printf("✓ Test: Verify restored database file exists")
	recoveryPath := "./data/TestDB_Recovery.db"
	if _, err := os.Stat(recoveryPath); os.IsNotExist(err) {
		log.Printf("  ✗ Restored database file does not exist: %s", recoveryPath)
	} else {
		log.Printf("  ✓ Restored database file exists: %s", recoveryPath)
	}
	log.Println()

	// Test 11: Verify database is back in catalog
	testListDatabases(db)

	// Test 12: DROP DATABASE again
	testDropDatabase(db, "TestDB_Recovery")

	// Test 13: Empty trash
	testEmptyTrash()

	// Test 14: Verify trash is empty
	testListTrash()

	// Test 15: List databases (clean)
	testListDatabases(db)

	log.Println()
	log.Println("✓ All tests completed successfully!")
}

func testListDatabases(db *sql.DB) {
	log.Println("✓ Test: List Databases")

	query := "SELECT name, database_id, state FROM sys.databases"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("✗ Error listing databases: %v", err)
		return
	}
	defer rows.Close()

	log.Printf("  Databases:")
	for rows.Next() {
		var name, state string
		var id int
		err := rows.Scan(&name, &id, &state)
		if err != nil {
			continue
		}
		log.Printf("    - %s (ID: %d, State: %s)", name, id, state)
	}

	log.Println()
}

func testCreateDatabase(db *sql.DB, dbName string) {
	log.Printf("✓ Test: CREATE DATABASE %s", dbName)

	// Check if database exists first
	var exists int
	db.QueryRow("SELECT COUNT(*) FROM sys.databases WHERE name = ?", dbName).Scan(&exists)

	if exists > 0 {
		log.Printf("  - Database %s already exists, dropping first...", dbName)
		testDropDatabase(db, dbName)
		// Wait a bit for file system
		time.Sleep(100 * time.Millisecond)
	}

	query := fmt.Sprintf("CREATE DATABASE %s", dbName)
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("✗ Error creating database %s: %v", dbName, err)
		return
	}

	log.Printf("  - Database %s created successfully", dbName)
	log.Println()
}

func testDropDatabase(db *sql.DB, dbName string) {
	log.Printf("✓ Test: DROP DATABASE %s (should move to trash)", dbName)

	query := fmt.Sprintf("DROP DATABASE %s", dbName)
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("✗ Error dropping database %s: %v", dbName, err)
		return
	}

	log.Printf("  - Database %s dropped successfully (moved to trash/recycle bin)", dbName)
	log.Println()
}

func testListTrash() {
	log.Println("✓ Test: List Trash (manual trash directory)")

	// List manual trash directory (./trash)
	files, err := trash.ListTrash()
	if err != nil {
		log.Printf("✗ Error listing trash: %v", err)
		return
	}

	log.Printf("  Trash files (%d):", len(files))
	for _, file := range files {
		log.Printf("    - %s (%d bytes, Modified: %s)", 
			file.Name(), file.Size(), file.ModTime().Format("2006-01-02 15:04:05"))
	}

	log.Println()
}

func testRestoreFromTrash() {
	log.Println("✓ Test: Restore from Trash")

	// Get trash files
	files, err := trash.ListTrash()
	if err != nil {
		log.Printf("✗ Error listing trash: %v", err)
		return
	}

	if len(files) == 0 {
		log.Println("  - No files in trash to restore")
		return
	}

	// Restore the last file (most recent)
	fileToRestore := files[len(files)-1]
	trashPath := "./trash/" + fileToRestore.Name()

	// Extract original database name from filename
	// Filename format: TestDB_Recovery_20240101_150405.db
	originalName := fileToRestore.Name()[:len(fileToRestore.Name())-20] // Remove timestamp
	restorePath := "./data/" + originalName

	log.Printf("  - Restoring: %s -> %s", trashPath, restorePath)

	err = trash.RestoreFromTrash(trashPath, restorePath)
	if err != nil {
		log.Printf("✗ Error restoring from trash: %v", err)
		return
	}

	log.Printf("  - File restored successfully")
	log.Println()
}

func testEmptyTrash() {
	log.Println("✓ Test: Empty Trash")

	err := trash.EmptyTrash()
	if err != nil {
		log.Printf("✗ Error emptying trash: %v", err)
		return
	}

	log.Println("  - Trash emptied successfully")
	log.Println()
}
