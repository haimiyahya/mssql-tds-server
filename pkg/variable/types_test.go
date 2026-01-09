package variable

import (
	"testing"
)

func TestContext_Declare(t *testing.T) {
	ctx := NewContext()

	// Test single variable declaration
	variable, err := ctx.Declare("@id", TypeInt, 0)
	if err != nil {
		t.Errorf("Declare() error = %v", err)
		return
	}

	if variable.Name != "ID" {
		t.Errorf("Declare() Name = %v, want ID", variable.Name)
	}

	if variable.Type != TypeInt {
		t.Errorf("Declare() Type = %v, want %v", variable.Type, TypeInt)
	}

	if !variable.IsNull {
		t.Error("Declare() IsNull should be true for new variable")
	}
}

func TestContext_Declare_Duplicate(t *testing.T) {
	ctx := NewContext()

	// Declare first variable
	_, err := ctx.Declare("@id", TypeInt, 0)
	if err != nil {
		t.Fatalf("Declare() error = %v", err)
	}

	// Try to declare duplicate
	_, err = ctx.Declare("@id", TypeInt, 0)
	if err == nil {
		t.Error("Declare() expected error for duplicate variable")
	}
}

func TestContext_Get(t *testing.T) {
	ctx := NewContext()

	// Declare variable
	_, err := ctx.Declare("@name", TypeVarchar, 50)
	if err != nil {
		t.Fatalf("Declare() error = %v", err)
	}

	// Get variable
	variable, exists := ctx.Get("@name")
	if !exists {
		t.Error("Get() variable should exist")
		return
	}

	if variable.Name != "NAME" {
		t.Errorf("Get() Name = %v, want NAME", variable.Name)
	}
}

func TestContext_Get_NotFound(t *testing.T) {
	ctx := NewContext()

	_, exists := ctx.Get("@nonexistent")
	if exists {
		t.Error("Get() variable should not exist")
	}
}

func TestContext_Set(t *testing.T) {
	ctx := NewContext()

	// Declare variable
	_, err := ctx.Declare("@id", TypeInt, 0)
	if err != nil {
		t.Fatalf("Declare() error = %v", err)
	}

	// Set value
	err = ctx.Set("@id", 123)
	if err != nil {
		t.Errorf("Set() error = %v", err)
	}

	// Get and verify
	variable, _ := ctx.Get("@id")
	if variable.Value != 123 {
		t.Errorf("Set() Value = %v, want 123", variable.Value)
	}

	if variable.IsNull {
		t.Error("Set() IsNull should be false after setting value")
	}
}

func TestContext_Set_Null(t *testing.T) {
	ctx := NewContext()

	// Declare variable
	_, err := ctx.Declare("@id", TypeInt, 0)
	if err != nil {
		t.Fatalf("Declare() error = %v", err)
	}

	// Set value
	err = ctx.Set("@id", 456)
	if err != nil {
		t.Errorf("Set() error = %v", err)
	}

	// Set to NULL
	err = ctx.Set("@id", nil)
	if err != nil {
		t.Errorf("Set() error = %v", err)
	}

	// Get and verify
	variable, _ := ctx.Get("@id")
	if variable.Value != nil {
		t.Errorf("Set() Value = %v, want nil", variable.Value)
	}

	if !variable.IsNull {
		t.Error("Set() IsNull should be true after setting to NULL")
	}
}

func TestContext_Set_Undeclared(t *testing.T) {
	ctx := NewContext()

	err := ctx.Set("@undeclared", 123)
	if err == nil {
		t.Error("Set() expected error for undeclared variable")
	}
}

func TestParseDeclaration(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantName  string
		wantType  Type
		wantLen   int
		wantErr   bool
	}{
		{
			name:     "INT variable",
			sql:      "DECLARE @id INT",
			wantName: "id",
			wantType: TypeInt,
			wantLen:  0,
			wantErr:  false,
		},
		{
			name:     "VARCHAR variable",
			sql:      "DECLARE @name VARCHAR(50)",
			wantName: "name",
			wantType: TypeVarchar,
			wantLen:  50,
			wantErr:  false,
		},
		{
			name:     "BIT variable",
			sql:      "DECLARE @active BIT",
			wantName: "active",
			wantType: TypeBit,
			wantLen:  0,
			wantErr:  false,
		},
		{
			name:    "Invalid syntax",
			sql:     "DECLARE",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variable, err := ParseDeclaration(tt.sql)

			if tt.wantErr {
				if err == nil {
					t.Error("ParseDeclaration() expected error but got none")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("ParseDeclaration() unexpected error = %v", err)
				return
			}

			if variable.Name != tt.wantName {
				t.Errorf("ParseDeclaration() Name = %v, want %v", variable.Name, tt.wantName)
			}

			if variable.Type != tt.wantType {
				t.Errorf("ParseDeclaration() Type = %v, want %v", variable.Type, tt.wantType)
			}

			if variable.Length != tt.wantLen {
				t.Errorf("ParseDeclaration() Length = %v, want %v", variable.Length, tt.wantLen)
			}
		})
	}
}

func TestParseSetAssignment(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		wantName   string
		wantValue  interface{}
		wantErr    bool
	}{
		{
			name:      "Integer value",
			sql:       "SET @id = 123",
			wantName:  "id",
			wantValue: 123,
			wantErr:   false,
		},
		{
			name:      "String value",
			sql:       "SET @name = 'Alice'",
			wantName:  "name",
			wantValue: "Alice",
			wantErr:   false,
		},
		{
			name:      "NULL value",
			sql:       "SET @value = NULL",
			wantName:  "value",
			wantValue: nil,
			wantErr:   false,
		},
		{
			name:      "Float value",
			sql:       "SET @price = 19.99",
			wantName:  "price",
			wantValue: 19.99,
			wantErr:   false,
		},
		{
			name:    "Invalid syntax",
			sql:     "SET @id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, value, err := ParseSetAssignment(tt.sql)

			if tt.wantErr {
				if err == nil {
					t.Error("ParseSetAssignment() expected error but got none")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("ParseSetAssignment() unexpected error = %v", err)
				return
			}

			if name != tt.wantName {
				t.Errorf("ParseSetAssignment() Name = %v, want %v", name, tt.wantName)
			}

			// Compare values loosely due to type differences
			if tt.wantValue != nil && value != tt.wantValue {
				t.Errorf("ParseSetAssignment() Value = %v (%T), want %v (%T)", value, value, tt.wantValue, tt.wantValue)
			}
		})
	}
}

func TestParseSelectAssignment(t *testing.T) {
	tests := []struct {
		name          string
		sql           string
		wantName      string
		wantExpr      string
		wantErr       bool
	}{
		{
			name:     "Simple assignment",
			sql:      "SELECT @id = id FROM users",
			wantName: "id",
			wantExpr: "id FROM users",
			wantErr:  false,
		},
		{
			name:     "Expression assignment",
			sql:      "SELECT @count = COUNT(*) FROM users",
			wantName: "count",
			wantExpr: "COUNT(*) FROM users",
			wantErr:  false,
		},
		{
			name:    "Not an assignment",
			sql:     "SELECT * FROM users",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, expr, err := ParseSelectAssignment(tt.sql)

			if tt.wantErr {
				if err == nil {
					t.Error("ParseSelectAssignment() expected error but got none")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("ParseSelectAssignment() unexpected error = %v", err)
				return
			}

			if name != tt.wantName {
				t.Errorf("ParseSelectAssignment() Name = %v, want %v", name, tt.wantName)
			}

			if expr != tt.wantExpr {
				t.Errorf("ParseSelectAssignment() Expression = %v, want %v", expr, tt.wantExpr)
			}
		})
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		varType  Type
		want     string
	}{
		{
			name:    "Integer",
			value:   123,
			varType: TypeInt,
			want:    "123",
		},
		{
			name:    "String",
			value:   "Alice",
			varType: TypeVarchar,
			want:    "'Alice'",
		},
		{
			name:    "String with quotes",
			value:   "O'Brien",
			varType: TypeVarchar,
			want:    "'O''Brien'",
		},
		{
			name:    "NULL",
			value:   nil,
			varType: TypeInt,
			want:    "NULL",
		},
		{
			name:    "Boolean true",
			value:   true,
			varType: TypeBit,
			want:    "1",
		},
		{
			name:    "Boolean false",
			value:   false,
			varType: TypeBit,
			want:    "0",
		},
		{
			name:    "Float",
			value:   19.99,
			varType: TypeFloat,
			want:    "19.990000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatValue(tt.value, tt.varType)
			if result != tt.want {
				t.Errorf("FormatValue() = %v, want %v", result, tt.want)
			}
		})
	}
}
