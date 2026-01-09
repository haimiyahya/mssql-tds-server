package sqlparser

import (
	"strings"
)

// parseCreateDatabase parses a CREATE DATABASE statement
func (p *Parser) parseCreateDatabase(query string) *Statement {
	// Format: CREATE DATABASE database_name

	upperQuery := strings.ToUpper(query)

	// Remove "CREATE DATABASE "
	query = strings.TrimSpace(strings.TrimPrefix(query, "CREATE DATABASE "))
	upperQuery = strings.TrimSpace(strings.TrimPrefix(upperQuery, "CREATE DATABASE "))

	// Extract database name
	databaseName := strings.TrimSpace(query)

	// Validate database name
	if databaseName == "" {
		return &Statement{
			Type:    StatementTypeCreateDatabase,
			RawQuery: query,
		}
	}

	return &Statement{
		Type: StatementTypeCreateDatabase,
		CreateDatabase: &CreateDatabaseStatement{
			DatabaseName: databaseName,
		},
		RawQuery: query,
	}
}

// parseDropDatabase parses a DROP DATABASE statement
func (p *Parser) parseDropDatabase(query string) *Statement {
	// Format: DROP DATABASE database_name

	upperQuery := strings.ToUpper(query)

	// Remove "DROP DATABASE "
	query = strings.TrimSpace(strings.TrimPrefix(query, "DROP DATABASE "))
	upperQuery = strings.TrimSpace(strings.TrimPrefix(upperQuery, "DROP DATABASE "))

	// Extract database name
	databaseName := strings.TrimSpace(query)

	// Validate database name
	if databaseName == "" {
		return &Statement{
			Type:    StatementTypeDropDatabase,
			RawQuery: query,
		}
	}

	return &Statement{
		Type: StatementTypeDropDatabase,
		DropDatabase: &DropDatabaseStatement{
			DatabaseName: databaseName,
		},
		RawQuery: query,
	}
}

// parseUseDatabase parses a USE statement
func (p *Parser) parseUseDatabase(query string) *Statement {
	// Format: USE database_name

	upperQuery := strings.ToUpper(query)

	// Remove "USE "
	query = strings.TrimSpace(strings.TrimPrefix(query, "USE "))
	upperQuery = strings.TrimSpace(strings.TrimPrefix(upperQuery, "USE "))

	// Extract database name
	databaseName := strings.TrimSpace(query)

	// Validate database name
	if databaseName == "" {
		return &Statement{
			Type:    StatementTypeUseDatabase,
			RawQuery: query,
		}
	}

	return &Statement{
		Type: StatementTypeUseDatabase,
		UseDatabase: &UseDatabaseStatement{
			DatabaseName: databaseName,
		},
		RawQuery: query,
	}
}
