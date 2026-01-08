package tds

import (
	"encoding/binary"
	"fmt"
)

// RPCParameter represents a parameter in an RPC call
type RPCParameter struct {
	Name      string
	Type      byte
	MaxLength uint16
	Value     interface{}
}

// RPCRequest represents an RPC (Remote Procedure Call) request
type RPCRequest struct {
	ProcName       string
	ProcID         uint16
	Options        uint16
	Params         []*RPCParameter
	TransactionID   uint16
}

// ParseRPCRequest parses an RPC packet from TDS data
func ParseRPCRequest(data []byte) (*RPCRequest, error) {
	req := &RPCRequest{}

	buf := NewBuffer(data)

	// Read total length (4 bytes)
	_, err := buf.ReadUint16()
	if err != nil {
		return nil, fmt.Errorf("error reading RPC total length: %w", err)
	}
	// Skip next 2 bytes (reserved)
	buf.ReadBytes(2)

	// Read parameter count (4 bytes)
	paramCount, err := buf.ReadUint16()
	if err != nil {
		return nil, fmt.Errorf("error reading parameter count: %w", err)
	}
	// Skip next 2 bytes (reserved)
	buf.ReadBytes(2)

	// Read option flags (4 bytes)
	options, err := buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("error reading RPC options: %w", err)
	}
	req.Options = uint16(options)

	// Read parameter data
	req.Params = make([]*RPCParameter, 0, paramCount)

	for i := uint16(0); i < paramCount; i++ {
		param, err := parseRPCParameter(buf)
		if err != nil {
			return nil, fmt.Errorf("error parsing parameter %d: %w", i, err)
		}
		req.Params = append(req.Params, param)
	}

	// Extract procedure name from first parameter (if it's the procedure name parameter)
	if len(req.Params) > 0 {
		// In TDS RPC, the first parameter often contains the procedure name
		if name, ok := req.Params[0].Value.(string); ok {
			req.ProcName = name
		}
	}

	// Try to read procedure ID from data stream
	if len(data) >= 8 {
		req.ProcID = binary.BigEndian.Uint16(data[0:2])
	}

	return req, nil
}

// parseRPCParameter parses a single RPC parameter
func parseRPCParameter(buf *Buffer) (*RPCParameter, error) {
	param := &RPCParameter{}

	// Read name length and name
	nameLen, err := buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("error reading parameter name length: %w", err)
	}

	if nameLen > 0 {
		nameBytes, err := buf.ReadBytes(int(nameLen))
		if err != nil {
			return nil, fmt.Errorf("error reading parameter name: %w", err)
		}
		param.Name = string(nameBytes)
	}

	// Read status flags
	status, err := buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("error reading parameter status: %w", err)
	}
	_ = status // Status not used for now

	// Read user type (4 bytes)
	userType, err := buf.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("error reading user type: %w", err)
	}
	_ = userType // Not used for now

	// Read data type (1 byte)
	dataType, err := buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("error reading data type: %w", err)
	}
	param.Type = dataType

	// Parse value based on data type
	value, err := parseRPCValue(buf, dataType)
	if err != nil {
		return nil, fmt.Errorf("error parsing parameter value: %w", err)
	}
	param.Value = value

	return param, nil
}

// parseRPCValue parses a parameter value based on data type
func parseRPCValue(buf *Buffer, dataType byte) (interface{}, error) {
	switch dataType {
	case 0x38: // VARCHAR
		length, err := buf.ReadUint16()
		if err != nil {
			return nil, fmt.Errorf("error reading varchar length: %w", err)
		}
		if length == 0xFFFF {
			// NULL value
			return nil, nil
		}
		valueBytes, err := buf.ReadBytes(int(length))
		if err != nil {
			return nil, fmt.Errorf("error reading varchar value: %w", err)
		}
		return string(valueBytes), nil

	case 0x68: // BIGINT
		valueBytes, err := buf.ReadBytes(8)
		if err != nil {
			return nil, fmt.Errorf("error reading bigint value: %w", err)
		}
		return int64(binary.BigEndian.Uint64(valueBytes)), nil

	case 0x30: // INT
		valueBytes, err := buf.ReadBytes(4)
		if err != nil {
			return nil, fmt.Errorf("error reading int value: %w", err)
		}
		return int32(binary.BigEndian.Uint32(valueBytes)), nil

	case 0x34: // SMALLINT
		valueBytes, err := buf.ReadBytes(2)
		if err != nil {
			return nil, fmt.Errorf("error reading smallint value: %w", err)
		}
		return int16(binary.BigEndian.Uint16(valueBytes)), nil

	default:
		return nil, fmt.Errorf("unsupported data type: %#02x", dataType)
	}
}

// BuildRPCResponse builds an RPC response packet with a result set
func BuildRPCResponse(rows [][]string) *Packet {
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

	return NewPacket(PacketTypeTabular, StatusEOM, 3, buf)
}
