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

	// Test 1: CREATE TABLES
	log.Println("\n=== Test 1: CREATE TABLES ===")
	testCreateTables(db)

	// Test 2: INSERT data
	log.Println("\n=== Test 2: INSERT data ===")
	testInsert(db)

	// Test 3: INNER JOIN
	log.Println("\n=== Test 3: INNER JOIN ===")
	testInnerJoin(db)

	// Test 4: LEFT JOIN
	log.Println("\n=== Test 4: LEFT JOIN ===")
	testLeftJoin(db)

	// Test 5: RIGHT JOIN (will fail with SQLite)
	log.Println("\n=== Test 5: RIGHT JOIN (expected to fail) ===")
	testRightJoin(db)

	// Test 6: FULL JOIN (will fail with SQLite)
	log.Println("\n=== Test 6: FULL JOIN (expected to fail) ===")
	testFullJoin(db)

	// Test 7: Multiple JOINs
	log.Println("\n=== Test 7: Multiple JOINs ===")
	testMultipleJoins(db)

	// Test 8: JOIN with WHERE
	log.Println("\n=== Test 8: JOIN with WHERE ===")
	testJoinWithWhere(db)

	// Test 9: JOIN with GROUP BY
	log.Println("\n=== Test 9: JOIN with GROUP BY ===")
	testJoinWithGroupBy(db)

	// Test 10: Self JOIN
	log.Println("\n=== Test 10: Self JOIN ===")
	testSelfJoin(db)

	// Test 11: JOIN with table alias
	log.Println("\n=== Test 11: JOIN with table alias ===")
	testJoinWithAlias(db)

	// Test 12: DROP TABLES
	log.Println("\n=== Test 12: DROP TABLES ===")
	testDropTables(db)

	log.Println("\n=== All Phase 11 Iteration 4 tests completed! ===")
}

func testCreateTables(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE departments (id INTEGER, name TEXT)")
	if err != nil {
		log.Printf("Error creating departments table: %v", err)
		return
	}
	log.Println("✓ Created table: departments")

	_, err = db.Exec("CREATE TABLE employees (id INTEGER, name TEXT, department_id INTEGER, salary REAL)")
	if err != nil {
		log.Printf("Error creating employees table: %v", err)
		return
	}
	log.Println("✓ Created table: employees")

	_, err = db.Exec("CREATE TABLE projects (id INTEGER, name TEXT, employee_id INTEGER)")
	if err != nil {
		log.Printf("Error creating projects table: %v", err)
		return
	}
	log.Println("✓ Created table: projects")
}

func testInsert(db *sql.DB) {
	// Insert departments
	deptQueries := []string{
		"INSERT INTO departments VALUES (1, 'Engineering')",
		"INSERT INTO departments VALUES (2, 'Marketing')",
		"INSERT INTO departments VALUES (3, 'HR')",
	}

	for _, query := range deptQueries {
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("Error inserting department: %v", err)
			continue
		}
		rowsAffected, _ := result.RowsAffected()
		log.Printf("✓ Inserted department: %d row(s)", rowsAffected)
	}

	// Insert employees
	empQueries := []string{
		"INSERT INTO employees VALUES (1, 'Alice', 1, 75000.00)",
		"INSERT INTO employees VALUES (2, 'Bob', 1, 80000.00)",
		"INSERT INTO employees VALUES (3, 'Charlie', 2, 65000.00)",
		"INSERT INTO employees VALUES (4, 'Diana', 2, 70000.00)",
		"INSERT INTO employees VALUES (5, 'Eve', 3, 60000.00)",
	}

	for _, query := range empQueries {
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("Error inserting employee: %v", err)
			continue
		}
		rowsAffected, _ := result.RowsAffected()
		log.Printf("✓ Inserted employee: %d row(s)", rowsAffected)
	}

	// Insert projects
	projectQueries := []string{
		"INSERT INTO projects VALUES (1, 'Project A', 1)",
		"INSERT INTO projects VALUES (2, 'Project B', 2)",
		"INSERT INTO projects VALUES (3, 'Project C', 1)",
	}

	for _, query := range projectQueries {
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("Error inserting project: %v", err)
			continue
		}
		rowsAffected, _ := result.RowsAffected()
		log.Printf("✓ Inserted project: %d row(s)", rowsAffected)
	}
}

func testInnerJoin(db *sql.DB) {
	// INNER JOIN: Only rows with matching values in both tables
	rows, err := db.Query("SELECT e.name, d.name FROM employees e INNER JOIN departments d ON e.department_id = d.id")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var empName, deptName string
		err := rows.Scan(&empName, &deptName)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s works in %s", empName, deptName)
	}

	log.Printf("✓ INNER JOIN test completed: %d row(s)", rowCount)
}

func testLeftJoin(db *sql.DB) {
	// LEFT JOIN: All rows from left table, matching rows from right table
	rows, err := db.Query("SELECT e.name, d.name FROM employees e LEFT JOIN departments d ON e.department_id = d.id")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var empName, deptName string
		err := rows.Scan(&empName, &deptName)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		if deptName == "" {
			log.Printf("  Result: %s works in (no department)", empName)
		} else {
			log.Printf("  Result: %s works in %s", empName, deptName)
		}
	}

	log.Printf("✓ LEFT JOIN test completed: %d row(s)", rowCount)
}

func testRightJoin(db *sql.DB) {
	// RIGHT JOIN: All rows from right table, matching rows from left table
	// Note: SQLite doesn't support RIGHT JOIN directly
	// This test is expected to fail
	rows, err := db.Query("SELECT e.name, d.name FROM employees e RIGHT JOIN departments d ON e.department_id = d.id")
	if err != nil {
		log.Printf("  ✗ RIGHT JOIN failed (expected): %v", err)
		log.Println("  Note: SQLite doesn't support RIGHT JOIN natively")
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var empName, deptName string
		err := rows.Scan(&empName, &deptName)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s works in %s", empName, deptName)
	}

	log.Printf("✓ RIGHT JOIN test completed: %d row(s)", rowCount)
}

func testFullJoin(db *sql.DB) {
	// FULL JOIN: All rows from both tables, matching where possible
	// Note: SQLite doesn't support FULL JOIN directly
	// This test is expected to fail
	rows, err := db.Query("SELECT e.name, d.name FROM employees e FULL JOIN departments d ON e.department_id = d.id")
	if err != nil {
		log.Printf("  ✗ FULL JOIN failed (expected): %v", err)
		log.Println("  Note: SQLite doesn't support FULL JOIN natively")
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var empName, deptName string
		err := rows.Scan(&empName, &deptName)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s works in %s", empName, deptName)
	}

	log.Printf("✓ FULL JOIN test completed: %d row(s)", rowCount)
}

func testMultipleJoins(db *sql.DB) {
	// Multiple JOINs: Join more than 2 tables
	rows, err := db.Query("SELECT e.name, d.name, p.name FROM employees e INNER JOIN departments d ON e.department_id = d.id INNER JOIN projects p ON e.id = p.employee_id")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var empName, deptName, projectName string
		err := rows.Scan(&empName, &deptName, &projectName)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s works in %s on %s", empName, deptName, projectName)
	}

	log.Printf("✓ Multiple JOINs test completed: %d row(s)", rowCount)
}

func testJoinWithWhere(db *sql.DB) {
	// JOIN with WHERE: Filter joined results
	rows, err := db.Query("SELECT e.name, d.name FROM employees e INNER JOIN departments d ON e.department_id = d.id WHERE e.salary > 70000")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var empName, deptName string
		err := rows.Scan(&empName, &deptName)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s works in %s (salary > 70000)", empName, deptName)
	}

	log.Printf("✓ JOIN with WHERE test completed: %d row(s)", rowCount)
}

func testJoinWithGroupBy(db *sql.DB) {
	// JOIN with GROUP BY: Group joined results
	rows, err := db.Query("SELECT d.name, COUNT(*), AVG(e.salary) FROM employees e INNER JOIN departments d ON e.department_id = d.id GROUP BY d.name")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var deptName string
		var count int64
		var avgSalary float64
		err := rows.Scan(&deptName, &count, &avgSalary)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s - Count: %d, Avg Salary: %.2f", deptName, count, avgSalary)
	}

	log.Printf("✓ JOIN with GROUP BY test completed: %d row(s)", rowCount)
}

func testSelfJoin(db *sql.DB) {
	// Self JOIN: Join table with itself
	// Find employees who work in the same department
	rows, err := db.Query("SELECT e1.name, e2.name FROM employees e1 INNER JOIN employees e2 ON e1.department_id = e2.department_id WHERE e1.id < e2.id")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var emp1Name, emp2Name string
		err := rows.Scan(&emp1Name, &emp2Name)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s and %s work together", emp1Name, emp2Name)
	}

	log.Printf("✓ Self JOIN test completed: %d row(s)", rowCount)
}

func testJoinWithAlias(db *sql.DB) {
	// JOIN with table alias: Use alias for table names
	rows, err := db.Query("SELECT e.name, d.name FROM employees AS e INNER JOIN departments AS d ON e.department_id = d.id")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var empName, deptName string
		err := rows.Scan(&empName, &deptName)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s works in %s", empName, deptName)
	}

	log.Printf("✓ JOIN with table alias test completed: %d row(s)", rowCount)
}

func testDropTables(db *sql.DB) {
	_, err := db.Exec("DROP TABLE projects")
	if err != nil {
		log.Printf("Error dropping projects table: %v", err)
		return
	}
	log.Println("✓ Dropped table: projects")

	_, err = db.Exec("DROP TABLE employees")
	if err != nil {
		log.Printf("Error dropping employees table: %v", err)
		return
	}
	log.Println("✓ Dropped table: employees")

	_, err = db.Exec("DROP TABLE departments")
	if err != nil {
		log.Printf("Error dropping departments table: %v", err)
		return
	}
	log.Println("✓ Dropped table: departments")
}
