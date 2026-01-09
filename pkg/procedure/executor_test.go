package procedure

import (
	"os"
	"testing"

	"github.com/factory/mssql-tds-server/pkg/sqlite"
)

func setupExecutor(t *testing.T) (*Executor, *Storage, func()) {
	// Create temporary database
	dbPath := "/tmp/test_executor.db"
	os.Remove(dbPath) // Clean up any existing file

	db, err := sqlite.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	err = db.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Create a test table for queries
	_, err = db.GetDB().Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER, name TEXT, department TEXT)")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Insert test data
	_, err = db.GetDB().Exec("INSERT INTO users (id, name, department) VALUES (1, 'Alice', 'Engineering')")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}
	_, err = db.GetDB().Exec("INSERT INTO users (id, name, department) VALUES (2, 'Bob', 'Marketing')")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	storage, err := NewStorage(db)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	executor, err := NewExecutor(db, storage)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Cleanup function
	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return executor, storage, cleanup
}

func TestExecutor_Execute(t *testing.T) {
	executor, storage, cleanup := setupExecutor(t)
	defer cleanup()

	// Create a procedure
	proc := &Procedure{
		Name: "GET_USER_BY_ID",
		Body: "SELECT * FROM users WHERE id = @id",
		Parameters: []Parameter{
			{Name: "id", Type: "INT"},
		},
	}

	err := storage.Create(proc)
	if err != nil {
		t.Fatalf("Failed to create procedure: %v", err)
	}

	// Execute the procedure
	paramValues := map[string]interface{}{
		"@id": 1,
	}

	results, err := executor.Execute("GET_USER_BY_ID", paramValues)
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	// Check results (should have header + 1 data row)
	if len(results) != 2 {
		t.Errorf("Execute() returned %d rows, want 2", len(results))
	}
}

func TestExecutor_Execute_MissingParameter(t *testing.T) {
	executor, storage, cleanup := setupExecutor(t)
	defer cleanup()

	// Create a procedure with required parameter
	proc := &Procedure{
		Name: "GET_USER",
		Body: "SELECT * FROM users WHERE id = @id",
		Parameters: []Parameter{
			{Name: "id", Type: "INT"},
		},
	}

	err := storage.Create(proc)
	if err != nil {
		t.Fatalf("Failed to create procedure: %v", err)
	}

	// Execute without required parameter
	paramValues := map[string]interface{}{}

	_, err = executor.Execute("GET_USER", paramValues)
	if err == nil {
		t.Error("Execute() expected error for missing required parameter")
	}
}

func TestExecutor_Execute_NonexistentProcedure(t *testing.T) {
	executor, _, cleanup := setupExecutor(t)
	defer cleanup()

	paramValues := map[string]interface{}{
		"@id": 1,
	}

	_, err := executor.Execute("NONEXISTENT", paramValues)
	if err == nil {
		t.Error("Execute() expected error for non-existent procedure")
	}
}

func TestExecutor_Execute_StringParameter(t *testing.T) {
	executor, storage, cleanup := setupExecutor(t)
	defer cleanup()

	// Create a procedure with string parameter
	proc := &Procedure{
		Name: "GET_USERS_BY_DEPT",
		Body: "SELECT * FROM users WHERE department = @dept",
		Parameters: []Parameter{
			{Name: "dept", Type: "VARCHAR"},
		},
	}

	err := storage.Create(proc)
	if err != nil {
		t.Fatalf("Failed to create procedure: %v", err)
	}

	// Execute with string parameter
	paramValues := map[string]interface{}{
		"@dept": "Engineering",
	}

	results, err := executor.Execute("GET_USERS_BY_DEPT", paramValues)
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	// Check results
	if len(results) < 2 {
		t.Errorf("Execute() returned %d rows, want at least 2", len(results))
	}
}

func TestExecutor_Execute_StringWithQuotes(t *testing.T) {
	executor, storage, cleanup := setupExecutor(t)
	defer cleanup()

	// Create a procedure
	proc := &Procedure{
		Name: "GET_BY_NAME",
		Body: "SELECT * FROM users WHERE name = @name",
		Parameters: []Parameter{
			{Name: "name", Type: "VARCHAR"},
		},
	}

	err := storage.Create(proc)
	if err != nil {
		t.Fatalf("Failed to create procedure: %v", err)
	}

	// Execute with name containing special characters
	paramValues := map[string]interface{}{
		"@name": "Alice",
	}

	results, err := executor.Execute("GET_BY_NAME", paramValues)
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	// Check results
	if len(results) < 2 {
		t.Errorf("Execute() returned %d rows, want at least 2", len(results))
	}
}

func TestExecutor_ValidateParameters(t *testing.T) {
	executor, _, cleanup := setupExecutor(t)
	defer cleanup()

	// Create a test procedure
	proc := &Procedure{
		Name: "TEST_PROC",
		Body: "SELECT 1",
		Parameters: []Parameter{
			{Name: "id", Type: "INT"},
			{Name: "name", Type: "VARCHAR"},
		},
	}

	tests := []struct {
		name       string
		params     map[string]interface{}
		wantErr    bool
		errContains string
	}{
		{
			name:    "All required parameters provided",
			params:  map[string]interface{}{"@id": 1, "@name": "Alice"},
			wantErr: false,
		},
		{
			name:       "Missing required parameter",
			params:     map[string]interface{}{"@id": 1},
			wantErr:    true,
			errContains: "missing required parameter",
		},
		{
			name:       "Extra parameter provided",
			params:     map[string]interface{}{"@id": 1, "@name": "Alice", "@extra": "value"},
			wantErr:    true,
			errContains: "unexpected parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := executor.validateParameters(proc, tt.params)

			if tt.wantErr {
				if err == nil {
					t.Errorf("validateParameters() expected error but got none")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("validateParameters() error = %v, want error containing %s", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("validateParameters() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestExecutor_ReplaceParameters(t *testing.T) {
	executor, _, cleanup := setupExecutor(t)
	defer cleanup()

	tests := []struct {
		name  string
		sql   string
		params map[string]interface{}
		want  string
	}{
		{
			name:  "Replace single parameter",
			sql:   "SELECT * FROM users WHERE id = @id",
			params: map[string]interface{}{"@id": 1},
			want:  "SELECT * FROM users WHERE id = 1",
		},
		{
			name:  "Replace string parameter",
			sql:   "SELECT * FROM users WHERE name = @name",
			params: map[string]interface{}{"@name": "Alice"},
			want:  "SELECT * FROM users WHERE name = 'Alice'",
		},
		{
			name:  "Replace multiple parameters",
			sql:   "SELECT * FROM users WHERE id = @id AND department = @dept",
			params: map[string]interface{}{"@id": 1, "@dept": "Engineering"},
			want:  "SELECT * FROM users WHERE id = 1 AND department = 'Engineering'",
		},
		{
			name:  "String with single quotes",
			sql:   "SELECT * FROM users WHERE name = @name",
			params: map[string]interface{}{"@name": "O'Brien"},
			want:  "SELECT * FROM users WHERE name = 'O''Brien'",
		},
		{
			name:  "NULL value",
			sql:   "SELECT * FROM users WHERE name = @name",
			params: map[string]interface{}{"@name": nil},
			want:  "SELECT * FROM users WHERE name = NULL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := executor.replaceParameters(tt.sql, tt.params)
			if err != nil {
				t.Errorf("replaceParameters() error = %v", err)
				return
			}
			if result != tt.want {
				t.Errorf("replaceParameters() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestExecutor_FormatValue(t *testing.T) {
	executor, _, cleanup := setupExecutor(t)
	defer cleanup()

	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{
			name:  "Integer",
			value: 123,
			want:  "123",
		},
		{
			name:  "String",
			value: "Alice",
			want:  "'Alice'",
		},
		{
			name:  "String with quotes",
			value: "O'Brien",
			want:  "'O''Brien'",
		},
		{
			name:  "Nil",
			value: nil,
			want:  "NULL",
		},
		{
			name:  "Boolean true",
			value: true,
			want:  "1",
		},
		{
			name:  "Boolean false",
			value: false,
			want:  "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.formatValue(tt.value)
			if result != tt.want {
				t.Errorf("formatValue() = %v, want %v", result, tt.want)
			}
		})
	}
}
