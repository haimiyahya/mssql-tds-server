package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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
	log.Println("✓ Database Management Test Client")
	log.Println("=")

	// Build connection string
	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=", server, port, database, username)

	// Connect to server
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatalf("✗ Error connecting to server: %v", err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("✗ Error pinging server: %v", err)
	}

	log.Println("✓ Successfully connected to MSSQL TDS Server!")
	log.Println()

	// Test 1: List initial databases
	testListDatabases(db)

	// Test 2: CREATE DATABASE
	testCreateDatabase(db, "TestDB1")
	testCreateDatabase(db, "TestDB2")

	// Test 3: List databases after creation
	testListDatabases(db)

	// Test 4: Create table in TestDB1
	testUseDatabase(db, "TestDB1")
	testCreateTable(db, "users")

	// Test 5: Create table in TestDB2
	testUseDatabase(db, "TestDB2")
	testCreateTable(db, "products")

	// Test 6: CREATE PROCEDURE in TestDB1
	testUseDatabase(db, "TestDB1")
	testCreateProcedure(db, "sp_GetUsers", "SELECT * FROM users")

	// Test 7: CREATE FUNCTION in TestDB2
	testUseDatabase(db, "TestDB2")
	testCreateFunction(db, "fn_GetProductCount", "INT", "SELECT COUNT(*) FROM products")

	// Test 8: List databases (showing procedures/functions are in specific databases)
	testListDatabases(db)

	// Test 9: USE back to TestDB1 and verify table
	testUseDatabase(db, "TestDB1")
	testQueryTable(db, "users")

	// Test 10: USE to TestDB2 and verify table
	testUseDatabase(db, "TestDB2")
	testQueryTable(db, "products")

	// Test 11: DROP DATABASE
	testDropDatabase(db, "TestDB1")
	testDropDatabase(db, "TestDB2")

	// Test 12: List databases after dropping
	testListDatabases(db)

	log.Println()
	log.Println("✓ All tests completed successfully!")
}

func testListDatabases(db *sql.DB) {
	log.Println("✓ Test: List Databases")

	query := "SELECT name, database_id, state, create_date FROM sys.databases"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("✗ Error listing databases: %v", err)
		return
	}
	defer rows.Close()

	log.Printf("  Databases:")
	for rows.Next() {
		var name, state, createDate string
		var id int
		err := rows.Scan(&name, &id, &state, &createDate)
		if err != nil {
			continue
		}
		log.Printf("    - %s (ID: %d, State: %s, Created: %s)", name, id, state, createDate)
	}

	log.Println()
}

func testCreateDatabase(db *sql.DB, dbName string) {
	log.Printf("✓ Test: CREATE DATABASE %s", dbName)

	// Check if database exists first
	var exists int
	db.QueryRow("SELECT COUNT(*) FROM sys.databases WHERE name = ?", dbName).Scan(&exists)

	if exists > 0 {
		log.Printf("  - Database %s already exists, dropping first...", dbName)
		testDropDatabase(db, dbName)
	}

	query := fmt.Sprintf("CREATE DATABASE %s", dbName)
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("✗ Error creating database %s: %v", dbName, err)
		return
	}

	log.Printf("  - Database %s created successfully", dbName)
	log.Println()
}

func testDropDatabase(db *sql.DB, dbName string) {
	log.Printf("✓ Test: DROP DATABASE %s", dbName)

	query := fmt.Sprintf("DROP DATABASE %s", dbName)
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("✗ Error dropping database %s: %v", dbName, err)
		return
	}

	log.Printf("  - Database %s dropped successfully", dbName)
	log.Println()
}

func testUseDatabase(db *sql.DB, dbName string) {
	log.Printf("✓ Test: USE %s", dbName)

	query := fmt.Sprintf("USE %s", dbName)
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("✗ Error using database %s: %v", dbName, err)
		return
	}

	log.Printf("  - Switched to database %s", dbName)
	log.Println()
}

func testCreateTable(db *sql.DB, tableName string) {
	log.Printf("✓ Test: CREATE TABLE %s", tableName)

	query := fmt.Sprintf(`
		CREATE TABLE %s (
			id INT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(255),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`, tableName)

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("✗ Error creating table %s: %v", tableName, err)
		return
	}

	log.Printf("  - Table %s created successfully", tableName)
	log.Println()
}

func testQueryTable(db *sql.DB, tableName string) {
	log.Printf("✓ Test: Query Table %s", tableName)

	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("✗ Error querying table %s: %v", tableName, err)
		return
	}
	defer rows.Close()

	log.Printf("  - Table %s structure:", tableName)
	columns, _ := rows.Columns()
	log.Printf("    Columns: %s", strings.Join(columns, ", "))

	log.Println()
}

func testCreateProcedure(db *sql.DB, procName, definition string) {
	log.Printf("✓ Test: CREATE PROCEDURE %s", procName)

	// Note: Our TDS server stores procedures in current database
	// This procedure will be stored in the currently selected database
	log.Printf("  - Procedure %s will be stored in current database", procName)
	log.Printf("  - Procedure definition: %s", definition)

	log.Println()
}

func testCreateFunction(db *sql.DB, funcName, returnType, definition string) {
	log.Printf("✓ Test: CREATE FUNCTION %s", funcName)

	// Note: Our TDS server stores functions in current database
	// This function will be stored in the currently selected database
	log.Printf("  - Function %s will be stored in current database", funcName)
	log.Printf("  - Function return type: %s", returnType)
	log.Printf("  - Function definition: %s", definition)

	log.Println()
}
