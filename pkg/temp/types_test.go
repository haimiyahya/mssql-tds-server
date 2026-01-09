package temp

import (
	"strings"
	"testing"
)

func TestIsTempTable(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "Local temp table",
			input: "#temp",
			want:  true,
		},
		{
			name:  "Global temp table",
			input: "##temp",
			want:  true,
		},
		{
			name:  "Regular table",
			input: "users",
			want:  false,
		},
		{
			name:  "Temp table with underscore",
			input: "#temp_table",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTempTable(tt.input)
			if result != tt.want {
				t.Errorf("IsTempTable() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestNormalizeTableName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Remove single #",
			input: "#temp",
			want:  "temp",
		},
		{
			name:  "Remove double ##",
			input: "##temp",
			want:  "temp",
		},
		{
			name:  "Regular table",
			input: "users",
			want:  "users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeTableName(tt.input)
			if result != tt.want {
				t.Errorf("NormalizeTableName() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestParseCreateTable(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantName  string
		wantCols  int
		wantErr   bool
	}{
		{
			name:     "Simple temp table",
			sql:      "CREATE TABLE #temp (id INT, name VARCHAR(50))",
			wantName: "temp",
			wantCols: 2,
			wantErr:  false,
		},
		{
			name:     "Temp table with nullable column",
			sql:      "CREATE TABLE #temp (id INT, name VARCHAR(50) NULL)",
			wantName: "temp",
			wantCols: 2,
			wantErr:  false,
		},
		{
			name:     "Not a temp table",
			sql:      "CREATE TABLE users (id INT, name VARCHAR(50))",
			wantName: "",
			wantCols: 0,
			wantErr:  true,
		},
		{
			name:     "Invalid syntax",
			sql:      "CREATE TABLE #temp",
			wantName: "",
			wantCols: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, columns, err := ParseCreateTable(tt.sql)

			if tt.wantErr {
				if err == nil {
					t.Error("ParseCreateTable() expected error but got none")
					return
				}
			} else {
				if err != nil {
					t.Errorf("ParseCreateTable() unexpected error = %v", err)
					return
				}

				if name != tt.wantName {
					t.Errorf("ParseCreateTable() name = %v, want %v", name, tt.wantName)
				}

				if len(columns) != tt.wantCols {
					t.Errorf("ParseCreateTable() columns count = %d, want %d", len(columns), tt.wantCols)
				}
			}
		})
	}
}

func TestManager(t *testing.T) {
	manager := NewManager()

	// Test session creation
	sessionID := manager.CreateSession()
	if sessionID == "" {
		t.Error("CreateSession() returned empty string")
	}

	// Test table creation
	columns := []Column{
		{Name: "id", Type: "INT", Nullable: false},
		{Name: "name", Type: "VARCHAR", Nullable: true},
	}

	table, err := manager.CreateTable(sessionID, "temp", columns)
	if err != nil {
		t.Errorf("CreateTable() error = %v", err)
		return
	}

	if table.Name != "temp" {
		t.Errorf("CreateTable() name = %v, want temp", table.Name)
	}

	if len(table.Columns) != 2 {
		t.Errorf("CreateTable() columns count = %d, want 2", len(table.Columns))
	}

	// Test insert
	row := Row{"id": 1, "name": "Alice"}
	err = manager.Insert(sessionID, "temp", row)
	if err != nil {
		t.Errorf("Insert() error = %v", err)
	}

	// Test select
	rows, err := manager.Select(sessionID, "temp")
	if err != nil {
		t.Errorf("Select() error = %v", err)
	}

	if len(rows) != 1 {
		t.Errorf("Select() rows count = %d, want 1", len(rows))
	}

	// Test update
	updates := Row{"name": "Bob"}
	count, err := manager.Update(sessionID, "temp", func(r Row) bool {
		return r["id"] == 1
	}, updates)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	if count != 1 {
		t.Errorf("Update() count = %d, want 1", count)
	}

	// Verify update
	rows, _ = manager.Select(sessionID, "temp")
	if rows[0]["name"] != "Bob" {
		t.Errorf("Update() name = %v, want Bob", rows[0]["name"])
	}

	// Test delete
	count, err = manager.Delete(sessionID, "temp", func(r Row) bool {
		return r["id"] == 1
	})
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	if count != 1 {
		t.Errorf("Delete() count = %d, want 1", count)
	}

	// Verify delete
	rows, _ = manager.Select(sessionID, "temp")
	if len(rows) != 0 {
		t.Errorf("Delete() rows count = %d, want 0", len(rows))
	}

	// Test drop table
	err = manager.DropTable(sessionID, "temp")
	if err != nil {
		t.Errorf("DropTable() error = %v", err)
	}

	// Verify drop
	_, err = manager.GetTable(sessionID, "temp")
	if err == nil {
		t.Error("DropTable() expected error but got none")
	}

	// Test drop session
	manager.DropSession(sessionID)

	// Verify session dropped
	_, err = manager.GetTable(sessionID, "temp")
	if err == nil {
		t.Error("DropSession() expected error but got none")
	}
}

func TestDetectTempTableReference(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		wantTables int
	}{
		{
			name:       "Single temp table reference",
			sql:        "SELECT * FROM #temp",
			wantTables: 1,
		},
		{
			name:       "Multiple temp table references",
			sql:        "SELECT * FROM #temp JOIN #temp2 ON #temp.id = #temp2.id",
			wantTables: 2,
		},
		{
			name:       "Duplicate references",
			sql:        "SELECT * FROM #temp WHERE id IN (SELECT id FROM #temp)",
			wantTables: 1,
		},
		{
			name:       "No temp table references",
			sql:        "SELECT * FROM users",
			wantTables: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tables := DetectTempTableReference(tt.sql)
			if len(tables) != tt.wantTables {
				t.Errorf("DetectTempTableReference() count = %d, want %d", len(tables), tt.wantTables)
			}
		})
	}
}

func TestReplaceTempTableNames(t *testing.T) {
	sessionID := "session_1"

	tests := []struct {
		name     string
		sql      string
		contains string
	}{
		{
			name:     "Single temp table",
			sql:      "SELECT * FROM #temp",
			contains: "session_1_temp",
		},
		{
			name:     "Multiple temp tables",
			sql:      "SELECT * FROM #temp JOIN #temp2 ON #temp.id = #temp2.id",
			contains: "session_1_temp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceTempTableNames(tt.sql, sessionID)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("ReplaceTempTableNames() result = %v, want to contain %v", result, tt.contains)
			}
		})
	}
}
