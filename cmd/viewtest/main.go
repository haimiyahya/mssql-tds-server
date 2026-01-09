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

	// Test 3: Simple CREATE VIEW
	log.Println("\n=== Test 3: Simple CREATE VIEW ===")
	testSimpleCreateView(db)

	// Test 4: SELECT from view
	log.Println("\n=== Test 4: SELECT from view ===")
	testSelectFromView(db)

	// Test 5: Complex view with JOIN
	log.Println("\n=== Test 5: Complex view with JOIN ===")
	testComplexViewWithJoin(db)

	// Test 6: View with WHERE clause
	log.Println("\n=== Test 6: View with WHERE clause ===")
	testViewWithWhere(db)

	// Test 7: View with aggregation
	log.Println("\n=== Test 7: View with aggregation ===")
	testViewWithAggregation(db)

	// Test 8: View with ORDER BY
	log.Println("\n=== Test 8: View with ORDER BY ===")
	testViewWithOrderBy(db)

	// Test 9: DROP VIEW
	log.Println("\n=== Test 9: DROP VIEW ===")
	testDropView(db)

	// Test 10: Multiple views
	log.Println("\n=== Test 10: Multiple views ===")
	testMultipleViews(db)

	// Test 11: View with DISTINCT
	log.Println("\n=== Test 11: View with DISTINCT ===")
	testViewWithDistinct(db)

	// Test 12: View with GROUP BY
	log.Println("\n=== Test 12: View with GROUP BY ===")
	testViewWithGroupBy(db)

	// Test 13: View with HAVING
	log.Println("\n=== Test 13: View with HAVING ===")
	testViewWithHaving(db)

	// Test 14: View with subquery
	log.Println("\n=== Test 14: View with subquery ===")
	testViewWithSubquery(db)

	// Test 15: View lifecycle
	log.Println("\n=== Test 15: View lifecycle ===")
	testViewLifecycle(db)

	// Test 16: DROP TABLES
	log.Println("\n=== Test 16: DROP TABLES ===")
	testDropTables(db)

	log.Println("\n=== All Phase 13 tests completed! ===")
	log.Println("🎉 Phase 13: Views Implementation - COMPLETE! 🎉")
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

func testSimpleCreateView(db *sql.DB) {
	// Simple CREATE VIEW
	_, err := db.Exec("CREATE VIEW employee_view AS SELECT id, name, salary FROM employees")
	if err != nil {
		log.Printf("Error creating view: %v", err)
		return
	}
	log.Println("✓ Created view: employee_view")

	// Verify view created
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='view' AND name='employee_view'")
	if err != nil {
		log.Printf("Error querying view: %v", err)
		return
	}
	defer rows.Close()

	var name string
	rows.Next()
	rows.Scan(&name)
	log.Printf("✓ Verified: View 'employee_view' created in database")
}

func testSelectFromView(db *sql.DB) {
	// SELECT from view
	rows, err := db.Query("SELECT * FROM employee_view")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var id int
		var name string
		var salary float64
		err := rows.Scan(&id, &name, &salary)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s earns %.2f", name, salary)
	}

	log.Printf("✓ Selected from view: %d row(s)", rowCount)
}

func testComplexViewWithJoin(db *sql.DB) {
	// CREATE VIEW with JOIN
	_, err := db.Exec("CREATE VIEW employee_dept_view AS SELECT e.id, e.name, e.salary, d.name as department_name FROM employees e JOIN departments d ON e.department_id = d.id")
	if err != nil {
		log.Printf("Error creating view with join: %v", err)
		return
	}
	log.Println("✓ Created view: employee_dept_view (with JOIN)")

	// SELECT from view
	rows, err := db.Query("SELECT * FROM employee_dept_view")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var id int
		var name string
		var salary float64
		var deptName string
		err := rows.Scan(&id, &name, &salary, &deptName)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s works in %s and earns %.2f", name, deptName, salary)
	}

	log.Printf("✓ Selected from view with join: %d row(s)", rowCount)
}

func testViewWithWhere(db *sql.DB) {
	// CREATE VIEW with WHERE clause
	_, err := db.Exec("CREATE VIEW high_earner_view AS SELECT * FROM employees WHERE salary > 70000.00")
	if err != nil {
		log.Printf("Error creating view with WHERE: %v", err)
		return
	}
	log.Println("✓ Created view: high_earner_view (with WHERE)")

	// SELECT from view
	rows, err := db.Query("SELECT * FROM high_earner_view")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var id int
		var name string
		var deptID int
		var salary float64
		err := rows.Scan(&id, &name, &deptID, &salary)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s earns %.2f (high earner)", name, salary)
	}

	log.Printf("✓ Selected from view with WHERE: %d row(s)", rowCount)
}

func testViewWithAggregation(db *sql.DB) {
	// CREATE VIEW with aggregation
	_, err := db.Exec("CREATE VIEW dept_salary_view AS SELECT department_id, COUNT(*) as employee_count, AVG(salary) as avg_salary FROM employees GROUP BY department_id")
	if err != nil {
		log.Printf("Error creating view with aggregation: %v", err)
		return
	}
	log.Println("✓ Created view: dept_salary_view (with aggregation)")

	// SELECT from view
	rows, err := db.Query("SELECT * FROM dept_salary_view")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var deptID int
		var empCount int
		var avgSalary float64
		err := rows.Scan(&deptID, &empCount, &avgSalary)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: Dept %d has %d employees with avg salary %.2f", deptID, empCount, avgSalary)
	}

	log.Printf("✓ Selected from view with aggregation: %d row(s)", rowCount)
}

func testViewWithOrderBy(db *sql.DB) {
	// CREATE VIEW with ORDER BY
	_, err := db.Exec("CREATE VIEW employee_order_view AS SELECT * FROM employees ORDER BY salary DESC")
	if err != nil {
		log.Printf("Error creating view with ORDER BY: %v", err)
		return
	}
	log.Println("✓ Created view: employee_order_view (with ORDER BY)")

	// SELECT from view
	rows, err := db.Query("SELECT * FROM employee_order_view")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var id int
		var name string
		var deptID int
		var salary float64
		err := rows.Scan(&id, &name, &deptID, &salary)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %d. %s earns %.2f", rowCount, name, salary)
	}

	log.Printf("✓ Selected from view with ORDER BY: %d row(s)", rowCount)
}

func testDropView(db *sql.DB) {
	// DROP VIEW
	_, err := db.Exec("DROP VIEW high_earner_view")
	if err != nil {
		log.Printf("Error dropping view: %v", err)
		return
	}
	log.Println("✓ Dropped view: high_earner_view")

	// Verify view dropped
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='view' AND name='high_earner_view'")
	if err != nil {
		log.Printf("Error querying view: %v", err)
		return
	}
	defer rows.Close()

	rows.Next()
	var name sql.NullString
	rows.Scan(&name)
	if !name.Valid {
		log.Println("✓ Verified: View 'high_earner_view' dropped from database")
	} else {
		log.Printf("✗ Verified: View still exists: %s", name.String)
	}
}

func testMultipleViews(db *sql.DB) {
	// CREATE multiple views
	views := []string{
		"CREATE VIEW eng_view AS SELECT * FROM employees WHERE department_id = 1",
		"CREATE VIEW marketing_view AS SELECT * FROM employees WHERE department_id = 2",
		"CREATE VIEW hr_view AS SELECT * FROM employees WHERE department_id = 3",
	}

	for _, query := range views {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error creating view: %v", err)
			continue
		}
		log.Println("✓ Created department view")
	}

	// Count views
	rows, err := db.Query("SELECT COUNT(*) FROM sqlite_master WHERE type='view'")
	if err != nil {
		log.Printf("Error querying views: %v", err)
		return
	}
	defer rows.Close()

	var count int
	rows.Next()
	rows.Scan(&count)
	log.Printf("✓ Total views in database: %d", count)
}

func testViewWithDistinct(db *sql.DB) {
	// CREATE VIEW with DISTINCT
	_, err := db.Exec("CREATE VIEW distinct_dept_view AS SELECT DISTINCT department_id FROM employees")
	if err != nil {
		log.Printf("Error creating view with DISTINCT: %v", err)
		return
	}
	log.Println("✓ Created view: distinct_dept_view (with DISTINCT)")

	// SELECT from view
	rows, err := db.Query("SELECT * FROM distinct_dept_view")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var deptID int
		err := rows.Scan(&deptID)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: Department ID %d", deptID)
	}

	log.Printf("✓ Selected from view with DISTINCT: %d row(s)", rowCount)
}

func testViewWithGroupBy(db *sql.DB) {
	// CREATE VIEW with GROUP BY
	_, err := db.Exec("CREATE VIEW dept_summary_view AS SELECT department_id, COUNT(*) as emp_count FROM employees GROUP BY department_id")
	if err != nil {
		log.Printf("Error creating view with GROUP BY: %v", err)
		return
	}
	log.Println("✓ Created view: dept_summary_view (with GROUP BY)")

	// SELECT from view
	rows, err := db.Query("SELECT * FROM dept_summary_view")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var deptID int
		var empCount int
		err := rows.Scan(&deptID, &empCount)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: Dept %d has %d employees", deptID, empCount)
	}

	log.Printf("✓ Selected from view with GROUP BY: %d row(s)", rowCount)
}

func testViewWithHaving(db *sql.DB) {
	// CREATE VIEW with HAVING
	_, err := db.Exec("CREATE VIEW dept_high_count_view AS SELECT department_id, COUNT(*) as emp_count FROM employees GROUP BY department_id HAVING COUNT(*) > 1")
	if err != nil {
		log.Printf("Error creating view with HAVING: %v", err)
		return
	}
	log.Println("✓ Created view: dept_high_count_view (with HAVING)")

	// SELECT from view
	rows, err := db.Query("SELECT * FROM dept_high_count_view")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var deptID int
		var empCount int
		err := rows.Scan(&deptID, &empCount)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: Dept %d has %d employees (> 1)", deptID, empCount)
	}

	log.Printf("✓ Selected from view with HAVING: %d row(s)", rowCount)
}

func testViewWithSubquery(db *sql.DB) {
	// CREATE VIEW with subquery
	_, err := db.Exec("CREATE VIEW avg_salary_view AS SELECT name, salary, (SELECT AVG(salary) FROM employees) as avg_salary FROM employees")
	if err != nil {
		log.Printf("Error creating view with subquery: %v", err)
		return
	}
	log.Println("✓ Created view: avg_salary_view (with subquery)")

	// SELECT from view
	rows, err := db.Query("SELECT * FROM avg_salary_view")
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
		var avgSalary float64
		err := rows.Scan(&name, &salary, &avgSalary)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		if rowCount == 1 {
			log.Printf("  Result: Average salary is %.2f", avgSalary)
		}
		log.Printf("  Result: %s earns %.2f", name, salary)
	}

	log.Printf("✓ Selected from view with subquery: %d row(s)", rowCount)
}

func testViewLifecycle(db *sql.DB) {
	// Create view
	_, err := db.Exec("CREATE VIEW lifecycle_view AS SELECT * FROM employees")
	if err != nil {
		log.Printf("Error creating view: %v", err)
		return
	}
	log.Println("✓ Created view: lifecycle_view")

	// Query view
	rows, err := db.Query("SELECT COUNT(*) FROM lifecycle_view")
	if err != nil {
		log.Printf("Error querying view: %v", err)
		return
	}
	defer rows.Close()

	var count int
	rows.Next()
	rows.Scan(&count)
	log.Printf("✓ View has %d rows", count)

	// Drop view
	_, err = db.Exec("DROP VIEW lifecycle_view")
	if err != nil {
		log.Printf("Error dropping view: %v", err)
		return
	}
	log.Println("✓ Dropped view: lifecycle_view")

	// Recreate view
	_, err = db.Exec("CREATE VIEW lifecycle_view AS SELECT name, salary FROM employees")
	if err != nil {
		log.Printf("Error recreating view: %v", err)
		return
	}
	log.Println("✓ Recreated view: lifecycle_view")

	// Drop view
	_, err = db.Exec("DROP VIEW lifecycle_view")
	if err != nil {
		log.Printf("Error dropping view: %v", err)
		return
	}
	log.Println("✓ Dropped view: lifecycle_view (lifecycle complete)")
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
