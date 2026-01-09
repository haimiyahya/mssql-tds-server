# Phase 11: Advanced SELECT Features - Implementation Plan

## Overview
Extend Phase 10 basic SELECT support to include advanced features commonly used in real-world applications: JOINs, GROUP BY, HAVING, ORDER BY, DISTINCT, aggregate functions, and subqueries.

## Current State
- Phase 10: Basic SELECT implemented
  - SELECT * FROM table
  - SELECT col1, col2 FROM table
  - WHERE clause support
- Parser recognizes basic SELECT syntax
- Executor executes simple SELECT queries
- Result sets include column headers

## Phase 11 Goals
1. **JOIN Support**
   - INNER JOIN
   - LEFT JOIN
   - RIGHT JOIN
   - FULL JOIN
   - ON clause parsing
   - Multiple table joins

2. **GROUP BY and HAVING**
   - GROUP BY clause parsing
   - HAVING clause parsing
   - Group aggregation logic
   - Filter groups with HAVING

3. **ORDER BY**
   - ASC and DESC ordering
   - Multiple column sorting
   - Expression-based sorting

4. **DISTINCT Keyword**
   - Remove duplicate rows
   - DISTINCT vs ALL (default)

5. **Aggregate Functions**
   - COUNT(*), COUNT(column)
   - SUM(column)
   - AVG(column)
   - MIN(column)
   - MAX(column)
   - Mixed aggregate and non-aggregate queries

6. **Subqueries**
   - Nested SELECT in WHERE clause
   - Subqueries in SELECT list
   - Subqueries in FROM clause

## Implementation Strategy

### Iteration 1: ORDER BY and DISTINCT (Highest Priority)

#### 1.1 Parser Extensions
**Files**: `pkg/sqlparser/parser.go`, `pkg/sqlparser/types.go`

**Tasks**:
- Add `OrderByClause` struct:
  ```go
  type OrderByClause struct {
      Column    string
      Direction string // "ASC" or "DESC"
  }
  ```

- Extend `SelectStatement` struct:
  ```go
  type SelectStatement struct {
      Columns     []string
      Table       string
      WhereClause string
      Distinct    bool
      OrderBy     []OrderByClause
  }
  ```

- Update `parseSelect()` to detect:
  - `DISTINCT` keyword
  - `ORDER BY` clause

- Implement `parseOrderBy()`:
  - Parse column name
  - Parse ASC/DESC direction (default ASC)
  - Handle multiple ORDER BY columns (comma-separated)

**Implementation Example**:
```go
// Check for DISTINCT
distinct := false
if strings.HasPrefix(upperQuery, "SELECT DISTINCT ") {
    distinct = true
    query = strings.TrimPrefix(query, "SELECT DISTINCT ")
    query = "SELECT " + query
}

// Check for ORDER BY
orderByIndex := strings.Index(upperQuery, " ORDER BY ")
if orderByIndex != -1 {
    orderByClause := strings.TrimSpace(query[orderByIndex+9:])
    stmt.OrderBy = parseOrderBy(orderByClause)
    query = strings.TrimSpace(query[:orderByIndex])
}
```

#### 1.2 Executor Extensions
**Files**: `pkg/sqlexecutor/executor.go`

**Tasks**:
- Add ORDER BY logic to `executeSelect()`:
  - After fetching all rows
  - Sort rows based on ORDER BY columns
  - Handle ASC (ascending) vs DESC (descending)

- Add DISTINCT logic:
  - After sorting or before sorting
  - Remove duplicate rows
  - Compare all column values

**Implementation Example**:
```go
// Apply ORDER BY
if len(stmt.OrderBy) > 0 {
    sortRows(resultRows, stmt.OrderBy, columns)
}

// Apply DISTINCT
if stmt.Distinct {
    resultRows = removeDuplicates(resultRows)
}
```

- Implement `sortRows()`:
  - Use Go's `sort.Slice()`
  - Compare based on column values
  - Handle string, numeric, and date comparisons

- Implement `removeDuplicates()`:
  - Use map to track seen rows
  - Or compare consecutive sorted rows

#### 1.3 Test Cases
**File**: `cmd/advancedselecttest/main.go`

**Tests**:
```go
// ORDER BY tests
db.Query("SELECT * FROM products ORDER BY name ASC")
db.Query("SELECT * FROM products ORDER BY price DESC")
db.Query("SELECT * FROM products ORDER BY price DESC, name ASC")

// DISTINCT tests
db.Query("SELECT DISTINCT department FROM employees")
db.Query("SELECT DISTINCT department, title FROM employees")

// Combined tests
db.Query("SELECT DISTINCT department FROM employees ORDER BY department")
```

**Success Criteria**:
- ORDER BY sorts correctly
- ASC sorts ascending
- DESC sorts descending
- Multiple ORDER BY columns work
- DISTINCT removes duplicates
- Combined DISTINCT + ORDER BY works

### Iteration 2: Aggregate Functions

#### 2.1 Parser Extensions
**Files**: `pkg/sqlparser/parser.go`, `pkg/sqlparser/types.go`

**Tasks**:
- Add `AggregateFunction` struct:
  ```go
  type AggregateFunction struct {
      Type    string // "COUNT", "SUM", "AVG", "MIN", "MAX"
      Column  string
      Alias   string // Optional AS alias
  }
  ```

- Extend `SelectStatement`:
  ```go
  type SelectStatement struct {
      Columns           []string
      Aggregates        []AggregateFunction
      IsAggregateQuery  bool
      // ... other fields
  }
  ```

- Update `parseSelect()` to detect aggregates:
  - Pattern: `COUNT(*)`, `COUNT(column)`, `SUM(column)`, etc.
  - Regex to match aggregate function calls

- Implement `parseAggregates()`:
  - Scan column list for aggregate functions
  - Extract function type and column name
  - Handle AS aliases

**Implementation Example**:
```go
// Detect aggregate functions
re := regexp.MustCompile(`(COUNT|SUM|AVG|MIN|MAX)\(\*?\w*\)`)
matches := re.FindAllString(columnsStr, -1)

if len(matches) > 0 {
    stmt.IsAggregateQuery = true
    for _, match := range matches {
        stmt.Aggregates = append(stmt.Aggregates, parseAggregate(match))
    }
}
```

#### 2.2 Executor Extensions
**Files**: `pkg/sqlexecutor/executor.go`

**Tasks**:
- Add aggregate function execution to `executeSelect()`:
  - For GROUP BY queries: Calculate aggregates per group
  - For non-GROUP BY: Calculate aggregates over all rows
  - Handle mixed aggregate + non-aggregate queries (requires GROUP BY)

- Implement aggregate functions:
  - `COUNT(*)`: Count all rows
  - `COUNT(column)`: Count non-NULL values
  - `SUM(column)`: Sum numeric values
  - `AVG(column)`: Average numeric values
  - `MIN(column)`: Minimum value
  - `MAX(column)`: Maximum value

**Implementation Example**:
```go
// Execute aggregates
if stmt.IsAggregateQuery {
    if len(stmt.GroupBy) > 0 {
        // Group aggregation (Iteration 3)
        resultRows = executeGroupByAggregates(resultRows, stmt)
    } else {
        // Simple aggregation over all rows
        resultRows = executeSimpleAggregates(resultRows, stmt)
    }
}
```

- Implement `executeSimpleAggregates()`:
  - Calculate each aggregate function
  - Return single row with aggregate results

- Implement `calculateAggregate()`:
  - Switch on aggregate type
  - Perform calculation
  - Return result

#### 2.3 Test Cases
**File**: `cmd/advancedselecttest/main.go`

**Tests**:
```go
// COUNT tests
db.Query("SELECT COUNT(*) FROM employees")
db.Query("SELECT COUNT(salary) FROM employees")
db.Query("SELECT COUNT(DISTINCT department) FROM employees")

// SUM tests
db.Query("SELECT SUM(salary) FROM employees")
db.Query("SELECT SUM(bonus) FROM employees")

// AVG tests
db.Query("SELECT AVG(salary) FROM employees")

// MIN/MAX tests
db.Query("SELECT MIN(salary), MAX(salary) FROM employees")

// Mixed tests
db.Query("SELECT COUNT(*), AVG(salary) FROM employees")
db.Query("SELECT department, COUNT(*) FROM employees GROUP BY department")
```

**Success Criteria**:
- All aggregate functions work correctly
- NULL handling in aggregates (COUNT excludes NULL, SUM/MIN/MAX ignore NULL)
- Aggregate queries without GROUP BY return single row
- Mixed aggregate + non-aggregate queries handled (require GROUP BY)

### Iteration 3: GROUP BY and HAVING

#### 3.1 Parser Extensions
**Files**: `pkg/sqlparser/parser.go`, `pkg/sqlparser/types.go`

**Tasks**:
- Extend `SelectStatement`:
  ```go
  type SelectStatement struct {
      // ... existing fields
      GroupByClause []string
      HavingClause  string
  }
  ```

- Update `parseSelect()` to detect:
  - `GROUP BY` clause
  - `HAVING` clause

- Implement `parseGroupBy()`:
  - Parse comma-separated column list
  - Handle expressions

- Implement `parseHaving()`:
  - Parse HAVING condition
  - Similar to WHERE clause syntax

**Implementation Example**:
```go
// Check for GROUP BY
groupByIndex := strings.Index(upperQuery, " GROUP BY ")
if groupByIndex != -1 {
    groupByClause := query[groupByIndex+10:]
    // Check for HAVING
    havingIndex := strings.Index(strings.ToUpper(groupByClause), " HAVING ")
    if havingIndex != -1 {
        stmt.HavingClause = strings.TrimSpace(groupByClause[havingIndex+7:])
        groupByClause = strings.TrimSpace(groupByClause[:havingIndex])
    }
    stmt.GroupByClause = parseColumns(groupByClause)
    query = strings.TrimSpace(query[:groupByIndex])
}
```

#### 3.2 Executor Extensions
**Files**: `pkg/sqlexecutor/executor.go`

**Tasks**:
- Implement GROUP BY logic:
  - Group rows by GROUP BY columns
  - Calculate aggregates per group
  - Return one row per group

- Implement `executeGroupByAggregates()`:
  - Create map of groups: key = GROUP BY column values, value = rows
  - For each group:
    - Calculate aggregates
    - Create result row with GROUP BY columns + aggregate results

- Implement HAVING filtering:
  - After GROUP BY aggregation
  - Apply HAVING clause to filter groups
  - Reuse WHERE clause parsing logic for HAVING

**Implementation Example**:
```go
func (e *Executor) executeGroupByAggregates(rows [][]interface{}, stmt *SelectStatement) [][]interface{} {
    // Create groups
    groups := make(map[string][]interface{})
    groupKeys := make([]string, 0)

    for _, row := range rows {
        // Build group key
        key := buildGroupKey(row, stmt.GroupByClause)
        if _, exists := groups[key]; !exists {
            groupKeys = append(groupKeys, key)
        }
        groups[key] = append(groups[key], row)
    }

    // Calculate aggregates per group
    resultRows := make([][]interface{}, 0)
    for _, key := range groupKeys {
        groupRows := groups[key]
        resultRow := make([]interface{}, 0)

        // Add GROUP BY columns
        resultRow = append(resultRow, getGroupByColumns(groupRows[0], stmt.GroupByClause)...)

        // Add aggregates
        for _, agg := range stmt.Aggregates {
            result := calculateAggregate(agg, groupRows)
            resultRow = append(resultRow, result)
        }

        resultRows = append(resultRows, resultRow)
    }

    // Apply HAVING
    if stmt.HavingClause != "" {
        resultRows = filterHaving(resultRows, stmt.HavingClause)
    }

    return resultRows
}
```

#### 3.3 Test Cases
**File**: `cmd/advancedselecttest/main.go`

**Tests**:
```go
// GROUP BY tests
db.Query("SELECT department, COUNT(*) FROM employees GROUP BY department")
db.Query("SELECT department, AVG(salary) FROM employees GROUP BY department")
db.Query("SELECT department, SUM(salary), COUNT(*) FROM employees GROUP BY department")

// HAVING tests
db.Query("SELECT department, COUNT(*) FROM employees GROUP BY department HAVING COUNT(*) > 5")
db.Query("SELECT department, AVG(salary) FROM employees GROUP BY department HAVING AVG(salary) > 70000")

// Combined tests
db.Query("SELECT department, COUNT(*), AVG(salary) FROM employees GROUP BY department ORDER BY department")
```

**Success Criteria**:
- GROUP BY groups rows correctly
- Aggregates calculated per group
- HAVING filters groups correctly
- GROUP BY + HAVING + ORDER BY works together

### Iteration 4: JOIN Support

#### 4.1 Parser Extensions
**Files**: `pkg/sqlparser/parser.go`, `pkg/sqlparser/types.go`

**Tasks**:
- Add `JoinClause` struct:
  ```go
  type JoinClause struct {
      Type      string // "INNER", "LEFT", "RIGHT", "FULL"
      Table     string
      OnClause  string
      Alias     string // Optional table alias
  }
  ```

- Extend `SelectStatement`:
  ```go
  type SelectStatement struct {
      Columns     []string
      Table       string
      Joins       []JoinClause
      WhereClause string
      // ... other fields
  }
  ```

- Update `parseSelect()` to detect JOINs:
  - Pattern: `table1 JOIN table2 ON condition`
  - Handle multiple JOINs
  - Parse JOIN type (INNER, LEFT, RIGHT, FULL)
  - Parse ON clause

- Implement `parseJoins()`:
  - Find JOIN keyword in query
  - Extract join type
  - Extract joined table name
  - Extract ON clause condition

**Implementation Example**:
```go
// Find JOIN clauses
joinPattern := regexp.MustCompile(`(INNER|LEFT|RIGHT|FULL)\s+JOIN\s+(\w+)\s+ON\s+(.+?)(?:\s+(?:WHERE|GROUP|ORDER|$))`)
joinMatches := joinPattern.FindAllStringSubmatch(query, -1)

for _, match := range joinMatches {
    join := JoinClause{
        Type:     match[1],
        Table:    match[2],
        OnClause: match[3],
    }
    stmt.Joins = append(stmt.Joins, join)
}
```

#### 4.2 Executor Extensions
**Files**: `pkg/sqlexecutor/executor.go`

**Tasks**:
- Implement JOIN logic:
  - For simplicity: Fetch all rows from all tables
  - Join rows in memory
  - Apply ON clause to filter joined rows
  - Handle different JOIN types (INNER, LEFT, RIGHT, FULL)

- Implement `executeJoin()`:
  - Fetch rows from primary table
  - Fetch rows from joined table
  - Join rows based on ON clause
  - Return joined result set

**Implementation Example**:
```go
func (e *Executor) executeJoin(stmt *SelectStatement) ([][]interface{}, error) {
    // Fetch primary table rows
    primaryRows, err := e.fetchAll(stmt.Table)
    if err != nil {
        return nil, err
    }

    // For each JOIN
    resultRows := primaryRows
    for _, join := range stmt.Joins {
        // Fetch joined table rows
        joinedRows, err := e.fetchAll(join.Table)
        if err != nil {
            return nil, err
        }

        // Join tables in memory
        resultRows = e.joinTables(resultRows, joinedRows, join)
    }

    return resultRows, nil
}

func (e *Executor) joinTables(left, right [][]interface{}, join JoinClause) [][]interface{} {
    result := make([][]interface{}, 0)

    for _, leftRow := range left {
        joined := false
        for _, rightRow := range right {
            // Check ON clause
            if e.evaluateOnClause(leftRow, rightRow, join.OnClause) {
                // Merge rows
                merged := append(append([]interface{}{}, leftRow...), rightRow...)
                result = append(result, merged)
                joined = true
            }
        }

        // Handle LEFT JOIN (include all left rows)
        if join.Type == "LEFT" && !joined {
            result = append(result, leftRow)
        }
    }

    return result
}
```

#### 4.3 Test Cases
**File**: `cmd/advancedselecttest/main.go`

**Tests**:
```go
// INNER JOIN tests
db.Query("SELECT * FROM users INNER JOIN orders ON users.id = orders.user_id")
db.Query("SELECT users.name, orders.total FROM users INNER JOIN orders ON users.id = orders.user_id")

// LEFT JOIN tests
db.Query("SELECT * FROM users LEFT JOIN orders ON users.id = orders.user_id")

// RIGHT JOIN tests
db.Query("SELECT * FROM users RIGHT JOIN orders ON users.id = orders.user_id")

// Multiple JOINs
db.Query("SELECT * FROM users INNER JOIN orders ON users.id = orders.user_id INNER JOIN products ON orders.product_id = products.id")
```

**Success Criteria**:
- INNER JOIN returns matching rows
- LEFT JOIN returns all left rows + matching right rows
- RIGHT JOIN returns matching left rows + all right rows
- Multiple JOINs work
- ON clause filters correctly
- Column names are properly qualified (table.column)

### Iteration 5: Subqueries (Basic)

#### 5.1 Parser Extensions
**Files**: `pkg/sqlparser/parser.go`, `pkg/sqlparser/types.go`

**Tasks**:
- Update `parseSelect()` to detect nested SELECT:
  - In WHERE clause: `WHERE id IN (SELECT id FROM ...)`
  - In SELECT list: `SELECT (SELECT MAX(salary) FROM employees) FROM departments`
  - In FROM clause: `SELECT * FROM (SELECT * FROM employees) AS e`

- Implement `parseSubquery()`:
  - Recursively parse nested SELECT
  - Handle parentheses matching

**Implementation Example**:
```go
// Detect subquery in WHERE clause
subqueryPattern := regexp.MustCompile(`\([^)]*SELECT[^)]*\)`)
subqueryMatches := subqueryPattern.FindAllString(whereClause, -1)

for _, subquery := range subqueryMatches {
    subStmt, err := parser.Parse(subquery)
    if err != nil {
        continue
    }
    // Execute subquery first
    subqueryResults, _ := e.executeSubquery(subStmt)
}
```

#### 5.2 Executor Extensions
**Files**: `pkg/sqlexecutor/executor.go`

**Tasks**:
- Implement subquery execution:
  - Execute subquery first
  - Use subquery results in main query
  - Handle IN, EXISTS, = comparisons

- Implement `executeSubquery()`:
  - Execute SELECT statement
  - Return results
  - Use results in WHERE clause evaluation

**Implementation Example**:
```go
func (e *Executor) executeSubquery(stmt *SelectStatement) ([][]interface{}, error) {
    // Execute subquery
    return e.ExecuteSelect(stmt.RawQuery)
}
```

#### 5.3 Test Cases
**File**: `cmd/advancedselecttest/main.go`

**Tests**:
```go
// Subquery in WHERE
db.Query("SELECT * FROM employees WHERE salary > (SELECT AVG(salary) FROM employees)")
db.Query("SELECT * FROM employees WHERE department_id IN (SELECT id FROM departments)")

// Subquery in SELECT list
db.Query("SELECT name, (SELECT COUNT(*) FROM orders WHERE user_id = users.id) AS order_count FROM users")
```

**Success Criteria**:
- Subqueries execute correctly
- Subquery results used in main query
- IN, EXISTS, = comparisons work
- Nested queries work (limited depth)

## Implementation Order

1. **Iteration 1** (2-3 hours): ORDER BY and DISTINCT
   - Parser extensions
   - Executor sorting logic
   - DISTINCT deduplication
   - Test cases

2. **Iteration 2** (2-3 hours): Aggregate Functions
   - Parser detection
   - Executor calculation
   - COUNT, SUM, AVG, MIN, MAX
   - Test cases

3. **Iteration 3** (2-3 hours): GROUP BY and HAVING
   - Parser extensions
   - Executor grouping logic
   - HAVING filtering
   - Test cases

4. **Iteration 4** (3-4 hours): JOIN Support
   - Parser detection
   - Executor join logic
   - INNER, LEFT, RIGHT joins
   - Test cases

5. **Iteration 5** (2-3 hours): Subqueries (Basic)
   - Parser detection
   - Executor subquery logic
   - Test cases

**Total Estimated Effort**: 11-16 hours

## Technical Considerations

### Sorting Implementation
- Use `sort.Slice()` for flexible sorting
- Handle different data types (string, int, float)
- Consider NULL values in sorting
- Multi-column sorting with secondary keys

### DISTINCT Implementation
- Use map to track seen rows
- Or compare consecutive sorted rows
- Handle case sensitivity for strings

### Aggregate Implementation
- Use accumulator pattern for SUM, AVG
- MIN/MAX need initial values
- COUNT(*) vs COUNT(column) difference
- Handle NULL values properly

### GROUP BY Implementation
- Use map for grouping
- Key = concatenated GROUP BY column values
- Value = slice of rows in group
- Calculate aggregates per group

### JOIN Implementation
- Simple approach: Join in memory
- Fetch all rows from both tables
- Nested loops for join (O(n*m) complexity)
- For production: Use database JOIN support
- Handle LEFT JOIN (include all left rows)
- Handle RIGHT JOIN (include all right rows)

### Subquery Implementation
- Recursive parsing
- Execute subquery first
- Use results in main query
- Limited depth to prevent infinite recursion

## SQLite Compatibility

### JOINs
- SQLite supports INNER, LEFT OUTER, CROSS joins
- RIGHT OUTER and FULL OUTER not directly supported
- May need workarounds for RIGHT/FULL joins

### Aggregates
- SQLite supports COUNT, SUM, AVG, MIN, MAX
- Can use SQLite's native aggregate functions
- Let SQLite handle aggregation when possible

### GROUP BY
- SQLite supports GROUP BY and HAVING
- Can use SQLite's native GROUP BY
- Let SQLite handle grouping when possible

## Test Coverage

### Unit Tests
- Parser tests for each feature
- Executor tests for each feature
- Edge cases (NULL values, empty sets, etc.)

### Integration Tests
- Combined features (ORDER BY + GROUP BY, etc.)
- Multiple tables and JOINs
- Complex queries

### Performance Tests
- Large datasets
- Multiple JOINs
- Complex grouping

## Success Criteria

### Iteration 1 (ORDER BY + DISTINCT)
- ✓ ORDER BY sorts correctly (ASC, DESC)
- ✓ Multiple ORDER BY columns work
- ✓ DISTINCT removes duplicates
- ✓ Combined DISTINCT + ORDER BY works

### Iteration 2 (Aggregate Functions)
- ✓ COUNT, SUM, AVG, MIN, MAX work
- ✓ Aggregate queries without GROUP BY return single row
- ✓ NULL handling in aggregates correct
- ✓ Mixed aggregate + non-aggregate queries handled

### Iteration 3 (GROUP BY + HAVING)
- ✓ GROUP BY groups rows correctly
- ✓ Aggregates calculated per group
- ✓ HAVING filters groups correctly
- ✓ Combined GROUP BY + HAVING + ORDER BY works

### Iteration 4 (JOINs)
- ✓ INNER JOIN returns matching rows
- ✓ LEFT JOIN returns all left rows
- ✓ RIGHT JOIN returns all right rows
- ✓ Multiple JOINs work
- ✓ ON clause filters correctly

### Iteration 5 (Subqueries)
- ✓ Subqueries execute correctly
- ✓ Subquery results used in main query
- ✓ IN, EXISTS, = comparisons work
- ✓ Nested queries work (limited depth)

## Example Usage After Implementation

### ORDER BY
```sql
SELECT * FROM products ORDER BY price DESC
SELECT * FROM products ORDER BY price DESC, name ASC
```

### DISTINCT
```sql
SELECT DISTINCT department FROM employees
SELECT DISTINCT department, title FROM employees
```

### Aggregates
```sql
SELECT COUNT(*) FROM employees
SELECT SUM(salary) FROM employees
SELECT AVG(salary) FROM employees
SELECT MIN(salary), MAX(salary) FROM employees
```

### GROUP BY
```sql
SELECT department, COUNT(*) FROM employees GROUP BY department
SELECT department, AVG(salary) FROM employees GROUP BY department
```

### HAVING
```sql
SELECT department, COUNT(*) FROM employees GROUP BY department HAVING COUNT(*) > 5
```

### JOINs
```sql
SELECT * FROM users INNER JOIN orders ON users.id = orders.user_id
SELECT * FROM users LEFT JOIN orders ON users.id = orders.user_id
```

### Combined Features
```sql
SELECT department, AVG(salary)
FROM employees
GROUP BY department
HAVING AVG(salary) > 70000
ORDER BY department
```

## Future Enhancements

- UNION and UNION ALL
- INTERSECT and EXCEPT
- Window functions (OVER, PARTITION BY, ROW_NUMBER)
- Common Table Expressions (CTEs)
- Subquery optimization (execute in database when possible)
- Full OUTER JOIN workarounds
- Better JOIN performance (use database JOIN when possible)
