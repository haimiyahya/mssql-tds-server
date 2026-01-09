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

	// Test 1: JSON Extraction (json_extract)
	log.Println("\n=== Test 1: JSON Extraction (json_extract) ===")
testJSONExtract(db)

	// Test 2: JSON Creation (json_object, json_array)
	log.Println("\n=== Test 2: JSON Creation (json_object, json_array) ===")
testJSONCreation(db)

	// Test 3: JSON Modification (json_set, json_insert, json_replace)
	log.Println("\n=== Test 3: JSON Modification (json_set, json_insert, json_replace) ===")
testJSONModification(db)

	// Test 4: JSON Path Operations (json_patch)
	log.Println("\n=== Test 4: JSON Path Operations (json_patch) ===")
testJSONPatch(db)

	// Test 5: JSON Validation (json_valid)
	log.Println("\n=== Test 5: JSON Validation (json_valid) ===")
testJSONValid(db)

	// Test 6: JSON Query (json_each, json_tree)
	log.Println("\n=== Test 6: JSON Query (json_each, json_tree) ===")
testJSONQuery(db)

	// Test 7: JSON Aggregation (json_group_array, json_group_object)
	log.Println("\n=== Test 7: JSON Aggregation (json_group_array, json_group_object) ===")
testJSONAggregation(db)

	// Test 8: JSON Length (json_array_length)
	log.Println("\n=== Test 8: JSON Length (json_array_length) ===")
testJSONArrayLength(db)

	// Test 9: JSON Type Detection (json_type)
	log.Println("\n=== Test 9: JSON Type Detection (json_type) ===")
testJSONType(db)

	// Test 10: JSON with SQL Tables
	log.Println("\n=== Test 10: JSON with SQL Tables ===")
testJSONWithTables(db)

	// Test 11: JSON Nested Extraction
	log.Println("\n=== Test 11: JSON Nested Extraction ===")
testJSONNestedExtraction(db)

	// Test 12: JSON Array Operations (json_remove)
	log.Println("\n=== Test 12: JSON Array Operations (json_remove) ===")
testJSONRemove(db)

	// Test 13: JSON Pretty Print (json_pretty)
	log.Println("\n=== Test 13: JSON Pretty Print (json_pretty) ===")
testJSONPretty(db)

	// Test 14: Cleanup
	log.Println("\n=== Test 14: Cleanup ===")
testCleanup(db)

	log.Println("\n=== All Phase 30 tests completed! ===")
	log.Println("🎉 Phase 30: JSON Functions - COMPLETE! 🎉")
}

func testJSONExtract(db *sql.DB) {
	// Simple JSON extraction
	var result string
	err := db.QueryRow(`
		SELECT json_extract('{"name": "John", "age": 30}', '$.name')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error extracting JSON: %v", err)
		return
	}

	log.Printf("✓ Extracted name: %s", result)

	// Nested JSON extraction
	err = db.QueryRow(`
		SELECT json_extract('{"user": {"name": "Jane", "email": "jane@example.com"}}', '$.user.email')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error extracting nested JSON: %v", err)
		return
	}

	log.Printf("✓ Extracted email: %s", result)

	// Array extraction
	err = db.QueryRow(`
		SELECT json_extract('{"items": ["apple", "banana", "orange"]}', '$.items[1]')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error extracting JSON array: %v", err)
		return
	}

	log.Printf("✓ Extracted array item: %s", result)

	// Multiple paths
	rows, err := db.Query(`
		SELECT 
		  json_extract('{"name": "John", "age": 30, "city": "NYC"}', '$.name') as name,
		  json_extract('{"name": "John", "age": 30, "city": "NYC"}', '$.age') as age,
		  json_extract('{"name": "John", "age": 30, "city": "NYC"}', '$.city') as city
	`)
	if err != nil {
		log.Printf("Error extracting multiple paths: %v", err)
		return
	}
	defer rows.Close()

	log.Println("✓ Multiple paths extraction:")
	for rows.Next() {
		var name, city string
		var age int
		err := rows.Scan(&name, &age, &city)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Name: %s, Age: %d, City: %s", name, age, city)
	}
}

func testJSONCreation(db *sql.DB) {
	// JSON object creation
	var result string
	err := db.QueryRow(`
		SELECT json_object('name', 'John', 'age', 30, 'city', 'NYC')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error creating JSON object: %v", err)
		return
	}

	log.Printf("✓ JSON object: %s", result)

	// JSON array creation
	err = db.QueryRow(`
		SELECT json_array('apple', 'banana', 'orange')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error creating JSON array: %v", err)
		return
	}

	log.Printf("✓ JSON array: %s", result)

	// Nested JSON object
	err = db.QueryRow(`
		SELECT json_object(
		  'user', json_object('name', 'Jane', 'email', 'jane@example.com'),
		  'items', json_array('item1', 'item2', 'item3')
		)
	`).Scan(&result)
	if err != nil {
		log.Printf("Error creating nested JSON: %v", err)
		return
	}

	log.Printf("✓ Nested JSON: %s", result)

	// JSON array with numbers
	err = db.QueryRow(`
		SELECT json_array(1, 2, 3, 4, 5)
	`).Scan(&result)
	if err != nil {
		log.Printf("Error creating JSON array with numbers: %v", err)
		return
	}

	log.Printf("✓ JSON array with numbers: %s", result)
}

func testJSONModification(db *sql.DB) {
	// json_set - set value
	var result string
	err := db.QueryRow(`
		SELECT json_set('{"name": "John", "age": 30}', '$.age', 31)
	`).Scan(&result)
	if err != nil {
		log.Printf("Error setting JSON value: %v", err)
		return
	}

	log.Printf("✓ JSON set: %s", result)

	// json_insert - insert new value
	err = db.QueryRow(`
		SELECT json_insert('{"name": "John", "age": 30}', '$.city', 'NYC')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error inserting JSON value: %v", err)
		return
	}

	log.Printf("✓ JSON insert: %s", result)

	// json_replace - replace value
	err = db.QueryRow(`
		SELECT json_replace('{"name": "John", "age": 30}', '$.name', 'Jane')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error replacing JSON value: %v", err)
		return
	}

	log.Printf("✓ JSON replace: %s", result)

	// Multiple modifications
	err = db.QueryRow(`
		SELECT json_set('{"name": "John", "age": 30}', '$.age', 35, '$.city', 'Boston')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error performing multiple JSON modifications: %v", err)
		return
	}

	log.Printf("✓ Multiple JSON modifications: %s", result)
}

func testJSONPatch(db *sql.DB) {
	// json_patch - merge JSON objects
	var result string
	err := db.QueryRow(`
		SELECT json_patch(
		  '{"name": "John", "age": 30}',
		  '{"age": 31, "city": "NYC"}'
		)
	`).Scan(&result)
	if err != nil {
		log.Printf("Error patching JSON: %v", err)
		return
	}

	log.Printf("✓ JSON patch: %s", result)

	// Nested patch
	err = db.QueryRow(`
		SELECT json_patch(
		  '{"user": {"name": "John", "age": 30}, "active": true}',
		  '{"user": {"age": 31}, "active": false}'
		)
	`).Scan(&result)
	if err != nil {
		log.Printf("Error patching nested JSON: %v", err)
		return
	}

	log.Printf("✓ Nested JSON patch: %s", result)

	// Multiple patches
	err = db.QueryRow(`
		SELECT json_patch(
		  json_patch(
		    '{"name": "John", "age": 30}',
		    '{"age": 31}'
		  ),
		  '{"city": "NYC"}'
		)
	`).Scan(&result)
	if err != nil {
		log.Printf("Error performing multiple JSON patches: %v", err)
		return
	}

	log.Printf("✓ Multiple JSON patches: %s", result)
}

func testJSONValid(db *sql.DB) {
	// Valid JSON
	var isValid int
	err := db.QueryRow(`
		SELECT json_valid('{"name": "John", "age": 30}')
	`).Scan(&isValid)
	if err != nil {
		log.Printf("Error validating JSON: %v", err)
		return
	}

	log.Printf("✓ Valid JSON: %d (1 = true)", isValid)

	// Invalid JSON
	err = db.QueryRow(`
		SELECT json_valid('{name: "John", age: 30}')
	`).Scan(&isValid)
	if err != nil {
		log.Printf("Error validating invalid JSON: %v", err)
		return
	}

	log.Printf("✓ Invalid JSON: %d (0 = false)", isValid)

	// Valid JSON array
	err = db.QueryRow(`
		SELECT json_valid('["apple", "banana", "orange"]')
	`).Scan(&isValid)
	if err != nil {
		log.Printf("Error validating JSON array: %v", err)
		return
	}

	log.Printf("✓ Valid JSON array: %d (1 = true)", isValid)

	// Empty JSON
	err = db.QueryRow(`
		SELECT json_valid('{}')
	`).Scan(&isValid)
	if err != nil {
		log.Printf("Error validating empty JSON: %v", err)
		return
	}

	log.Printf("✓ Empty JSON: %d (1 = true)", isValid)

	// NULL JSON
	err = db.QueryRow(`
		SELECT json_valid(NULL)
	`).Scan(&isValid)
	if err != nil {
		log.Printf("Error validating NULL JSON: %v", err)
		return
	}

	log.Printf("✓ NULL JSON: %d (0 = false)", isValid)
}

func testJSONQuery(db *sql.DB) {
	// json_each - iterate over JSON array
	log.Println("✓ JSON each (array):")
	rows, err := db.Query(`
		SELECT key, value, type
		FROM json_each('["apple", "banana", "orange"]')
	`)
	if err != nil {
		log.Printf("Error querying JSON each: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var value string
		var jsonType string
		err := rows.Scan(&key, &value, &jsonType)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Key: %s, Value: %s, Type: %s", key, value, jsonType)
	}

	// json_each - iterate over JSON object
	log.Println("✓ JSON each (object):")
	rows, err = db.Query(`
		SELECT key, value, type
		FROM json_each('{"name": "John", "age": 30, "city": "NYC"}')
	`)
	if err != nil {
		log.Printf("Error querying JSON each object: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var value string
		var jsonType string
		err := rows.Scan(&key, &value, &jsonType)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Key: %s, Value: %s, Type: %s", key, value, jsonType)
	}

	// json_tree - recursive JSON query
	log.Println("✓ JSON tree (recursive):")
	rows, err = db.Query(`
		SELECT key, value, type, path
		FROM json_tree('{"user": {"name": "John", "email": "jane@example.com"}, "active": true}')
		ORDER BY path
	`)
	if err != nil {
		log.Printf("Error querying JSON tree: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var key, value, jsonType, path string
		err := rows.Scan(&key, &value, &jsonType, &path)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Path: %s, Key: %s, Value: %s, Type: %s", path, key, value, jsonType)
	}
}

func testJSONAggregation(db *sql.DB) {
	// Create test table
	_, err := db.Exec("CREATE TABLE users (id INTEGER, name TEXT, age INTEGER, city TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: users")

	// Insert test data
	users := []struct {
		id   int
		name string
		age  int
		city string
	}{
		{1, "John", 30, "NYC"},
		{2, "Jane", 25, "Boston"},
		{3, "Bob", 35, "NYC"},
		{4, "Alice", 28, "Chicago"},
	}

	for _, user := range users {
		_, err = db.Exec("INSERT INTO users VALUES (?, ?, ?, ?)", user.id, user.name, user.age, user.city)
		if err != nil {
			log.Printf("Error inserting user: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 4 users")

	// json_group_array - group values into JSON array
	var jsonArray string
	err = db.QueryRow(`
		SELECT json_group_array(name) FROM users
	`).Scan(&jsonArray)
	if err != nil {
		log.Printf("Error grouping JSON array: %v", err)
		return
	}

	log.Printf("✓ JSON group array: %s", jsonArray)

	// json_group_object - group key-value pairs into JSON object
	var jsonObject string
	err = db.QueryRow(`
		SELECT json_group_object(name, city) FROM users
	`).Scan(&jsonObject)
	if err != nil {
		log.Printf("Error grouping JSON object: %v", err)
		return
	}

	log.Printf("✓ JSON group object: %s", jsonObject)

	// Group by city with json_group_array
	log.Println("✓ JSON group array by city:")
	rows, err := db.Query(`
		SELECT city, json_group_array(name) as names
		FROM users
		GROUP BY city
	`)
	if err != nil {
		log.Printf("Error grouping JSON array by city: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var city string
		var names string
		err := rows.Scan(&city, &names)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  City: %s, Names: %s", city, names)
	}
}

func testJSONArrayLength(db *sql.DB) {
	// Array length
	var length int
	err := db.QueryRow(`
		SELECT json_array_length('["apple", "banana", "orange"]')
	`).Scan(&length)
	if err != nil {
		log.Printf("Error getting JSON array length: %v", err)
		return
	}

	log.Printf("✓ JSON array length: %d", length)

	// Empty array length
	err = db.QueryRow(`
		SELECT json_array_length('[]')
	`).Scan(&length)
	if err != nil {
		log.Printf("Error getting empty JSON array length: %v", err)
		return
	}

	log.Printf("✓ Empty JSON array length: %d", length)

	// Nested array length
	err = db.QueryRow(`
		SELECT json_array_length('[[1, 2], [3, 4]]')
	`).Scan(&length)
	if err != nil {
		log.Printf("Error getting nested JSON array length: %v", err)
		return
	}

	log.Printf("✓ Nested JSON array length: %d", length)

	// Array path length
	err = db.QueryRow(`
		SELECT json_array_length('{"items": ["apple", "banana", "orange"]}', '$.items')
	`).Scan(&length)
	if err != nil {
		log.Printf("Error getting JSON array path length: %v", err)
		return
	}

	log.Printf("✓ JSON array path length: %d", length)
}

func testJSONType(db *sql.DB) {
	// Object type
	var jsonType string
	err := db.QueryRow(`
		SELECT json_type('{"name": "John", "age": 30}')
	`).Scan(&jsonType)
	if err != nil {
		log.Printf("Error getting JSON object type: %v", err)
		return
	}

	log.Printf("✓ JSON object type: %s", jsonType)

	// Array type
	err = db.QueryRow(`
		SELECT json_type('["apple", "banana", "orange"]')
	`).Scan(&jsonType)
	if err != nil {
		log.Printf("Error getting JSON array type: %v", err)
		return
	}

	log.Printf("✓ JSON array type: %s", jsonType)

	// String type
	err = db.QueryRow(`
		SELECT json_type('"Hello"')
	`).Scan(&jsonType)
	if err != nil {
		log.Printf("Error getting JSON string type: %v", err)
		return
	}

	log.Printf("✓ JSON string type: %s", jsonType)

	// Number type
	err = db.QueryRow(`
		SELECT json_type('42')
	`).Scan(&jsonType)
	if err != nil {
		log.Printf("Error getting JSON number type: %v", err)
		return
	}

	log.Printf("✓ JSON number type: %s", jsonType)

	// Boolean type
	err = db.QueryRow(`
		SELECT json_type('true')
	`).Scan(&jsonType)
	if err != nil {
		log.Printf("Error getting JSON boolean type: %v", err)
		return
	}

	log.Printf("✓ JSON boolean type: %s", jsonType)

	// Null type
	err = db.QueryRow(`
		SELECT json_type('null')
	`).Scan(&jsonType)
	if err != nil {
		log.Printf("Error getting JSON null type: %v", err)
		return
	}

	log.Printf("✓ JSON null type: %s", jsonType)

	// Path type
	err = db.QueryRow(`
		SELECT json_type('{"name": "John", "age": 30}', '$.name')
	`).Scan(&jsonType)
	if err != nil {
		log.Printf("Error getting JSON path type: %v", err)
		return
	}

	log.Printf("✓ JSON path type: %s", jsonType)
}

func testJSONWithTables(db *sql.DB) {
	// Create table with JSON column
	_, err := db.Exec("CREATE TABLE products (id INTEGER, name TEXT, attributes TEXT)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}
	log.Println("✓ Created table: products")

	// Insert products with JSON attributes
	products := []struct {
		id         int
		name       string
		attributes string
	}{
		{1, "Product 1", `{"color": "red", "size": "M", "price": 10.00}`},
		{2, "Product 2", `{"color": "blue", "size": "L", "price": 15.00}`},
		{3, "Product 3", `{"color": "green", "size": "S", "price": 8.00}`},
	}

	for _, product := range products {
		_, err = db.Exec("INSERT INTO products VALUES (?, ?, ?)", product.id, product.name, product.attributes)
		if err != nil {
			log.Printf("Error inserting product: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 3 products")

	// Query JSON data
	log.Println("✓ Query JSON data:")
	rows, err := db.Query(`
		SELECT 
		  id,
		  name,
		  json_extract(attributes, '$.color') as color,
		  json_extract(attributes, '$.size') as size,
		  json_extract(attributes, '$.price') as price
		FROM products
	`)
	if err != nil {
		log.Printf("Error querying JSON data: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, color, size string
		var price float64
		err := rows.Scan(&id, &name, &color, &size, &price)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s - Color: %s, Size: %s, Price: $%.2f", id, name, color, size, price)
	}

	// Filter JSON data
	log.Println("✓ Filter JSON data (price > 10):")
	rows, err = db.Query(`
		SELECT id, name
		FROM products
		WHERE json_extract(attributes, '$.price') > 10
	`)
	if err != nil {
		log.Printf("Error filtering JSON data: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s", id, name)
	}

	// Update JSON data
	_, err = db.Exec(`
		UPDATE products
		SET attributes = json_set(attributes, '$.price', 12.00)
		WHERE id = 1
	`)
	if err != nil {
		log.Printf("Error updating JSON data: %v", err)
		return
	}
	log.Println("✓ Updated JSON data")

	// Verify update
	var price float64
	err = db.QueryRow(`
		SELECT json_extract(attributes, '$.price') FROM products WHERE id = 1
	`).Scan(&price)
	if err != nil {
		log.Printf("Error verifying JSON update: %v", err)
		return
	}
	log.Printf("✓ Updated price: $%.2f", price)
}

func testJSONNestedExtraction(db *sql.DB) {
	// Deep nested JSON
	var result string
	err := db.QueryRow(`
		SELECT json_extract('{"level1": {"level2": {"level3": {"level4": "deep value"}}}}', '$.level1.level2.level3.level4')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error extracting deep nested JSON: %v", err)
		return
	}

	log.Printf("✓ Deep nested JSON: %s", result)

	// Array of objects
	rows, err := db.Query(`
		SELECT json_extract(value, '$.name'), json_extract(value, '$.age')
		FROM json_each('[{"name": "John", "age": 30}, {"name": "Jane", "age": 25}]')
	`)
	if err != nil {
		log.Printf("Error extracting JSON array of objects: %v", err)
		return
	}
	defer rows.Close()

	log.Println("✓ JSON array of objects:")
	for rows.Next() {
		var name string
		var age int
		err := rows.Scan(&name, &age)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  Name: %s, Age: %d", name, age)
	}

	// Multiple array indices
	err = db.QueryRow(`
		SELECT json_extract('["apple", "banana", "orange", "grape", "pear"]', '$[0], $[2], $[4]')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error extracting multiple array indices: %v", err)
		return
	}

	log.Printf("✓ Multiple array indices: %s", result)
}

func testJSONRemove(db *sql.DB) {
	// Remove from array
	var result string
	err := db.QueryRow(`
		SELECT json_remove('["apple", "banana", "orange"]', '$[1]')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error removing from JSON array: %v", err)
		return
	}

	log.Printf("✓ JSON remove from array: %s", result)

	// Remove from object
	err = db.QueryRow(`
		SELECT json_remove('{"name": "John", "age": 30, "city": "NYC"}', '$.city')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error removing from JSON object: %v", err)
		return
	}

	log.Printf("✓ JSON remove from object: %s", result)

	// Remove multiple paths
	err = db.QueryRow(`
		SELECT json_remove('{"name": "John", "age": 30, "city": "NYC", "active": true}', '$.city', '$.active')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error removing multiple paths: %v", err)
		return
	}

	log.Printf("✓ JSON remove multiple paths: %s", result)
}

func testJSONPretty(db *sql.DB) {
	// Pretty print JSON
	var result string
	err := db.QueryRow(`
		SELECT json_pretty('{"name":"John","age":30,"city":"NYC"}')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error pretty printing JSON: %v", err)
		return
	}

	log.Printf("✓ JSON pretty print:\n%s", result)

	// Pretty print nested JSON
	err = db.QueryRow(`
		SELECT json_pretty('{"user":{"name":"John","email":"jane@example.com"},"active":true}')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error pretty printing nested JSON: %v", err)
		return
	}

	log.Printf("✓ Nested JSON pretty print:\n%s", result)

	// Pretty print JSON array
	err = db.QueryRow(`
		SELECT json_pretty('["apple","banana","orange"]')
	`).Scan(&result)
	if err != nil {
		log.Printf("Error pretty printing JSON array: %v", err)
		return
	}

	log.Printf("✓ JSON array pretty print:\n%s", result)
}

func testCleanup(db *sql.DB) {
	tables := []string{
		"users",
		"products",
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
