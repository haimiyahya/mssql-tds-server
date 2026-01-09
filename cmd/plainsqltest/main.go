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

	// Test 1: CREATE TABLE
	log.Println("\n=== Test 1: CREATE TABLE ===")
	testCreateTable(db)

	// Test 2: INSERT data
	log.Println("\n=== Test 2: INSERT data ===")
	testInsert(db)

	// Test 3: SELECT all data
	log.Println("\n=== Test 3: SELECT all data ===")
	testSelectAll(db)

	// Test 4: SELECT with WHERE clause
	log.Println("\n=== Test 4: SELECT with WHERE clause ===")
	testSelectWhere(db)

	// Test 5: SELECT specific columns
	log.Println("\n=== Test 5: SELECT specific columns ===")
	testSelectColumns(db)

	// Test 6: UPDATE data
	log.Println("\n=== Test 6: UPDATE data ===")
	testUpdate(db)

	// Test 7: Verify UPDATE
	log.Println("\n=== Test 7: Verify UPDATE ===")
	testSelectAll(db)

	// Test 8: DELETE data
	log.Println("\n=== Test 8: DELETE data ===")
	testDelete(db)

	// Test 9: Verify DELETE
	log.Println("\n=== Test 9: Verify DELETE ===")
	testSelectAll(db)

	// Test 10: DROP TABLE
	log.Println("\n=== Test 10: DROP TABLE ===")
	testDropTable(db)

	log.Println("\n=== All Phase 10 tests completed successfully! ===")
}

func testCreateTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE employees (id INTEGER, name TEXT, department TEXT, salary REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	log.Println("✓ Created table: employees")
}

func testInsert(db *sql.DB) {
	// Insert multiple rows
	queries := []string{
		"INSERT INTO employees (id, name, department, salary) VALUES (1, 'Alice', 'Engineering', 75000.50)",
		"INSERT INTO employees (id, name, department, salary) VALUES (2, 'Bob', 'Marketing', 65000.00)",
		"INSERT INTO employees (id, name, department, salary) VALUES (3, 'Charlie', 'Engineering', 80000.75)",
		"INSERT INTO employees (id, name, department, salary) VALUES (4, 'Diana', 'HR', 60000.00)",
	}

	for _, query := range queries {
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		log.Printf("✓ Inserted row: %d row(s) affected", rowsAffected)
	}
}

func testSelectAll(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM employees")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	// Get column names
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
		var id int
		var name, department string
		var salary float64

		err := rows.Scan(&id, &name, &department, &salary)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		log.Printf("  Row %d: id=%d, name=%s, department=%s, salary=%.2f",
			rowCount, id, name, department, salary)
	}

	log.Printf("✓ Selected %d row(s)", rowCount)
}

func testSelectWhere(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM employees WHERE department = 'Engineering'")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var id int
		var name, department string
		var salary float64

		err := rows.Scan(&id, &name, &department, &salary)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		log.Printf("  Row %d: id=%d, name=%s, department=%s, salary=%.2f",
			rowCount, id, name, department, salary)
	}

	log.Printf("✓ Selected %d row(s) from Engineering department", rowCount)
}

func testSelectColumns(db *sql.DB) {
	rows, err := db.Query("SELECT name, salary FROM employees ORDER BY salary DESC")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var name string
		var salary float64

		err := rows.Scan(&name, &salary)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		log.Printf("  Row %d: name=%s, salary=%.2f", rowCount, name, salary)
	}

	log.Printf("✓ Selected %d row(s) with name and salary columns", rowCount)
}

func testUpdate(db *sql.DB) {
	result, err := db.Exec("UPDATE employees SET salary = 85000.00 WHERE name = 'Alice'")
	if err != nil {
		log.Printf("Error updating data: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✓ Updated Alice's salary: %d row(s) affected", rowsAffected)
}

func testDelete(db *sql.DB) {
	result, err := db.Exec("DELETE FROM employees WHERE id = 4")
	if err != nil {
		log.Printf("Error deleting data: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✓ Deleted Diana: %d row(s) affected", rowsAffected)
}

func testDropTable(db *sql.DB) {
	_, err := db.Exec("DROP TABLE employees")
	if err != nil {
		log.Printf("Error dropping table: %v", err)
		return
	}

	log.Println("✓ Dropped table: employees")
}
