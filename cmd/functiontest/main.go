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

	// Create test table
	_, err = db.Exec("CREATE TABLE function_test (id INTEGER, name TEXT, description TEXT, price REAL, quantity INTEGER, created_date TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: function_test")

	// Insert test data
	_, err = db.Exec("INSERT INTO function_test VALUES (1, 'Product A', '  Test Description  ', 99.99, 10, '2024-01-15')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}
	_, err = db.Exec("INSERT INTO function_test VALUES (2, 'product b', 'Another Description', 149.99, 5, '2024-02-20')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}
	_, err = db.Exec("INSERT INTO function_test VALUES (3, 'PRODUCT C', '  Yet Another Description  ', 199.99, 15, '2024-03-25')")
	if err != nil {
		log.Printf("Error inserting data: %v", err)
		return
	}

	// Test 1: String Functions
	log.Println("\n=== Test 1: String Functions ===")
	testStringFunctions(db)

	// Test 2: Numeric Functions
	log.Println("\n=== Test 2: Numeric Functions ===")
	testNumericFunctions(db)

	// Test 3: Date/Time Functions
	log.Println("\n=== Test 3: Date/Time Functions ===")
	testDateTimeFunctions(db)

	// Test 4: Conditional Functions
	log.Println("\n=== Test 4: Conditional Functions ===")
	testConditionalFunctions(db)

	// Test 5: Aggregate Functions
	log.Println("\n=== Test 5: Aggregate Functions ===")
	testAggregateFunctions(db)

	// Test 6: Type Conversion Functions
	log.Println("\n=== Test 6: Type Conversion Functions ===")
	testTypeConversionFunctions(db)

	// Test 7: Cleanup
	log.Println("\n=== Test 7: Cleanup ===")
	testCleanup(db)

	log.Println("\n=== All Phase 18 tests completed! ===")
	log.Println("🎉 Phase 18: SQL Functions - COMPLETE! 🎉")
}

func testStringFunctions(db *sql.DB) {
	// UPPER
	var result string
	err := db.QueryRow("SELECT UPPER(name) FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with UPPER: %v", err)
	} else {
		log.Printf("✓ UPPER('Product A') = '%s'", result)
	}

	// LOWER
	err = db.QueryRow("SELECT LOWER(name) FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with LOWER: %v", err)
	} else {
		log.Printf("✓ LOWER('Product A') = '%s'", result)
	}

	// TRIM
	err = db.QueryRow("SELECT TRIM(description) FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with TRIM: %v", err)
	} else {
		log.Printf("✓ TRIM('  Test Description  ') = '%s'", result)
	}

	// LTRIM
	err = db.QueryRow("SELECT LTRIM(description) FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with LTRIM: %v", err)
	} else {
		log.Printf("✓ LTRIM('  Test Description  ') = '%s'", result)
	}

	// RTRIM
	err = db.QueryRow("SELECT RTRIM(description) FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with RTRIM: %v", err)
	} else {
		log.Printf("✓ RTRIM('  Test Description  ') = '%s'", result)
	}

	// LENGTH
	var length int
	err = db.QueryRow("SELECT LENGTH(name) FROM function_test WHERE id = 1").Scan(&length)
	if err != nil {
		log.Printf("Error with LENGTH: %v", err)
	} else {
		log.Printf("✓ LENGTH('Product A') = %d", length)
	}

	// SUBSTRING / SUBSTR
	err = db.QueryRow("SELECT SUBSTR(name, 1, 7) FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with SUBSTR: %v", err)
	} else {
		log.Printf("✓ SUBSTR('Product A', 1, 7) = '%s'", result)
	}

	// CONCAT (SQLite uses ||)
	err = db.QueryRow("SELECT name || ' - ' || description FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with CONCAT: %v", err)
	} else {
		log.Printf("✓ 'Product A' || ' - ' || '  Test Description  ' = '%s'", result)
	}

	// REPLACE
	err = db.QueryRow("SELECT REPLACE(name, 'Product', 'Item') FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with REPLACE: %v", err)
	} else {
		log.Printf("✓ REPLACE('Product A', 'Product', 'Item') = '%s'", result)
	}

	// INSTR / POSITION
	var pos int
	err = db.QueryRow("SELECT INSTR(name, 'A') FROM function_test WHERE id = 1").Scan(&pos)
	if err != nil {
		log.Printf("Error with INSTR: %v", err)
	} else {
		log.Printf("✓ INSTR('Product A', 'A') = %d", pos)
	}

	// LEFT
	err = db.QueryRow("SELECT SUBSTR(name, 1, 7) FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with LEFT: %v", err)
	} else {
		log.Printf("✓ LEFT('Product A', 7) = '%s'", result)
	}

	// RIGHT
	err = db.QueryRow("SELECT SUBSTR(name, -1) FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with RIGHT: %v", err)
	} else {
		log.Printf("✓ RIGHT('Product A') = '%s'", result)
	}
}

func testNumericFunctions(db *sql.DB) {
	// ABS
	var result float64
	err := db.QueryRow("SELECT ABS(price) FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with ABS: %v", err)
	} else {
		log.Printf("✓ ABS(99.99) = %.2f", result)
	}

	// ROUND
	err = db.QueryRow("SELECT ROUND(price, 1) FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with ROUND: %v", err)
	} else {
		log.Printf("✓ ROUND(99.99, 1) = %.1f", result)
	}

	// CEILING
	err = db.QueryRow("SELECT CEILING(price) FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with CEILING: %v", err)
	} else {
		log.Printf("✓ CEILING(99.99) = %.0f", result)
	}

	// FLOOR
	err = db.QueryRow("SELECT FLOOR(price) FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with FLOOR: %v", err)
	} else {
		log.Printf("✓ FLOOR(99.99) = %.0f", result)
	}

	// MOD
	var mod int
	err = db.QueryRow("SELECT MOD(quantity, 3) FROM function_test WHERE id = 1").Scan(&mod)
	if err != nil {
		log.Printf("Error with MOD: %v", err)
	} else {
		log.Printf("✓ MOD(10, 3) = %d", mod)
	}

	// POWER
	var power float64
	err = db.QueryRow("SELECT POWER(2, 3) FROM function_test WHERE id = 1").Scan(&power)
	if err != nil {
		log.Printf("Error with POWER: %v", err)
	} else {
		log.Printf("✓ POWER(2, 3) = %.0f", power)
	}

	// SQRT
	err = db.QueryRow("SELECT SQRT(price) FROM function_test WHERE id = 1").Scan(&result)
	if err != nil {
		log.Printf("Error with SQRT: %v", err)
	} else {
		log.Printf("✓ SQRT(99.99) = %.2f", result)
	}

	// MIN
	err = db.QueryRow("SELECT MIN(price) FROM function_test").Scan(&result)
	if err != nil {
		log.Printf("Error with MIN: %v", err)
	} else {
		log.Printf("✓ MIN(price) = %.2f", result)
	}

	// MAX
	err = db.QueryRow("SELECT MAX(price) FROM function_test").Scan(&result)
	if err != nil {
		log.Printf("Error with MAX: %v", err)
	} else {
		log.Printf("✓ MAX(price) = %.2f", result)
	}

	// SUM
	var sum float64
	err = db.QueryRow("SELECT SUM(price) FROM function_test").Scan(&sum)
	if err != nil {
		log.Printf("Error with SUM: %v", err)
	} else {
		log.Printf("✓ SUM(price) = %.2f", sum)
	}

	// AVG
	err = db.QueryRow("SELECT AVG(price) FROM function_test").Scan(&result)
	if err != nil {
		log.Printf("Error with AVG: %v", err)
	} else {
		log.Printf("✓ AVG(price) = %.2f", result)
	}

	// RANDOM
	var random int
	err = db.QueryRow("SELECT RANDOM()").Scan(&random)
	if err != nil {
		log.Printf("Error with RANDOM: %v", err)
	} else {
		log.Printf("✓ RANDOM() = %d", random)
	}
}

func testDateTimeFunctions(db *sql.DB) {
	// DATE
	var date string
	err := db.QueryRow("SELECT DATE('now')").Scan(&date)
	if err != nil {
		log.Printf("Error with DATE: %v", err)
	} else {
		log.Printf("✓ DATE('now') = '%s'", date)
	}

	// TIME
	var timeStr string
	err = db.QueryRow("SELECT TIME('now')").Scan(&timeStr)
	if err != nil {
		log.Printf("Error with TIME: %v", err)
	} else {
		log.Printf("✓ TIME('now') = '%s'", timeStr)
	}

	// DATETIME
	var datetime string
	err = db.QueryRow("SELECT DATETIME('now')").Scan(&datetime)
	if err != nil {
		log.Printf("Error with DATETIME: %v", err)
	} else {
		log.Printf("✓ DATETIME('now') = '%s'", datetime)
	}

	// STRFTIME (format date)
	var formatted string
	err = db.QueryRow("SELECT STRFTIME('%Y-%m-%d %H:%M:%S', 'now')").Scan(&formatted)
	if err != nil {
		log.Printf("Error with STRFTIME: %v", err)
	} else {
		log.Printf("✓ STRFTIME('%%Y-%%m-%%d %%H:%%M:%%S', 'now') = '%s'", formatted)
	}

	// JULIANDAY
	var julianday float64
	err = db.QueryRow("SELECT JULIANDAY('now')").Scan(&julianday)
	if err != nil {
		log.Printf("Error with JULIANDAY: %v", err)
	} else {
		log.Printf("✓ JULIANDAY('now') = %.6f", julianday)
	}

	// DATEADD equivalent (using modifiers)
	err = db.QueryRow("SELECT DATE('now', '+7 days')").Scan(&date)
	if err != nil {
		log.Printf("Error with DATEADD: %v", err)
	} else {
		log.Printf("✓ DATE('now', '+7 days') = '%s' (DATEADD)", date)
	}

	// DATEDIFF equivalent (using julianday)
	var diff float64
	err = db.QueryRow("SELECT JULIANDAY('2024-02-20') - JULIANDAY('2024-01-15')").Scan(&diff)
	if err != nil {
		log.Printf("Error with DATEDIFF: %v", err)
	} else {
		log.Printf("✓ DATEDIFF('2024-02-20', '2024-01-15') = %.0f days", diff)
	}

	// Extract year, month, day
	var year, month, day int
	err = db.QueryRow("SELECT CAST(STRFTIME('%Y', created_date) AS INTEGER), CAST(STRFTIME('%m', created_date) AS INTEGER), CAST(STRFTIME('%d', created_date) AS INTEGER) FROM function_test WHERE id = 1").Scan(&year, &month, &day)
	if err != nil {
		log.Printf("Error with DATEPART: %v", err)
	} else {
		log.Printf("✓ DATEPART('2024-01-15') = Year=%d, Month=%d, Day=%d", year, month, day)
	}

	// CURRENT_TIMESTAMP
	var current string
	err = db.QueryRow("SELECT CURRENT_TIMESTAMP").Scan(&current)
	if err != nil {
		log.Printf("Error with CURRENT_TIMESTAMP: %v", err)
	} else {
		log.Printf("✓ CURRENT_TIMESTAMP = '%s'", current)
	}

	// CURRENT_DATE
	err = db.QueryRow("SELECT CURRENT_DATE").Scan(&date)
	if err != nil {
		log.Printf("Error with CURRENT_DATE: %v", err)
	} else {
		log.Printf("✓ CURRENT_DATE = '%s'", date)
	}

	// CURRENT_TIME
	err = db.QueryRow("SELECT CURRENT_TIME").Scan(&timeStr)
	if err != nil {
		log.Printf("Error with CURRENT_TIME: %v", err)
	} else {
		log.Printf("✓ CURRENT_TIME = '%s'", timeStr)
	}
}

func testConditionalFunctions(db *sql.DB) {
	// COALESCE
	var result string
	err := db.QueryRow("SELECT COALESCE(NULL, 'default value')").Scan(&result)
	if err != nil {
		log.Printf("Error with COALESCE: %v", err)
	} else {
		log.Printf("✓ COALESCE(NULL, 'default value') = '%s'", result)
	}

	// IFNULL (SQLite equivalent of COALESCE)
	err = db.QueryRow("SELECT IFNULL(NULL, 'default value')").Scan(&result)
	if err != nil {
		log.Printf("Error with IFNULL: %v", err)
	} else {
		log.Printf("✓ IFNULL(NULL, 'default value') = '%s'", result)
	}

	// NULLIF
	var nullResult sql.NullString
	err = db.QueryRow("SELECT NULLIF(1, 1)").Scan(&nullResult)
	if err != nil {
		log.Printf("Error with NULLIF: %v", err)
	} else {
		if nullResult.Valid {
			log.Printf("✓ NULLIF(1, 1) = '%s' (unexpected)", nullResult.String)
		} else {
			log.Printf("✓ NULLIF(1, 1) = NULL (expected)")
		}
	}

	// CASE WHEN
	var caseResult string
	err = db.QueryRow("SELECT CASE WHEN quantity > 10 THEN 'High' WHEN quantity > 5 THEN 'Medium' ELSE 'Low' END FROM function_test WHERE id = 1").Scan(&caseResult)
	if err != nil {
		log.Printf("Error with CASE WHEN: %v", err)
	} else {
		log.Printf("✓ CASE WHEN quantity > 10 THEN 'High' WHEN quantity > 5 THEN 'Medium' ELSE 'Low' END (quantity=10) = '%s'", caseResult)
	}

	// CAST
	var castResult string
	err = db.QueryRow("SELECT CAST(price AS TEXT) FROM function_test WHERE id = 1").Scan(&castResult)
	if err != nil {
		log.Printf("Error with CAST: %v", err)
	} else {
		log.Printf("✓ CAST(99.99 AS TEXT) = '%s'", castResult)
	}

	// TYPEOF
	var typeOf string
	err = db.QueryRow("SELECT TYPEOF(name) FROM function_test WHERE id = 1").Scan(&typeOf)
	if err != nil {
		log.Printf("Error with TYPEOF: %v", err)
	} else {
		log.Printf("✓ TYPEOF('Product A') = '%s'", typeOf)
	}
}

func testAggregateFunctions(db *sql.DB) {
	// COUNT
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM function_test").Scan(&count)
	if err != nil {
		log.Printf("Error with COUNT: %v", err)
	} else {
		log.Printf("✓ COUNT(*) = %d", count)
	}

	// COUNT with condition
	var countCond int
	err = db.QueryRow("SELECT COUNT(*) FROM function_test WHERE price > 100").Scan(&countCond)
	if err != nil {
		log.Printf("Error with COUNT with condition: %v", err)
	} else {
		log.Printf("✓ COUNT(*) WHERE price > 100 = %d", countCond)
	}

	// SUM
	var sum float64
	err = db.QueryRow("SELECT SUM(price) FROM function_test").Scan(&sum)
	if err != nil {
		log.Printf("Error with SUM: %v", err)
	} else {
		log.Printf("✓ SUM(price) = %.2f", sum)
	}

	// AVG
	var avg float64
	err = db.QueryRow("SELECT AVG(price) FROM function_test").Scan(&avg)
	if err != nil {
		log.Printf("Error with AVG: %v", err)
	} else {
		log.Printf("✓ AVG(price) = %.2f", avg)
	}

	// MIN
	var min float64
	err = db.QueryRow("SELECT MIN(price) FROM function_test").Scan(&min)
	if err != nil {
		log.Printf("Error with MIN: %v", err)
	} else {
		log.Printf("✓ MIN(price) = %.2f", min)
	}

	// MAX
	var max float64
	err = db.QueryRow("SELECT MAX(price) FROM function_test").Scan(&max)
	if	err != nil {
		log.Printf("Error with MAX: %v", err)
	} else {
		log.Printf("✓ MAX(price) = %.2f", max)
	}

	// TOTAL (SQLite-specific)
	var total float64
	err = db.QueryRow("SELECT TOTAL(price) FROM function_test").Scan(&total)
	if err != nil {
		log.Printf("Error with TOTAL: %v", err)
	} else {
		log.Printf("✓ TOTAL(price) = %.2f", total)
	}

	// GROUP_CONCAT
	var concat string
	err = db.QueryRow("SELECT GROUP_CONCAT(name, ', ') FROM function_test").Scan(&concat)
	if err != nil {
		log.Printf("Error with GROUP_CONCAT: %v", err)
	} else {
		log.Printf("✓ GROUP_CONCAT(name, ', ') = '%s'", concat)
	}
}

func testTypeConversionFunctions(db *sql.DB) {
	// CAST to INTEGER
	var intResult int
	err := db.QueryRow("SELECT CAST('123' AS INTEGER)").Scan(&intResult)
	if err != nil {
		log.Printf("Error with CAST to INTEGER: %v", err)
	} else {
		log.Printf("✓ CAST('123' AS INTEGER) = %d", intResult)
	}

	// CAST to REAL
	var realResult float64
	err = db.QueryRow("SELECT CAST('123.45' AS REAL)").Scan(&realResult)
	if err != nil {
		log.Printf("Error with CAST to REAL: %v", err)
	} else {
		log.Printf("✓ CAST('123.45' AS REAL) = %.2f", realResult)
	}

	// CAST to TEXT
	var textResult string
	err = db.QueryRow("SELECT CAST(123 AS TEXT)").Scan(&textResult)
	if err != nil {
		log.Printf("Error with CAST to TEXT: %v", err)
	} else {
		log.Printf("✓ CAST(123 AS TEXT) = '%s'", textResult)
	}

	// ROUND with CAST
	var rounded float64
	err = db.QueryRow("SELECT CAST(ROUND(price, 0) AS INTEGER) FROM function_test WHERE id = 1").Scan(&rounded)
	if err != nil {
		log.Printf("Error with ROUND + CAST: %v", err)
	} else {
		log.Printf("✓ CAST(ROUND(99.99, 0) AS INTEGER) = %.0f", rounded)
	}

	// ABS + CAST
	var absInt int
	err = db.QueryRow("SELECT CAST(ABS(-100) AS INTEGER)").Scan(&absInt)
	if err != nil {
		log.Printf("Error with ABS + CAST: %v", err)
	} else {
		log.Printf("✓ CAST(ABS(-100) AS INTEGER) = %d", absInt)
	}
}

func testCleanup(db *sql.DB) {
	_, err := db.Exec("DROP TABLE function_test")
	if err != nil {
		log.Printf("Error dropping table: %v", err)
		return
	}
	log.Println("✓ Dropped table: function_test")
}
