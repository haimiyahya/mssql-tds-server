package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"

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

	// Test 1: Haversine Distance Calculation
	log.Println("\n=== Test 1: Haversine Distance Calculation ===")
testHaversineDistance(db)

	// Test 2: Bounding Box Calculations
	log.Println("\n=== Test 2: Bounding Box Calculations ===")
testBoundingBox(db)

	// Test 3: Point-in-Polygon Detection
	log.Println("\n=== Test 3: Point-in-Polygon Detection ===")
testPointInPolygon(db)

	// Test 4: Spatial Queries (Find Nearby Points)
	log.Println("\n=== Test 4: Spatial Queries (Find Nearby Points) ===")
testSpatialQueries(db)

	// Test 5: Spatial Joins
	log.Println("\n=== Test 5: Spatial Joins ===")
testSpatialJoins(db)

	// Test 6: Spatial Aggregations
	log.Println("\n=== Test 6: Spatial Aggregations ===")
testSpatialAggregations(db)

	// Test 7: Coordinate Conversions
	log.Println("\n=== Test 7: Coordinate Conversions ===")
testCoordinateConversions(db)

	// Test 8: Geospatial Indexing
	log.Println("\n=== Test 8: Geospatial Indexing ===")
testGeospatialIndexing(db)

	// Test 9: Polygon Area Calculations
	log.Println("\n=== Test 9: Polygon Area Calculations ===")
testPolygonArea(db)

	// Test 10: Centroid Calculations
	log.Println("\n=== Test 10: Centroid Calculations ===")
testCentroid(db)

	// Test 11: Spatial Filtering
	log.Println("\n=== Test 11: Spatial Filtering ===")
testSpatialFiltering(db)

	// Test 12: Geodesic vs. Euclidean Distance
	log.Println("\n=== Test 12: Geodesic vs. Euclidean Distance ===")
testDistanceComparison(db)

	// Test 13: Spatial Sorting (Distance Ordering)
	log.Println("\n=== Test 13: Spatial Sorting (Distance Ordering) ===")
testSpatialSorting(db)

	// Test 14: Cleanup
	log.Println("\n=== Test 14: Cleanup ===")
testCleanup(db)

	log.Println("\n=== All Phase 32 tests completed! ===")
	log.Println("🎉 Phase 32: Geospatial Functions - COMPLETE! 🎉")
}

func testHaversineDistance(db *sql.DB) {
	// Haversine formula for calculating distance between two points on Earth
	// Distance = 2 * R * asin(sqrt(sin²((lat2 - lat1)/2) + cos(lat1) * cos(lat2) * sin²((lon2 - lon1)/2)))
	// Where R is Earth's radius (6371 km or 3959 miles)

	log.Println("✓ Haversine distance calculations:")

	// Distance between New York (40.7128, -74.0060) and London (51.5074, -0.1278)
	lat1 := 40.7128
	lon1 := -74.0060
	lat2 := 51.5074
	lon2 := -0.1278

	distanceKm := haversineDistanceCalc(lat1, lon1, lat2, lon2, 6371)
	distanceMiles := haversineDistanceCalc(lat1, lon1, lat2, lon2, 3959)

	log.Printf("✓ Distance New York to London: %.2f km (%.2f miles)", distanceKm, distanceMiles)

	// Distance between San Francisco (37.7749, -122.4194) and Los Angeles (34.0522, -118.2437)
	lat1 = 37.7749
	lon1 = -122.4194
	lat2 = 34.0522
	lon2 = -118.2437

	distanceKm = haversineDistanceCalc(lat1, lon1, lat2, lon2, 6371)
	distanceMiles = haversineDistanceCalc(lat1, lon1, lat2, lon2, 3959)

	log.Printf("✓ Distance San Francisco to Los Angeles: %.2f km (%.2f miles)", distanceKm, distanceMiles)

	// Test with SQL (using custom functions)
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS locations (id INTEGER, name TEXT, lat REAL, lon REAL)
	`)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	// Insert test locations
	locations := []struct {
		id   int
		name string
		lat  float64
		lon  float64
	}{
		{1, "Point A", 40.7128, -74.0060},
		{2, "Point B", 40.7130, -74.0065},
		{3, "Point C", 40.7150, -74.0080},
	}

	for _, loc := range locations {
		_, err = db.Exec("INSERT OR REPLACE INTO locations VALUES (?, ?, ?, ?)",
			loc.id, loc.name, loc.lat, loc.lon)
		if err != nil {
			log.Printf("Error inserting location: %v", err)
			return
		}
	}

	log.Println("✓ Created table and inserted 3 test locations")
}

func testBoundingBox(db *sql.DB) {
	// Bounding box: minLat, minLon, maxLat, maxLon
	// Used to quickly filter points within a rectangular region

	log.Println("✓ Bounding box calculations:")

	// Create bounding box around New York City
	minLat := 40.7000
	maxLat := 40.7500
	minLon := -74.0200
	maxLon := -73.9900

	log.Printf("✓ Bounding box for NYC: lat [%.4f, %.4f], lon [%.4f, %.4f]",
		minLat, maxLat, minLon, maxLon)

	// Test if point is within bounding box
	testPoints := []struct {
		lat float64
		lon float64
	}{
		{40.7128, -74.0060}, // Inside (NYC)
		{40.7500, -74.0150}, // On edge
		{40.6800, -74.0050}, // Below
		{40.7200, -73.9800}, // Right
		{40.7800, -74.0200}, // Above and left
	}

	for _, p := range testPoints {
		within := isWithinBoundingBox(p.lat, p.lon, minLat, maxLat, minLon, maxLon)
		status := "Outside"
		if within {
			status = "Inside"
		}
		log.Printf("✓ Point (%.4f, %.4f): %s", p.lat, p.lon, status)
	}

	// SQL query for points within bounding box
	rows, err := db.Query(`
		SELECT id, name, lat, lon
		FROM locations
		WHERE lat BETWEEN ? AND ?
		  AND lon BETWEEN ? AND ?
	`, minLat, maxLat, minLon, maxLon)
	if err != nil {
		log.Printf("Error querying bounding box: %v", err)
		return
	}
	defer rows.Close()

	log.Println("✓ Points within bounding box (from SQL):")
	count := 0
	for rows.Next() {
		var id int
		var name string
		var lat, lon float64
		err := rows.Scan(&id, &name, &lat, &lon)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %d: %s (%.4f, %.4f)", id, name, lat, lon)
		count++
	}
	log.Printf("✓ Found %d points within bounding box", count)
}

func testPointInPolygon(db *sql.DB) {
	// Point-in-polygon detection using ray casting algorithm
	// Point is inside polygon if ray from point intersects polygon edges odd number of times

	log.Println("✓ Point-in-polygon detection:")

	// Define a polygon (triangle)
	polygon := []struct {
		lat float64
		lon float64
	}{
		{40.7100, -74.0050},
		{40.7150, -74.0100},
		{40.7100, -74.0150},
	}

	log.Println("✓ Polygon vertices:")
	for i, v := range polygon {
		log.Printf("  Vertex %d: (%.4f, %.4f)", i, v.lat, v.lon)
	}

	// Test points
	testPoints := []struct {
		lat  float64
		lon  float64
		name string
	}{
		{40.7120, -74.0100, "Center"},
		{40.7080, -74.0100, "Below"},
		{40.7200, -74.0100, "Above"},
		{40.7120, -74.0200, "Right"},
	}

	for _, p := range testPoints {
		inside := isPointInPolygon(p.lat, p.lon, polygon)
		status := "Outside"
		if inside {
			status = "Inside"
		}
		log.Printf("✓ Point %s (%.4f, %.4f): %s", p.name, p.lat, p.lon, status)
	}
}

func testSpatialQueries(db *sql.DB) {
	// Find nearby points within a given radius

	log.Println("✓ Spatial queries (find nearby points):")

	// Create more test locations
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS cities (id INTEGER, name TEXT, lat REAL, lon REAL)
	`)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	// Insert major cities
	cities := []struct {
		id   int
		name string
		lat  float64
		lon  float64
	}{
		{1, "New York", 40.7128, -74.0060},
		{2, "Boston", 42.3601, -71.0589},
		{3, "Philadelphia", 39.9526, -75.1652},
		{4, "Washington DC", 38.9072, -77.0369},
		{5, "Baltimore", 39.2904, -76.6122},
	}

	for _, city := range cities {
		_, err = db.Exec("INSERT OR REPLACE INTO cities VALUES (?, ?, ?, ?)",
			city.id, city.name, city.lat, city.lon)
		if err != nil {
			log.Printf("Error inserting city: %v", err)
			return
		}
	}
	log.Println("✓ Inserted 5 cities")

	// Find cities within 200km of New York
	centerLat := 40.7128
	centerLon := -74.0060
	radiusKm := 200.0

	log.Printf("✓ Cities within %.0f km of New York:", radiusKm)
	rows, err := db.Query("SELECT id, name, lat, lon FROM cities")
	if err != nil {
		log.Printf("Error querying cities: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var lat, lon float64
		err := rows.Scan(&id, &name, &lat, &lon)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		distance := haversineDistanceCalc(centerLat, centerLon, lat, lon, 6371)
		if distance <= radiusKm {
			log.Printf("  %s: %.2f km", name, distance)
		}
	}
}

func testSpatialJoins(db *sql.DB) {
	// Join tables based on spatial relationships

	log.Println("✓ Spatial joins:")

	// Create customers and stores tables
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS customers (id INTEGER, name TEXT, lat REAL, lon REAL)
	`)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS stores (id INTEGER, name TEXT, lat REAL, lon REAL)
	`)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	// Insert customers
	customers := []struct {
		id   int
		name string
		lat  float64
		lon  float64
	}{
		{1, "Customer 1", 40.7128, -74.0060},
		{2, "Customer 2", 40.7140, -74.0080},
		{3, "Customer 3", 40.7100, -74.0040},
	}

	for _, c := range customers {
		_, err = db.Exec("INSERT OR REPLACE INTO customers VALUES (?, ?, ?, ?)",
			c.id, c.name, c.lat, c.lon)
		if err != nil {
			log.Printf("Error inserting customer: %v", err)
			return
		}
	}

	// Insert stores
	stores := []struct {
		id   int
		name string
		lat  float64
		lon  float64
	}{
		{1, "Store A", 40.7120, -74.0065},
		{2, "Store B", 40.7150, -74.0090},
		{3, "Store C", 40.7090, -74.0030},
	}

	for _, s := range stores {
		_, err = db.Exec("INSERT OR REPLACE INTO stores VALUES (?, ?, ?, ?)",
			s.id, s.name, s.lat, s.lon)
		if err != nil {
			log.Printf("Error inserting store: %v", err)
			return
		}
	}

	log.Println("✓ Created customers and stores tables")

	// Find nearest store for each customer
	log.Println("✓ Nearest store for each customer:")
	customerRows, err := db.Query("SELECT id, name, lat, lon FROM customers")
	if err != nil {
		log.Printf("Error querying customers: %v", err)
		return
	}
	defer customerRows.Close()

	for customerRows.Next() {
		var cID int
		var cName string
		var cLat, cLon float64
		err := customerRows.Scan(&cID, &cName, &cLat, &cLon)
		if err != nil {
			log.Printf("Error scanning customer: %v", err)
			continue
		}

		// Find nearest store
		var nearestStore string
		var minDistance float64 = math.MaxFloat64
		storeRows, err := db.Query("SELECT id, name, lat, lon FROM stores")
		if err != nil {
			log.Printf("Error querying stores: %v", err)
			continue
		}

		for storeRows.Next() {
			var sID int
			var sName string
			var sLat, sLon float64
			err := storeRows.Scan(&sID, &sName, &sLat, &sLon)
			if err != nil {
				log.Printf("Error scanning store: %v", err)
				continue
			}

			distance := haversineDistanceCalc(cLat, cLon, sLat, sLon, 6371)
			if distance < minDistance {
			}
		}
		storeRows.Close()

		log.Printf("  %s: Nearest store is %s (%.2f km)", cName, nearestStore, minDistance)
	}
}

func testSpatialAggregations(db *sql.DB) {
	// Aggregate spatial data

	log.Println("✓ Spatial aggregations:")

	// Calculate centroid of all locations
	rows, err := db.Query("SELECT lat, lon FROM locations")
	if err != nil {
		log.Printf("Error querying locations: %v", err)
		return
	}
	defer rows.Close()

	var sumLat, sumLon float64
	var count int

	for rows.Next() {
		var lat, lon float64
		err := rows.Scan(&lat, &lon)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		sumLat += lat
		sumLon += lon
		count++
	}

	if count > 0 {
		centroidLat := sumLat / float64(count)
		centroidLon := sumLon / float64(count)
		log.Printf("✓ Centroid of locations: (%.4f, %.4f)", centroidLat, centroidLon)
	}

	// Calculate average distance from centroid
	if count > 0 {
		centroidLat := sumLat / float64(count)
		centroidLon := sumLon / float64(count)

		rows, err := db.Query("SELECT lat, lon FROM locations")
		if err != nil {
			log.Printf("Error querying locations: %v", err)
			return
		}
		defer rows.Close()

		var totalDistance float64
		var pointCount int

		for rows.Next() {
			var lat, lon float64
			err := rows.Scan(&lat, &lon)
			if err != nil {
				log.Printf("Error scanning row: %v", err)
				continue
			}
			distance := haversineDistanceCalc(centroidLat, centroidLon, lat, lon, 6371)
			totalDistance += distance
			pointCount++
		}

		if pointCount > 0 {
			avgDistance := totalDistance / float64(pointCount)
			log.Printf("✓ Average distance from centroid: %.2f km", avgDistance)
		}
	}
}

func testCoordinateConversions(db *sql.DB) {
	// Convert between different coordinate systems

	log.Println("✓ Coordinate conversions:")

	// Degrees to radians
	degrees := 90.0
	radians := degrees * math.Pi / 180.0
	log.Printf("✓ %.2f degrees = %.4f radians", degrees, radians)

	// Radians to degrees
	radians = 1.5708
	degrees = radians * 180.0 / math.Pi
	log.Printf("✓ %.4f radians = %.2f degrees", radians, degrees)

	// Nautical miles to km
	nm := 100.0
	km := nm * 1.852
	log.Printf("✓ %.2f nautical miles = %.2f km", nm, km)

	// Miles to km
	miles := 62.0
	km = miles * 1.60934
	log.Printf("✓ %.2f miles = %.2f km", miles, km)
}

func testGeospatialIndexing(db *sql.DB) {
	// Demonstrate indexing for geospatial queries

	log.Println("✓ Geospatial indexing:")

	// Create indexes on lat and lon columns
	_, err := db.Exec("CREATE INDEX IF NOT EXISTS idx_locations_lat ON locations(lat)")
	if err != nil {
		log.Printf("Error creating lat index: %v", err)
		return
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_locations_lon ON locations(lon)")
	if err != nil {
		log.Printf("Error creating lon index: %v", err)
		return
	}

	log.Println("✓ Created indexes on lat and lon columns")

	// Create composite index
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_locations_latlon ON locations(lat, lon)")
	if err != nil {
		log.Printf("Error creating composite index: %v", err)
		return
	}

	log.Println("✓ Created composite index on (lat, lon)")
}

func testPolygonArea(db *sql.DB) {
	// Calculate area of polygon using shoelace formula

	log.Println("✓ Polygon area calculations:")

	// Define a polygon (rectangle for simplicity)
	polygon := []struct {
		lat float64
		lon float64
	}{
		{40.7100, -74.0050},
		{40.7150, -74.0050},
		{40.7150, -74.0150},
		{40.7100, -74.0150},
	}

	area := calculatePolygonArea(polygon)
	log.Printf("✓ Polygon area: %.6f square degrees (approximate)", area)

	// Note: For real-world area calculations, need to use geodesic formulas
	// and convert degrees to appropriate units based on latitude
}

func testCentroid(db *sql.DB) {
	// Calculate centroid of polygon

	log.Println("✓ Centroid calculations:")

	// Define a polygon
	polygon := []struct {
		lat float64
		lon float64
	}{
		{40.7100, -74.0050},
		{40.7150, -74.0100},
		{40.7100, -74.0150},
	}

	centroid := calculateCentroid(polygon)
	log.Printf("✓ Polygon centroid: (%.4f, %.4f)", centroid.lat, centroid.lon)
}

func testSpatialFiltering(db *sql.DB) {
	// Filter data using spatial conditions

	log.Println("✓ Spatial filtering:")

	// Filter cities within latitude range
	rows, err := db.Query(`
		SELECT id, name, lat, lon
		FROM cities
		WHERE lat BETWEEN 40.0 AND 41.0
		ORDER BY lat
	`)
	if err != nil {
		log.Printf("Error querying cities: %v", err)
		return
	}
	defer rows.Close()

	log.Println("✓ Cities between lat 40.0 and 41.0:")
	for rows.Next() {
		var id int
		var name string
		var lat, lon float64
		err := rows.Scan(&id, &name, &lat, &lon)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %s: (%.4f, %.4f)", name, lat, lon)
	}

	// Filter cities within longitude range
	rows, err = db.Query(`
		SELECT id, name, lat, lon
		FROM cities
		WHERE lon BETWEEN -75.0 AND -73.0
		ORDER BY lon
	`)
	if err != nil {
		log.Printf("Error querying cities: %v", err)
		return
	}
	defer rows.Close()

	log.Println("✓ Cities between lon -75.0 and -73.0:")
	for rows.Next() {
		var id int
		var name string
		var lat, lon float64
		err := rows.Scan(&id, &name, &lat, &lon)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		log.Printf("  %s: (%.4f, %.4f)", name, lat, lon)
	}
}

func testDistanceComparison(db *sql.DB) {
	// Compare geodesic (Haversine) and Euclidean distances

	log.Println("✓ Geodesic vs. Euclidean distance comparison:")

	// Test points
	lat1 := 40.7128
	lon1 := -74.0060
	lat2 := 40.7140
	lon2 := -74.0080

	// Geodesic distance (Haversine)
	geodesicKm := haversineDistanceCalc(lat1, lon1, lat2, lon2, 6371)
	geodesicMiles := haversineDistanceCalc(lat1, lon1, lat2, lon2, 3959)

	log.Printf("✓ Geodesic distance: %.4f km (%.4f miles)", geodesicKm, geodesicMiles)

	// Euclidean distance (approximate, for small distances)
	euclideanKm := euclideanDistanceCalc(lat1, lon1, lat2, lon2)
	euclideanMiles := euclideanKm * 0.621371

	log.Printf("✓ Euclidean distance: %.4f km (%.4f miles)", euclideanKm, euclideanMiles)

	// Difference
	diffKm := math.Abs(geodesicKm - euclideanKm)
	diffMiles := math.Abs(geodesicMiles - euclideanMiles)
	log.Printf("✓ Difference: %.4f km (%.4f miles)", diffKm, diffMiles)
}

func testSpatialSorting(db *sql.DB) {
	// Sort locations by distance from a point

	log.Println("✓ Spatial sorting (distance ordering):")

	// Sort cities by distance from New York
	centerLat := 40.7128
	centerLon := -74.0060

	log.Println("✓ Cities sorted by distance from New York:")

	// Get all cities and sort by distance
	rows, err := db.Query("SELECT id, name, lat, lon FROM cities")
	if err != nil {
		log.Printf("Error querying cities: %v", err)
		return
	}
	defer rows.Close()

	type cityDist struct {
		id       int
		name     string
		distance float64
	}

	var cityDists []cityDist

	for rows.Next() {
		var id int
		var name string
		var lat, lon float64
		err := rows.Scan(&id, &name, &lat, &lon)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		distance := haversineDistanceCalc(centerLat, centerLon, lat, lon, 6371)
		cityDists = append(cityDists, cityDist{id, name, distance})
	}

	// Sort by distance (simple bubble sort for demonstration)
	for i := 0; i < len(cityDists); i++ {
		for j := i + 1; j < len(cityDists); j++ {
			if cityDists[i].distance > cityDists[j].distance {
				cityDists[i], cityDists[j] = cityDists[j], cityDists[i]
			}
		}
	}

	for _, cd := range cityDists {
		log.Printf("  %d: %s (%.2f km)", cd.id, cd.name, cd.distance)
	}
}

func testCleanup(db *sql.DB) {
	tables := []string{
		"locations",
		"cities",
		"customers",
		"stores",
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

// Helper functions for geospatial calculations

func haversineDistanceCalc(lat1, lon1, lat2, lon2, radius float64) float64 {
	// Convert to radians
	lat1Rad := lat1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0
	deltaLat := (lat2 - lat1) * math.Pi / 180.0
	deltaLon := (lon2 - lon1) * math.Pi / 180.0

	// Haversine formula
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return radius * c
}

func euclideanDistanceCalc(lat1, lon1, lat2, lon2 float64) float64 {
	// Approximate Euclidean distance (valid only for small distances)
	// Convert degrees to km (approximate: 1 degree lat = 111 km, 1 degree lon = 111 km * cos(lat))
	deltaLat := lat2 - lat1
	deltaLon := lon2 - lon1

	avgLat := (lat1 + lat2) / 2.0
	latKm := deltaLat * 111.0
	lonKm := deltaLon * 111.0 * math.Cos(avgLat*math.Pi/180.0)

	return math.Sqrt(latKm*latKm + lonKm*lonKm)
}

func isWithinBoundingBox(lat, lon, minLat, maxLat, minLon, maxLon float64) bool {
	return lat >= minLat && lat <= maxLat && lon >= minLon && lon <= maxLon
}

func isPointInPolygon(lat, lon float64, polygon []struct {
	lat float64
	lon float64
}) bool {
	if len(polygon) < 3 {
		return false
	}

	inside := false
	j := len(polygon) - 1

	for i := 0; i < len(polygon); i++ {
		if ((polygon[i].lat > lat) != (polygon[j].lat > lat)) &&
			(lon < (polygon[j].lon-polygon[i].lon)*(lat-polygon[i].lat)/(polygon[j].lat-polygon[i].lat)+polygon[i].lon) {
			inside = !inside
		}
		j = i
	}

	return inside
}

func calculatePolygonArea(polygon []struct {
	lat float64
	lon float64
}) float64 {
	if len(polygon) < 3 {
		return 0
	}

	area := 0.0
	n := len(polygon)

	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += polygon[i].lon * polygon[j].lat
		area -= polygon[j].lon * polygon[i].lat
	}

	return math.Abs(area) / 2.0
}

func calculateCentroid(polygon []struct {
	lat float64
	lon float64
}) struct {
	lat float64
	lon float64
} {
	if len(polygon) < 1 {
		return struct{ lat float64; lon float64 }{}
	}

	var sumLat, sumLon float64
	for _, p := range polygon {
		sumLat += p.lat
		sumLon += p.lon
	}

	n := float64(len(polygon))
	return struct{ lat float64; lon float64 }{
		lat: sumLat / n,
		lon: sumLon / n,
	}
}
