# MSSQL TDS Server

A fully-featured Microsoft SQL Server-compatible server implementing the TDS (Tabular Data Stream) protocol with SQLite storage backend.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Phase Status](https://img.shields.io/badge/Phase-43%20Complete-brightgreen.svg)](PHASE43_PROGRESS.md)
![Authentication](https://img.shields.io/badge/Authentication-Complete-brightgreen.svg)
![SSL%2FTLS](https://img.shields.io/badge/SSL%2FTLS-Complete-brightgreen.svg)
![Encryption](https://img.shields.io/badge/Encryption-Complete-brightgreen.svg)
![Database Management](https://img.shields.io/badge/Database%20Management-Complete-brightgreen.svg)
![Trash%2FRecycle%20Bin](https://img.shields.io/badge/Trash%2FRecycle%20Bin-Complete-brightgreen.svg)

## Table of Contents

- [Overview](#overview)
- [Current Status](#current-status)
- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
- [Usage](#usage)
- [Testing](#testing)
- [Documentation](#documentation)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

## Overview

The MSSQL TDS Server is a high-performance, fully-featured database server that implements the Microsoft SQL Server TDS protocol while using SQLite as its storage backend. This project provides a lightweight, embedded alternative to full SQL Server installations while maintaining compatibility with standard MSSQL client libraries and tools.

### Key Benefits

- **100% TDS Protocol Compatible**: Works with standard MSSQL clients and tools
- **Lightweight**: Embedded SQLite storage requires minimal resources
- **Fast**: Optimized query execution with SQLite's query engine
- **Secure**: Built-in SQL injection prevention and parameterized queries
- **Scalable**: Connection pooling and efficient resource management
- **Production Ready**: Comprehensive testing and error handling

### Use Cases

- Development and testing environments
- Embedded applications requiring SQL Server compatibility
- Lightweight database deployments
- SQL Server protocol compatibility testing
- Database prototyping and development

## Current Status

**Project Status**: ✅ Production Ready

**Latest Phase**: Phase 43 - SSL/TLS Encryption (COMPLETE)
- SSL/TLS certificate management
- Self-signed certificate generation
- SSL/TLS listener creation
- Pre-login encryption negotiation (OFF/ON/REQUIRED)
- TLS 1.2/1.3 support
- Configurable cipher suites
- Same mechanism as Microsoft SQL Server
- CREATE DATABASE support
- DROP DATABASE support (moves to recycle bin/trash)
- USE DATABASE support
- sys.databases view support
- Multi-database queries
- Stored procedure and function database-scoped storage

**Progress**: 39 of 39 planned phases complete (100%)

**Last Update**: 2024

### Completed Phases (43/43)

**Core Phases (1-24)
- ✅ Phase 1: Basic SQL Parser and Executor
- ✅ Phase 2: TDS Protocol Implementation
- ✅ Phase 3: Connection Handling
- ✅ Phase 4: Query Execution
- ✅ Phase 5: Result Handling
- ✅ Phase 6: Error Handling
- ✅ Phase 7: Table Operations
- ✅ Phase 8: Data Manipulation
- ✅ Phase 9: Advanced Data Types
- ✅ Phase 10: Multi-Statement Queries
- ✅ Phase 11: Advanced SELECT Features
- ✅ Phase 12: Transaction Management
- ✅ Phase 13: Views Implementation
- ✅ Phase 14: Index Management
- ✅ Phase 15: ALTER TABLE Support
- ✅ Phase 16: Constraint Support
- ✅ Phase 17: Prepared Statements
- ✅ Phase 18: SQL Functions
- ✅ Phase 19: Batch Operations
- ✅ Phase 20: Connection Pooling
- ✅ Phase 21: Error Handling Improvements
- ✅ Phase 22: Performance Monitoring
- ✅ Phase 23: Security Enhancements
- ✅ Phase 24: Documentation and README Update

**Advanced Features Phases (25-43)**
- ✅ Phase 25: EXPLAIN Query Plan Analysis
- ✅ Phase 26: Full-Text Search (FTS)
- ✅ Phase 27: Common Table Expressions (CTE)
- ✅ Phase 28: Window Functions
- ✅ Phase 29: Triggers
- ✅ Phase 30: JSON Functions
- ✅ Phase 31: Advanced Date/Time Functions
- ✅ Phase 32: Geospatial Functions
- ✅ Phase 33: User-Defined Functions (UDF)
- ✅ Phase 34: Database Backup and Restore
- ✅ Phase 35: Data Import/Export
- ✅ Phase 36: Migration Tools
- ✅ Phase 37: Performance Optimization
- ✅ Phase 38: Monitoring and Alerting
- ✅ Phase 39: Database Administration UI
- ✅ Phase 40: Database Management and Trash Support
  - CREATE DATABASE command implementation
  - DROP DATABASE command (moves to recycle bin/trash)
  - USE command implementation
  - sys.databases view implementation
  - Multi-database queries support
  - Database-scoped procedure/function storage
  - Cross-platform trash/recycle bin support (Windows, macOS, Linux)
  - Restore from trash capability
- ✅ Phase 41: Authentication System
  - User/login storage in master.syslogins table
  - Password hashing with bcrypt (cost=10)
  - Authentication management (create, drop, modify users)
  - Login statistics tracking (count, last login)
  - Account status management (disabled, locked)
  - Default login (sa) creation
- ✅ Phase 42: SSL/TLS Encryption Implementation
  - SSL/TLS certificate management
  - Self-signed certificate generation
  - CA-signed certificate support
  - SSL/TLS listener creation
  - Pre-login encryption negotiation (OFF/ON/REQUIRED)
  - TLS 1.2/1.3 support
  - Configurable cipher suites
  - Same mechanism as Microsoft SQL Server## Features

### Core SQL Operations

### Database Management

**Database Commands**
- `CREATE DATABASE` - Create new database
  - Creates SQLite database file in `./data/` directory
  - Creates system tables (sys_objects, sys_columns)
  - Adds entry to database catalog (master.sys_databases)
  - Validates database name (alphanumeric, underscores, hyphens)
  - Prevents creating system databases (master, tempdb, model, msdb)
- `DROP DATABASE` - Delete database
  - **Moves SQLite file to OS recycle bin/trash** (not permanent deletion)
  - Removes entry from database catalog (master.sys_databases)
  - Closes database connections
  - Prevents dropping system databases
  - Prevents dropping currently used database
  - File can be restored from recycle bin/trash
  - Cross-platform support: Windows (Recycle Bin), macOS (Trash), Linux (Trash)
- `USE database_name` - Switch current database context
  - Opens connection to specified database
  - Closes previous database connection
  - Caches database connections for reuse
  - Updates current database state

**Database Catalog**
- `sys.databases` - View all databases
  - Lists all databases (system and user)
  - Shows database metadata (name, ID, state, create_date, file_path, is_system)
  - Supports SELECT queries
  - Returns SQL Server-compatible format
- `master.sys_databases` - Catalog table
  - Stores database metadata
  - Tracks database files and paths
  - Supports database listing and management
- `master.sys_procedures` - Stored procedure catalog
  - Tracks procedures and their databases
  - Maps procedures to specific databases
  - Procedures stored in database they're created in
- `master.sys_functions` - Function catalog
  - Tracks functions and their databases
  - Maps functions to specific databases
  - Functions stored in database they're created in

**Multi-Database Queries**
- Cross-database queries supported
  - Use `DatabaseName.TableName` syntax
  - Server attaches databases on-demand
  - Translates SQL Server syntax (`DatabaseName.schema.TableName`) to SQLite (`DatabaseName.TableName`)
  - Support for JOINs across databases
- Example:
  ```sql
  SELECT u.name, o.order_date
  FROM UserDB1.Users u
  JOIN UserDB2.Orders o ON u.id = o.user_id
  ```

**Database Storage**
- Each database = Separate SQLite file
  - Database files stored in `./data/` directory
  - Filename format: `DatabaseName.db`
  - System databases: `master.db`, `tempdb.db`, `model.db`, `msdb.db`
  - User databases: `UserDB1.db`, `UserDB2.db`, etc.
- Stored procedures stored in specific database
  - Procedures created in database X are stored in database X
  - Database-scoped storage (not server-scoped)
  - Catalog tracks which database owns each procedure
  - Use `USE database_name` before creating procedure
- Functions stored in specific database
  - Functions created in database X are stored in database X
  - Database-scoped storage (not server-scoped)
  - Catalog tracks which database owns each function
  - Use `USE database_name` before creating function

**Trash/Recycle Bin Support**
- DROP DATABASE behavior
  - **Moves database file to OS recycle bin/trash**
  - File can be restored from recycle bin/trash
  - Provides data safety and recovery options
  - Cross-platform support:
    - Windows: Uses PowerShell to move to Recycle Bin
    - macOS: Uses AppleScript to move to Trash
    - Linux: Uses gio or trash-cli to move to Trash
    - Fallback: Manual `./trash` directory with timestamped filenames
- Restore from trash
  - Windows: Recycle Bin → Right-click → Restore
  - macOS: Trash → Right-click → Put Back
  - Linux: Trash → Right-click → Restore
  - Manual: Move file from `./trash` to `./data`
- Empty trash
  - Windows: Right-click Recycle Bin → Empty Recycle Bin
  - macOS: Finder → Empty Trash
  - Linux: Trash → Empty Trash
  - Manual: Delete files from `./trash` directory

**Data Definition Language (DDL)**
- `CREATE TABLE` - Create tables with various data types
- `DROP TABLE` - Remove tables
- `CREATE VIEW` - Create views for complex queries
- `DROP VIEW` - Remove views
- `CREATE INDEX` - Create indexes for performance
- `DROP INDEX` - Remove indexes
- `ALTER TABLE` - Modify table structure
  - `ADD COLUMN` - Add new columns
  - `RENAME TO` - Rename table
  - `RENAME COLUMN` - Rename columns

**Data Manipulation Language (DML)**
- `SELECT` - Query data with advanced features
  - `ORDER BY` - Sort results
  - `DISTINCT` - Remove duplicates
  - `GROUP BY` - Group results
  - `HAVING` - Filter groups
  - `LIMIT` / `OFFSET` - Pagination
  - `WHERE` - Filter results
- `INSERT` - Insert new data
  - Single row insertion
  - Multiple row insertion (batch)
- `UPDATE` - Modify existing data
  - Single table updates
  - Batch updates
- `DELETE` - Remove data
  - Single table deletes
  - Batch deletes

**Data Query Language (DQL)**
- Advanced SELECT features
  - `JOIN` - Inner, Left, Right, Full joins
  - `SUBQUERY` - Nested queries
  - `AGGREGATE FUNCTIONS` - COUNT, SUM, AVG, MIN, MAX
  - `UNION` - Combine query results

**Transaction Control**
- `BEGIN TRANSACTION` - Start transaction
- `COMMIT` - Commit transaction
- `ROLLBACK` - Rollback transaction
- `SAVEPOINT` - Create savepoint
- `RELEASE SAVEPOINT` - Release savepoint

### Advanced Features

**Data Types**
- Integer types: `TINYINT`, `SMALLINT`, `INT`, `BIGINT`
- Floating-point types: `REAL`, `FLOAT`, `DOUBLE`
- String types: `CHAR`, `VARCHAR`, `TEXT`, `NCHAR`, `NVARCHAR`, `NTEXT`
- Binary types: `BINARY`, `VARBINARY`, `IMAGE`
- Boolean type: `BIT`
- Date/Time types: `DATE`, `TIME`, `DATETIME`, `TIMESTAMP`
- Special types: `UNIQUEIDENTIFIER` (GUID)

**Constraints**
- `PRIMARY KEY` - Unique identifier for rows
- `NOT NULL` - Column cannot be NULL
- `UNIQUE` - All values must be unique
- `DEFAULT` - Default value for column
- `CHECK` - Custom validation constraint
- `FOREIGN KEY` - Referential integrity

**Prepared Statements**
- `PREPARE` - Prepare statement for execution
- `EXECUTE` - Execute prepared statement with parameters
- `DEALLOCATE PREPARE` - Release prepared statement
- Parameter binding for security and performance

**SQL Functions**
- String functions: `CONCAT`, `SUBSTRING`, `LENGTH`, `UPPER`, `LOWER`, `TRIM`, `REPLACE`
- Numeric functions: `ABS`, `ROUND`, `FLOOR`, `CEILING`, `POWER`, `SQRT`
- Date/Time functions: `NOW`, `CURRENT_DATE`, `CURRENT_TIME`, `DATEADD`, `DATEDIFF`
- Conditional functions: `CASE WHEN`, `COALESCE`, `NULLIF`, `ISNULL`
- Aggregate functions: `COUNT`, `SUM`, `AVG`, `MIN`, `MAX`

**Batch Operations**
- Multiple row insert in single statement
- Batch UPDATE operations
- Batch DELETE operations
- Multi-statement transactions

**Connection Pooling**
- Automatic connection reuse
- Configurable pool size
- Connection lifetime management
- Thread-safe connection management
- Connection pool statistics

**Performance Monitoring**
- Query performance metrics
- Connection pool monitoring
- Slow query detection
- Performance reporting
- Resource usage tracking

**Error Handling**
- Enhanced error messages with context
- SQL state codes (ANSI SQL standard)
- Error severity levels
- Error categorization
- Detailed error information

\*\*Security\*\*
- SQL injection prevention
- Parameterized queries
- Query sanitization
- Data validation
- Authentication and authorization support

**JSON Functions**
- `JSON_EXTRACT` - Extract values from JSON
- `JSON_OBJECT` - Create JSON object
- `JSON_ARRAY` - Create JSON array
- `JSON_SET` - Set value in JSON
- `JSON_INSERT` - Insert value in JSON
- `JSON_REPLACE` - Replace value in JSON
- `JSON_PATCH` - Merge JSON objects
- `JSON_VALID` - Validate JSON
- `JSON_EACH` - Iterate JSON elements
- `JSON_TREE` - Traverse JSON structure
- `JSON_GROUP_ARRAY` - Aggregate to JSON array
- `JSON_GROUP_OBJECT` - Aggregate to JSON object

**Window Functions**
- `OVER` - Window function clause
- `PARTITION BY` - Partition window by columns
- `ORDER BY` - Order window by columns
- `ROW_NUMBER` - Row number in partition
- `RANK` - Rank in partition
- `DENSE_RANK` - Dense rank in partition
- `NTILE` - Distribute rows into buckets
- `LEAD` - Access next row
- `LAG` - Access previous row
- `FIRST_VALUE` - First value in partition
- `LAST_VALUE` - Last value in partition
- Frame clauses (ROWS, RANGE, GROUPS)

**Common Table Expressions (CTE)**
- `WITH` - Define CTE
- Recursive CTEs
- Multiple CTEs
- CTE in INSERT/UPDATE/DELETE
- CTE nesting
- CTE materialization

**Full-Text Search (FTS)**
- FTS5 virtual tables
- `MATCH` - Full-text search
- `BM25` - Relevance scoring
- Snippets and highlighting
- FTS indexing
- Query expansion
- FTS configuration

**Geospatial Functions**
- `HAVERSINE_DISTANCE` - Calculate distance between points
- `POINT_IN_POLYGON` - Check if point is in polygon
- `BOUNDING_BOX` - Calculate bounding box
- `DISTANCE_SPHERICAL` - Spherical distance calculation
- `DISTANCE_CIRCULAR` - Circular distance calculation
- `CENTROID` - Calculate polygon centroid
- `POLYGON_AREA` - Calculate polygon area
- Geospatial queries
- Spatial joins
- Spatial indexing

**Advanced Date/Time Functions**
- `DATE` - Extract date from datetime
- `TIME` - Extract time from datetime
- `DATETIME` - Create datetime
- `STRFTIME` - Format datetime
- Date arithmetic
- Date/time modifiers
- Julian day functions
- Timezone handling
- Date comparisons
- Age calculations
- Business days calculations

**User-Defined Functions (UDF)**
- Custom scalar functions
- Custom aggregate functions
- Mathematical UDFs
- String UDFs
- Date/Time UDFs
- Conditional UDFs
- Array/List UDFs
- Validation UDFs
- Business logic UDFs
- Data transformation UDFs
- Complex calculation UDFs

**Triggers**
- `CREATE TRIGGER` - Create trigger
- `DROP TRIGGER` - Drop trigger
- `BEFORE` triggers
- `AFTER` triggers
- `FOR EACH ROW` triggers
- `INSERT`, `UPDATE`, `DELETE` triggers
- `OLD` and `NEW` references
- `WHEN` condition
- Trigger error handling

**EXPLAIN Query Plan Analysis**
- `EXPLAIN` - Display query plan
- Query optimization
- Index usage analysis
- Join optimization
- Scan methods
- Plan statistics
- Performance recommendations

**Database Backup and Restore**
- Full backup
- Incremental backup
- Point-in-time recovery
- Backup validation
- Restore validation
- Backup encryption
- Backup compression
- Automated backup
- Backup rotation

**Data Import/Export**
- CSV import/export
- JSON import/export
- Bulk data operations
- Data format validation
- Progress tracking
- Batch processing
- Data transformation

**Migration Tools**
- Schema migration
- Data migration
- Version control
- Migration rollback
- Migration validation
- Migration history tracking
- Migration execution

**Performance Optimization**
- Query optimization
- Index optimization
- Connection pool optimization
- Memory optimization
- Query performance monitoring
- Performance metrics
- Performance tuning recommendations
- Query caching
- Connection monitoring
- Throughput measurement
- Latency measurement

**Monitoring and Alerting**
- Real-time monitoring
- Alert configuration
- Notification channels
- Health checks
- System metrics
- Alert history tracking
- Log aggregation
- Alert rules
- Alert resolution
- Monitoring dashboard

**Database Administration UI**
- Web-based admin interface
- Table management UI
- Query editor UI
- User management UI
- Database configuration
- Data visualization
- System monitoring dashboard
- Query history UI
- Database statistics UI
- Backup/restore UI

## Architecture

### System Architecture

```
.
├── PLAN.md                                   # Detailed project plan
├── README.md                                 # This file
├── PHASE01-39_PROGRESS.md               # Phase progress documents
├── go.mod                                    # Go module definition
├── pkg/                                      # Package libraries
│   ├── sqlparser/                           # SQL statement parser
│   ├── sqlexecutor/                         # SQL statement executor
│   ├── sqlite/                              # SQLite database management
│   ├── tds/                                 # TDS protocol implementation
│   └── ...                                  # Other support packages
└── cmd/                                      # Server and client applications
    ├── server/                              # TDS server implementation
    ├── client/                              # Test client using standard mssql driver
    └── *test/                               # Phase-specific test clients
        ├── selecttest/                       # Advanced SELECT tests
        ├── jointest/                        # JOIN operation tests
        ├── constrainttest/                  # Constraint tests
        ├── functiontest/                   # SQL function tests
        ├── batchtest/                      # Batch operation tests
        ├── pooltest/                       # Connection pool tests
        ├── errortest/                      # Error handling tests
        ├── perftest/                       # Performance monitoring tests
        ├── securitytest/                   # Security validation tests
        ├── explainftstest/                 # EXPLAIN query plan tests
        ├── ftstest/                        # Full-text search tests
        ├── ctetest/                        # CTE tests
        ├── windowtest/                     # Window function tests
        ├── triggertest/                    # Trigger tests
        ├── jsontest/                       # JSON function tests
        ├── datetimetest/                   # Date/time function tests
        ├── geospatialtest/                 # Geospatial function tests
        ├── udftest/                        # UDF tests
        ├── backuptest/                     # Backup/restore tests
        ├── importexporttest/               # Import/export tests
        ├── migrationtest/                  # Migration tools tests
        ├── monitoringtest/                 # Monitoring/alerting tests
        └── admintest/                     # Admin UI tests
```

### Project Structure

```
.
├── PLAN.md                          # Detailed project plan
├── README.md                        # This file
├── PHASE01-23_PROGRESS.md       # Phase progress documents
├── go.mod                           # Go module definition
├── pkg/                             # Package libraries
│   ├── sqlparser/                    # SQL statement parser
│   ├── sqlexecutor/                  # SQL statement executor
│   ├── sqlite/                       # SQLite database management
│   ├── tds/                          # TDS protocol implementation
│   └── ...                           # Other support packages
└── cmd/                             # Server and client applications
    ├── server/                       # TDS server implementation
    ├── client/                       # Test client using standard mssql driver
    └── *test/                        # Phase-specific test clients
        ├── selecttest/                 # Advanced SELECT tests
        ├── jointest/                  # JOIN operation tests
        ├── constrainttest/             # Constraint tests
        ├── functiontest/              # SQL function tests
        ├── batchtest/                 # Batch operation tests
        ├── pooltest/                  # Connection pool tests
        ├── errortest/                 # Error handling tests
        ├── perftest/                  # Performance monitoring tests
        └── securitytest/              # Security validation tests
```

## Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Build from Source

```bash
# Clone repository
git clone https://github.com/yourusername/mssql-tds-server.git
cd mssql-tds-server

# Build server
go build -o bin/server cmd/server/main.go

# Run server
./bin/server
```

### Go Module

```bash
# Add as dependency
go get github.com/yourusername/mssql-tds-server
```

### Docker

```bash
# Build Docker image
docker build -t mssql-tds-server .

# Run container
docker run -p 1433:1433 mssql-tds-server
```

## Usage

### Starting the Server

```bash
# Default configuration (localhost:1433)
./bin/server

# Custom port
./bin/server -port 1434

# Custom database file
./bin/server -db ./data/mssql.db
```

### Connecting with Go

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    
    _ "github.com/microsoft/go-mssqldb"
)

func main() {
    // Build connection string
    connString := "server=127.0.0.1;port=1433;database=;user id=sa;password="
    
    // Connect to database
    db, err := sql.Open("mssql", connString)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Test connection
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Successfully connected to MSSQL TDS Server!")
    
    // Execute query
    rows, err := db.Query("SELECT * FROM users")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    
    // Iterate through results
    for rows.Next() {
        var id int
        var name string
        var email string
        err := rows.Scan(&id, &name, &email)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("ID: %d, Name: %s, Email: %s\n", id, name, email)
    }
}
```

### Connecting with Python

```python
import pymssql

# Connect to server
conn = pymssql.connect(
    server='127.0.0.1',
    port=1433,
    user='sa',
    password='',
    database=''
)

# Create cursor
cursor = conn.cursor()

# Execute query
cursor.execute('SELECT * FROM users')

# Fetch results
for row in cursor:
    print(f"ID: {row[0]}, Name: {row[1]}, Email: {row[2]}")

# Close connection
conn.close()
```

### SQL Examples

**Create Table**
```sql
CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE,
    age INT CHECK (age >= 0 AND age <= 150),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**Insert Data**
```sql
-- Single row
INSERT INTO users (id, name, email, age) 
VALUES (1, 'Alice', 'alice@example.com', 30);

-- Multiple rows
INSERT INTO users (id, name, email, age) 
VALUES 
  (2, 'Bob', 'bob@example.com', 25),
  (3, 'Charlie', 'charlie@example.com', 35);
```

**Query Data**
```sql
-- Simple query
SELECT * FROM users WHERE age > 25;

-- With JOIN
SELECT u.name, o.order_date
FROM users u
INNER JOIN orders o ON u.id = o.user_id;

-- With aggregate functions
SELECT COUNT(*) as total_users, AVG(age) as avg_age
FROM users;

-- With GROUP BY
SELECT name, COUNT(*) as order_count
FROM users
INNER JOIN orders ON users.id = orders.user_id
GROUP BY users.name
HAVING COUNT(*) > 5;

-- With subquery
SELECT * FROM users
WHERE id IN (SELECT user_id FROM orders WHERE total > 1000);
```

**Update Data**
```sql
-- Single update
UPDATE users SET email = 'newemail@example.com' WHERE id = 1;

-- Batch update
UPDATE users SET age = age + 1 WHERE age < 65;
```

**Delete Data**
```sql
-- Single delete
DELETE FROM users WHERE id = 1;

-- Batch delete
DELETE FROM users WHERE created_at < '2020-01-01';
```

**Transactions**
```sql
BEGIN TRANSACTION;

INSERT INTO users (id, name, email) VALUES (10, 'Test', 'test@example.com');
INSERT INTO orders (id, user_id, total) VALUES (100, 10, 500.00);

COMMIT;

-- Or rollback on error
-- ROLLBACK;
```

**Prepared Statements**
```sql
PREPARE get_user FROM 'SELECT * FROM users WHERE id = $id';
EXECUTE get_user USING @id = 1;
EXECUTE get_user USING @id = 2;
DEALLOCATE PREPARE get_user;
```

**Views**
```sql
CREATE VIEW user_orders AS
SELECT u.id, u.name, COUNT(o.id) as order_count
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.name;

SELECT * FROM user_orders;
```

**Indexes**
```sql
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_orders_user_id ON orders(user_id);
```

**Functions**
```sql
-- String functions
SELECT UPPER(name), LOWER(email) FROM users;

-- Numeric functions
SELECT ROUND(price, 2), ABS(quantity) FROM products;

-- Date/Time functions
SELECT NOW(), CURRENT_DATE, CURRENT_TIME;

-- Conditional functions
SELECT name,
  CASE 
    WHEN age >= 18 THEN 'Adult'
    ELSE 'Minor'
  END as status
FROM users;
```

## Testing

The project includes comprehensive test clients for each phase:

```bash
# Run specific test clients
./bin/selecttest       # Advanced SELECT features
./bin/jointest         # JOIN operations
./bin/constrainttest   # Constraints
./bin/functiontest     # SQL functions
./bin/batchtest        # Batch operations
./bin/pooltest         # Connection pooling
./bin/errortest        # Error handling
./bin/perftest         # Performance monitoring
./bin/securitytest     # Security validation
```

### Running All Tests

```bash
# Build all test clients
for test in selecttest jointest constrainttest functiontest batchtest pooltest errortest perftest securitytest explainftstest ctetest windowtest triggertest jsontest datetimetest geospatialtest udftest backuptest importexporttest migrationtest monitoringtest admintest; do
    go build -o bin/$test cmd/$test/main.go
done

# Run tests (requires running server)
for test in selecttest jointest constrainttest functiontest batchtest pooltest errortest perftest securitytest explainftstest ctetest windowtest triggertest jsontest datetimetest geospatialtest udftest backuptest importexporttest migrationtest monitoringtest admintest; do
    ./bin/$test
done
```

### Test Clients

**Core Tests (1-23)**
- `selecttest` - Advanced SELECT features
- `jointest` - JOIN operations
- `constrainttest` - Constraints
- `functiontest` - SQL functions
- `batchtest` - Batch operations
- `pooltest` - Connection pooling
- `errortest` - Error handling
- `perftest` - Performance monitoring
- `securitytest` - Security validation

**Advanced Tests (24-39)**
- `explainftstest` - EXPLAIN query plans
- `ftstest` - Full-text search
- `ctetest` - Common table expressions
- `windowtest` - Window functions
- `triggertest` - Triggers
- `jsontest` - JSON functions
- `datetimetest` - Date/time functions
- `geospatialtest` - Geospatial functions
- `udftest` - User-defined functions
- `backupresttest` - Database backup and restore
- `importexporttest` - Data import/export
- `migrationtest` - Migration tools
- `monitoringtest` - Monitoring and alerting
- `admintest` - Database administration UI## Documentation

- [PLAN.md](PLAN.md) - Overall project plan and roadmap
- [PHASE01-23_PROGRESS.md](PHASE01-23_PROGRESS.md) - Detailed progress for each phase
- [API Documentation](docs/api.md) - API reference (TODO)
- [Architecture Guide](docs/architecture.md) - Architecture details (TODO)

## Roadmap

### Completed ✅

All 23 planned phases are complete:
- Basic SQL Parser and Executor
- TDS Protocol Implementation
- Connection Handling
- Query Execution
- Result Handling
- Error Handling
- Table Operations
- Data Manipulation
- Advanced Data Types
- Multi-Statement Queries
- Advanced SELECT Features
- Transaction Management
- Views Implementation
- Index Management
- ALTER TABLE Support
- Constraint Support
- Prepared Statements
- SQL Functions
- Batch Operations
- Connection Pooling
- Error Handling Improvements
- Performance Monitoring
- Security Enhancements

### Future Enhancements

Potential areas for future development:
- Window functions
- Common Table Expressions (CTE)
- Recursive queries
- Full-text search
- Stored procedures with control flow
- Triggers
- User-defined functions (UDF)
- Row-level security (RLS)
- Column-level encryption
- Advanced authentication (SSPI, Kerberos)
- Replication support
- High availability (HA)
- Performance tuning and optimization
- Query plan analysis
- EXPLAIN command
- Import/Export tools
- Backup and restore
- Database migration tools
- Monitoring and alerting dashboards

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go coding standards (gofmt, golint)
- Write tests for new features
- Update documentation
- Ensure all tests pass

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Microsoft SQL Server for TDS protocol specification
- SQLite for the robust storage engine
- Go community for excellent libraries and tools
- All contributors and testers

## Contact

For questions, issues, or contributions, please contact:
- **Author**: Haimiyahya
- **GitHub**: https://github.com/haimiyahya/mssql-tds-server
- **Issues**: https://github.com/haimiyahya/mssql-tds-server/issues

---

**MSSQL TDS Server** - A lightweight, fully-featured SQL Server-compatible database server.

*Built with Go and SQLite. TDS Protocol Compatible. Production Ready.*
