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

	// Test 1: PREPARE and EXECUTE for SELECT
	log.Println("\n=== Test 1: PREPARE and EXECUTE for SELECT ===")
	testPrepareExecuteSelect(db)

	// Test 2: PREPARE and EXECUTE for INSERT
	log.Println("\n=== Test 2: PREPARE and EXECUTE for INSERT ===")
	testPrepareExecuteInsert(db)

	// Test 3: PREPARE and EXECUTE with multiple parameters
	log.Println("\n=== Test 3: PREPARE and EXECUTE with multiple parameters ===")
	testPrepareExecuteMultipleParams(db)

	// Test 4: PREPARE and EXECUTE with named parameters
	log.Println("\n=== Test 4: PREPARE and EXECUTE with named parameters ===")
	testPrepareExecuteNamedParams(db)

	// Test 5: PREPARE and EXECUTE for UPDATE
	log.Println("\n=== Test 5: PREPARE and EXECUTE for UPDATE ===")
	testPrepareExecuteUpdate(db)

	// Test 6: PREPARE and EXECUTE for DELETE
	log.Println("\n=== Test 6: PREPARE and EXECUTE for DELETE ===")
	testPrepareExecuteDelete(db)

	// Test 7: PREPARE with FROM clause
	log.Println("\n=== Test 7: PREPARE with FROM clause ===")
	testPrepareFromClause(db)

	// Test 8: PREPARE with AS clause
	log.Println("\n=== Test 8: PREPARE with AS clause ===")
	testPrepareAsClause(db)

	// Test 9: DEALLOCATE PREPARE
	log.Println("\n=== Test 9: DEALLOCATE PREPARE ===")
	testDeallocatePrepare(db)

	// Test 10: Error handling
	log.Println("\n=== Test 10: Error handling ===")
	testPrepareErrors(db)

	// Test 11: Cleanup
	log.Println("\n=== Test 11: Cleanup ===")
	testCleanup(db)

	log.Println("\n=== All Phase 17 tests completed! ===")
	log.Println("🎉 Phase 17: Prepared Statements - COMPLETE! 🎉")
}

func testPrepareExecuteSelect(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE users (id INTEGER, name TEXT, email TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: users")

	// Insert test data
	_, err = db.Exec("INSERT INTO users VALUES (1, 'Alice', 'alice@example.com')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}
	_, err = db.Exec("INSERT INTO users VALUES (2, 'Bob', 'bob@example.com')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	// PREPARE SELECT statement
	_, err = db.Exec("PREPARE get_user FROM 'SELECT * FROM users WHERE id = $id'")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return
	}
	log.Println("✓ Prepared statement: get_user")

	// EXECUTE prepared statement
	rows, err := db.Query("EXECUTE get_user USING @id = 1")
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return
	}
	defer rows.Close()

	log.Println("✓ Executed prepared statement with @id = 1")

	// Query results
	for rows.Next() {
		var id int
		var name string
		var email string
		err := rows.Scan(&id, &name, &email)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: ID=%d, Name=%s, Email=%s", id, name, email)
	}
}

func testPrepareExecuteInsert(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE products (id INTEGER, name TEXT, price REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: products")

	// PREPARE INSERT statement
	_, err = db.Exec("PREPARE insert_product FROM 'INSERT INTO products VALUES ($id, $name, $price)'")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return
	}
	log.Println("✓ Prepared statement: insert_product")

	// EXECUTE prepared statement
	result, err := db.Exec("EXECUTE insert_product USING @id = 1, @name = 'Product A', @price = 99.99")
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✓ Executed prepared statement: %d row(s) inserted", rowsAffected)

	// Verify data
	var id int
	var name string
	var price float64
	err = db.QueryRow("SELECT * FROM products WHERE id = 1").Scan(&id, &name, &price)
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	log.Printf("  Verified: ID=%d, Name=%s, Price=%.2f", id, name, price)
}

func testPrepareExecuteMultipleParams(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE orders (id INTEGER, user_id INTEGER, product_id INTEGER, quantity INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: orders")

	// PREPARE INSERT statement with multiple parameters
	_, err = db.Exec("PREPARE insert_order FROM 'INSERT INTO orders VALUES ($id, $user_id, $product_id, $quantity)'")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return
	}
	log.Println("✓ Prepared statement: insert_order")

	// EXECUTE with multiple parameters
	_, err = db.Exec("EXECUTE insert_order USING @id = 1, @user_id = 100, @product_id = 200, @quantity = 5")
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return
	}
	log.Println("✓ Executed prepared statement with multiple parameters")

	// EXECUTE with different parameters
	result, err := db.Exec("EXECUTE insert_order USING @id = 2, @user_id = 101, @product_id = 201, @quantity = 10")
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("✓ Executed prepared statement again: %d row(s) inserted", rowsAffected)
}

func testPrepareExecuteNamedParams(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE employees (id INTEGER, first_name TEXT, last_name TEXT, department TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: employees")

	// PREPARE INSERT with descriptive parameter names
	_, err = db.Exec("PREPARE insert_employee FROM 'INSERT INTO employees VALUES ($id, $first_name, $last_name, $department)'")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return
	}
	log.Println("✓ Prepared statement: insert_employee")

	// EXECUTE with descriptive parameter names
	_, err = db.Exec("EXECUTE insert_employee USING @id = 1, @first_name = 'John', @last_name = 'Doe', @department = 'Engineering'")
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return
	}
	log.Println("✓ Executed prepared statement with named parameters")
}

func testPrepareExecuteUpdate(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE inventory (id INTEGER, name TEXT, quantity INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: inventory")

	// Insert initial data
	_, err = db.Exec("INSERT INTO inventory VALUES (1, 'Product A', 10)")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	// PREPARE UPDATE statement
	_, err = db.Exec("PREPARE update_inventory FROM 'UPDATE inventory SET quantity = $quantity WHERE id = $id'")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return
	}
	log.Println("✓ Prepared statement: update_inventory")

	// EXECUTE UPDATE
	result, err := db.Exec("EXECUTE update_inventory USING @id = 1, @quantity = 20")
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✓ Executed prepared statement: %d row(s) updated", rowsAffected)

	// Verify update
	var quantity int
	err = db.QueryRow("SELECT quantity FROM inventory WHERE id = 1").Scan(&quantity)
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	log.Printf("  Verified: Quantity = %d", quantity)
}

func testPrepareExecuteDelete(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE logs (id INTEGER, message TEXT, level TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: logs")

	// Insert test data
	_, err = db.Exec("INSERT INTO logs VALUES (1, 'Info message', 'INFO')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}
	_, err = db.Exec("INSERT INTO logs VALUES (2, 'Error message', 'ERROR')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}
	_, err = db.Exec("INSERT INTO logs VALUES (3, 'Warning message', 'WARNING')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	// PREPARE DELETE statement
	_, err = db.Exec("PREPARE delete_log FROM 'DELETE FROM logs WHERE id = $id'")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return
	}
	log.Println("✓ Prepared statement: delete_log")

	// EXECUTE DELETE
	result, err := db.Exec("EXECUTE delete_log USING @id = 2")
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✓ Executed prepared statement: %d row(s) deleted", rowsAffected)

	// Verify remaining logs
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM logs").Scan(&count)
	if err != nil {
		log.Printf("Error querying count: %v", err)
		return
	}
	log.Printf("  Verified: %d log(s) remaining", count)
}

func testPrepareFromClause(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE customers (id INTEGER, name TEXT, email TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	// PREPARE with FROM clause (quoted SQL)
	_, err = db.Exec("PREPARE find_customer FROM 'SELECT * FROM customers WHERE email = $email'")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return
	}
	log.Println("✓ Prepared statement with FROM clause")
}

func testPrepareAsClause(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE items (id INTEGER, name TEXT, category TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	// PREPARE with AS clause (unquoted SQL)
	_, err = db.Exec("PREPARE get_items AS SELECT * FROM items WHERE category = $category")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return
	}
	log.Println("✓ Prepared statement with AS clause")
}

func testDeallocatePrepare(db *sql.DB) {
	// PREPARE a statement
	_, err := db.Exec("PREPARE temp_stmt FROM 'SELECT 1'")
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
		return
	}
	log.Println("✓ Prepared temporary statement")

	// DEALLOCATE prepared statement
	_, err = db.Exec("DEALLOCATE PREPARE temp_stmt")
	if err != nil {
		log.Printf("Error deallocating statement: %v", err)
		return
	}
	log.Println("✓ Deallocated prepared statement")

	// Try to execute deallocated statement (should fail)
	_, err = db.Exec("EXECUTE temp_stmt")
	if err != nil {
		log.Printf("✓ Correctly rejected execution of deallocated statement: %v", err)
		return
	}
	log.Printf("✗ Failed: Executed deallocated statement")
}

func testPrepareErrors(db *sql.DB) {
	// Try to execute non-existent prepared statement
	_, err := db.Exec("EXECUTE non_existent_stmt")
	if err != nil {
		log.Printf("✓ Correctly rejected non-existent statement: %v", err)
	}

	// Try to deallocate non-existent prepared statement
	_, err = db.Exec("DEALLOCATE PREPARE non_existent_stmt")
	if err != nil {
		log.Printf("✓ Correctly rejected deallocation of non-existent statement: %v", err)
	}
}

func testCleanup(db *sql.DB) {
	tables := []string{
		"users",
		"products",
		"orders",
		"employees",
		"inventory",
		"logs",
		"customers",
		"items",
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
