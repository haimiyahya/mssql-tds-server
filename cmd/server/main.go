package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/factory/mssql-tds-server/pkg/auth"
	"github.com/factory/mssql-tds-server/pkg/database"
	"github.com/factory/mssql-tds-server/pkg/procedure"
	"github.com/factory/mssql-tds-server/pkg/sqlexecutor"
	"github.com/factory/mssql-tds-server/pkg/sqlite"
	"github.com/factory/mssql-tds-server/pkg/tds"
)

const (
	defaultPort = 1433
)

type Server struct {
	addr                 string
	db                   *sqlite.Database
	catalog              *database.Catalog
	procedureStorage      *procedure.Storage
	procedureExecutor     *procedure.Executor
	queryProcessor       *tds.QueryProcessor
	storedProcedureHandler *tds.StoredProcedureHandler
	sqlExecutor          *sqlexecutor.Executor
	authManager          *auth.AuthManager
}

func NewServer(port int, dbPath string) (*Server, error) {
	// Initialize SQLite database
	db, err := sqlite.NewDatabase(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// Initialize database tables
	err = db.Initialize()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create procedure storage
	procStorage, err := procedure.NewStorage(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create procedure storage: %w", err)
	}

	// Create procedure executor
	procExecutor, err := procedure.NewExecutor(db, procStorage)
	if err != nil {
		return nil, fmt.Errorf("failed to create procedure executor: %w", err)
	}

	// Create SQL executor for plain SQL execution
	// Initialize database catalog
	catalog := database.NewCatalog("./data", db.GetDB())
	if catalog == nil {
		return nil, fmt.Errorf("failed to create database catalog")
	}

	sqlExec := sqlexecutor.NewExecutor(db.GetDB(), catalog)

	// Create query processor and set SQL executor
	queryProc := tds.NewQueryProcessor()
	queryProc.SetExecutor(sqlExec)

	// Create authentication manager
	authMgr, err := auth.NewAuthManager(db.GetDB())
	if err != nil {
		return nil, fmt.Errorf("failed to create authentication manager: %w", err)
	}

	// Initialize default logins (sa)
	err = authMgr.InitializeDefaultLogins()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize default logins: %w", err)
	}

	return &Server{
		addr:                 fmt.Sprintf(":%d", port),
		db:                   db,
		procedureStorage:      procStorage,
		procedureExecutor:     procExecutor,
		queryProcessor:       queryProc,
		storedProcedureHandler: tds.NewStoredProcedureHandler(),
		sqlExecutor:          sqlExec,
		authManager:          authMgr,
	}, nil
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer listener.Close()

	log.Printf("TDS Server listening on %s", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("New connection from %s", conn.RemoteAddr())

	// Read first packet
	packet, err := s.readPacket(conn)
	if err != nil {
		log.Printf("Error reading packet: %v", err)
		return
	}

	log.Printf("Received packet: Type=%#02x, Status=%#02x, Length=%d",
		packet.Header.Type, packet.Header.Status, packet.Header.Length)

	// Handle pre-login request
	if packet.Header.Type == tds.PacketTypePreLogin {
		err = s.handlePreLogin(conn, packet)
		if err != nil {
			log.Printf("Error handling pre-login: %v", err)
			return
		}
	} else if packet.Header.Type == tds.PacketTypeLogin || packet.Header.Type == 0x12 {
		// Handle direct login (no pre-login)
		log.Printf("Handling direct login packet (type=%#02x)", packet.Header.Type)
		err = s.handleLogin(conn, packet)
		if err != nil {
			log.Printf("Error handling login: %v", err)
			return
		}
	}

	// Read subsequent packets
	for {
		packet, err = s.readPacket(conn)
		if err != nil {
			log.Printf("Error reading packet: %v", err)
			break
		}

		log.Printf("Received packet: Type=%#02x, Status=%#02x, Length=%d",
			packet.Header.Type, packet.Header.Status, packet.Header.Length)

		log.Printf("PacketTypeLogin=%#02x, comparing with %#02x", tds.PacketTypeLogin, packet.Header.Type)

		// Handle login request (TDS 7.0+ uses 0x10, some clients use 0x12)
		if packet.Header.Type == tds.PacketTypeLogin || packet.Header.Type == 0x12 {
			log.Printf("Handling login packet (type=%#02x)", packet.Header.Type)
			err = s.handleLogin(conn, packet)
			if err != nil {
				log.Printf("Error handling login: %v", err)
				break
			}
		} else if packet.Header.Type == tds.PacketTypeRPC {
			// Handle RPC (Remote Procedure Call)
			log.Printf("Handling RPC packet")
			err = s.handleRPC(conn, packet)
			if err != nil {
				log.Printf("Error handling RPC: %v", err)
				break
			}
		} else if packet.Header.Type == tds.PacketTypeSQLBatch {
			// Handle SQL batch
			err = s.handleSQLBatch(conn, packet)
			if err != nil {
				log.Printf("Error handling SQL batch: %v", err)
				break
			}
		} else {
			log.Printf("Unknown packet type %#02x, skipping...", packet.Header.Type)
		}
	}

	log.Printf("Connection closed from %s", conn.RemoteAddr())
}

func (s *Server) readPacket(conn net.Conn) (*tds.Packet, error) {
	// Read packet header first
	headerBuf := make([]byte, 8)
	_, err := conn.Read(headerBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	header, err := tds.ParseHeader(headerBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse header: %w", err)
	}

	// Read packet data
	dataSize := int(header.Length) - 8
	dataBuf := make([]byte, dataSize)
	if dataSize > 0 {
		_, err = conn.Read(dataBuf)
		if err != nil {
			return nil, fmt.Errorf("failed to read data: %w", err)
		}
	}

	return &tds.Packet{
		Header: header,
		Data:   dataBuf,
	}, nil
}

func (s *Server) writePacket(conn net.Conn, packet *tds.Packet) error {
	data := packet.Serialize()
	_, err := conn.Write(data)
	return err
}

func (s *Server) handlePreLogin(conn net.Conn, packet *tds.Packet) error {
	log.Println("Handling pre-login request")

	// Parse pre-login request
	req, err := tds.ParsePreLoginRequest(packet.Data)
	if err != nil {
		return fmt.Errorf("failed to parse pre-login request: %w", err)
	}

	log.Printf("Pre-login request: Version=%#v, Encryption=%#02x, Instance=%s",
		req.Version, req.Encryption, req.Instance)

	// Create pre-login response
	resp := tds.DefaultPreLoginResponse()

	// Serialize response
	respData := tds.SerializePreLoginResponse(resp)

	// Send response packet
	respPacket := tds.NewPacket(tds.PacketTypePreLogin, tds.StatusEOM, 1, respData)
	err = s.writePacket(conn, respPacket)
	if err != nil {
		return fmt.Errorf("failed to send pre-login response: %w", err)
	}

	log.Println("Sent pre-login response")
	return nil
}

func (s *Server) handleLogin(conn net.Conn, packet *tds.Packet) error {
	log.Println("Handling login request")

	// Parse login packet to extract username and password
	// Note: TDS login packet format is complex
	// For now, we'll use a simple approach:
	// - Allow all logins (for compatibility)
	// - In production, parse login packet and authenticate
	
	// TODO: Parse TDS login packet (0x10 or 0x12) to extract:
	// - Username (offset + length)
	// - Password (offset + length)
	// - Then authenticate using s.authManager.AuthenticateLogin(username, password)
	
	// For now, we'll accept any login (backwards compatibility)
	// You can create logins using:
	// INSERT INTO master.syslogins (name, password_hash) VALUES ('user', 'hash')
	
	// Send login acknowledgment
	loginAck := s.buildLoginAckPacket()

	err := s.writePacket(conn, loginAck)
	if err != nil {
		return fmt.Errorf("failed to send login ack: %w", err)
	}

	// Send DONE token to indicate login complete
	donePacket := s.buildDonePacket()
	err = s.writePacket(conn, donePacket)
	if err != nil {
		return fmt.Errorf("failed to send done packet: %w", err)
	}

	log.Println("Sent login acknowledgment and done token")
	return nil
}

func (s *Server) buildLoginAckPacket() *tds.Packet {
	// Build a simple login acknowledgment
	// Format: [TokenType(0xAD)][Length][Interface][TDSVersion][ProgName][Version][...]
	var buf []byte

	buf = append(buf, 0xAD) // Token type: LOGINACK
	buf = append(buf, 0x01) // Length (byte)
	buf = append(buf, 0x00) // Interface (Native SQL)
	// TDS Version
	buf = append(buf, 0x09, 0x00, 0x00, 0x00) // TDS 7.3
	// ProgName length
	buf = append(buf, 0x00, 0x00) // 0 bytes
	// Version
	buf = append(buf, 0x09, 0x00, 0x00, 0x00)

	return tds.NewPacket(tds.PacketTypeTabular, tds.StatusEOM, 2, buf)
}

func (s *Server) buildDonePacket() *tds.Packet {
	// Build a DONE token to indicate login success
	var buf []byte

	buf = append(buf, 0xFD) // Token type: DONE
	buf = append(buf, 0x00, 0x10) // Status: FINAL
	buf = append(buf, 0x00, 0x00) // Current command
	buf = append(buf, 0x00, 0x00, 0x00, 0x00) // Row count (0)

	return tds.NewPacket(tds.PacketTypeTabular, tds.StatusEOM, 2, buf)
}

func (s *Server) buildErrorPacket(err error) *tds.Packet {
	var buf []byte

	// Token type: ERROR (0xAA)
	buf = append(buf, 0xAA)
	// Token length (variable)
	// Error number (4 bytes)
	buf = append(buf, 0x00, 0x00, 0x00, 0x01) // Error number 1
	// Error state (1 byte)
	buf = append(buf, 0x01)
	// Error severity (1 byte)
	buf = append(buf, 0x0A) // Severity 10
	// Error message length (2 bytes)
	errorMsg := err.Error()
	msgLen := uint16(len(errorMsg))
	buf = append(buf, byte(msgLen>>8), byte(msgLen))
	// Error message
	buf = append(buf, []byte(errorMsg)...)
	// Server name length (1 byte)
	buf = append(buf, 0x00) // Empty server name
	// Procedure name length (1 byte)
	buf = append(buf, 0x00) // Empty procedure name
	// Line number (4 bytes)
	buf = append(buf, 0x00, 0x00, 0x00, 0x00)

	return tds.NewPacket(tds.PacketTypeTabular, tds.StatusEOM, 3, buf)
}

func (s *Server) handleSQLBatch(conn net.Conn, packet *tds.Packet) error {
	query := string(packet.Data)
	log.Printf("Handling SQL batch: %s", query)

	// Normalize query
	query = strings.TrimSpace(query)
	queryUpper := strings.ToUpper(query)

	// Check for CREATE PROCEDURE
	if strings.HasPrefix(queryUpper, "CREATE PROCEDURE") || strings.HasPrefix(queryUpper, "CREATE PROC") {
		return s.handleCreateProcedure(conn, query)
	}

	// Check for DROP PROCEDURE
	if strings.HasPrefix(queryUpper, "DROP PROCEDURE") || strings.HasPrefix(queryUpper, "DROP PROC") {
		return s.handleDropProcedure(conn, query)
	}

	// Check for EXEC command
	if strings.HasPrefix(queryUpper, "EXEC ") || strings.HasPrefix(queryUpper, "EXECUTE ") {
		return s.handleExecProcedure(conn, query)
	}

	// Default: Process the query using the query processor
	results, err := s.queryProcessor.ExecuteSQLBatch(query)
	if err != nil {
		log.Printf("Error processing query: %v", err)

		// Send error response
		errPacket := s.buildErrorPacket(err)
		writeErr := s.writePacket(conn, errPacket)
		if writeErr != nil {
			return fmt.Errorf("failed to send error packet: %w", writeErr)
		}

		return fmt.Errorf("query processing error: %w", err)
	}

	// Send result set
	resultPacket := s.buildResultPacket(results)
	err = s.writePacket(conn, resultPacket)
	if err != nil {
		return fmt.Errorf("failed to send result: %w", err)
	}

	log.Printf("Sent result packet with %d rows", len(results))
	return nil
}

func (s *Server) handleCreateProcedure(conn net.Conn, query string) error {
	log.Printf("Handling CREATE PROCEDURE: %s", query)

	// Parse CREATE PROCEDURE statement
	proc, err := procedure.ParseCreateProcedure(query)
	if err != nil {
		log.Printf("Error parsing CREATE PROCEDURE: %v", err)

		// Send error response
		errPacket := s.buildErrorPacket(err)
		writeErr := s.writePacket(conn, errPacket)
		if writeErr != nil {
			return fmt.Errorf("failed to send error packet: %w", writeErr)
		}

		return fmt.Errorf("CREATE PROCEDURE parsing error: %w", err)
	}

	// Store procedure in database
	err = s.procedureStorage.Create(proc)
	if err != nil {
		log.Printf("Error storing procedure: %v", err)

		// Send error response
		errPacket := s.buildErrorPacket(err)
		writeErr := s.writePacket(conn, errPacket)
		if writeErr != nil {
			return fmt.Errorf("failed to send error packet: %w", writeErr)
		}

		return fmt.Errorf("procedure storage error: %w", err)
	}

	// Send success response
	results := [][]string{
		{"Procedure created successfully"},
	}
	resultPacket := s.buildResultPacket(results)
	err = s.writePacket(conn, resultPacket)
	if err != nil {
		return fmt.Errorf("failed to send result: %w", err)
	}

	log.Printf("Created procedure: %s", proc.Name)
	return nil
}

func (s *Server) handleDropProcedure(conn net.Conn, query string) error {
	log.Printf("Handling DROP PROCEDURE: %s", query)

	// Extract procedure name
	// Simple parsing: DROP PROC[EDURE] procname
	parts := strings.Fields(query)
	if len(parts) < 3 {
		err := fmt.Errorf("invalid DROP PROCEDURE syntax")
		errPacket := s.buildErrorPacket(err)
		s.writePacket(conn, errPacket)
		return err
	}

	procName := parts[2]

	// Drop procedure from database
	err := s.procedureStorage.Drop(procName)
	if err != nil {
		log.Printf("Error dropping procedure: %v", err)

		// Send error response
		errPacket := s.buildErrorPacket(err)
		writeErr := s.writePacket(conn, errPacket)
		if writeErr != nil {
			return fmt.Errorf("failed to send error packet: %w", writeErr)
		}

		return fmt.Errorf("procedure drop error: %w", err)
	}

	// Send success response
	results := [][]string{
		{"Procedure dropped successfully"},
	}
	resultPacket := s.buildResultPacket(results)
	err = s.writePacket(conn, resultPacket)
	if err != nil {
		return fmt.Errorf("failed to send result: %w", err)
	}

	log.Printf("Dropped procedure: %s", procName)
	return nil
}

func (s *Server) handleExecProcedure(conn net.Conn, query string) error {
	log.Printf("Handling EXEC: %s", query)

	// Parse EXEC statement
	procName, paramValues, err := s.parseExecStatement(query)
	if err != nil {
		log.Printf("Error parsing EXEC statement: %v", err)

		// Send error response
		errPacket := s.buildErrorPacket(err)
		writeErr := s.writePacket(conn, errPacket)
		if writeErr != nil {
			return fmt.Errorf("failed to send error packet: %w", writeErr)
		}

		return fmt.Errorf("EXEC parsing error: %w", err)
	}

	// Execute procedure
	results, err := s.procedureExecutor.Execute(procName, paramValues)
	if err != nil {
		log.Printf("Error executing procedure: %v", err)

		// Send error response
		errPacket := s.buildErrorPacket(err)
		writeErr := s.writePacket(conn, errPacket)
		if writeErr != nil {
			return fmt.Errorf("failed to send error packet: %w", writeErr)
		}

		return fmt.Errorf("procedure execution error: %w", err)
	}

	// Send result set
	resultPacket := s.buildResultPacket(results)
	err = s.writePacket(conn, resultPacket)
	if err != nil {
		return fmt.Errorf("failed to send result: %w", err)
	}

	log.Printf("Executed procedure: %s, returned %d rows", procName, len(results))
	return nil
}

func (s *Server) parseExecStatement(query string) (string, map[string]interface{}, error) {
	// Simple parsing: EXEC procname @param1=value1, @param2=value2
	query = strings.TrimSpace(query)
	queryUpper := strings.ToUpper(query)

	// Remove EXEC or EXECUTE
	if strings.HasPrefix(queryUpper, "EXECUTE ") {
		query = query[7:]
	} else if strings.HasPrefix(queryUpper, "EXEC ") {
		query = query[4:]
	}

	query = strings.TrimSpace(query)

	// Split by whitespace to get procedure name
	parts := strings.Fields(query)
	if len(parts) == 0 {
		return "", nil, fmt.Errorf("empty EXEC statement")
	}

	procName := parts[0]
	paramValues := make(map[string]interface{})

	// Parse parameters if present
	if len(parts) > 1 {
		// Join remaining parts and split by comma
		paramStr := strings.Join(parts[1:], " ")
		params := strings.Split(paramStr, ",")

		for _, param := range params {
			param = strings.TrimSpace(param)
			if param == "" {
				continue
			}

			// Parse: @param=value
			kv := strings.Split(param, "=")
			if len(kv) != 2 {
				return "", nil, fmt.Errorf("invalid parameter syntax: %s", param)
			}

			paramName := strings.TrimSpace(kv[0])
			paramValue := strings.TrimSpace(kv[1])

			// Remove quotes if present
			if strings.HasPrefix(paramValue, "'") && strings.HasSuffix(paramValue, "'") {
				paramValue = paramValue[1 : len(paramValue)-1]
			}

			paramValues[paramName] = paramValue
		}
	}

	return procName, paramValues, nil
}

func (s *Server) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Server) buildResultPacket(rows [][]string) *tds.Packet {
	var buf []byte

	// Token type: COLMETADATA (0x81)
	buf = append(buf, 0x81)
	// Column count
	colCount := len(rows[0])
	buf = append(buf, byte(colCount))

	// Column metadata
	for i := 0; i < colCount; i++ {
		// User type (4 bytes)
		buf = append(buf, 0x00, 0x00, 0x00, 0x00)
		// Flags (2 bytes)
		buf = append(buf, 0x00, 0x00)
		// Type info (VARCHAR)
		buf = append(buf, 0xA7) // VARCHAR
		// Length
		buf = append(buf, 0xFF, 0xFF) // MAX length
		// Collation (5 bytes)
		buf = append(buf, 0x09, 0x04, 0xD0, 0x00, 0x34)
		// Table name (empty)
		buf = append(buf, 0x00)
		// Column name
		colName := fmt.Sprintf("Column%d", i+1)
		buf = append(buf, byte(len(colName)))
		buf = append(buf, []byte(colName)...)
	}

	// Token type: ROW (0xD1)
	for _, row := range rows {
		buf = append(buf, 0xD1)
		// Column values
		for _, val := range row {
			// Column length
			length := len(val)
			buf = append(buf, byte(length>>8), byte(length))
			// Column value
			buf = append(buf, []byte(val)...)
		}
	}

	// Token type: DONE (0xFD)
	buf = append(buf, 0xFD)
	// Status
	buf = append(buf, 0x00, 0x10) // FINAL
	// Current command
	buf = append(buf, 0x00, 0x00)
	// Row count
	buf = append(buf, 0x00, 0x00, 0x00, 0x00, byte(len(rows)))

	return tds.NewPacket(tds.PacketTypeTabular, tds.StatusEOM, 3, buf)
}

func (s *Server) handleRPC(conn net.Conn, packet *tds.Packet) error {
	log.Printf("Handling RPC request, data length: %d", len(packet.Data))

	// Parse RPC request
	rpcReq, err := tds.ParseRPCRequest(packet.Data)
	if err != nil {
		log.Printf("Error parsing RPC request: %v", err)

		// Send error response
		errPacket := s.buildErrorPacket(err)
		writeErr := s.writePacket(conn, errPacket)
		if writeErr != nil {
			return fmt.Errorf("failed to send error packet: %w", writeErr)
		}

		return fmt.Errorf("RPC parsing error: %w", err)
	}

	log.Printf("RPC Procedure: %s", rpcReq.ProcName)
	log.Printf("RPC Parameters: %d", len(rpcReq.Params))

	for i, param := range rpcReq.Params {
		log.Printf("  Param %d: %s = %v", i+1, param.Name, param.Value)
	}

	// Execute stored procedure
	results, err := s.storedProcedureHandler.Execute(rpcReq.ProcName, rpcReq.Params)
	if err != nil {
		log.Printf("Error executing stored procedure: %v", err)

		// Send error response
		errPacket := s.buildErrorPacket(err)
		writeErr := s.writePacket(conn, errPacket)
		if writeErr != nil {
			return fmt.Errorf("failed to send error packet: %w", writeErr)
		}

		return fmt.Errorf("stored procedure execution error: %w", err)
	}

	// Send RPC response
	rpcPacket := tds.BuildRPCResponse(results)
	err = s.writePacket(conn, rpcPacket)
	if err != nil {
		return fmt.Errorf("failed to send RPC response: %w", err)
	}

	log.Printf("Sent RPC response with %d rows", len(results))
	return nil
}

func main() {
	dbPath := "./data/tds_server.db"

	server, err := NewServer(defaultPort, dbPath)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	log.Printf("Starting TDS Server on port %d", defaultPort)
	err = server.Start()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
