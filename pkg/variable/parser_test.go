package variable

import (
	"testing"
)

func TestParseStatement(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantType StatementType
	}{
		{
			name:     "DECLARE statement",
			sql:      "DECLARE @id INT",
			wantType: StatementDeclare,
		},
		{
			name:     "SET statement",
			sql:      "SET @id = 123",
			wantType: StatementSet,
		},
		{
			name:     "SELECT assignment",
			sql:      "SELECT @id = id FROM users",
			wantType: StatementSelectAssignment,
		},
		{
			name:     "SELECT query",
			sql:      "SELECT * FROM users",
			wantType: StatementQuery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseStatement(tt.sql)
			if result != tt.wantType {
				t.Errorf("ParseStatement() = %v, want %v", result, tt.wantType)
			}
		})
	}
}

func TestParseProcedureBody(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		wantStatements  int
		wantErr        bool
	}{
		{
			name:          "Single statement",
			body:          "SELECT * FROM users",
			wantStatements: 1,
			wantErr:       false,
		},
		{
			name:          "Multiple statements",
			body:          "DECLARE @id INT; SET @id = 1; SELECT * FROM users WHERE id = @id",
			wantStatements: 3,
			wantErr:       false,
		},
		{
			name:          "Statement with string containing semicolon",
			body:          "SELECT * FROM users WHERE name = 'O\\'Brien'; SELECT * FROM users",
			wantStatements: 2,
			wantErr:       false,
		},
		{
			name:          "Statements with parentheses",
			body:          "SELECT COUNT(*) FROM users; SELECT * FROM users WHERE id IN (1, 2, 3)",
			wantStatements: 2,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statements, err := ParseProcedureBody(tt.body)

			if tt.wantErr {
				if err == nil {
					t.Error("ParseProcedureBody() expected error but got none")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("ParseProcedureBody() unexpected error = %v", err)
				return
			}

			if len(statements) != tt.wantStatements {
				t.Errorf("ParseProcedureBody() count = %d, want %d", len(statements), tt.wantStatements)
			}
		})
	}
}

func TestFindVariableReferences(t *testing.T) {
	tests := []struct {
		name           string
		sql            string
		wantReferences int
	}{
		{
			name:           "Single variable",
			sql:            "SELECT * FROM users WHERE id = @id",
			wantReferences: 1,
		},
		{
			name:           "Multiple variables",
			sql:            "SELECT * FROM users WHERE id = @id AND name = @name",
			wantReferences: 2,
		},
		{
			name:           "Duplicate variable references",
			sql:            "SELECT * FROM users WHERE id = @id AND parent_id = @id",
			wantReferences: 1, // Should deduplicate
		},
		{
			name:           "No variables",
			sql:            "SELECT * FROM users",
			wantReferences: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			references := FindVariableReferences(tt.sql)
			if len(references) != tt.wantReferences {
				t.Errorf("FindVariableReferences() count = %d, want %d", len(references), tt.wantReferences)
			}
		})
	}
}

func TestReplaceVariables(t *testing.T) {
	ctx := NewContext()

	// Declare variables
	ctx.Declare("@id", TypeInt, 0)
	ctx.Declare("@name", TypeVarchar, 50)

	// Set values
	ctx.Set("@id", 1)
	ctx.Set("@name", "Alice")

	tests := []struct {
		name    string
		sql     string
		want    string
		wantErr bool
	}{
		{
			name:    "Replace single variable",
			sql:     "SELECT * FROM users WHERE id = @id",
			want:    "SELECT * FROM users WHERE id = 1",
			wantErr: false,
		},
		{
			name:    "Replace multiple variables",
			sql:     "SELECT * FROM users WHERE id = @id AND name = @name",
			want:    "SELECT * FROM users WHERE id = 1 AND name = 'Alice'",
			wantErr: false,
		},
		{
			name:    "Undeclared variable",
			sql:     "SELECT * FROM users WHERE id = @undeclared",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ReplaceVariables(tt.sql, ctx)

			if tt.wantErr {
				if err == nil {
					t.Error("ReplaceVariables() expected error but got none")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("ReplaceVariables() unexpected error = %v", err)
				return
			}

			if result != tt.want {
				t.Errorf("ReplaceVariables() = %v, want %v", result, tt.want)
			}
		})
	}
}
