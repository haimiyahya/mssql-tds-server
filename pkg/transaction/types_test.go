package transaction

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestParseStatement(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		want    TransactionType
	}{
		{
			name: "BEGIN TRANSACTION",
			sql:  "BEGIN TRANSACTION",
			want: TransactionBegin,
		},
		{
			name: "BEGIN TRAN",
			sql:  "BEGIN TRAN",
			want: TransactionBegin,
		},
		{
			name: "START TRANSACTION",
			sql:  "START TRANSACTION",
			want: TransactionBegin,
		},
		{
			name: "COMMIT",
			sql:  "COMMIT",
			want: TransactionCommit,
		},
		{
			name: "COMMIT TRANSACTION",
			sql:  "COMMIT TRANSACTION",
			want: TransactionCommit,
		},
		{
			name: "ROLLBACK",
			sql:  "ROLLBACK",
			want: TransactionRollback,
		},
		{
			name: "ROLLBACK TRANSACTION",
			sql:  "ROLLBACK TRANSACTION",
			want: TransactionRollback,
		},
		{
			name: "SELECT statement",
			sql:  "SELECT * FROM users",
			want: TransactionUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseStatement(tt.sql)
			if result != tt.want {
				t.Errorf("ParseStatement() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestIsTransactionStatement(t *testing.T) {
	tests := []struct {
		name  string
		sql   string
		want  bool
	}{
		{
			name: "Transaction statement",
			sql:  "BEGIN TRANSACTION",
			want:  true,
		},
		{
			name: "Non-transaction statement",
			sql:  "SELECT * FROM users",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTransactionStatement(tt.sql)
			if result != tt.want {
				t.Errorf("IsTransactionStatement() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestDetectTransactionUsage(t *testing.T) {
	tests := []struct {
		name  string
		sql   string
		want  bool
	}{
		{
			name: "Procedure with transactions",
			sql:  "BEGIN TRANSACTION; SELECT * FROM users; COMMIT",
			want:  true,
		},
		{
			name: "Procedure without transactions",
			sql:  "SELECT * FROM users",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectTransactionUsage(tt.sql)
			if result != tt.want {
				t.Errorf("DetectTransactionUsage() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestContext(t *testing.T) {
	// Create in-memory database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create test table
	_, err = db.Exec("CREATE TABLE test (id INTEGER, value TEXT)")
	if err != nil {
		t.Fatalf("Error creating table: %v", err)
	}

	ctx := NewContext()

	// Test 1: Begin transaction
	tx, err := ctx.Begin(db)
	if err != nil {
		t.Errorf("Begin() error = %v", err)
		return
	}

	if !ctx.IsActive() {
		t.Error("Begin() did not set active transaction")
	}

	if ctx.GetLevel() != 1 {
		t.Errorf("GetLevel() = %d, want 1", ctx.GetLevel())
	}

	// Insert data in transaction
	_, err = tx.Exec("INSERT INTO test (id, value) VALUES (1, 'test')")
	if err != nil {
		t.Errorf("Error inserting data: %v", err)
	}

	// Test 2: Commit transaction
	err = ctx.Commit()
	if err != nil {
		t.Errorf("Commit() error = %v", err)
		return
	}

	if ctx.IsActive() {
		t.Error("Commit() did not clear active transaction")
	}

	if ctx.GetLevel() != 0 {
		t.Errorf("GetLevel() = %d, want 0", ctx.GetLevel())
	}

	// Verify data was committed
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM test").Scan(&count)
	if err != nil {
		t.Errorf("Error querying count: %v", err)
	}

	if count != 1 {
		t.Errorf("Count = %d, want 1", count)
	}

	// Test 3: Rollback transaction
	tx, err = ctx.Begin(db)
	if err != nil {
		t.Errorf("Begin() error = %v", err)
		return
	}

	// Insert data in transaction
	_, err = tx.Exec("INSERT INTO test (id, value) VALUES (2, 'test2')")
	if err != nil {
		t.Errorf("Error inserting data: %v", err)
	}

	// Rollback
	err = ctx.Rollback()
	if err != nil {
		t.Errorf("Rollback() error = %v", err)
		return
	}

	// Verify data was not committed
	err = db.QueryRow("SELECT COUNT(*) FROM test").Scan(&count)
	if err != nil {
		t.Errorf("Error querying count: %v", err)
	}

	if count != 1 {
		t.Errorf("Count = %d, want 1 (rollback should have discarded insert)", count)
	}
}

func TestGetCurrentTx(t *testing.T) {
	ctx := NewContext()

	// Test 1: No active transaction
	tx := ctx.GetCurrentTx()
	if tx != nil {
		t.Error("GetCurrentTx() should return nil when no transaction is active")
	}

	// Note: We can't test with actual transaction without a real database
	// The Context struct is tested indirectly through other tests
}
