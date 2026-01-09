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

	// Test 1: Basic EXPLAIN Command
	log.Println("\n=== Test 1: Basic EXPLAIN Command ===")
	testBasicExplain(db)

	// Test 2: EXPLAIN with Index
	log.Println("\n=== Test 2: EXPLAIN with Index ===")
testExplainWithIndex(db)

	// Test 3: EXPLAIN with JOIN
	log.Println("\n=== Test 3: EXPLAIN with JOIN ===")
testExplainWithJoin(db)

	// Test 4: EXPLAIN with Subquery
	log.Println("\n=== Test 4: EXPLAIN with Subquery ===")
testExplainWithSubquery(db)

	// Test 5: EXPLAIN with GROUP BY
	log.Println("\n=== Test 5: EXPLAIN with GROUP BY ===")
testExplainWithGroupBy(db)

	// Test 6: EXPLAIN with ORDER BY
	log.Println("\n=== Test 6: EXPLAIN with ORDER BY ===")
testExplainWithOrderBy(db)

	// Test 7: EXPLAIN QUERY PLAN Command
	log.Println("\n=== Test 7: EXPLAIN QUERY PLAN Command ===")
testExplainQueryPlan(db)

	// Test 8: EXPLAIN with WHERE Clause
	log.Println("\n=== Test 8: EXPLAIN with WHERE Clause ===")
testExplainWithWhere(db)

	// Test 9: EXPLAIN with LIMIT
	log.Println("\n=== Test 9: EXPLAIN with LIMIT ===")
testExplainWithLimit(db)

	// Test 10: EXPLAIN with Aggregate Functions
	log.Println("\n=== Test 10: EXPLAIN with Aggregate Functions ===")
testExplainWithAggregate(db)

	// Test 11: Cleanup
	log.Println("\n=== Test 11: Cleanup ===")
testCleanup(db)

	log.Println("\n=== All Phase 25 tests completed! ===")
	log.Println("🎉 Phase 25: EXPLAIN Query Plan Analysis - COMPLETE! 🎉")
}

func testBasicExplain(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE explain_test (id INTEGER, name TEXT, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: explain_test")

	// Insert test data
	for i := 1; i <= 10; i++ {
		_, err = db.Exec("INSERT INTO explain_test VALUES (?, ?, ?)", i, fmt.Sprintf("Item %d", i), i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 10 rows")

	// EXPLAIN SELECT
	log.Println("✓ EXPLAIN SELECT * FROM explain_test:")
	rows, err := db.Query("EXPLAIN SELECT * FROM explain_test")
	if err != nil {
		log.Printf("Error executing EXPLAIN: %v", err)
		return
	}
	defer rows.Close()

	// Display EXPLAIN results
	for rows.Next() {
		var addr int
		var opcode string
		var p1 int
		var p2 int
		var p3 int
		var p4 interface{}
		var comment sql.NullString

		err := rows.Scan(&addr, &opcode, &p1, &p2, &p3, &p4, &comment)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		commentStr := ""
		if comment.Valid {
			commentStr = fmt.Sprintf(" ; %s", comment.String)
		}

		log.Printf("  %d | %s %d %d %d %d%s", addr, opcode, p1, p2, p3, p4, commentStr)
	}
}

func testExplainWithIndex(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE explain_index (id INTEGER PRIMARY KEY, name TEXT, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: explain_index (with PRIMARY KEY)")

	// Create index
	_, err = db.Exec("CREATE INDEX idx_value ON explain_index(value)")
	if err != nil {
		log.Printf("Error creating index: %v", err)
		return
	}
	log.Println("✓ Created index: idx_value on value column")

	// Insert test data
	for i := 1; i <= 100; i++ {
		_, err = db.Exec("INSERT INTO explain_index VALUES (?, ?, ?)", i, fmt.Sprintf("Item %d", i), i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 100 rows")

	// EXPLAIN SELECT with WHERE clause
	log.Println("✓ EXPLAIN SELECT * FROM explain_index WHERE value = 50:")
	rows, err := db.Query("EXPLAIN SELECT * FROM explain_index WHERE value = 50")
	if err != nil {
		log.Printf("Error executing EXPLAIN: %v", err)
		return
	}
	defer rows.Close()

	// Display EXPLAIN results
	for rows.Next() {
		var addr int
		var opcode string
		var p1 int
		var p2 int
		var p3 int
		var p4 interface{}
		var comment sql.NullString

		err := rows.Scan(&addr, &opcode, &p1, &p2, &p3, &p4, &comment)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		commentStr := ""
		if comment.Valid {
			commentStr = fmt.Sprintf(" ; %s", comment.String)
		}

		log.Printf("  %d | %s %d %d %d %d%s", addr, opcode, p1, p2, p3, p4, commentStr)
	}
}

func testExplainWithJoin(db *sql.DB) {
	// Create test tables
	_, err := db.Exec("CREATE TABLE explain_users (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	_, err = db.Exec("CREATE TABLE explain_orders (id INTEGER PRIMARY KEY, user_id INTEGER, total REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created tables: explain_users, explain_orders")

	// Insert test data
	for i := 1; i <= 10; i++ {
		_, err = db.Exec("INSERT INTO explain_users VALUES (?, ?)", i, fmt.Sprintf("User %d", i))
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	for i := 1; i <= 50; i++ {
		userID := (i % 10) + 1
		_, err = db.Exec("INSERT INTO explain_orders VALUES (?, ?, ?)", i, userID, float64(i)*100)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted test data")

	// EXPLAIN SELECT with JOIN
	log.Println("✓ EXPLAIN SELECT u.name, o.total FROM explain_users u JOIN explain_orders o ON u.id = o.user_id:")
	rows, err := db.Query("EXPLAIN SELECT u.name, o.total FROM explain_users u JOIN explain_orders o ON u.id = o.user_id")
	if err != nil {
		log.Printf("Error executing EXPLAIN: %v", err)
		return
	}
	defer rows.Close()

	// Display EXPLAIN results
	for rows.Next() {
		var addr int
		var opcode string
		var p1 int
		var p2 int
		var p3 int
		var p4 interface{}
		var comment sql.NullString

		err := rows.Scan(&addr, &opcode, &p1, &p2, &p3, &p4, &comment)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		commentStr := ""
		if comment.Valid {
			commentStr = fmt.Sprintf(" ; %s", comment.String)
		}

		log.Printf("  %d | %s %d %d %d %d%s", addr, opcode, p1, p2, p3, p4, commentStr)
	}
}

func testExplainWithSubquery(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE explain_subquery (id INTEGER PRIMARY KEY, category TEXT, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: explain_subquery")

	// Insert test data
	for i := 1; i <= 20; i++ {
		category := fmt.Sprintf("Category %d", (i%5)+1)
		_, err = db.Exec("INSERT INTO explain_subquery VALUES (?, ?, ?)", i, category, i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 20 rows")

	// EXPLAIN SELECT with subquery
	log.Println("✓ EXPLAIN SELECT * FROM explain_subquery WHERE value > (SELECT AVG(value) FROM explain_subquery):")
	rows, err := db.Query("EXPLAIN SELECT * FROM explain_subquery WHERE value > (SELECT AVG(value) FROM explain_subquery)")
	if err != nil {
		log.Printf("Error executing EXPLAIN: %v", err)
		return
	}
	defer rows.Close()

	// Display EXPLAIN results
	for rows.Next() {
		var addr int
		var opcode string
		var p1 int
		var p2 int
		var p3 int
		var p4 interface{}
		var comment sql.NullString

		err := rows.Scan(&addr, &opcode, &p1, &p2, &p3, &p4, &comment)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		commentStr := ""
		if comment.Valid {
			commentStr = fmt.Sprintf(" ; %s", comment.String)
		}

		log.Printf("  %d | %s %d %d %d %d%s", addr, opcode, p1, p2, p3, p4, commentStr)
	}
}

func testExplainWithGroupBy(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE explain_groupby (id INTEGER, category TEXT, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: explain_groupby")

	// Insert test data
	for i := 1; i <= 50; i++ {
		category := fmt.Sprintf("Category %d", (i%5)+1)
		_, err = db.Exec("INSERT INTO explain_groupby VALUES (?, ?, ?)", i, category, i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 50 rows")

	// EXPLAIN SELECT with GROUP BY
	log.Println("✓ EXPLAIN SELECT category, COUNT(*), SUM(value) FROM explain_groupby GROUP BY category:")
	rows, err := db.Query("EXPLAIN SELECT category, COUNT(*), SUM(value) FROM explain_groupby GROUP BY category")
	if err != nil {
		log.Printf("Error executing EXPLAIN: %v", err)
		return
	}
	defer rows.Close()

	// Display EXPLAIN results
	for rows.Next() {
		var addr int
		var opcode string
		var p1 int
		var p2 int
		var p3 int
		var p4 interface{}
		var comment sql.NullString

		err := rows.Scan(&addr, &opcode, &p1, &p2, &p3, &p4, &comment)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		commentStr := ""
		if comment.Valid {
			commentStr = fmt.Sprintf(" ; %s", comment.String)
		}

		log.Printf("  %d | %s %d %d %d %d%s", addr, opcode, p1, p2, p3, p4, commentStr)
	}
}

func testExplainWithOrderBy(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE explain_orderby (id INTEGER, name TEXT, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: explain_orderby")

	// Insert test data
	for i := 1; i <= 30; i++ {
		_, err = db.Exec("INSERT INTO explain_orderby VALUES (?, ?, ?)", i, fmt.Sprintf("Item %d", i), i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 30 rows")

	// EXPLAIN SELECT with ORDER BY
	log.Println("✓ EXPLAIN SELECT * FROM explain_orderby ORDER BY value DESC:")
	rows, err := db.Query("EXPLAIN SELECT * FROM explain_orderby ORDER BY value DESC")
	if err != nil {
		log.Printf("Error executing EXPLAIN: %v", err)
		return
	}
	defer rows.Close()

	// Display EXPLAIN results
	for rows.Next() {
		var addr int
		var opcode string
		var p1 int
		var p2 int
		var p3 int
		var p4 interface{}
		var comment sql.NullString

		err := rows.Scan(&addr, &opcode, &p1, &p2, &p3, &p4, &comment)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		commentStr := ""
		if comment.Valid {
			commentStr = fmt.Sprintf(" ; %s", comment.String)
		}

		log.Printf("  %d | %s %d %d %d %d%s", addr, opcode, p1, p2, p3, p4, commentStr)
	}
}

func testExplainQueryPlan(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE explain_plan (id INTEGER PRIMARY KEY, name TEXT, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: explain_plan")

	// Create index
	_, err = db.Exec("CREATE INDEX idx_explain_plan_value ON explain_plan(value)")
	if err != nil {
		log.Printf("Error creating index: %v", err)
		return
	}
	log.Println("✓ Created index: idx_explain_plan_value")

	// Insert test data
	for i := 1; i <= 50; i++ {
		_, err = db.Exec("INSERT INTO explain_plan VALUES (?, ?, ?)", i, fmt.Sprintf("Item %d", i), i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 50 rows")

	// EXPLAIN QUERY PLAN SELECT
	log.Println("✓ EXPLAIN QUERY PLAN SELECT * FROM explain_plan WHERE value = 50:")
	rows, err := db.Query("EXPLAIN QUERY PLAN SELECT * FROM explain_plan WHERE value = 50")
	if err != nil {
		log.Printf("Error executing EXPLAIN QUERY PLAN: %v", err)
		return
	}
	defer rows.Close()

	// Display EXPLAIN QUERY PLAN results
	for rows.Next() {
		var id int
		var parent int
		var notused int
		var detail sql.NullString

		err := rows.Scan(&id, &parent, &notused, &detail)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		detailStr := "N/A"
		if detail.Valid {
			detailStr = detail.String
		}

		log.Printf("  id=%d, parent=%d, detail=%s", id, parent, detailStr)
	}
}

func testExplainWithWhere(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE explain_where (id INTEGER, name TEXT, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: explain_where")

	// Insert test data
	for i := 1; i <= 20; i++ {
		_, err = db.Exec("INSERT INTO explain_where VALUES (?, ?, ?)", i, fmt.Sprintf("Item %d", i), i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 20 rows")

	// EXPLAIN SELECT with WHERE clause
	log.Println("✓ EXPLAIN SELECT * FROM explain_where WHERE value > 50:")
	rows, err := db.Query("EXPLAIN SELECT * FROM explain_where WHERE value > 50")
	if err != nil {
		log.Printf("Error executing EXPLAIN: %v", err)
		return
	}
	defer rows.Close()

	// Display EXPLAIN results
	for rows.Next() {
		var addr int
		var opcode string
		var p1 int
		var p2 int
		var p3 int
		var p4 interface{}
		var comment sql.NullString

		err := rows.Scan(&addr, &opcode, &p1, &p2, &p3, &p4, &comment)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		commentStr := ""
		if comment.Valid {
			commentStr = fmt.Sprintf(" ; %s", comment.String)
		}

		log.Printf("  %d | %s %d %d %d %d%s", addr, opcode, p1, p2, p3, p4, commentStr)
	}
}

func testExplainWithLimit(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE explain_limit (id INTEGER, name TEXT, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: explain_limit")

	// Insert test data
	for i := 1; i <= 100; i++ {
		_, err = db.Exec("INSERT INTO explain_limit VALUES (?, ?, ?)", i, fmt.Sprintf("Item %d", i), i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 100 rows")

	// EXPLAIN SELECT with LIMIT
	log.Println("✓ EXPLAIN SELECT * FROM explain_limit LIMIT 10:")
	rows, err := db.Query("EXPLAIN SELECT * FROM explain_limit LIMIT 10")
	if err != nil {
		log.Printf("Error executing EXPLAIN: %v", err)
		return
	}
	defer rows.Close()

	// Display EXPLAIN results
	for rows.Next() {
		var addr int
		var opcode string
		var p1 int
		var p2 int
		var p3 int
		var p4 interface{}
		var comment sql.NullString

		err := rows.Scan(&addr, &opcode, &p1, &p2, &p3, &p4, &comment)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		commentStr := ""
		if comment.Valid {
			commentStr = fmt.Sprintf(" ; %s", comment.String)
		}

		log.Printf("  | %s %d %d %d %d%s", opcode, p1, p2, p3, p4, commentStr)
	}
}

func testExplainWithAggregate(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE explain_agg (id INTEGER, category TEXT, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: explain_agg")

	// Insert test data
	for i := 1; i <= 50; i++ {
		category := fmt.Sprintf("Category %d", (i%5)+1)
		_, err = db.Exec("INSERT INTO explain_agg VALUES (?, ?, ?)", i, category, i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 50 rows")

	// EXPLAIN SELECT with aggregate functions
	log.Println("✓ EXPLAIN SELECT COUNT(*), SUM(value), AVG(value), MIN(value), MAX(value) FROM explain_agg:")
	rows, err := db.Query("EXPLAIN SELECT COUNT(*), SUM(value), AVG(value), MIN(value), MAX(value) FROM explain_agg")
	if err != nil {
		log.Printf("Error executing EXPLAIN: %v", err)
		return
	}
	defer rows.Close()

	// Display EXPLAIN results
	for rows.Next() {
		var addr int
		var opcode string
		var p1 int
		var p2 int
		var p3 int
		var p4 interface{}
		var comment sql.NullString

		err := rows.Scan(&addr, &opcode, &p1, &p2, &p3, &p4, &comment)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		commentStr := ""
		if comment.Valid {
			commentStr = fmt.Sprintf(" ; %s", comment.String)
		}

		log.Printf("  %d | %s %d %d %d %d%s", addr, opcode, p1, p2, p3, p4, commentStr)
	}
}

func testCleanup(db *sql.DB) {
	tables := []string{
		"explain_test",
		"explain_index",
		"explain_users",
		"explain_orders",
		"explain_subquery",
		"explain_groupby",
		"explain_orderby",
		"explain_plan",
		"explain_where",
		"explain_limit",
		"explain_agg",
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
