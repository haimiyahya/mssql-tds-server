package procedure

import (
	"os"
	"testing"

	"github.com/factory/mssql-tds-server/pkg/sqlite"
)


func setupTestDB(t *testing.T) (*sqlite.Database, *Storage, func()) {
	debugLog(t, "setupTestDB: START")
	
	// Create temporary database
	dbPath := "/tmp/test_procedures.db"
	debugLog(t, "setupTestDB: Removing db file: %s", dbPath)
	os.Remove(dbPath) // Clean up any existing file

	debugLog(t, "setupTestDB: Creating database")
	db, err := sqlite.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	debugLog(t, "setupTestDB: Database created")

	debugLog(t, "setupTestDB: Initializing database")
	err = db.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	debugLog(t, "setupTestDB: Database initialized")

	debugLog(t, "setupTestDB: Creating storage")
	storage, err := NewStorage(db)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	debugLog(t, "setupTestDB: Storage created")

	// Cleanup function
	cleanup := func() {
		debugLog(t, "setupTestDB: Cleanup called")
		db.Close()
		os.Remove(dbPath)
	}

	debugLog(t, "setupTestDB: END")
	return db, storage, cleanup
}

func TestStorage_Create(t *testing.T) {
	debugLog(t, "TestStorage_Create: START")
	
	_, storage, cleanup := setupTestDB(t)
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
	debugLog(t, "TestStorage_Create: END")
}

func TestStorage_Get(t *testing.T) {
	debugLog(t, "TestStorage_Get: START")
	
	_, storage, cleanup := setupTestDB(t)
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
	debugLog(t, "TestStorage_Get: Procedure created")

	// Retrieve procedure
	debugLog(t, "TestStorage_Get: Retrieving procedure")
	retrieved, err := storage.Get("GET_USER")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	debugLog(t, "TestStorage_Get: Procedure retrieved")

	if retrieved.Name != proc.Name {
		t.Errorf("Get() Name = %v, want %v", retrieved.Name, proc.Name)
	}

	if retrieved.Body != proc.Body {
		t.Errorf("Get() Body = %v, want %v", retrieved.Body, proc.Body)
	}

	if len(retrieved.Parameters) != len(proc.Parameters) {
		t.Errorf("Get() Parameters count = %v, want %v", len(retrieved.Parameters), len(proc.Parameters))
	}
	debugLog(t, "TestStorage_Get: END")
}

func TestStorage_Get_NotFound(t *testing.T) {
	debugLog(t, "TestStorage_Get_NotFound: START")
	
	_, storage, cleanup := setupTestDB(t)
	defer cleanup()

	debugLog(t, "TestStorage_Get_NotFound: Retrieving non-existent procedure")
	_, err := storage.Get("NONEXISTENT")
	if err == nil {
		t.Error("Get() expected error for non-existent procedure")
	}
	debugLog(t, "TestStorage_Get_NotFound: END")
}

func TestStorage_List(t *testing.T) {
	debugLog(t, "TestStorage_List: START")
	
	_, storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Create multiple procedures
	procedures := []*Procedure{
		{Name: "PROC1", Body: "SELECT 1", Parameters: []Parameter{}},
		{Name: "PROC2", Body: "SELECT 2", Parameters: []Parameter{}},
		{Name: "PROC3", Body: "SELECT 3", Parameters: []Parameter{}},
	}

	for _, proc := range procedures {
		debugLog(t, "TestStorage_List: Creating procedure: %s", proc.Name)
		err := storage.Create(proc)
		if err != nil {
			t.Fatalf("Failed to create procedure: %v", err)
		}
	}

	// List all procedures
	debugLog(t, "TestStorage_List: Listing procedures")
	list, err := storage.List()
	if err != nil {
		t.Errorf("List() error = %v", err)
	}

	if len(list) != len(procedures) {
		t.Errorf("List() count = %v, want %v", len(list), len(procedures))
	}
	debugLog(t, "TestStorage_List: END")
}

func TestStorage_Drop(t *testing.T) {
	debugLog(t, "TestStorage_Drop: START")
	
	_, storage, cleanup := setupTestDB(t)
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

	// Drop procedure
	debugLog(t, "TestStorage_Drop: Dropping procedure")
	err = storage.Drop("TO_DROP")
	if err != nil {
		t.Errorf("Drop() error = %v", err)
	}
	debugLog(t, "TestStorage_Drop: Procedure dropped")

	// Verify it's gone
	debugLog(t, "TestStorage_Drop: Verifying procedure is gone")
	_, err = storage.Get("TO_DROP")
	if err == nil {
		t.Error("Drop() procedure still exists after being dropped")
	}
	debugLog(t, "TestStorage_Drop: END")
}

func TestStorage_Drop_NotFound(t *testing.T) {
	debugLog(t, "TestStorage_Drop_NotFound: START")
	
	_, storage, cleanup := setupTestDB(t)
	defer cleanup()

	err := storage.Drop("NONEXISTENT")
	if err == nil {
		t.Error("Drop() expected error for non-existent procedure")
	}
	debugLog(t, "TestStorage_Drop_NotFound: END")
}

func TestStorage_Exists(t *testing.T) {
	debugLog(t, "TestStorage_Exists: START")
	
	_, storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Test non-existent procedure
	debugLog(t, "TestStorage_Exists: Checking non-existent procedure")
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
	debugLog(t, "TestStorage_Exists: Checking existing procedure")
	exists, err = storage.Exists("CHECK_EXISTS")
	if err != nil {
		t.Errorf("Exists() error = %v", err)
	}
	if !exists {
		t.Error("Exists() returned false for existing procedure")
	}
	debugLog(t, "TestStorage_Exists: END")
}

func TestStorage_Duplicate(t *testing.T) {
	debugLog(t, "TestStorage_Duplicate: START")
	
	_, storage, cleanup := setupTestDB(t)
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
	debugLog(t, "TestStorage_Duplicate: Creating duplicate procedure")
	err = storage.Create(proc)
	if err == nil {
		t.Error("Create() expected error for duplicate procedure name")
	}
	debugLog(t, "TestStorage_Duplicate: END")
}
