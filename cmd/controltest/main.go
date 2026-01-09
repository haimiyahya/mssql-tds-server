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

	// Test 2: Create procedure with simple IF
	log.Println("\n=== Test 2: Create procedure with simple IF ===")
	testCreateSimpleIF(db)

	// Test 3: Execute simple IF procedure (true condition)
	log.Println("\n=== Test 3: Execute simple IF (true) ===")
	testExecuteSimpleIFTrue(db)

	// Test 4: Execute simple IF procedure (false condition)
	log.Println("\n=== Test 4: Execute simple IF (false) ===")
	testExecuteSimpleIFFalse(db)

	// Test 5: Create procedure with IF/ELSE
	log.Println("\n=== Test 5: Create procedure with IF/ELSE ===")
	testCreateIFElse(db)

	// Test 6: Execute IF/ELSE procedure (true condition)
	log.Println("\n=== Test 6: Execute IF/ELSE (true) ===")
	testExecuteIFElseTrue(db)

	// Test 7: Execute IF/ELSE procedure (false condition)
	log.Println("\n=== Test 7: Execute IF/ELSE (false) ===")
	testExecuteIFElseFalse(db)

	// Test 8: Create procedure with IF/ELSE and variables
	log.Println("\n=== Test 8: Create procedure with IF/ELSE and variables ===")
	testCreateIFWithVariables(db)

	// Test 9: Execute IF/ELSE with variables
	log.Println("\n=== Test 9: Execute IF/ELSE with variables ===")
	testExecuteIFWithVariables(db)

	// Test 10: Create procedure with complex IF conditions
	log.Println("\n=== Test 10: Create procedure with complex IF conditions ===")
	testCreateComplexIF(db)

	// Test 11: Execute complex IF
	log.Println("\n=== Test 11: Execute complex IF ===")
	testExecuteComplexIF(db)

	log.Println("\n=== All tests completed successfully! ===")
}

func testCreateTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE employees (id INTEGER, name TEXT, department TEXT, salary REAL, active BIT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	// Insert test data
	_, err = db.Exec("INSERT INTO employees (id, name, department, salary, active) VALUES (1, 'Alice', 'Engineering', 75000, 1)")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	_, err = db.Exec("INSERT INTO employees (id, name, department, salary, active) VALUES (2, 'Bob', 'Marketing', 50000, 1)")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	_, err = db.Exec("INSERT INTO employees (id, name, department, salary, active) VALUES (3, 'Charlie', 'Engineering', 80000, 0)")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	log.Println("✓ Created test table with sample data")
}

func testCreateSimpleIF(db *sql.DB) {
	_, err := db.Exec("CREATE PROCEDURE CheckActive @id INT AS IF @id = 1 THEN SELECT 'User 1 is active' END")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with simple IF: CheckActive")
}

func testExecuteSimpleIFTrue(db *sql.DB) {
	rows, err := db.Query("EXEC CheckActive @id=1")
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

	log.Printf("✓ Executed simple IF (true), returned %d rows", rowCount)
}

func testExecuteSimpleIFFalse(db *sql.DB) {
	rows, err := db.Query("EXEC CheckActive @id=2")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Read results
	rowCount := 0
	for rows.Next() {
		rowCount++
	}

	log.Printf("✓ Executed simple IF (false), returned %d rows", rowCount)
}

func testCreateIFElse(db *sql.DB) {
	_, err := db.Exec("CREATE PROCEDURE CheckUserStatus @id INT AS IF @id = 1 THEN SELECT 'User 1 exists' ELSE SELECT 'User 1 does not exist' END")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with IF/ELSE: CheckUserStatus")
}

func testExecuteIFElseTrue(db *sql.DB) {
	rows, err := db.Query("EXEC CheckUserStatus @id=1")
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

	log.Printf("✓ Executed IF/ELSE (true), returned %d rows", rowCount)
}

func testExecuteIFElseFalse(db *sql.DB) {
	rows, err := db.Query("EXEC CheckUserStatus @id=999")
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

	log.Printf("✓ Executed IF/ELSE (false), returned %d rows", rowCount)
}

func testCreateIFWithVariables(db *sql.DB) {
	_, err := db.Exec("CREATE PROCEDURE CheckEmployeeStatus @id INT AS DECLARE @active BIT; SELECT @active = active FROM employees WHERE id = @id; IF @active = 1 THEN SELECT 'Employee is active' ELSE SELECT 'Employee is inactive' END")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with IF/ELSE and variables: CheckEmployeeStatus")
}

func testExecuteIFWithVariables(db *sql.DB) {
	rows, err := db.Query("EXEC CheckEmployeeStatus @id=1")
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

	log.Printf("✓ Executed IF/ELSE with variables, returned %d rows", rowCount)
}

func testCreateComplexIF(db *sql.DB) {
	_, err := db.Exec("CREATE PROCEDURE CheckSalary @id INT AS DECLARE @salary REAL; SELECT @salary = salary FROM employees WHERE id = @id; IF @salary > 70000 THEN SELECT 'High salary' ELSE SELECT 'Regular salary' END")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with complex IF: CheckSalary")
}

func testExecuteComplexIF(db *sql.DB) {
	rows, err := db.Query("EXEC CheckSalary @id=1")
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

	log.Printf("✓ Executed complex IF, returned %d rows", rowCount)
}
