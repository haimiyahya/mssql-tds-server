package controlflow

import (
	"testing"
)

func TestParseCondition(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		wantErr  bool
	}{
		{
			name:    "Simple equality",
			expr:    "@id = 1",
			wantErr: false,
		},
		{
			name:    "String comparison",
			expr:    "@name = 'Alice'",
			wantErr: false,
		},
		{
			name:    "Greater than",
			expr:    "@price > 100",
			wantErr: false,
		},
		{
			name:    "Less than or equal",
			expr:    "@stock <= 50",
			wantErr: false,
		},
		{
			name:    "Not equal",
			expr:    "@status <> 'INACTIVE'",
			wantErr: false,
		},
		{
			name:    "AND condition",
			expr:    "@id = 1 AND @active = 1",
			wantErr: false,
		},
		{
			name:    "OR condition",
			expr:    "@status = 'ACTIVE' OR @status = 'PENDING'",
			wantErr: false,
		},
		{
			name:    "Complex condition",
			expr:    "(@id = 1 OR @id = 2) AND @active = 1",
			wantErr: false,
		},
		{
			name:    "No operator",
			expr:    "@id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCondition(tt.expr)

			if tt.wantErr {
				if err == nil {
					t.Error("ParseCondition() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ParseCondition() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name      string
		expr      string
		variables map[string]interface{}
		want      bool
		wantErr   bool
	}{
		{
			name:      "Integer equality true",
			expr:      "@id = 1",
			variables: map[string]interface{}{"@id": 1},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "Integer equality false",
			expr:      "@id = 2",
			variables: map[string]interface{}{"@id": 1},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "String equality true",
			expr:      "@name = 'Alice'",
			variables: map[string]interface{}{"@name": "Alice"},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "Greater than true",
			expr:      "@price > 100",
			variables: map[string]interface{}{"@price": 150.5},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "Less than true",
			expr:      "@stock < 50",
			variables: map[string]interface{}{"@stock": 25},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "Not equal true",
			expr:      "@status <> 'INACTIVE'",
			variables: map[string]interface{}{"@status": "ACTIVE"},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "AND condition true",
			expr:      "@id = 1 AND @active = 1",
			variables: map[string]interface{}{"@id": 1, "@active": 1},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "AND condition false",
			expr:      "@id = 1 AND @active = 0",
			variables: map[string]interface{}{"@id": 1, "@active": 0},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "OR condition true",
			expr:      "@status = 'ACTIVE' OR @status = 'PENDING'",
			variables: map[string]interface{}{"@status": "PENDING"},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "OR condition false",
			expr:      "@status = 'ACTIVE' OR @status = 'PENDING'",
			variables: map[string]interface{}{"@status": "INACTIVE"},
			want:      false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			condition, err := ParseCondition(tt.expr)
			if err != nil {
				t.Fatalf("ParseCondition() error = %v", err)
			}

			result, err := Evaluate(condition, tt.variables)

			if tt.wantErr {
				if err == nil {
					t.Error("Evaluate() expected error but got none")
					return
				}
			} else {
				if err != nil {
					t.Errorf("Evaluate() unexpected error = %v", err)
					return
				}

				if result != tt.want {
					t.Errorf("Evaluate() = %v, want %v", result, tt.want)
				}
			}
		})
	}
}

func TestParseIFStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{
			name:    "Simple IF",
			sql:     "IF @id = 1 THEN SELECT 'Match' END",
			wantErr: false,
		},
		{
			name:    "IF with ELSE",
			sql:     "IF @status = 'ACTIVE' THEN SELECT 'Active' ELSE SELECT 'Inactive' END",
			wantErr: false,
		},
		{
			name:    "IF with complex condition",
			sql:     "IF @id = 1 AND @active = 1 THEN SELECT 'Match' END",
			wantErr: false,
		},
		{
			name:    "Missing THEN",
			sql:     "IF @id = 1 SELECT 'Match' END",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := ParseIFStatement(tt.sql)

			if tt.wantErr {
				if err == nil {
					t.Error("ParseIFStatement() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ParseIFStatement() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestCompareValues(t *testing.T) {
	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		want     int // -1 if a < b, 0 if a == b, 1 if a > b
	}{
		{
			name: "Integers less",
			a:    1,
			b:    2,
			want:  -1,
		},
		{
			name: "Integers equal",
			a:    5,
			b:    5,
			want:  0,
		},
		{
			name: "Integers greater",
			a:    10,
			b:    5,
			want:  1,
		},
		{
			name: "Strings less",
			a:    "Alice",
			b:    "Bob",
			want:  -1,
		},
		{
			name: "Strings equal",
			a:    "Test",
			b:    "Test",
			want:  0,
		},
		{
			name: "Strings greater",
			a:    "Charlie",
			b:    "Alice",
			want:  1,
		},
		{
			name: "Floats less",
			a:    1.5,
			b:    2.5,
			want:  -1,
		},
		{
			name: "Floats greater",
			a:    3.14,
			b:    2.71,
			want:  1,
		},
		{
			name: "Nil vs value",
			a:    nil,
			b:    1,
			want:  -1,
		},
		{
			name: "Value vs nil",
			a:    "test",
			b:    nil,
			want:  1,
		},
		{
			name: "Both nil",
			a:    nil,
			b:    nil,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareValues(tt.a, tt.b)
			if result != tt.want {
				t.Errorf("compareValues() = %d, want %d", result, tt.want)
			}
		})
	}
}
