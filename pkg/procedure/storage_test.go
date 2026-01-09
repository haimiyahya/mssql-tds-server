package procedure

import (
	"os"
	"testing"

	"github.com/factory/mssql-tds-server/pkg/sqlite"
)

func setupTestDB(t *testing.T) (*sqlite.Database, *Storage, func()) {
	// Create temporary database
	dbPath := "/tmp/test_procedures.db"
	os.Remove(dbPath) // Clean up any existing file

	db, err := sqlite.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	err = db.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	storage, err := NewStorage(db)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Cleanup function
	cleanup := func() {
		db.Close()
		os.Remove(dbPath)
	}

	return db, storage, cleanup
}

func TestStorage_Create(t *testing.T) {
	db, storage, cleanup := setupTestDB(t)
	defer cleanup()

	proc := &Procedure{
		Name: "TEST_PROC",
		Body: "SELECT * FROM users WHERE id = @id",
		Parameters: []Parameter{
			{Name: "id", Type: "INT"},
		},
	}

	err := storage.Create(proc)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	if proc.ID == 0 {
		t.Error("Create() did not set ID")
	}
}

func TestStorage_Get(t *testing.T) {
	db, storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a procedure
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

	// Retrieve the procedure
	retrieved, err := storage.Get("GET_USER")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}

	if retrieved.Name != proc.Name {
		t.Errorf("Get() Name = %v, want %v", retrieved.Name, proc.Name)
	}

	if retrieved.Body != proc.Body {
		t.Errorf("Get() Body = %v, want %v", retrieved.Body, proc.Body)
	}

	if len(retrieved.Parameters) != len(proc.Parameters) {
		t.Errorf("Get() Parameters count = %v, want %v", len(retrieved.Parameters), len(proc.Parameters))
	}
}

func TestStorage_Get_NotFound(t *testing.T) {
	db, storage, cleanup := setupTestDB(t)
	defer cleanup()

	_, err := storage.Get("NONEXISTENT")
	if err == nil {
		t.Error("Get() expected error for non-existent procedure")
	}
}

func TestStorage_List(t *testing.T) {
	db, storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Create multiple procedures
	procedures := []*Procedure{
		{Name: "PROC1", Body: "SELECT 1", Parameters: []Parameter{}},
		{Name: "PROC2", Body: "SELECT 2", Parameters: []Parameter{}},
		{Name: "PROC3", Body: "SELECT 3", Parameters: []Parameter{}},
	}

	for _, proc := range procedures {
		err := storage.Create(proc)
		if err != nil {
			t.Fatalf("Failed to create procedure: %v", err)
		}
	}

	// List all procedures
	list, err := storage.List()
	if err != nil {
		t.Errorf("List() error = %v", err)
	}

	if len(list) != len(procedures) {
		t.Errorf("List() count = %v, want %v", len(list), len(procedures))
	}
}

func TestStorage_Drop(t *testing.T) {
	db, storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a procedure
	proc := &Procedure{
		Name: "TO_DROP",
		Body: "SELECT 1",
		Parameters: []Parameter{},
	}

	err := storage.Create(proc)
	if err != nil {
		t.Fatalf("Failed to create procedure: %v", err)
	}

	// Drop the procedure
	err = storage.Drop("TO_DROP")
	if err != nil {
		t.Errorf("Drop() error = %v", err)
	}

	// Verify it's gone
	_, err = storage.Get("TO_DROP")
	if err == nil {
		t.Error("Drop() procedure still exists after being dropped")
	}
}

func TestStorage_Drop_NotFound(t *testing.T) {
	db, storage, cleanup := setupTestDB(t)
	defer cleanup()

	err := storage.Drop("NONEXISTENT")
	if err == nil {
		t.Error("Drop() expected error for non-existent procedure")
	}
}

func TestStorage_Exists(t *testing.T) {
	db, storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Test non-existent procedure
	exists, err := storage.Exists("NONEXISTENT")
	if err != nil {
		t.Errorf("Exists() error = %v", err)
	}
	if exists {
		t.Error("Exists() returned true for non-existent procedure")
	}

	// Create a procedure
	proc := &Procedure{
		Name: "CHECK_EXISTS",
		Body: "SELECT 1",
		Parameters: []Parameter{},
	}

	err = storage.Create(proc)
	if err != nil {
		t.Fatalf("Failed to create procedure: %v", err)
	}

	// Test existing procedure
	exists, err = storage.Exists("CHECK_EXISTS")
	if err != nil {
		t.Errorf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Exists() returned false for existing procedure")
	}
}

func TestStorage_Duplicate(t *testing.T) {
	db, storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a procedure
	proc := &Procedure{
		Name: "DUPLICATE",
		Body: "SELECT 1",
		Parameters: []Parameter{},
	}

	err := storage.Create(proc)
	if err != nil {
		t.Fatalf("Failed to create procedure: %v", err)
	}

	// Try to create duplicate
	err = storage.Create(proc)
	if err == nil {
		t.Error("Create() expected error for duplicate procedure name")
	}
}
