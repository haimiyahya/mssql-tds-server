package procedure

import (
	"fmt"
	"testing"
)

func debugLog(t *testing.T, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	t.Log("DEBUG:", msg)
	fmt.Println("DEBUG:", msg)
}

func TestParseCreateProcedure(t *testing.T) {
	debugLog(t, "TestParseCreateProcedure: START")
	
	tests := []struct {
		name        string
		sql         string
		wantName    string
		wantBody    string
		wantParams  int
		wantErr     bool
		errContains string
	}{
		{
			name:       "Simple procedure without parameters",
			sql:        "CREATE PROCEDURE GetUser AS SELECT * FROM users",
			wantName:   "GETUSER",
			wantBody:   "SELECT * FROM users",
			wantParams: 0,
			wantErr:    false,
		},
		{
			name:       "Procedure with single parameter",
			sql:        "CREATE PROCEDURE GetUserById @id INT AS SELECT * FROM users WHERE id = @id",
			wantName:   "GETUSERBYID",
			wantBody:   "SELECT * FROM users WHERE id = @id",
			wantParams: 1,
			wantErr:    false,
		},
		{
			name:       "Procedure with multiple parameters",
			sql:        "CREATE PROCEDURE GetUserByDept @dept VARCHAR(50), @active BIT AS SELECT * FROM users WHERE department = @dept AND active = @active",
			wantName:   "GETUSERBYDEPT",
			wantBody:   "SELECT * FROM users WHERE department = @dept AND active = @active",
			wantParams: 2,
			wantErr:    false,
		},
		{
			name:       "Procedure with default value",
			sql:        "CREATE PROCEDURE GetUsers @limit INT DEFAULT 10 AS SELECT TOP @limit * FROM users",
			wantName:   "GETUSERS",
			wantBody:   "SELECT TOP @limit * FROM users",
			wantParams: 1,
			wantErr:    false,
		},
		{
			name:        "Not a CREATE PROCEDURE statement",
			sql:         "SELECT * FROM users",
			wantErr:     true,
			errContains: "not a CREATE PROCEDURE",
		},
		{
			name:        "Invalid syntax",
			sql:         "CREATE PROCEDURE GetUser AS",
			wantErr:     true,
			errContains: "invalid CREATE PROCEDURE syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugLog(t, "TestParseCreateProcedure/%s: START", tt.name)
			proc, err := ParseCreateProcedure(tt.sql)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseCreateProcedure() expected error but got none")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ParseCreateProcedure() error = %v, want error containing %s", err, tt.errContains)
				}
				debugLog(t, "TestParseCreateProcedure/%s: END (expected error)", tt.name)
				return
			}

			if err != nil {
				t.Errorf("ParseCreateProcedure() unexpected error = %v", err)
				return
			}

			if proc.Name != tt.wantName {
				t.Errorf("ParseCreateProcedure() Name = %v, want %v", proc.Name, tt.wantName)
			}

			if proc.Body != tt.wantBody {
				t.Errorf("ParseCreateProcedure() Body = %v, want %v", proc.Body, tt.wantBody)
			}

			if len(proc.Parameters) != tt.wantParams {
				t.Errorf("ParseCreateProcedure() Parameters count = %v, want %v", len(proc.Parameters), tt.wantParams)
			}
			debugLog(t, "TestParseCreateProcedure/%s: END", tt.name)
		})
	}
	debugLog(t, "TestParseCreateProcedure: END")
}

func TestParseParameters(t *testing.T) {
	debugLog(t, "TestParseParameters: START")
	
	tests := []struct {
		name      string
		paramsStr string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Empty parameters",
			paramsStr:  "",
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "Single parameter",
			paramsStr:  "@id INT",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "Multiple parameters",
			paramsStr:  "@id INT, @name VARCHAR(50), @active BIT",
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "Parameter with default",
			paramsStr:  "@limit INT DEFAULT 10",
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "Parameter with length",
			paramsStr:  "@name VARCHAR(100)",
			wantCount: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugLog(t, "TestParseParameters/%s: START", tt.name)
			params, err := parseParameters(tt.paramsStr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseParameters() expected error but got none")
				}
				debugLog(t, "TestParseParameters/%s: END (expected error)", tt.name)
				return
			}

			if err != nil {
				t.Errorf("parseParameters() unexpected error = %v", err)
				return
			}

			if len(params) != tt.wantCount {
				t.Errorf("parseParameters() count = %v, want %v", len(params), tt.wantCount)
			}
			debugLog(t, "TestParseParameters/%s: END", tt.name)
		})
	}
	debugLog(t, "TestParseParameters: END")
}

func TestParametersToJSON(t *testing.T) {
	debugLog(t, "TestParametersToJSON: START")
	
	params := []Parameter{
		{Name: "id", Type: "INT"},
		{Name: "name", Type: "VARCHAR", Length: 50},
		{Name: "active", Type: "BIT", HasDefault: true, Default: "1"},
	}

	json, err := ParametersToJSON(params)
	if err != nil {
		t.Fatalf("ParametersToJSON() error = %v", err)
	}

	if json == "" {
		t.Error("ParametersToJSON() returned empty string")
	}

	// Try to unmarshal back
	params2, err := ParametersFromJSON(json)
	if err != nil {
		t.Fatalf("ParametersFromJSON() error = %v", err)
	}

	if len(params2) != len(params) {
		t.Errorf("Parameters count mismatch: got %d, want %d", len(params2), len(params))
	}
	debugLog(t, "TestParametersToJSON: END")
}

func TestParametersFromJSON(t *testing.T) {
	debugLog(t, "TestParametersFromJSON: START")
	
	tests := []struct {
		name      string
		json      string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Empty JSON array",
			json:      "[]",
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "Empty string",
			json:      "",
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "Valid JSON",
			json:      `[{"name":"id","type":"INT"},{"name":"name","type":"VARCHAR"}]`,
			wantCount: 2,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugLog(t, "TestParametersFromJSON/%s: START", tt.name)
			params, err := ParametersFromJSON(tt.json)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParametersFromJSON() expected error but got none")
				}
				debugLog(t, "TestParametersFromJSON/%s: END (expected error)", tt.name)
				return
			}

			if err != nil {
				t.Errorf("ParametersFromJSON() unexpected error = %v", err)
				return
			}

			if len(params) != tt.wantCount {
				t.Errorf("ParametersFromJSON() count = %v, want %v", len(params), tt.wantCount)
			}
			debugLog(t, "TestParametersFromJSON/%s: END", tt.name)
		})
	}
	debugLog(t, "TestParametersFromJSON: END")
}

func TestSplitByCommaOutsideParentheses(t *testing.T) {
	debugLog(t, "TestSplitByCommaOutsideParentheses: START")
	
	tests := []struct {
		name     string
		input    string
		wantLen  int
	}{
		{
			name:    "Simple comma separated",
			input:   "@id INT, @name VARCHAR",
			wantLen: 2,
		},
		{
			name:    "With parentheses",
			input:   "@name VARCHAR(50), @id INT",
			wantLen: 2,
		},
		{
			name:    "Complex with nested parentheses",
			input:   "@param1 VARCHAR(50), @param2 INT, @param3 FLOAT(10,2)",
			wantLen: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugLog(t, "TestSplitByCommaOutsideParentheses/%s: START", tt.name)
			result := splitByCommaOutsideParentheses(tt.input)
			if len(result) != tt.wantLen {
				t.Errorf("splitByCommaOutsideParentheses() = %v, want len %v", result, tt.wantLen)
			}
			debugLog(t, "TestSplitByCommaOutsideParentheses/%s: END", tt.name)
		})
	}
	debugLog(t, "TestSplitByCommaOutsideParentheses: END")
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || len(s) > len(substr) && containsString(s[1:], substr)
}
