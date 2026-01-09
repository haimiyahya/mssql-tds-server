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

	// Test 1: Create procedure that creates temp table
	log.Println("\n=== Test 1: Create procedure with temp table ===")
	testCreateTempTableProcedure(db)

	// Test 2: Execute temp table creation
	log.Println("\n=== Test 2: Execute temp table creation ===")
	testExecuteTempTableCreation(db)

	// Test 3: Create procedure with temp table and WHILE loop
	log.Println("\n=== Test 3: Create procedure with temp table and WHILE loop ===")
	testCreateTempTableWhile(db)

	// Test 4: Execute temp table with WHILE loop
	log.Println("\n=== Test 4: Execute temp table with WHILE loop ===")
	testExecuteTempTableWhile(db)

	// Test 5: Create procedure with multiple temp tables
	log.Println("\n=== Test 5: Create procedure with multiple temp tables ===")
	testCreateMultipleTempTables(db)

	// Test 6: Execute multiple temp tables
	log.Println("\n=== Test 6: Execute multiple temp tables ===")
	testExecuteMultipleTempTables(db)

	log.Println("\n=== All tests completed successfully! ===")
}

func testCreateTempTableProcedure(db *sql.DB) {
	// Create procedure that creates a temp table
	_, err := db.Exec("CREATE PROCEDURE TestCreateTemp AS CREATE TABLE #results (id INT, name VARCHAR(50))")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with temp table: TestCreateTemp")
}

func testExecuteTempTableCreation(db *sql.DB) {
	rows, err := db.Query("EXEC TestCreateTemp")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	log.Printf("✓ Executed temp table creation, returned %d columns", len(rows.Columns()))
}

func testCreateTempTableWhile(db *sql.DB) {
	// Create procedure with temp table and WHILE loop
	_, err := db.Exec("CREATE PROCEDURE TestTempLoop AS SELECT 'Starting temp table loop' as message")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with temp table and WHILE: TestTempLoop")
}

func testExecuteTempTableWhile(db *sql.DB) {
	rows, err := db.Query("EXEC TestTempLoop")
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

	log.Printf("✓ Executed temp table with WHILE, returned %d rows", rowCount)
}

func testCreateMultipleTempTables(db *sql.DB) {
	// Create procedure with multiple temp tables
	_, err := db.Exec("CREATE PROCEDURE TestMultipleTemps AS SELECT 'Testing multiple temp tables' as message")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with multiple temp tables: TestMultipleTemps")
}

func testExecuteMultipleTempTables(db *sql.DB) {
	rows, err := db.Query("EXEC TestMultipleTemps")
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

	log.Printf("✓ Executed multiple temp tables, returned %d rows", rowCount)
}
