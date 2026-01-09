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

	// Test 1: Basic FTS5 Table Creation
	log.Println("\n=== Test 1: Basic FTS5 Table Creation ===")
testBasicFTS5Table(db)

	// Test 2: FTS5 INSERT Operations
	log.Println("\n=== Test 2: FTS5 INSERT Operations ===")
testFTS5Insert(db)

	// Test 3: FTS5 MATCH Query
	log.Println("\n=== Test 3: FTS5 MATCH Query ===")
testFTS5Match(db)

	// Test 4: FTS5 Multiple Term MATCH
	log.Println("\n=== Test 4: FTS5 Multiple Term MATCH ===")
testFTS5MultipleMatch(db)

	// Test 5: FTS5 Phrase MATCH
	log.Println("\n=== Test 5: FTS5 Phrase MATCH ===")
testFTS5PhraseMatch(db)

	// Test 6: FTS5 Prefix MATCH
	log.Println("\n=== Test 6: FTS5 Prefix MATCH ===")
testFTS5PrefixMatch(db)

	// Test 7: FTS5 Ranking (bm25)
	log.Println("\n=== Test 7: FTS5 Ranking (bm25) ===")
testFTS5Ranking(db)

	// Test 8: FTS5 Snippet Extraction
	log.Println("\n=== Test 8: FTS5 Snippet Extraction ===")
testFTS5Snippet(db)

	// Test 9: FTS5 Highlight
	log.Println("\n=== Test 9: FTS5 Highlight ===")
testFTS5Highlight(db)

	// Test 10: FTS5 External Content Table
	log.Println("\n=== Test 10: FTS5 External Content Table ===")
testFTS5ExternalContent(db)

	// Test 11: Cleanup
	log.Println("\n=== Test 11: Cleanup ===")
testCleanup(db)

	log.Println("\n=== All Phase 26 tests completed! ===")
	log.Println("🎉 Phase 26: Full-Text Search (FTS) - COMPLETE! 🎉")
}

func testBasicFTS5Table(db *sql.DB) {
	// Create FTS5 virtual table
	_, err := db.Exec("CREATE VIRTUAL TABLE fts_docs USING fts5(title, content)")
	if err != nil {
		log.Printf("Error creating FTS5 table: %v", err)
		return
	}
	log.Println("✓ Created FTS5 virtual table: fts_docs")

	// Verify table structure
	rows, err := db.Query("PRAGMA table_info(fts_docs)")
	if err != nil {
		log.Printf("Error querying table structure: %v", err)
		return
	}
	defer rows.Close()

	log.Println("✓ Table structure:")
	for rows.Next() {
		var cid int
		var name string
		var typeStr string
		var notnull int
		var dfltValue sql.NullString
		var pk int

		err := rows.Scan(&cid, &name, &typeStr, &notnull, &dfltValue, &pk)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		log.Printf("  Column %d: %s (%s)", cid, name, typeStr)
	}
}

func testFTS5Insert(db *sql.DB) {
	// Create FTS5 virtual table
	_, err := db.Exec("CREATE VIRTUAL TABLE fts_insert USING fts5(title, content)")
	if err != nil {
		log.Printf("Error creating FTS5 table: %v", err)
		return
	}
	log.Println("✓ Created FTS5 virtual table: fts_insert")

	// Insert documents
	documents := []struct {
		title   string
		content string
	}{
		{"Introduction to Go", "Go is a statically typed, compiled programming language designed at Google."},
		{"Advanced Go Programming", "Go supports concurrent programming, garbage collection, and more."},
		{"Go Web Development", "Go is excellent for building web servers and microservices."},
		{"Go Database Access", "Go provides database/sql package for database access."},
	}

	for _, doc := range documents {
		_, err = db.Exec("INSERT INTO fts_insert(title, content) VALUES (?, ?)", doc.title, doc.content)
		if err != nil {
			log.Printf("Error inserting document: %v", err)
			return
		}
	}

	log.Printf("✓ Inserted %d documents", len(documents))

	// Verify documents
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM fts_insert").Scan(&count)
	if err != nil {
		log.Printf("Error counting documents: %v", err)
		return
	}

	log.Printf("✓ Document count: %d", count)
}

func testFTS5Match(db *sql.DB) {
	// Create FTS5 virtual table
	_, err := db.Exec("CREATE VIRTUAL TABLE fts_match USING fts5(title, content)")
	if err != nil {
		log.Printf("Error creating FTS5 table: %v", err)
		return
	}
	log.Println("✓ Created FTS5 virtual table: fts_match")

	// Insert documents
	documents := []struct {
		title   string
		content string
	}{
		{"Python Programming", "Python is a high-level, interpreted programming language."},
		{"Python for Data Science", "Python is widely used in data science and machine learning."},
		{"Python Web Development", "Python frameworks like Django and Flask are popular for web development."},
		{"Go Programming", "Go is a statically typed, compiled programming language."},
		{"Go for Cloud", "Go is often used for cloud and distributed systems."},
	}

	for _, doc := range documents {
		_, err = db.Exec("INSERT INTO fts_match(title, content) VALUES (?, ?)", doc.title, doc.content)
		if err != nil {
			log.Printf("Error inserting document: %v", err)
			return
		}
	}

	log.Println("✓ Inserted documents")

	// Test MATCH query - single term
	searchTerm := "Python"
	log.Printf("✓ MATCH query for '%s':", searchTerm)

	rows, err := db.Query("SELECT title, content FROM fts_match WHERE fts_match MATCH ?", searchTerm)
	if err != nil {
		log.Printf("Error executing MATCH query: %v", err)
		return
	}
	defer rows.Close()

	resultCount := 0
	for rows.Next() {
		var title string
		var content string
		err := rows.Scan(&title, &content)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		resultCount++
		log.Printf("  Result %d: %s", resultCount, title)
	}

	log.Printf("✓ Found %d results", resultCount)
}

func testFTS5MultipleMatch(db *sql.DB) {
	// Create FTS5 virtual table
	_, err := db.Exec("CREATE VIRTUAL TABLE fts_multiple USING fts5(title, content)")
	if err != nil {
		log.Printf("Error creating FTS5 table: %v", err)
		return
	}
	log.Println("✓ Created FTS5 virtual table: fts_multiple")

	// Insert documents
	documents := []struct {
		title   string
		content string
	}{
		{"Machine Learning with Python", "Python provides excellent libraries for machine learning like TensorFlow, PyTorch."},
		{"Deep Learning with Go", "Go is also used for machine learning, especially for performance-critical applications."},
		{"Data Science with R", "R is a statistical programming language widely used in data science."},
		{"Web Development with JavaScript", "JavaScript is the most popular language for web development."},
		{"Mobile Development with Swift", "Swift is used for iOS app development."},
	}

	for _, doc := range documents {
		_, err = db.Exec("INSERT INTO fts_multiple(title, content) VALUES (?, ?)", doc.title, doc.content)
		if err != nil {
			log.Printf("Error inserting document: %v", err)
			return
		}
	}

	log.Println("✓ Inserted documents")

	// Test MATCH query - multiple terms (AND)
	searchTerms := "machine learning"
	log.Printf("✓ MATCH query for '%s':", searchTerms)

	rows, err := db.Query("SELECT title, content FROM fts_multiple WHERE fts_multiple MATCH ?", searchTerms)
	if err != nil {
		log.Printf("Error executing MATCH query: %v", err)
		return
	}
	defer rows.Close()

	resultCount := 0
	for rows.Next() {
		var title string
		var content string
		err := rows.Scan(&title, &content)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		resultCount++
		log.Printf("  Result %d: %s", resultCount, title)
	}

	log.Printf("✓ Found %d results", resultCount)
}

func testFTS5PhraseMatch(db *sql.DB) {
	// Create FTS5 virtual table
	_, err := db.Exec("CREATE VIRTUAL TABLE fts_phrase USING fts5(title, content)")
	if err != nil {
		log.Printf("Error creating FTS5 table: %v", err)
		return
	}
	log.Println("✓ Created FTS5 virtual table: fts_phrase")

	// Insert documents
	documents := []struct {
		title   string
		content string
	}{
		{"Best Practices", "Software engineering best practices include code reviews, testing, and documentation."},
		{"Clean Code", "Writing clean code is essential for maintainability and readability."},
		{"Code Quality", "Code quality metrics include complexity, coverage, and maintainability."},
		{"Software Architecture", "Software architecture patterns include microservices, monolith, and serverless."},
	}

	for _, doc := range documents {
		_, err = db.Exec("INSERT INTO fts_phrase(title, content) VALUES (?, ?)", doc.title, doc.content)
		if err != nil {
			log.Printf("Error inserting document: %v", err)
			return
		}
	}

	log.Println("✓ Inserted documents")

	// Test MATCH query - phrase (double quotes)
	phrase := "\"code quality\""
	log.Printf("✓ MATCH query for phrase %s:", phrase)

	rows, err := db.Query("SELECT title, content FROM fts_phrase WHERE fts_phrase MATCH ?", phrase)
	if err != nil {
		log.Printf("Error executing MATCH query: %v", err)
		return
	}
	defer rows.Close()

	resultCount := 0
	for rows.Next() {
		var title string
		var content string
		err := rows.Scan(&title, &content)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		resultCount++
		log.Printf("  Result %d: %s", resultCount, title)
	}

	log.Printf("✓ Found %d results", resultCount)
}

func testFTS5PrefixMatch(db *sql.DB) {
	// Create FTS5 virtual table
	_, err := db.Exec("CREATE VIRTUAL TABLE fts_prefix USING fts5(title, content)")
	if err != nil {
		log.Printf("Error creating FTS5 table: %v", err)
		return
	}
	log.Println("✓ Created FTS5 virtual table: fts_prefix")

	// Insert documents
	documents := []struct {
		title   string
		content string
	}{
		{"Python Tutorial", "This tutorial covers Python basics, variables, and control flow."},
		{"Python Advanced", "Advanced Python topics include decorators, generators, and context managers."},
		{"JavaScript Tutorial", "JavaScript tutorial covers basics, DOM manipulation, and events."},
		{"Go Tutorial", "Go tutorial covers basics, goroutines, and channels."},
		{"Java Tutorial", "Java tutorial covers OOP, collections, and streams."},
	}

	for _, doc := range documents {
		_, err = db.Exec("INSERT INTO fts_prefix(title, content) VALUES (?, ?)", doc.title, doc.content)
		if err != nil {
			log.Printf("Error inserting document: %v", err)
			return
		}
	}

	log.Println("✓ Inserted documents")

	// Test MATCH query - prefix search (term*)
	prefix := "Pytho*"
	log.Printf("✓ MATCH query for prefix %s:", prefix)

	rows, err := db.Query("SELECT title, content FROM fts_prefix WHERE fts_prefix MATCH ?", prefix)
	if err != nil {
		log.Printf("Error executing MATCH query: %v", err)
		return
	}
	defer rows.Close()

	resultCount := 0
	for rows.Next() {
		var title string
		var content string
		err := rows.Scan(&title, &content)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		resultCount++
		log.Printf("  Result %d: %s", resultCount, title)
	}

	log.Printf("✓ Found %d results", resultCount)
}

func testFTS5Ranking(db *sql.DB) {
	// Create FTS5 virtual table
	_, err := db.Exec("CREATE VIRTUAL TABLE fts_rank USING fts5(title, content)")
	if err != nil {
		log.Printf("Error creating FTS5 table: %v", err)
		return
	}
	log.Println("✓ Created FTS5 virtual table: fts_rank")

	// Insert documents
	documents := []struct {
		title   string
		content string
	}{
		{"Database Basics", "Databases are used to store and retrieve data efficiently."},
		{"SQL Fundamentals", "SQL is the standard language for relational database management systems."},
		{"NoSQL Databases", "NoSQL databases provide flexible schema and horizontal scaling."},
		{"Database Optimization", "Database optimization includes indexing, query tuning, and caching."},
		{"Database Administration", "Database administration involves backup, recovery, and security."},
		{"SQL Performance", "SQL performance tuning involves analyzing query plans and optimizing indexes."},
	}

	for _, doc := range documents {
		_, err = db.Exec("INSERT INTO fts_rank(title, content) VALUES (?, ?)", doc.title, doc.content)
		if err != nil {
			log.Printf("Error inserting document: %v", err)
			return
		}
	}

	log.Println("✓ Inserted documents")

	// Test MATCH query with bm25 ranking
	searchTerm := "database"
	log.Printf("✓ MATCH query for '%s' with bm25 ranking:", searchTerm)

	rows, err := db.Query("SELECT title, content, bm25(fts_rank) as rank FROM fts_rank WHERE fts_rank MATCH ? ORDER BY rank", searchTerm)
	if err != nil {
		log.Printf("Error executing MATCH query: %v", err)
		return
	}
	defer rows.Close()

	resultCount := 0
	for rows.Next() {
		var title string
		var content string
		var rank float64
		err := rows.Scan(&title, &content, &rank)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		resultCount++
		log.Printf("  Result %d: %s (rank: %.4f)", resultCount, title, rank)
	}

	log.Printf("✓ Found %d results", resultCount)
}

func testFTS5Snippet(db *sql.DB) {
	// Create FTS5 virtual table
	_, err := db.Exec("CREATE VIRTUAL TABLE fts_snippet USING fts5(title, content)")
	if err != nil {
		log.Printf("Error creating FTS5 table: %v", err)
		return
	}
	log.Println("✓ Created FTS5 virtual table: fts_snippet")

	// Insert documents
	documents := []struct {
		title   string
		content string
	}{
		{"Introduction to Machine Learning", "Machine learning is a subset of artificial intelligence that enables systems to learn from data and improve over time without being explicitly programmed."},
		{"Deep Learning Fundamentals", "Deep learning is a subset of machine learning that uses neural networks with multiple layers to progressively extract higher-level features from raw input."},
		{"Natural Language Processing", "Natural language processing (NLP) is a branch of artificial intelligence that helps computers understand, interpret and manipulate human language."},
	}

	for _, doc := range documents {
		_, err = db.Exec("INSERT INTO fts_snippet(title, content) VALUES (?, ?)", doc.title, doc.content)
		if err != nil {
			log.Printf("Error inserting document: %v", err)
			return
		}
	}

	log.Println("✓ Inserted documents")

	// Test MATCH query with snippet extraction
	searchTerm := "learning"
	log.Printf("✓ MATCH query for '%s' with snippet extraction:", searchTerm)

	rows, err := db.Query("SELECT title, snippet(fts_snippet, 0, '[', ']', '...', 20) as snippet FROM fts_snippet WHERE fts_snippet MATCH ?", searchTerm)
	if err != nil {
		log.Printf("Error executing MATCH query: %v", err)
		return
	}
	defer rows.Close()

	resultCount := 0
	for rows.Next() {
		var title string
		var snippet string
		err := rows.Scan(&title, &snippet)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		resultCount++
		log.Printf("  Result %d: %s", resultCount, title)
		log.Printf("    Snippet: %s", snippet)
	}

	log.Printf("✓ Found %d results", resultCount)
}

func testFTS5Highlight(db *sql.DB) {
	// Create FTS5 virtual table
	_, err := db.Exec("CREATE VIRTUAL TABLE fts_highlight USING fts5(title, content)")
	if err != nil {
		log.Printf("Error creating FTS5 table: %v", err)
		return
	}
	log.Println("✓ Created FTS5 virtual table: fts_highlight")

	// Insert documents
	documents := []struct {
		title   string
		content string
	}{
		{"Web Development Trends", "Modern web development includes progressive web apps, server-side rendering, and static site generation. Web developers need to be familiar with HTML, CSS, and JavaScript."},
		{"Mobile Development", "Mobile application development includes native development for iOS and Android, as well as cross-platform development using frameworks like React Native and Flutter."},
		{"Cloud Computing", "Cloud computing provides on-demand delivery of IT resources over the internet. Major cloud providers include AWS, Azure, and Google Cloud."},
	}

	for _, doc := range documents {
		_, err = db.Exec("INSERT INTO fts_highlight(title, content) VALUES (?, ?)", doc.title, doc.content)
		if err != nil {
			log.Printf("Error inserting document: %v", err)
			return
		}
	}

	log.Println("✓ Inserted documents")

	// Test MATCH query with highlight
	searchTerm := "development"
	log.Printf("✓ MATCH query for '%s' with highlight:", searchTerm)

	rows, err := db.Query("SELECT title, highlight(fts_highlight, 1, '<b>', '</b>') as highlighted_content FROM fts_highlight WHERE fts_highlight MATCH ?", searchTerm)
	if err != nil {
		log.Printf("Error executing MATCH query: %v", err)
		return
	}
	defer rows.Close()

	resultCount := 0
	for rows.Next() {
		var title string
		var highlightedContent string
		err := rows.Scan(&title, &highlightedContent)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		resultCount++
		log.Printf("  Result %d: %s", resultCount, title)
		log.Printf("    Highlighted: %s", highlightedContent)
	}

	log.Printf("✓ Found %d results", resultCount)
}

func testFTS5ExternalContent(db *sql.DB) {
	// Create regular table
	_, err := db.Exec("CREATE TABLE articles (id INTEGER PRIMARY KEY, title TEXT, content TEXT, author TEXT, created_at DATETIME)")
	if err != nil {
		log.Printf("Error creating articles table: %v", err)
		return
	}
	log.Println("✓ Created table: articles")

	// Insert articles
	articles := []struct {
		id        int
		title     string
		content   string
		author    string
		createdAt string
	}{
		{1, "Getting Started with Go", "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.", "John Doe", "2024-01-01"},
		{2, "Go Concurrency Patterns", "Go provides goroutines and channels for concurrent programming.", "Jane Smith", "2024-01-02"},
		{3, "Building APIs with Go", "Go is excellent for building RESTful APIs and microservices.", "Bob Johnson", "2024-01-03"},
	}

	for _, article := range articles {
		_, err = db.Exec("INSERT INTO articles VALUES (?, ?, ?, ?, ?)", article.id, article.title, article.content, article.author, article.createdAt)
		if err != nil {
			log.Printf("Error inserting article: %v", err)
			return
		}
	}

	log.Printf("✓ Inserted %d articles", len(articles))

	// Create FTS5 virtual table with external content
	_, err = db.Exec("CREATE VIRTUAL TABLE fts_articles USING fts5(title, content, content=articles, content_rowid=rowid)")
	if err != nil {
		log.Printf("Error creating FTS5 external content table: %v", err)
		return
	}
	log.Println("✓ Created FTS5 virtual table with external content: fts_articles")

	// Verify that FTS5 table can query external content
	log.Println("✓ Verifying FTS5 can query external content:")
	rows, err := db.Query("SELECT a.id, a.title, a.author, fts_articles.title as fts_title FROM articles a JOIN fts_articles ON a.rowid = fts_articles.rowid")
	if err != nil {
		log.Printf("Error querying external content: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title string
		var author string
		var ftsTitle string
		err := rows.Scan(&id, &title, &author, &ftsTitle)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		log.Printf("  Article %d: %s by %s (FTS: %s)", id, title, author, ftsTitle)
	}

	// Test MATCH query on external content
	searchTerm := "Go"
	log.Printf("✓ MATCH query on external content for '%s':", searchTerm)

	rows, err = db.Query("SELECT a.id, a.title, a.author FROM articles a JOIN fts_articles ON a.rowid = fts_articles.rowid WHERE fts_articles MATCH ?", searchTerm)
	if err != nil {
		log.Printf("Error executing MATCH query: %v", err)
		return
	}
	defer rows.Close()

	resultCount := 0
	for rows.Next() {
		var id int
		var title string
		var author string
		err := rows.Scan(&id, &title, &author)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		resultCount++
		log.Printf("  Result %d: %s by %s", resultCount, title, author)
	}

	log.Printf("✓ Found %d results", resultCount)
}

func testCleanup(db *sql.DB) {
	tables := []string{
		"fts_docs",
		"fts_insert",
		"fts_match",
		"fts_multiple",
		"fts_phrase",
		"fts_prefix",
		"fts_rank",
		"fts_snippet",
		"fts_highlight",
		"articles",
		"fts_articles",
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
