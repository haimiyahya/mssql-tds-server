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

	// Test 1: Current Date/Time (date, time, datetime)
	log.Println("\n=== Test 1: Current Date/Time (date, time, datetime) ===")
testCurrentDateTime(db)

	// Test 2: Date/Time Formatting (strftime)
	log.Println("\n=== Test 2: Date/Time Formatting (strftime) ===")
testDateTimeFormatting(db)

	// Test 3: Date Arithmetic (adding days, months, etc.)
	log.Println("\n=== Test 3: Date Arithmetic ===")
testDateArithmetic(db)

	// Test 4: Date/Time Modifiers (start of, weekdays, etc.)
	log.Println("\n=== Test 4: Date/Time Modifiers ===")
testDateTimeModifiers(db)

	// Test 5: Julian Day Functions
	log.Println("\n=== Test 5: Julian Day Functions ===")
testJulianDay(db)

	// Test 6: Date Comparisons
	log.Println("\n=== Test 6: Date Comparisons ===")
testDateComparisons(db)

	// Test 7: Date/Time with SQL Tables
	log.Println("\n=== Test 7: Date/Time with SQL Tables ===")
testDateTimeWithTables(db)

	// Test 8: Timezone Handling
	log.Println("\n=== Test 8: Timezone Handling ===")
testTimezoneHandling(db)

	// Test 9: Date Parsing
	log.Println("\n=== Test 9: Date Parsing ===")
testDateParsing(db)

	// Test 10: Age Calculations
	log.Println("\n=== Test 10: Age Calculations ===")
testAgeCalculations(db)

	// Test 11: Business Days Calculations
	log.Println("\n=== Test 11: Business Days Calculations ===")
testBusinessDays(db)

	// Test 12: Date/Time in Different Formats
	log.Println("\n=== Test 12: Date/Time in Different Formats ===")
testDateTimeFormats(db)

	// Test 13: Date/Time with Aggregates
	log.Println("\n=== Test 13: Date/Time with Aggregates ===")
testDateTimeWithAggregates(db)

	// Test 14: Cleanup
	log.Println("\n=== Test 14: Cleanup ===")
testCleanup(db)

	log.Println("\n=== All Phase 31 tests completed! ===")
	log.Println("🎉 Phase 31: Advanced Date/Time Functions - COMPLETE! 🎉")
}

func testCurrentDateTime(db *sql.DB) {
	// Current date
	var currentDate string
	err := db.QueryRow("SELECT date('now')").Scan(&currentDate)
	if err != nil {
		log.Printf("Error getting current date: %v", err)
		return
	}
	log.Printf("✓ Current date: %s", currentDate)

	// Current time
	var currentTime string
	err = db.QueryRow("SELECT time('now')").Scan(&currentTime)
	if err != nil {
		log.Printf("Error getting current time: %v", err)
		return
	}
	log.Printf("✓ Current time: %s", currentTime)

	// Current datetime
	var currentDateTime string
	err = db.QueryRow("SELECT datetime('now')").Scan(&currentDateTime)
	if err != nil {
		log.Printf("Error getting current datetime: %v", err)
		return
	}
	log.Printf("✓ Current datetime: %s", currentDateTime)

	// Multiple datetime functions
	rows, err := db.Query(`
		SELECT 
		  date('now') as d,
		  time('now') as t,
		  datetime('now') as dt
	`)
	if err != nil {
		log.Printf("Error getting multiple datetime values: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var d, t, dt string
		err := rows.Scan(&d, &t, &dt)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("✓ Date: %s, Time: %s, DateTime: %s", d, t, dt)
	}
}

func testDateTimeFormatting(db *sql.DB) {
	// Various date formats
	var result string
	err := db.QueryRow("SELECT strftime('%Y-%m-%d', 'now')").Scan(&result)
	if err != nil {
		log.Printf("Error formatting date: %v", err)
		return
	}
	log.Printf("✓ Date format (YYYY-MM-DD): %s", result)

	// Time format
	err = db.QueryRow("SELECT strftime('%H:%M:%S', 'now')").Scan(&result)
	if err != nil {
		log.Printf("Error formatting time: %v", err)
		return
	}
	log.Printf("✓ Time format (HH:MM:SS): %s", result)

	// Custom format
	err = db.QueryRow("SELECT strftime('%A, %B %d, %Y', 'now')").Scan(&result)
	if err != nil {
		log.Printf("Error formatting custom date: %v", err)
		return
	}
	log.Printf("✓ Custom format: %s", result)

	// Unix timestamp
	var unixTime int
	err = db.QueryRow("SELECT strftime('%s', 'now')").Scan(&unixTime)
	if err != nil {
		log.Printf("Error getting unix time: %v", err)
		return
	}
	log.Printf("✓ Unix timestamp: %d", unixTime)

	// Multiple format specifiers
	rows, err := db.Query(`
		SELECT 
		  strftime('%Y', 'now') as year,
		  strftime('%m', 'now') as month,
		  strftime('%d', 'now') as day,
		  strftime('%H', 'now') as hour,
		  strftime('%M', 'now') as minute,
		  strftime('%S', 'now') as second
	`)
	if err != nil {
		log.Printf("Error getting multiple format specifiers: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var year, month, day, hour, minute, second string
		err := rows.Scan(&year, &month, &day, &hour, &minute, &second)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("✓ Year: %s, Month: %s, Day: %s, Hour: %s, Minute: %s, Second: %s",
			year, month, day, hour, minute, second)
	}
}

func testDateArithmetic(db *sql.DB) {
	// Add days
	var result string
	err := db.QueryRow("SELECT date('now', '+7 days')").Scan(&result)
	if err != nil {
		log.Printf("Error adding days: %v", err)
		return
	}
	log.Printf("✓ Date + 7 days: %s", result)

	// Subtract days
	err = db.QueryRow("SELECT date('now', '-7 days')").Scan(&result)
	if err != nil {
		log.Printf("Error subtracting days: %v", err)
		return
	}
	log.Printf("✓ Date - 7 days: %s", result)

	// Add months
	err = db.QueryRow("SELECT date('now', '+3 months')").Scan(&result)
	if err != nil {
		log.Printf("Error adding months: %v", err)
		return
	}
	log.Printf("✓ Date + 3 months: %s", result)

	// Add years
	err = db.QueryRow("SELECT date('now', '+1 year')").Scan(&result)
	if err != nil {
		log.Printf("Error adding years: %v", err)
		return
	}
	log.Printf("✓ Date + 1 year: %s", result)

	// Add hours to datetime
	err = db.QueryRow("SELECT datetime('now', '+5 hours')").Scan(&result)
	if err != nil {
		log.Printf("Error adding hours: %v", err)
		return
	}
	log.Printf("✓ DateTime + 5 hours: %s", result)

	// Add minutes
	err = db.QueryRow("SELECT datetime('now', '+30 minutes')").Scan(&result)
	if err != nil {
		log.Printf("Error adding minutes: %v", err)
		return
	}
	log.Printf("✓ DateTime + 30 minutes: %s", result)

	// Add seconds
	err = db.QueryRow("SELECT datetime('now', '+90 seconds')").Scan(&result)
	if err != nil {
		log.Printf("Error adding seconds: %v", err)
		return
	}
	log.Printf("✓ DateTime + 90 seconds: %s", result)

	// Complex arithmetic
	err = db.QueryRow("SELECT date('now', '+1 year', '+2 months', '+3 days')").Scan(&result)
	if err != nil {
		log.Printf("Error with complex arithmetic: %v", err)
		return
	}
	log.Printf("✓ Date + 1 year + 2 months + 3 days: %s", result)
}

func testDateTimeModifiers(db *sql.DB) {
	// Start of day
	var result string
	err := db.QueryRow("SELECT date('now', 'start of day')").Scan(&result)
	if err != nil {
		log.Printf("Error getting start of day: %v", err)
		return
	}
	log.Printf("✓ Start of day: %s", result)

	// Start of month
	err = db.QueryRow("SELECT date('now', 'start of month')").Scan(&result)
	if err != nil {
		log.Printf("Error getting start of month: %v", err)
		return
	}
	log.Printf("✓ Start of month: %s", result)

	// Start of year
	err = db.QueryRow("SELECT date('now', 'start of year')").Scan(&result)
	if err != nil {
		log.Printf("Error getting start of year: %v", err)
		return
	}
	log.Printf("✓ Start of year: %s", result)

	// Next Monday
	err = db.QueryRow("SELECT date('now', 'weekday 0')").Scan(&result)
	if err != nil {
		log.Printf("Error getting next Monday: %v", err)
		return
	}
	log.Printf("✓ Next Monday (weekday 0): %s", result)

	// Next Tuesday
	err = db.QueryRow("SELECT date('now', 'weekday 1')").Scan(&result)
	if err != nil {
		log.Printf("Error getting next Tuesday: %v", err)
		return
	}
	log.Printf("✓ Next Tuesday (weekday 1): %s", result)

	// End of month
	err = db.QueryRow("SELECT date('now', 'start of month', '+1 month', '-1 day')").Scan(&result)
	if err != nil {
		log.Printf("Error getting end of month: %v", err)
		return
	}
	log.Printf("✓ End of month: %s", result)

	// First day of next month
	err = db.QueryRow("SELECT date('now', 'start of month', '+1 month')").Scan(&result)
	if err != nil {
		log.Printf("Error getting first day of next month: %v", err)
		return
	}
	log.Printf("✓ First day of next month: %s", result)
}

func testJulianDay(db *sql.DB) {
	// Current Julian day
	var julianDay float64
	err := db.QueryRow("SELECT julianday('now')").Scan(&julianDay)
	if err != nil {
		log.Printf("Error getting Julian day: %v", err)
		return
	}
	log.Printf("✓ Current Julian day: %.6f", julianDay)

	// Julian day for specific date
	err = db.QueryRow("SELECT julianday('2024-01-01')").Scan(&julianDay)
	if err != nil {
		log.Printf("Error getting Julian day for date: %v", err)
		return
	}
	log.Printf("✓ Julian day for 2024-01-01: %.6f", julianDay)

	// Julian day difference
	var julianDay1, julianDay2 float64
	err = db.QueryRow("SELECT julianday('2024-01-01')").Scan(&julianDay1)
	if err != nil {
		log.Printf("Error getting Julian day 1: %v", err)
		return
	}
	err = db.QueryRow("SELECT julianday('2024-12-31')").Scan(&julianDay2)
	if err != nil {
		log.Printf("Error getting Julian day 2: %v", err)
		return
	}
	daysDiff := julianDay2 - julianDay1
	log.Printf("✓ Days difference between 2024-01-01 and 2024-12-31: %.1f", daysDiff)

	// Convert Julian day back to date
	var date string
	err = db.QueryRow("SELECT date(julianday('2024-01-01'))").Scan(&date)
	if err != nil {
		log.Printf("Error converting Julian day to date: %v", err)
		return
	}
	log.Printf("✓ Julian day back to date: %s", date)
}
func testDateTimeWithTables(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE appointments (id INTEGER, name TEXT, appointment_date DATETIME)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: appointments")

	// Insert test data
	appointments := []struct {
		id   int
		name string
		date string
	}{
		{1, "Appointment 1", "2024-01-15 09:00:00"},
		{2, "Appointment 2", "2024-02-20 14:30:00"},
		{3, "Appointment 3", "2024-03-25 10:15:00"},
		{4, "Appointment 4", "2024-04-10 16:45:00"},
	}

	for _, appt := range appointments {
		_, err = db.Exec("INSERT INTO appointments VALUES (?, ?, ?)", appt.id, appt.name, appt.date)
		if err != nil {
			log.Printf("Error inserting appointment: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 4 appointments")

	// Query appointments with datetime extraction
	log.Println("✓ Appointments with datetime extraction:")
	rows, err := db.Query(`
		SELECT 
		  id,
		  name,
		  appointment_date,
		  strftime('%Y', appointment_date) as year,
		  strftime('%m', appointment_date) as month,
		  strftime('%d', appointment_date) as day,
		  strftime('%H', appointment_date) as hour,
		  strftime('%M', appointment_date) as minute
		FROM appointments
		ORDER BY appointment_date
	`)
	if err != nil {
		log.Printf("Error querying appointments: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, appointmentDate, year, month, day, hour, minute string
		err := rows.Scan(&id, &name, &appointmentDate, &year, &month, &day, &hour, &minute)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - %s (Year: %s, Month: %s, Day: %s, Time: %s:%s)",
			id, name, appointmentDate, year, month, day, hour, minute)
	}

	// Future appointments
	log.Println("✓ Future appointments:")
	rows, err = db.Query(`
		SELECT id, name, appointment_date
		FROM appointments
		WHERE appointment_date > datetime('now')
		ORDER BY appointment_date
	`)
	if err != nil {
		log.Printf("Error querying future appointments: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, appointmentDate string
		err := rows.Scan(&id, &name, &appointmentDate)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - %s", id, name, appointmentDate)
	}

	// Update appointment date
	_, err = db.Exec(`
		UPDATE appointments
		SET appointment_date = datetime(appointment_date, '+1 day')
		WHERE id = 1
	`)
	if err != nil {
		log.Printf("Error updating appointment: %v", err)
		return
	}
	log.Println("✓ Updated appointment 1 date +1 day")

	// Verify update
	var updatedDate string
	err = db.QueryRow("SELECT appointment_date FROM appointments WHERE id = 1").Scan(&updatedDate)
	if err != nil {
		log.Printf("Error verifying update: %v", err)
		return
	}
	log.Printf("✓ Updated appointment date: %s", updatedDate)
}

func testDateComparisons(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE events (id INTEGER, name TEXT, event_date DATE)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: events")

	// Insert test data
	events := []struct {
		id   int
		name string
		date string
	}{
		{1, "Event 1", "2024-01-15"},
		{2, "Event 2", "2024-02-20"},
		{3, "Event 3", "2024-03-25"},
	}

	for _, event := range events {
		_, err = db.Exec("INSERT INTO events VALUES (?, ?, ?)", event.id, event.name, event.date)
		if err != nil {
			log.Printf("Error inserting event: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 3 events")

	// Events after date
	log.Println("✓ Events after 2024-02-01:")
	rows, err := db.Query("SELECT id, name, event_date FROM events WHERE event_date > ?", "2024-02-01")
	if err != nil {
		log.Printf("Error querying events: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, eventDate string
		err := rows.Scan(&id, &name, &eventDate)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - %s", id, name, eventDate)
	}

	// Events between dates
	log.Println("✓ Events between 2024-01-01 and 2024-03-01:")
	rows, err = db.Query("SELECT id, name, event_date FROM events WHERE event_date BETWEEN ? AND ?", "2024-01-01", "2024-03-01")
	if err != nil {
		log.Printf("Error querying events between: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, eventDate string
		err := rows.Scan(&id, &name, &eventDate)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - %s", id, name, eventDate)
	}

	// Events in current month
	log.Println("✓ Events in current month:")
	rows, err = db.Query(`
		SELECT id, name, event_date 
		FROM events 
		WHERE strftime('%Y-%m', event_date) = strftime('%Y-%m', 'now')
	`)
	if err != nil {
		log.Printf("Error querying events in current month: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, eventDate string
		err := rows.Scan(&id, &name, &eventDate)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - %s", id, name, eventDate)
	}
}

func testTimezoneHandling(db *sql.DB) {
	// Current time in UTC
	var utcTime string
	err := db.QueryRow("SELECT datetime('now', 'utc')").Scan(&utcTime)
	if err != nil {
		log.Printf("Error getting UTC time: %v", err)
		return
	}
	log.Printf("✓ UTC time: %s", utcTime)

	// Current time in local timezone
	var localTime string
	err = db.QueryRow("SELECT datetime('now', 'localtime')").Scan(&localTime)
	if err != nil {
		log.Printf("Error getting local time: %v", err)
		return
	}
	log.Printf("✓ Local time: %s", localTime)

	// Convert UTC to local
	var convertedTime string
	err = db.QueryRow("SELECT time('now', 'utc', 'localtime')").Scan(&convertedTime)
	if err != nil {
		log.Printf("Error converting UTC to local: %v", err)
		return
	}
	log.Printf("✓ UTC converted to local: %s", convertedTime)

	// Add timezone offset
	var offsetTime string
	err = db.QueryRow("SELECT datetime('now', '+8 hours')").Scan(&offsetTime)
	if err != nil {
		log.Printf("Error adding timezone offset: %v", err)
		return
	}
	log.Printf("✓ Time +8 hours (UTC+8): %s", offsetTime)
}

func testDateParsing(db *sql.DB) {
	// Parse ISO format
	var parsedDate string
	err := db.QueryRow("SELECT date('2024-01-15')").Scan(&parsedDate)
	if err != nil {
		log.Printf("Error parsing ISO date: %v", err)
		return
	}
	log.Printf("✓ Parsed ISO date (2024-01-15): %s", parsedDate)

	// Parse datetime
	var parsedDateTime string
	err = db.QueryRow("SELECT datetime('2024-01-15 14:30:00')").Scan(&parsedDateTime)
	if err != nil {
		log.Printf("Error parsing datetime: %v", err)
		return
	}
	log.Printf("✓ Parsed datetime (2024-01-15 14:30:00): %s", parsedDateTime)

	// Parse time
	var parsedTime string
	err = db.QueryRow("SELECT time('14:30:00')").Scan(&parsedTime)
	if err != nil {
		log.Printf("Error parsing time: %v", err)
		return
	}
	log.Printf("✓ Parsed time (14:30:00): %s", parsedTime)

	// Parse custom format
	var customDate string
	err = db.QueryRow("SELECT strftime('%Y-%m-%d', 'January 15, 2024')").Scan(&customDate)
	if err != nil {
		log.Printf("Error parsing custom date format: %v", err)
		return
	}
	log.Printf("✓ Parsed custom format: %s", customDate)
}

func testAgeCalculations(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE people (id INTEGER, name TEXT, birth_date DATE)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: people")

	// Insert test data
	people := []struct {
		id   int
		name string
		date string
	}{
		{1, "Person 1", "1990-01-15"},
		{2, "Person 2", "1985-06-20"},
		{3, "Person 3", "2000-11-30"},
	}

	for _, person := range people {
		_, err = db.Exec("INSERT INTO people VALUES (?, ?, ?)", person.id, person.name, person.date)
		if err != nil {
			log.Printf("Error inserting person: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 3 people")

	// Calculate age
	log.Println("✓ Age calculations:")
	rows, err := db.Query(`
		SELECT 
		  id, 
		  name, 
		  birth_date,
		  (julianday('now') - julianday(birth_date)) / 365.25 as age
		FROM people
		ORDER BY birth_date
	`)
	if err != nil {
		log.Printf("Error calculating age: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, birthDate string
		var age float64
		err := rows.Scan(&id, &name, &birthDate, &age)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - Born: %s, Age: %.1f", id, name, birthDate, age)
	}

	// Next birthday
	log.Println("✓ Next birthday:")
	rows, err = db.Query(`
		SELECT 
		  id, 
		  name, 
		  strftime('%Y-%m-%d', birth_date, '+' || (strftime('%Y', 'now') - strftime('%Y', birth_date)) || ' years') as next_birthday,
		  CASE
		    WHEN strftime('%Y-%m-%d', birth_date, '+' || (strftime('%Y', 'now') - strftime('%Y', birth_date)) || ' years') < date('now')
		    THEN strftime('%Y-%m-%d', birth_date, '+' || (strftime('%Y', 'now') - strftime('%Y', birth_date) + 1) || ' years')
		    ELSE strftime('%Y-%m-%d', birth_date, '+' || (strftime('%Y', 'now') - strftime('%Y', birth_date)) || ' years')
		  END as actual_next_birthday
		FROM people
		ORDER BY actual_next_birthday
	`)
	if err != nil {
		log.Printf("Error calculating next birthday: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, nextBirthday, actualNextBirthday string
		err := rows.Scan(&id, &name, &nextBirthday, &actualNextBirthday)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - Next birthday: %s", id, name, actualNextBirthday)
	}
}

func testBusinessDays(db *sql.DB) {
	// Calculate business days between two dates
	// Assuming Monday-Friday are business days
	var startDate, endDate string
	startDate = "2024-01-01"
	endDate = "2024-01-31"

	log.Printf("✓ Business days between %s and %s:", startDate, endDate)

	// Count weekdays (Monday-Friday)
	rows, err := db.Query(fmt.Sprintf(`
		SELECT 
		  strftime('%%Y-%%m-%%d', date('%s', '+' || (n-1) || ' days')) as day,
		  strftime('%%w', date('%s', '+' || (n-1) || ' days')) as weekday
		FROM (
		  SELECT 1 as n UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5
		  UNION SELECT 6 UNION SELECT 7 UNION SELECT 8 UNION SELECT 9 UNION SELECT 10
		  UNION SELECT 11 UNION SELECT 12 UNION SELECT 13 UNION SELECT 14 UNION SELECT 15
		  UNION SELECT 16 UNION SELECT 17 UNION SELECT 18 UNION SELECT 19 UNION SELECT 20
		  UNION SELECT 21 UNION SELECT 22 UNION SELECT 23 UNION SELECT 24 UNION SELECT 25
		  UNION SELECT 26 UNION SELECT 27 UNION SELECT 28 UNION SELECT 29 UNION SELECT 30 UNION SELECT 31
		) as days
		WHERE date('%s', '+' || (n-1) || ' days') <= '%s'
	`, startDate, startDate, startDate, endDate))
	if err != nil {
		log.Printf("Error querying business days: %v", err)
		return
	}
	defer rows.Close()

	businessDays := 0
	for rows.Next() {
		var day string
		var weekday string
		err := rows.Scan(&day, &weekday)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		if weekday != "0" && weekday != "6" { // Not Sunday (0) or Saturday (6)
			businessDays++
		}
	}

	log.Printf("✓ Business days count: %d", businessDays)

	// Working days calculation (excluding weekends)
	var workingDays int
	err = db.QueryRow(`
		WITH RECURSIVE days(day) AS (
		  SELECT date('2024-01-01')
		  UNION ALL
		  SELECT date(day, '+1 day')
		  FROM days
		  WHERE day < '2024-01-31'
		)
		SELECT COUNT(*) FROM days
		WHERE strftime('%w', day) NOT IN ('0', '6')
	`).Scan(&workingDays)
	if err != nil {
		log.Printf("Error counting working days: %v", err)
		return
	}
	log.Printf("✓ Working days (using CTE): %d", workingDays)
}

func testDateTimeFormats(db *sql.DB) {
	// Various date formats
	formats := []struct {
		name   string
		format string
	}{
		{"ISO Date", "%Y-%m-%d"},
		{"US Date", "%m/%d/%Y"},
		{"European Date", "%d/%m/%Y"},
		{"Full Date", "%A, %B %d, %Y"},
		{"Short Date", "%Y-%m-%d"},
		{"Month Name", "%B"},
		{"Day Name", "%A"},
		{"Week of Year", "%W"},
		{"Day of Year", "%j"},
	}

	for _, f := range formats {
		var result string
		err := db.QueryRow(fmt.Sprintf("SELECT strftime('%s', 'now')", f.format)).Scan(&result)
		if err != nil {
			log.Printf("Error formatting date (%s): %v", f.name, err)
			continue
		}
		log.Printf("✓ %s: %s", f.name, result)
	}

	// Various time formats
	timeFormats := []struct {
		name   string
		format string
	}{
		{"24 Hour Time", "%H:%M:%S"},
		{"12 Hour Time", "%I:%M:%S %p"},
		{"Hour", "%H"},
		{"Minute", "%M"},
		{"Second", "%S"},
		{"AM/PM", "%p"},
	}

	for _, f := range timeFormats {
		var result string
		err := db.QueryRow(fmt.Sprintf("SELECT strftime('%s', 'now')", f.format)).Scan(&result)
		if err != nil {
			log.Printf("Error formatting time (%s): %v", f.name, err)
			continue
		}
		log.Printf("✓ %s: %s", f.name, result)
	}
}

func testDateTimeWithAggregates(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE sales (id INTEGER, sale_date DATE, amount REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: sales")

	// Insert test data
	sales := []struct {
		id     int
		date   string
		amount float64
	}{
		{1, "2024-01-15", 100.00},
		{2, "2024-01-20", 150.00},
		{3, "2024-02-10", 200.00},
		{4, "2024-02-25", 180.00},
		{5, "2024-03-05", 220.00},
	}

	for _, sale := range sales {
		_, err = db.Exec("INSERT INTO sales VALUES (?, ?, ?)", sale.id, sale.date, sale.amount)
		if err != nil {
			log.Printf("Error inserting sale: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 5 sales")

	// Group by month
	log.Println("✓ Sales grouped by month:")
	rows, err := db.Query(`
		SELECT 
		  strftime('%Y-%m', sale_date) as month,
		  COUNT(*) as sales_count,
		  SUM(amount) as total_sales,
		  AVG(amount) as avg_sale
		FROM sales
		GROUP BY month
		ORDER BY month
	`)
	if err != nil {
		log.Printf("Error grouping sales by month: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var month string
		var salesCount int
		var totalSales, avgSale float64
		err := rows.Scan(&month, &salesCount, &totalSales, &avgSale)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %s: Count=%d, Total=$%.2f, Avg=$%.2f", month, salesCount, totalSales, avgSale)
	}

	// Year to date
	log.Println("✓ Year to date sales:")
	var ytdSales float64
	err = db.QueryRow(`
		SELECT SUM(amount) 
		FROM sales 
		WHERE strftime('%Y', sale_date) = strftime('%Y', 'now')
	`).Scan(&ytdSales)
	if err != nil {
		log.Printf("Error calculating YTD sales: %v", err)
		return
	}
	log.Printf("✓ Year to date sales: $%.2f", ytdSales)
}

func testCleanup(db *sql.DB) {
	tables := []string{
		"events",
		"people",
		"appointments",
		"sales",
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
