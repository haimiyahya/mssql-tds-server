package tds

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// PacketType represents TDS packet types
type PacketType byte

const (
	PacketTypePreLogin  PacketType = 0x04
	PacketTypeLogin     PacketType = 0x10
	PacketTypeSQLBatch  PacketType = 0x01
	PacketTypeRPC       PacketType = 0x03
	PacketTypeTabular   PacketType = 0x04
	PacketTypeAttention PacketType = 0x06
	PacketTypeBulkLoad  PacketType = 0x07
	PacketTypeTransMgr  PacketType = 0x0E
	PacketTypeManager   PacketType = 0x0F
)

// PacketStatus represents TDS packet status flags
type PacketStatus byte

const (
	StatusEOM      PacketStatus = 0x01 // End of message
	StatusIgnore   PacketStatus = 0x02
	StatusReset    PacketStatus = 0x08
	StatusResetExp PacketStatus = 0x10
)

// PacketHeader represents the 8-byte TDS packet header
type PacketHeader struct {
	Type     PacketType
	Status   PacketStatus
	Length   uint16
	SPID     uint16
	PacketID uint8
	Window   uint8
}

// ParseHeader parses a TDS packet header from a byte slice
func ParseHeader(data []byte) (*PacketHeader, error) {
	if len(data) < 8 {
		return nil, errors.New("invalid packet header: too short")
	}

	header := &PacketHeader{
		Type:     PacketType(data[0]),
		Status:   PacketStatus(data[1]),
		Length:   binary.BigEndian.Uint16(data[2:4]),
		SPID:     binary.BigEndian.Uint16(data[4:6]),
		PacketID: data[6],
		Window:   data[7],
	}

	return header, nil
}

// Serialize serializes a packet header to bytes
func (h *PacketHeader) Serialize() []byte {
	buf := make([]byte, 8)
	buf[0] = byte(h.Type)
	buf[1] = byte(h.Status)
	binary.BigEndian.PutUint16(buf[2:4], h.Length)
	binary.BigEndian.PutUint16(buf[4:6], h.SPID)
	buf[6] = h.PacketID
	buf[7] = h.Window
	return buf
}

// Packet represents a complete TDS packet
type Packet struct {
	Header *PacketHeader
	Data   []byte
}

// ParsePacket parses a complete TDS packet from a byte slice
func ParsePacket(data []byte) (*Packet, error) {
	header, err := ParseHeader(data)
	if err != nil {
		return nil, err
	}

	if len(data) < int(header.Length) {
		return nil, fmt.Errorf("incomplete packet: expected %d bytes, got %d", header.Length, len(data))
	}

	return &Packet{
		Header: header,
		Data:   data[8:header.Length],
	}, nil
}

// Serialize serializes a complete TDS packet to bytes
func (p *Packet) Serialize() []byte {
	p.Header.Length = uint16(8 + len(p.Data))
	buf := make([]byte, p.Header.Length)
	copy(buf, p.Header.Serialize())
	copy(buf[8:], p.Data)
	return buf
}

// NewPacket creates a new TDS packet
func NewPacket(packetType PacketType, status PacketStatus, packetID uint8, data []byte) *Packet {
	return &Packet{
		Header: &PacketHeader{
			Type:     packetType,
			Status:   status,
			PacketID: packetID,
		},
		Data: data,
	}
}

// Buffer provides methods for reading/writing TDS data
type Buffer struct {
	buf *bytes.Buffer
}

// NewBuffer creates a new TDS buffer
func NewBuffer(data []byte) *Buffer {
	return &Buffer{
		buf: bytes.NewBuffer(data),
	}
}

// ReadByte reads a single byte
func (b *Buffer) ReadByte() (byte, error) {
	return b.buf.ReadByte()
}

// ReadBytes reads n bytes
func (b *Buffer) ReadBytes(n int) ([]byte, error) {
	data := make([]byte, n)
	_, err := b.buf.Read(data)
	return data, err
}

// ReadUint16 reads a big-endian uint16
func (b *Buffer) ReadUint16() (uint16, error) {
	data, err := b.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(data), nil
}

// ReadUint32 reads a big-endian uint32
func (b *Buffer) ReadUint32() (uint32, error) {
	data, err := b.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(data), nil
}

// WriteByte writes a single byte
func (b *Buffer) WriteByte(v byte) {
	b.buf.WriteByte(v)
}

// WriteBytes writes bytes
func (b *Buffer) WriteBytes(data []byte) {
	b.buf.Write(data)
}

// WriteUint16 writes a big-endian uint16
func (b *Buffer) WriteUint16(v uint16) {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, v)
	b.WriteBytes(data)
}

// WriteUint32 writes a big-endian uint32
func (b *Buffer) WriteUint32(v uint32) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, v)
	b.WriteBytes(data)
}

// Bytes returns the buffer contents
func (b *Buffer) Bytes() []byte {
	return b.buf.Bytes()
}

// Len returns the buffer length
func (b *Buffer) Len() int {
	return b.buf.Len()
}
