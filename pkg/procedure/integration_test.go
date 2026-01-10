package procedure

import (
	"os"
	"testing"

	"github.com/factory/mssql-tds-server/pkg/sqlite"
)

// TestIntegration_ProcedureLifecycle tests complete stored procedure lifecycle
// Create procedure → Execute (INSERT) → Validate SELECT → Execute (UPDATE) → Validate SELECT → Execute (DELETE) → Validate deletion
func TestIntegration_ProcedureLifecycle(t *testing.T) {
	debugLog(t, "TestIntegration_ProcedureLifecycle: START")
	
	// Create temporary database
	dbPath := "/tmp/test_procedure_lifecycle.db"
	os.Remove(dbPath) // Clean up any existing file

	db, err := sqlite.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove(dbPath)
	}()

	err = db.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Create test table
	_, err = db.GetDB().Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			name TEXT, 
			email TEXT, 
			age INTEGER
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	storage, err := NewStorage(db)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	executor, err := NewExecutor(db, storage)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	var insertedID int64

	// Step 1: Create INSERT procedure
	t.Run("Step1_CreateInsertProcedure", func(t *testing.T) {
		proc := &Procedure{
			Name: "InsertUser",
			Body: "INSERT INTO users (name, email, age) VALUES (@name, @email, @age)",
			Parameters: []Parameter{
				{Name: "name", Type: "VARCHAR", Length: 100},
				{Name: "email", Type: "VARCHAR", Length: 255},
				{Name: "age", Type: "INT"},
			},
		}
		err := storage.Create(proc)
		if err != nil {
			t.Fatalf("Failed to create InsertUser procedure: %v", err)
		}
	})

	// Step 2: Execute INSERT procedure
	t.Run("Step2_ExecuteInsertProcedure", func(t *testing.T) {
		// Get last inserted ID before insert
		var maxID int64
		db.GetDB().QueryRow("SELECT MAX(id) FROM users").Scan(&maxID)
		
		params := map[string]interface{}{
			"@name":  "Alice",
			"@email": "alice@test.com",
			"@age":   30,
		}
		
		_, err := executor.Execute("InsertUser", params)
		if err != nil {
			t.Fatalf("Failed to execute InsertUser: %v", err)
		}
		
		// Get last inserted ID after insert
		db.GetDB().QueryRow("SELECT MAX(id) FROM users").Scan(&insertedID)
		if insertedID <= maxID {
			t.Errorf("Expected new ID > %d, got %d", maxID, insertedID)
		}
		t.Logf("Inserted user with ID: %d", insertedID)
	})

	// Step 3: Validate SELECT
	t.Run("Step3_ValidateInsert", func(t *testing.T) {
		// Create SELECT procedure
		selectProc := &Procedure{
			Name: "GetUserByID",
			Body: "SELECT * FROM users WHERE id = @id",
			Parameters: []Parameter{
				{Name: "id", Type: "INT"},
			},
		}
		err := storage.Create(selectProc)
		if err != nil {
			t.Fatalf("Failed to create GetUserByID procedure: %v", err)
		}

		// Execute SELECT
		params := map[string]interface{}{
			"@id": int(insertedID),
		}
		
		results, err := executor.Execute("GetUserByID", params)
		if err != nil {
			t.Fatalf("Failed to execute GetUserByID: %v", err)
		}

		// Results are [][]string: [header row, data rows]
		if len(results) < 2 {
			t.Fatalf("Expected at least 2 results (header + data), got %d", len(results))
		}

		dataRow := results[1] // Skip header row
		// Columns: id, name, email, age
		if dataRow[1] != "Alice" {
			t.Errorf("Expected name='Alice', got '%s'", dataRow[1])
		}
		if dataRow[2] != "alice@test.com" {
			t.Errorf("Expected email='alice@test.com', got '%s'", dataRow[2])
		}
		if dataRow[3] != "30" {
			t.Errorf("Expected age=30, got '%s'", dataRow[3])
		}
	})

	// Step 4: Create UPDATE procedure
	t.Run("Step4_CreateUpdateProcedure", func(t *testing.T) {
		updateProc := &Procedure{
			Name: "UpdateUserAge",
			Body: "UPDATE users SET age = @age, email = @email WHERE id = @id",
			Parameters: []Parameter{
				{Name: "id", Type: "INT"},
				{Name: "age", Type: "INT"},
				{Name: "email", Type: "VARCHAR", Length: 255},
			},
		}
		err := storage.Create(updateProc)
		if err != nil {
			t.Fatalf("Failed to create UpdateUserAge procedure: %v", err)
		}
	})

	// Step 5: Execute UPDATE procedure
	t.Run("Step5_ExecuteUpdateProcedure", func(t *testing.T) {
		params := map[string]interface{}{
			"@id":    int(insertedID),
			"@age":   31,
			"@email": "alice.updated.com",
		}
		
		_, err := executor.Execute("UpdateUserAge", params)
		if err != nil {
			t.Fatalf("Failed to execute UpdateUserAge: %v", err)
		}
	})

	// Step 6: Validate UPDATE
	t.Run("Step6_ValidateUpdate", func(t *testing.T) {
		params := map[string]interface{}{
			"@id": int(insertedID),
		}
		
		results, err := executor.Execute("GetUserByID", params)
		if err != nil {
			t.Fatalf("Failed to execute GetUserByID: %v", err)
		}

		if len(results) < 2 {
			t.Fatalf("Expected at least 2 results, got %d", len(results))
		}

		dataRow := results[1]
		if dataRow[3] != "31" {
			t.Errorf("Expected age=31 after update, got '%s'", dataRow[3])
		}
		if dataRow[2] != "alice.updated.com" {
			t.Errorf("Expected email='alice@updated.com', got '%s'", dataRow[2])
		}
		if dataRow[1] != "Alice" {
			t.Errorf("Expected name='Alice' (unchanged), got '%s'", dataRow[1])
		}
	})

	// Step 7: Create DELETE procedure
	t.Run("Step7_CreateDeleteProcedure", func(t *testing.T) {
		deleteProc := &Procedure{
			Name: "DeleteUser",
			Body: "DELETE FROM users WHERE id = @id",
			Parameters: []Parameter{
				{Name: "id", Type: "INT"},
			},
		}
		err := storage.Create(deleteProc)
		if err != nil {
			t.Fatalf("Failed to create DeleteUser procedure: %v", err)
		}
	})

	// Step 8: Execute DELETE procedure
	t.Run("Step8_ExecuteDeleteProcedure", func(t *testing.T) {
		params := map[string]interface{}{
			"@id": int(insertedID),
		}
		
		_, err := executor.Execute("DeleteUser", params)
		if err != nil {
			t.Fatalf("Failed to execute DeleteUser: %v", err)
		}
	})

	// Step 9: Validate DELETE
	t.Run("Step9_ValidateDelete", func(t *testing.T) {
		params := map[string]interface{}{
			"@id": int(insertedID),
		}
		
		results, err := executor.Execute("GetUserByID", params)
		if err != nil {
			t.Fatalf("Failed to execute GetUserByID: %v", err)
		}

		// Should only have header row, no data
		if len(results) != 1 {
			t.Errorf("Expected 1 result (header only), got %d", len(results))
		}
	})

	debugLog(t, "TestIntegration_ProcedureLifecycle: END")
}

// TestIntegration_MultiProcedureChain tests multiple procedures working together
func TestIntegration_MultiProcedureChain(t *testing.T) {
	debugLog(t, "TestIntegration_MultiProcedureChain: START")
	
	// Create temporary database
	dbPath := "/tmp/test_multi_procedure.db"
	os.Remove(dbPath)

	db, err := sqlite.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove(dbPath)
	}()

	err = db.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Create tables
	_, err = db.GetDB().Exec(`
		CREATE TABLE departments (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT);
		CREATE TABLE employees (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, dept_id INTEGER);
	`)
	if err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}

	storage, err := NewStorage(db)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	executor, err := NewExecutor(db, storage)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	var deptID int64

	// Create and execute department insertion
	t.Run("CreateDepartmentProcedure", func(t *testing.T) {
		proc := &Procedure{
			Name: "CreateDepartment",
			Body: "INSERT INTO departments (name) VALUES (@name)",
			Parameters: []Parameter{
				{Name: "name", Type: "VARCHAR", Length: 100},
			},
		}
		err := storage.Create(proc)
		if err != nil {
			t.Fatalf("Failed to create CreateDepartment procedure: %v", err)
		}

		// Get max ID before
		var maxID int64
		db.GetDB().QueryRow("SELECT MAX(id) FROM departments").Scan(&maxID)

		params := map[string]interface{}{
			"@name": "IT",
		}
		
		_, err = executor.Execute("CreateDepartment", params)
		if err != nil {
			t.Fatalf("Failed to execute CreateDepartment: %v", err)
		}

		// Get max ID after
		db.GetDB().QueryRow("SELECT MAX(id) FROM departments").Scan(&deptID)
		if deptID <= maxID {
			t.Errorf("Expected new ID > %d, got %d", maxID, deptID)
		}
	})

	// Create and execute employee insertion
	t.Run("CreateEmployeeProcedure", func(t *testing.T) {
		proc := &Procedure{
			Name: "CreateEmployee",
			Body: "INSERT INTO employees (name, dept_id) VALUES (@name, @dept_id)",
			Parameters: []Parameter{
				{Name: "name", Type: "VARCHAR", Length: 100},
				{Name: "dept_id", Type: "INT"},
			},
		}
		err := storage.Create(proc)
		if err != nil {
			t.Fatalf("Failed to create CreateEmployee procedure: %v", err)
		}

		params := map[string]interface{}{
			"@name":    "Alice",
			"@dept_id": int(deptID),
		}
		
		_, err = executor.Execute("CreateEmployee", params)
		if err != nil {
			t.Fatalf("Failed to execute CreateEmployee: %v", err)
		}
	})

	// Validate join using procedure
	t.Run("CreateJoinProcedure", func(t *testing.T) {
		proc := &Procedure{
			Name: "GetEmployeesByDepartment",
			Body: "SELECT e.name AS emp_name, d.name AS dept_name FROM employees e JOIN departments d ON e.dept_id = d.id WHERE d.name = @dept_name",
			Parameters: []Parameter{
				{Name: "dept_name", Type: "VARCHAR", Length: 100},
			},
		}
		err := storage.Create(proc)
		if err != nil {
			t.Fatalf("Failed to create GetEmployeesByDepartment procedure: %v", err)
		}

		params := map[string]interface{}{
			"@dept_name": "IT",
		}
		
		results, err := executor.Execute("GetEmployeesByDepartment", params)
		if err != nil {
			t.Fatalf("Failed to execute GetEmployeesByDepartment: %v", err)
		}

		// Should have header + 1 data row
		if len(results) < 2 {
			t.Errorf("Expected at least 2 results, got %d", len(results))
		}

		// Validate employee data
		dataRow := results[1]
		if dataRow[0] != "Alice" {
			t.Errorf("Expected emp_name='Alice', got '%s'", dataRow[0])
		}
		if dataRow[1] != "IT" {
			t.Errorf("Expected dept_name='IT', got '%s'", dataRow[1])
		}
	})

	debugLog(t, "TestIntegration_MultiProcedureChain: END")
}

// TestIntegration_BulkOperations tests bulk operations with stored procedures
func TestIntegration_BulkOperations(t *testing.T) {
	debugLog(t, "TestIntegration_BulkOperations: START")
	
	dbPath := "/tmp/test_bulk_procedure.db"
	os.Remove(dbPath)

	db, err := sqlite.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove(dbPath)
	}()

	err = db.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	_, err = db.GetDB().Exec("CREATE TABLE orders (id INTEGER PRIMARY KEY AUTOINCREMENT, customer_id INTEGER, amount DECIMAL)")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	storage, err := NewStorage(db)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	executor, err := NewExecutor(db, storage)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Create insert procedure
	t.Run("CreateInsertOrderProcedure", func(t *testing.T) {
		proc := &Procedure{
			Name: "InsertOrder",
			Body: "INSERT INTO orders (customer_id, amount) VALUES (@customer_id, @amount)",
			Parameters: []Parameter{
				{Name: "customer_id", Type: "INT"},
				{Name: "amount", Type: "DECIMAL"},
			},
		}
		err := storage.Create(proc)
		if err != nil {
			t.Fatalf("Failed to create InsertOrder procedure: %v", err)
		}
	})

	// Bulk insert
	t.Run("BulkInsert", func(t *testing.T) {
		_, err := executor.Execute("InsertOrder", map[string]interface{}{"@customer_id": 1, "@amount": 100.00})
		if err != nil {
			t.Fatalf("Failed to insert order 1: %v", err)
		}
		_, err = executor.Execute("InsertOrder", map[string]interface{}{"@customer_id": 1, "@amount": 200.00})
		if err != nil {
			t.Fatalf("Failed to insert order 2: %v", err)
		}
		_, err = executor.Execute("InsertOrder", map[string]interface{}{"@customer_id": 2, "@amount": 150.00})
		if err != nil {
			t.Fatalf("Failed to insert order 3: %v", err)
		}
	})

	// Validate count and sum using raw SQL
	t.Run("ValidateCountAndSum", func(t *testing.T) {
		var count int64
		db.GetDB().QueryRow("SELECT COUNT(id) FROM orders WHERE customer_id = 1").Scan(&count)
		if count != 2 {
			t.Errorf("Expected 2 orders, got %d", count)
		}

		var sum float64
		db.GetDB().QueryRow("SELECT SUM(amount) FROM orders WHERE customer_id = 1").Scan(&sum)
		if sum != 300.00 {
			t.Errorf("Expected sum=300.00, got %f", sum)
		}
	})

	// Create update procedure
	t.Run("CreateUpdateOrderProcedure", func(t *testing.T) {
		proc := &Procedure{
			Name: "UpdateOrderAmount",
			Body: "UPDATE orders SET amount = amount * @multiplier WHERE customer_id = @customer_id",
			Parameters: []Parameter{
				{Name: "customer_id", Type: "INT"},
				{Name: "multiplier", Type: "DECIMAL"},
			},
		}
		err := storage.Create(proc)
		if err != nil {
			t.Fatalf("Failed to create UpdateOrderAmount procedure: %v", err)
		}
	})

	// Update and validate
	t.Run("UpdateAndValidate", func(t *testing.T) {
		_, err := executor.Execute("UpdateOrderAmount", map[string]interface{}{"@customer_id": 1, "@multiplier": 1.1})
		if err != nil {
			t.Fatalf("Failed to update orders: %v", err)
		}

		var sum float64
		db.GetDB().QueryRow("SELECT SUM(amount) FROM orders WHERE customer_id = 1").Scan(&sum)
		if sum < 329.99 || sum > 330.01 {
			t.Errorf("Expected sum=330.00 after 10%% increase, got %f", sum)
		}
	})

	debugLog(t, "TestIntegration_BulkOperations: END")
}
