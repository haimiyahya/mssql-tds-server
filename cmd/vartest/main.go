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

	// Test 2: Create procedure with DECLARE
	log.Println("\n=== Test 2: Create procedure with DECLARE ===")
	testCreateProcedureWithDeclare(db)

	// Test 3: Execute procedure with DECLARE
	log.Println("\n=== Test 3: Execute procedure with DECLARE ===")
	testExecuteProcedureWithDeclare(db)

	// Test 4: Create procedure with SET
	log.Println("\n=== Test 4: Create procedure with SET ===")
	testCreateProcedureWithSet(db)

	// Test 5: Execute procedure with SET
	log.Println("\n=== Test 5: Execute procedure with SET ===")
	testExecuteProcedureWithSet(db)

	// Test 6: Create procedure with SELECT assignment
	log.Println("\n=== Test 6: Create procedure with SELECT assignment ===")
	testCreateProcedureWithSelectAssignment(db)

	// Test 7: Execute procedure with SELECT assignment
	log.Println("\n=== Test 7: Execute procedure with SELECT assignment ===")
	testExecuteProcedureWithSelectAssignment(db)

	// Test 8: Create complex procedure with multiple variables
	log.Println("\n=== Test 8: Create complex procedure ===")
	testCreateComplexProcedure(db)

	// Test 9: Execute complex procedure
	log.Println("\n=== Test 9: Execute complex procedure ===")
	testExecuteComplexProcedure(db)

	// Test 10: Procedure with parameter and variables
	log.Println("\n=== Test 10: Procedure with parameter and variables ===")
	testProcedureWithParamAndVars(db)

	log.Println("\n=== All tests completed successfully! ===")
}

func testCreateTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE products (id INTEGER, name TEXT, price REAL, stock INT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	// Insert test data
	_, err = db.Exec("INSERT INTO products (id, name, price, stock) VALUES (1, 'Laptop', 999.99, 10)")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	_, err = db.Exec("INSERT INTO products (id, name, price, stock) VALUES (2, 'Mouse', 29.99, 50)")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	_, err = db.Exec("INSERT INTO products (id, name, price, stock) VALUES (3, 'Keyboard', 79.99, 25)")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	log.Println("✓ Created test table with sample data")
}

func testCreateProcedureWithDeclare(db *sql.DB) {
	_, err := db.Exec("CREATE PROCEDURE TestDeclare AS DECLARE @count INT; SELECT @count")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with DECLARE: TestDeclare")
}

func testExecuteProcedureWithDeclare(db *sql.DB) {
	rows, err := db.Query("EXEC TestDeclare")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Read results
	rowCount := 0
	for rows.Next() {
		rowCount++
		var value interface{}
		rows.Scan(&value)
		log.Printf("  Row %d: %v", rowCount, value)
	}

	log.Printf("✓ Executed procedure with DECLARE, returned %d rows", rowCount)
}

func testCreateProcedureWithSet(db *sql.DB) {
	_, err := db.Exec("CREATE PROCEDURE TestSet AS DECLARE @price REAL; SET @price = 999.99; SELECT @price")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with SET: TestSet")
}

func testExecuteProcedureWithSet(db *sql.DB) {
	rows, err := db.Query("EXEC TestSet")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Read results
	rowCount := 0
	for rows.Next() {
		rowCount++
		var value interface{}
		rows.Scan(&value)
		log.Printf("  Row %d: %v", rowCount, value)
	}

	log.Printf("✓ Executed procedure with SET, returned %d rows", rowCount)
}

func testCreateProcedureWithSelectAssignment(db *sql.DB) {
	_, err := db.Exec("CREATE PROCEDURE TestSelectAssign AS DECLARE @total INT; SELECT @total = COUNT(*) FROM products; SELECT @total")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with SELECT assignment: TestSelectAssign")
}

func testExecuteProcedureWithSelectAssignment(db *sql.DB) {
	rows, err := db.Query("EXEC TestSelectAssign")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Read results
	rowCount := 0
	for rows.Next() {
		rowCount++
		var value interface{}
		rows.Scan(&value)
		log.Printf("  Row %d: %v", rowCount, value)
	}

	log.Printf("✓ Executed procedure with SELECT assignment, returned %d rows", rowCount)
}

func testCreateComplexProcedure(db *sql.DB) {
	_, err := db.Exec("CREATE PROCEDURE ComplexProcedure AS DECLARE @id INT; DECLARE @name VARCHAR(50); SET @id = 2; SELECT @name = name FROM products WHERE id = @id; SELECT @name")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created complex procedure: ComplexProcedure")
}

func testExecuteComplexProcedure(db *sql.DB) {
	rows, err := db.Query("EXEC ComplexProcedure")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Read results
	rowCount := 0
	for rows.Next() {
		rowCount++
		var value interface{}
		rows.Scan(&value)
		log.Printf("  Row %d: %v", rowCount, value)
	}

	log.Printf("✓ Executed complex procedure, returned %d rows", rowCount)
}

func testProcedureWithParamAndVars(db *sql.DB) {
	_, err := db.Exec("CREATE PROCEDURE GetProductStock @productId INT AS DECLARE @productName VARCHAR(50); SELECT @productName = name FROM products WHERE id = @productId; SELECT @productName")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with parameter and variables: GetProductStock")

	// Execute with parameter
	rows, err := db.Query("EXEC GetProductStock @productId=1")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Read results
	rowCount := 0
	for rows.Next() {
		rowCount++
		var value interface{}
		rows.Scan(&value)
		log.Printf("  Row %d: %v", rowCount, value)
	}

	log.Printf("✓ Executed procedure with parameter and variables, returned %d rows", rowCount)
}
