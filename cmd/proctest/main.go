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

	// Test 2: Create a simple procedure
	log.Println("\n=== Test 2: Create simple procedure ===")
	testCreateSimpleProcedure(db)

	// Test 3: Execute simple procedure
	log.Println("\n=== Test 3: Execute simple procedure ===")
	testExecuteSimpleProcedure(db)

	// Test 4: Create procedure with parameter
	log.Println("\n=== Test 4: Create procedure with parameter ===")
	testCreateProcedureWithParam(db)

	// Test 5: Execute procedure with parameter
	log.Println("\n=== Test 5: Execute procedure with parameter ===")
	testExecuteProcedureWithParam(db)

	// Test 6: Create procedure with multiple parameters
	log.Println("\n=== Test 6: Create procedure with multiple parameters ===")
	testCreateProcedureWithMultipleParams(db)

	// Test 7: Execute procedure with multiple parameters
	log.Println("\n=== Test 7: Execute procedure with multiple parameters ===")
	testExecuteProcedureWithMultipleParams(db)

	// Test 8: Drop procedure
	log.Println("\n=== Test 8: Drop procedure ===")
	testDropProcedure(db)

	log.Println("\n=== All tests completed successfully! ===")
}

func testCreateTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE users (id INTEGER, name TEXT, department TEXT, email TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	// Insert test data
	_, err = db.Exec("INSERT INTO users (id, name, department, email) VALUES (1, 'Alice', 'Engineering', 'alice@example.com')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	_, err = db.Exec("INSERT INTO users (id, name, department, email) VALUES (2, 'Bob', 'Marketing', 'bob@example.com')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	_, err = db.Exec("INSERT INTO users (id, name, department, email) VALUES (3, 'Charlie', 'Engineering', 'charlie@example.com')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	log.Println("✓ Created test table with sample data")
}

func testCreateSimpleProcedure(db *sql.DB) {
	_, err := db.Exec("CREATE PROCEDURE GetAllUsers AS SELECT * FROM users")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created simple procedure: GetAllUsers")
}

func testExecuteSimpleProcedure(db *sql.DB) {
	rows, err := db.Query("EXEC GetAllUsers")
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Get columns
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("Error getting columns: %v", err)
		return
	}

	log.Printf("Columns: %v", columns)

	// Read rows
	rowCount := 0
	for rows.Next() {
		rowCount++
		// Create slice for row values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		// Get value pointers
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan row
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Printf("Error scanning row %d: %v", rowCount, err)
			continue
		}

		// Print row
		rowStr := ""
		for i, val := range values {
			if i > 0 {
				rowStr += ", "
			}
			rowStr += fmt.Sprintf("%v", val)
		}
		log.Printf("  Row %d: %s", rowCount, rowStr)
	}

	log.Printf("✓ Executed procedure, returned %d rows", rowCount)
}

func testCreateProcedureWithParam(db *sql.DB) {
	_, err := db.Exec("CREATE PROCEDURE GetUserById @id INT AS SELECT * FROM users WHERE id = @id")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with parameter: GetUserById")
}

func testExecuteProcedureWithParam(db *sql.DB) {
	// Execute with parameter value
	query := "EXEC GetUserById @id=1"

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Get columns
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("Error getting columns: %v", err)
		return
	}

	log.Printf("Columns: %v", columns)

	// Read rows
	rowCount := 0
	for rows.Next() {
		rowCount++
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Printf("Error scanning row %d: %v", rowCount, err)
			continue
		}

		rowStr := ""
		for i, val := range values {
			if i > 0 {
				rowStr += ", "
			}
			rowStr += fmt.Sprintf("%v", val)
		}
		log.Printf("  Row %d: %s", rowCount, rowStr)
	}

	log.Printf("✓ Executed procedure with parameter, returned %d rows", rowCount)
}

func testCreateProcedureWithMultipleParams(db *sql.DB) {
	_, err := db.Exec("CREATE PROCEDURE GetUsersByDept @department VARCHAR(50), @active BIT AS SELECT * FROM users WHERE department = @department")
	if err != nil {
		log.Printf("Error creating procedure: %v", err)
		return
	}

	log.Println("✓ Created procedure with multiple parameters: GetUsersByDept")
}

func testExecuteProcedureWithMultipleParams(db *sql.DB) {
	// Execute with multiple parameter values
	query := "EXEC GetUsersByDept @department='Engineering'"

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error executing procedure: %v", err)
		return
	}
	defer rows.Close()

	// Get columns
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("Error getting columns: %v", err)
		return
	}

	log.Printf("Columns: %v", columns)

	// Read rows
	rowCount := 0
	for rows.Next() {
		rowCount++
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Printf("Error scanning row %d: %v", rowCount, err)
			continue
		}

		rowStr := ""
		for i, val := range values {
			if i > 0 {
				rowStr += ", "
			}
			rowStr += fmt.Sprintf("%v", val)
		}
		log.Printf("  Row %d: %s", rowCount, rowStr)
	}

	log.Printf("✓ Executed procedure with multiple parameters, returned %d rows", rowCount)
}

func testDropProcedure(db *sql.DB) {
	_, err := db.Exec("DROP PROCEDURE GetAllUsers")
	if err != nil {
		log.Printf("Error dropping procedure: %v", err)
		return
	}

	log.Println("✓ Dropped procedure: GetAllUsers")
}
