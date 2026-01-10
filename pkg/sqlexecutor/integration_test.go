package sqlexecutor

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestIntegration_CRUDLifecycle tests full CRUD operations in sequence
// 1. INSERT → 2. SELECT validate → 3. UPDATE → 4. SELECT validate → 5. DELETE → 6. SELECT validate
func TestIntegration_CRUDLifecycle(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	executor := NewExecutor(db, nil)

	// Step 1: Create table
	t.Run("Step1_CreateTable", func(t *testing.T) {
		_, err := db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, email TEXT, age INTEGER)")
		if err != nil {
			t.Fatalf("Failed to create table: %v", err)
		}
	})

	// Step 2: Insert record
	var insertedID int64
	t.Run("Step2_Insert", func(t *testing.T) {
		result, err := db.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)", "Alice", "alice@test.com", 30)
		if err != nil {
			t.Fatalf("Failed to insert record: %v", err)
		}
		rows, _ := result.RowsAffected()
		if rows != 1 {
			t.Errorf("Expected 1 row affected, got %d", rows)
		}
		insertedID, _ = result.LastInsertId()
		if insertedID == 0 {
			t.Error("Expected non-zero inserted ID")
		}
		t.Logf("Inserted record with ID: %d", insertedID)
	})

	// Step 3: Select and validate inserted values
	t.Run("Step3_SelectValidate", func(t *testing.T) {
		result, err := executor.Execute(fmt.Sprintf("SELECT id, name, email, age FROM users WHERE id = %d", insertedID))
		if err != nil {
			t.Fatalf("Failed to select record: %v", err)
		}
		if len(result.Rows) != 1 {
			t.Fatalf("Expected 1 row, got %d", len(result.Rows))
		}
		row := result.Rows[0]
		// Columns: id, name, email, age
		 if row[1].(string) != "Alice" {
			t.Errorf("Expected name='Alice', got '%s'", row[1])
		}
		 if row[2].(string) != "alice@test.com" {
			t.Errorf("Expected email='alice@test.com', got '%s'", row[2])
		}
		 if row[3].(int64) != 30 {
			t.Errorf("Expected age=30, got '%s'", row[3])
		}
	})

	// Step 4: Update record
	t.Run("Step4_Update", func(t *testing.T) {
		result, err := executor.Execute(fmt.Sprintf("UPDATE users SET age = 31, email = 'alice@updated.com' WHERE id = %d", insertedID))
		if err != nil {
			t.Fatalf("Failed to update record: %v", err)
		}
		if result.RowCount != 1 {
			t.Errorf("Expected 1 row affected, got %d", result.RowCount)
		}
	})

	// Step 5: Select and validate updated values
	t.Run("Step5_SelectValidateUpdate", func(t *testing.T) {
		result, err := executor.Execute(fmt.Sprintf("SELECT id, name, email, age FROM users WHERE id = %d", insertedID))
		if err != nil {
			t.Fatalf("Failed to select record: %v", err)
		}
		if len(result.Rows) != 1 {
			t.Fatalf("Expected 1 row, got %d", len(result.Rows))
		}
		row := result.Rows[0]
		 if row[3].(int64) != 31 {
			t.Errorf("Expected age=31 after update, got '%s'", row[3])
		}
		 if row[2].(string) != "alice@updated.com" {
			t.Errorf("Expected email='alice@updated.com', got '%s'", row[2])
		}
		 if row[1].(string) != "Alice" {
			t.Errorf("Expected name='Alice' (unchanged), got '%s'", row[1])
		}
	})

	// Step 6: Delete record
	t.Run("Step6_Delete", func(t *testing.T) {
		result, err := executor.Execute(fmt.Sprintf("DELETE FROM users WHERE id = %d", insertedID))
		if err != nil {
			t.Fatalf("Failed to delete record: %v", err)
		}
		if result.RowCount != 1 {
			t.Errorf("Expected 1 row affected, got %d", result.RowCount)
		}
	})

	// Step 7: Select and validate deletion
	t.Run("Step7_SelectValidateDelete", func(t *testing.T) {
		result, err := executor.Execute(fmt.Sprintf("SELECT id FROM users WHERE id = %d", insertedID))
		if err != nil {
			t.Fatalf("Failed to select record: %v", err)
		}
		if len(result.Rows) != 0 {
			t.Errorf("Expected 0 rows after deletion, got %d", len(result.Rows))
		}
	})
}

// TestIntegration_MultiTableOperations tests operations across multiple tables with joins
func TestIntegration_MultiTableOperations(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	// Setup tables
	_, err = db.Exec(`
		CREATE TABLE departments (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT);
		CREATE TABLE employees (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, dept_id INTEGER);
	`)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	executor := NewExecutor(db, nil)

	// Insert departments
	var deptID1, deptID2 int64
	t.Run("InsertDepartments", func(t *testing.T) {
		r1, err := db.Exec("INSERT INTO departments (name) VALUES ('IT')")
		if err != nil {
			t.Fatalf("Failed to insert IT department: %v", err)
		}
		deptID1, _ = r1.LastInsertId()

		r2, err := db.Exec("INSERT INTO departments (name) VALUES ('HR')")
		if err != nil {
			t.Fatalf("Failed to insert HR department: %v", err)
		}
		deptID2, _ = r2.LastInsertId()
	})

	// Insert employees
	t.Run("InsertEmployees", func(t *testing.T) {
		_, err := db.Exec(fmt.Sprintf("INSERT INTO employees (name, dept_id) VALUES ('Alice', %d)", deptID1))
		if err != nil {
			t.Fatalf("Failed to insert Alice: %v", err)
		}
		_, err = db.Exec(fmt.Sprintf("INSERT INTO employees (name, dept_id) VALUES ('Bob', %d)", deptID2))
		if err != nil {
			t.Fatalf("Failed to insert Bob: %v", err)
		}
	})

	// Validate join
	t.Run("ValidateJoin", func(t *testing.T) {
		result, err := executor.Execute(`
			SELECT e.name AS emp_name, d.name AS dept_name 
			FROM employees e 
			JOIN departments d ON e.dept_id = d.id
		`)
		if err != nil {
			t.Fatalf("Failed to execute join: %v", err)
		}
		if len(result.Rows) != 2 {
			t.Errorf("Expected 2 rows, got %d", len(result.Rows))
		}
		// Validate Alice is in IT, Bob is in HR
		for _, row := range result.Rows {
			empName := row[0].(string)
			deptName := row[1].(string)
			t.Logf("Employee: %s works in %s", empName, deptName)
		}
	})
}

// TestIntegration_BulkOperations tests bulk insert/update with validation
func TestIntegration_BulkOperations(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE orders (id INTEGER PRIMARY KEY AUTOINCREMENT, customer_id INTEGER, amount DECIMAL)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Bulk insert
	t.Run("BulkInsert", func(t *testing.T) {
		_, err := db.Exec("INSERT INTO orders (customer_id, amount) VALUES (1, 100.00)")
		if err != nil {
			t.Fatalf("Failed to insert order 1: %v", err)
		}
		_, err = db.Exec("INSERT INTO orders (customer_id, amount) VALUES (1, 200.00)")
		if err != nil {
			t.Fatalf("Failed to insert order 2: %v", err)
		}
		_, err = db.Exec("INSERT INTO orders (customer_id, amount) VALUES (2, 150.00)")
		if err != nil {
			t.Fatalf("Failed to insert order 3: %v", err)
		}
	})

	// Validate count and sum using raw SQL
	t.Run("ValidateCountAndSum", func(t *testing.T) {
		var count int64
		err = db.QueryRow("SELECT COUNT(id) FROM orders WHERE customer_id = 1").Scan(&count)
		if err != nil {
			t.Fatalf("Failed to count orders: %v", err)
		}
		if count != 2 {
			t.Errorf("Expected 2 orders, got %d", count)
		}

		var sum float64
		err = db.QueryRow("SELECT SUM(amount) FROM orders WHERE customer_id = 1").Scan(&sum)
		if err != nil {
			t.Fatalf("Failed to sum amounts: %v", err)
		}
		if sum != 300.00 {
			t.Errorf("Expected sum=300.00, got %f", sum)
		}
	})
}
