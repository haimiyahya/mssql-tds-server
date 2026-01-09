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

	// Test 3: COUNT(*)
	log.Println("\n=== Test 3: COUNT(*) ===")
	testCountAll(db)

	// Test 4: COUNT(column)
	log.Println("\n=== Test 4: COUNT(column) ===")
	testCountColumn(db)

	// Test 5: SUM(column)
	log.Println("\n=== Test 5: SUM(column) ===")
	testSum(db)

	// Test 6: AVG(column)
	log.Println("\n=== Test 6: AVG(column) ===")
	testAvg(db)

	// Test 7: MIN(column)
	log.Println("\n=== Test 7: MIN(column) ===")
	testMin(db)

	// Test 8: MAX(column)
	log.Println("\n=== Test 8: MAX(column) ===")
	testMax(db)

	// Test 9: COUNT(DISTINCT column)
	log.Println("\n=== Test 9: COUNT(DISTINCT column) ===")
	testCountDistinct(db)

	// Test 10: Multiple aggregates
	log.Println("\n=== Test 10: Multiple aggregates ===")
	testMultipleAggregates(db)

	// Test 11: Aggregates with WHERE
	log.Println("\n=== Test 11: Aggregates with WHERE ===")
	testAggregatesWithWhere(db)

	// Test 12: DROP TABLE
	log.Println("\n=== Test 12: DROP TABLE ===")
	testDropTable(db)

	log.Println("\n=== All Phase 11 Iteration 2 tests completed! ===")
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
		"INSERT INTO employees (id, name, department, salary) VALUES (5, 'Eve', 'Engineering', 85000.00)",
		"INSERT INTO employees (id, name, department, salary) VALUES (6, 'Frank', 'Marketing', 70000.00)",
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

func testCountAll(db *sql.DB) {
	rows, err := db.Query("SELECT COUNT(*) FROM employees")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var count int64
		err := rows.Scan(&count)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: COUNT(*) = %d", count)
	}

	log.Printf("✓ COUNT(*) test completed: %d row(s)", rowCount)
}

func testCountColumn(db *sql.DB) {
	rows, err := db.Query("SELECT COUNT(department) FROM employees")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var count int64
		err := rows.Scan(&count)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: COUNT(department) = %d", count)
	}

	log.Printf("✓ COUNT(column) test completed: %d row(s)", rowCount)
}

func testSum(db *sql.DB) {
	rows, err := db.Query("SELECT SUM(salary) FROM employees")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var sum float64
		err := rows.Scan(&sum)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: SUM(salary) = %.2f", sum)
	}

	log.Printf("✓ SUM test completed: %d row(s)", rowCount)
}

func testAvg(db *sql.DB) {
	rows, err := db.Query("SELECT AVG(salary) FROM employees")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var avg float64
		err := rows.Scan(&avg)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: AVG(salary) = %.2f", avg)
	}

	log.Printf("✓ AVG test completed: %d row(s)", rowCount)
}

func testMin(db *sql.DB) {
	rows, err := db.Query("SELECT MIN(salary) FROM employees")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var min float64
		err := rows.Scan(&min)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: MIN(salary) = %.2f", min)
	}

	log.Printf("✓ MIN test completed: %d row(s)", rowCount)
}

func testMax(db *sql.DB) {
	rows, err := db.Query("SELECT MAX(salary) FROM employees")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var max float64
		err := rows.Scan(&max)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: MAX(salary) = %.2f", max)
	}

	log.Printf("✓ MAX test completed: %d row(s)", rowCount)
}

func testCountDistinct(db *sql.DB) {
	rows, err := db.Query("SELECT COUNT(DISTINCT department) FROM employees")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var count int64
		err := rows.Scan(&count)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: COUNT(DISTINCT department) = %d", count)
	}

	log.Printf("✓ COUNT(DISTINCT) test completed: %d row(s)", rowCount)
}

func testMultipleAggregates(db *sql.DB) {
	rows, err := db.Query("SELECT COUNT(*), SUM(salary), AVG(salary), MIN(salary), MAX(salary) FROM employees")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var count int64
		var sum, avg, min, max float64
		err := rows.Scan(&count, &sum, &avg, &min, &max)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result:")
		log.Printf("    COUNT(*) = %d", count)
		log.Printf("    SUM(salary) = %.2f", sum)
		log.Printf("    AVG(salary) = %.2f", avg)
		log.Printf("    MIN(salary) = %.2f", min)
		log.Printf("    MAX(salary) = %.2f", max)
	}

	log.Printf("✓ Multiple aggregates test completed: %d row(s)", rowCount)
}

func testAggregatesWithWhere(db *sql.DB) {
	rows, err := db.Query("SELECT COUNT(*), AVG(salary) FROM employees WHERE department = 'Engineering'")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var count int64
		var avg float64
		err := rows.Scan(&count, &avg)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result:")
		log.Printf("    COUNT(*) for Engineering = %d", count)
		log.Printf("    AVG(salary) for Engineering = %.2f", avg)
	}

	log.Printf("✓ Aggregates with WHERE test completed: %d row(s)", rowCount)
}

func testDropTable(db *sql.DB) {
	_, err := db.Exec("DROP TABLE employees")
	if err != nil {
		log.Printf("Error dropping table: %v", err)
		return
	}

	log.Println("✓ Dropped table: employees")
}
