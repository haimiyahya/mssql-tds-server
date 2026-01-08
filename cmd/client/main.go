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

	// Execute a simple query
	rows, err := db.Query("SELECT hello FROM world")
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return
	}
	defer rows.Close()

	// Read results
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("Error getting columns: %v", err)
		return
	}

	log.Printf("Columns: %v", columns)

	for rows.Next() {
		var result string
		err := rows.Scan(&result)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("Result: %s", result)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("Rows error: %v", err)
	}

	log.Println("Client test completed successfully!")
}
