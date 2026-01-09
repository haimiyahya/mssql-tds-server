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

	// Test 1: ROW_NUMBER()
	log.Println("\n=== Test 1: ROW_NUMBER() ===")
testRowNumber(db)

	// Test 2: RANK()
	log.Println("\n=== Test 2: RANK() ===")
testRank(db)

	// Test 3: DENSE_RANK()
	log.Println("\n=== Test 3: DENSE_RANK() ===")
testDenseRank(db)

	// Test 4: SUM() OVER (Running Total)
	log.Println("\n=== Test 4: SUM() OVER (Running Total) ===")
testSumOver(db)

	// Test 5: AVG() OVER (Moving Average)
	log.Println("\n=== Test 5: AVG() OVER (Moving Average) ===")
testAvgOver(db)

	// Test 6: COUNT() OVER
	log.Println("\n=== Test 6: COUNT() OVER ===")
testCountOver(db)

	// Test 7: MIN() OVER and MAX() OVER
	log.Println("\n=== Test 7: MIN() OVER and MAX() OVER ===")
testMinMaxOver(db)

	// Test 8: PARTITION BY
	log.Println("\n=== Test 8: PARTITION BY ===")
testPartitionBy(db)

	// Test 9: ROWS BETWEEN (Frame Clause)
	log.Println("\n=== Test 9: ROWS BETWEEN (Frame Clause) ===")
testRowsBetween(db)

	// Test 10: RANGE BETWEEN (Frame Clause)
	log.Println("\n=== Test 10: RANGE BETWEEN (Frame Clause) ===")
testRangeBetween(db)

	// Test 11: Multiple Window Functions
	log.Println("\n=== Test 11: Multiple Window Functions ===")
testMultipleWindowFunctions(db)

	// Test 12: Window Functions with Aggregates
	log.Println("\n=== Test 12: Window Functions with Aggregates ===")
testWindowWithAggregates(db)

	// Test 13: Cleanup
	log.Println("\n=== Test 13: Cleanup ===")
testCleanup(db)

	log.Println("\n=== All Phase 28 tests completed! ===")
	log.Println("🎉 Phase 28: Window Functions - COMPLETE! 🎉")
}

func testRowNumber(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE employees (id INTEGER, name TEXT, department TEXT, salary INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: employees")

	// Insert test data
	employees := []struct {
		id         int
		name       string
		department string
		salary     int
	}{
		{1, "Alice", "Sales", 50000},
		{2, "Bob", "Sales", 55000},
		{3, "Charlie", "Engineering", 60000},
		{4, "Diana", "Engineering", 65000},
		{5, "Eve", "Engineering", 70000},
	}

	for _, emp := range employees {
		_, err = db.Exec("INSERT INTO employees VALUES (?, ?, ?, ?)", emp.id, emp.name, emp.department, emp.salary)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 5 employees")

	// ROW_NUMBER()
	log.Println("✓ ROW_NUMBER() - Employee Ranking:")
	rows, err := db.Query(`
		SELECT id, name, department, salary, 
		       ROW_NUMBER() OVER (ORDER BY salary DESC) as row_num
		FROM employees
	`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var department string
		var salary int
		var rowNum int
		err := rows.Scan(&id, &name, &department, &salary, &rowNum)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - %s - $%d (Row Num: %d)", id, name, department, salary, rowNum)
	}
}

func testRank(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE scores (id INTEGER, player TEXT, score INTEGER, game_date TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: scores")

	// Insert test data (with ties)
	scores := []struct {
		id        int
		player    string
		score     int
		gameDate  string
	}{
		{1, "Alice", 100, "2024-01-01"},
		{2, "Bob", 90, "2024-01-01"},
		{3, "Charlie", 100, "2024-01-02"},
		{4, "Diana", 85, "2024-01-02"},
		{5, "Eve", 95, "2024-01-03"},
	}

	for _, s := range scores {
		_, err = db.Exec("INSERT INTO scores VALUES (?, ?, ?, ?)", s.id, s.player, s.score, s.gameDate)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 5 scores")

	// RANK()
	log.Println("✓ RANK() - Score Ranking (with ties):")
	rows, err := db.Query(`
		SELECT id, player, score, 
		       RANK() OVER (ORDER BY score DESC) as rank
		FROM scores
	`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var player string
		var score int
		var rank int
		err := rows.Scan(&id, &player, &score, &rank)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - Score: %d (Rank: %d)", id, player, score, rank)
	}
}

func testDenseRank(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE products (id INTEGER, name TEXT, price REAL, category TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: products")

	// Insert test data (with ties)
	products := []struct {
		id       int
		name     string
		price    float64
		category string
	}{
		{1, "Product1", 10.00, "Electronics"},
		{2, "Product2", 20.00, "Electronics"},
		{3, "Product3", 20.00, "Electronics"},
		{4, "Product4", 30.00, "Electronics"},
		{5, "Product5", 30.00, "Electronics"},
	}

	for _, p := range products {
		_, err = db.Exec("INSERT INTO products VALUES (?, ?, ?, ?)", p.id, p.name, p.price, p.category)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 5 products")

	// DENSE_RANK()
	log.Println("✓ DENSE_RANK() - Price Ranking (without gaps):")
	rows, err := db.Query(`
		SELECT id, name, price, 
		       DENSE_RANK() OVER (ORDER BY price ASC) as dense_rank
		FROM products
	`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var price float64
		var denseRank int
		err := rows.Scan(&id, &name, &price, &denseRank)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - Price: $%.2f (Dense Rank: %d)", id, name, price, denseRank)
	}
}

func testSumOver(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE sales (id INTEGER, sales_date TEXT, amount REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: sales")

	// Insert test data
	sales := []struct {
		id        int
		salesDate string
		amount    float64
	}{
		{1, "2024-01-01", 1000},
		{2, "2024-01-02", 1500},
		{3, "2024-01-03", 2000},
		{4, "2024-01-04", 1200},
		{5, "2024-01-05", 1800},
	}

	for _, s := range sales {
		_, err = db.Exec("INSERT INTO sales VALUES (?, ?, ?)", s.id, s.salesDate, s.amount)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 5 sales")

	// SUM() OVER (Running Total)
	log.Println("✓ SUM() OVER - Running Total:")
	rows, err := db.Query(`
		SELECT id, sales_date, amount,
		       SUM(amount) OVER (ORDER BY id) as running_total
		FROM sales
	`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var salesDate string
		var amount float64
		var runningTotal float64
		err := rows.Scan(&id, &salesDate, &amount, &runningTotal)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - $%.2f (Running Total: $%.2f)", id, salesDate, amount, runningTotal)
	}
}

func testAvgOver(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE stock_prices (id INTEGER, price_date TEXT, price REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: stock_prices")

	// Insert test data
	prices := []struct {
		id         int
		priceDate  string
		price      float64
	}{
		{1, "2024-01-01", 100.00},
		{2, "2024-01-02", 105.00},
		{3, "2024-01-03", 110.00},
		{4, "2024-01-04", 108.00},
		{5, "2024-01-05", 112.00},
	}

	for _, p := range prices {
		_, err = db.Exec("INSERT INTO stock_prices VALUES (?, ?, ?)", p.id, p.priceDate, p.price)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 5 stock prices")

	// AVG() OVER (Moving Average - 3-day)
	log.Println("✓ AVG() OVER - 3-Day Moving Average:")
	rows, err := db.Query(`
		SELECT id, price_date, price,
		       AVG(price) OVER (
		         ORDER BY id 
		         ROWS BETWEEN 2 PRECEDING AND CURRENT ROW
		       ) as moving_avg
		FROM stock_prices
	`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var priceDate string
		var price float64
		var movingAvg float64
		err := rows.Scan(&id, &priceDate, &price, &movingAvg)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - $%.2f (3-Day Avg: $%.2f)", id, priceDate, price, movingAvg)
	}
}

func testCountOver(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE visits (id INTEGER, visit_date TEXT, page_views INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: visits")

	// Insert test data
	visits := []struct {
		id         int
		visitDate  string
		pageViews  int
	}{
		{1, "2024-01-01", 100},
		{2, "2024-01-02", 150},
		{3, "2024-01-03", 200},
		{4, "2024-01-04", 120},
		{5, "2024-01-05", 180},
	}

	for _, v := range visits {
		_, err = db.Exec("INSERT INTO visits VALUES (?, ?, ?)", v.id, v.visitDate, v.pageViews)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 5 visits")

	// COUNT() OVER
	log.Println("✓ COUNT() OVER - Cumulative Count:")
	rows, err := db.Query(`
		SELECT id, visit_date, page_views,
		       COUNT(*) OVER (ORDER BY id) as cumulative_count
		FROM visits
	`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var visitDate string
		var pageViews int
		var cumulativeCount int
		err := rows.Scan(&id, &visitDate, &pageViews, &cumulativeCount)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - %d views (Cumulative Count: %d)", id, visitDate, pageViews, cumulativeCount)
	}
}

func testMinMaxOver(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE temperatures (id INTEGER, date TEXT, temperature REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: temperatures")

	// Insert test data
	temps := []struct {
		id           int
		date         string
		temperature  float64
	}{
		{1, "2024-01-01", 20.0},
		{2, "2024-01-02", 22.0},
		{3, "2024-01-03", 18.0},
		{4, "2024-01-04", 25.0},
		{5, "2024-01-05", 21.0},
	}

	for _, t := range temps {
		_, err = db.Exec("INSERT INTO temperatures VALUES (?, ?, ?)", t.id, t.date, t.temperature)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 5 temperatures")

	// MIN() OVER and MAX() OVER
	log.Println("✓ MIN() OVER and MAX() OVER - Temperature Range:")
	rows, err := db.Query(`
		SELECT id, date, temperature,
		       MIN(temperature) OVER (
		         ORDER BY id 
		         ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
		       ) as min_temp,
		       MAX(temperature) OVER (
		         ORDER BY id 
		         ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
		       ) as max_temp
		FROM temperatures
	`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var date string
		var temperature float64
		var minTemp float64
		var maxTemp float64
		err := rows.Scan(&id, &date, &temperature, &minTemp, &maxTemp)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - %.1f°C (Min: %.1f°C, Max: %.1f°C)", id, date, temperature, minTemp, maxTemp)
	}
}

func testPartitionBy(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE orders (id INTEGER, customer_id INTEGER, order_date TEXT, total REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: orders")

	// Insert test data
	orders := []struct {
		id         int
		customerID int
		orderDate  string
		total      float64
	}{
		{1, 1, "2024-01-01", 100.00},
		{2, 1, "2024-01-02", 150.00},
		{3, 2, "2024-01-01", 200.00},
		{4, 2, "2024-01-02", 250.00},
		{5, 3, "2024-01-01", 300.00},
	}

	for _, o := range orders {
		_, err = db.Exec("INSERT INTO orders VALUES (?, ?, ?, ?)", o.id, o.customerID, o.orderDate, o.total)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 5 orders")

	// PARTITION BY
	log.Println("✓ PARTITION BY - Customer Order Ranking:")
	rows, err := db.Query(`
		SELECT id, customer_id, order_date, total,
		       ROW_NUMBER() OVER (PARTITION BY customer_id ORDER BY order_date) as order_rank
		FROM orders
	`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var customerID int
		var orderDate string
		var total float64
		var orderRank int
		err := rows.Scan(&id, &customerID, &orderDate, &total, &orderRank)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: Customer %d - %s - $%.2f (Order Rank: %d)", id, customerID, orderDate, total, orderRank)
	}
}

func testRowsBetween(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE monthly_sales (id INTEGER, month TEXT, sales REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: monthly_sales")

	// Insert test data
	sales := []struct {
		id    int
		month string
		sales float64
	}{
		{1, "2024-01", 10000},
		{2, "2024-02", 15000},
		{3, "2024-03", 20000},
		{4, "2024-04", 12000},
		{5, "2024-05", 18000},
	}

	for _, s := range sales {
		_, err = db.Exec("INSERT INTO monthly_sales VALUES (?, ?, ?)", s.id, s.month, s.sales)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 5 months of sales")

	// ROWS BETWEEN (3-month moving average)
	log.Println("✓ ROWS BETWEEN - 3-Month Moving Average:")
	rows, err := db.Query(`
		SELECT id, month, sales,
		       AVG(sales) OVER (
		         ORDER BY id 
		         ROWS BETWEEN 2 PRECEDING AND CURRENT ROW
		       ) as moving_avg,
		       SUM(sales) OVER (
		         ORDER BY id 
		         ROWS BETWEEN 2 PRECEDING AND CURRENT ROW
		       ) as moving_sum
		FROM monthly_sales
	`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var month string
		var sales float64
		var movingAvg float64
		var movingSum float64
		err := rows.Scan(&id, &month, &sales, &movingAvg, &movingSum)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - $%.0f (3-Month Avg: $%.0f, 3-Month Sum: $%.0f)", id, month, sales, movingAvg, movingSum)
	}
}

func testRangeBetween(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE student_scores (id INTEGER, student TEXT, score INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: student_scores")

	// Insert test data
	scores := []struct {
		id      int
		student string
		score   int
	}{
		{1, "Alice", 95},
		{2, "Bob", 88},
		{3, "Charlie", 92},
		{4, "Diana", 90},
		{5, "Eve", 88},
	}

	for _, s := range scores {
		_, err = db.Exec("INSERT INTO student_scores VALUES (?, ?, ?)", s.id, s.student, s.score)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 5 student scores")

	// RANGE BETWEEN (scores within 5 points)
	log.Println("✓ RANGE BETWEEN - Students with similar scores (±5 points):")
	rows, err := db.Query(`
		SELECT id, student, score,
		       COUNT(*) OVER (
		         ORDER BY score
		         RANGE BETWEEN 5 PRECEDING AND 5 FOLLOWING
		       ) as similar_count
		FROM student_scores
	`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var student string
		var score int
		var similarCount int
		err := rows.Scan(&id, &student, &score, &similarCount)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - Score: %d (Similar count: %d)", id, student, score, similarCount)
	}
}

func testMultipleWindowFunctions(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE revenue (id INTEGER, year INTEGER, quarter INTEGER, amount REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: revenue")

	// Insert test data
	revenue := []struct {
		id     int
		year   int
		quarter int
		amount float64
	}{
		{1, 2024, 1, 10000},
		{2, 2024, 2, 12000},
		{3, 2024, 3, 15000},
		{4, 2024, 4, 18000},
		{5, 2025, 1, 11000},
	}

	for _, r := range revenue {
		_, err = db.Exec("INSERT INTO revenue VALUES (?, ?, ?, ?)", r.id, r.year, r.quarter, r.amount)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 5 revenue records")

	// Multiple window functions
	log.Println("✓ Multiple Window Functions - Revenue Analysis:")
	rows, err := db.Query(`
		SELECT id, year, quarter, amount,
		       ROW_NUMBER() OVER (ORDER BY amount) as row_num,
		       RANK() OVER (ORDER BY amount) as rank_num,
		       SUM(amount) OVER (PARTITION BY year ORDER BY quarter) as ytd,
		       AVG(amount) OVER (PARTITION BY year) as avg_quarter
		FROM revenue
	`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var year int
		quarter := 0
		var amount float64
		var rowNum int
		var rankNum int
		var ytd float64
		var avgQuarter float64
		err := rows.Scan(&id, &year, &quarter, &amount, &rowNum, &rankNum, &ytd, &avgQuarter)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %d-Q%d - $%.0f (Row: %d, Rank: %d, YTD: $%.0f, Avg: $%.0f)",
			id, year, quarter, amount, rowNum, rankNum, ytd, avgQuarter)
	}
}

func testWindowWithAggregates(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE sales_summary (id INTEGER, region TEXT, product TEXT, sales REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: sales_summary")

	// Insert test data
	sales := []struct {
		id      int
		region  string
		product string
		sales   float64
	}{
		{1, "East", "A", 100},
		{2, "East", "B", 150},
		{3, "West", "A", 200},
		{4, "West", "B", 250},
		{5, "East", "A", 120},
	}

	for _, s := range sales {
		_, err = db.Exec("INSERT INTO sales_summary VALUES (?, ?, ?, ?)", s.id, s.region, s.product, s.sales)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 5 sales records")

	// Window functions with GROUP BY and HAVING
	log.Println("✓ Window Functions with Aggregates - Regional Sales:")
	rows, err := db.Query(`
		SELECT region, product, sales,
		       SUM(sales) OVER (PARTITION BY region) as region_total,
		       sales * 100.0 / SUM(sales) OVER (PARTITION BY region) as pct_of_region
		FROM sales_summary
		WHERE sales > 100
		ORDER BY region, sales DESC
	`)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var region string
		var product string
		var sales float64
		var regionTotal float64
		var pctOfRegion float64
		err := rows.Scan(&region, &product, &sales, &regionTotal, &pctOfRegion)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %s - %s - $%.0f (Region Total: $%.0f, %.1f%%)", region, product, sales, regionTotal, pctOfRegion)
	}
}

func testCleanup(db *sql.DB) {
	tables := []string{
		"employees",
		"scores",
		"products",
		"sales",
		"stock_prices",
		"visits",
		"temperatures",
		"orders",
		"monthly_sales",
		"student_scores",
		"revenue",
		"sales_summary",
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
