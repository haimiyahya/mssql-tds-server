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

## Future Project Scope (To Be Implemented Separately)
The following phases will be implemented in a future project that will use this project as a foundation:

- **Phase 4**: SQLite Integration
  - Connect to SQLite database
  - Implement basic SQL parsing
  - Map SQL queries to SQLite operations

- **Phase 5**: TDS Feature Expansion
  - Implement TDS login authentication
  - Support for multiple result sets
  - Proper error code handling
  - Transaction support

- **Phase 6**: Advanced SQL Features
  - SELECT queries
  - INSERT/UPDATE/DELETE operations
  - JOIN support
  - WHERE clauses and filtering

- **Phase 7**: Data Type Support
  - Integer types (INT, BIGINT, SMALLINT, TINYINT)
  - String types (VARCHAR, NVARCHAR, CHAR)
  - Date/Time types
  - Numeric/Decimal types

- **Phase 8**: Production Readiness
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
