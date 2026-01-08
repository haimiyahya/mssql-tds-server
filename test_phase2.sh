#!/bin/bash

# Test Phase 2: Basic request/response communication

export PATH=/usr/local/go/bin:$PATH
cd /root/projects/mssql-tds-server

echo "=== Phase 2 Test: Basic Request/Response ==="
echo ""

# Start server in background
echo "Starting server..."
./bin/server > /tmp/server_phase2.log 2>&1 &
SERVER_PID=$!

echo "Server started with PID: $SERVER_PID"
echo "Waiting for server to initialize..."
sleep 2

# Check if server is running
if ps -p $SERVER_PID > /dev/null; then
    echo "✓ Server is running"
else
    echo "✗ Server failed to start"
    cat /tmp/server_phase2.log
    exit 1
fi

echo ""
echo "=== Testing Phase 2 functionality ==="
echo ""
echo "The server should now be able to:"
echo "1. Accept connections from mssql client"
echo "2. Process SQL batch queries"
echo "3. Return results in uppercase (echo functionality)"
echo ""
echo "Server is ready for client connections."
echo "Server logs (first 20 lines):"
head -20 /tmp/server_phase2.log

echo ""
echo "=== Test Setup Complete ==="
echo "Server is running on port 1433"
echo "You can test with the client using:"
echo "  ./bin/client"
echo ""
echo "Press Ctrl+C to stop the server"

# Wait for user interrupt
trap "kill $SERVER_PID; echo ''; echo 'Server stopped'; exit 0" INT TERM
wait $SERVER_PID
