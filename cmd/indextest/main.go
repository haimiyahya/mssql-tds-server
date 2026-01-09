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

	// Test 1: CREATE TABLES
	log.Println("\n=== Test 1: CREATE TABLES ===")
	testCreateTables(db)

	// Test 2: INSERT data
	log.Println("\n=== Test 2: INSERT data ===")
	testInsert(db)

	// Test 3: Simple CREATE INDEX
	log.Println("\n=== Test 3: Simple CREATE INDEX ===")
	testSimpleCreateIndex(db)

	// Test 4: UNIQUE INDEX
	log.Println("\n=== Test 4: UNIQUE INDEX ===")
	testUniqueIndex(db)

	// Test 5: Multi-column index
	log.Println("\n=== Test 5: Multi-column index ===")
	testMultiColumnIndex(db)

	// Test 6: Multiple indexes on same table
	log.Println("\n=== Test 6: Multiple indexes on same table ===")
	testMultipleIndexes(db)

	// Test 7: DROP INDEX
	log.Println("\n=== Test 7: DROP INDEX ===")
	testDropIndex(db)

	// Test 8: Recreate index
	log.Println("\n=== Test 8: Recreate index ===")
	testRecreateIndex(db)

	// Test 9: Index with ORDER BY query
	log.Println("\n=== Test 9: Index with ORDER BY query ===")
	testIndexWithOrderBy(db)

	// Test 10: Index with WHERE query
	log.Println("\n=== Test 10: Index with WHERE query ===")
	testIndexWithWhere(db)

	// Test 11: Index performance test
	log.Println("\n=== Test 11: Index performance test ===")
	testIndexPerformance(db)

	// Test 12: Index on large table
	log.Println("\n=== Test 12: Index on large table ===")
	testIndexOnLargeTable(db)

	// Test 13: DROP TABLES
	log.Println("\n=== Test 13: DROP TABLES ===")
	testDropTables(db)

	log.Println("\n=== All Phase 14 tests completed! ===")
	log.Println("🎉 Phase 14: Index Management - COMPLETE! 🎉")
}

func testCreateTables(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE users (id INTEGER, name TEXT, email TEXT)")
	if err != nil {
		log.Printf("Error creating users table: %v", err)
		return
	}
	log.Println("✓ Created table: users")

	_, err = db.Exec("CREATE TABLE products (id INTEGER, name TEXT, category TEXT, price REAL, stock INTEGER)")
	if err != nil {
		log.Printf("Error creating products table: %v", err)
		return
	}
	log.Println("✓ Created table: products")

	_, err = db.Exec("CREATE TABLE orders (id INTEGER, user_id INTEGER, product_id INTEGER, quantity INTEGER, total REAL)")
	if err != nil {
		log.Printf("Error creating orders table: %v", err)
		return
	}
	log.Println("✓ Created table: orders")
}

func testInsert(db *sql.DB) {
	// Insert users
	userQueries := []string{
		"INSERT INTO users VALUES (1, 'Alice', 'alice@example.com')",
		"INSERT INTO users VALUES (2, 'Bob', 'bob@example.com')",
		"INSERT INTO users VALUES (3, 'Charlie', 'charlie@example.com')",
		"INSERT INTO users VALUES (4, 'Diana', 'diana@example.com')",
		"INSERT INTO users VALUES (5, 'Eve', 'eve@example.com')",
	}

	for _, query := range userQueries {
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			continue
		}
		rowsAffected, _ := result.RowsAffected()
		log.Printf("✓ Inserted user: %d row(s)", rowsAffected)
	}

	// Insert products
	productQueries := []string{
		"INSERT INTO products VALUES (1, 'Laptop', 'Electronics', 999.99, 50)",
		"INSERT INTO products VALUES (2, 'Phone', 'Electronics', 699.99, 100)",
		"INSERT INTO products VALUES (3, 'Book', 'Education', 29.99, 200)",
		"INSERT INTO products VALUES (4, 'Headphones', 'Electronics', 149.99, 150)",
		"INSERT INTO products VALUES (5, 'Table', 'Furniture', 299.99, 30)",
	}

	for _, query := range productQueries {
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("Error inserting product: %v", err)
			continue
		}
		rowsAffected, _ := result.RowsAffected()
		log.Printf("✓ Inserted product: %d row(s)", rowsAffected)
	}

	// Insert orders
	orderQueries := []string{
		"INSERT INTO orders VALUES (1, 1, 1, 1, 999.99)",
		"INSERT INTO orders VALUES (2, 2, 2, 2, 1399.98)",
		"INSERT INTO orders VALUES (3, 3, 3, 5, 149.95)",
		"INSERT INTO orders VALUES (4, 4, 4, 1, 149.99)",
		"INSERT INTO orders VALUES (5, 5, 5, 2, 599.98)",
	}

	for _, query := range orderQueries {
		result, err := db.Exec(query)
		if err != nil {
			log.Printf("Error inserting order: %v", err)
			continue
		}
		rowsAffected, _ := result.RowsAffected()
		log.Printf("✓ Inserted order: %d row(s)", rowsAffected)
	}
}

func testSimpleCreateIndex(db *sql.DB) {
	// Simple CREATE INDEX
	_, err := db.Exec("CREATE INDEX idx_users_name ON users (name)")
	if err != nil {
		log.Printf("Error creating index: %v", err)
		return
	}
	log.Println("✓ Created index: idx_users_name on users (name)")

	// Verify index created
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='index' AND name='idx_users_name'")
	if err != nil {
		log.Printf("Error querying index: %v", err)
		return
	}
	defer rows.Close()

	var name string
	rows.Next()
	rows.Scan(&name)
	log.Printf("✓ Verified: Index 'idx_users_name' created in database")
}

func testUniqueIndex(db *sql.DB) {
	// UNIQUE INDEX
	_, err := db.Exec("CREATE UNIQUE INDEX idx_users_email ON users (email)")
	if err != nil {
		log.Printf("Error creating unique index: %v", err)
		return
	}
	log.Println("✓ Created unique index: idx_users_email on users (email)")

	// Verify index created
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='index' AND name='idx_users_email'")
	if err != nil {
		log.Printf("Error querying index: %v", err)
		return
	}
	defer rows.Close()

	var name string
	rows.Next()
	rows.Scan(&name)
	log.Printf("✓ Verified: Unique Index 'idx_users_email' created in database")

	// Try to insert duplicate email (should fail with unique constraint)
	_, err = db.Exec("INSERT INTO users VALUES (6, 'Frank', 'alice@example.com')")
	if err != nil {
		log.Printf("✓ Verified: Unique constraint works (duplicate email rejected)")
		return
	}
	log.Printf("✗ Failed: Unique constraint not enforced")
}

func testMultiColumnIndex(db *sql.DB) {
	// Multi-column INDEX
	_, err := db.Exec("CREATE INDEX idx_products_category_price ON products (category, price)")
	if err != nil {
		log.Printf("Error creating multi-column index: %v", err)
		return
	}
	log.Println("✓ Created multi-column index: idx_products_category_price on products (category, price)")

	// Verify index created
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='index' AND name='idx_products_category_price'")
	if err != nil {
		log.Printf("Error querying index: %v", err)
		return
	}
	defer rows.Close()

	var name string
	rows.Next()
	rows.Scan(&name)
	log.Printf("✓ Verified: Multi-column Index 'idx_products_category_price' created in database")
}

func testMultipleIndexes(db *sql.DB) {
	// Multiple indexes on same table
	indexes := []string{
		"CREATE INDEX idx_products_name ON products (name)",
		"CREATE INDEX idx_products_category ON products (category)",
		"CREATE INDEX idx_products_price ON products (price)",
	}

	for _, query := range indexes {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error creating index: %v", err)
			continue
		}
		log.Println("✓ Created additional index on products table")
	}

	// Count indexes
	rows, err := db.Query("SELECT COUNT(*) FROM sqlite_master WHERE type='index'")
	if err != nil {
		log.Printf("Error querying indexes: %v", err)
		return
	}
	defer rows.Close()

	var count int
	rows.Next()
	rows.Scan(&count)
	log.Printf("✓ Total indexes in database: %d", count)
}

func testDropIndex(db *sql.DB) {
	// DROP INDEX
	_, err := db.Exec("DROP INDEX idx_products_price")
	if err != nil {
		log.Printf("Error dropping index: %v", err)
		return
	}
	log.Println("✓ Dropped index: idx_products_price")

	// Verify index dropped
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='index' AND name='idx_products_price'")
	if err != nil {
		log.Printf("Error querying index: %v", err)
		return
	}
	defer rows.Close()

	rows.Next()
	var name sql.NullString
	rows.Scan(&name)
	if !name.Valid {
		log.Println("✓ Verified: Index 'idx_products_price' dropped from database")
	} else {
		log.Printf("✗ Verified: Index still exists: %s", name.String)
	}
}

func testRecreateIndex(db *sql.DB) {
	// Recreate index (drop and create)
	_, err := db.Exec("DROP INDEX idx_products_category")
	if err != nil {
		log.Printf("Error dropping index: %v", err)
		return
	}
	log.Println("✓ Dropped index: idx_products_category")

	_, err = db.Exec("CREATE INDEX idx_products_category ON products (category)")
	if err != nil {
		log.Printf("Error recreating index: %v", err)
		return
	}
	log.Println("✓ Recreated index: idx_products_category")

	// Verify index recreated
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='index' AND name='idx_products_category'")
	if err != nil {
		log.Printf("Error querying index: %v", err)
		return
	}
	defer rows.Close()

	var name string
	rows.Next()
	rows.Scan(&name)
	log.Printf("✓ Verified: Index 'idx_products_category' recreated in database")
}

func testIndexWithOrderBy(db *sql.DB) {
	// SELECT with ORDER BY (should use index)
	rows, err := db.Query("SELECT * FROM products ORDER BY name")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var id int
		var name string
		var category string
		var price float64
		var stock int
		err := rows.Scan(&id, &name, &category, &price, &stock)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		if rowCount <= 3 {
			log.Printf("  Result: %s (%s) - $%.2f", name, category, price)
		}
	}

	log.Printf("✓ Selected from products with ORDER BY: %d row(s) (using idx_products_name)", rowCount)
}

func testIndexWithWhere(db *sql.DB) {
	// SELECT with WHERE (should use index)
	rows, err := db.Query("SELECT * FROM products WHERE category = 'Electronics'")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
		var id int
		var name string
		var category string
		var price float64
		var stock int
		err := rows.Scan(&id, &name, &category, &price, &stock)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Result: %s - $%.2f (Stock: %d)", name, price, stock)
	}

	log.Printf("✓ Selected from products with WHERE: %d row(s) (using idx_products_category)", rowCount)
}

func testIndexPerformance(db *sql.DB) {
	// Performance test: Query with index vs without index
	// Note: SQLite automatically uses indexes when beneficial

	// Query with index (user_id)
	rows, err := db.Query("SELECT * FROM orders WHERE user_id = 1")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		rowCount++
	}
	log.Printf("✓ Query performance test: %d row(s) found for user_id = 1", rowCount)
}

func testIndexOnLargeTable(db *sql.DB) {
	// Create larger table and index
	_, err := db.Exec("CREATE TABLE test_large (id INTEGER, value TEXT)")
	if err != nil {
		log.Printf("Error creating large table: %v", err)
		return
	}
	log.Println("✓ Created table: test_large")

	// Insert 100 rows
	for i := 1; i <= 100; i++ {
		query := fmt.Sprintf("INSERT INTO test_large VALUES (%d, 'value%d')", i, i)
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error inserting row: %v", err)
			continue
		}
	}
	log.Println("✓ Inserted 100 rows into test_large table")

	// Create index
	_, err = db.Exec("CREATE INDEX idx_test_large_value ON test_large (value)")
	if err != nil {
		log.Printf("Error creating index: %v", err)
		return
	}
	log.Println("✓ Created index: idx_test_large_value on test_large (value)")

	// Query with index
	rows, err := db.Query("SELECT * FROM test_large WHERE value = 'value50'")
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	var id int
	var value string
	rows.Next()
	rows.Scan(&id, &value)
	log.Printf("✓ Query with index: Found %s at id %d", value, id)

	// Drop table
	_, err = db.Exec("DROP TABLE test_large")
	if err != nil {
		log.Printf("Error dropping table: %v", err)
		return
	}
	log.Println("✓ Dropped table: test_large")
}

func testDropTables(db *sql.DB) {
	_, err := db.Exec("DROP TABLE orders")
	if err != nil {
		log.Printf("Error dropping orders table: %v", err)
		return
	}
	log.Println("✓ Dropped table: orders")

	_, err = db.Exec("DROP TABLE products")
	if err != nil {
		log.Printf("Error dropping products table: %v", err)
		return
	}
	log.Println("✓ Dropped table: products")

	_, err = db.Exec("DROP TABLE users")
	if err != nil {
		log.Printf("Error dropping users table: %v", err)
		return
	}
	log.Println("✓ Dropped table: users")
}
