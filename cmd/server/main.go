package main

import (
	"fmt"
	"log"
	"net"

	"github.com/factory/mssql-tds-server/pkg/tds"
)

const (
	defaultPort = 1433
)

type Server struct {
	addr string
}

func NewServer(port int) *Server {
	return &Server{
		addr: fmt.Sprintf(":%d", port),
	}
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

	// Read first packet (pre-login)
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

		// Handle login request
		if packet.Header.Type == tds.PacketTypeLogin {
			err = s.handleLogin(conn, packet)
			if err != nil {
				log.Printf("Error handling login: %v", err)
				break
			}
		}

		// Handle SQL batch
		if packet.Header.Type == tds.PacketTypeSQLBatch {
			err = s.handleSQLBatch(conn, packet)
			if err != nil {
				log.Printf("Error handling SQL batch: %v", err)
				break
			}
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
	log.Printf("Packet data length: %d", len(packet.Data))
	log.Printf("Packet data (hex): %#x", packet.Data)

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

	// For now, just acknowledge the login
	// In a full implementation, we would parse the login packet and authenticate

	// Send login acknowledgment
	// For Phase 1, we'll just send a simple success response
	loginAck := s.buildLoginAckPacket()

	err := s.writePacket(conn, loginAck)
	if err != nil {
		return fmt.Errorf("failed to send login ack: %w", err)
	}

	log.Println("Sent login acknowledgment")
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

func (s *Server) handleSQLBatch(conn net.Conn, packet *tds.Packet) error {
	log.Printf("Handling SQL batch: %s", string(packet.Data))

	// For Phase 2, we'll implement simple query processing
	// For Phase 1, we just log it

	// Send a simple result set
	resultPacket := s.buildResultPacket([][]string{{"HELLO", "WORLD"}})

	err := s.writePacket(conn, resultPacket)
	if err != nil {
		return fmt.Errorf("failed to send result: %w", err)
	}

	log.Println("Sent result packet")
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

func main() {
	server := NewServer(defaultPort)

	log.Printf("Starting TDS Server on port %d", defaultPort)
	err := server.Start()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
