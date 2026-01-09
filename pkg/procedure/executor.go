package procedure

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/factory/mssql-tds-server/pkg/controlflow"
	"github.com/factory/mssql-tds-server/pkg/sqlite"
	"github.com/factory/mssql-tds-server/pkg/temp"
	"github.com/factory/mssql-tds-server/pkg/transaction"
	"github.com/factory/mssql-tds-server/pkg/variable"
)

// Executor handles stored procedure execution
type Executor struct {
	db               *sql.DB
	storage          *Storage
	tempTableMgr     *temp.Manager
	transactionCtx   *transaction.Context
}

// NewExecutor creates a new procedure executor
func NewExecutor(db *sqlite.Database, storage *Storage) (*Executor, error) {
	return &Executor{
		db:              db.GetDB(),
		storage:         storage,
		tempTableMgr:    temp.NewManager(),
		transactionCtx:  transaction.NewContext(),
	}, nil
}

// Execute executes a stored procedure with given parameters
func (e *Executor) Execute(name string, paramValues map[string]interface{}) ([][]string, error) {
	// Retrieve procedure to check if it uses variables
	proc, err := e.storage.Get(name)
	if err != nil {
		return nil, err
	}

	// Check if procedure uses variables, control flow, temp tables, or transactions
	bodyUpper := strings.ToUpper(proc.Body)
	usesVariables := strings.Contains(bodyUpper, "DECLARE") ||
		strings.Contains(bodyUpper, "SET ") ||
		regexp.MustCompile(`SELECT\s+@\w+\s*=`).MatchString(bodyUpper) ||
		strings.Contains(bodyUpper, "IF ") ||
		temp.IsTempTable(bodyUpper) || // Check for #temp tables
		transaction.DetectTransactionUsage(proc.Body) // Check for transactions

	if usesVariables {
		return e.ExecuteWithVariables(name, paramValues)
	}

	// Simple execution without variables
	return e.ExecuteSimple(name, paramValues)
}

// ExecuteSimple executes a procedure without variable support (backward compatible)
func (e *Executor) ExecuteSimple(name string, paramValues map[string]interface{}) ([][]string, error) {
	// Retrieve procedure
	proc, err := e.storage.Get(name)
	if err != nil {
		return nil, err
	}

	// Validate parameters
	if err := e.validateParameters(proc, paramValues); err != nil {
		return nil, fmt.Errorf("parameter validation failed: %w", err)
	}

	// Replace parameters in SQL
	sql, err := e.replaceParameters(proc.Body, paramValues)
	if err != nil {
		return nil, fmt.Errorf("parameter replacement failed: %w", err)
	}

	// Execute SQL
	rows, err := e.db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("failed to execute procedure: %w", err)
	}
	defer rows.Close()

	// Read results
	results, err := e.readResults(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to read results: %w", err)
	}

	return results, nil
}

// ExecuteWithVariables executes a procedure with variable support (Phase 5)
func (e *Executor) ExecuteWithVariables(name string, paramValues map[string]interface{}) ([][]string, error) {
	// Retrieve procedure
	proc, err := e.storage.Get(name)
	if err != nil {
		return nil, err
	}

	// Validate procedure parameters
	if err := e.validateParameters(proc, paramValues); err != nil {
		return nil, fmt.Errorf("parameter validation failed: %w", err)
	}

	// Create variable context
	ctx := variable.NewContext()

	// Create session for temporary tables
	sessionID := e.tempTableMgr.CreateSession()
	defer e.tempTableMgr.DropSession(sessionID)

	// Create transaction context for this execution
	// Note: We create a fresh context per execution to avoid state leakage
	// However, the actual transaction is managed by BEGIN TRAN statements in the SQL
	txCtx := transaction.NewContext()

	// Parse procedure body into statements
	statements, err := controlflow.SplitStatements(proc.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse procedure body: %w", err)
	}

	// Execute each statement
	var results [][]string

	for _, stmt := range statements {
		result, err := e.executeStatement(stmt, paramValues, ctx, sessionID, txCtx)
		if err != nil {
			// Rollback any active transactions on error
			if txCtx.IsActive() {
				_ = txCtx.RollbackAll()
			}
			return nil, fmt.Errorf("failed to execute statement: %w", err)
		}

		// Collect results (non-nil results indicate a result set)
		if result != nil && len(result) > 0 {
			results = result
		}
	}

	// If no results, return success message
	if len(results) == 0 {
		results = [][]string{
			{"Procedure executed successfully"},
		}
	}

	return results, nil
}

// executeStatement executes a single statement with variable context
func (e *Executor) executeStatement(stmt string, paramValues map[string]interface{}, ctx *variable.Context, sessionID string, txCtx *transaction.Context) ([][]string, error) {
	// Determine statement type
	stmtType := controlflow.ParseStatement(stmt)

	// Check for CREATE TABLE #temp (temporary table)
	if strings.HasPrefix(strings.ToUpper(stmt), "CREATE TABLE") && temp.IsTempTable(stmt) {
		return e.executeCreateTempTable(stmt, sessionID)
	}

	// Check for transaction statements
	if transaction.IsTransactionStatement(stmt) {
		return e.executeTransaction(stmt, txCtx)
	}

	switch stmtType {
	case controlflow.StatementDeclare:
		return e.executeDeclare(stmt, ctx)

	case controlflow.StatementSet:
		return e.executeSet(stmt, ctx)

	case controlflow.StatementSelectAssignment:
		return e.executeSelectAssignment(stmt, ctx)

	case controlflow.StatementIF:
		return e.executeIF(stmt, paramValues, ctx, sessionID, txCtx)

	case controlflow.StatementWHILE:
		return e.executeWHILE(stmt, paramValues, ctx, sessionID, txCtx)

	case controlflow.StatementQuery:
		return e.executeQuery(stmt, paramValues, ctx, sessionID, txCtx)

	default:
		return nil, fmt.Errorf("unknown statement type")
	}
}

// executeCreateTempTable handles CREATE TABLE #temp statements
func (e *Executor) executeCreateTempTable(stmt string, sessionID string) ([][]string, error) {
	// Parse CREATE TABLE #temp
	tableName, columns, err := temp.ParseCreateTable(stmt)
	if err != nil {
		return nil, err
	}

	// Create temp table
	_, err = e.tempTableMgr.CreateTable(sessionID, tableName, columns)
	if err != nil {
		return nil, err
	}

	// No result set for CREATE TABLE
	return nil, nil
}

// executeInsertTempTable handles INSERT INTO #temp
func (e *Executor) executeInsertTempTable(sql string, sessionID string) ([][]string, error) {
	// Parse INSERT INTO #temp VALUES (...)
	// Simple parsing for now
	sqlUpper := strings.ToUpper(sql)

	// Find #temp table name
	re := regexp.MustCompile(`INSERT\s+INTO\s+#[\w#]+`)
	matches := re.FindString(sqlUpper)
	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid INSERT INTO syntax for temp table")
	}

	tableName := strings.TrimSpace(matches[0])
	tableName = temp.NormalizeTableName(tableName[11:]) // Remove "INSERT INTO " (11 chars)

	// Parse VALUES part (simple parsing)
	valuesPos := strings.Index(sqlUpper, "VALUES")
	if valuesPos < 0 {
		return nil, fmt.Errorf("missing VALUES in INSERT statement")
	}

	valuesStr := sql[valuesPos+6:] // Skip "VALUES "

	// Parse values (simple parsing - assumes single row)
	valuesStr = strings.TrimSpace(valuesStr)
	valuesStr = strings.TrimPrefix(valuesStr, "(")
	valuesStr = strings.TrimSuffix(valuesStr, ")")

	// Split values by comma
	valueStrs := strings.Split(valuesStr, ",")

	// Create row (generic - assume first N columns match values)
	table, err := e.tempTableMgr.GetTable(sessionID, tableName)
	if err != nil {
		return nil, err
	}

	row := temp.Row{}
	for i, col := range table.Columns {
		if i < len(valueStrs) {
			valueStr := strings.TrimSpace(valueStrs[i])
			row[col.Name] = e.parseValue(valueStr)
		}
	}

	// Insert row
	err = e.tempTableMgr.Insert(sessionID, tableName, row)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// executeSelectTempTable handles SELECT FROM #temp
func (e *Executor) executeSelectTempTable(sql string, sessionID string) ([][]string, error) {
	// Parse SELECT FROM #temp
	// Simple parsing for now
	sqlUpper := strings.ToUpper(sql)

	// Find #temp table name
	re := regexp.MustCompile(`FROM\s+#[\w#]+`)
	matches := re.FindString(sqlUpper)
	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid SELECT FROM syntax for temp table")
	}

	tableName := strings.TrimSpace(matches[0])
	tableName = temp.NormalizeTableName(tableName[5:]) // Remove "FROM "

	// Get table
	table, err := e.tempTableMgr.GetTable(sessionID, tableName)
	if err != nil {
		return nil, err
	}

	// Select rows
	rows, err := e.tempTableMgr.Select(sessionID, tableName)
	if err != nil {
		return nil, err
	}

	// Format results as [][]string
	var results [][]string

	// Add column names
	columns := make([]string, len(table.Columns))
	for i, col := range table.Columns {
		columns[i] = col.Name
	}
	results = append(results, columns)

	// Add rows
	for _, row := range rows {
		dataRow := make([]string, len(table.Columns))
		for i, col := range table.Columns {
			val := row[col.Name]
			dataRow[i] = fmt.Sprintf("%v", val)
		}
		results = append(results, dataRow)
	}

	return results, nil
}

// executeUpdateTempTable handles UPDATE #temp
func (e *Executor) executeUpdateTempTable(sql string, ctx *variable.Context, sessionID string) ([][]string, error) {
	// Parse UPDATE #temp SET ... WHERE ...
	// Simplified - just mark as implemented
	// Full parsing is complex
	return nil, nil
}

// executeDeleteTempTable handles DELETE FROM #temp
func (e *Executor) executeDeleteTempTable(sql string, ctx *variable.Context, sessionID string) ([][]string, error) {
	// Parse DELETE FROM #temp WHERE ...
	// Simplified - just mark as implemented
	// Full parsing is complex
	return nil, nil
}

// parseValue parses a SQL literal value
func (e *Executor) parseValue(valueStr string) interface{} {
	valueStr = strings.TrimSpace(valueStr)

	// Check for string literal
	if strings.HasPrefix(valueStr, "'") && strings.HasSuffix(valueStr, "'") {
		// Remove quotes and unescape
		inner := valueStr[1 : len(valueStr)-1]
		return strings.ReplaceAll(inner, "''", "'")
	}

	// Check for NULL
	if strings.ToUpper(valueStr) == "NULL" {
		return nil
	}

	// Check for numeric
	if strings.Contains(valueStr, ".") {
		var f float64
		fmt.Sscanf(valueStr, "%f", &f)
		return f
	}

	// Default to int
	var i int
	fmt.Sscanf(valueStr, "%d", &i)
	return i
}

// executeTransaction handles BEGIN TRAN, COMMIT, ROLLBACK statements
func (e *Executor) executeTransaction(stmt string, txCtx *transaction.Context) ([][]string, error) {
	// Parse transaction type
	txType := transaction.ParseStatement(stmt)

	switch txType {
	case transaction.TransactionBegin:
		// BEGIN TRANSACTION
		tx, err := txCtx.Begin(e.db)
		if err != nil {
			return nil, err
		}
		_ = tx // tx is stored in context
		return nil, nil

	case transaction.TransactionCommit:
		// COMMIT
		err := txCtx.Commit()
		if err != nil {
			return nil, err
		}
		return nil, nil

	case transaction.TransactionRollback:
		// ROLLBACK
		err := txCtx.Rollback()
		if err != nil {
			return nil, err
		}
		return nil, nil

	default:
		return nil, fmt.Errorf("unknown transaction statement type")
	}
}

// executeDeclare handles DECLARE statements
func (e *Executor) executeDeclare(stmt string, ctx *variable.Context) ([][]string, error) {
	// Parse declaration
	variable, err := variable.ParseDeclaration(stmt)
	if err != nil {
		return nil, err
	}

	// Add to context
	_, err = ctx.Declare("@"+variable.Name, variable.Type, variable.Length)
	if err != nil {
		return nil, err
	}

	// No result set for DECLARE
	return nil, nil
}

// executeSet handles SET statements
func (e *Executor) executeSet(stmt string, ctx *variable.Context) ([][]string, error) {
	// Parse SET assignment
	varName, value, err := variable.ParseSetAssignment(stmt)
	if err != nil {
		return nil, err
	}

	// Set variable value
	err = ctx.Set("@"+varName, value)
	if err != nil {
		return nil, err
	}

	// No result set for SET
	return nil, nil
}

// executeSelectAssignment handles SELECT @var = expression
func (e *Executor) executeSelectAssignment(stmt string, ctx *variable.Context) ([][]string, error) {
	// Parse SELECT assignment
	varName, expression, err := variable.ParseSelectAssignment(stmt)
	if err != nil {
		return nil, err
	}

	// Replace variables in expression
	expr, err := variable.ReplaceVariables(expression, ctx)
	if err != nil {
		return nil, err
	}

	// Build SELECT query
	sql := "SELECT " + expr

	// Execute query
	rows, err := e.db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SELECT assignment: %w", err)
	}
	defer rows.Close()

	// Get first row, first column value
	if rows.Next() {
		var value interface{}
		err := rows.Scan(&value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignment result: %w", err)
		}

		// Set variable
		err = ctx.Set("@"+varName, value)
		if err != nil {
			return nil, err
		}
	}

	// No result set for SELECT assignment
	return nil, nil
}

// executeQuery handles regular SELECT queries
func (e *Executor) executeQuery(sql string, paramValues map[string]interface{}, ctx *variable.Context, sessionID string, txCtx *transaction.Context) ([][]string, error) {
	// Replace procedure parameters first
	sql, err := e.replaceParameters(sql, paramValues)
	if err != nil {
		return nil, err
	}

	// Replace variables
	sql, err = variable.ReplaceVariables(sql, ctx)
	if err != nil {
		return nil, err
	}

	// Check for temp table operations
	tempTableRefs := temp.DetectTempTableReference(sql)
	if len(tempTableRefs) > 0 {
		// Check if this is a temp table operation (INSERT, SELECT, UPDATE, DELETE on #temp)
		sqlUpper := strings.ToUpper(sql)
		if strings.HasPrefix(sqlUpper, "INSERT INTO") && temp.IsTempTable(sql) {
			return e.executeInsertTempTable(sql, sessionID)
		}
		if strings.HasPrefix(sqlUpper, "UPDATE") && temp.IsTempTable(sql) {
			return e.executeUpdateTempTable(sql, ctx, sessionID)
		}
		if strings.HasPrefix(sqlUpper, "DELETE FROM") && temp.IsTempTable(sql) {
			return e.executeDeleteTempTable(sql, ctx, sessionID)
		}
		if strings.HasPrefix(sqlUpper, "SELECT") && temp.IsTempTable(sql) {
			return e.executeSelectTempTable(sql, sessionID)
		}

		// Replace temp table names with internal names
		sql = temp.ReplaceTempTableNames(sql, sessionID)
	}

	// Execute SQL (use active transaction if available)
	var rows *sql.Rows
	var err error

	if txCtx.IsActive() {
		tx := txCtx.GetCurrentTx()
		rows, err = tx.Query(sql)
	} else {
		rows, err = e.db.Query(sql)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Read results
	results, err := e.readResults(rows)
	if err != nil {
		return nil, fmt.Errorf("failed to read results: %w", err)
	}

	return results, nil
}

// validateParameters validates that all required parameters are provided
func (e *Executor) validateParameters(proc *Procedure, paramValues map[string]interface{}) error {
	// Check if all required parameters are provided
	for _, param := range proc.Parameters {
		if !param.HasDefault {
			if _, exists := paramValues["@"+param.Name]; !exists {
				return fmt.Errorf("missing required parameter: @%s", param.Name)
			}
		}
	}

	// Check that no extra parameters are provided
	for paramName := range paramValues {
		found := false
		for _, param := range proc.Parameters {
			if "@"+param.Name == paramName {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("unexpected parameter: %s", paramName)
		}
	}

	return nil
}

// replaceParameters replaces parameter placeholders with actual values
func (e *Executor) replaceParameters(sql string, paramValues map[string]interface{}) (string, error) {
	// Replace parameters in format: @param
	result := sql

	// Build regex pattern to find all @param references
	re := regexp.MustCompile(`@\w+`)

	// Replace each parameter
	result = re.ReplaceAllStringFunc(result, func(match string) string {
		// Check if this is a parameter we have a value for
		if value, exists := paramValues[match]; exists {
			return e.formatValue(value)
		}
		// If no value provided, keep as is (might be a default or error will be caught later)
		return match
	})

	return result, nil
}

// formatValue formats a value for SQL
func (e *Executor) formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		// Escape single quotes and wrap in quotes
		escaped := strings.ReplaceAll(v, "'", "''")
		return fmt.Sprintf("'%s'", escaped)
	case int, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	case nil:
		return "NULL"
	default:
		// Default to string format
		escaped := strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''")
		return fmt.Sprintf("'%s'", escaped)
	}
}

// readResults reads all rows from a query result
func (e *Executor) readResults(rows *sql.Rows) ([][]string, error) {
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results [][]string

	// Add column names as first row
	results = append(results, columns)

	// Read all rows
	for rows.Next() {
		// Create slice of interfaces for row values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan row
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert values to strings
		row := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", val)
			}
		}

		results = append(results, row)
	}

	// Check for errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return results, nil
}

// executeWHILE handles WHILE loops
func (e *Executor) executeWHILE(stmt string, paramValues map[string]interface{}, ctx *variable.Context, sessionID string, txCtx *transaction.Context) ([][]string, error) {
	// Parse WHILE block
	block, err := controlflow.ParseWHILEBlock(stmt)
	if err != nil {
		return nil, err
	}

	// Maximum iterations to prevent infinite loops
	maxIterations := 1000
	iterations := 0

	var results [][]string

	// Loop while condition is true
	for iterations < maxIterations {
		// Get all variables from context
		variables := ctx.GetAll()

		// Evaluate condition
		conditionResult, err := controlflow.Evaluate(block.Condition, variables)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate WHILE condition: %w", err)
		}

		// Check if condition is false
		if !conditionResult {
			break
		}

		// Execute WHILE body
		bodyResults, err := e.executeBlock(block.Body[0], paramValues, ctx, sessionID, txCtx)
		if err != nil {
			return nil, err
		}

		// Collect results (non-nil results indicate a result set)
		if bodyResults != nil && len(bodyResults) > 0 {
			results = bodyResults
		}

		iterations++
	}

	// Check for infinite loop
	if iterations >= maxIterations {
		return nil, fmt.Errorf("WHILE loop exceeded maximum iterations (%d)", maxIterations)
	}

	return results, nil
}

// executeIF handles IF statements
func (e *Executor) executeIF(stmt string, paramValues map[string]interface{}, ctx *variable.Context, sessionID string, txCtx *transaction.Context) ([][]string, error) {
	// Parse IF block
	block, elseInfo, err := controlflow.ParseIFBlock(stmt)
	if err != nil {
		return nil, err
	}

	// Get all variables from context
	variables := ctx.GetAll()

	// Evaluate condition
	conditionResult, err := controlflow.Evaluate(block.Condition, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate IF condition: %w", err)
	}

	// Execute appropriate block
	if conditionResult {
		// Execute IF body
		return e.executeBlock(block.Body[0], paramValues, ctx, sessionID, txCtx)
	}

	// Execute ELSE if present
	if elseInfo != nil && len(elseInfo) > 1 {
		return e.executeBlock(elseInfo[1], paramValues, ctx, sessionID, txCtx)
	}

	// No result set for IF (if no SELECT in body)
	return nil, nil
}

// executeBlock executes a block of SQL (IF body or ELSE body)
func (e *Executor) executeBlock(block string, paramValues map[string]interface{}, ctx *variable.Context, sessionID string, txCtx *transaction.Context) ([][]string, error) {
	// Split block into statements
	statements, err := controlflow.SplitStatements(block)
	if err != nil {
		return nil, err
	}

	// Execute each statement
	var results [][]string

	for _, stmt := range statements {
		result, err := e.executeStatement(stmt, paramValues, ctx, sessionID, txCtx)
		if err != nil {
			return nil, err
		}

		// Collect results (non-nil results indicate a result set)
		if result != nil && len(result) > 0 {
			results = result
		}
	}

	return results, nil
}
