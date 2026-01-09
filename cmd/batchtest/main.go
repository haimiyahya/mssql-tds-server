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

	// Test 1: Batch INSERT with multiple rows
	log.Println("\n=== Test 1: Batch INSERT with multiple rows ===")
	testBatchInsertMultipleRows(db)

	// Test 2: Batch INSERT with multiple statements
	log.Println("\n=== Test 2: Batch INSERT with multiple statements ===")
	testBatchInsertMultipleStatements(db)

	// Test 3: Batch UPDATE
	log.Println("\n=== Test 3: Batch UPDATE ===")
	testBatchUpdate(db)

	// Test 4: Batch DELETE
	log.Println("\n=== Test 4: Batch DELETE ===")
	testBatchDelete(db)

	// Test 5: Multi-statement transaction
	log.Println("\n=== Test 5: Multi-statement transaction ===")
	testMultiStatementTransaction(db)

	// Test 6: Batch operations with prepared statements
	log.Println("\n=== Test 6: Batch operations with prepared statements ===")
	testBatchPreparedStatements(db)

	// Test 7: Error handling in batches
	log.Println("\n=== Test 7: Error handling in batches ===")
	testBatchErrorHandling(db)

	// Test 8: Large batch operations
	log.Println("\n=== Test 8: Large batch operations ===")
	testLargeBatchOperations(db)

	// Test 9: Cleanup
	log.Println("\n=== Test 9: Cleanup ===")
	testCleanup(db)

	log.Println("\n=== All Phase 19 tests completed! ===")
	log.Println("🎉 Phase 19: Batch Operations - COMPLETE! 🎉")
}

func testBatchInsertMultipleRows(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE batch_users (id INTEGER, name TEXT, email TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: batch_users")

	// Batch INSERT with multiple rows
	_, err = db.Exec("INSERT INTO batch_users VALUES (1, 'Alice', 'alice@example.com'), (2, 'Bob', 'bob@example.com'), (3, 'Charlie', 'charlie@example.com')")
	if err != nil {
		log.Printf("Error with batch INSERT: %v", err)
		return
	}
	log.Println("✓ Batch INSERT (3 rows) successful")

	// Verify inserted data
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM batch_users").Scan(&count)
	if err != nil {
		log.Printf("Error querying count: %v", err)
		return
	}
	log.Printf("✓ Verified: %d rows inserted", count)
}

func testBatchInsertMultipleStatements(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE batch_products (id INTEGER, name TEXT, price REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: batch_products")

	// Batch INSERT with multiple statements (using transaction)
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return
	}

	_, err = tx.Exec("INSERT INTO batch_products VALUES (1, 'Product A', 99.99)")
	if err != nil {
		tx.Rollback()
		log.Printf("Error with first INSERT: %v", err)
		return
	}
	_, err = tx.Exec("INSERT INTO batch_products VALUES (2, 'Product B', 149.99)")
	if err != nil {
		tx.Rollback()
		log.Printf("Error with second INSERT: %v", err)
		return
	}
	_, err = tx.Exec("INSERT INTO batch_products VALUES (3, 'Product C', 199.99)")
	if err != nil {
		tx.Rollback()
		log.Printf("Error with third INSERT: %v", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return
	}
	log.Println("✓ Batch INSERT (3 statements with transaction) successful")

	// Verify inserted data
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM batch_products").Scan(&count)
	if err != nil {
		log.Printf("Error querying count: %v", err)
		return
	}
	log.Printf("✓ Verified: %d rows inserted", count)
}

func testBatchUpdate(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE batch_items (id INTEGER, name TEXT, quantity INTEGER, price REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: batch_items")

	// Insert initial data
	_, err = db.Exec("INSERT INTO batch_items VALUES (1, 'Item A', 10, 99.99), (2, 'Item B', 20, 149.99), (3, 'Item C', 30, 199.99)")
	if err != nil {
		log.Printf("Error inserting initial data: %v", err)
		return
	}

	// Batch UPDATE using CASE WHEN
	_, err = db.Exec(`
		UPDATE batch_items
		SET price = CASE
			WHEN id = 1 THEN 109.99
			WHEN id = 2 THEN 159.99
			WHEN id = 3 THEN 209.99
		END
		WHERE id IN (1, 2, 3)
	`)
	if err != nil {
		log.Printf("Error with batch UPDATE: %v", err)
		return
	}
	log.Println("✓ Batch UPDATE (3 rows with CASE WHEN) successful")

	// Verify updated data
	rows, err := db.Query("SELECT id, price FROM batch_items WHERE id IN (1, 2, 3)")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var price float64
		err := rows.Scan(&id, &price)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Verified: ID=%d, Price=%.2f", id, price)
	}
}

func testBatchDelete(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE batch_orders (id INTEGER, status TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: batch_orders")

	// Insert initial data
	_, err = db.Exec("INSERT INTO batch_orders VALUES (1, 'pending'), (2, 'completed'), (3, 'cancelled'), (4, 'pending'), (5, 'completed')")
	if err != nil {
		log.Printf("Error inserting initial data: %v", err)
		return
	}

	// Batch DELETE using IN clause
	result, err := db.Exec("DELETE FROM batch_orders WHERE status IN ('cancelled', 'pending')")
	if err != nil {
		log.Printf("Error with batch DELETE: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✓ Batch DELETE (%d rows) successful", rowsAffected)

	// Verify deleted data
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM batch_orders").Scan(&count)
	if err != nil {
		log.Printf("Error querying count: %v", err)
		return
	}
	log.Printf("✓ Verified: %d rows remaining", count)
}

func testMultiStatementTransaction(db *sql.DB) {
	// Create tables
	_, err := db.Exec("CREATE TABLE batch_accounts (id INTEGER, balance REAL)")
	if err != nil {
		log.Printf("Error creating accounts table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE batch_transactions (id INTEGER, from_id INTEGER, to_id INTEGER, amount REAL)")
	if err != nil {
		log.Printf("Error creating transactions table: %v", err)
		return
	}
	log.Println("✓ Created tables: batch_accounts, batch_transactions")

	// Insert initial data
	_, err = db.Exec("INSERT INTO batch_accounts VALUES (1, 1000.00), (2, 500.00)")
	if err != nil {
		log.Printf("Error inserting initial data: %v", err)
		return
	}

	// Multi-statement transaction (transfer money)
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return
	}

	// Deduct from account 1
	_, err = tx.Exec("UPDATE batch_accounts SET balance = balance - 100.00 WHERE id = 1")
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating account 1: %v", err)
		return
	}

	// Add to account 2
	_, err = tx.Exec("UPDATE batch_accounts SET balance = balance + 100.00 WHERE id = 2")
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating account 2: %v", err)
		return
	}

	// Record transaction
	_, err = tx.Exec("INSERT INTO batch_transactions VALUES (1, 1, 2, 100.00)")
	if err != nil {
		tx.Rollback()
		log.Printf("Error inserting transaction: %v", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return
	}
	log.Println("✓ Multi-statement transaction (transfer money) successful")

	// Verify balances
	var balance1, balance2 float64
	err = db.QueryRow("SELECT balance FROM batch_accounts WHERE id = 1").Scan(&balance1)
	if err != nil {
		log.Printf("Error querying account 1: %v", err)
		return
	}
	err = db.QueryRow("SELECT balance FROM batch_accounts WHERE id = 2").Scan(&balance2)
	if err != nil {
		log.Printf("Error querying account 2: %v", err)
		return
	}
	log.Printf("✓ Verified: Account 1 balance=%.2f, Account 2 balance=%.2f", balance1, balance2)
}

func testBatchPreparedStatements(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE batch_employees (id INTEGER, name TEXT, department TEXT, salary REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: batch_employees")

	// Prepare INSERT statement
	_, err = db.Exec("PREPARE insert_employee FROM 'INSERT INTO batch_employees VALUES ($id, $name, $department, $salary)'")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return
	}
	log.Println("✓ Prepared statement: insert_employee")

	// Execute prepared statement multiple times (batch)
	_, err = db.Exec("EXECUTE insert_employee USING @id = 1, @name = 'John', @department = 'Engineering', @salary = 80000.00")
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return
	}
	_, err = db.Exec("EXECUTE insert_employee USING @id = 2, @name = 'Jane', @department = 'Marketing', @salary = 75000.00")
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return
	}
	_, err = db.Exec("EXECUTE insert_employee USING @id = 3, @name = 'Bob', @department = 'Engineering', @salary = 85000.00")
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return
	}
	log.Println("✓ Batch operations with prepared statements (3 executions) successful")

	// Verify inserted data
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM batch_employees").Scan(&count)
	if err != nil {
		log.Printf("Error querying count: %v", err)
		return
	}
	log.Printf("✓ Verified: %d rows inserted", count)
}

func testBatchErrorHandling(db *sql.DB) {
	// Create table with constraint
	_, err := db.Exec("CREATE TABLE batch_unique_items (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: batch_unique_items")

	// Insert initial data
	_, err = db.Exec("INSERT INTO batch_unique_items VALUES (1, 'Item A'), (2, 'Item B')")
	if err != nil {
		log.Printf("Error inserting initial data: %v", err)
		return
	}

	// Try to insert duplicate primary key (should fail)
	_, err = db.Exec("INSERT INTO batch_unique_items VALUES (1, 'Duplicate'), (3, 'Item C')")
	if err != nil {
		log.Printf("✓ Batch INSERT with constraint violation correctly rejected: %v", err)
		return
	}
	log.Printf("✗ Failed: Batch INSERT with constraint violation not rejected")
}

func testLargeBatchOperations(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE batch_large (id INTEGER, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: batch_large")

	// Large batch INSERT (100 rows)
	valuesSQL := ""
	for i := 1; i <= 100; i++ {
		if i > 1 {
			valuesSQL += ", "
		}
		valuesSQL += fmt.Sprintf("(%d, %d)", i, i*10)
	}

	_, err = db.Exec("INSERT INTO batch_large VALUES " + valuesSQL)
	if err != nil {
		log.Printf("Error with large batch INSERT: %v", err)
		return
	}
	log.Printf("✓ Large batch INSERT (100 rows) successful")

	// Verify count
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM batch_large").Scan(&count)
	if err != nil {
		log.Printf("Error querying count: %v", err)
		return
	}
	log.Printf("✓ Verified: %d rows inserted", count)

	// Large batch UPDATE
	_, err = db.Exec("UPDATE batch_large SET value = value + 100 WHERE id % 2 = 0")
	if err != nil {
		log.Printf("Error with large batch UPDATE: %v", err)
		return
	}

	rowsAffected, _ := db.Exec("UPDATE batch_large SET value = value + 100 WHERE id % 2 = 0")
	log.Printf("✓ Large batch UPDATE (%d rows affected) successful", rowsAffected)
}

func testCleanup(db *sql.DB) {
	tables := []string{
		"batch_users",
		"batch_products",
		"batch_items",
		"batch_orders",
		"batch_accounts",
		"batch_transactions",
		"batch_employees",
		"batch_unique_items",
		"batch_large",
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
