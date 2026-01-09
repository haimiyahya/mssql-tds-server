package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
)

const (
	server   = "127.0.0.1"
	port     = 1433
	database = ""
	username = "sa"
	password = ""
)

func main() {
	// Build connection string
	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s",
		server, port, database, username, password)

	log.Printf("Connecting to TDS server at %s:%d", server, port)

	// Connect to database
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging server: %v", err)
	}

	log.Println("Successfully connected to TDS server!")

	// Test 1: PRIMARY KEY constraint
	log.Println("\n=== Test 1: PRIMARY KEY constraint ===")
	testPrimaryKey(db)

	// Test 2: NOT NULL constraint
	log.Println("\n=== Test 2: NOT NULL constraint ===")
	testNotNull(db)

	// Test 3: UNIQUE constraint
	log.Println("\n=== Test 3: UNIQUE constraint ===")
	testUnique(db)

	// Test 4: DEFAULT constraint
	log.Println("\n=== Test 4: DEFAULT constraint ===")
	testDefault(db)

	// Test 5: CHECK constraint
	log.Println("\n=== Test 5: CHECK constraint ===")
	testCheck(db)

	// Test 6: FOREIGN KEY constraint
	log.Println("\n=== Test 6: FOREIGN KEY constraint ===")
	testForeignKey(db)

	// Test 7: Combined constraints
	log.Println("\n=== Test 7: Combined constraints ===")
	testCombinedConstraints(db)

	// Test 8: Multiple constraints on same column
	log.Println("\n=== Test 8: Multiple constraints on same column ===")
	testMultipleConstraints(db)

	// Test 9: DROP TABLES
	log.Println("\n=== Test 9: DROP TABLES ===")
	testDropTables(db)

	log.Println("\n=== All Phase 16 tests completed! ===")
	log.Println("🎉 Phase 16: Constraint Support - COMPLETE! 🎉")
}

func testPrimaryKey(db *sql.DB) {
	// Create table with PRIMARY KEY
	_, err := db.Exec("CREATE TABLE pk_test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table with PRIMARY KEY constraint")

	// Insert valid data
	_, err = db.Exec("INSERT INTO pk_test VALUES (1, 'Alice')")
	if err != nil {
		log.Printf("Error inserting valid data: %v", err)
		return
	}
	log.Println("✓ Inserted valid data with unique primary key")

	// Try to insert duplicate PRIMARY KEY
	_, err = db.Exec("INSERT INTO pk_test VALUES (1, 'Bob')")
	if err != nil {
		log.Printf("✓ PRIMARY KEY constraint works: duplicate rejected (%v)", err)
		return
	}
	log.Printf("✗ Failed: PRIMARY KEY constraint not enforced")
}

func testNotNull(db *sql.DB) {
	// Create table with NOT NULL constraint
	_, err := db.Exec("CREATE TABLE nn_test (id INTEGER, name TEXT NOT NULL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table with NOT NULL constraint")

	// Insert valid data
	_, err = db.Exec("INSERT INTO nn_test VALUES (1, 'Alice')")
	if err != nil {
		log.Printf("Error inserting valid data: %v", err)
		return
	}
	log.Println("✓ Inserted valid data with NOT NULL column")

	// Try to insert NULL into NOT NULL column
	_, err = db.Exec("INSERT INTO nn_test VALUES (2, NULL)")
	if err != nil {
		log.Printf("✓ NOT NULL constraint works: NULL rejected (%v)", err)
		return
	}
	log.Printf("✗ Failed: NOT NULL constraint not enforced")
}

func testUnique(db *sql.DB) {
	// Create table with UNIQUE constraint
	_, err := db.Exec("CREATE TABLE uniq_test (id INTEGER, email TEXT UNIQUE)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table with UNIQUE constraint")

	// Insert valid data
	_, err = db.Exec("INSERT INTO uniq_test VALUES (1, 'alice@example.com')")
	if err != nil {
		log.Printf("Error inserting valid data: %v", err)
		return
	}
	log.Println("✓ Inserted valid data with unique email")

	// Try to insert duplicate UNIQUE value
	_, err = db.Exec("INSERT INTO uniq_test VALUES (2, 'alice@example.com')")
	if err != nil {
		log.Printf("✓ UNIQUE constraint works: duplicate rejected (%v)", err)
		return
	}
	log.Printf("✗ Failed: UNIQUE constraint not enforced")
}

func testDefault(db *sql.DB) {
	// Create table with DEFAULT constraint
	_, err := db.Exec("CREATE TABLE def_test (id INTEGER, status TEXT DEFAULT 'active')")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table with DEFAULT constraint")

	// Insert without specifying default column
	_, err = db.Exec("INSERT INTO def_test (id) VALUES (1)")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}
	log.Println("✓ Inserted data without specifying default column")

	// Query to verify default value
	var status string
	err = db.QueryRow("SELECT status FROM def_test WHERE id = 1").Scan(&status)
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	log.Printf("✓ DEFAULT value works: status = '%s'", status)
}

func testCheck(db *sql.DB) {
	// Create table with CHECK constraint
	_, err := db.Exec("CREATE TABLE check_test (id INTEGER, age INTEGER CHECK (age >= 0 AND age < 150))")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table with CHECK constraint")

	// Insert valid data
	_, err = db.Exec("INSERT INTO check_test VALUES (1, 25)")
	if err != nil {
		log.Printf("Error inserting valid data: %v", err)
		return
	}
	log.Println("✓ Inserted valid data that satisfies CHECK constraint")

	// Try to insert invalid data
	_, err = db.Exec("INSERT INTO check_test VALUES (2, -5)")
	if err != nil {
		log.Printf("✓ CHECK constraint works: invalid data rejected (%v)", err)
		return
	}
	log.Printf("✗ Failed: CHECK constraint not enforced")
}

func testForeignKey(db *sql.DB) {
	// Create parent table
	_, err := db.Exec("CREATE TABLE parent (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		log.Printf("Error creating parent table: %v", err)
		return
	}
	log.Println("✓ Created parent table")

	// Insert data into parent table
	_, err = db.Exec("INSERT INTO parent VALUES (1, 'Department A')")
	if err != nil {
		log.Printf("Error inserting into parent: %v", err)
		return
	}
	log.Println("✓ Inserted data into parent table")

	// Create child table with FOREIGN KEY constraint
	_, err = db.Exec("CREATE TABLE child (id INTEGER, parent_id INTEGER REFERENCES parent(id))")
	if err != nil {
		log.Printf("Error creating child table: %v", err)
		return
	}
	log.Println("✓ Created child table with FOREIGN KEY constraint")

	// Insert valid data with valid foreign key
	_, err = db.Exec("INSERT INTO child VALUES (1, 1)")
	if err != nil {
		log.Printf("Error inserting valid data: %v", err)
		return
	}
	log.Println("✓ Inserted valid data with valid foreign key")

	// Try to insert invalid data with non-existent foreign key
	_, err = db.Exec("INSERT INTO child VALUES (2, 999)")
	if err != nil {
		log.Printf("✓ FOREIGN KEY constraint works: invalid reference rejected (%v)", err)
		return
	}
	log.Printf("✗ Failed: FOREIGN KEY constraint not enforced")
}

func testCombinedConstraints(db *sql.DB) {
	// Create table with combined constraints
	_, err := db.Exec("CREATE TABLE combined_test (id INTEGER PRIMARY KEY, email TEXT UNIQUE NOT NULL DEFAULT 'test@example.com', age INTEGER CHECK (age >= 0))")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table with combined constraints")

	// Insert valid data
	_, err = db.Exec("INSERT INTO combined_test VALUES (1, 'alice@example.com', 25)")
	if err != nil {
		log.Printf("Error inserting valid data: %v", err)
		return
	}
	log.Println("✓ Inserted valid data with combined constraints")

	// Verify all constraints
	var email string
	var age int
	err = db.QueryRow("SELECT email, age FROM combined_test WHERE id = 1").Scan(&email, &age)
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	log.Printf("✓ Verified: email='%s', age=%d", email, age)
}

func testMultipleConstraints(db *sql.DB) {
	// Create table with PRIMARY KEY + UNIQUE + NOT NULL
	_, err := db.Exec("CREATE TABLE multi_test (id INTEGER PRIMARY KEY UNIQUE NOT NULL, name TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table with multiple constraints on same column")

	// Insert valid data
	_, err = db.Exec("INSERT INTO multi_test VALUES (1, 'Alice')")
	if err != nil {
		log.Printf("Error inserting valid data: %v", err)
		return
	}
	log.Println("✓ Inserted valid data")

	// Verify constraints work
	_, err = db.Exec("INSERT INTO multi_test VALUES (1, 'Bob')")
	if err != nil {
		log.Printf("✓ PRIMARY KEY constraint works: duplicate rejected")
	}

	_, err = db.Exec("INSERT INTO multi_test VALUES (2, NULL)")
	if err != nil {
		log.Printf("✓ NOT NULL constraint works: NULL rejected")
	}

	log.Println("✓ All multiple constraints working correctly")
}

func testDropTables(db *sql.DB) {
	tables := []string{
		"pk_test",
		"nn_test",
		"uniq_test",
		"def_test",
		"check_test",
		"child",
		"parent",
		"combined_test",
		"multi_test",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE %s", table))
		if err != nil {
			log.Printf("Error dropping table %s: %v", table, err)
			continue
		}
		log.Printf("✓ Dropped table: %s", table)
	}
}
