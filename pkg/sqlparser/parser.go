package sqlparser

import (
	"regexp"
	"strings"
)

// Parser handles SQL statement parsing
type Parser struct{}

// NewParser creates a new SQL parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses a SQL query and returns a Statement
func (p *Parser) Parse(query string) (*Statement, error) {
	// Trim and normalize query
	query = strings.TrimSpace(query)

	// Remove trailing semicolons
	query = strings.TrimSuffix(query, ";")

	// Detect statement type
	upperQuery := strings.ToUpper(query)

	var stmt *Statement

	if strings.HasPrefix(upperQuery, "SELECT ") {
		stmt = p.parseSelect(query)
	} else if strings.HasPrefix(upperQuery, "INSERT INTO ") {
		stmt = p.parseInsert(query)
	} else if strings.HasPrefix(upperQuery, "UPDATE ") {
		stmt = p.parseUpdate(query)
	} else if strings.HasPrefix(upperQuery, "DELETE FROM ") {
		stmt = p.parseDelete(query)
	} else if strings.HasPrefix(upperQuery, "CREATE TABLE ") {
		stmt = p.parseCreateTable(query)
	} else if strings.HasPrefix(upperQuery, "DROP TABLE ") {
		stmt = p.parseDropTable(query)
	} else {
		// Unknown statement type
		stmt = &Statement{
			Type:    StatementTypeUnknown,
			RawQuery: query,
		}
	}

	return stmt, nil
}

// parseSelect parses a SELECT statement
func (p *Parser) parseSelect(query string) *Statement {
	// Extract columns and table name
	// Format: SELECT [DISTINCT] [columns] FROM [table] [JOIN ...] [WHERE clause] [GROUP BY clause] [HAVING clause] [ORDER BY clause]

	upperQuery := strings.ToUpper(query)

	// Check for DISTINCT keyword
	distinct := false
	if strings.HasPrefix(upperQuery, "SELECT DISTINCT ") {
		distinct = true
		query = strings.TrimPrefix(query, "SELECT DISTINCT ")
		query = "SELECT " + query
		upperQuery = strings.ToUpper(query)
	}

	// Find FROM keyword
	fromIndex := strings.Index(upperQuery, " FROM ")
	if fromIndex == -1 {
		return &Statement{
			Type:    StatementTypeSelect,
			RawQuery: query,
		}
	}

	// Extract columns (between SELECT and FROM)
	columnsPart := query[7:fromIndex] // 7 = len("SELECT ")
	columns := p.parseColumns(columnsPart)

	// Check for aggregate functions in columns
	aggregates := p.parseAggregates(columnsPart)
	isAggregateQuery := len(aggregates) > 0

	// Check for subqueries
	hasSubqueries := p.detectSubqueries(query)

	// Find ORDER BY clause (optional) - parse this first as it comes last
	orderByIndex := strings.Index(upperQuery, " ORDER BY ")
	var orderByClause []OrderByClause
	var havingClause string
	var groupByClause []GroupByClause
	var whereClause string
	var joins []JoinClause
	var tableName string

	// Parse ORDER BY first (comes after HAVING or GROUP BY or WHERE or JOINs or FROM)
	if orderByIndex != -1 {
		orderByClause = p.parseOrderBy(query[orderByIndex+9:]) // 9 = len(" ORDER BY ")
		query = strings.TrimSpace(query[:orderByIndex])
		upperQuery = strings.ToUpper(query)
	}

	// Find HAVING clause (optional) - comes after GROUP BY
	havingIndex := strings.Index(upperQuery, " HAVING ")
	if havingIndex != -1 {
		havingClause = strings.TrimSpace(query[havingIndex+7:]) // 7 = len(" HAVING ")
		query = strings.TrimSpace(query[:havingIndex])
		upperQuery = strings.ToUpper(query)
	}

	// Find GROUP BY clause (optional) - comes after WHERE or JOINs
	groupByIndex := strings.Index(upperQuery, " GROUP BY ")
	if groupByIndex != -1 {
		groupByClause = p.parseGroupBy(query[groupByIndex+9:]) // 9 = len(" GROUP BY ")
		query = strings.TrimSpace(query[:groupByIndex])
		upperQuery = strings.ToUpper(query)
	}

	// Find WHERE clause (optional) - comes after JOINs
	whereIndex := strings.Index(upperQuery, " WHERE ")
	if whereIndex != -1 {
		whereClause = strings.TrimSpace(query[whereIndex+7:]) // 7 = len(" WHERE ")
		query = strings.TrimSpace(query[:whereIndex])
		upperQuery = strings.ToUpper(query)
	}

	// Parse JOIN clauses (optional) - comes after FROM
	joins = p.parseJoins(query, upperQuery, fromIndex)
	// Remove JOINs from query to extract table name
	// We need to recalculate query after removing JOINs
	// For now, we'll extract table name from before the first JOIN

	// Extract table name (from FROM to first JOIN or WHERE or end)
	tableEndIndex := len(query)
	// Find first JOIN, GROUP BY, HAVING, ORDER BY, or WHERE
	firstClauseIndex := -1

	joinKeywords := []string{" INNER JOIN ", " LEFT JOIN ", " RIGHT JOIN ", " FULL JOIN ", " JOIN "}
	for _, keyword := range joinKeywords {
		if idx := strings.Index(upperQuery, keyword); idx != -1 && (firstClauseIndex == -1 || idx < firstClauseIndex) {
			firstClauseIndex = idx
		}
	}

	if whereIndex != -1 && whereIndex < tableEndIndex {
		tableEndIndex = whereIndex
	}
	if groupByIndex != -1 && groupByIndex < tableEndIndex {
		tableEndIndex = groupByIndex
	}
	if havingIndex != -1 && havingIndex < tableEndIndex {
		tableEndIndex = havingIndex
	}
	if orderByIndex != -1 && orderByIndex < tableEndIndex {
		tableEndIndex = orderByIndex
	}
	if firstClauseIndex != -1 && firstClauseIndex < tableEndIndex {
		tableEndIndex = firstClauseIndex
	}

	// Extract table name
	tablePart := query[fromIndex+6 : tableEndIndex] // 6 = len(" FROM ")
	tableName = strings.TrimSpace(tablePart)

	// Clean up table name (remove trailing stuff)
	tableName = p.extractTableName(tableName)

	return &Statement{
		Type: StatementTypeSelect,
		Select: &SelectStatement{
			Columns:          columns,
			Table:            tableName,
			Joins:            joins,
			WhereClause:       whereClause,
			Distinct:         distinct,
			OrderBy:           orderByClause,
			Aggregates:       aggregates,
			IsAggregateQuery:  isAggregateQuery,
			GroupBy:          groupByClause,
			HavingClause:     havingClause,
			HasSubqueries:    hasSubqueries,
		},
		RawQuery: query,
	}
}

// parseOrderBy parses an ORDER BY clause
func (p *Parser) parseOrderBy(orderByClause string) []OrderByClause {
	// Remove trailing semicolon if present
	orderByClause = strings.TrimSpace(strings.TrimSuffix(orderByClause, ";"))

	// Split by comma
	columnParts := strings.Split(orderByClause, ",")
	orderByClauses := make([]OrderByClause, 0, len(columnParts))

	for _, part := range columnParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Check for ASC or DESC direction
		upperPart := strings.ToUpper(part)
		direction := "ASC" // Default direction

		if strings.HasSuffix(upperPart, " ASC") {
			part = strings.TrimSpace(part[:len(part)-4])
		} else if strings.HasSuffix(upperPart, " DESC") {
			part = strings.TrimSpace(part[:len(part)-5])
			direction = "DESC"
		}

		orderByClauses = append(orderByClauses, OrderByClause{
			Column:    part,
			Direction: direction,
		})
	}

	return orderByClauses
}

// parseGroupBy parses a GROUP BY clause
func (p *Parser) parseGroupBy(groupByClause string) []GroupByClause {
	// Remove trailing semicolon if present
	groupByClause = strings.TrimSpace(strings.TrimSuffix(groupByClause, ";"))

	// Split by comma
	columnParts := strings.Split(groupByClause, ",")
	groupByClauses := make([]GroupByClause, 0, len(columnParts))

	for _, part := range columnParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		groupByClauses = append(groupByClauses, GroupByClause{
			Column: part,
		})
	}

	return groupByClauses
}

// parseJoins parses JOIN clauses from query
func (p *Parser) parseJoins(query string, upperQuery string, fromIndex int) []JoinClause {
	// Find all JOINs in the query
	joins := make([]JoinClause, 0)
	currentIndex := fromIndex

	for currentIndex < len(query) {
		// Find next JOIN keyword
		joinKeywords := []struct {
			keyword string
			joinType string
		}{
			{" INNER JOIN ", "INNER"},
			{" LEFT JOIN ", "LEFT"},
			{" RIGHT JOIN ", "RIGHT"},
			{" FULL JOIN ", "FULL"},
			{" JOIN ", "INNER"}, // Default JOIN type is INNER
		}

		bestIndex := -1
		bestType := ""
		bestKeyword := ""

		for _, jk := range joinKeywords {
			idx := strings.Index(upperQuery[currentIndex:], jk.keyword)
			if idx != -1 && (bestIndex == -1 || idx < bestIndex) {
				bestIndex = idx + currentIndex
				bestType = jk.joinType
				bestKeyword = jk.keyword
			}
		}

		if bestIndex == -1 {
			// No more JOINs found
			break
		}

		// Extract JOIN part (from JOIN keyword to next clause or end)
		joinStart := bestIndex + len(bestKeyword)
		joinPart := query[joinStart:]

		// Find ON keyword
		onIndex := strings.Index(strings.ToUpper(joinPart), " ON ")
		if onIndex == -1 {
			// Invalid JOIN syntax, skip
			currentIndex = bestIndex + len(bestKeyword)
			continue
		}

		// Extract table name (from JOIN to ON)
		tablePart := strings.TrimSpace(joinPart[:onIndex])

		// Check for AS alias
		var alias string
		aliasIndex := strings.Index(strings.ToUpper(tablePart), " AS ")
		if aliasIndex != -1 {
			alias = strings.TrimSpace(tablePart[aliasIndex+4:])
			tablePart = strings.TrimSpace(tablePart[:aliasIndex])
		}

		// Extract ON clause (from ON to next clause or end)
		onClausePart := joinPart[onIndex+4:] // 4 = len(" ON ")

		// Find where this ON clause ends (next JOIN, WHERE, GROUP BY, HAVING, ORDER BY)
		onEndKeywords := []string{" INNER JOIN ", " LEFT JOIN ", " RIGHT JOIN ", " FULL JOIN ", " JOIN ", " WHERE ", " GROUP BY ", " HAVING ", " ORDER BY "}
		onEndIndex := -1

		for _, keyword := range onEndKeywords {
			idx := strings.Index(strings.ToUpper(onClausePart), keyword)
			if idx != -1 && (onEndIndex == -1 || idx < onEndIndex) {
				onEndIndex = idx
			}
		}

		if onEndIndex != -1 {
			onClause := strings.TrimSpace(onClausePart[:onEndIndex])
			joins = append(joins, JoinClause{
				Type:     bestType,
				Table:    tablePart,
				OnClause: onClause,
				Alias:    alias,
			})
			currentIndex = joinStart + onEndIndex
		} else {
			// ON clause goes to end of query
			onClause := strings.TrimSpace(onClausePart)
			joins = append(joins, JoinClause{
				Type:     bestType,
				Table:    tablePart,
				OnClause: onClause,
				Alias:    alias,
			})
			break
		}
	}

	return joins
}

// parseAggregates parses aggregate functions from column list
func (p *Parser) parseAggregates(columnsPart string) []AggregateFunction {
	// Aggregate function patterns: COUNT(*), COUNT(column), SUM(column), AVG(column), MIN(column), MAX(column)
	re := regexp.MustCompile(`(COUNT|SUM|AVG|MIN|MAX)\(\*?[a-zA-Z0-9_]*\)`)

	matches := re.FindAllString(columnsPart, -1)
	aggregates := make([]AggregateFunction, 0, len(matches))

	for _, match := range matches {
		// Extract function type
		upperMatch := strings.ToUpper(match)
		var funcType string
		var column string

		switch {
		case strings.HasPrefix(upperMatch, "COUNT("):
			funcType = "COUNT"
			column = strings.TrimSuffix(strings.TrimPrefix(match, "COUNT("), ")")
		case strings.HasPrefix(upperMatch, "SUM("):
			funcType = "SUM"
			column = strings.TrimSuffix(strings.TrimPrefix(match, "SUM("), ")")
		case strings.HasPrefix(upperMatch, "AVG("):
			funcType = "AVG"
			column = strings.TrimSuffix(strings.TrimPrefix(match, "AVG("), ")")
		case strings.HasPrefix(upperMatch, "MIN("):
			funcType = "MIN"
			column = strings.TrimSuffix(strings.TrimPrefix(match, "MIN("), ")")
		case strings.HasPrefix(upperMatch, "MAX("):
			funcType = "MAX"
			column = strings.TrimSuffix(strings.TrimPrefix(match, "MAX("), ")")
		}

		// Check for AS alias (e.g., COUNT(*) AS total)
		alias := ""
		aliasPattern := regexp.MustCompile(strings.ToUpper(match) + `\s+AS\s+(\w+)`)
		aliasMatches := aliasPattern.FindStringSubmatch(columnsPart)
		if len(aliasMatches) > 1 {
			alias = aliasMatches[1]
		}

		aggregates = append(aggregates, AggregateFunction{
			Type:   funcType,
			Column: column,
			Alias:  alias,
		})
	}

	return aggregates
}

// detectSubqueries checks if a query contains subqueries
func (p *Parser) detectSubqueries(query string) bool {
	upperQuery := strings.ToUpper(query)

	// Common subquery patterns
	subqueryPatterns := []string{
		" IN (SELECT ",
		" NOT IN (SELECT ",
		" EXISTS (SELECT ",
		" NOT EXISTS (SELECT ",
		" = (SELECT ",
		" != (SELECT ",
		" <> (SELECT ",
		" > (SELECT ",
		" < (SELECT ",
		" >= (SELECT ",
		" <= (SELECT ",
	}

	for _, pattern := range subqueryPatterns {
		if strings.Contains(upperQuery, pattern) {
			return true
		}
	}

	return false
}

// parseInsert parses an INSERT statement
func (p *Parser) parseInsert(query string) *Statement {
	// Format: INSERT INTO table (columns) VALUES (values)

	upperQuery := strings.ToUpper(query)

	// Remove "INSERT INTO "
	query = strings.TrimPrefix(query, "INSERT INTO ")
	upperQuery = strings.TrimPrefix(upperQuery, "INSERT INTO ")

	// Extract table name
	spaceIndex := strings.Index(query, " ")
	if spaceIndex == -1 {
		return &Statement{
			Type:    StatementTypeInsert,
			RawQuery: query,
		}
	}

	tableName := strings.TrimSpace(query[:spaceIndex])

	// Check for column list
	var columns []string
	valuesStart := spaceIndex + 1

	if query[spaceIndex+1] == '(' {
		// Find closing parenthesis
		closingParen := strings.Index(query, ")")
		if closingParen != -1 {
			columnsPart := query[spaceIndex+2 : closingParen]
			columns = p.parseColumns(columnsPart)
			valuesStart = closingParen + 1
		}
	}

	// Find VALUES keyword
	valuesIndex := strings.Index(strings.ToUpper(query[valuesStart:]), " VALUES ")
	if valuesIndex == -1 {
		return &Statement{
			Type:    StatementTypeInsert,
			RawQuery: query,
		}
	}

	actualValuesStart := valuesStart + valuesIndex + 8 // 8 = len(" VALUES ")

	// Parse values (handle single or multiple row inserts)
	values, _ := p.parseValues(query[actualValuesStart:])

	return &Statement{
		Type: StatementTypeInsert,
		Insert: &InsertStatement{
			Table:  tableName,
			Columns: columns,
			Values: values,
		},
		RawQuery: query,
	}
}

// parseUpdate parses an UPDATE statement
func (p *Parser) parseUpdate(query string) *Statement {
	// Format: UPDATE table SET clause [WHERE clause]

	upperQuery := strings.ToUpper(query)

	// Remove "UPDATE "
	query = strings.TrimPrefix(query, "UPDATE ")
	upperQuery = strings.TrimPrefix(upperQuery, "UPDATE ")

	// Extract table name
	spaceIndex := strings.Index(query, " ")
	if spaceIndex == -1 {
		return &Statement{
			Type:    StatementTypeUpdate,
			RawQuery: query,
		}
	}

	tableName := strings.TrimSpace(query[:spaceIndex])

	// Find SET keyword
	setIndex := strings.Index(upperQuery, " SET ")
	if setIndex == -1 {
		return &Statement{
			Type:    StatementTypeUpdate,
			RawQuery: query,
		}
	}

	// Find WHERE clause (optional)
	whereIndex := strings.Index(upperQuery, " WHERE ")
	var setClause string
	var whereClause string

	if whereIndex != -1 {
		// Extract SET clause
		setClause = strings.TrimSpace(query[setIndex+5 : whereIndex])

		// Extract WHERE clause
		whereClause = strings.TrimSpace(query[whereIndex+7:])
	} else {
		// No WHERE clause
		setClause = strings.TrimSpace(query[setIndex+5:])
	}

	return &Statement{
		Type: StatementTypeUpdate,
		Update: &UpdateStatement{
			Table:       tableName,
			SetClause:   setClause,
			WhereClause: whereClause,
		},
		RawQuery: query,
	}
}

// parseDelete parses a DELETE statement
func (p *Parser) parseDelete(query string) *Statement {
	// Format: DELETE FROM table [WHERE clause]

	upperQuery := strings.ToUpper(query)

	// Remove "DELETE FROM "
	query = strings.TrimPrefix(query, "DELETE FROM ")
	upperQuery = strings.TrimPrefix(upperQuery, "DELETE FROM ")

	// Find WHERE clause (optional)
	whereIndex := strings.Index(upperQuery, " WHERE ")
	var tableName string
	var whereClause string

	if whereIndex != -1 {
		// Extract table name
		tableName = strings.TrimSpace(query[:whereIndex])

		// Extract WHERE clause
		whereClause = strings.TrimSpace(query[whereIndex+7:])
	} else {
		// No WHERE clause
		tableName = strings.TrimSpace(query)
	}

	// Clean up table name
	tableName = p.extractTableName(tableName)

	return &Statement{
		Type: StatementTypeDelete,
		Delete: &DeleteStatement{
			Table:       tableName,
			WhereClause: whereClause,
		},
		RawQuery: query,
	}
}

// parseCreateTable parses a CREATE TABLE statement
func (p *Parser) parseCreateTable(query string) *Statement {
	// Format: CREATE TABLE table (column1 type, column2 type, ...)

	upperQuery := strings.ToUpper(query)

	// Remove "CREATE TABLE "
	query = strings.TrimPrefix(query, "CREATE TABLE ")
	upperQuery = strings.TrimPrefix(upperQuery, "CREATE TABLE ")

	// Extract table name
	parenIndex := strings.Index(query, "(")
	if parenIndex == -1 {
		return &Statement{
			Type:    StatementTypeCreateTable,
			RawQuery: query,
		}
	}

	tableName := strings.TrimSpace(query[:parenIndex])

	// Extract column definitions
	closingParen := strings.LastIndex(query, ")")
	if closingParen == -1 {
		return &Statement{
			Type:    StatementTypeCreateTable,
			RawQuery: query,
		}
	}

	columnsPart := query[parenIndex+1 : closingParen]
	columns := p.parseColumnDefinitions(columnsPart)

	return &Statement{
		Type: StatementTypeCreateTable,
		CreateTable: &CreateTableStatement{
			TableName: tableName,
			Columns:   columns,
		},
		RawQuery: query,
	}
}

// parseDropTable parses a DROP TABLE statement
func (p *Parser) parseDropTable(query string) *Statement {
	// Format: DROP TABLE table

	upperQuery := strings.ToUpper(query)

	// Remove "DROP TABLE "
	query = strings.TrimPrefix(query, "DROP TABLE ")
	upperQuery = strings.TrimPrefix(upperQuery, "DROP TABLE ")

	// Extract table name
	tableName := strings.TrimSpace(query)

	return &Statement{
		Type: StatementTypeDropTable,
		DropTable: &DropTableStatement{
			TableName: tableName,
		},
		RawQuery: query,
	}
}

// parseColumns parses a comma-separated list of columns
func (p *Parser) parseColumns(columnsStr string) []string {
	columnsStr = strings.TrimSpace(columnsStr)

	// Check if it's "SELECT *"
	if columnsStr == "*" {
		return []string{"*"}
	}

	// Split by comma
	columns := strings.Split(columnsStr, ",")
	result := make([]string, 0, len(columns))

	for _, col := range columns {
		col = strings.TrimSpace(col)
		if col != "" {
			// Handle aliases (e.g., "name AS username")
			if asIndex := strings.Index(strings.ToUpper(col), " AS "); asIndex != -1 {
				col = strings.TrimSpace(col[:asIndex])
			}
			result = append(result, col)
		}
	}

	return result
}

// parseValues parses VALUES clause
func (p *Parser) parseValues(valuesStr string) ([][]interface{}, error) {
	valuesStr = strings.TrimSpace(valuesStr)

	// Handle single row: VALUES (val1, val2, ...)
	if !strings.HasPrefix(valuesStr, "(") {
		return nil, nil
	}

	var values [][]interface{}

	// Simple parsing: extract values from parentheses
	// This is a basic implementation - for production, use a proper SQL parser

	// Remove outer parentheses if present
	if strings.HasPrefix(valuesStr, "(") && strings.HasSuffix(valuesStr, ")") {
		// Single row
		valuesStr = valuesStr[1 : len(valuesStr)-1]
	}

	// Split by comma (basic - doesn't handle strings with commas)
	valueStrs := strings.Split(valuesStr, ",")
	rowValues := make([]interface{}, 0, len(valueStrs))

	for _, val := range valueStrs {
		val = strings.TrimSpace(val)

		// Handle strings
		if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
			val = val[1 : len(val)-1]
		}

		rowValues = append(rowValues, val)
	}

	values = append(values, rowValues)

	return values, nil
}

// parseColumnDefinitions parses column definitions from CREATE TABLE
func (p *Parser) parseColumnDefinitions(columnsStr string) []ColumnDefinition {
	columnsStr = strings.TrimSpace(columnsStr)

	// Split by comma
	columnDefs := strings.Split(columnsStr, ",")
	result := make([]ColumnDefinition, 0, len(columnDefs))

	for _, colDef := range columnDefs {
		colDef = strings.TrimSpace(colDef)
		if colDef == "" {
			continue
		}

		// Split by first space to get column name and type
		parts := strings.Fields(colDef)
		if len(parts) >= 2 {
			columnName := parts[0]
			columnType := strings.Join(parts[1:], " ")

			// Remove constraints (PRIMARY KEY, NOT NULL, etc.)
			for _, constraint := range []string{"PRIMARY KEY", "NOT NULL", "UNIQUE", "NULL", "DEFAULT"} {
				if idx := strings.Index(strings.ToUpper(columnType), constraint); idx != -1 {
					columnType = strings.TrimSpace(columnType[:idx])
				}
			}

			result = append(result, ColumnDefinition{
				Name: columnName,
				Type: columnType,
			})
		}
	}

	return result
}

// extractTableName extracts clean table name from SQL
func (p *Parser) extractTableName(tableName string) string {
	// Remove schema prefix if present (e.g., dbo.users -> users)
	if dotIndex := strings.Index(tableName, "."); dotIndex != -1 {
		tableName = tableName[dotIndex+1:]
	}

	// Remove square brackets if present (e.g., [users] -> users)
	tableName = strings.TrimPrefix(tableName, "[")
	tableName = strings.TrimSuffix(tableName, "]")

	// Remove backticks if present (e.g., `users` -> users)
	tableName = strings.TrimPrefix(tableName, "`")
	tableName = strings.TrimSuffix(tableName, "`")

	// Remove double quotes if present (e.g., "users" -> users)
	tableName = strings.TrimPrefix(tableName, "\"")
	tableName = strings.TrimSuffix(tableName, "\"")

	tableName = strings.TrimSpace(tableName)

	return tableName
}

// ParseStatementType detects the statement type from a query
func ParseStatementType(query string) StatementType {
	query = strings.TrimSpace(query)
	query = strings.TrimSuffix(query, ";")
	upperQuery := strings.ToUpper(query)

	if strings.HasPrefix(upperQuery, "SELECT ") {
		return StatementTypeSelect
	} else if strings.HasPrefix(upperQuery, "INSERT INTO ") {
		return StatementTypeInsert
	} else if strings.HasPrefix(upperQuery, "UPDATE ") {
		return StatementTypeUpdate
	} else if strings.HasPrefix(upperQuery, "DELETE FROM ") {
		return StatementTypeDelete
	} else if strings.HasPrefix(upperQuery, "CREATE TABLE ") {
		return StatementTypeCreateTable
	} else if strings.HasPrefix(upperQuery, "DROP TABLE ") {
		return StatementTypeDropTable
	}

	return StatementTypeUnknown
}

// IsStatementTypeSupported checks if a statement type is supported
func IsStatementTypeSupported(stmtType StatementType) bool {
	return stmtType >= StatementTypeSelect && stmtType <= StatementTypeDropTable
}

// ExtractTableName extracts table name from various statement types
func ExtractTableName(query string) string {
	parser := NewParser()
	stmt, _ := parser.Parse(query)

	switch stmt.Type {
	case StatementTypeSelect:
		return stmt.Select.Table
	case StatementTypeInsert:
		return stmt.Insert.Table
	case StatementTypeUpdate:
		return stmt.Update.Table
	case StatementTypeDelete:
		return stmt.Delete.Table
	case StatementTypeCreateTable:
		return stmt.CreateTable.TableName
	case StatementTypeDropTable:
		return stmt.DropTable.TableName
	default:
		return ""
	}
}

// HasWhereClause checks if a query has a WHERE clause
func HasWhereClause(query string) bool {
	upperQuery := strings.ToUpper(query)
	return strings.Contains(upperQuery, " WHERE ")
}

// ExtractWhereClause extracts WHERE clause from query
func ExtractWhereClause(query string) string {
	upperQuery := strings.ToUpper(query)
	whereIndex := strings.Index(upperQuery, " WHERE ")
	if whereIndex == -1 {
		return ""
	}

	whereClause := strings.TrimSpace(query[whereIndex+7:])

	// Remove trailing ORDER BY, GROUP BY, etc.
	for _, terminator := range []string{" ORDER BY ", " GROUP BY ", " LIMIT "} {
		if idx := strings.Index(strings.ToUpper(whereClause), terminator); idx != -1 {
			whereClause = strings.TrimSpace(whereClause[:idx])
		}
	}

	return whereClause
}

// StripComments removes SQL comments from a query
func StripComments(query string) string {
	// Remove single-line comments (--)
	re := regexp.MustCompile(`--[^\n]*\n`)
	query = re.ReplaceAllString(query, " ")

	// Remove multi-line comments (/* */)
	re = regexp.MustCompile(`/\*.*?\*/`)
	query = re.ReplaceAllString(query, " ")

	return strings.TrimSpace(query)
}
