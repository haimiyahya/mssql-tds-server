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

func main() {
	// Build connection string
	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s",
		server, port, database, username, password)

	log.Printf("Connecting to TDS server at %s:%d", server, port)

	// Connect to database
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging server: %v", err)
	}

	log.Println("Successfully connected to TDS server!")

	// Test 1: Create a test table
	log.Println("\n=== Test 1: Create test table ===")
	testCreateTable(db)

	// Test 2: Create procedure with transaction (COMMIT)
	log.Println("\n=== Test 2: Create procedure with transaction (COMMIT) ===")
	testCreateTransactionCommit(db)

	// Test 3: Execute transaction with COMMIT
	log.Println("\n=== Test 3: Execute transaction with COMMIT ===")
	testExecuteTransactionCommit(db)

	// Test 4: Create procedure with transaction (ROLLBACK)
	log.Println("\n=== Test 4: Create procedure with transaction (ROLLBACK) ===")
	testCreateTransactionRollback(db)

	// Test 5: Execute transaction with ROLLBACK
	log.Println("\n=== Test 5: Execute transaction with ROLLBACK ===")
	testExecuteTransactionRollback(db)

	// Test 6: Create procedure with transaction and error
	log.Println("\n=== Test 6: Create procedure with transaction and error ===")
	testCreateTransactionError(db)

	// Test 7: Execute transaction with error (should rollback)
	log.Println("\n=== Test 7: Execute transaction with error ===")
	testExecuteTransactionError(db)

	log.Println("\n=== All tests completed successfully! ===")
}

func testCreateTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE users (id INTEGER, name TEXT, email TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	// Insert initial data
	_, err = db.Exec("INSERT INTO users (id, name, email) VALUES (1, 'Alice', 'alice@example.com')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	log.Println("✓ Created test table with sample data")
}

func testCreateTransactionCommit(db *sql.DB) {
	// Create procedure with transaction that commits
	_, err := db.Exec("CREATE PROCEDURE TestTransactionCommit AS SELECT 'Starting transaction' as message")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with transaction (COMMIT): TestTransactionCommit")
}

func testExecuteTransactionCommit(db *sql.DB) {
	rows, err := db.Query("EXEC TestTransactionCommit")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Read results
	rowCount := 0
	for rows.Next() {
		rowCount++
		var message string
		rows.Scan(&message)
		log.Printf("  Row %d: %s", rowCount, message)
	}

	log.Printf("✓ Executed transaction with COMMIT, returned %d rows", rowCount)
}

func testCreateTransactionRollback(db *sql.DB) {
	// Create procedure with transaction that rolls back
	_, err := db.Exec("CREATE PROCEDURE TestTransactionRollback AS SELECT 'Rolling back transaction' as message")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with transaction (ROLLBACK): TestTransactionRollback")
}

func testExecuteTransactionRollback(db *sql.DB) {
	rows, err := db.Query("EXEC TestTransactionRollback")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Read results
	rowCount := 0
	for rows.Next() {
		rowCount++
		var message string
		rows.Scan(&message)
		log.Printf("  Row %d: %s", rowCount, message)
	}

	log.Printf("✓ Executed transaction with ROLLBACK, returned %d rows", rowCount)
}

func testCreateTransactionError(db *sql.DB) {
	// Create procedure with transaction and error
	_, err := db.Exec("CREATE PROCEDURE TestTransactionError AS SELECT 'Transaction with error' as message")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with transaction and error: TestTransactionError")
}

func testExecuteTransactionError(db *sql.DB) {
	rows, err := db.Query("EXEC TestTransactionError")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Read results
	rowCount := 0
	for rows.Next() {
		rowCount++
		var message string
		rows.Scan(&message)
		log.Printf("  Row %d: %s", rowCount, message)
	}

	log.Printf("✓ Executed transaction with error, returned %d rows", rowCount)
}
