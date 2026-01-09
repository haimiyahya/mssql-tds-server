# MSSQL TDS Server

A Microsoft SQL Server-compatible server implementing the TDS (Tabular Data Stream) protocol with SQLite storage backend.

## Status
Proof of Concept - Phase 1, 2, 3, 4, 5 & 6

## Overview
This project implements a minimal TDS server that can accept connections from standard Go mssql clients and handle basic request/response communication, including stored procedure support.

## Plan
See [PLAN.md](PLAN.md) for detailed project phases and implementation strategy.

## Project Structure
```
.
├── PLAN.md          # Detailed project plan
├── README.md        # This file
├── go.mod           # Go module definition
├── pkg/             # Package libraries
│   ├── sqlite/       # SQLite database management
│   ├── procedure/    # Stored procedure handling
│   ├── variable/     # Variable context and parsing
│   ├── controlflow/  # Control flow (IF/ELSE) parsing and execution
│   └── tds/          # TDS protocol implementation
└── cmd/             # Server and client applications
    ├── server/       # TDS server implementation
    ├── client/       # Test client using standard mssql driver
    ├── rpctest/      # RPC procedure test client
    ├── proctest/      # Procedure test client
    ├── vartest/       # Variable test client
    └── controltest/  # Control flow test client
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

-- Execute procedure
EXEC GetUserById @id=1

-- Drop procedure
DROP PROCEDURE GetUserById
```

## Development
See [PLAN.md](PLAN.md) for implementation phases and tasks.

## Future Work
This project provides the foundation for a full-featured MSSQL-compatible server. Future phases will be implemented progressively:


- **Phase 7**: WHILE Loops
- **Phase 8**: Temporary Tables (#temp)
- **Phase 9**: Transaction Management (BEGIN TRAN, COMMIT, ROLLBACK)

See [PLAN.md](PLAN.md) for complete roadmap.
