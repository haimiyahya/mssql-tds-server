# Database Object Storage Locations

## Overview

This document describes where database objects are stored in MSSQL TDS Server, including table names, structures, and file locations.

---

## Master Database Tables

### Storage Location: `./data/master.db`

The master database stores all system-level catalog tables that track databases, logins, stored procedures, and functions.

---

### 1. `syslogins` - User Login Storage

**Purpose**: Stores user login credentials and authentication information

**Table Structure**:
```sql
CREATE TABLE syslogins (
    sid INTEGER PRIMARY KEY AUTOINCREMENT,        -- Login ID (Security ID)
    name TEXT NOT NULL UNIQUE,                   -- Login username
    password_hash TEXT NOT NULL,                  -- Password hash (bcrypt)
    type TEXT DEFAULT 'SQL_SERVER',              -- Authentication type
    default_database_name TEXT DEFAULT 'master',    -- Default database
    default_language TEXT DEFAULT 'english',        -- Default language
    created_date DATETIME DEFAULT CURRENT_TIMESTAMP, -- Creation date
    modified_date DATETIME DEFAULT CURRENT_TIMESTAMP, -- Last modified date
    is_disabled BOOLEAN DEFAULT 0,               -- Disabled flag
    is_locked BOOLEAN DEFAULT 0,                 -- Locked flag
    login_count INTEGER DEFAULT 0,                -- Login count
    last_login_date DATETIME,                    -- Last login date
    description TEXT                             -- Description
);
```

**Access**: `pkg/auth/auth.go` - `AuthManager` struct

**Stored Information**:
- ✅ Login username
- ✅ Password hash (bcrypt)
- ✅ Authentication type (SQL_SERVER, WINDOWS, MIXED)
- ✅ Default database
- ✅ Default language
- ✅ Account status (disabled, locked)
- ✅ Login statistics (count, last login)

---

### 2. `sys_databases` - Database Catalog

**Purpose**: Catalog of all databases on the server

**Table Structure**:
```sql
CREATE TABLE sys_databases (
    database_id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    state TEXT DEFAULT 'ONLINE',
    create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    file_path TEXT,
    is_system BOOLEAN DEFAULT 0
);
```

**Access**: `pkg/database/catalog.go` - `Catalog` struct

**Stored Information**:
- ✅ Database name
- ✅ Database ID
- ✅ Database state (ONLINE, OFFLINE)
- ✅ Creation date
- ✅ File path (`./data/DatabaseName.db`)
- ✅ System database flag (master, tempdb, model, msdb)

**Example Entries**:
```
database_id | name        | state   | create_date        | file_path                    | is_system
-------------|-------------|---------|--------------------|------------------------------|-----------
1            | master      | ONLINE   | 2023-01-01 00:00:00 | /root/projects/.../data/master.db | 1
2            | tempdb      | ONLINE   | 2023-01-01 00:00:00 | /root/projects/.../data/tempdb.db | 1
3            | model       | ONLINE   | 2023-01-01 00:00:00 | /root/projects/.../data/model.db | 1
4            | msdb        | ONLINE   | 2023-01-01 00:00:00 | /root/projects/.../data/msdb.db | 1
5            | MyDatabase  | ONLINE   | 2024-01-15 10:30:00 | /root/projects/.../data/MyDatabase.db | 0
```

---

### 3. `sys_procedures` - Stored Procedure Catalog

**Purpose**: Catalog of all stored procedures across databases

**Table Structure**:
```sql
CREATE TABLE sys_procedures (
    procedure_id INTEGER PRIMARY KEY,
    database_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    definition TEXT NOT NULL,
    create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (database_id) REFERENCES sys_databases(database_id)
);
```

**Access**: `pkg/database/catalog.go` - `Catalog` struct

**Stored Information**:
- ✅ Procedure name
- ✅ Procedure ID
- ✅ Database ID (which database owns this procedure)
- ✅ Procedure definition (SQL body)
- ✅ Creation date

**Behavior**: 
- Procedures are **database-scoped** (stored in database they're created in)
- Catalog tracks which database owns each procedure
- Use `USE database_name` before creating procedure

**Example Entries**:
```
procedure_id | database_id | name              | definition               | create_date
-------------|-------------|-------------------|-------------------------|--------------------
1            | 6           | sp_GetUsers       | SELECT * FROM Users     | 2024-01-15 10:30:00
2            | 7           | sp_GetProducts    | SELECT * FROM Products  | 2024-01-15 10:35:00
```

---

### 4. `sys_functions` - Function Catalog

**Purpose**: Catalog of all functions across databases

**Table Structure**:
```sql
CREATE TABLE sys_functions (
    function_id INTEGER PRIMARY KEY,
    database_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    definition TEXT NOT NULL,
    return_type TEXT,
    create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (database_id) REFERENCES sys_databases(database_id)
);
```

**Access**: `pkg/database/catalog.go` - `Catalog` struct

**Stored Information**:
- ✅ Function name
- ✅ Function ID
- ✅ Database ID (which database owns this function)
- ✅ Function definition (SQL body)
- ✅ Return type (data type)
- ✅ Creation date

**Behavior**:
- Functions are **database-scoped** (stored in database they're created in)
- Catalog tracks which database owns each function
- Use `USE database_name` before creating function

**Example Entries**:
```
function_id | database_id | name              | return_type | definition                     | create_date
------------|-------------|-------------------|-------------|--------------------------------|--------------------
1           | 6           | fn_GetUserName     | VARCHAR(100) | SELECT Name FROM Users...  | 2024-01-15 10:40:00
2           | 7           | fn_GetTotalPrice   | DECIMAL      | SUM(price) OVER(...)          | 2024-01-15 10:45:00
```

---

## Per-Database System Tables

### Storage Location: `./data/DatabaseName.db`

Each database (including master.db) contains system tables that track objects and columns within that database.

---

### 5. `sys_objects` - Object Catalog (Per Database)

**Purpose**: Catalog of objects (tables, views, etc.) in a database

**Table Structure**:
```sql
CREATE TABLE sys_objects (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    create_date DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**Access**: Created automatically when database is created

**Stored Information**:
- ✅ Object name
- ✅ Object ID
- ✅ Object type (TABLE, VIEW, INDEX, etc.)
- ✅ Creation date

**Behavior**:
- Exists in **each database**
- Created automatically when database is created
- Tracks all objects in that database

**Example Entries**:
```
id | name        | type   | create_date
----|-------------|---------|--------------------
1   | Users       | TABLE   | 2024-01-15 10:30:00
2   | Products    | TABLE   | 2024-01-15 10:35:00
3   | Orders      | TABLE   | 2024-01-15 10:40:00
4   | UserView    | VIEW    | 2024-01-15 10:45:00
```

---

### 6. `sys_columns` - Column Catalog (Per Database)

**Purpose**: Catalog of columns for each object in a database

**Table Structure**:
```sql
CREATE TABLE sys_columns (
    id INTEGER PRIMARY KEY,
    object_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    is_nullable BOOLEAN DEFAULT 1,
    FOREIGN KEY (object_id) REFERENCES sys_objects(id)
);
```

**Access**: Created automatically when database is created

**Stored Information**:
- ✅ Column name
- ✅ Column ID
- ✅ Object ID (which object owns this column)
- ✅ Column type (INT, VARCHAR, etc.)
- ✅ Nullable flag
- ✅ Foreign key to sys_objects

**Behavior**:
- Exists in **each database**
- Created automatically when database is created
- Tracks all columns for all objects in that database

**Example Entries**:
```
id | object_id | name        | type            | is_nullable
----|------------|-------------|----------------|-------------
1   | 1          | id          | INT             | 0
2   | 1          | name        | VARCHAR(100)    | 0
3   | 1          | email       | VARCHAR(255)    | 0
4   | 2          | id          | INT             | 0
5   | 2          | name        | VARCHAR(200)    | 0
```

---

## File System Structure

### Data Directory: `./data/`

**Database Files**:
```
./data/
  ├── master.db        # Master database (contains catalog tables)
  ├── tempdb.db       # Temporary database
  ├── model.db        # Template database
  ├── msdb.db         # Management database
  └── UserDB1.db      # User databases
```

**Catalog Tables** (in `master.db`):
- `syslogins` - User login storage
- `sys_databases` - Database catalog
- `sys_procedures` - Stored procedure catalog
- `sys_functions` - Function catalog

**System Tables** (in each database):
- `sys_objects` - Object catalog
- `sys_columns` - Column catalog

---

## Object Storage Summary

| Object Type | Storage Location | Table Name | Database | File |
|--------------|------------------|--------------|---------|
| **User Logins** | Master DB | `syslogins` | `./data/master.db` |
| **Database Catalog** | Master DB | `sys_databases` | `./data/master.db` |
| **Stored Procedures** | Master DB | `sys_procedures` | `./data/master.db` |
| **Functions** | Master DB | `sys_functions` | `./data/master.db` |
| **Objects** (Tables, Views) | Per Database | `sys_objects` | `./data/DatabaseName.db` |
| **Columns** | Per Database | `sys_columns` | `./data/DatabaseName.db` |

---

## Access and Management

### Programmatic Access

**Master Database Catalog**:
```go
// Create catalog
catalog := database.NewCatalog("./data", masterDB)

// Access database catalog
catalog.CreateDatabase("MyDatabase")
catalog.DropDatabase("MyDatabase")
catalog.ListDatabases()

// Access procedure catalog
catalog.CreateProcedure("MyDatabase", "sp_GetUsers", "SELECT * FROM Users")

// Access function catalog
catalog.CreateFunction("MyDatabase", "fn_GetUserName", "SELECT Name FROM Users...", "VARCHAR(100)")
```

**User Authentication**:
```go
// Create auth manager
authManager := auth.NewAuthManager(masterDB)

// Access login catalog
authManager.CreateLogin("username", "password", auth.AuthTypeSQLServer)
authManager.AuthenticateLogin("username", "password")
authManager.ListLogins()
```

### SQL Access

**Master Database** (via `USE master`):
```sql
-- View all databases
SELECT * FROM master.sys_databases;

-- View all logins
SELECT name, database_id, is_disabled FROM master.syslogins;

-- View all procedures
SELECT name, database_id FROM master.sys_procedures;

-- View all functions
SELECT name, database_id, return_type FROM master.sys_functions;
```

**Per-Database** (via `USE MyDatabase`):
```sql
-- View all objects
SELECT name, type, create_date FROM sys_objects;

-- View all columns
SELECT name, type, is_nullable FROM sys_columns;
```

---

## Storage Behavior

### Database-Scoped Storage

**Procedures and Functions**:
- Stored in **specific database** they're created in
- Catalog in `master.sys_procedures` / `master.sys_functions` tracks ownership
- Use `USE database_name` before creating procedure/function

**Objects and Columns**:
- Stored in **database** they're created in
- Each database has its own `sys_objects` and `sys_columns` tables
- Objects exist only in their database

### Cross-Database Queries

**Catalog Tables**: In `master.db`
- Track databases, procedures, functions across all databases
- Enable cross-database queries
- Support database attachment

**System Tables**: In each database
- Track objects and columns within that database
- Enable object management within database

---

## Security Considerations

### Access Control

**syslogins Table**:
- Passwords hashed with bcrypt
- Disabled/locked flags prevent access
- Login statistics for auditing

**Catalog Tables** (sys_databases, sys_procedures, sys_functions):
- In master database
- Require database access to read
- Foreign key constraints maintain integrity

**System Tables** (sys_objects, sys_columns):
- In each database
- Require database access to read
- Foreign key constraints maintain integrity

---

## References

- [README.md](README.md) - Project README
- [PHASE41_PROGRESS.md](PHASE41_PROGRESS.md) - Database Management Documentation
- [PHASE42_PROGRESS.md](PHASE42_PROGRESS.md) - Authentication System Documentation
- [AUTHENTICATION.md](AUTHENTICATION.md) - Authentication and Password Security
- [pkg/database/catalog.go](pkg/database/catalog.go) - Database Catalog Implementation
- [pkg/auth/auth.go](pkg/auth/auth.go) - Authentication Implementation

---

## Summary

### Storage Locations:

| Object Type | Table Name | Database | File Path |
|--------------|-------------|------------|------------|
| **User Logins** | `syslogins` | Master | `./data/master.db` |
| **Database Catalog** | `sys_databases` | Master | `./data/master.db` |
| **Stored Procedures** | `sys_procedures` | Master | `./data/master.db` |
| **Functions** | `sys_functions` | Master | `./data/master.db` |
| **Objects** | `sys_objects` | Per Database | `./data/DatabaseName.db` |
| **Columns** | `sys_columns` | Per Database | `./data/DatabaseName.db` |

### Key Points:

1. **Master Database**: Contains all catalog tables
2. **syslogins**: Stores user logins and credentials
3. **sys_databases**: Catalog of all databases
4. **sys_procedures**: Catalog of all stored procedures
5. **sys_functions**: Catalog of all functions
6. **sys_objects**: Catalog of objects in each database
7. **sys_columns**: Catalog of columns in each database
8. **Database-Scoped**: Procedures, functions, objects stored in their database
9. **Catalog Tables**: Track objects across all databases
10. **File Structure**: Each database = SQLite file in `./data/`

---

*Database Object Storage - Fully Documented.*
