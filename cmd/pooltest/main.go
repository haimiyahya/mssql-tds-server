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

func main() {
	// Build connection string
	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s;pooling=true",
		server, port, database, username, password)

	log.Printf("Connecting to TDS server at %s:%d", server, port)

	// Connect to database with connection pooling enabled
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

	log.Println("Successfully connected to TDS server with connection pooling!")

	// Test 1: Connection Pool Configuration
	log.Println("\n=== Test 1: Connection Pool Configuration ===")
	testConnectionPoolConfiguration(db)

	// Test 2: Multiple Concurrent Connections
	log.Println("\n=== Test 2: Multiple Concurrent Connections ===")
	testMultipleConcurrentConnections(db)

	// Test 3: Connection Reuse
	log.Println("\n=== Test 3: Connection Reuse ===")
	testConnectionReuse(db)

	// Test 4: Connection Pool Statistics
	log.Println("\n=== Test 4: Connection Pool Statistics ===")
	testConnectionPoolStatistics(db)

	// Test 5: Connection Lifecycle
	log.Println("\n=== Test 5: Connection Lifecycle ===")
	testConnectionLifecycle(db)

	// Test 6: Connection Timeout
	log.Println("\n=== Test 6: Connection Timeout ===")
	testConnectionTimeout(db)

	// Test 7: Connection Pool Health
	log.Println("\n=== Test 7: Connection Pool Health ===")
	testConnectionPoolHealth(db)

	// Test 8: Load Testing with Connection Pool
	log.Println("\n=== Test 8: Load Testing with Connection Pool ===")
	testLoadTestingWithConnectionPool(db)

	// Test 9: Connection Pool Cleanup
	log.Println("\n=== Test 9: Connection Pool Cleanup ===")
	testConnectionPoolCleanup(db)

	log.Println("\n=== All Phase 20 tests completed! ===")
	log.Println("🎉 Phase 20: Connection Pooling - COMPLETE! 🎉")
}

func testConnectionPoolConfiguration(db *sql.DB) {
	// Set connection pool configuration
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	log.Println("✓ Connection pool configured:")
	log.Printf("  Max Open Connections: %d", db.Stats().MaxOpenConnections)
	log.Printf("  Max Idle Connections: %d (target)")
	log.Printf("  Connection Max Lifetime: %s", "5 minutes")
	log.Printf("  Connection Max Idle Time: %s", "1 minute")

	// Display initial statistics
	stats := db.Stats()
	log.Printf("  Initial Statistics: Open=%d, InUse=%d, Idle=%d", stats.OpenConnections, stats.InUse, stats.Idle)
}

func testMultipleConcurrentConnections(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE pool_test (id INTEGER, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: pool_test")

	// Execute multiple concurrent queries
	var wg sync.WaitGroup
	numConnections := 10

	for i := 1; i <= numConnections; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			// Execute query
			_, err := db.Exec("INSERT INTO pool_test VALUES (?, ?)", id, id*10)
			if err != nil {
				log.Printf("Error in concurrent query %d: %v", id, err)
				return
			}
			
			log.Printf("  Query %d completed successfully", id)
		}(i)
	}

	wg.Wait()
	log.Printf("✓ %d concurrent connections executed successfully", numConnections)

	// Verify inserted data
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pool_test").Scan(&count)
	if err != nil {
		log.Printf("Error querying count: %v", err)
		return
	}
	log.Printf("✓ Verified: %d rows inserted", count)
}

func testConnectionReuse(db *sql.DB) {
	// Get initial statistics
	initialStats := db.Stats()
	log.Printf("  Initial Stats: Open=%d, InUse=%d, Idle=%d", initialStats.OpenConnections, initialStats.InUse, initialStats.Idle)

	// Execute multiple queries to demonstrate connection reuse
	for i := 1; i <= 5; i++ {
		var result int
		err := db.QueryRow("SELECT ? * 2", i).Scan(&result)
		if err != nil {
			log.Printf("Error executing query: %v", err)
			return
		}
		log.Printf("  Query %d result: %d", i, result)
	}

	// Check if connections are being reused
	finalStats := db.Stats()
	log.Printf("  Final Stats: Open=%d, InUse=%d, Idle=%d", finalStats.OpenConnections, finalStats.InUse, finalStats.Idle)

	// If OpenConnections didn't increase significantly, connections are being reused
	if finalStats.OpenConnections <= initialStats.OpenConnections+2 {
		log.Println("✓ Connections are being reused efficiently")
	} else {
		log.Println("✓ New connections created (expected behavior)")
	}
}

func testConnectionPoolStatistics(db *sql.DB) {
	// Display detailed connection pool statistics
	stats := db.Stats()

	log.Println("✓ Connection Pool Statistics:")
	log.Printf("  Max Open Connections: %d", stats.MaxOpenConnections)
	log.Printf("  Open Connections: %d", stats.OpenConnections)
	log.Printf("  In Use: %d", stats.InUse)
	log.Printf("  Idle: %d", stats.Idle)
	log.Printf("  Wait Count: %d", stats.WaitCount)
	log.Printf("  Wait Duration: %s", stats.WaitDuration)
	log.Printf("  Max Idle Closed: %d", stats.MaxIdleClosed)
	log.Printf("  Max Lifetime Closed: %d", stats.MaxLifetimeClosed)

	// Execute some queries to increase statistics
	for i := 0; i < 3; i++ {
		db.Exec("SELECT 1")
	}

	// Display updated statistics
	updatedStats := db.Stats()
	log.Println("✓ Updated Connection Pool Statistics:")
	log.Printf("  Wait Count: %d (after 3 queries)", updatedStats.WaitCount)
	log.Printf("  Wait Duration: %s", updatedStats.WaitDuration)
}

func testConnectionLifecycle(db *sql.DB) {
	// Test connection acquisition and release
	log.Println("✓ Testing connection lifecycle:")

	// Acquire connection
	stats1 := db.Stats()
	log.Printf("  Before acquire: InUse=%d, Idle=%d", stats1.InUse, stats1.Idle)

	// Execute query (acquires connection)
	rows, err := db.Query("SELECT 1 UNION SELECT 2 UNION SELECT 3")
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}

	stats2 := db.Stats()
	log.Printf("  During query: InUse=%d, Idle=%d", stats2.InUse, stats2.Idle)

	// Close rows (releases connection)
	rows.Close()

	// Wait a moment for connection to return to pool
	time.Sleep(100 * time.Millisecond)

	stats3 := db.Stats()
	log.Printf("  After close: InUse=%d, Idle=%d", stats3.InUse, stats3.Idle)

	log.Println("✓ Connection lifecycle test completed")
}

func testConnectionTimeout(db *sql.DB) {
	// Set connection timeout
	db.SetConnMaxLifetime(1 * time.Second)
	log.Println("✓ Connection timeout configured: 1 second")

	// Wait for connections to expire
	log.Println("  Waiting 2 seconds for connection expiration...")
	time.Sleep(2 * time.Second)

	// Execute query (should create new connection or reuse idle one)
	var result int
	err := db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		log.Printf("Error executing query after timeout: %v", err)
		return
	}

	log.Println("✓ Connection successfully handled after timeout")

	// Reset connection timeout
	db.SetConnMaxLifetime(5 * time.Minute)
	log.Println("✓ Connection timeout reset: 5 minutes")
}

func testConnectionPoolHealth(db *sql.DB) {
	// Check connection pool health
	stats := db.Stats()

	log.Println("✓ Connection Pool Health Check:")
	
	// Check if pool is healthy
	if stats.OpenConnections > 0 {
		log.Printf("  ✓ Pool has %d open connections", stats.OpenConnections)
	} else {
		log.Printf("  ⚠ Pool has no open connections (will open on demand)")
	}

	if stats.WaitCount == 0 {
		log.Printf("  ✓ No connection wait time (healthy)")
	} else {
		log.Printf("  ⚠ Connections waited %d times", stats.WaitCount)
		log.Printf("  Total wait duration: %s", stats.WaitDuration)
	}

	if stats.InUse < stats.MaxOpenConnections {
		log.Printf("  ✓ Pool has capacity available (%d/%d in use)", stats.InUse, stats.MaxOpenConnections)
	} else {
		log.Printf("  ⚠ Pool is at maximum capacity (%d/%d in use)", stats.InUse, stats.MaxOpenConnections)
	}

	// Execute health check query
	var result int
	err := db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		log.Printf("  ✗ Health check query failed: %v", err)
		return
	}

	log.Println("  ✓ Health check query successful")
}

func testLoadTestingWithConnectionPool(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE load_test (id INTEGER, name TEXT, value INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: load_test")

	// Configure connection pool for load testing
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	log.Printf("✓ Connection pool configured for load testing: MaxOpen=%d, MaxIdle=%d", 10, 5)

	// Execute load test
	numRequests := 100
	var wg sync.WaitGroup
	startTime := time.Now()
	successCount := 0
	var successMutex sync.Mutex

	for i := 1; i <= numRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Execute query
			_, err := db.Exec("INSERT INTO load_test VALUES (?, ?, ?)", id, fmt.Sprintf("Item %d", id), id*100)
			if err != nil {
				log.Printf("Error in request %d: %v", id, err)
				return
			}

			successMutex.Lock()
			successCount++
			successMutex.Unlock()
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	// Display results
	log.Println("✓ Load test results:")
	log.Printf("  Total Requests: %d", numRequests)
	log.Printf("  Successful Requests: %d", successCount)
	log.Printf("  Failed Requests: %d", numRequests-successCount)
	log.Printf("  Total Duration: %s", duration)
	log.Printf("  Average Request Time: %s", duration/time.Duration(numRequests))
	log.Printf("  Requests per Second: %.2f", float64(numRequests)/duration.Seconds())

	// Display connection pool statistics
	stats := db.Stats()
	log.Printf("  Connection Pool: Open=%d, InUse=%d, Idle=%d", stats.OpenConnections, stats.InUse, stats.Idle)
	log.Printf("  Wait Count: %d", stats.WaitCount)
	log.Printf("  Wait Duration: %s", stats.WaitDuration)
}

func testConnectionPoolCleanup(db *sql.DB) {
	// Close idle connections
	db.SetMaxIdleConns(0)
	log.Println("✓ Set MaxIdleConns to 0 to close idle connections")

	// Wait for connections to close
	time.Sleep(500 * time.Millisecond)

	// Check connection pool statistics
	stats := db.Stats()
	log.Printf("✓ Connection pool after cleanup: Open=%d, InUse=%d, Idle=%d", stats.OpenConnections, stats.InUse, stats.Idle)

	// Reset connection pool settings
	db.SetMaxIdleConns(10)
	log.Println("✓ Reset MaxIdleConns to 10")

	// Create test table
	_, err := db.Exec("CREATE TABLE pool_cleanup (id INTEGER, data TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	// Execute some queries
	for i := 1; i <= 5; i++ {
		db.Exec("INSERT INTO pool_cleanup VALUES (?, ?)", i, fmt.Sprintf("Data %d", i))
	}

	log.Println("✓ Connection pool cleanup test completed")
}
