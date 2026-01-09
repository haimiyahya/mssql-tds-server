# MSSQL TDS Server Project Plan

## Project Overview
Build a Microsoft SQL Server-compatible server that implements the TDS (Tabular Data Stream) protocol but uses SQLite as the storage backend. The server will be built in Go, and we'll use the standard Go mssql client library for testing.

## Architecture
- **Server Language**: Go (Golang)
- **Client Language**: Go (using standard mssql driver)
- **Storage Backend**: SQLite
- **Protocol**: TDS (Tabular Data Stream)

## Project Phases

### Phase 1: TDS Protocol Foundation (Proof of Concept)
**Goal**: Establish basic TDS connection between client and server

**Tasks**:
1. Research TDS protocol specifications and packet structure
2. Implement basic TDS packet parser
3. Create a minimal TDS server that:
   - Listens on TCP port (default 1433)
   - Accepts incoming TCP connections
   - Handles TDS pre-login handshake
   - Responds with basic TDS acknowledgment
4. Implement simple TDS packet serialization/deserialization

**Success Criteria**:
- Server starts and listens on configured port
- Standard Go mssql client can establish a connection to the server
- No connection errors or timeouts

### Phase 2: Basic Request/Response Communication
**Goal**: Implement simple bidirectional communication

**Tasks**:
1. Extend server to handle TDS query packets
2. Implement TDS SQL batch command processing
3. Create a simple request handler:
   - Client sends a simple query/string (e.g., "hello")
   - Server processes the request
   - Server returns response in caps (e.g., "HELLO")
4. Implement TDS result set formatting for simple responses
5. Add proper error handling and logging

**Success Criteria**:
- Client can send a simple query to the server
- Server receives and processes the query
- Server returns a properly formatted TDS response
- Client successfully receives and displays the response

### Phase 3: Stored Procedure Handling
**Goal**: Implement RPC (Remote Procedure Call) handling for stored procedure execution

**Tasks**:
1. Implement TDS RPC packet parsing (packet type 0x03)
2. Create stored procedure handler to receive procedure name and parameters
3. Extract parameters from RPC packets
4. Implement parameter data type parsing (int, varchar, etc.)
5. Generate result sets for stored procedure calls
6. Add support for common stored procedure patterns

**Success Criteria**:
- Server receives and parses RPC packets
- Extracts stored procedure name correctly
- Extracts parameters with correct data types
- Returns properly formatted result sets
- Client can successfully call stored procedures

## Stored Procedure Implementation Strategy

### Chosen Approach: Option 2 - Parse at EXEC Time

**Strategy Decision**: Parse and execute stored procedures at EXEC time (not at CREATE time).

**Rationale**:
- Simpler CREATE PROCEDURE implementation (just store string)
- More flexible for dynamic SQL
- Don't need full T-SQL parser upfront
- Easier to implement iteratively
- Allows for progressive enhancement

**Implementation Flow**:
1. CREATE PROCEDURE: Store raw T-SQL string in SQLite + basic validation
2. EXEC: Retrieve procedure body → Parse T-SQL → Replace parameters → Execute SQL

### Stored Procedure Phases

#### **Phase 4: Simple Procedures (MVP)**
**Goal**: Basic single-statement procedures with parameters only

**Subtasks**:
4.1 Initialize SQLite database for procedure storage
   - Create procedures table (name, body, parameters, created_at)
   - Add basic metadata tracking

4.2 Implement CREATE PROCEDURE parser (syntax validation)
   - Parse procedure name
   - Extract parameter list (name, type, default)
   - Basic syntax checking (no full T-SQL parsing yet)
   - Validate parameter names are valid

4.3 Implement procedure storage
   - Store procedure metadata in SQLite
   - Store raw T-SQL body
   - Handle procedure creation errors
   - Support DROP PROCEDURE

4.4 Implement simple parameter replacement engine
   - Parse procedure body for @param references
   - Replace with actual parameter values
   - Handle NULL values
   - Support basic data types (VARCHAR, INT)

4.5 Implement EXEC command handling
   - Parse EXEC syntax: EXEC procname @param1=value1
   - Retrieve procedure from SQLite
   - Validate parameter count and types
   - Replace parameters in SQL
   - Execute SQL in SQLite
   - Return result sets

4.6 Add unit tests
   - Test procedure storage and retrieval
   - Test parameter replacement
   - Test simple SELECT procedures
   - Test error handling

**Success Criteria**:
- Can CREATE PROCEDURE with simple SELECT statement
- Can EXEC procedure with parameters
- Parameters are correctly replaced and executed
- Results are returned properly
- No support for: variables, IF/ELSE, WHILE, #temp tables

**Example procedures supported**:
```sql
CREATE PROCEDURE GetUserById @id INT
AS
SELECT * FROM users WHERE id = @id

CREATE PROCEDURE GetUsersByDept @dept VARCHAR(50)
AS
SELECT name, email FROM users WHERE department = @dept
```

#### **Phase 5: Variables Support**
**Goal**: Add DECLARE and SET/SELECT variable assignment

**Subtasks**:
5.1 Implement variable declaration parser
   - Parse DECLARE @var TYPE statements
   - Support basic types: INT, VARCHAR, BIGINT
   - Initialize variable context

5.2 Implement SET and SELECT variable assignment
   - Parse SET @var = value
   - Parse SELECT @var = column FROM table
   - Store variable values in execution context

5.3 Implement variable reference in SQL
   - Replace @var with actual values
   - Handle variable in WHERE clauses
   - Handle variable in SELECT lists

5.4 Add unit tests
   - Test variable declarations
   - Test variable assignments
   - Test variable usage in procedures

**Example procedures supported**:
```sql
CREATE PROCEDURE GetUserCount @dept VARCHAR(50)
AS
    DECLARE @count INT
    SELECT @count = COUNT(*) FROM users WHERE department = @dept
    SELECT @count as user_count
```

#### **Phase 6: Basic Control Flow**
**Goal**: Add IF/ELSE support (no nesting)

**Subtasks**:
6.1 Implement IF statement parser
   - Parse IF condition THEN/END structure
   - Support basic comparisons (=, >, <, >=, <=, <>)
   - Support AND/OR operators

6.2 Implement ELSE block support
   - Parse ELSE blocks
   - Handle conditional execution paths

6.3 Implement conditional expression evaluator
   - Evaluate conditions with variables
   - Return boolean results

6.4 Add unit tests
   - Test IF statements
   - Test ELSE blocks
   - Test complex conditions

**Example procedures supported**:
```sql
CREATE PROCEDURE CheckUserStatus @id INT
AS
    DECLARE @status VARCHAR(20)
    SELECT @status = status FROM users WHERE id = @id

    IF @status = 'ACTIVE'
        SELECT 'User is active' as message
    ELSE
        SELECT 'User is inactive' as message
```

#### **Phase 7: WHILE Loops**
**Goal**: Add basic WHILE loop support

**Subtasks**:
7.1 Implement WHILE loop parser
   - Parse WHILE condition DO/END structure
   - Handle loop termination

7.2 Implement loop execution engine
   - Execute loop body while condition is true
   - Prevent infinite loops (max iterations)
   - Handle BREAK and CONTINUE

7.3 Add unit tests
   - Test simple WHILE loops
   - Test loop with variable updates
   - Test BREAK/CONTINUE

**Example procedures supported**:
```sql
CREATE PROCEDURE GenerateSampleData @count INT
AS
    DECLARE @i INT
    SET @i = 0

    WHILE @i < @count
        BEGIN
            INSERT INTO sample (name) VALUES ('User ' + CAST(@i AS VARCHAR))
            SET @i = @i + 1
        END
```

#### **Phase 8: Temporary Tables**
**Goal**: Add #temp table support

**Subtasks**:
8.1 Implement temporary table manager
   - Handle CREATE TABLE #temp
   - Track temp table lifetime
   - Support INSERT/SELECT on temp tables

8.2 Integrate temp tables with SQLite
   - Use SQLite in-memory databases or temp tables
   - Handle temp table scope (connection-level)

8.3 Add unit tests
   - Test temp table creation
   - Test temp table operations
   - Test temp table cleanup

**Example procedures supported**:
```sql
CREATE PROCEDURE GetTopUsers @limit INT
AS
    CREATE TABLE #top_users (id INT, name VARCHAR(100), score INT)
    INSERT INTO #top_users
    SELECT TOP @limit id, name, score FROM users
    ORDER BY score DESC
    SELECT * FROM #top_users
```

#### **Phase 9: Transaction Management**
**Goal**: Add BEGIN TRAN, COMMIT, ROLLBACK support

**Subtasks**:
9.1 Implement transaction manager
   - Track transaction state
   - Handle nested transactions

9.2 Integrate with SQLite transactions
   - Map T-SQL transactions to SQLite
   - Handle ROLLBACK

9.3 Add unit tests
   - Test transaction begin/commit
   - Test rollback on error
   - Test nested transactions

**Example procedures supported**:
```sql
CREATE PROCEDURE UpdateUser @id INT, @name VARCHAR(100)
AS
    BEGIN TRAN
        UPDATE users SET name = @name WHERE id = @id
        IF @@ERROR = 0
            COMMIT TRAN
        ELSE
            ROLLBACK TRAN
```

## Future Project Scope (To Be Implemented Separately)
The following phases will be implemented in a future project that will use this project as a foundation:

- **Phase 10**: Exception Handling
  - TRY/CATCH blocks
  - RAISERROR
  - Custom error messages

- **Phase 11**: Cursor Support
  - DECLARE CURSOR
  - FETCH operations
  - CLOSE/DEALLOCATE

- **Phase 12**: Dynamic SQL
  - EXEC('...') statements
  - sp_executesql

- **Phase 13**: TDS Feature Expansion
  - Implement TDS login authentication
  - Support for multiple result sets
  - Proper error code handling

- **Phase 14**: Production Readiness
  - Connection pooling
  - Concurrent request handling
  - Performance optimization
  - Security hardening

## Development Approach

### Proof of Concept Strategy
1. **Start Simple**: Begin with minimal TDS implementation
2. **Iterative**: Add features incrementally
3. **Test-Driven**: Use Go mssql client as primary testing tool
4. **Document Everything**: Keep detailed notes on TDS protocol behavior

### Key Resources
- Microsoft TDS Protocol Documentation
- TDS packet format specifications
- Go mssql driver source code (for understanding client expectations)
- Existing open-source TDS implementations (for reference)

## Technical Considerations

### TDS Protocol Complexity
- TDS is a binary protocol with multiple versions
- Different TDS versions have different features
- Need to support at least TDS 7.0+ for modern mssql clients

### Connection States
1. Pre-login
2. Login
3. Post-login (query processing)
4. Transaction states (future phases)

### Packet Types to Implement
- Pre-login packet
- Login acknowledgment
- SQL batch command
- RPC (Remote Procedure Call) command
- Result set responses
- Error messages (basic)

## Success Metrics
- Phase 1: Client can connect without errors
- Phase 2: Simple echo functionality works end-to-end
- Phase 3: Stored procedure calls work with parameter extraction

## Notes
- This is a proof of concept - focus on functionality over optimization
- Use existing libraries where possible (e.g., for TDS packet parsing if available)
- Prioritize getting the basic connection working before adding complex features
- Keep detailed logs of TDS packet exchanges for debugging
