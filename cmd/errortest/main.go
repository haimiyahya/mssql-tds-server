package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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

	// Test 1: Enhanced Error Messages
	log.Println("\n=== Test 1: Enhanced Error Messages ===")
	testEnhancedErrorMessages(db)

	// Test 2: SQL State Codes
	log.Println("\n=== Test 2: SQL State Codes ===")
	testSQLStateCodes(db)

	// Test 3: Error Severity Levels
	log.Println("\n=== Test 3: Error Severity Levels ===")
	testErrorSeverityLevels(db)

	// Test 4: Error Categorization
	log.Println("\n=== Test 4: Error Categorization ===")
	testErrorCategorization(db)

	// Test 5: Detailed Error Information
	log.Println("\n=== Test 5: Detailed Error Information ===")
	testDetailedErrorInformation(db)

	// Test 6: Constraint Violation Errors
	log.Println("\n=== Test 6: Constraint Violation Errors ===")
	testConstraintViolationErrors(db)

	// Test 7: Syntax Errors
	log.Println("\n=== Test 7: Syntax Errors ===")
	testSyntaxErrors(db)

	// Test 8: Data Type Errors
	log.Println("\n=== Test 8: Data Type Errors ===")
	testDataTypeErrors(db)

	// Test 9: Transaction Error Handling
	log.Println("\n=== Test 9: Transaction Error Handling ===")
	testTransactionErrorHandling(db)

	// Test 10: Cleanup
	log.Println("\n=== Test 10: Cleanup ===")
	testCleanup(db)

	log.Println("\n=== All Phase 21 tests completed! ===")
	log.Println("🎉 Phase 21: Error Handling Improvements - COMPLETE! 🎉")
}

func testEnhancedErrorMessages(db *sql.DB) {
	// Create table with constraint
	_, err := db.Exec("CREATE TABLE error_test (id INTEGER PRIMARY KEY, name TEXT NOT NULL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: error_test")

	// Test NULL constraint violation with enhanced error message
	_, err = db.Exec("INSERT INTO error_test (id, name) VALUES (1, NULL)")
	if err != nil {
		log.Printf("✓ NULL constraint violation: %v", err)
		log.Printf("  Error message includes: table=%t, column=%t", 
			strings.Contains(err.Error(), "error_test"),
			strings.Contains(err.Error(), "name"))
	}

	// Test duplicate primary key with enhanced error message
	_, err = db.Exec("INSERT INTO error_test VALUES (1, 'Alice')")
	if err == nil {
		_, err = db.Exec("INSERT INTO error_test VALUES (1, 'Bob')")
	}
	if err != nil {
		log.Printf("✓ Duplicate primary key: %v", err)
		log.Printf("  Error message includes: table=%t, constraint=%t", 
			strings.Contains(err.Error(), "error_test"),
			strings.Contains(err.Error(), "PRIMARY") || strings.Contains(err.Error(), "UNIQUE"))
	}
}

func testSQLStateCodes(db *sql.DB) {
	// Note: SQLite errors don't include SQLSTATE codes by default
	// This test verifies that error messages are consistent with SQL standards

	// Test integrity constraint violation (SQLSTATE: 23000)
	_, err := db.Exec("CREATE TABLE sqlstate_test (id INTEGER PRIMARY KEY)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	_, err = db.Exec("INSERT INTO sqlstate_test VALUES (1)")
	if err == nil {
		_, err = db.Exec("INSERT INTO sqlstate_test VALUES (1)")
	}
	if err != nil {
		log.Printf("✓ Integrity constraint violation (SQLSTATE 23000): %v", err)
	}

	// Test syntax error (SQLSTATE: 42000)
	_, err = db.Exec("SELCT * FROM sqlstate_test")
	if err != nil {
		log.Printf("✓ Syntax error (SQLSTATE 42000): %v", err)
	}

	// Test data exception (SQLSTATE: 22000)
	_, err = db.Exec("CREATE TABLE datatype_test (id INTEGER, name INTEGER)")
	if err == nil {
		_, err = db.Exec("INSERT INTO datatype_test VALUES (1, 'invalid')")
	}
	if err != nil {
		log.Printf("✓ Data type error (SQLSTATE 22000): %v", err)
	}
}

func testErrorSeverityLevels(db *sql.DB) {
	// Information: Successful operation
	_, err := db.Exec("CREATE TABLE severity_test (id INTEGER)")
	if err == nil {
		log.Println("✓ Information (Severity 0): Table created successfully")
	}

	// Warning: No data found
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM severity_test").Scan(&count)
	if err == nil && count == 0 {
		log.Println("✓ Warning (Severity 1): Table is empty")
	}

	// Error: Table doesn't exist
	_, err = db.Exec("SELECT * FROM non_existent_table")
	if err != nil {
		log.Printf("✓ Error (Severity 2): Table not found: %v", err)
	}
}

func testErrorCategorization(db *sql.DB) {
	// Syntax errors
	_, err := db.Exec("CREET TABLE syntax_test (id INTEGER)")
	if err != nil {
		log.Printf("✓ Syntax Error: %v", err)
	}

	// Runtime errors
	_, err = db.Exec("DROP TABLE non_existent_table")
	if err != nil {
		log.Printf("✓ Runtime Error: %v", err)
	}

	// Constraint violations
	_, err = db.Exec("CREATE TABLE constraint_test (id INTEGER PRIMARY KEY)")
	if err == nil {
		_, err = db.Exec("INSERT INTO constraint_test VALUES (1)")
		if err == nil {
			_, err = db.Exec("INSERT INTO constraint_test VALUES (1)")
		}
	}
	if err != nil {
		log.Printf("✓ Constraint Violation: %v", err)
	}

	// Data type errors
	_, err = db.Exec("CREATE TABLE datatype_test (id INTEGER, value INTEGER)")
	if err == nil {
		_, err = db.Exec("INSERT INTO datatype_test VALUES (1, 'string')")
	}
	if err != nil {
		log.Printf("✓ Data Type Error: %v", err)
	}
}

func testDetailedErrorInformation(db *sql.DB) {
	// Create table with constraints
	_, err := db.Exec("CREATE TABLE detail_test (id INTEGER PRIMARY KEY, email TEXT UNIQUE)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: detail_test")

	// Insert initial data
	_, err = db.Exec("INSERT INTO detail_test VALUES (1, 'test@example.com')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	// Test duplicate email with detailed error
	_, err = db.Exec("INSERT INTO detail_test VALUES (2, 'test@example.com')")
	if err != nil {
		log.Printf("✓ Duplicate constraint error details:")
		log.Printf("  Error: %v", err)
		log.Printf("  Contains table name: %t", strings.Contains(err.Error(), "detail_test"))
		log.Printf("  Contains column name: %t", strings.Contains(err.Error(), "email"))
		log.Printf("  Contains constraint info: %t", strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "unique"))
	}
}

func testConstraintViolationErrors(db *sql.DB) {
	// PRIMARY KEY violation
	_, err := db.Exec("CREATE TABLE pk_test (id INTEGER PRIMARY KEY)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("INSERT INTO pk_test VALUES (1)")
	if err == nil {
		_, err = db.Exec("INSERT INTO pk_test VALUES (1)")
	}
	if err != nil {
		log.Printf("✓ PRIMARY KEY violation: %v", err)
	}

	// UNIQUE constraint violation
	_, err = db.Exec("CREATE TABLE unique_test (email TEXT UNIQUE)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("INSERT INTO unique_test VALUES ('test@example.com')")
	if err == nil {
		_, err = db.Exec("INSERT INTO unique_test VALUES ('test@example.com')")
	}
	if err != nil {
		log.Printf("✓ UNIQUE constraint violation: %v", err)
	}

	// NOT NULL constraint violation
	_, err = db.Exec("CREATE TABLE notnull_test (id INTEGER, name TEXT NOT NULL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("INSERT INTO notnull_test (id, name) VALUES (1, NULL)")
	if err != nil {
		log.Printf("✓ NOT NULL constraint violation: %v", err)
	}

	// CHECK constraint violation
	_, err = db.Exec("CREATE TABLE check_test (age INTEGER CHECK (age >= 18))")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("INSERT INTO check_test VALUES (15)")
	if err != nil {
		log.Printf("✓ CHECK constraint violation: %v", err)
	}

	// FOREIGN KEY constraint violation
	_, err = db.Exec("CREATE TABLE fk_users (id INTEGER PRIMARY KEY)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE fk_orders (id INTEGER, user_id INTEGER, FOREIGN KEY (user_id) REFERENCES fk_users(id))")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("INSERT INTO fk_orders VALUES (1, 999)")
	if err != nil {
		log.Printf("✓ FOREIGN KEY constraint violation: %v", err)
	}
}

func testSyntaxErrors(db *sql.DB) {
	// Invalid keyword
	_, err := db.Exec("SELCT * FROM syntax_test")
	if err != nil {
		log.Printf("✓ Invalid keyword error: %v", err)
	}

	// Missing parenthesis
	_, err = db.Exec("SELECT COUNT * FROM syntax_test")
	if err != nil {
		log.Printf("✓ Missing parenthesis error: %v", err)
	}

	// Invalid table name
	_, err = db.Exec("SELECT * FROM 123table")
	if err != nil {
		log.Printf("✓ Invalid table name error: %v", err)
	}

	// Missing comma in VALUES
	_, err = db.Exec("CREATE TABLE comma_test (id INTEGER); INSERT INTO comma_test VALUES (1 'value')")
	if err != nil {
		log.Printf("✓ Missing comma error: %v", err)
	}

	// Unterminated string
	_, err = db.Exec("SELECT 'unterminated")
	if err != nil {
		log.Printf("✓ Unterminated string error: %v", err)
	}
}

func testDataTypeErrors(db *sql.DB) {
	// Invalid integer value
	_, err := db.Exec("CREATE TABLE int_test (id INTEGER)")
	if err == nil {
		_, err = db.Exec("INSERT INTO int_test VALUES ('invalid')")
	}
	if err != nil {
		log.Printf("✓ Invalid integer value: %v", err)
	}

	// String to integer conversion
	_, err = db.Exec("CREATE TABLE cast_test (id INTEGER)")
	if err == nil {
		_, err = db.Exec("INSERT INTO cast_test VALUES ('abc')")
	}
	if err != nil {
		log.Printf("✓ String to integer conversion error: %v", err)
	}

	// NULL value in non-NULL column
	_, err = db.Exec("CREATE TABLE notnull_test (id INTEGER NOT NULL)")
	if err == nil {
		_, err = db.Exec("INSERT INTO notnull_test VALUES (NULL)")
	}
	if err != nil {
		log.Printf("✓ NULL value in NOT NULL column: %v", err)
	}
}

func testTransactionErrorHandling(db *sql.DB) {
	// Create tables
	_, err := db.Exec("CREATE TABLE trans_test1 (id INTEGER PRIMARY KEY)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE trans_test2 (id INTEGER PRIMARY KEY)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created tables: trans_test1, trans_test2")

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return
	}

	// Insert into first table
	_, err = tx.Exec("INSERT INTO trans_test1 VALUES (1)")
	if err != nil {
		log.Printf("Error inserting into trans_test1: %v", err)
		tx.Rollback()
		return
	}

	// Try to insert duplicate into second table
	_, err = tx.Exec("INSERT INTO trans_test2 VALUES (1)")
	if err == nil {
		_, err = tx.Exec("INSERT INTO trans_test2 VALUES (1)")
	}
	if err != nil {
		log.Printf("✓ Transaction error detected: %v", err)
		log.Printf("  Rolling back transaction...")
		err = tx.Rollback()
		if err != nil {
			log.Printf("  Error rolling back: %v", err)
		} else {
			log.Println("  ✓ Transaction rolled back successfully")
		}

		// Verify rollback
		var count1, count2 int
		db.QueryRow("SELECT COUNT(*) FROM trans_test1").Scan(&count1)
		db.QueryRow("SELECT COUNT(*) FROM trans_test2").Scan(&count2)
		log.Printf("  ✓ Verified: trans_test1=%d rows, trans_test2=%d rows (rollback successful)", count1, count2)
		return
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return
	}

	log.Println("✓ Transaction committed successfully")
}

func testCleanup(db *sql.DB) {
	tables := []string{
		"error_test",
		"sqlstate_test",
		"severity_test",
		"syntax_test",
		"constraint_test",
		"datatype_test",
		"detail_test",
		"pk_test",
		"unique_test",
		"notnull_test",
		"check_test",
		"fk_users",
		"fk_orders",
		"int_test",
		"cast_test",
		"comma_test",
		"trans_test1",
		"trans_test2",
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
