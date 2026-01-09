package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
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

// Performance metrics
type PerformanceMetrics struct {
	QueryCount        int
	TotalDuration     time.Duration
	AverageDuration   time.Duration
	MinDuration      time.Duration
	MaxDuration      time.Duration
	SlowQueryCount   int
	SlowQueryThreshold time.Duration
}

// Connection pool metrics
type PoolMetrics struct {
	OpenConnections int
	InUse          int
	Idle           int
	WaitCount      int64
	WaitDuration   time.Duration
}

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

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging server: %v", err)
	}

	log.Println("Successfully connected to TDS server!")

	// Test 1: Query Performance Metrics
	log.Println("\n=== Test 1: Query Performance Metrics ===")
	testQueryPerformanceMetrics(db)

	// Test 2: Connection Pool Monitoring
	log.Println("\n=== Test 2: Connection Pool Monitoring ===")
	testConnectionPoolMonitoring(db)

	// Test 3: Slow Query Detection
	log.Println("\n=== Test 3: Slow Query Detection ===")
	testSlowQueryDetection(db)

	// Test 4: Concurrent Query Performance
	log.Println("\n=== Test 4: Concurrent Query Performance ===")
	testConcurrentQueryPerformance(db)

	// Test 5: Performance Reporting
	log.Println("\n=== Test 5: Performance Reporting ===")
	testPerformanceReporting(db)

	// Test 6: Resource Usage Tracking
	log.Println("\n=== Test 6: Resource Usage Tracking ===")
	testResourceUsageTracking(db)

	// Test 7: Performance Under Load
	log.Println("\n=== Test 7: Performance Under Load ===")
	testPerformanceUnderLoad(db)

	// Test 8: Connection Pool Performance
	log.Println("\n=== Test 8: Connection Pool Performance ===")
	testConnectionPoolPerformance(db)

	// Test 9: Batch Operation Performance
	log.Println("\n=== Test 9: Batch Operation Performance ===")
	testBatchOperationPerformance(db)

	// Test 10: Cleanup
	log.Println("\n=== Test 10: Cleanup ===")
	testCleanup(db)

	log.Println("\n=== All Phase 22 tests completed! ===")
	log.Println("🎉 Phase 22: Performance Monitoring - COMPLETE! 🎉")
}

func testQueryPerformanceMetrics(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE perf_test (id INTEGER, name TEXT, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: perf_test")

	// Insert test data
	for i := 1; i <= 100; i++ {
		_, err = db.Exec("INSERT INTO perf_test VALUES (?, ?, ?)", i, fmt.Sprintf("Item %d", i), i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 100 rows")

	// Measure query performance
	metrics := &PerformanceMetrics{
		SlowQueryThreshold: 10 * time.Millisecond,
	}

	// Execute queries and measure performance
	for i := 1; i <= 10; i++ {
		start := time.Now()
		rows, err := db.Query("SELECT * FROM perf_test WHERE value > ?", i*100)
		duration := time.Since(start)

		metrics.QueryCount++
		metrics.TotalDuration += duration

		if metrics.MinDuration == 0 || duration < metrics.MinDuration {
			metrics.MinDuration = duration
		}
		if duration > metrics.MaxDuration {
			metrics.MaxDuration = duration
		}

		if duration > metrics.SlowQueryThreshold {
			metrics.SlowQueryCount++
		}

		if err != nil {
			log.Printf("Error executing query: %v", err)
			return
		}
		rows.Close()
	}

	// Calculate average
	metrics.AverageDuration = metrics.TotalDuration / time.Duration(metrics.QueryCount)

	// Display metrics
	log.Println("✓ Query Performance Metrics:")
	log.Printf("  Total Queries: %d", metrics.QueryCount)
	log.Printf("  Total Duration: %v", metrics.TotalDuration)
	log.Printf("  Average Duration: %v", metrics.AverageDuration)
	log.Printf("  Min Duration: %v", metrics.MinDuration)
	log.Printf("  Max Duration: %v", metrics.MaxDuration)
	log.Printf("  Slow Queries (> %v): %d", metrics.SlowQueryThreshold, metrics.SlowQueryCount)
}

func testConnectionPoolMonitoring(db *sql.DB) {
	// Get initial pool statistics
	initialStats := db.Stats()
	log.Println("✓ Initial Pool Statistics:")
	log.Printf("  Open Connections: %d", initialStats.OpenConnections)
	log.Printf("  In Use: %d", initialStats.InUse)
	log.Printf("  Idle: %d", initialStats.Idle)
	log.Printf("  Max Open Connections: %d", initialStats.MaxOpenConnections)
	log.Printf("  Wait Count: %d", initialStats.WaitCount)
	log.Printf("  Wait Duration: %v", initialStats.WaitDuration)

	// Execute queries to monitor pool
	for i := 1; i <= 5; i++ {
		var result int
		start := time.Now()
		db.QueryRow("SELECT ? * 2", i).Scan(&result)
		duration := time.Since(start)
		log.Printf("  Query %d: Duration=%v, Result=%d", i, duration, result)
	}

	// Get final pool statistics
	finalStats := db.Stats()
	log.Println("✓ Final Pool Statistics:")
	log.Printf("  Open Connections: %d", finalStats.OpenConnections)
	log.Printf("  In Use: %d", finalStats.InUse)
	log.Printf("  Idle: %d", finalStats.Idle)
	log.Printf("  Wait Count: %d", finalStats.WaitCount)
	log.Printf("  Wait Duration: %v", finalStats.WaitDuration)
	log.Printf("  Max Idle Closed: %d", finalStats.MaxIdleClosed)
	log.Printf("  Max Lifetime Closed: %d", finalStats.MaxLifetimeClosed)
}

func testSlowQueryDetection(db *sql.DB) {
	// Create test table with data
	_, err := db.Exec("CREATE TABLE slow_test (id INTEGER, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: slow_test")

	// Insert test data
	for i := 1; i <= 1000; i++ {
		_, err = db.Exec("INSERT INTO slow_test VALUES (?, ?)", i, i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 1000 rows")

	// Define slow query threshold
	threshold := 5 * time.Millisecond
	log.Printf("✓ Slow Query Threshold: %v", threshold)

	// Execute queries and detect slow ones
	slowQueries := 0
	for i := 1; i <= 10; i++ {
		start := time.Now()
		rows, err := db.Query("SELECT * FROM slow_test WHERE value > ?", i*1000)
		duration := time.Since(start)

		if err != nil {
			log.Printf("Error executing query: %v", err)
			return
		}
		rows.Close()

		isSlow := duration > threshold
		if isSlow {
			slowQueries++
			log.Printf("  ⚠ Query %d: SLOW (%v > %v)", i, duration, threshold)
		} else {
			log.Printf("  ✓ Query %d: Fast (%v)", i, duration)
		}
	}

	log.Printf("✓ Slow Query Detection: %d slow queries detected out of 10", slowQueries)
}

func testConcurrentQueryPerformance(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE concurrent_test (id INTEGER, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: concurrent_test")

	// Insert test data
	for i := 1; i <= 100; i++ {
		_, err = db.Exec("INSERT INTO concurrent_test VALUES (?, ?)", i, i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 100 rows")

	// Execute concurrent queries
	var wg sync.WaitGroup
	numGoroutines := 10
	queriesPerGoroutine := 10
	totalQueries := numGoroutines * queriesPerGoroutine

	start := time.Now()

	for g := 1; g <= numGoroutines; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for q := 1; q <= queriesPerGoroutine; q++ {
				var count int
				db.QueryRow("SELECT COUNT(*) FROM concurrent_test WHERE value > ?", q*50).Scan(&count)
			}
		}(g)
	}

	wg.Wait()
	totalDuration := time.Since(start)
	avgDuration := totalDuration / time.Duration(totalQueries)
	queriesPerSecond := float64(totalQueries) / totalDuration.Seconds()

	log.Println("✓ Concurrent Query Performance:")
	log.Printf("  Total Goroutines: %d", numGoroutines)
	log.Printf("  Total Queries: %d", totalQueries)
	log.Printf("  Total Duration: %v", totalDuration)
	log.Printf("  Average Query Duration: %v", avgDuration)
	log.Printf("  Queries Per Second: %.2f", queriesPerSecond)

	// Display connection pool statistics
	stats := db.Stats()
	log.Printf("  Pool Stats: Open=%d, InUse=%d, Idle=%d", stats.OpenConnections, stats.InUse, stats.Idle)
}

func testPerformanceReporting(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE report_test (id INTEGER, category TEXT, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: report_test")

	// Insert test data
	for i := 1; i <= 50; i++ {
		category := fmt.Sprintf("Category %d", (i%5)+1)
		_, err = db.Exec("INSERT INTO report_test VALUES (?, ?, ?)", i, category, i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 50 rows")

	// Generate performance report
	log.Println("✓ Performance Report:")

	// Query 1: COUNT
	start := time.Now()
	var count int
	db.QueryRow("SELECT COUNT(*) FROM report_test").Scan(&count)
	duration1 := time.Since(start)
	log.Printf("  COUNT(*): %v (%d rows)", duration1, count)

	// Query 2: AGGREGATE
	start = time.Now()
	var total sql.NullInt64
	db.QueryRow("SELECT SUM(value) FROM report_test").Scan(&total)
	duration2 := time.Since(start)
	log.Printf("  SUM(value): %v (Total=%d)", duration2, total.Int64)

	// Query 3: GROUP BY
	start = time.Now()
	rows, err := db.Query("SELECT category, COUNT(*) FROM report_test GROUP BY category")
	duration3 := time.Since(start)
	if err != nil {
		log.Printf("Error with GROUP BY: %v", err)
	} else {
		log.Printf("  GROUP BY: %v", duration3)
		rows.Close()
	}

	// Query 4: SELECT with filter
	start = time.Now()
	rows, err = db.Query("SELECT * FROM report_test WHERE value > ?", 250)
	duration4 := time.Since(start)
	if err != nil {
		log.Printf("Error with SELECT filter: %v", err)
	} else {
		log.Printf("  SELECT with filter: %v", duration4)
		rows.Close()
	}

	// Summary
	totalDuration := duration1 + duration2 + duration3 + duration4
	avgDuration := totalDuration / 4
	log.Printf("  Summary: Total=%v, Average=%v", totalDuration, avgDuration)
}

func testResourceUsageTracking(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE resource_test (id INTEGER, data TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: resource_test")

	// Insert varying amounts of data
	dataSizes := []int{100, 500, 1000, 5000, 10000}
	for i, size := range dataSizes {
		data := string(make([]byte, size))
		start := time.Now()
		_, err = db.Exec("INSERT INTO resource_test VALUES (?, ?)", i+1, data)
		duration := time.Since(start)

		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}

		log.Printf("  Insert %d bytes: Duration=%v", size, duration)
	}

	// Query varying amounts of data
	start := time.Now()
	rows, err := db.Query("SELECT * FROM resource_test")
	duration := time.Since(start)
	if err != nil {
		log.Printf("Error querying data: %v", err)
		return
	}
	defer rows.Close()

	rowCount := 0
	totalBytes := 0
	for rows.Next() {
		var id int
		var data string
		rows.Scan(&id, &data)
		rowCount++
		totalBytes += len(data)
	}

	log.Printf("  Query %d rows (%d total bytes): Duration=%v", rowCount, totalBytes, duration)

	// Display connection pool statistics
	stats := db.Stats()
	log.Printf("  Pool Stats: Open=%d, InUse=%d, Idle=%d", stats.OpenConnections, stats.InUse, stats.Idle)
}

func testPerformanceUnderLoad(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE load_test (id INTEGER, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: load_test")

	// Insert test data
	for i := 1; i <= 1000; i++ {
		_, err = db.Exec("INSERT INTO load_test VALUES (?, ?)", i, i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 1000 rows")

	// Execute load test
	var wg sync.WaitGroup
	numGoroutines := 20
	queriesPerGoroutine := 20
	totalQueries := numGoroutines * queriesPerGoroutine

	log.Printf("✓ Starting Load Test: %d goroutines, %d queries/goroutine", numGoroutines, queriesPerGoroutine)

	start := time.Now()

	for g := 1; g <= numGoroutines; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for q := 1; q <= queriesPerGoroutine; q++ {
				var count int
				db.QueryRow("SELECT COUNT(*) FROM load_test WHERE value > ?", q*50).Scan(&count)
			}
		}(g)
	}

	wg.Wait()
	totalDuration := time.Since(start)
	avgDuration := totalDuration / time.Duration(totalQueries)
	queriesPerSecond := float64(totalQueries) / totalDuration.Seconds()

	log.Println("✓ Load Test Results:")
	log.Printf("  Total Goroutines: %d", numGoroutines)
	log.Printf("  Total Queries: %d", totalQueries)
	log.Printf("  Total Duration: %v", totalDuration)
	log.Printf("  Average Query Duration: %v", avgDuration)
	log.Printf("  Queries Per Second: %.2f", queriesPerSecond)

	// Display connection pool statistics
	stats := db.Stats()
	log.Printf("  Pool Stats: Open=%d, InUse=%d, Idle=%d", stats.OpenConnections, stats.InUse, stats.Idle)
	log.Printf("  Wait Count: %d", stats.WaitCount)
	log.Printf("  Wait Duration: %v", stats.WaitDuration)
}

func testConnectionPoolPerformance(db *sql.DB) {
	// Test different pool configurations
	configurations := []struct {
		maxOpen  int
		maxIdle  int
		name      string
	}{
		{5, 2, "Small Pool"},
		{10, 5, "Medium Pool"},
		{25, 10, "Large Pool"},
	}

	for _, config := range configurations {
		// Configure pool
		db.SetMaxOpenConns(config.maxOpen)
		db.SetMaxIdleConns(config.maxIdle)

		// Reset connection pool (close all connections)
		db.SetMaxIdleConns(0)
		db.SetMaxIdleConns(config.maxIdle)

		// Execute queries
		start := time.Now()
		numQueries := 50
		for i := 1; i <= numQueries; i++ {
			var result int
			db.QueryRow("SELECT ? * 2", i).Scan(&result)
		}
		duration := time.Since(start)

		// Get pool statistics
		stats := db.Stats()

		log.Printf("✓ %s (MaxOpen=%d, MaxIdle=%d):", config.name, config.maxOpen, config.maxIdle)
		log.Printf("  Duration: %v", duration)
		log.Printf("  Pool Stats: Open=%d, InUse=%d, Idle=%d", stats.OpenConnections, stats.InUse, stats.Idle)
	}
}

func testBatchOperationPerformance(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE batch_perf_test (id INTEGER, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: batch_perf_test")

	// Test individual INSERT vs batch INSERT
	numRows := 100

	// Individual INSERT
	start := time.Now()
	for i := 1; i <= numRows; i++ {
		_, err = db.Exec("INSERT INTO batch_perf_test VALUES (?, ?)", i, i*10)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			return
		}
	}
	individualDuration := time.Since(start)

	log.Printf("✓ Individual INSERT (%d rows): Duration=%v", numRows, individualDuration)

	// Clean table
	_, err = db.Exec("DELETE FROM batch_perf_test")
	if err != nil {
		log.Printf("Error deleting data: %v", err)
		return
	}

	// Batch INSERT (multiple rows in single statement)
	valuesSQL := ""
	for i := 1; i <= numRows; i++ {
		if i > 1 {
			valuesSQL += ", "
		}
		valuesSQL += fmt.Sprintf("(%d, %d)", i, i*10)
	}

	start = time.Now()
	_, err = db.Exec("INSERT INTO batch_perf_test VALUES " + valuesSQL)
	batchDuration := time.Since(start)
	if err != nil {
		log.Printf("Error with batch INSERT: %v", err)
		return
	}

	log.Printf("✓ Batch INSERT (%d rows): Duration=%v", numRows, batchDuration)

	// Calculate performance improvement
	speedup := float64(individualDuration) / float64(batchDuration)
	percentageImprovement := ((float64(individualDuration) - float64(batchDuration)) / float64(individualDuration)) * 100

	log.Printf("✓ Performance Comparison:")
	log.Printf("  Speedup: %.2fx", speedup)
	log.Printf("  Improvement: %.1f%%", percentageImprovement)
}

func testCleanup(db *sql.DB) {
	tables := []string{
		"perf_test",
		"slow_test",
		"concurrent_test",
		"report_test",
		"resource_test",
		"load_test",
		"batch_perf_test",
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
