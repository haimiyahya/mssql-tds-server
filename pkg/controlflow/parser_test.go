package controlflow

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
			name:     "IF statement",
			sql:      "IF @id = 1 THEN SELECT 'Match' END",
			wantType: StatementIF,
		},
		{
			name:     "DECLARE statement",
			sql:      "DECLARE @id INT",
			wantType: StatementDeclare,
		},
		{
			name:     "SET statement",
			sql:      "SET @id = 1",
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

func TestSplitStatements(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		wantCount     int
		wantErr       bool
	}{
		{
			name:      "Single statement",
			body:      "SELECT * FROM users",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "Multiple statements with semicolons",
			body:      "DECLARE @id INT; SET @id = 1; SELECT @id",
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "IF block",
			body:      "IF @id = 1 THEN SELECT 'Match' END",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "IF with ELSE",
			body:      "IF @id = 1 THEN SELECT 'Match' ELSE SELECT 'No Match' END",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "IF block followed by semicolon statement",
			body:      "IF @id = 1 THEN SELECT 'Match' END; SELECT * FROM users",
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "Statements with string containing semicolon",
			body:      "SELECT * FROM users WHERE name = 'O\\'Brien'; SELECT * FROM users",
			wantCount: 2,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statements, err := SplitStatements(tt.body)

			if tt.wantErr {
				if err == nil {
					t.Error("SplitStatements() expected error but got none")
					return
				}
			} else {
				if err != nil {
					t.Errorf("SplitStatements() unexpected error = %v", err)
					return
				}

				if len(statements) != tt.wantCount {
					t.Errorf("SplitStatements() count = %d, want %d", len(statements), tt.wantCount)
				}
			}
		})
	}
}

func TestParseIFBlock(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantErr   bool
		hasElse   bool
	}{
		{
			name:    "Simple IF",
			sql:     "IF @id = 1 THEN SELECT 'Match' END",
			wantErr: false,
			hasElse: false,
		},
		{
			name:    "IF with ELSE",
			sql:     "IF @id = 1 THEN SELECT 'Match' ELSE SELECT 'No Match' END",
			wantErr: false,
			hasElse: true,
		},
		{
			name:    "IF with complex condition",
			sql:     "IF @id = 1 AND @active = 1 THEN SELECT 'Match' END",
			wantErr: false,
			hasElse: false,
		},
		{
			name:    "Missing THEN",
			sql:     "IF @id = 1 SELECT 'Match' END",
			wantErr: true,
		},
		{
			name:    "Missing END",
			sql:     "IF @id = 1 THEN SELECT 'Match'",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block, elseInfo, err := ParseIFBlock(tt.sql)

			if tt.wantErr {
				if err == nil {
					t.Error("ParseIFBlock() expected error but got none")
					return
				}
			} else {
				if err != nil {
					t.Errorf("ParseIFBlock() unexpected error = %v", err)
					return
				}

				if block == nil {
					t.Error("ParseIFBlock() returned nil block")
					return
				}

				hasElse := elseInfo != nil
				if hasElse != tt.hasElse {
					t.Errorf("ParseIFBlock() hasElse = %v, want %v", hasElse, tt.hasElse)
				}
			}
		})
	}
}

func TestFindOutermostEND(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantPos  int
		wantErr  bool
	}{
		{
			name:    "Simple END",
			sql:     "IF @id = 1 THEN SELECT 'Match' END",
			wantPos: len("IF @id = 1 THEN SELECT 'Match'"),
			wantErr: false,
		},
		{
			name:    "END in string",
			sql:     "IF @id = 1 THEN SELECT 'END' END",
			wantPos: len("IF @id = 1 THEN SELECT 'END'"),
			wantErr: false,
		},
		{
			name:    "END after parentheses",
			sql:     "IF (@id = 1) THEN SELECT 'Match' END",
			wantPos: len("IF (@id = 1) THEN SELECT 'Match'"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := findOutermostEND(tt.sql)
			if pos != tt.wantPos {
				t.Errorf("findOutermostEND() = %d, want %d", pos, tt.wantPos)
			}
		})
	}
}

func TestFindELSEInBody(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantPos int
	}{
		{
			name:    "No ELSE",
			body:    "SELECT 'Match'",
			wantPos: -1,
		},
		{
			name:    "Simple ELSE",
			body:    "SELECT 'Match' ELSE SELECT 'No Match'",
			wantPos: len("SELECT 'Match'"),
		},
		{
			name:    "ELSE in string",
			body:    "SELECT 'ELSE' ELSE SELECT 'No Match'",
			wantPos: len("SELECT 'ELSE'"),
		},
		{
			name:    "ELSE after parentheses",
			body:    "SELECT 'Match' ELSE (SELECT 'No Match')",
			wantPos: len("SELECT 'Match'"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := findELSEInBody(tt.body)
			if pos != tt.wantPos {
				t.Errorf("findELSEInBody() = %d, want %d", pos, tt.wantPos)
			}
		})
	}
}
