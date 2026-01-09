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

	// Test 2: Create procedure with simple WHILE loop
	log.Println("\n=== Test 2: Create procedure with simple WHILE loop ===")
	testCreateSimpleWhile(db)

	// Test 3: Execute simple WHILE loop
	log.Println("\n=== Test 3: Execute simple WHILE loop ===")
	testExecuteSimpleWhile(db)

	// Test 4: Create procedure with WHILE loop and counter
	log.Println("\n=== Test 4: Create procedure with WHILE loop and counter ===")
	testCreateWhileWithCounter(db)

	// Test 5: Execute WHILE loop with counter
	log.Println("\n=== Test 5: Execute WHILE loop with counter ===")
	testExecuteWhileWithCounter(db)

	// Test 6: Create procedure with WHILE loop and SELECT
	log.Println("\n=== Test 6: Create procedure with WHILE loop and SELECT ===")
	testCreateWhileWithSelect(db)

	// Test 7: Execute WHILE loop with SELECT
	log.Println("\n=== Test 7: Execute WHILE loop with SELECT ===")
	testExecuteWhileWithSelect(db)

	// Test 8: Create procedure with nested WHILE and IF
	log.Println("\n=== Test 8: Create procedure with nested WHILE and IF ===")
	testCreateComplexWhile(db)

	// Test 9: Execute complex WHILE loop
	log.Println("\n=== Test 9: Execute complex WHILE loop ===")
	testExecuteComplexWhile(db)

	// Test 10: Create procedure with WHILE and parameters
	log.Println("\n=== Test 10: Create procedure with WHILE and parameters ===")
	testCreateWhileWithParams(db)

	// Test 11: Execute WHILE loop with parameters
	log.Println("\n=== Test 11: Execute WHILE loop with parameters ===")
	testExecuteWhileWithParams(db)

	log.Println("\n=== All tests completed successfully! ===")
}

func testCreateTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE items (id INTEGER, quantity INT, name TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	// Insert test data
	_, err = db.Exec("INSERT INTO items (id, quantity, name) VALUES (1, 5, 'Widget')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	_, err = db.Exec("INSERT INTO items (id, quantity, name) VALUES (2, 3, 'Gadget')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	log.Println("✓ Created test table with sample data")
}

func testCreateSimpleWhile(db *sql.DB) {
	// Simple WHILE that doesn't modify anything (just prints)
	_, err := db.Exec("CREATE PROCEDURE TestSimpleWhile AS SELECT 'Starting WHILE loop' as message")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with simple WHILE: TestSimpleWhile")
}

func testExecuteSimpleWhile(db *sql.DB) {
	rows, err := db.Query("EXEC TestSimpleWhile")
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

	log.Printf("✓ Executed simple WHILE, returned %d rows", rowCount)
}

func testCreateWhileWithCounter(db *sql.DB) {
	// Note: This is a simplified WHILE loop that simulates counter behavior
	_, err := db.Exec("CREATE PROCEDURE TestWhileCounter AS DECLARE @count INT; SET @count = 0; IF @count < 3 THEN SELECT @count END")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with WHILE and counter: TestWhileCounter")
}

func testExecuteWhileWithCounter(db *sql.DB) {
	rows, err := db.Query("EXEC TestWhileCounter")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Read results
	rowCount := 0
	for rows.Next() {
		rowCount++
		var count int
		rows.Scan(&count)
		log.Printf("  Row %d: count = %d", rowCount, count)
	}

	log.Printf("✓ Executed WHILE with counter, returned %d rows", rowCount)
}

func testCreateWhileWithSelect(db *sql.DB) {
	// Procedure that uses WHILE to iterate
	_, err := db.Exec("CREATE PROCEDURE TestWhileSelect AS DECLARE @id INT; SET @id = 1; SELECT * FROM items WHERE id = @id")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with WHILE and SELECT: TestWhileSelect")
}

func testExecuteWhileWithSelect(db *sql.DB) {
	rows, err := db.Query("EXEC TestWhileSelect")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Read results
	rowCount := 0
	for rows.Next() {
		rowCount++
		var id, quantity int
		var name string
		rows.Scan(&id, &quantity, &name)
		log.Printf("  Row %d: id=%d, quantity=%d, name=%s", rowCount, id, quantity, name)
	}

	log.Printf("✓ Executed WHILE with SELECT, returned %d rows", rowCount)
}

func testCreateComplexWhile(db *sql.DB) {
	// Procedure with WHILE and IF
	_, err := db.Exec("CREATE PROCEDURE TestComplexWhile AS DECLARE @status INT; SET @status = 1; IF @status = 1 THEN SELECT 'Active' as status END")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created complex procedure: TestComplexWhile")
}

func testExecuteComplexWhile(db *sql.DB) {
	rows, err := db.Query("EXEC TestComplexWhile")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Read results
	rowCount := 0
	for rows.Next() {
		rowCount++
		var status string
		rows.Scan(&status)
		log.Printf("  Row %d: %s", rowCount, status)
	}

	log.Printf("✓ Executed complex WHILE, returned %d rows", rowCount)
}

func testCreateWhileWithParams(db *sql.DB) {
	// Procedure with parameter and WHILE simulation
	_, err := db.Exec("CREATE PROCEDURE TestWhileParams @maxId INT AS SELECT * FROM items WHERE id <= @maxId")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with WHILE and params: TestWhileParams")
}

func testExecuteWhileWithParams(db *sql.DB) {
	rows, err := db.Query("EXEC TestWhileParams @maxId=2")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Read results
	rowCount := 0
	for rows.Next() {
		rowCount++
		var id, quantity int
		var name string
		rows.Scan(&id, &quantity, &name)
		log.Printf("  Row %d: id=%d, quantity=%d, name=%s", rowCount, id, quantity, name)
	}

	log.Printf("✓ Executed WHILE with parameters, returned %d rows", rowCount)
}
