package sqlparser

import (
	"database/sql"
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
	} else if strings.HasPrefix(upperQuery, "ALTER TABLE ") {
		stmt = p.parseAlterTable(query)
	} else if strings.HasPrefix(upperQuery, "CREATE VIEW ") {
		stmt = p.parseCreateView(query)
	} else if strings.HasPrefix(upperQuery, "DROP VIEW ") {
		stmt = p.parseDropView(query)
	} else if strings.HasPrefix(upperQuery, "CREATE INDEX ") || strings.HasPrefix(upperQuery, "CREATE UNIQUE INDEX ") {
		stmt = p.parseCreateIndex(query)
	} else if strings.HasPrefix(upperQuery, "DROP INDEX ") {
		stmt = p.parseDropIndex(query)
	} else if strings.HasPrefix(upperQuery, "PREPARE ") {
		stmt = p.parsePrepare(query)
	} else if strings.HasPrefix(upperQuery, "EXECUTE ") {
		stmt = p.executeStatement(query)
	} else if strings.HasPrefix(upperQuery, "DEALLOCATE PREPARE ") {
		stmt = p.parseDeallocatePrepare(query)
	} else if strings.HasPrefix(upperQuery, "BEGIN TRANSACTION") || strings.HasPrefix(upperQuery, "BEGIN") || strings.HasPrefix(upperQuery, "START TRANSACTION") {
		stmt = p.parseBeginTransaction(query)
	} else if strings.HasPrefix(upperQuery, "COMMIT") || strings.HasPrefix(upperQuery, "COMMIT TRAN") {
		stmt = p.parseCommit(query)
	} else if strings.HasPrefix(upperQuery, "ROLLBACK") || strings.HasPrefix(upperQuery, "ROLLBACK TRAN") {
		stmt = p.parseRollback(query)
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
	// Format: CREATE TABLE table (column1 type constraints, column2 type constraints, ...)

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

	// Parse table-level constraints (will be passed through in RawQuery)
	constraints := make([]TableConstraint, 0)

	return &Statement{
		Type: StatementTypeCreateTable,
		CreateTable: &CreateTableStatement{
			TableName:   tableName,
			Columns:     columns,
			Constraints: constraints,
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

// parseAlterTable parses an ALTER TABLE statement
func (p *Parser) parseAlterTable(query string) *Statement {
	// Format: ALTER TABLE table_name ADD COLUMN column_name column_type
	// Format: ALTER TABLE table_name RENAME TO new_table_name
	// Format: ALTER TABLE table_name RENAME COLUMN old_name TO new_name
	// Note: SQLite doesn't support DROP COLUMN or ALTER COLUMN natively

	upperQuery := strings.ToUpper(query)

	// Remove "ALTER TABLE "
	query = strings.TrimPrefix(query, "ALTER TABLE ")
	upperQuery = strings.TrimPrefix(upperQuery, "ALTER TABLE ")

	// Find first space (separates table name from action)
	spaceIndex := strings.Index(upperQuery, " ")
	if spaceIndex == -1 {
		return &Statement{
			Type:    StatementTypeAlterTable,
			RawQuery: query,
		}
	}

	// Extract table name
	tableName := strings.TrimSpace(query[:spaceIndex])

	// Extract action part
	actionPart := strings.TrimSpace(query[spaceIndex:])
	upperActionPart := strings.ToUpper(actionPart)

	// Determine action type
	alterTableStmt := &AlterTableStatement{
		TableName: tableName,
	}

	if strings.HasPrefix(upperActionPart, "ADD ") || strings.HasPrefix(upperActionPart, "ADD COLUMN ") {
		// ALTER TABLE table_name ADD COLUMN column_name column_type
		alterTableStmt.Action = "ADD"
		actionPart = strings.TrimPrefix(actionPart, "ADD COLUMN ")
		actionPart = strings.TrimPrefix(actionPart, "ADD ")
		upperActionPart = strings.ToUpper(actionPart)

		// Extract column name and type
		spaceIndex = strings.Index(upperActionPart, " ")
		if spaceIndex == -1 {
			return &Statement{
				Type:    StatementTypeAlterTable,
				RawQuery: query,
			}
		}

		alterTableStmt.Column = strings.TrimSpace(actionPart[:spaceIndex])
		alterTableStmt.Type = strings.TrimSpace(actionPart[spaceIndex:])

	} else if strings.HasPrefix(upperActionPart, "RENAME TO ") {
		// ALTER TABLE table_name RENAME TO new_table_name
		alterTableStmt.Action = "RENAME TO"
		actionPart = strings.TrimPrefix(actionPart, "RENAME TO ")
		alterTableStmt.NewName = strings.TrimSpace(actionPart)

	} else if strings.HasPrefix(upperActionPart, "RENAME COLUMN ") {
		// ALTER TABLE table_name RENAME COLUMN old_name TO new_name
		alterTableStmt.Action = "RENAME COLUMN"
		actionPart = strings.TrimPrefix(actionPart, "RENAME COLUMN ")
		upperActionPart = strings.ToUpper(actionPart)

		// Find TO keyword
		toIndex := strings.Index(upperActionPart, " TO ")
		if toIndex == -1 {
			return &Statement{
				Type:    StatementTypeAlterTable,
				RawQuery: query,
			}
		}

		alterTableStmt.Column = strings.TrimSpace(actionPart[:toIndex])
		alterTableStmt.NewName = strings.TrimSpace(actionPart[toIndex+4:]) // 4 = len(" TO ")
	} else {
		// Unsupported action
		return &Statement{
			Type:    StatementTypeAlterTable,
			RawQuery: query,
		}
	}

	return &Statement{
		Type: StatementTypeAlterTable,
		AlterTable: alterTableStmt,
		RawQuery: query,
	}
}

// parseCreateView parses a CREATE VIEW statement
func (p *Parser) parseCreateView(query string) *Statement {
	// Format: CREATE VIEW view_name AS SELECT ...
	// Format: CREATE VIEW view_name AS (SELECT ...)

	upperQuery := strings.ToUpper(query)

	// Remove "CREATE VIEW "
	query = strings.TrimPrefix(query, "CREATE VIEW ")
	upperQuery = strings.TrimPrefix(upperQuery, "CREATE VIEW ")

	// Find AS keyword
	asIndex := strings.Index(upperQuery, " AS ")
	if asIndex == -1 {
		return &Statement{
			Type:    StatementTypeCreateView,
			RawQuery: query,
		}
	}

	// Extract view name
	viewName := strings.TrimSpace(query[:asIndex])

	// Extract SELECT query (everything after AS)
	selectQuery := strings.TrimSpace(query[asIndex+4:]) // 4 = len(" AS ")

	// Remove parentheses if present (SQL Server style)
	selectQuery = strings.TrimPrefix(selectQuery, "(")
	selectQuery = strings.TrimSuffix(selectQuery, ")")

	return &Statement{
		Type: StatementTypeCreateView,
		CreateView: &CreateViewStatement{
			ViewName:  viewName,
			SelectQuery: selectQuery,
		},
		RawQuery: query,
	}
}

// parseDropView parses a DROP VIEW statement
func (p *Parser) parseDropView(query string) *Statement {
	// Format: DROP VIEW view_name

	upperQuery := strings.ToUpper(query)

	// Remove "DROP VIEW "
	query = strings.TrimPrefix(query, "DROP VIEW ")
	upperQuery = strings.TrimPrefix(upperQuery, "DROP VIEW ")

	// Extract view name
	viewName := strings.TrimSpace(query)

	return &Statement{
		Type: StatementTypeDropView,
		DropView: &DropViewStatement{
			ViewName: viewName,
		},
		RawQuery: query,
	}
}

// parseCreateIndex parses a CREATE INDEX statement
func (p *Parser) parseCreateIndex(query string) *Statement {
	// Format: CREATE INDEX index_name ON table_name (col1, col2, ...)
	// Format: CREATE UNIQUE INDEX index_name ON table_name (col1, col2, ...)

	upperQuery := strings.ToUpper(query)

	// Check for UNIQUE keyword
	unique := false
	if strings.HasPrefix(upperQuery, "CREATE UNIQUE INDEX ") {
		unique = true
		query = strings.TrimPrefix(query, "CREATE UNIQUE INDEX ")
		upperQuery = strings.TrimPrefix(upperQuery, "CREATE UNIQUE INDEX ")
	} else {
		query = strings.TrimPrefix(query, "CREATE INDEX ")
		upperQuery = strings.TrimPrefix(upperQuery, "CREATE INDEX ")
	}

	// Find ON keyword
	onIndex := strings.Index(upperQuery, " ON ")
	if onIndex == -1 {
		return &Statement{
			Type:    StatementTypeCreateIndex,
			RawQuery: query,
		}
	}

	// Extract index name
	indexName := strings.TrimSpace(query[:onIndex])

	// Extract table name and columns part
	tableColumnsPart := strings.TrimSpace(query[onIndex+4:]) // 4 = len(" ON ")

	// Find opening parenthesis
	openParenIndex := strings.Index(tableColumnsPart, "(")
	if openParenIndex == -1 {
		return &Statement{
			Type:    StatementTypeCreateIndex,
			RawQuery: query,
		}
	}

	// Extract table name
	tableName := strings.TrimSpace(tableColumnsPart[:openParenIndex])

	// Extract columns part
	columnsPart := strings.TrimSpace(tableColumnsPart[openParenIndex+1:]) // +1 to skip '('

	// Find closing parenthesis
	closeParenIndex := strings.Index(columnsPart, ")")
	if closeParenIndex == -1 {
		return &Statement{
			Type:    StatementTypeCreateIndex,
			RawQuery: query,
		}
	}

	// Extract columns
	columnsStr := strings.TrimSpace(columnsPart[:closeParenIndex])
	columns := p.parseColumns(columnsStr)

	return &Statement{
		Type: StatementTypeCreateIndex,
		CreateIndex: &CreateIndexStatement{
			IndexName: indexName,
			TableName: tableName,
			Columns:   columns,
			Unique:    unique,
		},
		RawQuery: query,
	}
}

// parseDropIndex parses a DROP INDEX statement
func (p *Parser) parseDropIndex(query string) *Statement {
	// Format: DROP INDEX index_name ON table_name

	upperQuery := strings.ToUpper(query)

	// Remove "DROP INDEX "
	query = strings.TrimPrefix(query, "DROP INDEX ")
	upperQuery = strings.TrimPrefix(upperQuery, "DROP INDEX ")

	// Find ON keyword
	onIndex := strings.Index(upperQuery, " ON ")
	if onIndex == -1 {
		// SQLite doesn't require ON clause for DROP INDEX
		return &Statement{
			Type: StatementTypeDropIndex,
			DropIndex: &DropIndexStatement{
				IndexName: strings.TrimSpace(query),
			},
			RawQuery: query,
		}
	}

	// Extract index name
	indexName := strings.TrimSpace(query[:onIndex])

	// Extract table name
	tableName := strings.TrimSpace(query[onIndex+4:]) // 4 = len(" ON ")

	return &Statement{
		Type: StatementTypeDropIndex,
		DropIndex: &DropIndexStatement{
			IndexName: indexName,
			TableName: tableName,
		},
		RawQuery: query,
	}
}

// parsePrepare parses a PREPARE statement
func (p *Parser) parsePrepare(query string) *Statement {
	// Format: PREPARE statement_name FROM 'sql_statement'
	// Format: PREPARE statement_name AS sql_statement

	upperQuery := strings.ToUpper(query)

	// Remove "PREPARE "
	query = strings.TrimPrefix(query, "PREPARE ")
	upperQuery = strings.TrimPrefix(upperQuery, "PREPARE ")

	// Extract statement name
	spaceIndex := strings.Index(query, " ")
	if spaceIndex == -1 {
		return &Statement{
			Type:    StatementTypePrepare,
			RawQuery: query,
		}
	}

	statementName := strings.TrimSpace(query[:spaceIndex])

	// Extract SQL statement
	sqlPart := strings.TrimSpace(query[spaceIndex:])
	upperSQLPart := strings.ToUpper(sqlPart)

	var sqlStatement string
	parameters := make([]string, 0)

	if strings.HasPrefix(upperSQLPart, "FROM ") {
		// Format: PREPARE statement_name FROM 'sql_statement'
		sqlPart = strings.TrimPrefix(sqlPart, "FROM ")
		sqlStatement = strings.TrimSpace(sqlPart)

		// Remove quotes if present
		if (strings.HasPrefix(sqlStatement, "'") && strings.HasSuffix(sqlStatement, "'")) ||
		   (strings.HasPrefix(sqlStatement, "\"") && strings.HasSuffix(sqlStatement, "\"")) {
			sqlStatement = strings.Trim(sqlStatement, "'\"")
		}
	} else if strings.HasPrefix(upperSQLPart, "AS ") {
		// Format: PREPARE statement_name AS sql_statement
		sqlPart = strings.TrimPrefix(sqlPart, "AS ")
		sqlStatement = strings.TrimSpace(sqlPart)
	} else {
		// Assume SQL statement follows directly
		sqlStatement = sqlPart
	}

	// Extract parameter placeholders (e.g., $1, $2, @param1, @param2)
	re := regexp.MustCompile(`[$@]([\w]+)`)
	matches := re.FindAllStringSubmatch(sqlStatement, -1)
	for _, match := range matches {
		if len(match) > 1 {
			parameters = append(parameters, match[1])
		}
	}

	return &Statement{
		Type: StatementTypePrepare,
		Prepare: &PrepareStatement{
			Name:       statementName,
			SQL:        sqlStatement,
			Parameters: parameters,
		},
		RawQuery: query,
	}
}

// executeStatement parses an EXECUTE statement
func (p *Parser) executeStatement(query string) *Statement {
	// Format: EXECUTE statement_name USING @param1 = value1, @param2 = value2
	// Format: EXECUTE statement_name

	upperQuery := strings.ToUpper(query)

	// Remove "EXECUTE "
	query = strings.TrimPrefix(query, "EXECUTE ")
	upperQuery = strings.TrimPrefix(upperQuery, "EXECUTE ")

	// Extract statement name
	spaceIndex := strings.Index(query, " ")
	if spaceIndex == -1 {
		// No USING clause
		return &Statement{
			Type: StatementTypeExecute,
			Execute: &ExecuteStatement{
				Name:       strings.TrimSpace(query),
				Parameters: make(map[string]interface{}),
			},
			RawQuery: query,
		}
	}

	statementName := strings.TrimSpace(query[:spaceIndex])

	// Extract USING clause
	usingPart := strings.TrimSpace(query[spaceIndex:])
	upperUsingPart := strings.ToUpper(usingPart)

	if !strings.HasPrefix(upperUsingPart, "USING ") {
		// Invalid USING clause
		return &Statement{
			Type: StatementTypeExecute,
			Execute: &ExecuteStatement{
				Name:       statementName,
				Parameters: make(map[string]interface{}),
			},
			RawQuery: query,
		}
	}

	usingPart = strings.TrimPrefix(usingPart, "USING ")

	// Parse parameters
	parameters := make(map[string]interface{})
	paramPairs := strings.Split(usingPart, ",")
	for _, pair := range paramPairs {
		pair = strings.TrimSpace(pair)
		equalIndex := strings.Index(pair, "=")
		if equalIndex == -1 {
			continue
		}

		paramName := strings.TrimSpace(pair[:equalIndex])
		paramValue := strings.TrimSpace(pair[equalIndex+1:])

		// Remove @ or $ prefix from parameter name
		paramName = strings.TrimPrefix(paramName, "@")
		paramName = strings.TrimPrefix(paramName, "$")

		// Remove quotes from parameter value
		if (strings.HasPrefix(paramValue, "'") && strings.HasSuffix(paramValue, "'")) ||
		   (strings.HasPrefix(paramValue, "\"") && strings.HasSuffix(paramValue, "\"")) {
			paramValue = strings.Trim(paramValue, "'\"")
		}

		parameters[paramName] = paramValue
	}

	return &Statement{
		Type: StatementTypeExecute,
		Execute: &ExecuteStatement{
			Name:       statementName,
			Parameters: parameters,
		},
		RawQuery: query,
	}
}

// parseDeallocatePrepare parses a DEALLOCATE PREPARE statement
func (p *Parser) parseDeallocatePrepare(query string) *Statement {
	// Format: DEALLOCATE PREPARE statement_name

	upperQuery := strings.ToUpper(query)

	// Remove "DEALLOCATE PREPARE "
	query = strings.TrimPrefix(query, "DEALLOCATE PREPARE ")
	upperQuery = strings.TrimPrefix(upperQuery, "DEALLOCATE PREPARE ")

	// Extract statement name
	statementName := strings.TrimSpace(query)

	return &Statement{
		Type: StatementTypeDeallocatePrepare,
		DeallocatePrepare: &DeallocatePrepareStatement{
			Name: statementName,
		},
		RawQuery: query,
	}
}

// parseBeginTransaction parses a BEGIN TRANSACTION statement
func (p *Parser) parseBeginTransaction(query string) *Statement {
	// Format: BEGIN TRANSACTION [name]
	// Format: BEGIN [name]
	// Format: START TRANSACTION [name]

	upperQuery := strings.ToUpper(query)

	// Remove "BEGIN TRANSACTION", "BEGIN", or "START TRANSACTION"
	query = strings.TrimPrefix(query, "BEGIN TRANSACTION ")
	query = strings.TrimPrefix(query, "BEGIN ")
	query = strings.TrimPrefix(query, "START TRANSACTION ")
	upperQuery = strings.TrimPrefix(upperQuery, "BEGIN TRANSACTION ")
	upperQuery = strings.TrimPrefix(upperQuery, "BEGIN ")
	upperQuery = strings.TrimPrefix(upperQuery, "START TRANSACTION ")

	// Extract transaction name (optional)
	transactionName := strings.TrimSpace(query)
	if upperQuery != "" {
		transactionName = strings.TrimSpace(upperQuery)
	}

	return &Statement{
		Type: StatementTypeBeginTransaction,
		BeginTransaction: &BeginTransactionStatement{
			Name: transactionName,
		},
		RawQuery: query,
	}
}

// parseCommit parses a COMMIT statement
func (p *Parser) parseCommit(query string) *Statement {
	// Format: COMMIT [name]
	// Format: COMMIT TRAN [name]

	upperQuery := strings.ToUpper(query)

	// Remove "COMMIT" or "COMMIT TRAN"
	query = strings.TrimPrefix(query, "COMMIT ")
	query = strings.TrimPrefix(query, "COMMIT TRAN ")
	upperQuery = strings.TrimPrefix(upperQuery, "COMMIT ")
	upperQuery = strings.TrimPrefix(upperQuery, "COMMIT TRAN ")

	// Extract transaction name (optional)
	transactionName := strings.TrimSpace(query)
	if upperQuery != "" {
		transactionName = strings.TrimSpace(upperQuery)
	}

	return &Statement{
		Type: StatementTypeCommit,
		Commit: &CommitStatement{
			Name: transactionName,
		},
		RawQuery: query,
	}
}

// parseRollback parses a ROLLBACK statement
func (p *Parser) parseRollback(query string) *Statement {
	// Format: ROLLBACK [name]
	// Format: ROLLBACK TRAN [name]
	// Format: ROLLBACK TO SAVEPOINT savepoint_name

	upperQuery := strings.ToUpper(query)

	// Check for "ROLLBACK TO SAVEPOINT"
	toSavepointIndex := strings.Index(upperQuery, " TO SAVEPOINT ")

	if toSavepointIndex != -1 {
		// ROLLBACK TO SAVEPOINT savepoint_name
		savepointName := strings.TrimSpace(query[toSavepointIndex+15:]) // 15 = len(" TO SAVEPOINT ")

		return &Statement{
			Type: StatementTypeRollback,
			Rollback: &RollbackStatement{
				SavepointName: savepointName,
			},
			RawQuery: query,
		}
	}

	// Remove "ROLLBACK" or "ROLLBACK TRAN"
	query = strings.TrimPrefix(query, "ROLLBACK ")
	query = strings.TrimPrefix(query, "ROLLBACK TRAN ")
	upperQuery = strings.TrimPrefix(upperQuery, "ROLLBACK ")
	upperQuery = strings.TrimPrefix(upperQuery, "ROLLBACK TRAN ")

	// Extract transaction name (optional)
	transactionName := strings.TrimSpace(query)
	if upperQuery != "" {
		transactionName = strings.TrimSpace(upperQuery)
	}

	return &Statement{
		Type: StatementTypeRollback,
		Rollback: &RollbackStatement{
			Name: transactionName,
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

		// Check if it's a table-level constraint (starts with PRIMARY, UNIQUE, FOREIGN, CHECK)
		upperColDef := strings.ToUpper(colDef)
		if strings.HasPrefix(upperColDef, "PRIMARY KEY") ||
		   strings.HasPrefix(upperColDef, "UNIQUE") ||
		   strings.HasPrefix(upperColDef, "FOREIGN KEY") ||
		   strings.HasPrefix(upperColDef, "CHECK") {
			// Skip table-level constraints (will be parsed separately)
			continue
		}

		// Split by first space to get column name and rest
		spaceIndex := strings.Index(colDef, " ")
		if spaceIndex == -1 {
			continue
		}

		columnName := strings.TrimSpace(colDef[:spaceIndex])
		columnPart := strings.TrimSpace(colDef[spaceIndex:])
		upperColumnPart := strings.ToUpper(columnPart)

		// Extract column type (first word)
		typeSpaceIndex := strings.Index(upperColumnPart, " ")
		var columnType string
		if typeSpaceIndex == -1 {
			columnType = columnPart
		} else {
			columnType = strings.TrimSpace(columnPart[:typeSpaceIndex])
			columnPart = strings.TrimSpace(columnPart[typeSpaceIndex:])
			upperColumnPart = strings.ToUpper(columnPart)
		}

		// Parse constraints
		colDefStruct := ColumnDefinition{
			Name: columnName,
			Type: columnType,
		}

		// PRIMARY KEY
		if strings.Contains(upperColumnPart, "PRIMARY KEY") {
			colDefStruct.PrimaryKey = true
		}

		// NOT NULL
		if strings.Contains(upperColumnPart, "NOT NULL") {
			colDefStruct.NotNull = true
		}

		// UNIQUE
		if strings.Contains(upperColumnPart, "UNIQUE") {
			colDefStruct.Unique = true
		}

		// DEFAULT value
		if idx := strings.Index(upperColumnPart, "DEFAULT "); idx != -1 {
			defaultPart := strings.TrimSpace(columnPart[idx+8:]) // 8 = len("DEFAULT ")
			// Remove quotes if present
			if (strings.HasPrefix(defaultPart, "'") && strings.HasSuffix(defaultPart, "'")) ||
			   (strings.HasPrefix(defaultPart, "\"") && strings.HasSuffix(defaultPart, "\"")) {
				defaultPart = strings.Trim(defaultPart, "'\"")
			}
			colDefStruct.DefaultValue = sql.NullString{String: defaultPart, Valid: true}
		}

		// CHECK constraint
		if idx := strings.Index(upperColumnPart, "CHECK "); idx != -1 {
			checkPart := strings.TrimSpace(columnPart[idx+6:]) // 6 = len("CHECK ")
			// Remove parentheses
			if strings.HasPrefix(checkPart, "(") && strings.HasSuffix(checkPart, ")") {
				checkPart = strings.Trim(checkPart, "()")
			}
			colDefStruct.Check = checkPart
		}

		// FOREIGN KEY
		if idx := strings.Index(upperColumnPart, "REFERENCES "); idx != -1 {
			refPart := strings.TrimSpace(columnPart[idx+11:]) // 11 = len("REFERENCES ")
			// Parse table_name(column_name)
			openParen := strings.Index(refPart, "(")
			if openParen != -1 {
				closeParen := strings.LastIndex(refPart, ")")
				if closeParen != -1 {
					refTable := strings.TrimSpace(refPart[:openParen])
					refColumn := strings.TrimSpace(refPart[openParen+1 : closeParen])
					colDefStruct.ForeignKey = &ForeignKeyConstraint{
						ReferenceTable:  refTable,
						ReferenceColumn: refColumn,
					}
				}
			}
		}

		result = append(result, colDefStruct)
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
	return stmtType >= StatementTypeSelect && stmtType <= StatementTypeRollback
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
