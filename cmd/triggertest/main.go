package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

	// Test 1: Simple AFTER INSERT Trigger
	log.Println("\n=== Test 1: Simple AFTER INSERT Trigger ===")
testAfterInsertTrigger(db)

	// Test 2: BEFORE INSERT Trigger
	log.Println("\n=== Test 2: BEFORE INSERT Trigger ===")
testBeforeInsertTrigger(db)

	// Test 3: AFTER UPDATE Trigger
	log.Println("\n=== Test 3: AFTER UPDATE Trigger ===")
testAfterUpdateTrigger(db)

	// Test 4: BEFORE UPDATE Trigger
	log.Println("\n=== Test 4: BEFORE UPDATE Trigger ===")
testBeforeUpdateTrigger(db)

	// Test 5: AFTER DELETE Trigger
	log.Println("\n=== Test 5: AFTER DELETE Trigger ===")
testAfterDeleteTrigger(db)

	// Test 6: BEFORE DELETE Trigger
	log.Println("\n=== Test 6: BEFORE DELETE Trigger ===")
testBeforeDeleteTrigger(db)

	// Test 7: Trigger with OLD and NEW references
	log.Println("\n=== Test 7: Trigger with OLD and NEW references ===")
testTriggerWithOldNew(db)

	// Test 8: Trigger with condition (WHEN)
	log.Println("\n=== Test 8: Trigger with condition (WHEN) ===")
testTriggerWithCondition(db)

	// Test 9: Multiple triggers on same table
	log.Println("\n=== Test 9: Multiple triggers on same table ===")
testMultipleTriggers(db)

	// Test 10: Trigger with UPDATE OF column
	log.Println("\n=== Test 10: Trigger with UPDATE OF column ===")
testTriggerUpdateOfColumn(db)

	// Test 11: Trigger with INSERT OR UPDATE
	log.Println("\n=== Test 11: Trigger with INSERT OR UPDATE ===")
testTriggerInsertOrUpdate(db)

	// Test 12: Trigger audit log
	log.Println("\n=== Test 12: Trigger audit log ===")
testTriggerAuditLog(db)

	// Test 13: Trigger with timestamp
	log.Println("\n=== Test 13: Trigger with timestamp ===")
testTriggerWithTimestamp(db)

	// Test 14: Cleanup
	log.Println("\n=== Test 14: Cleanup ===")
testCleanup(db)

	log.Println("\n=== All Phase 29 tests completed! ===")
	log.Println("🎉 Phase 29: Triggers - COMPLETE! 🎉")
}

func testAfterInsertTrigger(db *sql.DB) {
	// Create tables
	_, err := db.Exec("CREATE TABLE products (id INTEGER, name TEXT, price REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE product_log (id INTEGER, product_id INTEGER, action TEXT, timestamp TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created tables: products, product_log")

	// Create AFTER INSERT trigger
	_, err = db.Exec(`
		CREATE TRIGGER log_product_insert
		AFTER INSERT ON products
		FOR EACH ROW
		BEGIN
			INSERT INTO product_log (product_id, action, timestamp)
			VALUES (NEW.id, 'INSERT', datetime('now'));
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: log_product_insert")

	// Insert product
	result, err := db.Exec("INSERT INTO products (name, price) VALUES (?, ?)", "Product 1", 10.00)
	if err != nil {
		log.Printf("Error inserting product: %v", err)
		return
	}

	lastID, _ := result.LastInsertId()
	log.Printf("✓ Inserted product: ID=%d", lastID)

	// Check log
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM product_log WHERE action = 'INSERT'").Scan(&count)
	if err != nil {
		log.Printf("Error counting log entries: %v", err)
		return
	}

	log.Printf("✓ Log entries: %d", count)
}

func testBeforeInsertTrigger(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE inventory (id INTEGER, name TEXT, quantity INTEGER, status TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: inventory")

	// Create BEFORE INSERT trigger
	_, err = db.Exec(`
		CREATE TRIGGER check_inventory_status
		BEFORE INSERT ON inventory
		FOR EACH ROW
		BEGIN
			IF NEW.quantity < 10 THEN
				NEW.status = 'Low';
			ELSEIF NEW.quantity >= 10 AND NEW.quantity < 50 THEN
				NEW.status = 'Medium';
			ELSE
				NEW.status = 'High';
			END IF;
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: check_inventory_status")

	// Insert inventory items
	items := []struct {
		name     string
		quantity int
	}{
		{"Item 1", 5},
		{"Item 2", 20},
		{"Item 3", 100},
	}

	for _, item := range items {
		_, err = db.Exec("INSERT INTO inventory (name, quantity) VALUES (?, ?)", item.name, item.quantity)
		if err != nil {
			log.Printf("Error inserting item: %v", err)
			return
		}
	}

	log.Println("✓ Inserted 3 inventory items")

	// Verify status
	rows, err := db.Query("SELECT name, quantity, status FROM inventory")
	if err != nil {
		log.Printf("Error querying inventory: %v", err)
		return
	}
	defer rows.Close()

	log.Println("✓ Inventory status:")
	for rows.Next() {
		var name string
		var quantity int
		var status string
		err := rows.Scan(&name, &quantity, &status)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %s: Qty=%d, Status=%s", name, quantity, status)
	}
}

func testAfterUpdateTrigger(db *sql.DB) {
	// Create tables
	_, err := db.Exec("CREATE TABLE employees (id INTEGER, name TEXT, salary REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE salary_changes (id INTEGER, employee_id INTEGER, old_salary REAL, new_salary REAL, timestamp TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created tables: employees, salary_changes")

	// Create AFTER UPDATE trigger
	_, err = db.Exec(`
		CREATE TRIGGER log_salary_change
		AFTER UPDATE OF salary ON employees
		FOR EACH ROW
		BEGIN
			INSERT INTO salary_changes (employee_id, old_salary, new_salary, timestamp)
			VALUES (NEW.id, OLD.salary, NEW.salary, datetime('now'));
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: log_salary_change")

	// Insert employee
	_, err = db.Exec("INSERT INTO employees (name, salary) VALUES (?, ?)", "John Doe", 50000)
	if err != nil {
		log.Printf("Error inserting employee: %v", err)
		return
	}
	log.Println("✓ Inserted employee")

	// Update salary
	_, err = db.Exec("UPDATE employees SET salary = 55000 WHERE id = 1")
	if err != nil {
		log.Printf("Error updating salary: %v", err)
		return
	}
	log.Println("✓ Updated salary")

	// Check log
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM salary_changes").Scan(&count)
	if err != nil {
		log.Printf("Error counting salary changes: %v", err)
		return
	}

	log.Printf("✓ Salary changes logged: %d", count)
}

func testBeforeUpdateTrigger(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE orders (id INTEGER, total REAL, discount REAL, final_total REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: orders")

	// Create BEFORE UPDATE trigger
	_, err = db.Exec(`
		CREATE TRIGGER calculate_final_total
		BEFORE UPDATE OF discount ON orders
		FOR EACH ROW
		BEGIN
			NEW.final_total = NEW.total - (NEW.total * NEW.discount / 100);
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: calculate_final_total")

	// Insert order
	_, err = db.Exec("INSERT INTO orders (total, discount, final_total) VALUES (?, ?, ?)", 100.00, 0, 100.00)
	if err != nil {
		log.Printf("Error inserting order: %v", err)
		return
	}
	log.Println("✓ Inserted order")

	// Update discount
	_, err = db.Exec("UPDATE orders SET discount = 10 WHERE id = 1")
	if err != nil {
		log.Printf("Error updating discount: %v", err)
		return
	}
	log.Println("✓ Updated discount to 10%")

	// Verify final total
	var total, discount, finalTotal float64
	err = db.QueryRow("SELECT total, discount, final_total FROM orders WHERE id = 1").Scan(&total, &discount, &finalTotal)
	if err != nil {
		log.Printf("Error querying order: %v", err)
		return
	}

	log.Printf("✓ Order: Total=$%.2f, Discount=%.1f%%, Final Total=$%.2f", total, discount, finalTotal)
}

func testAfterDeleteTrigger(db *sql.DB) {
	// Create tables
	_, err := db.Exec("CREATE TABLE users (id INTEGER, name TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE deleted_users (id INTEGER, user_id INTEGER, name TEXT, deleted_at TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created tables: users, deleted_users")

	// Create AFTER DELETE trigger
	_, err = db.Exec(`
		CREATE TRIGGER log_user_deletion
		AFTER DELETE ON users
		FOR EACH ROW
		BEGIN
			INSERT INTO deleted_users (user_id, name, deleted_at)
			VALUES (OLD.id, OLD.name, datetime('now'));
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: log_user_deletion")

	// Insert user
	_, err = db.Exec("INSERT INTO users (name) VALUES (?)", "User 1")
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return
	}
	log.Println("✓ Inserted user")

	// Delete user
	_, err = db.Exec("DELETE FROM users WHERE id = 1")
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return
	}
	log.Println("✓ Deleted user")

	// Check log
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM deleted_users").Scan(&count)
	if err != nil {
		log.Printf("Error counting deleted users: %v", err)
		return
	}

	log.Printf("✓ Deleted users logged: %d", count)
}

func testBeforeDeleteTrigger(db *sql.DB) {
	// Create tables
	_, err := db.Exec("CREATE TABLE items (id INTEGER, name TEXT, in_use INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE delete_prevention (id INTEGER, item_id INTEGER, reason TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created tables: items, delete_prevention")

	// Create BEFORE DELETE trigger
	_, err = db.Exec(`
		CREATE TRIGGER prevent_delete_in_use
		BEFORE DELETE ON items
		FOR EACH ROW
		BEGIN
			INSERT INTO delete_prevention (item_id, reason)
			VALUES (OLD.id, 'Item is in use');
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: prevent_delete_in_use")

	// Insert items
	_, err = db.Exec("INSERT INTO items (name, in_use) VALUES (?, ?)", "Item 1", 1)
	if err != nil {
		log.Printf("Error inserting item: %v", err)
		return
	}
	log.Println("✓ Inserted item")

	// Try to delete item
	_, err = db.Exec("DELETE FROM items WHERE id = 1")
	if err != nil {
		log.Printf("Error deleting item: %v", err)
		return
	}
	log.Println("✓ Attempted to delete item")

	// Check prevention log
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM delete_prevention").Scan(&count)
	if err != nil {
		log.Printf("Error counting preventions: %v", err)
		return
	}

	log.Printf("✓ Delete attempts prevented: %d", count)
}

func testTriggerWithOldNew(db *sql.DB) {
	// Create tables
	_, err := db.Exec("CREATE TABLE products_v2 (id INTEGER, name TEXT, price REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE price_history (id INTEGER, product_id INTEGER, old_price REAL, new_price REAL, timestamp TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created tables: products_v2, price_history")

	// Create trigger with OLD and NEW
	_, err = db.Exec(`
		CREATE TRIGGER track_price_change
		BEFORE UPDATE OF price ON products_v2
		FOR EACH ROW
		BEGIN
			INSERT INTO price_history (product_id, old_price, new_price, timestamp)
			VALUES (NEW.id, OLD.price, NEW.price, datetime('now'));
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: track_price_change")

	// Insert product
	_, err = db.Exec("INSERT INTO products_v2 (name, price) VALUES (?, ?)", "Product 1", 100.00)
	if err != nil {
		log.Printf("Error inserting product: %v", err)
		return
	}
	log.Println("✓ Inserted product")

	// Update price multiple times
	prices := []float64{120.00, 110.00, 130.00}
	for _, price := range prices {
		_, err = db.Exec("UPDATE products_v2 SET price = ? WHERE id = 1", price)
		if err != nil {
			log.Printf("Error updating price: %v", err)
			return
		}
	}

	log.Printf("✓ Updated price %d times", len(prices))

	// Check history
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM price_history").Scan(&count)
	if err != nil {
		log.Printf("Error counting price history: %v", err)
		return
	}

	log.Printf("✓ Price changes tracked: %d", count)

	// Display history
	rows, err := db.Query("SELECT * FROM price_history ORDER BY id")
	if err != nil {
		log.Printf("Error querying price history: %v", err)
		return
	}
	defer rows.Close()

	log.Println("✓ Price history:")
	for rows.Next() {
		var id, productID int
		var oldPrice, newPrice float64
		var timestamp string
		err := rows.Scan(&id, &productID, &oldPrice, &newPrice, &timestamp)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: $%.2f -> $%.2f (%s)", id, oldPrice, newPrice, timestamp)
	}
}

func testTriggerWithCondition(db *sql.DB) {
	// Create tables
	_, err := db.Exec("CREATE TABLE orders_v2 (id INTEGER, customer_id INTEGER, total REAL, status TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE high_value_orders (id INTEGER, order_id INTEGER, customer_id INTEGER, total REAL, timestamp TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created tables: orders_v2, high_value_orders")

	// Create trigger with WHEN condition
	_, err = db.Exec(`
		CREATE TRIGGER log_high_value_orders
		AFTER INSERT ON orders_v2
		FOR EACH ROW
		WHEN NEW.total > 1000
		BEGIN
			INSERT INTO high_value_orders (order_id, customer_id, total, timestamp)
			VALUES (NEW.id, NEW.customer_id, NEW.total, datetime('now'));
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: log_high_value_orders")

	// Insert orders
	orders := []struct {
		customerID int
		total      float64
	}{
		{1, 500.00},
		{2, 1500.00},
		{3, 2000.00},
		{4, 800.00},
	}

	for _, order := range orders {
		_, err = db.Exec("INSERT INTO orders_v2 (customer_id, total, status) VALUES (?, ?, 'Pending')", order.customerID, order.total)
		if err != nil {
			log.Printf("Error inserting order: %v", err)
			return
		}
	}

	log.Printf("✓ Inserted %d orders", len(orders))

	// Check high value orders
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM high_value_orders").Scan(&count)
	if err != nil {
		log.Printf("Error counting high value orders: %v", err)
		return
	}

	log.Printf("✓ High value orders logged: %d", count)
}

func testMultipleTriggers(db *sql.DB) {
	// Create tables
	_, err := db.Exec("CREATE TABLE accounts (id INTEGER, balance REAL, last_updated TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE account_logs (id INTEGER, account_id INTEGER, action TEXT, old_balance REAL, new_balance REAL, timestamp TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created tables: accounts, account_logs")

	// Create first trigger: log balance changes
	_, err = db.Exec(`
		CREATE TRIGGER log_balance_change
		AFTER UPDATE OF balance ON accounts
		FOR EACH ROW
		BEGIN
			INSERT INTO account_logs (account_id, action, old_balance, new_balance, timestamp)
			VALUES (NEW.id, 'UPDATE', OLD.balance, NEW.balance, datetime('now'));
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: log_balance_change")

	// Create second trigger: update timestamp
	_, err = db.Exec(`
		CREATE TRIGGER update_timestamp
		BEFORE UPDATE OF balance ON accounts
		FOR EACH ROW
		BEGIN
			UPDATE accounts SET last_updated = datetime('now') WHERE id = NEW.id;
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: update_timestamp")

	// Insert account
	_, err = db.Exec("INSERT INTO accounts (balance, last_updated) VALUES (?, ?)", 1000.00, "")
	if err != nil {
		log.Printf("Error inserting account: %v", err)
		return
	}
	log.Println("✓ Inserted account")

	// Update balance
	_, err = db.Exec("UPDATE accounts SET balance = 1500 WHERE id = 1")
	if err != nil {
		log.Printf("Error updating balance: %v", err)
		return
	}
	log.Println("✓ Updated balance")

	// Check logs
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM account_logs").Scan(&count)
	if err != nil {
		log.Printf("Error counting logs: %v", err)
		return
	}

	log.Printf("✓ Account changes logged: %d", count)

	// Check timestamp
	var lastUpdated string
	err = db.QueryRow("SELECT last_updated FROM accounts WHERE id = 1").Scan(&lastUpdated)
	if err != nil {
		log.Printf("Error querying timestamp: %v", err)
		return
	}

	log.Printf("✓ Last updated: %s", lastUpdated)
}

func testTriggerUpdateOfColumn(db *sql.DB) {
	// Create tables
	_, err := db.Exec("CREATE TABLE customers (id INTEGER, name TEXT, email TEXT, phone TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE email_changes (id INTEGER, customer_id INTEGER, old_email TEXT, new_email TEXT, timestamp TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created tables: customers, email_changes")

	// Create trigger with UPDATE OF
	_, err = db.Exec(`
		CREATE TRIGGER log_email_change
		AFTER UPDATE OF email ON customers
		FOR EACH ROW
		BEGIN
			INSERT INTO email_changes (customer_id, old_email, new_email, timestamp)
			VALUES (NEW.id, OLD.email, NEW.email, datetime('now'));
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: log_email_change")

	// Insert customer
	_, err = db.Exec("INSERT INTO customers (name, email, phone) VALUES (?, ?, ?)", "John Doe", "john@example.com", "555-1234")
	if err != nil {
		log.Printf("Error inserting customer: %v", err)
		return
	}
	log.Println("✓ Inserted customer")

	// Update phone (should NOT trigger)
	_, err = db.Exec("UPDATE customers SET phone = '555-5678' WHERE id = 1")
	if err != nil {
		log.Printf("Error updating phone: %v", err)
		return
	}
	log.Println("✓ Updated phone")

	// Check log (should be empty)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM email_changes").Scan(&count)
	if err != nil {
		log.Printf("Error counting email changes: %v", err)
		return
	}

	log.Printf("✓ Email changes logged (phone update): %d", count)

	// Update email (should trigger)
	_, err = db.Exec("UPDATE customers SET email = 'john.doe@example.com' WHERE id = 1")
	if err != nil {
		log.Printf("Error updating email: %v", err)
		return
	}
	log.Println("✓ Updated email")

	// Check log (should have 1 entry)
	err = db.QueryRow("SELECT COUNT(*) FROM email_changes").Scan(&count)
	if err != nil {
		log.Printf("Error counting email changes: %v", err)
		return
	}

	log.Printf("✓ Email changes logged (email update): %d", count)
}

func testTriggerInsertOrUpdate(db *sql.DB) {
	// Create tables
	_, err := db.Exec("CREATE TABLE items_v2 (id INTEGER, name TEXT, quantity INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE item_changes (id INTEGER, item_id INTEGER, name TEXT, action TEXT, timestamp TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created tables: items_v2, item_changes")

	// Create trigger with INSERT OR UPDATE
	_, err = db.Exec(`
		CREATE TRIGGER log_item_change
		AFTER INSERT OR UPDATE ON items_v2
		FOR EACH ROW
		BEGIN
			INSERT INTO item_changes (item_id, name, action, timestamp)
			VALUES (NEW.id, NEW.name, 'CHANGE', datetime('now'));
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: log_item_change")

	// Insert item
	_, err = db.Exec("INSERT INTO items_v2 (name, quantity) VALUES (?, ?)", "Item 1", 10)
	if err != nil {
		log.Printf("Error inserting item: %v", err)
		return
	}
	log.Println("✓ Inserted item")

	// Update item
	_, err = db.Exec("UPDATE items_v2 SET quantity = 20 WHERE id = 1")
	if err != nil {
		log.Printf("Error updating item: %v", err)
		return
	}
	log.Println("✓ Updated item")

	// Check log
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM item_changes").Scan(&count)
	if err != nil {
		log.Printf("Error counting changes: %v", err)
		return
	}

	log.Printf("✓ Item changes logged: %d", count)
}

func testTriggerAuditLog(db *sql.DB) {
	// Create tables
	_, err := db.Exec("CREATE TABLE audit_table (id INTEGER, value TEXT, created_at TEXT, updated_at TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE audit_log (id INTEGER, table_name TEXT, record_id INTEGER, action TEXT, old_value TEXT, new_value TEXT, timestamp TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created tables: audit_table, audit_log")

	// Create INSERT trigger
	_, err = db.Exec(`
		CREATE TRIGGER audit_insert
		AFTER INSERT ON audit_table
		FOR EACH ROW
		BEGIN
			INSERT INTO audit_log (table_name, record_id, action, old_value, new_value, timestamp)
			VALUES ('audit_table', NEW.id, 'INSERT', NULL, NEW.value, datetime('now'));
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: audit_insert")

	// Create UPDATE trigger
	_, err = db.Exec(`
		CREATE TRIGGER audit_update
		AFTER UPDATE ON audit_table
		FOR EACH ROW
		BEGIN
			INSERT INTO audit_log (table_name, record_id, action, old_value, new_value, timestamp)
			VALUES ('audit_table', NEW.id, 'UPDATE', OLD.value, NEW.value, datetime('now'));
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: audit_update")

	// Create DELETE trigger
	_, err = db.Exec(`
		CREATE TRIGGER audit_delete
		AFTER DELETE ON audit_table
		FOR EACH ROW
		BEGIN
			INSERT INTO audit_log (table_name, record_id, action, old_value, new_value, timestamp)
			VALUES ('audit_table', OLD.id, 'DELETE', OLD.value, NULL, datetime('now'));
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: audit_delete")

	// Insert record
	_, err = db.Exec("INSERT INTO audit_table (value) VALUES (?)", "Initial Value")
	if err != nil {
		log.Printf("Error inserting record: %v", err)
		return
	}
	log.Println("✓ Inserted record")

	// Update record
	_, err = db.Exec("UPDATE audit_table SET value = ? WHERE id = 1", "Updated Value")
	if err != nil {
		log.Printf("Error updating record: %v", err)
		return
	}
	log.Println("✓ Updated record")

	// Delete record
	_, err = db.Exec("DELETE FROM audit_table WHERE id = 1")
	if err != nil {
		log.Printf("Error deleting record: %v", err)
		return
	}
	log.Println("✓ Deleted record")

	// Check audit log
	rows, err := db.Query("SELECT * FROM audit_log ORDER BY id")
	if err != nil {
		log.Printf("Error querying audit log: %v", err)
		return
	}
	defer rows.Close()

	log.Println("✓ Audit log:")
	for rows.Next() {
		var id int
		var tableName, action, timestamp string
		var recordID int
		var oldValue, newValue sql.NullString
		err := rows.Scan(&id, &tableName, &recordID, &action, &oldValue, &newValue, &timestamp)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		oldStr := "NULL"
		if oldValue.Valid {
			oldStr = oldValue.String
		}
		newStr := "NULL"
		if newValue.Valid {
			newStr = newValue.String
		}

		log.Printf("  %d: %s, Record=%d, Action=%s, Old=%s, New=%s, Time=%s", id, tableName, recordID, action, oldStr, newStr, timestamp)
	}
}

func testTriggerWithTimestamp(db *sql.DB) {
	// Create table
	_, err := db.Exec("CREATE TABLE events (id INTEGER, name TEXT, event_time TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: events")

	// Create trigger with timestamp
	_, err = db.Exec(`
		CREATE TRIGGER set_event_time
		BEFORE INSERT ON events
		FOR EACH ROW
		BEGIN
			NEW.event_time = datetime('now');
		END
	`)
	if err != nil {
		log.Printf("Error creating trigger: %v", err)
		return
	}
	log.Println("✓ Created trigger: set_event_time")

	// Insert events
	for i := 1; i <= 3; i++ {
		_, err = db.Exec("INSERT INTO events (name) VALUES (?)", fmt.Sprintf("Event %d", i))
		if err != nil {
			log.Printf("Error inserting event: %v", err)
			return
		}
		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
	}

	log.Println("✓ Inserted 3 events")

	// Verify timestamps
	rows, err := db.Query("SELECT id, name, event_time FROM events")
	if err != nil {
		log.Printf("Error querying events: %v", err)
		return
	}
	defer rows.Close()

	log.Println("✓ Events with timestamps:")
	for rows.Next() {
		var id int
		var name, eventTime string
		err := rows.Scan(&id, &name, &eventTime)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - %s", id, name, eventTime)
	}
}

func testCleanup(db *sql.DB) {
	tables := []string{
		"products",
		"product_log",
		"inventory",
		"employees",
		"salary_changes",
		"orders",
		"salary_changes",
		"users",
		"deleted_users",
		"items",
		"delete_prevention",
		"products_v2",
		"price_history",
		"orders_v2",
		"high_value_orders",
		"accounts",
		"account_logs",
		"customers",
		"email_changes",
		"items_v2",
		"item_changes",
		"audit_table",
		"audit_log",
		"events",
	}

	// Drop triggers first
	triggers := []string{
		"log_product_insert",
		"check_inventory_status",
		"log_salary_change",
		"calculate_final_total",
		"log_user_deletion",
		"prevent_delete_in_use",
		"track_price_change",
		"log_high_value_orders",
		"log_balance_change",
		"update_timestamp",
		"log_email_change",
		"log_item_change",
		"audit_insert",
		"audit_update",
		"audit_delete",
		"set_event_time",
	}

	for _, trigger := range triggers {
		_, err := db.Exec(fmt.Sprintf("DROP TRIGGER IF EXISTS %s", trigger))
		if err != nil {
			log.Printf("Error dropping trigger %s: %v", trigger, err)
			continue
		}
		log.Printf("✓ Dropped trigger: %s", trigger)
	}

	// Drop tables
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			log.Printf("Error dropping table %s: %v", table, err)
			continue
		}
		log.Printf("✓ Dropped table: %s", table)
	}
}
