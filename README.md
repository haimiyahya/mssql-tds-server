# MSSQL TDS Server

A Microsoft SQL Server-compatible server implementing the TDS (Tabular Data Stream) protocol with SQLite storage backend.

## Status
Proof of Concept - Phase 1, 2 & 3

## Overview
This project implements a minimal TDS server that can accept connections from standard Go mssql clients and handle basic request/response communication.

## Plan
See [PLAN.md](PLAN.md) for detailed project phases and implementation strategy.

## Project Structure
```
.
├── PLAN.md          # Detailed project plan
├── README.md        # This file
├── go.mod           # Go module definition
└── cmd/             # Server and client applications
    ├── server/      # TDS server implementation
    └── client/      # Test client using standard mssql driver
```

## Development
See [PLAN.md](PLAN.md) for implementation phases and tasks.

## Future Work
This project provides the foundation for a full-featured MSSQL-compatible server. Future phases will be implemented in a separate project building on this codebase.
