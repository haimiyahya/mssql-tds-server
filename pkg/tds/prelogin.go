package tds

import "fmt"

// PreLoginToken represents TDS pre-login token types
type PreLoginToken byte

const (
	TokenVersion      PreLoginToken = 0x00
	TokenEncryption   PreLoginToken = 0x01
	TokenInstance     PreLoginToken = 0x02
	TokenThreadID     PreLoginToken = 0x03
	TokenMARS         PreLoginToken = 0x04
	TokenTraceID      PreLoginToken = 0x05
	TokenFedAuth      PreLoginToken = 0x06
	TokenNonce        PreLoginToken = 0x07
	TokenTerminator   PreLoginToken = 0xFF
)

// Encryption level constants
const (
	EncryptionOff      = 0x00 // No encryption (cleartext)
	EncryptionOn       = 0x01 // SSL/TLS encryption if supported
	EncryptionRequired = 0x02 // Encryption required (reject if not supported)
)

// PreLoginOption represents a pre-login option
type PreLoginOption struct {
	TokenType PreLoginToken
	Offset    uint16
	Length    uint16
}

// PreLoginRequest represents a TDS pre-login request
type PreLoginRequest struct {
	Version    []byte
	Encryption byte
	Instance   []byte
	ThreadID   []byte
	MARS       byte
}

// PreLoginResponse represents a TDS pre-login response
type PreLoginResponse struct {
	Version    []byte
	Encryption byte
	Instance   []byte
	ThreadID   []byte
	MARS       byte
}

// ParsePreLoginRequest parses a pre-login request packet
func ParsePreLoginRequest(data []byte) (*PreLoginRequest, error) {
	req := &PreLoginRequest{}

	buf := NewBuffer(data)

	// Parse options
	var options []*PreLoginOption
	for i := 0; i < 10; i++ { // Limit iterations to prevent infinite loop
		tokenType, err := buf.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("error reading token type at iteration %d: %w", i, err)
		}

		if PreLoginToken(tokenType) == TokenTerminator {
			break
		}

		offset, err := buf.ReadUint16()
		if err != nil {
			return nil, fmt.Errorf("error reading offset at iteration %d: %w", i, err)
		}

		length, err := buf.ReadUint16()
		if err != nil {
			return nil, fmt.Errorf("error reading length at iteration %d: %w", i, err)
		}

		options = append(options, &PreLoginOption{
			TokenType: PreLoginToken(tokenType),
			Offset:    offset,
			Length:    length,
		})
	}

	// Read option data
	for _, opt := range options {
		if opt.Offset+opt.Length > uint16(len(data)) {
			continue
		}
		optData := data[opt.Offset : opt.Offset+opt.Length]

		switch opt.TokenType {
		case TokenVersion:
			req.Version = optData
		case TokenEncryption:
			if len(optData) > 0 {
				req.Encryption = optData[0]
			}
		case TokenInstance:
			req.Instance = optData
		case TokenThreadID:
			req.ThreadID = optData
		case TokenMARS:
			if len(optData) > 0 {
				req.MARS = optData[0]
			}
		}
	}

	return req, nil
}

// SerializePreLoginResponse serializes a pre-login response
func SerializePreLoginResponse(resp *PreLoginResponse) []byte {

	// Build data section
	dataOffset := uint16(0)

	// Version (6 bytes)
	versionData := []byte{0x09, 0x00, 0x00, 0x00, 0x00, 0x00} // TDS 7.3
	if resp.Version != nil {
		versionData = resp.Version
	}
	versionOffset := dataOffset
	dataOffset += uint16(len(versionData))

	// Encryption (1 byte)
	encryptionData := []byte{resp.Encryption}
	encryptionOffset := dataOffset
	dataOffset += uint16(len(encryptionData))

	// Instance (variable)
	instanceData := resp.Instance
	if len(instanceData) == 0 {
		instanceData = []byte("MSSQLServer")
	}
	instanceOffset := dataOffset
	dataOffset += uint16(len(instanceData))

	// Thread ID (4 bytes)
	threadIDData := make([]byte, 4)
	// Use thread ID from request or default
	if len(resp.ThreadID) == 4 {
		threadIDData = resp.ThreadID
	}
	threadIDOffset := dataOffset
	dataOffset += uint16(len(threadIDData))

	// MARS (1 byte)
	marsData := []byte{resp.MARS}
	marsOffset := dataOffset
	dataOffset += uint16(len(marsData))

	// Calculate header size
	headerSize := 0
	headerSize += 5 // Version option
	headerSize += 5 // Encryption option
	headerSize += 5 // Instance option
	headerSize += 5 // Thread ID option
	headerSize += 5 // MARS option
	headerSize += 1 // Terminator

	totalSize := headerSize + int(dataOffset)
	result := make([]byte, totalSize)
	offset := 0

	// Write header
	result[offset] = byte(TokenVersion)
	offset++
	result[offset] = byte(versionOffset >> 8)
	offset++
	result[offset] = byte(versionOffset)
	offset++
	result[offset] = byte(len(versionData) >> 8)
	offset++
	result[offset] = byte(len(versionData))
	offset++

	result[offset] = byte(TokenEncryption)
	offset++
	result[offset] = byte(encryptionOffset >> 8)
	offset++
	result[offset] = byte(encryptionOffset)
	offset++
	result[offset] = byte(len(encryptionData) >> 8)
	offset++
	result[offset] = byte(len(encryptionData))
	offset++

	result[offset] = byte(TokenInstance)
	offset++
	result[offset] = byte(instanceOffset >> 8)
	offset++
	result[offset] = byte(instanceOffset)
	offset++
	result[offset] = byte(len(instanceData) >> 8)
	offset++
	result[offset] = byte(len(instanceData))
	offset++

	result[offset] = byte(TokenThreadID)
	offset++
	result[offset] = byte(threadIDOffset >> 8)
	offset++
	result[offset] = byte(threadIDOffset)
	offset++
	result[offset] = byte(len(threadIDData) >> 8)
	offset++
	result[offset] = byte(len(threadIDData))
	offset++

	result[offset] = byte(TokenMARS)
	offset++
	result[offset] = byte(marsOffset >> 8)
	offset++
	result[offset] = byte(marsOffset)
	offset++
	result[offset] = byte(len(marsData) >> 8)
	offset++
	result[offset] = byte(len(marsData))
	offset++

	result[offset] = byte(TokenTerminator)
	offset++

	// Write data
	copy(result[offset:], versionData)
	offset += len(versionData)
	copy(result[offset:], encryptionData)
	offset += len(encryptionData)
	copy(result[offset:], instanceData)
	offset += len(instanceData)
	copy(result[offset:], threadIDData)
	offset += len(threadIDData)
	copy(result[offset:], marsData)

	return result
}

// DefaultPreLoginResponse creates a default pre-login response
func DefaultPreLoginResponse(encryption byte) *PreLoginResponse {
	return &PreLoginResponse{
		Version:    []byte{0x09, 0x00, 0x00, 0x00, 0x00, 0x00},
		Encryption: encryption, // Encryption level (0x00=OFF, 0x01=ON, 0x02=REQUIRED)
		Instance:   []byte("MSSQLServer"),
		ThreadID:   []byte{0x00, 0x00, 0x00, 0x00},
		MARS:       0x00, // No MARS
	}
}
