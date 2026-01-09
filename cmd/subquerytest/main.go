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

	// Test 3: Subquery in WHERE clause with IN
	log.Println("\n=== Test 3: Subquery in WHERE clause with IN ===")
	testSubqueryWhereIn(db)

	// Test 4: Subquery in WHERE clause with NOT IN
	log.Println("\n=== Test 4: Subquery in WHERE clause with NOT IN ===")
	testSubqueryWhereNotIn(db)

	// Test 5: Subquery in WHERE clause with EXISTS
	log.Println("\n=== Test 5: Subquery in WHERE clause with EXISTS ===")
	testSubqueryWhereExists(db)

	// Test 6: Subquery in WHERE clause with NOT EXISTS
	log.Println("\n=== Test 6: Subquery in WHERE clause with NOT EXISTS ===")
	testSubqueryWhereNotExists(db)

	// Test 7: Subquery in WHERE clause with =
	log.Println("\n=== Test 7: Subquery in WHERE clause with = ===")
	testSubqueryWhereEquals(db)

	// Test 8: Subquery in WHERE clause with >
	log.Println("\n=== Test 8: Subquery in WHERE clause with > ===")
	testSubqueryWhereGreater(db)

	// Test 9: Subquery in SELECT list
	log.Println("\n=== Test 9: Subquery in SELECT list ===")
	testSubquerySelectList(db)

	// Test 10: Subquery in FROM clause (derived table)
	log.Println("\n=== Test 10: Subquery in FROM clause (derived table) ===")
	testSubqueryFrom(db)

	// Test 11: Correlated subquery
	log.Println("\n=== Test 11: Correlated subquery ===")
	testCorrelatedSubquery(db)

	// Test 12: Nested subquery (subquery within subquery)
	log.Println("\n=== Test 12: Nested subquery ===")
	testNestedSubquery(db)

	// Test 13: Subquery with JOIN
	log.Println("\n=== Test 13: Subquery with JOIN ===")
	testSubqueryWithJoin(db)

	// Test 14: Subquery with GROUP BY
	log.Println("\n=== Test 14: Subquery with GROUP BY ===")
	testSubqueryWithGroupBy(db)

	// Test 15: DROP TABLES
	log.Println("\n=== Test 15: DROP TABLES ===")
	testDropTables(db)

	log.Println("\n=== All Phase 11 Iteration 5 tests completed! ===")
	log.Println("🎉 Phase 11: Advanced SELECT Features - COMPLETE! 🎉")
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

	_, err = db.Exec("CREATE TABLE salaries (employee_id INTEGER, amount REAL)")
	if err != nil {
		log.Printf("Error creating salaries table: %v", err)
		return
	}
	log.Println("✓ Created table: salaries")
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
		"INSERT INTO employees VALUES (6, 'Frank', 1, 85000.00)",
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

	// Insert salaries
	salaryQueries := []string{
		"INSERT INTO salaries VALUES (1, 75000.00)",
		"INSERT INTO salaries VALUES (2, 80000.00)",
		"INSERT INTO salaries VALUES (3, 65000.00)",
		"INSERT INTO salaries VALUES (4, 70000.00)",
		"INSERT INTO salaries VALUES (5, 60000.00)",
		"INSERT INTO salaries VALUES (6, 85000.00)",
	}

	for _, query := range salaryQueries {
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("Error inserting salary: %v", err)
			continue
		}
		rowsAffected, _ := result.RowsAffected()
		log.Printf("✓ Inserted salary: %d row(s)", rowsAffected)
	}
}

func testSubqueryWhereIn(db *sql.DB) {
	// Subquery in WHERE clause with IN
	rows, err := db.Query("SELECT name, salary FROM employees WHERE department_id IN (SELECT id FROM departments WHERE name = 'Engineering')")
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
		log.Printf("  Result: %s earns %.2f (in Engineering)", name, salary)
	}

	log.Printf("✓ Subquery WHERE IN test completed: %d row(s)", rowCount)
}

func testSubqueryWhereNotIn(db *sql.DB) {
	// Subquery in WHERE clause with NOT IN
	rows, err := db.Query("SELECT name, salary FROM employees WHERE department_id NOT IN (SELECT id FROM departments WHERE name = 'HR')")
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
		log.Printf("  Result: %s earns %.2f (not in HR)", name, salary)
	}

	log.Printf("✓ Subquery WHERE NOT IN test completed: %d row(s)", rowCount)
}

func testSubqueryWhereExists(db *sql.DB) {
	// Subquery in WHERE clause with EXISTS
	rows, err := db.Query("SELECT name FROM employees WHERE EXISTS (SELECT * FROM departments WHERE id = employees.department_id AND name = 'Marketing')")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var name string
		err := rows.Scan(&name)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s works in Marketing", name)
	}

	log.Printf("✓ Subquery WHERE EXISTS test completed: %d row(s)", rowCount)
}

func testSubqueryWhereNotExists(db *sql.DB) {
	// Subquery in WHERE clause with NOT EXISTS
	rows, err := db.Query("SELECT name FROM employees WHERE NOT EXISTS (SELECT * FROM salaries WHERE employee_id = employees.id AND amount > 80000)")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var name string
		err := rows.Scan(&name)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s has salary <= 80000", name)
	}

	log.Printf("✓ Subquery WHERE NOT EXISTS test completed: %d row(s)", rowCount)
}

func testSubqueryWhereEquals(db *sql.DB) {
	// Subquery in WHERE clause with =
	rows, err := db.Query("SELECT name, salary FROM employees WHERE salary = (SELECT MAX(salary) FROM employees)")
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
		log.Printf("  Result: %s earns %.2f (highest salary)", name, salary)
	}

	log.Printf("✓ Subquery WHERE = test completed: %d row(s)", rowCount)
}

func testSubqueryWhereGreater(db *sql.DB) {
	// Subquery in WHERE clause with >
	rows, err := db.Query("SELECT name, salary FROM employees WHERE salary > (SELECT AVG(salary) FROM employees)")
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
		log.Printf("  Result: %s earns %.2f (above average)", name, salary)
	}

	log.Printf("✓ Subquery WHERE > test completed: %d row(s)", rowCount)
}

func testSubquerySelectList(db *sql.DB) {
	// Subquery in SELECT list
	rows, err := db.Query("SELECT name, (SELECT AVG(salary) FROM employees) as avg_salary FROM employees")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var name string
		var avgSalary float64
		err := rows.Scan(&name, &avgSalary)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		if rowCount == 1 {
			log.Printf("  Result: Average salary is %.2f", avgSalary)
		}
		log.Printf("  Result: %s (company avg: %.2f)", name, avgSalary)
	}

	log.Printf("✓ Subquery SELECT list test completed: %d row(s)", rowCount)
}

func testSubqueryFrom(db *sql.DB) {
	// Subquery in FROM clause (derived table)
	rows, err := db.Query("SELECT * FROM (SELECT name, salary FROM employees WHERE salary > 70000) as high_earners")
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
		log.Printf("  Result: %s earns %.2f", name, salary)
	}

	log.Printf("✓ Subquery FROM clause test completed: %d row(s)", rowCount)
}

func testCorrelatedSubquery(db *sql.DB) {
	// Correlated subquery
	rows, err := db.Query("SELECT name, salary FROM employees e1 WHERE salary > (SELECT AVG(salary) FROM employees e2 WHERE e2.department_id = e1.department_id)")
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
		log.Printf("  Result: %s earns %.2f (above dept average)", name, salary)
	}

	log.Printf("✓ Correlated subquery test completed: %d row(s)", rowCount)
}

func testNestedSubquery(db *sql.DB) {
	// Nested subquery (subquery within subquery)
	rows, err := db.Query("SELECT name FROM employees WHERE department_id IN (SELECT id FROM departments WHERE id IN (SELECT department_id FROM employees WHERE salary > 75000))")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var name string
		err := rows.Scan(&name)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s (in dept with high earners)", name)
	}

	log.Printf("✓ Nested subquery test completed: %d row(s)", rowCount)
}

func testSubqueryWithJoin(db *sql.DB) {
	// Subquery with JOIN
	rows, err := db.Query("SELECT e.name FROM employees e WHERE e.salary > (SELECT AVG(s.amount) FROM salaries s JOIN employees emp ON emp.id = s.employee_id WHERE emp.department_id = e.department_id)")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var name string
		err := rows.Scan(&name)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s (above salary table avg for dept)", name)
	}

	log.Printf("✓ Subquery with JOIN test completed: %d row(s)", rowCount)
}

func testSubqueryWithGroupBy(db *sql.DB) {
	// Subquery with GROUP BY
	rows, err := db.Query("SELECT name FROM employees WHERE salary > (SELECT MAX(avg_salary) FROM (SELECT AVG(salary) as avg_salary, department_id FROM employees GROUP BY department_id))")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var name string
		err := rows.Scan(&name)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s (above highest dept average)", name)
	}

	log.Printf("✓ Subquery with GROUP BY test completed: %d row(s)", rowCount)
}

func testDropTables(db *sql.DB) {
	_, err := db.Exec("DROP TABLE salaries")
	if err != nil {
		log.Printf("Error dropping salaries table: %v", err)
		return
	}
	log.Println("✓ Dropped table: salaries")

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
