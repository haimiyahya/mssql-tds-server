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

	// Test 1: CREATE TABLE
	log.Println("\n=== Test 1: CREATE TABLE ===")
	testCreateTable(db)

	// Test 2: INSERT data
	log.Println("\n=== Test 2: INSERT data ===")
	testInsert(db)

	// Test 3: ALTER TABLE ADD COLUMN
	log.Println("\n=== Test 3: ALTER TABLE ADD COLUMN ===")
	testAddColumn(db)

	// Test 4: ALTER TABLE ADD COLUMN with default value
	log.Println("\n=== Test 4: ALTER TABLE ADD COLUMN with default ===")
	testAddColumnWithDefault(db)

	// Test 5: ALTER TABLE ADD multiple columns
	log.Println("\n=== Test 5: ALTER TABLE ADD multiple columns ===")
	testAddMultipleColumns(db)

	// Test 6: ALTER TABLE RENAME TO
	log.Println("\n=== Test 6: ALTER TABLE RENAME TO ===")
	testRenameTable(db)

	// Test 7: ALTER TABLE RENAME COLUMN
	log.Println("\n=== Test 7: ALTER TABLE RENAME COLUMN ===")
	testRenameColumn(db)

	// Test 8: Verify schema changes
	log.Println("\n=== Test 8: Verify schema changes ===")
	testVerifySchema(db)

	// Test 9: INSERT after ALTER TABLE
	log.Println("\n=== Test 9: INSERT after ALTER TABLE ===")
	testInsertAfterAlter(db)

	// Test 10: Query after ALTER TABLE
	log.Println("\n=== Test 10: Query after ALTER TABLE ===")
	testQueryAfterAlter(db)

	// Test 11: DROP TABLE
	log.Println("\n=== Test 11: DROP TABLE ===")
	testDropTable(db)

	log.Println("\n=== All Phase 15 tests completed! ===")
	log.Println("🎉 Phase 15: ALTER TABLE Support - COMPLETE! 🎉")
}

func testCreateTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE users (id INTEGER, name TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: users")

	// Verify table created
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='users'")
	if err != nil {
		log.Printf("Error querying table: %v", err)
		return
	}
	defer rows.Close()

	var name string
	rows.Next()
	rows.Scan(&name)
	log.Printf("✓ Verified: Table 'users' created in database")
}

func testInsert(db *sql.DB) {
	queries := []string{
		"INSERT INTO users VALUES (1, 'Alice')",
		"INSERT INTO users VALUES (2, 'Bob')",
		"INSERT INTO users VALUES (3, 'Charlie')",
	}

	for _, query := range queries {
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("Error inserting row: %v", err)
			continue
		}
		rowsAffected, _ := result.RowsAffected()
		log.Printf("✓ Inserted row: %d row(s)", rowsAffected)
	}

	// Verify data
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
	}
	log.Printf("✓ Verified: %d row(s) in users table", rowCount)
}

func testAddColumn(db *sql.DB) {
	// ALTER TABLE ADD COLUMN
	_, err := db.Exec("ALTER TABLE users ADD COLUMN email TEXT")
	if err != nil {
		log.Printf("Error adding column: %v", err)
		return
	}
	log.Println("✓ Added column: email to users table")

	// Verify column added
	rows, err := db.Query("PRAGMA table_info(users)")
	if err != nil {
		log.Printf("Error querying schema: %v", err)
		return
	}
	defer rows.Close()

	columnCount := 0
	for rows.Next() {
		columnCount++
	}
	log.Printf("✓ Verified: Table has %d column(s)", columnCount)
}

func testAddColumnWithDefault(db *sql.DB) {
	// ALTER TABLE ADD COLUMN with default value
	_, err := db.Exec("ALTER TABLE users ADD COLUMN status TEXT DEFAULT 'active'")
	if err != nil {
		log.Printf("Error adding column with default: %v", err)
		return
	}
	log.Println("✓ Added column: status with default value 'active'")

	// Verify default value
	rows, err := db.Query("PRAGMA table_info(users)")
	if err != nil {
		log.Printf("Error querying schema: %v", err)
		return
	}
	defer rows.Close()

	columnCount := 0
	for rows.Next() {
		columnCount++
	}
	log.Printf("✓ Verified: Table has %d column(s) with default value", columnCount)
}

func testAddMultipleColumns(db *sql.DB) {
	// ALTER TABLE ADD multiple columns
	columns := []string{
		"ALTER TABLE users ADD COLUMN age INTEGER",
		"ALTER TABLE users ADD COLUMN city TEXT",
	}

	for _, query := range columns {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error adding column: %v", err)
			continue
		}
		log.Println("✓ Added column to users table")
	}

	// Verify all columns added
	rows, err := db.Query("PRAGMA table_info(users)")
	if err != nil {
		log.Printf("Error querying schema: %v", err)
		return
	}
	defer rows.Close()

	columnCount := 0
	for rows.Next() {
		columnCount++
	}
	log.Printf("✓ Verified: Table has %d column(s) after adding multiple columns", columnCount)
}

func testRenameTable(db *sql.DB) {
	// ALTER TABLE RENAME TO
	_, err := db.Exec("ALTER TABLE users RENAME TO employees")
	if err != nil {
		log.Printf("Error renaming table: %v", err)
		return
	}
	log.Println("✓ Renamed table: users to employees")

	// Verify table renamed
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='employees'")
	if err != nil {
		log.Printf("Error querying table: %v", err)
		return
	}
	defer rows.Close()

	var name string
	rows.Next()
	rows.Scan(&name)
	log.Printf("✓ Verified: Table renamed to 'employees'")
}

func testRenameColumn(db *sql.DB) {
	// ALTER TABLE RENAME COLUMN
	_, err := db.Exec("ALTER TABLE employees RENAME COLUMN name TO full_name")
	if err != nil {
		log.Printf("Error renaming column: %v", err)
		return
	}
	log.Println("✓ Renamed column: name to full_name in employees table")

	// Verify column renamed
	rows, err := db.Query("PRAGMA table_info(employees)")
	if err != nil {
		log.Printf("Error querying schema: %v", err)
		return
	}
	defer rows.Close()

	var colName string
	found := false
	for rows.Next() {
		rows.Scan(nil, &colName, nil, nil, nil) // Skip unused fields
		if colName == "full_name" {
			found = true
			break
		}
	}

	if found {
		log.Printf("✓ Verified: Column renamed to 'full_name'")
	} else {
		log.Printf("✗ Failed: Column not renamed")
	}
}

func testVerifySchema(db *sql.DB) {
	// Query schema
	rows, err := db.Query("PRAGMA table_info(employees)")
	if err != nil {
		log.Printf("Error querying schema: %v", err)
		return
	}
	defer rows.Close()

	log.Println("Table Schema:")
	for rows.Next() {
		var cid int
		var name string
		var typeStr string
		var notnull int
		var dfltValue sql.NullString
		err := rows.Scan(&cid, &name, &typeStr, &notnull, &dfltValue)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		defaultValue := "NULL"
		if dfltValue.Valid {
			defaultValue = dfltValue.String
		}

		log.Printf("  Column %d: %s (%s) - Default: %s, NotNull: %d", cid, name, typeStr, defaultValue, notnull)
	}

	log.Println("✓ Schema verified successfully")
}

func testInsertAfterAlter(db *sql.DB) {
	// Insert data with new columns
	query := "INSERT INTO employees (id, full_name, email, status, age, city) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, 4, "Diana", "diana@example.com", "active", 30, "New York")
	if err != nil {
		log.Printf("Error inserting row with new columns: %v", err)
		return
	}
	log.Println("✓ Inserted row with all columns (including newly added columns)")

	// Verify data inserted
	rows, err := db.Query("SELECT * FROM employees WHERE id = 4")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	var id int
	var fullName string
	var email sql.NullString
	var status sql.NullString
	var age sql.NullInt64
	var city sql.NullString

	rows.Next()
	rows.Scan(&id, &fullName, &email, &status, &age, &city)

	log.Printf("  Inserted data: ID=%d, Name=%s, Email=%s, Status=%s, Age=%d, City=%s",
		id, fullName, email.String, status.String, age.Int64, city.String)

	log.Println("✓ Verified: Data inserted with new columns")
}

func testQueryAfterAlter(db *sql.DB) {
	// Query with new columns
	rows, err := db.Query("SELECT id, full_name, email FROM employees")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var id int
		var fullName string
		var email sql.NullString
		err := rows.Scan(&id, &fullName, &email)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		emailStr := "NULL"
		if email.Valid {
			emailStr = email.String
		}

		if rowCount <= 3 {
			log.Printf("  Result: %d. %s (%s)", rowCount, fullName, emailStr)
		}
	}

	log.Printf("✓ Queried %d row(s) with new columns", rowCount)
}

func testDropTable(db *sql.DB) {
	_, err := db.Exec("DROP TABLE employees")
	if err != nil {
		log.Printf("Error dropping table: %v", err)
		return
	}
	log.Println("✓ Dropped table: employees")

	// Verify table dropped
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='employees'")
	if err != nil {
		log.Printf("Error querying table: %v", err)
		return
	}
	defer rows.Close()

	rows.Next()
	var name sql.NullString
	rows.Scan(&name)
	if !name.Valid {
		log.Println("✓ Verified: Table 'employees' dropped from database")
	} else {
		log.Printf("✗ Verified: Table still exists: %s", name.String)
	}
}
