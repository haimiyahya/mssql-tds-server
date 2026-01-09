package sqlexecutor

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	return db
}

func TestConvertValueToString(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{nil, "NULL"},
		{"hello", "hello"},
		{123, "123"},
		{45.67, "45.67"},
		{true, "true"},
		{[]byte("bytes"), "bytes"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := ConvertValueToString(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertValueToString(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestReplaceAll(t *testing.T) {
	tests := []struct {
		input    string
		old      string
		new      string
		expected string
	}{
		{"hello world", "world", "universe", "hello universe"},
		{"[users]", "[", "", "users]"},
		{"[users]", "]", "", "[users"},
		{"no match", "xyz", "abc", "no match"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := replaceAll(tt.input, tt.old, tt.new)
			if result != tt.expected {
				t.Errorf("replaceAll(%q, %q, %q) = %q, want %q", tt.input, tt.old, tt.new, result, tt.expected)
			}
		})
	}
}

func TestFindSubstring(t *testing.T) {
	tests := []struct {
		input    string
		substr   string
		expected int
	}{
		{"hello world", "world", 6},
		{"hello world", "hello", 0},
		{"hello world", "xyz", -1},
		{"", "test", -1},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := findSubstring(tt.input, tt.substr)
			if result != tt.expected {
				t.Errorf("findSubstring(%q, %q) = %d, want %d", tt.input, tt.substr, result, tt.expected)
			}
		})
	}
}

func TestExecutorSelect(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test table
	_, err := db.Exec("CREATE TABLE test (id INTEGER, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Insert test data
	_, err = db.Exec("INSERT INTO test (id, name) VALUES (1, 'Alice')")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	executor := NewExecutor(db)
	result, err := executor.Execute("SELECT * FROM test")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.RowCount != 1 {
		t.Errorf("RowCount = %d, want 1", result.RowCount)
	}

	if !result.IsQuery {
		t.Error("IsQuery should be true for SELECT")
	}

	if len(result.Rows) != 1 {
		t.Errorf("Rows length = %d, want 1", len(result.Rows))
	}
}

func TestExecutorInsert(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test table
	_, err := db.Exec("CREATE TABLE test (id INTEGER, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	executor := NewExecutor(db)
	result, err := executor.Execute("INSERT INTO test (id, name) VALUES (1, 'Alice')")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.RowCount != 1 {
		t.Errorf("RowCount = %d, want 1", result.RowCount)
	}

	if result.IsQuery {
		t.Error("IsQuery should be false for INSERT")
	}

	if result.Message == "" {
		t.Error("Message should not be empty for INSERT")
	}
}

func TestExecutorUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test table
	_, err := db.Exec("CREATE TABLE test (id INTEGER, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Insert test data
	_, err = db.Exec("INSERT INTO test (id, name) VALUES (1, 'Alice')")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	executor := NewExecutor(db)
	result, err := executor.Execute("UPDATE test SET name = 'Bob' WHERE id = 1")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.RowCount != 1 {
		t.Errorf("RowCount = %d, want 1", result.RowCount)
	}

	if result.IsQuery {
		t.Error("IsQuery should be false for UPDATE")
	}
}

func TestExecutorDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test table
	_, err := db.Exec("CREATE TABLE test (id INTEGER, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Insert test data
	_, err = db.Exec("INSERT INTO test (id, name) VALUES (1, 'Alice')")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	executor := NewExecutor(db)
	result, err := executor.Execute("DELETE FROM test WHERE id = 1")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.RowCount != 1 {
		t.Errorf("RowCount = %d, want 1", result.RowCount)
	}

	if result.IsQuery {
		t.Error("IsQuery should be false for DELETE")
	}
}

func TestExecutorCreateTable(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	executor := NewExecutor(db)
	result, err := executor.Execute("CREATE TABLE test (id INTEGER, name TEXT)")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.RowCount != 0 {
		t.Errorf("RowCount = %d, want 0", result.RowCount)
	}

	if result.IsQuery {
		t.Error("IsQuery should be false for CREATE TABLE")
	}

	if result.Message == "" {
		t.Error("Message should not be empty for CREATE TABLE")
	}
}

func TestExecutorDropTable(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create test table
	_, err := db.Exec("CREATE TABLE test (id INTEGER, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	executor := NewExecutor(db)
	result, err := executor.Execute("DROP TABLE test")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if result.RowCount != 0 {
		t.Errorf("RowCount = %d, want 0", result.RowCount)
	}

	if result.IsQuery {
		t.Error("IsQuery should be false for DROP TABLE")
	}

	if result.Message == "" {
		t.Error("Message should not be empty for DROP TABLE")
	}
}

func TestConvertCreateTable(t *testing.T) {
	executor := NewExecutor(nil)

	tests := []struct {
		input    string
		expected string
	}{
		{
			"CREATE TABLE [users] (id INT)",
			"CREATE TABLE users (id INT)",
		},
		{
			"CREATE TABLE dbo.users (id INT)",
			"CREATE TABLE dbo.users (id INT)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := executor.convertCreateTable(tt.input)
			if result != tt.expected {
				t.Errorf("convertCreateTable(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
