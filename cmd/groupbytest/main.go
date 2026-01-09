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

	// Test 3: GROUP BY with COUNT
	log.Println("\n=== Test 3: GROUP BY with COUNT ===")
	testGroupByCount(db)

	// Test 4: GROUP BY with SUM
	log.Println("\n=== Test 4: GROUP BY with SUM ===")
	testGroupBySum(db)

	// Test 5: GROUP BY with AVG
	log.Println("\n=== Test 5: GROUP BY with AVG ===")
	testGroupByAvg(db)

	// Test 6: GROUP BY with multiple aggregates
	log.Println("\n=== Test 6: GROUP BY with multiple aggregates ===")
	testGroupByMultiple(db)

	// Test 7: GROUP BY with WHERE
	log.Println("\n=== Test 7: GROUP BY with WHERE ===")
	testGroupByWhere(db)

	// Test 8: GROUP BY with HAVING
	log.Println("\n=== Test 8: GROUP BY with HAVING ===")
	testGroupByHaving(db)

	// Test 9: GROUP BY with ORDER BY
	log.Println("\n=== Test 9: GROUP BY with ORDER BY ===")
	testGroupByOrderBy(db)

	// Test 10: GROUP BY + HAVING + ORDER BY (combined)
	log.Println("\n=== Test 10: GROUP BY + HAVING + ORDER BY (combined) ===")
	testGroupByCombined(db)

	// Test 11: GROUP BY with multiple columns
	log.Println("\n=== Test 11: GROUP BY with multiple columns ===")
	testGroupByMultipleColumns(db)

	// Test 12: DROP TABLE
	log.Println("\n=== Test 12: DROP TABLE ===")
	testDropTable(db)

	log.Println("\n=== All Phase 11 Iteration 3 tests completed! ===")
}

func testCreateTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE sales (id INTEGER, product TEXT, category TEXT, quantity INTEGER, price REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	log.Println("✓ Created table: sales")
}

func testInsert(db *sql.DB) {
	// Insert multiple rows
	queries := []string{
		"INSERT INTO sales (id, product, category, quantity, price) VALUES (1, 'Laptop', 'Electronics', 5, 1200.00)",
		"INSERT INTO sales (id, product, category, quantity, price) VALUES (2, 'Mouse', 'Electronics', 20, 25.00)",
		"INSERT INTO sales (id, product, category, quantity, price) VALUES (3, 'Chair', 'Furniture', 10, 150.00)",
		"INSERT INTO sales (id, product, category, quantity, price) VALUES (4, 'Desk', 'Furniture', 5, 300.00)",
		"INSERT INTO sales (id, product, category, quantity, price) VALUES (5, 'Keyboard', 'Electronics', 15, 75.00)",
		"INSERT INTO sales (id, product, category, quantity, price) VALUES (6, 'Monitor', 'Electronics', 8, 250.00)",
		"INSERT INTO sales (id, product, category, quantity, price) VALUES (7, 'Table', 'Furniture', 3, 400.00)",
		"INSERT INTO sales (id, product, category, quantity, price) VALUES (8, 'Headphones', 'Electronics', 12, 100.00)",
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

func testGroupByCount(db *sql.DB) {
	rows, err := db.Query("SELECT category, COUNT(*) FROM sales GROUP BY category")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var category string
		var count int64
		err := rows.Scan(&category, &count)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s - Count = %d", category, count)
	}

	log.Printf("✓ GROUP BY COUNT test completed: %d row(s)", rowCount)
}

func testGroupBySum(db *sql.DB) {
	rows, err := db.Query("SELECT category, SUM(quantity) FROM sales GROUP BY category")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var category string
		var sum int64
		err := rows.Scan(&category, &sum)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s - SUM(quantity) = %d", category, sum)
	}

	log.Printf("✓ GROUP BY SUM test completed: %d row(s)", rowCount)
}

func testGroupByAvg(db *sql.DB) {
	rows, err := db.Query("SELECT category, AVG(price) FROM sales GROUP BY category")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var category string
		var avg float64
		err := rows.Scan(&category, &avg)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s - AVG(price) = %.2f", category, avg)
	}

	log.Printf("✓ GROUP BY AVG test completed: %d row(s)", rowCount)
}

func testGroupByMultiple(db *sql.DB) {
	rows, err := db.Query("SELECT category, COUNT(*), SUM(quantity), AVG(price), MIN(price), MAX(price) FROM sales GROUP BY category")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var category string
		var count int64
		var sum int64
		var avg, min, max float64
		err := rows.Scan(&category, &count, &sum, &avg, &min, &max)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s", category)
		log.Printf("    COUNT(*) = %d", count)
		log.Printf("    SUM(quantity) = %d", sum)
		log.Printf("    AVG(price) = %.2f", avg)
		log.Printf("    MIN(price) = %.2f", min)
		log.Printf("    MAX(price) = %.2f", max)
	}

	log.Printf("✓ GROUP BY multiple aggregates test completed: %d row(s)", rowCount)
}

func testGroupByWhere(db *sql.DB) {
	rows, err := db.Query("SELECT category, COUNT(*), SUM(quantity) FROM sales WHERE quantity > 10 GROUP BY category")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var category string
		var count int64
		var sum int64
		err := rows.Scan(&category, &count, &sum)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s - COUNT(*) = %d, SUM(quantity) = %d", category, count, sum)
	}

	log.Printf("✓ GROUP BY with WHERE test completed: %d row(s)", rowCount)
}

func testGroupByHaving(db *sql.DB) {
	rows, err := db.Query("SELECT category, COUNT(*) FROM sales GROUP BY category HAVING COUNT(*) > 3")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var category string
		var count int64
		err := rows.Scan(&category, &count)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s - COUNT(*) = %d", category, count)
	}

	log.Printf("✓ GROUP BY with HAVING test completed: %d row(s)", rowCount)
}

func testGroupByOrderBy(db *sql.DB) {
	rows, err := db.Query("SELECT category, COUNT(*) FROM sales GROUP BY category ORDER BY category DESC")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var category string
		var count int64
		err := rows.Scan(&category, &count)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s - COUNT(*) = %d", category, count)
	}

	log.Printf("✓ GROUP BY with ORDER BY test completed: %d row(s)", rowCount)
}

func testGroupByCombined(db *sql.DB) {
	rows, err := db.Query("SELECT category, COUNT(*), AVG(price) FROM sales GROUP BY category HAVING COUNT(*) > 3 ORDER BY category DESC")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var category string
		var count int64
		var avg float64
		err := rows.Scan(&category, &count, &avg)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s - COUNT(*) = %d, AVG(price) = %.2f", category, count, avg)
	}

	log.Printf("✓ GROUP BY + HAVING + ORDER BY combined test completed: %d row(s)", rowCount)
}

func testGroupByMultipleColumns(db *sql.DB) {
	// Create a more complex table for multiple column GROUP BY
	_, err := db.Exec("CREATE TABLE employees (id INTEGER, name TEXT, department TEXT, title TEXT, salary REAL)")
	if err != nil {
		log.Printf("Error creating employees table: %v", err)
		return
	}

	// Insert data
	_, err = db.Exec("INSERT INTO employees VALUES (1, 'Alice', 'Engineering', 'Developer', 75000.00)")
	if err != nil {
		log.Printf("Error inserting Alice: %v", err)
	}
	_, err = db.Exec("INSERT INTO employees VALUES (2, 'Bob', 'Engineering', 'Manager', 90000.00)")
	if err != nil {
		log.Printf("Error inserting Bob: %v", err)
	}
	_, err = db.Exec("INSERT INTO employees VALUES (3, 'Charlie', 'Engineering', 'Developer', 80000.00)")
	if err != nil {
		log.Printf("Error inserting Charlie: %v", err)
	}
	_, err = db.Exec("INSERT INTO employees VALUES (4, 'Diana', 'Marketing', 'Analyst', 65000.00)")
	if err != nil {
		log.Printf("Error inserting Diana: %v", err)
	}
	_, err = db.Exec("INSERT INTO employees VALUES (5, 'Eve', 'Marketing', 'Manager', 85000.00)")
	if err != nil {
		log.Printf("Error inserting Eve: %v", err)
	}

	// Test GROUP BY with multiple columns
	rows, err := db.Query("SELECT department, title, COUNT(*), AVG(salary) FROM employees GROUP BY department, title")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var department, title string
		var count int64
		var avg float64
		err := rows.Scan(&department, &title, &count, &avg)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s, %s - COUNT(*) = %d, AVG(salary) = %.2f", department, title, count, avg)
	}

	log.Printf("✓ GROUP BY with multiple columns test completed: %d row(s)", rowCount)

	// Drop the temporary table
	_, _ = db.Exec("DROP TABLE employees")
}

func testDropTable(db *sql.DB) {
	_, err := db.Exec("DROP TABLE sales")
	if err != nil {
		log.Printf("Error dropping table: %v", err)
		return
	}

	log.Println("✓ Dropped table: sales")
}
