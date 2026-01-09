package sqlparser

import (
	"testing"
)

func TestParseSelect(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name     string
		query    string
		expected StatementType
		table    string
		columns  []string
	}{
		{
			name:     "simple select",
			query:    "SELECT * FROM users",
			expected: StatementTypeSelect,
			table:    "users",
			columns:  []string{"*"},
		},
		{
			name:     "select with columns",
			query:    "SELECT name, email FROM users",
			expected: StatementTypeSelect,
			table:    "users",
			columns:  []string{"name", "email"},
		},
		{
			name:     "select with where",
			query:    "SELECT * FROM users WHERE id = 1",
			expected: StatementTypeSelect,
			table:    "users",
			columns:  []string{"*"},
		},
		{
			name:     "select with schema",
			query:    "SELECT * FROM dbo.users",
			expected: StatementTypeSelect,
			table:    "users",
			columns:  []string{"*"},
		},
		{
			name:     "select with brackets",
			query:    "SELECT * FROM [users]",
			expected: StatementTypeSelect,
			table:    "users",
			columns:  []string{"*"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := parser.Parse(tt.query)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if stmt.Type != tt.expected {
				t.Errorf("Type = %v, want %v", stmt.Type, tt.expected)
			}

			if stmt.Select == nil {
				t.Fatal("Select should not be nil")
			}

			if stmt.Select.Table != tt.table {
				t.Errorf("Table = %v, want %v", stmt.Select.Table, tt.table)
			}

			if len(stmt.Select.Columns) != len(tt.columns) {
				t.Errorf("Columns length = %v, want %v", len(stmt.Select.Columns), len(tt.columns))
			}
		})
	}
}

func TestParseInsert(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name     string
		query    string
		expected StatementType
		table    string
		columns  []string
	}{
		{
			name:     "simple insert",
			query:    "INSERT INTO users (name, email) VALUES ('John', 'john@example.com')",
			expected: StatementTypeInsert,
			table:    "users",
			columns:  []string{"name", "email"},
		},
		{
			name:     "insert without columns",
			query:    "INSERT INTO users VALUES ('John', 'john@example.com')",
			expected: StatementTypeInsert,
			table:    "users",
			columns:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := parser.Parse(tt.query)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if stmt.Type != tt.expected {
				t.Errorf("Type = %v, want %v", stmt.Type, tt.expected)
			}

			if stmt.Insert == nil {
				t.Fatal("Insert should not be nil")
			}

			if stmt.Insert.Table != tt.table {
				t.Errorf("Table = %v, want %v", stmt.Insert.Table, tt.table)
			}

			if tt.columns != nil && len(stmt.Insert.Columns) != len(tt.columns) {
				t.Errorf("Columns length = %v, want %v", len(stmt.Insert.Columns), len(tt.columns))
			}
		})
	}
}

func TestParseUpdate(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name     string
		query    string
		expected StatementType
		table    string
	}{
		{
			name:     "simple update",
			query:    "UPDATE users SET name = 'Jane' WHERE id = 1",
			expected: StatementTypeUpdate,
			table:    "users",
		},
		{
			name:     "update without where",
			query:    "UPDATE users SET name = 'Jane'",
			expected: StatementTypeUpdate,
			table:    "users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := parser.Parse(tt.query)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if stmt.Type != tt.expected {
				t.Errorf("Type = %v, want %v", stmt.Type, tt.expected)
			}

			if stmt.Update == nil {
				t.Fatal("Update should not be nil")
			}

			if stmt.Update.Table != tt.table {
				t.Errorf("Table = %v, want %v", stmt.Update.Table, tt.table)
			}
		})
	}
}

func TestParseDelete(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name     string
		query    string
		expected StatementType
		table    string
	}{
		{
			name:     "simple delete",
			query:    "DELETE FROM users WHERE id = 1",
			expected: StatementTypeDelete,
			table:    "users",
		},
		{
			name:     "delete without where",
			query:    "DELETE FROM users",
			expected: StatementTypeDelete,
			table:    "users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := parser.Parse(tt.query)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if stmt.Type != tt.expected {
				t.Errorf("Type = %v, want %v", stmt.Type, tt.expected)
			}

			if stmt.Delete == nil {
				t.Fatal("Delete should not be nil")
			}

			if stmt.Delete.Table != tt.table {
				t.Errorf("Table = %v, want %v", stmt.Delete.Table, tt.table)
			}
		})
	}
}

func TestParseCreateTable(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name        string
		query       string
		expected    StatementType
		tableName   string
		columnCount int
	}{
		{
			name:        "simple create table",
			query:       "CREATE TABLE users (id INT, name VARCHAR(50))",
			expected:    StatementTypeCreateTable,
			tableName:   "users",
			columnCount: 2,
		},
		{
			name:        "create table with constraints",
			query:       "CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(50) NOT NULL)",
			expected:    StatementTypeCreateTable,
			tableName:   "users",
			columnCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := parser.Parse(tt.query)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if stmt.Type != tt.expected {
				t.Errorf("Type = %v, want %v", stmt.Type, tt.expected)
			}

			if stmt.CreateTable == nil {
				t.Fatal("CreateTable should not be nil")
			}

			if stmt.CreateTable.TableName != tt.tableName {
				t.Errorf("TableName = %v, want %v", stmt.CreateTable.TableName, tt.tableName)
			}

			if len(stmt.CreateTable.Columns) != tt.columnCount {
				t.Errorf("Columns length = %v, want %v", len(stmt.CreateTable.Columns), tt.columnCount)
			}
		})
	}
}

func TestParseDropTable(t *testing.T) {
	parser := NewParser()

	query := "DROP TABLE users"
	stmt, err := parser.Parse(query)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if stmt.Type != StatementTypeDropTable {
		t.Errorf("Type = %v, want %v", stmt.Type, StatementTypeDropTable)
	}

	if stmt.DropTable == nil {
		t.Fatal("DropTable should not be nil")
	}

	if stmt.DropTable.TableName != "users" {
		t.Errorf("TableName = %v, want users", stmt.DropTable.TableName)
	}
}

func TestParseStatementType(t *testing.T) {
	tests := []struct {
		query    string
		expected StatementType
	}{
		{"SELECT * FROM users", StatementTypeSelect},
		{"INSERT INTO users VALUES (1)", StatementTypeInsert},
		{"UPDATE users SET name = 'John'", StatementTypeUpdate},
		{"DELETE FROM users", StatementTypeDelete},
		{"CREATE TABLE test (id INT)", StatementTypeCreateTable},
		{"DROP TABLE test", StatementTypeDropTable},
		{"UNKNOWN QUERY", StatementTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			result := ParseStatementType(tt.query)
			if result != tt.expected {
				t.Errorf("ParseStatementType(%q) = %v, want %v", tt.query, result, tt.expected)
			}
		})
	}
}

func TestExtractTableName(t *testing.T) {
	tests := []struct {
		query    string
		expected string
	}{
		{"SELECT * FROM users", "users"},
		{"INSERT INTO users VALUES (1)", "users"},
		{"UPDATE users SET name = 'John'", "users"},
		{"DELETE FROM users", "users"},
		{"CREATE TABLE users (id INT)", "users"},
		{"DROP TABLE users", "users"},
		{"SELECT * FROM dbo.users", "users"},
		{"SELECT * FROM [users]", "users"},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			result := ExtractTableName(tt.query)
			if result != tt.expected {
				t.Errorf("ExtractTableName(%q) = %v, want %v", tt.query, result, tt.expected)
			}
		})
	}
}

func TestHasWhereClause(t *testing.T) {
	tests := []struct {
		query    string
		expected bool
	}{
		{"SELECT * FROM users WHERE id = 1", true},
		{"SELECT * FROM users", false},
		{"UPDATE users SET name = 'John' WHERE id = 1", true},
		{"DELETE FROM users", false},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			result := HasWhereClause(tt.query)
			if result != tt.expected {
				t.Errorf("HasWhereClause(%q) = %v, want %v", tt.query, result, tt.expected)
			}
		})
	}
}

func TestStripComments(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"SELECT * FROM users -- comment", "SELECT * FROM users"},
		{"SELECT * FROM users /* multi-line comment */", "SELECT * FROM users"},
		{"SELECT * FROM -- comment\nusers", "SELECT * FROM users"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := StripComments(tt.input)
			if result != tt.expected {
				t.Errorf("StripComments(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
