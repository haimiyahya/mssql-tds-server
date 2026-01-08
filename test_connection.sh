#!/bin/bash

# Simple test to send a pre-login packet to the server

# Pre-login packet (minimal)
# Header: Type=0x04 (PreLogin), Status=0x01 (EOM), Length=0x0011 (17 bytes), PacketID=1
# Data: Version token (0x00, offset=0x0005, length=0x0006), Terminator (0xFF)

# Build pre-login packet
echo -ne '\x04\x01\x00\x11\x00\x00\x00\x01\x00' > /tmp/prelogin.bin
echo -ne '\x00\x00\x05\x00\x06' >> /tmp/prelogin.bin
echo -ne '\x01\x00\x0b\x00\x01' >> /tmp/prelogin.bin  # Encryption token
echo -ne '\x02\x00\x0c\x00\x0b' >> /tmp/prelogin.bin  # Instance token
echo -ne '\xff' >> /tmp/prelogin.bin
echo -ne '\x09\x00\x00\x00\x00\x00' >> /tmp/prelogin.bin  # Version data
echo -ne '\x00' >> /tmp/prelogin.bin  # Encryption data
echo -ne '\x4d\x53\x53\x51\x4c\x53\x65\x72\x76\x65\x72\x00' >> /tmp/prelogin.bin  # Instance data "MSSQLServer"

echo "Sending pre-login packet to server..."
cat /tmp/prelogin.bin | nc localhost 1433
echo ""
echo "Response received"
