# Plain T-SQL Script Execution - Implementation Plan

## Overview
Add support for executing plain T-SQL scripts directly without wrapping them in stored procedures. This is a fundamental MSSQL Server usage pattern that's currently missing.

## Current State
- The `QueryProcessor` only implements a simple echo function (converts queries to uppercase)
- SQL batch handling exists but doesn't execute actual SQL
- The `Executor` already has robust SQL execution logic for stored procedures
- No direct SELECT, INSERT, UPDATE, DELETE execution support

## Problem Statement
Users expect to be able to:
```sql
-- Direct table queries
SELECT * FROM users WHERE id = 1

-- Data modifications
INSERT INTO users (name, email) VALUES ('John', 'john@example.com')
UPDATE users SET name = 'Jane' WHERE id = 1
DELETE FROM users WHERE id = 1

-- DDL operations
CREATE TABLE test (id INT, name VARCHAR(50))
DROP TABLE test
```

Currently, these only work if wrapped in stored procedures, which is not a natural usage pattern.

## Implementation Strategy

### Phase 10: Basic SQL Statement Execution (MVP)

#### 10.1 SQL Statement Parser
**Goal**: Parse T-SQL statements to determine statement type and content

**Tasks**:
- Create `pkg/sqlparser` package
- Implement statement type detection:
  - SELECT statements
  - INSERT statements
  - UPDATE statements
  - DELETE statements
  - CREATE TABLE statements
  - DROP TABLE statements
- Extract table names and column names
- Extract WHERE clause conditions
- Handle statement termination (GO, semicolons)

**Files to Create**:
- `pkg/sqlparser/parser.go` - Main parser logic
- `pkg/sqlparser/types.go` - Statement type definitions
- `pkg/sqlparser/parser_test.go` - Unit tests

**Files to Modify**:
- `pkg/tds/query.go` - Update QueryProcessor

#### 10.2 SQL Statement Executor
**Goal**: Execute parsed SQL statements against SQLite

**Tasks**:
- Create `pkg/sqlexecutor` package
- Implement SELECT statement execution:
  - Parse SELECT query
  - Execute against SQLite
  - Format result sets as TDS responses
- Implement INSERT statement execution:
  - Parse INSERT statement
  - Execute against SQLite
  - Return affected row count
- Implement UPDATE statement execution:
  - Parse UPDATE statement
  - Execute against SQLite
  - Return affected row count
- Implement DELETE statement execution:
  - Parse DELETE statement
  - Execute against SQLite
  - Return affected row count
- Implement CREATE TABLE execution:
  - Parse CREATE TABLE statement
  - Execute against SQLite
  - Return success message
- Implement DROP TABLE execution:
  - Parse DROP TABLE statement
  - Execute against SQLite
  - Return success message

**Files to Create**:
- `pkg/sqlexecutor/executor.go` - Main executor logic
- `pkg/sqlexecutor/select.go` - SELECT statement handling
- `pkg/sqlexecutor/insert.go` - INSERT statement handling
- `pkg/sqlexecutor/update.go` - UPDATE statement handling
- `pkg/sqlexecutor/delete.go` - DELETE statement handling
- `pkg/sqlexecutor/ddl.go` - DDL statement handling
- `pkg/sqlexecutor/executor_test.go` - Unit tests

**Files to Modify**:
- `pkg/tds/query.go` - Update QueryProcessor to use executor

#### 10.3 Update QueryProcessor
**Goal**: Replace echo functionality with actual SQL execution

**Tasks**:
- Remove simple echo functionality
- Integrate SQL parser
- Integrate SQL executor
- Handle different statement types
- Return properly formatted TDS result sets
- Add comprehensive error handling

**Files to Modify**:
- `pkg/tds/query.go` - Complete rewrite of QueryProcessor

#### 10.4 Result Set Formatting
**Goal**: Format SQLite results as TDS result sets

**Tasks**:
- Convert SQLite column types to TDS column types
- Format column metadata (COLMETADATA token)
- Format row data (ROW token)
- Handle NULL values
- Handle different data types:
  - INTEGER → INT
  - TEXT → VARCHAR/NVARCHAR
  - REAL → FLOAT/DECIMAL
  - BLOB → VARBINARY
- Format row count (DONE token)

**Files to Modify**:
- `pkg/tds/packet.go` - Add result set formatting helpers
- `cmd/server/main.go` - Update buildResultPacket to handle more formats

#### 10.5 Error Handling
**Goal**: Provide meaningful error messages for SQL execution failures

**Tasks**:
- Parse SQLite errors
- Convert to TDS ERROR tokens
- Handle syntax errors
- Handle runtime errors (table not found, constraint violations, etc.)
- Format error messages in MSSQL-compatible format

**Files to Modify**:
- `pkg/tds/packet.go` - Improve error packet building
- `pkg/sqlexecutor/executor.go` - Add error mapping

#### 10.6 Integration Testing
**Goal**: Create comprehensive test suite for plain SQL execution

**Tasks**:
- Create test client application (`cmd/plainsqltest/main.go`)
- Test SELECT statements
- Test INSERT statements
- Test UPDATE statements
- Test DELETE statements
- Test CREATE TABLE statements
- Test DROP TABLE statements
- Test error handling
- Test with various data types

**Files to Create**:
- `cmd/plainsqltest/main.go` - Test client
- `test_phase10.sh` - Test script

**Files to Modify**:
- `PLAN.md` - Add Phase 10 documentation
- `README.md` - Update with new capabilities

## Success Criteria

### Phase 10 Success Metrics
- ✓ SELECT queries work and return proper result sets
- ✓ INSERT statements work and return affected row count
- ✓ UPDATE statements work and return affected row count
- ✓ DELETE statements work and return affected row count
- ✓ CREATE TABLE works
- ✓ DROP TABLE works
- ✓ Errors are properly formatted and returned
- ✓ All data types are handled correctly
- ✓ Result sets match MSSQL format expectations

## Implementation Order

### Iteration 1: Core SELECT Support (Highest Priority)
1. Create basic SQL parser for SELECT statements
2. Implement SELECT executor
3. Update QueryProcessor
4. Create basic test for SELECT
5. Verify SELECT works with existing test tables

### Iteration 2: DML Support (INSERT, UPDATE, DELETE)
1. Extend parser for INSERT, UPDATE, DELETE
2. Implement INSERT executor
3. Implement UPDATE executor
4. Implement DELETE executor
5. Add tests for DML operations

### Iteration 3: DDL Support (CREATE, DROP)
1. Extend parser for DDL statements
2. Implement CREATE TABLE executor
3. Implement DROP TABLE executor
4. Add tests for DDL operations

### Iteration 4: Polish and Integration
1. Improve error handling
2. Add comprehensive result set formatting
3. Run all existing tests to ensure no regressions
4. Add documentation
5. Create final test suite

## Technical Considerations

### SQL Dialect Differences
- T-SQL uses square brackets for identifiers: `[table]` vs SQLite's `table` or `"table"`
- T-SQL uses `IDENTITY` vs SQLite's `AUTOINCREMENT`
- T-SQL uses `DATETIME` vs SQLite's `TEXT` or `INTEGER` (timestamps)
- Need to handle these differences in the parser/executor

### Data Type Mapping
```
T-SQL Type → SQLite Type
INT         → INTEGER
VARCHAR(n)  → TEXT
NVARCHAR(n) → TEXT
CHAR(n)     → TEXT
NCHAR(n)    → TEXT
BIT         → INTEGER
FLOAT       → REAL
REAL        → REAL
DECIMAL     → REAL
MONEY       → REAL
DATETIME    → TEXT (ISO8601 format)
DATE        → TEXT (ISO8601 format)
TIME        → TEXT
VARBINARY   → BLOB
```

### Result Set Format
TDS result sets use tokens:
- `COLMETADATA` (0x81) - Column definitions
- `ROW` (0xD1) - Row data
- `DONE` (0xFD) - End of result set

Need to properly format:
- Column names
- Column types
- Column values
- Row count

## Future Enhancements (Beyond Phase 10)

### Phase 10.5: Advanced SELECT Features
- JOIN support
- GROUP BY, HAVING
- ORDER BY
- DISTINCT
- Aggregate functions (COUNT, SUM, AVG, etc.)
- Subqueries

### Phase 10.6: Complex DML
- Multi-row INSERT
- INSERT with SELECT
- UPDATE with JOIN
- DELETE with JOIN

### Phase 10.7: Transaction Support in Plain SQL
- BEGIN TRANSACTION
- COMMIT
- ROLLBACK
- Named transactions

### Phase 10.8: Schema Operations
- ALTER TABLE
- CREATE INDEX
- DROP INDEX
- Schema queries (INFORMATION_SCHEMA)

### Phase 10.9: Advanced Features
- CTEs (Common Table Expressions)
- Views
- Functions
- Triggers

## Example Usage After Implementation

### Basic SELECT
```go
db.Query("SELECT * FROM users")
// Returns all rows from users table
```

### SELECT with WHERE
```go
db.Query("SELECT name, email FROM users WHERE id = 1")
// Returns specific columns for user with id=1
```

### INSERT
```go
db.Exec("INSERT INTO users (name, email) VALUES ('John', 'john@example.com')")
// Returns: 1 row affected
```

### UPDATE
```go
db.Exec("UPDATE users SET name = 'Jane' WHERE id = 1")
// Returns: 1 row affected
```

### DELETE
```go
db.Exec("DELETE FROM users WHERE id = 1")
// Returns: 1 row affected
```

### CREATE TABLE
```go
db.Exec("CREATE TABLE products (id INT, name VARCHAR(50), price REAL)")
// Returns: Table created successfully
```

### DROP TABLE
```go
db.Exec("DROP TABLE products")
// Returns: Table dropped successfully
```

## Dependencies
- No new external dependencies required
- Uses existing `database/sql` and `github.com/mattn/go-sqlite3`

## Estimated Effort
- Iteration 1 (SELECT): 4-6 hours
- Iteration 2 (DML): 3-4 hours
- Iteration 3 (DDL): 2-3 hours
- Iteration 4 (Polish): 2-3 hours

**Total**: 11-16 hours

## Backward Compatibility
- All existing stored procedure functionality remains unchanged
- Existing tests continue to work
- No breaking changes to existing APIs
- Pure additive feature

## Testing Strategy
1. Unit tests for parser
2. Unit tests for executor
3. Integration tests with test client
4. Regression tests for existing functionality
5. Manual testing with standard mssql client tools

## Notes
- Start with minimal viable implementation (basic SELECT)
- Focus on common use cases first
- Add complexity incrementally
- Test thoroughly at each iteration
- Keep code simple and maintainable
- Reuse existing patterns from stored procedure executor
