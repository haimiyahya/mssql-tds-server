#!/bin/bash

# Start server in background and capture output
export PATH=/usr/local/go/bin:$PATH
cd /root/projects/mssql-tds-server

# Start server in background
./bin/server > /tmp/server.log 2>&1 &
SERVER_PID=$!

echo "Server started with PID: $SERVER_PID"
echo "Waiting for server to initialize..."
sleep 2

# Check if server is running
if ps -p $SERVER_PID > /dev/null; then
    echo "Server is running"
else
    echo "Server failed to start"
    cat /tmp/server.log
    exit 1
fi

# Test connection with a simple TCP client
echo "Testing connection..."

# Build a proper pre-login packet
# Header: Type=0x04, Status=0x01, Length=0x001A (26 bytes), SPID=0, PacketID=1, Window=0
# Data section (18 bytes):
#   Token 0x00 (Version) at offset 0x0000, length 0x0006
#   Token 0x01 (Encryption) at offset 0x0006, length 0x0001
#   Terminator 0xFF
#   Version data: 0x09,0x00,0x00,0x00,0x00,0x00 (TDS 7.3)
#   Encryption data: 0x00 (No encryption)
# Total: 8 (header) + 18 (data) = 26 bytes

printf '\x04\x01\x00\x1a\x00\x00\x00\x01\x00\x00\x00\x00\x00\x06\x01\x00\x06\x00\x01\xff\x09\x00\x00\x00\x00\x00\x00' | od -An -tx1 | head -5

printf '\x04\x01\x00\x1a\x00\x00\x00\x01\x00\x00\x00\x00\x00\x06\x01\x00\x06\x00\x01\xff\x09\x00\x00\x00\x00\x00\x00' | nc localhost 1433 &
NC_PID=$!

sleep 2

# Check server logs
echo "Server logs:"
cat /tmp/server.log

# Cleanup
kill $SERVER_PID $NC_PID 2>/dev/null
wait $SERVER_PID $NC_PID 2>/dev/null

echo "Test completed"
