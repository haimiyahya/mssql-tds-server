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

	// Test 1: Simple CTE
	log.Println("\n=== Test 1: Simple CTE ===")
testSimpleCTE(db)

	// Test 2: Multiple Column CTE
	log.Println("\n=== Test 2: Multiple Column CTE ===")
testMultipleColumnCTE(db)

	// Test 3: CTE with WHERE
	log.Println("\n=== Test 3: CTE with WHERE ===")
testCTEWithWhere(db)

	// Test 4: CTE with JOIN
	log.Println("\n=== Test 4: CTE with JOIN ===")
testCTEWithJoin(db)

	// Test 5: Multiple CTEs
	log.Println("\n=== Test 5: Multiple CTEs ===")
testMultipleCTEs(db)

	// Test 6: CTE with GROUP BY
	log.Println("\n=== Test 6: CTE with GROUP BY ===")
testCTEWithGroupBy(db)

	// Test 7: CTE with ORDER BY
	log.Println("\n=== Test 7: CTE with ORDER BY ===")
testCTEWithOrderBy(db)

	// Test 8: Recursive CTE
	log.Println("\n=== Test 8: Recursive CTE ===")
testRecursiveCTE(db)

	// Test 9: CTE in INSERT
	log.Println("\n=== Test 9: CTE in INSERT ===")
testCTEInInsert(db)

	// Test 10: CTE in UPDATE
	log.Println("\n=== Test 10: CTE in UPDATE ===")
testCTEInUpdate(db)

	// Test 11: Cleanup
	log.Println("\n=== Test 11: Cleanup ===")
testCleanup(db)

	log.Println("\n=== All Phase 27 tests completed! ===")
	log.Println("🎉 Phase 27: Common Table Expressions (CTE) - COMPLETE! 🎉")
}

func testSimpleCTE(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE employees (id INTEGER, name TEXT, department TEXT, salary INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: employees")

	// Insert test data
	for i := 1; i <= 10; i++ {
		dept := fmt.Sprintf("Department %d", (i%3)+1)
		salary := (i * 10000)
		_, err = db.Exec("INSERT INTO employees VALUES (?, ?, ?, ?)", i, fmt.Sprintf("Employee %d", i), dept, salary)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 10 employees")

	// Simple CTE - high salary employees
	log.Println("✓ Simple CTE - High Salary Employees:")
	rows, err := db.Query(`
		WITH HighSalary AS (
			SELECT * FROM employees WHERE salary >= 50000
		)
		SELECT * FROM HighSalary
	`)
	if err != nil {
		log.Printf("Error executing CTE query: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name string
		var department string
		var salary int
		err := rows.Scan(&id, &name, &department, &salary)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		count++
		log.Printf("  %d: %s - %s - $%d", id, name, department, salary)
	}

	log.Printf("✓ Found %d high salary employees", count)
}

func testMultipleColumnCTE(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE products (id INTEGER, name TEXT, price REAL, quantity INTEGER, category TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: products")

	// Insert test data
	for i := 1; i <= 10; i++ {
		category := fmt.Sprintf("Category %d", (i%4)+1)
		price := float64(i) * 10.5
		quantity := 100 - (i * 5)
		_, err = db.Exec("INSERT INTO products VALUES (?, ?, ?, ?, ?)", i, fmt.Sprintf("Product %d", i), price, quantity, category)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 10 products")

	// CTE with multiple columns - inventory value
	log.Println("✓ CTE with Multiple Columns - Inventory Value:")
	rows, err := db.Query(`
		WITH InventoryValue AS (
			SELECT id, name, price, quantity, (price * quantity) as total_value
			FROM products
		)
		SELECT id, name, total_value FROM InventoryValue WHERE total_value > 500
	`)
	if err != nil {
		log.Printf("Error executing CTE query: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name string
		var totalValue float64
		err := rows.Scan(&id, &name, &totalValue)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		count++
		log.Printf("  %d: %s - Total Value: $%.2f", id, name, totalValue)
	}

	log.Printf("✓ Found %d products with total value > $500", count)
}

func testCTEWithWhere(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE sales (id INTEGER, product_id INTEGER, quantity INTEGER, sale_date TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: sales")

	// Insert test data
	for i := 1; i <= 20; i++ {
		productID := (i % 5) + 1
		quantity := (i % 10) + 1
		date := fmt.Sprintf("2024-01-%02d", (i%28)+1)
		_, err = db.Exec("INSERT INTO sales VALUES (?, ?, ?, ?)", i, productID, quantity, date)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 20 sales")

	// CTE with WHERE clause - recent sales
	log.Println("✓ CTE with WHERE Clause - Recent Sales:")
	rows, err := db.Query(`
		WITH RecentSales AS (
			SELECT * FROM sales WHERE sale_date >= '2024-01-15'
		)
		SELECT * FROM RecentSales ORDER BY sale_date
	`)
	if err != nil {
		log.Printf("Error executing CTE query: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var productID int
		var quantity int
		var saleDate string
		err := rows.Scan(&id, &productID, &quantity, &saleDate)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		count++
		if count <= 5 {
			log.Printf("  %d: Product %d - Qty: %d - Date: %s", id, productID, quantity, saleDate)
		}
	}

	log.Printf("✓ Found %d recent sales", count)
}

func testCTEWithJoin(db *sql.DB) {
	// Create test tables
	_, err := db.Exec("CREATE TABLE orders (id INTEGER, customer_id INTEGER, order_date TEXT, total REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE customers (id INTEGER, name TEXT, city TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created tables: orders, customers")

	// Insert test data
	for i := 1; i <= 10; i++ {
		city := fmt.Sprintf("City %d", (i%3)+1)
		_, err = db.Exec("INSERT INTO customers VALUES (?, ?, ?)", i, fmt.Sprintf("Customer %d", i), city)
		if err != nil {
			log.Printf("Error inserting customer: %v", err)
			return
		}
	}
	for i := 1; i <= 20; i++ {
		customerID := (i % 10) + 1
		total := float64(i) * 100.5
		date := fmt.Sprintf("2024-01-%02d", (i%28)+1)
		_, err = db.Exec("INSERT INTO orders VALUES (?, ?, ?, ?)", i, customerID, date, total)
		if err != nil {
			log.Printf("Error inserting order: %v", err)
			return
		}
	}
	log.Println("✓ Inserted test data")

	// CTE with JOIN - customer orders summary
	log.Println("✓ CTE with JOIN - Customer Orders Summary:")
	rows, err := db.Query(`
		WITH CustomerOrders AS (
			SELECT c.id, c.name, COUNT(o.id) as order_count, SUM(o.total) as total_spent
			FROM customers c
			LEFT JOIN orders o ON c.id = o.customer_id
			GROUP BY c.id, c.name
		)
		SELECT * FROM CustomerOrders WHERE order_count > 0
	`)
	if err != nil {
		log.Printf("Error executing CTE query: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name string
		var orderCount int
		var totalSpent sql.NullFloat64
		err := rows.Scan(&id, &name, &orderCount, &totalSpent)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		count++
		spent := 0.0
		if totalSpent.Valid {
			spent = totalSpent.Float64
		}
		log.Printf("  %d: %s - Orders: %d - Total: $%.2f", id, name, orderCount, spent)
	}

	log.Printf("✓ Found %d customers with orders", count)
}

func testMultipleCTEs(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE inventory (id INTEGER, product TEXT, quantity INTEGER, location TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: inventory")

	// Insert test data
	for i := 1; i <= 15; i++ {
		location := fmt.Sprintf("Location %d", (i%4)+1)
		quantity := 10 + (i * 5)
		_, err = db.Exec("INSERT INTO inventory VALUES (?, ?, ?, ?)", i, fmt.Sprintf("Product %d", i), quantity, location)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 15 inventory items")

	// Multiple CTEs - low stock and high stock
	log.Println("✓ Multiple CTEs - Low Stock and High Stock:")
	rows, err := db.Query(`
		WITH LowStock AS (
			SELECT * FROM inventory WHERE quantity < 30
		),
		HighStock AS (
			SELECT * FROM inventory WHERE quantity >= 60
		)
		SELECT 'Low Stock' as category, * FROM LowStock
		UNION ALL
		SELECT 'High Stock' as category, * FROM HighStock
	`)
	if err != nil {
		log.Printf("Error executing CTE query: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var category string
		var id int
		var product string
		var quantity int
		var location string
		err := rows.Scan(&category, &id, &product, &quantity, &location)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		count++
		log.Printf("  %s: %d - %s - Qty: %d - %s", category, id, product, quantity, location)
	}

	log.Printf("✓ Found %d items (low stock or high stock)", count)
}

func testCTEWithGroupBy(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE transactions (id INTEGER, account TEXT, amount REAL, transaction_date TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: transactions")

	// Insert test data
	accounts := []string{"Account1", "Account2", "Account3"}
	for i := 1; i <= 30; i++ {
		account := accounts[i%3]
		amount := float64((i%10)+1) * 100.0
		date := fmt.Sprintf("2024-01-%02d", (i%28)+1)
		_, err = db.Exec("INSERT INTO transactions VALUES (?, ?, ?, ?)", i, account, amount, date)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 30 transactions")

	// CTE with GROUP BY - account summary
	log.Println("✓ CTE with GROUP BY - Account Summary:")
	rows, err := db.Query(`
		WITH AccountSummary AS (
			SELECT account, COUNT(*) as transaction_count, SUM(amount) as total_amount
			FROM transactions
			GROUP BY account
		)
		SELECT * FROM AccountSummary ORDER BY total_amount DESC
	`)
	if err != nil {
		log.Printf("Error executing CTE query: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var account string
		var transactionCount int
		var totalAmount float64
		err := rows.Scan(&account, &transactionCount, &totalAmount)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		count++
		log.Printf("  %s: Transactions: %d - Total: $%.2f", account, transactionCount, totalAmount)
	}

	log.Printf("✓ Found %d accounts", count)
}

func testCTEWithOrderBy(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE scores (id INTEGER, player TEXT, score INTEGER, game_date TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: scores")

	// Insert test data
	players := []string{"Alice", "Bob", "Charlie", "Diana"}
	for i := 1; i <= 20; i++ {
		player := players[i%4]
		score := (i * 10) + ((i % 7) * 50)
		date := fmt.Sprintf("2024-01-%02d", (i%28)+1)
		_, err = db.Exec("INSERT INTO scores VALUES (?, ?, ?, ?)", i, player, score, date)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 20 scores")

	// CTE with ORDER BY - top scores
	log.Println("✓ CTE with ORDER BY - Top 5 Scores:")
	rows, err := db.Query(`
		WITH TopScores AS (
			SELECT * FROM scores ORDER BY score DESC
		)
		SELECT * FROM TopScores LIMIT 5
	`)
	if err != nil {
		log.Printf("Error executing CTE query: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var player string
		var score int
		var gameDate string
		err := rows.Scan(&id, &player, &score, &gameDate)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		count++
		log.Printf("  %d. %s - Score: %d - Date: %s", count, player, score, gameDate)
	}

	log.Printf("✓ Found top %d scores", count)
}

func testRecursiveCTE(db *sql.DB) {
	// Create test table for hierarchy
	_, err := db.Exec("CREATE TABLE employees_hierarchy (id INTEGER, name TEXT, manager_id INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: employees_hierarchy")

	// Insert test data (hierarchy)
	employees := []struct {
		id       int
		name     string
		managerID sql.NullInt64
	}{
		{1, "CEO", sql.NullInt64{}},
		{2, "VP1", sql.NullInt64{Int64: 1, Valid: true}},
		{3, "VP2", sql.NullInt64{Int64: 1, Valid: true}},
		{4, "Manager1", sql.NullInt64{Int64: 2, Valid: true}},
		{5, "Manager2", sql.NullInt64{Int64: 2, Valid: true}},
		{6, "Manager3", sql.NullInt64{Int64: 3, Valid: true}},
		{7, "Employee1", sql.NullInt64{Int64: 4, Valid: true}},
		{8, "Employee2", sql.NullInt64{Int64: 4, Valid: true}},
		{9, "Employee3", sql.NullInt64{Int64: 5, Valid: true}},
	}

	for _, emp := range employees {
		var managerID interface{}
		if emp.managerID.Valid {
			managerID = emp.managerID.Int64
		} else {
			managerID = nil
		}
		_, err = db.Exec("INSERT INTO employees_hierarchy VALUES (?, ?, ?)", emp.id, emp.name, managerID)
		if err != nil {
			log.Printf("Error inserting employee: %v", err)
			return
		}
	}
	log.Println("✓ Inserted employee hierarchy")

	// Recursive CTE - find all employees under CEO
	log.Println("✓ Recursive CTE - Organization Hierarchy:")
	rows, err := db.Query(`
		WITH RECURSIVE OrgHierarchy AS (
			-- Base case: CEO
			SELECT id, name, manager_id, 0 as level
			FROM employees_hierarchy
			WHERE manager_id IS NULL

			UNION ALL

			-- Recursive case: employees under previous level
			SELECT e.id, e.name, e.manager_id, h.level + 1
			FROM employees_hierarchy e
			JOIN OrgHierarchy h ON e.manager_id = h.id
		)
		SELECT * FROM OrgHierarchy ORDER BY level, id
	`)
	if err != nil {
		log.Printf("Error executing recursive CTE query: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name string
		var managerID sql.NullInt64
		var level int
		err := rows.Scan(&id, &name, &managerID, &level)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		managerStr := "NULL"
		if managerID.Valid {
			managerStr = fmt.Sprintf("%d", managerID.Int64)
		}

		indent := ""
		for i := 0; i < level; i++ {
			indent += "  "
		}

		count++
		log.Printf("  %s%d: %s (Manager: %s, Level: %d)", indent, id, name, managerStr, level)
	}

	log.Printf("✓ Found %d employees in hierarchy", count)
}

func testCTEInInsert(db *sql.DB) {
	// Create test tables
	_, err := db.Exec("CREATE TABLE source_products (id INTEGER, name TEXT, price REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE premium_products (id INTEGER, name TEXT, price REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created tables: source_products, premium_products")

	// Insert test data
	for i := 1; i <= 10; i++ {
		price := float64(i) * 25.0
		_, err = db.Exec("INSERT INTO source_products VALUES (?, ?, ?)", i, fmt.Sprintf("Product %d", i), price)
		if err != nil {
			log.Printf("Error inserting product: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 10 source products")

	// CTE in INSERT - insert premium products
	log.Println("✓ CTE in INSERT - Insert Premium Products:")
	result, err := db.Exec(`
		INSERT INTO premium_products (id, name, price)
		WITH Premium AS (
			SELECT * FROM source_products WHERE price >= 100
		)
		SELECT * FROM Premium
	`)
	if err != nil {
		log.Printf("Error executing CTE in INSERT: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✓ Inserted %d premium products", rowsAffected)

	// Verify
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM premium_products").Scan(&count)
	if err != nil {
		log.Printf("Error counting premium products: %v", err)
		return
	}
	log.Printf("✓ Premium products count: %d", count)
}

func testCTEInUpdate(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE items (id INTEGER, name TEXT, quantity INTEGER, status TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: items")

	// Insert test data
	for i := 1; i <= 10; i++ {
		quantity := 5 + (i * 2)
		_, err = db.Exec("INSERT INTO items VALUES (?, ?, ?, ?)", i, fmt.Sprintf("Item %d", i), quantity, "In Stock")
		if err != nil {
			log.Printf("Error inserting item: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 10 items")

	// CTE in UPDATE - update low stock items
	log.Println("✓ CTE in UPDATE - Update Low Stock Items:")
	result, err := db.Exec(`
		UPDATE items SET status = 'Low Stock'
		WHERE id IN (
			WITH LowStockItems AS (
				SELECT id FROM items WHERE quantity < 15
			)
			SELECT id FROM LowStockItems
		)
	`)
	if err != nil {
		log.Printf("Error executing CTE in UPDATE: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✓ Updated %d items to 'Low Stock'", rowsAffected)

	// Verify
	log.Println("✓ Current item statuses:")
	rows, err := db.Query("SELECT id, name, quantity, status FROM items ORDER BY id")
	if err != nil {
		log.Printf("Error querying items: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var quantity int
		var status string
		err := rows.Scan(&id, &name, &quantity, &status)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - Qty: %d - Status: %s", id, name, quantity, status)
	}
}

func testCleanup(db *sql.DB) {
	tables := []string{
		"employees",
		"products",
		"sales",
		"orders",
		"customers",
		"inventory",
		"transactions",
		"scores",
		"employees_hierarchy",
		"source_products",
		"premium_products",
		"items",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			log.Printf("Error dropping table %s: %v", table, err)
			continue
		}
		log.Printf("✓ Dropped table: %s", table)
	}
}
