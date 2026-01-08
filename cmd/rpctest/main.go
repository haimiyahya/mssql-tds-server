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

	// Test 1: Call SP_HELLO (no parameters)
	log.Println("\n=== Test 1: SP_HELLO (no parameters) ===")
	testStoredProcedure(db, "SP_HELLO", nil)

	// Test 2: Call SP_HELLO with parameter
	log.Println("\n=== Test 2: SP_HELLO (with parameter) ===")
	testStoredProcedure(db, "SP_HELLO", []interface{}{"Bob"})

	// Test 3: Call SP_ECHO with multiple parameters
	log.Println("\n=== Test 3: SP_ECHO (multiple parameters) ===")
	testStoredProcedure(db, "SP_ECHO", []interface{}{"hello", 123, "world"})

	// Test 4: Call SP_GET_DATA (returns multiple rows)
	log.Println("\n=== Test 4: SP_GET_DATA (returns dataset) ===")
	testStoredProcedure(db, "SP_GET_DATA", nil)

	// Test 5: Call SP_GET_DATA with filter parameter
	log.Println("\n=== Test 5: SP_GET_DATA (with filter) ===")
	testStoredProcedure(db, "SP_GET_DATA", []interface{}{"ENGINEERING"})

	log.Println("\n=== All tests completed successfully! ===")
}

func testStoredProcedure(db *sql.DB, procName string, params []interface{}) {
	log.Printf("Calling stored procedure: %s", procName)

	// Use simple SQL EXEC command to trigger RPC
	// The mssql client library handles RPC conversion internally
	query := fmt.Sprintf("EXEC %s", procName)
	var rows *sql.Rows
	var err error

	if len(params) == 0 {
		rows, err = db.Query(query)
	} else {
		rows, err = db.Query(query, params...)
	}

	if err != nil {
		log.Printf("Error calling stored procedure: %v", err)
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

	err = rows.Err()
	if err != nil {
		log.Printf("Rows error: %v", err)
	}

	log.Printf("Returned %d rows", rowCount)
}
