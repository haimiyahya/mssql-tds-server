# MSSQL TDS Server

A Microsoft SQL Server-compatible server implementing the TDS (Tabular Data Stream) protocol with SQLite storage backend.

## Status
Proof of Concept - Phase 1-9 (ALL PHASES COMPLETE)
Phase 10 (Plain T-SQL Script Execution) - COMPLETE
Phase 11 (Advanced SELECT Features) - IN PROGRESS

## Overview
This project implements a minimal TDS server that can accept connections from standard Go mssql clients and handle basic request/response communication, including stored procedure support.

## Plan
See [PLAN.md](PLAN.md) for detailed project phases and implementation strategy.

## Project Structure
```
.
├── PLAN.md                          # Detailed project plan
├── README.md                        # This file
├── PLAIN_SQL_EXECUTION_PLAN.md       # Phase 10 implementation plan
├── ADVANCED_SELECT_PLAN.md           # Phase 11 implementation plan
├── go.mod                           # Go module definition
├── pkg/                             # Package libraries
│   ├── sqlite/                       # SQLite database management
│   ├── procedure/                    # Stored procedure handling
│   ├── variable/                     # Variable context and parsing
│   ├── controlflow/                  # Control flow (IF/ELSE) parsing and execution
│   ├── temp/                         # Temporary table management
│   ├── transaction/                  # Transaction management
│   ├── sqlparser/                    # SQL statement parser (Phase 10)
│   ├── sqlexecutor/                  # SQL statement executor (Phase 10)
│   └── tds/                         # TDS protocol implementation
└── cmd/                             # Server and client applications
    ├── server/                       # TDS server implementation
    ├── client/                       # Test client using standard mssql driver
    ├── rpctest/                      # RPC procedure test client
    ├── proctest/                     # Procedure test client
    ├── vartest/                      # Variable test client
    ├── controltest/                  # Control flow test client
    ├── whiletest/                    # WHILE loop test client
    ├── temptest/                     # Temporary table test client
    ├── trantest/                     # Transaction test client
    └── plainsqltest/                 # Plain SQL test client (Phase 10)
```

## Completed Phases

### Phase 1: TDS Protocol Foundation ✅
- Basic TDS connection and packet parsing
- Pre-login handshake handling
- Login acknowledgment

### Phase 2: Basic Request/Response Communication ✅
- SQL batch command handling
- Simple query processor with echo functionality
- Result set formatting

### Phase 3: Stored Procedure Handling (RPC) ✅
- RPC packet parsing and handling
- Parameter extraction with data type parsing
- Multiple stored procedure implementations (SP_HELLO, SP_ECHO, SP_GET_DATA)

### Phase 4: Simple Stored Procedures (MVP) ✅
- SQLite database initialization for procedure storage
- CREATE PROCEDURE parser (syntax validation)
- Procedure storage in SQLite
- Simple parameter replacement engine
- EXEC command handling
- DROP PROCEDURE support

### Phase 5: Variables Support ✅
- Variable declaration parser (DECLARE @var TYPE)
- Variable context management for procedure execution
- SET variable assignment
- SELECT variable assignment
- Variable reference in SQL replacement
- Support for basic types (INT, VARCHAR, BIGINT, BIT, etc.)
- Error handling for undeclared/duplicate variables

**Current Capabilities:**
- Create stored procedures with simple SQL statements
- Execute procedures with parameters (INT, VARCHAR, BIT)
- Parameter replacement in SQL statements
- Drop stored procedures
- Error handling for missing/invalid procedures
- DECLARE variables with various data types
- SET variable values
- SELECT variable assignment from queries
- Variable reference in SQL (WHERE, SELECT lists, etc.)

### Phase 6: Basic Control Flow (IF/ELSE) ✅
- IF statement parsing (IF condition THEN statements END)
- ELSE block support (IF condition THEN statements ELSE statements END)
- Condition evaluation with variables
- Support for comparison operators (=, <>, >, <, >=, <=)
- Support for logical operators (AND, OR)
- Conditional block execution
- Error handling for invalid conditions

### Phase 7: WHILE Loops ✅
- WHILE statement parsing (WHILE condition statements END)
- Loop condition evaluation with variables
- Loop body execution
- Maximum iteration protection (1000 iterations)
- Support for BREAK/CONTINUE (basic)
- Error handling for infinite loops

### Phase 8: Temporary Tables (#temp) ✅
- CREATE TABLE #temp statement parsing
- Temporary table creation and in-memory storage
- INSERT INTO #temp operations
- SELECT FROM #temp operations
- Session management for temp tables
- Automatic cleanup on session end
- Temp table name resolution (#temp → internal)
- Basic UPDATE and DELETE support

### Phase 9: Transaction Management (BEGIN TRAN, COMMIT, ROLLBACK) ✅
- Transaction statement parsing (BEGIN TRAN, COMMIT, ROLLBACK)
- Transaction context management
- SQLite transaction support
- BEGIN TRANSACTION handling
- COMMIT handling
- ROLLBACK handling
- Automatic rollback on errors
- Transaction isolation

**Example Usage:**
```sql
-- Simple procedure (Phase 4)
CREATE PROCEDURE GetUserById @id INT
AS SELECT * FROM users WHERE id = @id

-- Procedure with variables (Phase 5)
CREATE PROCEDURE GetUserCount @dept VARCHAR(50)
AS
    DECLARE @count INT
    SELECT @count = COUNT(*) FROM users WHERE department = @dept
    SELECT @count as user_count

-- Procedure with IF/ELSE (Phase 6)
CREATE PROCEDURE CheckUserStatus @id INT
AS
    DECLARE @status VARCHAR(20)
    SELECT @status = status FROM users WHERE id = @id

    IF @status = 'ACTIVE'
        SELECT 'User is active' as message
    ELSE
        SELECT 'User is inactive' as message

-- Procedure with WHILE loop (Phase 7)
CREATE PROCEDURE CountToTen AS
    DECLARE @i INT
    SET @i = 1

    WHILE @i <= 10
        SELECT @i as number
        SET @i = @i + 1

-- Procedure with temporary tables (Phase 8)
CREATE PROCEDURE ProcessResults AS
    CREATE TABLE #results (id INT, name VARCHAR(50))

    INSERT INTO #results VALUES (1, 'Alice')
    INSERT INTO #results VALUES (2, 'Bob')

    SELECT * FROM #results

-- Procedure with transactions (Phase 9)
CREATE PROCEDURE SafeInsert AS
    BEGIN TRANSACTION

    INSERT INTO users (id, name, email) VALUES (100, 'Test', 'test@example.com')
    SELECT 'User inserted in transaction' as message

    COMMIT

-- Execute procedure
EXEC GetUserById @id=1

-- Drop procedure
DROP PROCEDURE GetUserById
```

## Development
See [PLAN.md](PLAN.md) for implementation phases and tasks.

## Future Work
This project provides the foundation for a full-featured MSSQL-compatible server. Future phases will be implemented progressively:

**All planned phases (1-9) are now complete!**

The server currently supports:
- TDS protocol communication
- Stored procedures with parameters
- Variables (DECLARE, SET, SELECT @var)
- Control flow (IF/ELSE, WHILE loops)
- Temporary tables (#temp)
- Transactions (BEGIN TRAN, COMMIT, ROLLBACK)

**Phase 10 - Plain T-SQL Script Execution (In Progress)**

Adding support for executing plain T-SQL scripts directly without wrapping them in stored procedures. This includes:
- Direct SELECT queries
- INSERT, UPDATE, DELETE statements
- CREATE TABLE, DROP TABLE statements
- Proper result set formatting
- Error handling for SQL execution

See [PLAIN_SQL_EXECUTION_PLAN.md](PLAIN_SQL_EXECUTION_PLAN.md) for detailed implementation plan.

See [PLAN.md](PLAN.md) for complete roadmap including future phases.
