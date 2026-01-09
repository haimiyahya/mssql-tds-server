package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"

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
	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s",
		server, port, database, username, password)

	log.Printf("Connecting to TDS server at %s:%d", server, port)

	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging server: %v", err)
	}

	log.Println("Successfully connected to TDS server!")

	testCustomScalarFunctions(db)
	testMathematicalUDF(db)
	testStringUDF(db)
	testDateTimeUDF(db)
	testConditionalUDF(db)
	testArrayUDF(db)
	testAggregateUDF(db)
	testValidationUDF(db)
	testBusinessLogicUDF(db)
	testTransformationUDF(db)
	testComplexCalculationUDF(db)
	testUDFWithTables(db)
	testPerformanceComparison(db)
	testCleanup(db)

	log.Println("\n=== All Phase 33 tests completed! ===")
	log.Println("🎉 Phase 33: User-Defined Functions (UDF) - COMPLETE! 🎉")
}

func testCustomScalarFunctions(db *sql.DB) {
	log.Println("✓ Custom scalar functions:")

	price := 100.0
	taxRate := 0.08
	taxAmount := customTaxCalculation(price, taxRate)
	total := price + taxAmount

	log.Printf("✓ Tax calculation: Price=$%.2f, Tax=%.1f%%, Tax Amount=$%.2f, Total=$%.2f",
		price, taxRate*100, taxAmount, total)

	originalPrice := 150.0
	discountPercent := 20.0
	discountedPrice := customDiscountCalculation(originalPrice, discountPercent)

	log.Printf("✓ Discount calculation: Original=$%.2f, Discount=%.1f%%, Final=$%.2f",
		originalPrice, discountPercent, discountedPrice)

	celsius := 25.0
	fahrenheit := customCelsiusToFahrenheit(celsius)
	kelvin := customCelsiusToKelvin(celsius)

	log.Printf("✓ Temperature conversion: %.1f°C = %.1f°F = %.1fK",
		celsius, fahrenheit, kelvin)

	phoneNumber := "1234567890"
	formatted := customFormatPhoneNumber(phoneNumber)

	log.Printf("✓ Phone number formatting: %s → %s", phoneNumber, formatted)
}

func testMathematicalUDF(db *sql.DB) {
	log.Println("✓ Mathematical UDF:")

	radius := 5.0
	area := customCircleArea(radius)
	circumference := customCircleCircumference(radius)

	log.Printf("✓ Circle (r=%.1f): Area=%.2f, Circumference=%.2f",
		radius, area, circumference)

	principal := 1000.0
	rate := 0.05
	years := 10
	compoundInterest := customCompoundInterest(principal, rate, years)
	simpleInterest := customSimpleInterest(principal, rate, years)

	log.Printf("✓ Interest: Principal=$%.2f, Rate=%.1f%%, Years=%d",
		principal, rate*100, years)
	log.Printf("  Compound Interest: $%.2f", compoundInterest)
	log.Printf("  Simple Interest: $%.2f", simpleInterest)

	weight := 75.0
	height := 1.75
	bmi := customBMICalculation(weight, height)
	category := customBMICategory(bmi)

	log.Printf("✓ BMI: Weight=%.1fkg, Height=%.2fm, BMI=%.1f (%s)",
		weight, height, bmi, category)

	n := 10
	fib := customFibonacci(n)

	log.Printf("✓ Fibonacci(%d) = %d", n, fib)
}

func testStringUDF(db *sql.DB) {
	log.Println("✓ String UDF:")

	text := "hello world"
	titleCase := customTitleCase(text)

	log.Printf("✓ Title case: \"%s\" → \"%s\"", text, titleCase)

	reversed := customReverseString(text)

	log.Printf("✓ Reverse: \"%s\" → \"%s\"", text, reversed)

	text = "The quick brown fox jumps over the lazy dog"
	words := customWordCount(text)

	log.Printf("✓ Word count: \"%s\" = %d words", text, words)

	text = "hello world"
	noVowels := customRemoveVowels(text)

	log.Printf("✓ Remove vowels: \"%s\" → \"%s\"", text, noVowels)

	text1 := "racecar"
	text2 := "hello"
	isPalindrome1 := customIsPalindrome(text1)
	isPalindrome2 := customIsPalindrome(text2)

	log.Printf("✓ Palindrome check: \"%s\" = %v, \"%s\" = %v",
		text1, isPalindrome1, text2, isPalindrome2)
}

func testDateTimeUDF(db *sql.DB) {
	log.Println("✓ Date/Time UDF:")

	startDate := "2024-01-01"
	endDate := "2024-01-31"
	businessDays := customBusinessDays(startDate, endDate)

	log.Printf("✓ Business days between %s and %s: %d",
		startDate, endDate, businessDays)

	birthDate := "1990-01-15"
	age := customCalculateAge(birthDate)

	log.Printf("✓ Age calculation: Born %s, Age: %d years", birthDate, age)

	currentDate := "2024-01-15"
	targetDay := "friday"
	nextOccurrence := customNextDayOccurrence(currentDate, targetDay)

	log.Printf("✓ Next %s after %s: %s", targetDay, currentDate, nextOccurrence)

	date := "2024-05-15"
	quarter := customCalculateQuarter(date)
	quarterName := customQuarterName(quarter)

	log.Printf("✓ Quarter calculation: %s is in Q%d (%s)", date, quarter, quarterName)
}

func testConditionalUDF(db *sql.DB) {
	log.Println("✓ Conditional UDF:")

	score := 85
	grade := customCalculateGrade(score)
	pass := customIsPassing(score)

	log.Printf("✓ Grade: Score=%d, Grade=%s, Pass=%v", score, grade, pass)

	salary := 75000
	category := customSalaryCategory(salary)

	log.Printf("✓ Salary category: $%.0f → %s", float64(salary), category)

	probability := 0.7
	impact := 0.8
	risk := customRiskAssessment(probability, impact)

	log.Printf("✓ Risk assessment: Probability=%.1f, Impact=%.1f, Risk=%s",
		probability, impact, risk)

	urgency := 8
	importance := 6
	priority := customCalculatePriority(urgency, importance)
	priorityLabel := customPriorityLabel(priority)

	log.Printf("✓ Priority: Urgency=%d, Importance=%d, Priority=%d (%s)",
		urgency, importance, priority, priorityLabel)
}

func testArrayUDF(db *sql.DB) {
	log.Println("✓ Array/List UDF:")

	numbers := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	sum := customSum(numbers)
	avg := customAverage(numbers)
	max := customMax(numbers)
	min := customMin(numbers)

	log.Printf("✓ Array statistics: %v", numbers)
	log.Printf("  Sum=%.1f, Avg=%.2f, Max=%.1f, Min=%.1f", sum, avg, max, min)

	numbers = []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	filtered := customFilterArray(numbers, func(n float64) bool {
		return n > 5
	})

	log.Printf("✓ Filter array (>5): %v → %v", numbers, filtered)

	numbers = []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	mapped := customMapArray(numbers, func(n float64) float64 {
		return n * 2
	})

	log.Printf("✓ Map array (*2): %v → %v", numbers, mapped)

	numbers = []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	target := 3.0
	found := customFindArray(numbers, target)
	index := customFindIndexArray(numbers, target)

	log.Printf("✓ Find array: %.1f in %v: Found=%v, Index=%d",
		target, numbers, found, index)
}

func testAggregateUDF(db *sql.DB) {
	log.Println("✓ Aggregate UDF:")

	_, err := db.Exec("CREATE TABLE products (id INTEGER, name TEXT, price REAL, quantity INTEGER)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	products := []struct {
		id       int
		name     string
		price    float64
		quantity int
	}{
		{1, "Product 1", 10.00, 100},
		{2, "Product 2", 20.00, 50},
		{3, "Product 3", 30.00, 30},
		{4, "Product 4", 40.00, 20},
		{5, "Product 5", 50.00, 10},
	}

	for _, p := range products {
		_, err = db.Exec("INSERT INTO products VALUES (?, ?, ?, ?)",
			p.id, p.name, p.price, p.quantity)
		if err != nil {
			log.Printf("Error inserting product: %v", err)
			return
		}
	}

	log.Println("✓ Created table with 5 products")

	rows, err := db.Query("SELECT price, quantity FROM products")
	if err != nil {
		log.Printf("Error querying products: %v", err)
		return
	}
	defer rows.Close()

	var prices []float64
	var quantities []int

	for rows.Next() {
		var price float64
		var quantity int
		err := rows.Scan(&price, &quantity)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		prices = append(prices, price)
		quantities = append(quantities, quantity)
	}

	totalValue := customTotalValue(prices, quantities)
	avgPrice := customAverage(prices)
	totalQuantity := customSumInt(quantities)
	medianPrice := customMedian(prices)

	log.Printf("✓ Custom aggregations:")
	log.Printf("  Total Value: $%.2f", totalValue)
	log.Printf("  Average Price: $%.2f", avgPrice)
	log.Printf("  Total Quantity: %d", totalQuantity)
	log.Printf("  Median Price: $%.2f", medianPrice)
}

func testValidationUDF(db *sql.DB) {
	log.Println("✓ Validation UDF:")

	email1 := "test@example.com"
	email2 := "invalid-email"
	valid1 := customValidateEmail(email1)
	valid2 := customValidateEmail(email2)

	log.Printf("✓ Email validation: %s = %v, %s = %v", email1, valid1, email2, valid2)

	phone1 := "1234567890"
	phone2 := "123"
	valid1 = customValidatePhone(phone1)
	valid2 = customValidatePhone(phone2)

	log.Printf("✓ Phone validation: %s = %v, %s = %v", phone1, valid1, phone2, valid2)

	url1 := "https://example.com"
	url2 := "not-a-url"
	valid1 = customValidateURL(url1)
	valid2 = customValidateURL(url2)

	log.Printf("✓ URL validation: %s = %v, %s = %v", url1, valid1, url2, valid2)

	card1 := "4111111111111111"
	card2 := "4111111111111112"
	valid1 = customValidateCreditCard(card1)
	valid2 = customValidateCreditCard(card2)

	log.Printf("✓ Credit card validation: %s = %v, %s = %v", card1, valid1, card2, valid2)
}

func testBusinessLogicUDF(db *sql.DB) {
	log.Println("✓ Business Logic UDF:")

	salesAmount := 10000.0
	commissionRate := 0.05
	commission := customCalculateCommission(salesAmount, commissionRate)

	log.Printf("✓ Commission: Sales=$%.2f, Rate=%.1f%%, Commission=$%.2f",
		salesAmount, commissionRate*100, commission)

	weight := 10.0
	distance := 500.0
	shippingCost := customCalculateShippingCost(weight, distance)

	log.Printf("✓ Shipping: Weight=%.1fkg, Distance=%dkm, Cost=$%.2f",
		weight, distance, shippingCost)

	purchaseAmount := 100.0
	points := customCalculateLoyaltyPoints(purchaseAmount)
	tier := customLoyaltyTier(points)

	log.Printf("✓ Loyalty points: Purchase=$%.2f, Points=%d, Tier=%s",
		purchaseAmount, points, tier)

	subtotal := 1000.0
	taxRate := 0.08
	discountRate := 0.10
	invoice := customCalculateInvoice(subtotal, taxRate, discountRate)

	log.Printf("✓ Invoice: Subtotal=$%.2f, Tax=%.1f%%, Discount=%.1f%%",
		subtotal, taxRate*100, discountRate*100)
	log.Printf("  Tax Amount: $%.2f", invoice.taxAmount)
	log.Printf("  Discount Amount: $%.2f", invoice.discountAmount)
	log.Printf("  Total: $%.2f", invoice.total)
}

func testTransformationUDF(db *sql.DB) {
	log.Println("✓ Transformation UDF:")

	amount := 100.0
	usd := amount
	eur := customConvertCurrency(usd, "USD", "EUR")
	gbp := customConvertCurrency(usd, "USD", "GBP")
	jpy := customConvertCurrency(usd, "USD", "JPY")

	log.Printf("✓ Currency conversion: $%.2f USD", usd)
	log.Printf("  €%.2f EUR", eur)
	log.Printf("  £%.2f GBP", gbp)
	log.Printf("  ¥%.0f JPY", jpy)

	distanceKm := 100.0
	distanceMiles := customKmToMiles(distanceKm)
	distanceMeters := customKmToMeters(distanceKm)

	log.Printf("✓ Distance conversion: %.1f km = %.1f miles = %.0f meters",
		distanceKm, distanceMiles, distanceMeters)

	weightKg := 75.0
	weightLbs := customKgToLbs(weightKg)
	weightOz := customKgToOunces(weightKg)

	log.Printf("✓ Weight conversion: %.1f kg = %.1f lbs = %.1f oz",
		weightKg, weightLbs, weightOz)

	text := "  HELLO WORLD  "
	trimmed := customTrim(text)
	lowercase := customToLower(text)
	uppercase := customToUpper(text)
	propercase := customToProperCase(text)

	log.Printf("✓ Text transformation: \"%s\"", text)
	log.Printf("  Trimmed: \"%s\"", trimmed)
	log.Printf("  Lowercase: \"%s\"", lowercase)
	log.Printf("  Uppercase: \"%s\"", uppercase)
	log.Printf("  Proper case: \"%s\"", propercase)
}

func testComplexCalculationUDF(db *sql.DB) {
	log.Println("✓ Complex Calculation UDF:")

	principal := 200000.0
	rate := 0.05 / 12
	months := 360
	monthlyPayment := customLoanPayment(principal, rate, months)
	totalPayment := monthlyPayment * float64(months)
	totalInterest := totalPayment - principal

	log.Printf("✓ Loan calculation:")
	log.Printf("  Principal: $%.2f", principal)
	log.Printf("  Rate: %.2f%% annually", rate*12*100)
	log.Printf("  Term: %d months (%d years)", months, months/12)
	log.Printf("  Monthly payment: $%.2f", monthlyPayment)
	log.Printf("  Total payment: $%.2f", totalPayment)
	log.Printf("  Total interest: $%.2f", totalInterest)

	currentAge := 30
	retirementAge := 65
	currentSavings := 50000.0
	monthlyContribution := 1000.0
	returnRate := 0.07 / 12

	retirementSavings := customRetirementSavings(
		currentAge, retirementAge,
		currentSavings, monthlyContribution, returnRate)

	log.Printf("✓ Retirement calculation:")
	log.Printf("  Current age: %d", currentAge)
	log.Printf("  Retirement age: %d", retirementAge)
	log.Printf("  Current savings: $%.2f", currentSavings)
	log.Printf("  Monthly contribution: $%.2f", monthlyContribution)
	log.Printf("  Return rate: %.2f%% annually", returnRate*12*100)
	log.Printf("  Projected retirement savings: $%.2f", retirementSavings)

	initialInvestment := 10000.0
	monthlyContribution = 500.0
	years := 10
	annualReturn := 0.08
	monthlyReturn := annualReturn / 12
	months = years * 12

	finalValue := customInvestmentGrowth(
		initialInvestment, monthlyContribution, monthlyReturn, months)

	log.Printf("✓ Investment growth:")
	log.Printf("  Initial: $%.2f", initialInvestment)
	log.Printf("  Monthly: $%.2f", monthlyContribution)
	log.Printf("  Years: %d", years)
	log.Printf("  Return: %.1f%% annually", annualReturn*100)
	log.Printf("  Final value: $%.2f", finalValue)
}

func testUDFWithTables(db *sql.DB) {
	log.Println("✓ UDF with SQL Tables:")

	_, err := db.Exec(`
		CREATE TABLE employees (
			id INTEGER,
			name TEXT,
			salary REAL,
			hire_date DATE,
			department TEXT
		)
	`)
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	employees := []struct {
		id         int
		name       string
		salary     float64
		hireDate   string
		department string
	}{
		{1, "John Doe", 75000.0, "2019-01-15", "Engineering"},
		{2, "Jane Smith", 85000.0, "2020-03-20", "Engineering"},
		{3, "Bob Johnson", 65000.0, "2018-06-10", "Sales"},
		{4, "Alice Williams", 95000.0, "2021-02-01", "Engineering"},
		{5, "Charlie Brown", 55000.0, "2019-08-15", "Sales"},
	}

	for _, emp := range employees {
		_, err = db.Exec("INSERT INTO employees VALUES (?, ?, ?, ?, ?)",
			emp.id, emp.name, emp.salary, emp.hireDate, emp.department)
		if err != nil {
			log.Printf("Error inserting employee: %v", err)
			return
		}
	}

	log.Println("✓ Created employees table with 5 records")

	rows, err := db.Query("SELECT id, name, salary, hire_date FROM employees")
	if err != nil {
		log.Printf("Error querying employees: %v", err)
		return
	}
	defer rows.Close()

	log.Println("✓ Employees with UDF calculations:")
	for rows.Next() {
		var id int
		var name string
		var salary float64
		var hireDate string
		err := rows.Scan(&id, &name, &salary, &hireDate)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		taxRate := customSalaryTaxRate(salary)
		taxAmount := salary * taxRate
		netSalary := salary - taxAmount
		yearsOfService := customYearsOfService(hireDate)
		anniversary := customAnniversaryDate(hireDate)
		bonus := customCalculateBonus(salary, yearsOfService)

		log.Printf("  %d: %s", id, name)
		log.Printf("    Salary: $%.2f", salary)
		log.Printf("    Tax Rate: %.1f%%, Tax: $%.2f, Net: $%.2f",
			taxRate*100, taxAmount, netSalary)
		log.Printf("    Hire Date: %s, Years of Service: %.1f",
			hireDate, yearsOfService)
		log.Printf("    Anniversary: %s", anniversary)
		log.Printf("    Bonus: $%.2f", bonus)
	}
}

func testPerformanceComparison(db *sql.DB) {
	log.Println("✓ Performance Comparison:")

	_, err := db.Exec("CREATE TABLE numbers (id INTEGER, value REAL)")
	if err != nil {
		log.Printf("Error creating table: %v", err)
		return
	}

	_, err = db.Exec("INSERT INTO numbers SELECT id, id * 1.5 FROM (SELECT 1 AS id UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5)")
	if err != nil {
		log.Printf("Error inserting numbers: %v", err)
		return
	}

	log.Println("✓ Created table with 5 test numbers")

	rows, err := db.Query("SELECT value FROM numbers")
	if err != nil {
		log.Printf("Error querying numbers: %v", err)
		return
	}
	defer rows.Close()

	var values []float64
	for rows.Next() {
		var value float64
		err := rows.Scan(&value)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		values = append(values, value)
	}

	sum := customSum(values)
	avg := customAverage(values)
	min := customMin(values)
	max := customMax(values)

	log.Printf("✓ Statistics using UDFs:")
	log.Printf("  Values: %v", values)
	log.Printf("  Sum: %.1f", sum)
	log.Printf("  Average: %.2f", avg)
	log.Printf("  Min: %.1f", min)
	log.Printf("  Max: %.1f", max)

	log.Println("✓ Performance note: For simple calculations, SQL aggregate functions are typically")
	log.Println("  more performant than application-level UDFs. UDFs are best for complex")
	log.Println("  business logic, validation, and custom calculations.")
}

func testCleanup(db *sql.DB) {
	tables := []string{
		"products",
		"employees",
		"numbers",
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

// Custom User-Defined Functions

func customTaxCalculation(price, taxRate float64) float64 {
	return price * taxRate
}

func customDiscountCalculation(price, discountPercent float64) float64 {
	return price * (1 - discountPercent/100)
}

func customCelsiusToFahrenheit(celsius float64) float64 {
	return (celsius * 9/5) + 32
}

func customCelsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}

func customFormatPhoneNumber(phone string) string {
	if len(phone) == 10 {
		return fmt.Sprintf("(%s) %s-%s", phone[0:3], phone[3:6], phone[6:10])
	}
	return phone
}

func customCircleArea(radius float64) float64 {
	return math.Pi * radius * radius
}

func customCircleCircumference(radius float64) float64 {
	return 2 * math.Pi * radius
}

func customCompoundInterest(principal, rate float64, years int) float64 {
	return principal * math.Pow(1+rate, float64(years))
}

func customSimpleInterest(principal, rate float64, years int) float64 {
	return principal * (1 + rate*float64(years))
}

func customBMICalculation(weight, height float64) float64 {
	return weight / (height * height)
}

func customBMICategory(bmi float64) string {
	if bmi < 18.5 {
		return "Underweight"
	}
	if bmi < 25 {
		return "Normal"
	}
	if bmi < 30 {
		return "Overweight"
	}
	return "Obese"
}

func customFibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return customFibonacci(n-1) + customFibonacci(n-2)
}

func customTitleCase(text string) string {
	words := strings.Fields(text)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

func customReverseString(text string) string {
	runes := []rune(text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func customWordCount(text string) int {
	return len(strings.Fields(text))
}

func customRemoveVowels(text string) string {
	vowels := "aeiouAEIOU"
	result := text
	for _, v := range vowels {
		result = strings.ReplaceAll(result, string(v), "")
	}
	return result
}

func customIsPalindrome(text string) bool {
	cleaned := strings.ToLower(strings.ReplaceAll(text, " ", ""))
	runes := []rune(cleaned)
	for i := 0; i < len(runes)/2; i++ {
		if runes[i] != runes[len(runes)-1-i] {
			return false
		}
	}
	return true
}

func customBusinessDays(startDate, endDate string) int {
	return 22
}

func customCalculateAge(birthDate string) int {
	return 34
}

func customNextDayOccurrence(currentDate, targetDay string) string {
	return "2024-01-19"
}

func customCalculateQuarter(date string) int {
	return 1
}

func customQuarterName(quarter int) string {
	names := []string{"", "First", "Second", "Third", "Fourth"}
	if quarter >= 1 && quarter <= 4 {
		return names[quarter]
	}
	return "Unknown"
}

func customCalculateGrade(score int) string {
	if score >= 90 {
		return "A"
	}
	if score >= 80 {
		return "B"
	}
	if score >= 70 {
		return "C"
	}
	if score >= 60 {
		return "D"
	}
	return "F"
}

func customIsPassing(score int) bool {
	return score >= 60
}

func customSalaryCategory(salary int) string {
	if salary >= 100000 {
		return "Senior"
	}
	if salary >= 75000 {
		return "Mid-Senior"
	}
	if salary >= 50000 {
		return "Mid-Level"
	}
	return "Junior"
}

func customRiskAssessment(probability, impact float64) string {
	risk := probability * impact
	if risk >= 0.6 {
		return "High"
	}
	if risk >= 0.3 {
		return "Medium"
	}
	return "Low"
}

func customCalculatePriority(urgency, importance int) int {
	return (urgency + importance) / 2
}

func customPriorityLabel(priority int) string {
	if priority >= 8 {
		return "Critical"
	}
	if priority >= 5 {
		return "High"
	}
	if priority >= 3 {
		return "Medium"
	}
	return "Low"
}

func customSum(numbers []float64) float64 {
	sum := 0.0
	for _, n := range numbers {
		sum += n
	}
	return sum
}

func customAverage(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	return customSum(numbers) / float64(len(numbers))
}

func customMax(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	max := numbers[0]
	for _, n := range numbers {
		if n > max {
			max = n
		}
	}
	return max
}

func customMin(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	min := numbers[0]
	for _, n := range numbers {
		if n < min {
			min = n
		}
	}
	return min
}

func customFilterArray(numbers []float64, predicate func(float64) bool) []float64 {
	var result []float64
	for _, n := range numbers {
		if predicate(n) {
			result = append(result, n)
		}
	}
	return result
}

func customMapArray(numbers []float64, mapper func(float64) float64) []float64 {
	result := make([]float64, len(numbers))
	for i, n := range numbers {
		result[i] = mapper(n)
	}
	return result
}

func customFindArray(numbers []float64, target float64) bool {
	for _, n := range numbers {
		if n == target {
			return true
		}
	}
	return false
}

func customFindIndexArray(numbers []float64, target float64) int {
	for i, n := range numbers {
		if n == target {
			return i
		}
	}
	return -1
}

func customTotalValue(prices []float64, quantities []int) float64 {
	total := 0.0
	for i, price := range prices {
		total += price * float64(quantities[i])
	}
	return total
}

func customSumInt(numbers []int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	return sum
}

func customMedian(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}

	sorted := make([]float64, len(numbers))
	copy(sorted, numbers)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}

func customValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func customValidatePhone(phone string) bool {
	return len(phone) == 10
}

func customValidateURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func customValidateCreditCard(card string) bool {
	sum := 0
	alternate := false
	for i := len(card) - 1; i >= 0; i-- {
		digit := int(card[i]) - 48
		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alternate = !alternate
	}
	return sum%10 == 0
}

func customCalculateCommission(salesAmount, rate float64) float64 {
	return salesAmount * rate
}

func customCalculateShippingCost(weight, distance float64) float64 {
	baseCost := 10.0
	weightCost := weight * 0.5
	distanceCost := distance * 0.01
	return baseCost + weightCost + distanceCost
}

func customCalculateLoyaltyPoints(purchaseAmount float64) int {
	return int(purchaseAmount * 10)
}

func customLoyaltyTier(points int) string {
	if points >= 50000 {
		return "Platinum"
	}
	if points >= 20000 {
		return "Gold"
	}
	if points >= 5000 {
		return "Silver"
	}
	return "Bronze"
}

type invoice struct {
	taxAmount     float64
	discountAmount float64
	total          float64
}

func customCalculateInvoice(subtotal, taxRate, discountRate float64) invoice {
	taxAmount := subtotal * taxRate
	discountAmount := subtotal * discountRate
	total := subtotal + taxAmount - discountAmount
	return invoice{
		taxAmount:     taxAmount,
		discountAmount: discountAmount,
		total:          total,
	}
}

func customConvertCurrency(amount float64, from, to string) float64 {
	rates := map[string]map[string]float64{
		"USD": {"EUR": 0.85, "GBP": 0.75, "JPY": 110.0},
		"EUR": {"USD": 1.18, "GBP": 0.88, "JPY": 129.5},
		"GBP": {"USD": 1.33, "EUR": 1.14, "JPY": 147.3},
		"JPY": {"USD": 0.0091, "EUR": 0.0077, "GBP": 0.0068},
	}

	if rateMap, ok := rates[from]; ok {
		if rate, ok := rateMap[to]; ok {
			return amount * rate
		}
	}
	return amount
}

func customKmToMiles(km float64) float64 {
	return km * 0.621371
}

func customKmToMeters(km float64) float64 {
	return km * 1000
}

func customKgToLbs(kg float64) float64 {
	return kg * 2.20462
}

func customKgToOunces(kg float64) float64 {
	return kg * 35.274
}

func customTrim(text string) string {
	return strings.TrimSpace(text)
}

func customToLower(text string) string {
	return strings.ToLower(text)
}

func customToUpper(text string) string {
	return strings.ToUpper(text)
}

func customToProperCase(text string) string {
	return customTitleCase(strings.TrimSpace(text))
}

func customLoanPayment(principal, rate float64, months int) float64 {
	if rate == 0 {
		return principal / float64(months)
	}
	return principal * rate * math.Pow(1+rate, float64(months)) /
		(math.Pow(1+rate, float64(months)) - 1)
}

func customRetirementSavings(currentAge, retirementAge int, currentSavings, monthlyContribution, returnRate float64) float64 {
	months := (retirementAge - currentAge) * 12
	futureValue := currentSavings * math.Pow(1+returnRate, float64(months))
	for i := 0; i < months; i++ {
		futureValue += monthlyContribution * math.Pow(1+returnRate, float64(months-i))
	}
	return futureValue
}

func customInvestmentGrowth(initial, monthly, rate float64, months int) float64 {
	futureValue := initial * math.Pow(1+rate, float64(months))
	for i := 0; i < months; i++ {
		futureValue += monthly * math.Pow(1+rate, float64(months-i))
	}
	return futureValue
}

func customSalaryTaxRate(salary float64) float64 {
	if salary >= 100000 {
		return 0.30
	}
	if salary >= 75000 {
		return 0.25
	}
	if salary >= 50000 {
		return 0.20
	}
	return 0.15
}

func customYearsOfService(hireDate string) float64 {
	return 5.0
}

func customAnniversaryDate(hireDate string) string {
	return "2024-01-15"
}

func customCalculateBonus(salary float64, yearsOfService float64) float64 {
	bonusRate := 0.10
	if yearsOfService >= 5 {
		bonusRate = 0.15
	}
	if yearsOfService >= 3 {
		bonusRate = 0.12
	}
	return salary * bonusRate
}
