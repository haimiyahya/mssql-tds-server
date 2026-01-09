package temp

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// Column represents a table column
type Column struct {
	Name     string
	Type     string
	Nullable bool
}

// Row represents a table row (map of column name to value)
type Row map[string]interface{}

// Table represents a temporary table
type Table struct {
	Name       string
	Columns   []Column
	Rows       []Row
	SessionID  string
}

// Manager manages temporary tables for sessions
type Manager struct {
	mu       sync.RWMutex
	tables   map[string]*Table // key: sessionID + tableName
	sessions  map[string]bool   // track active sessions
	nextID    int
}

// NewManager creates a new temporary table manager
func NewManager() *Manager {
	return &Manager{
		tables:  make(map[string]*Table),
		sessions: make(map[string]bool),
		nextID:   1,
	}
}

// CreateSession creates a new session for temporary tables
func (m *Manager) CreateSession() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	sessionID := fmt.Sprintf("session_%d", m.nextID)
	m.nextID++

	m.sessions[sessionID] = true

	return sessionID
}

// DropSession removes all temporary tables for a session
func (m *Manager) DropSession(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Remove all tables for this session
	for key := range m.tables {
		if strings.HasPrefix(key, sessionID+"_") {
			delete(m.tables, key)
		}
	}

	// Remove session
	delete(m.sessions, sessionID)
}

// CreateTable creates a temporary table
func (m *Manager) CreateTable(sessionID, name string, columns []Column) (*Table, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate session
	if !m.sessions[sessionID] {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	// Check if table already exists
	key := sessionID + "_" + name
	if _, exists := m.tables[key]; exists {
		return nil, fmt.Errorf("temporary table already exists: %s", name)
	}

	// Create table
	table := &Table{
		Name:      name,
		Columns:   columns,
		Rows:      []Row{},
		SessionID: sessionID,
	}

	m.tables[key] = table

	return table, nil
}

// DropTable drops a temporary table
func (m *Manager) DropTable(sessionID, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := sessionID + "_" + name
	if _, exists := m.tables[key]; !exists {
		return fmt.Errorf("temporary table not found: %s", name)
	}

	delete(m.tables, key)

	return nil
}

// GetTable gets a temporary table
func (m *Manager) GetTable(sessionID, name string) (*Table, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := sessionID + "_" + name
	table, exists := m.tables[key]
	if !exists {
		return nil, fmt.Errorf("temporary table not found: %s", name)
	}

	return table, nil
}

// Insert inserts a row into a temporary table
func (m *Manager) Insert(sessionID, name string, row Row) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := sessionID + "_" + name
	table, exists := m.tables[key]
	if !exists {
		return fmt.Errorf("temporary table not found: %s", name)
	}

	// Validate row columns
	for colName := range row {
		found := false
		for _, col := range table.Columns {
			if col.Name == colName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("column not found: %s", colName)
		}
	}

	table.Rows = append(table.Rows, row)

	return nil
}

// Select selects rows from a temporary table
func (m *Manager) Select(sessionID, name string) ([]Row, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := sessionID + "_" + name
	table, exists := m.tables[key]
	if !exists {
		return nil, fmt.Errorf("temporary table not found: %s", name)
	}

	return table.Rows, nil
}

// Update updates rows in a temporary table
func (m *Manager) Update(sessionID, name string, condition func(Row) bool, updates Row) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := sessionID + "_" + name
	table, exists := m.tables[key]
	if !exists {
		return 0, fmt.Errorf("temporary table not found: %s", name)
	}

	count := 0
	for _, row := range table.Rows {
		if condition(row) {
			for colName, value := range updates {
				// Validate column exists
				found := false
				for _, col := range table.Columns {
					if col.Name == colName {
						found = true
						break
					}
				}
				if found {
					row[colName] = value
				}
			}
			count++
		}
	}

	return count, nil
}

// Delete deletes rows from a temporary table
func (m *Manager) Delete(sessionID, name string, condition func(Row) bool) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := sessionID + "_" + name
	table, exists := m.tables[key]
	if !exists {
		return 0, fmt.Errorf("temporary table not found: %s", name)
	}

	var newRows []Row
	deletedCount := 0

	for _, row := range table.Rows {
		if condition(row) {
			deletedCount++
		} else {
			newRows = append(newRows, row)
		}
	}

	table.Rows = newRows

	return deletedCount, nil
}

// IsTempTable checks if a table name is a temporary table
func IsTempTable(name string) bool {
	return strings.HasPrefix(name, "#") || strings.HasPrefix(name, "##")
}

// NormalizeTableName normalizes a temporary table name (removes # prefix)
func NormalizeTableName(name string) string {
	return strings.TrimLeft(name, "#")
}

// ParseCreateTable parses a CREATE TABLE statement for temporary tables
// Format: CREATE TABLE #temp (col1 type, col2 type, ...)
func ParseCreateTable(sql string) (name string, columns []Column, err error) {
	sql = strings.TrimSpace(sql)
	sqlUpper := strings.ToUpper(sql)

	// Check if it's a CREATE TABLE statement
	if !strings.HasPrefix(sqlUpper, "CREATE TABLE") {
		return "", nil, fmt.Errorf("not a CREATE TABLE statement")
	}

	// Remove CREATE TABLE
	afterCreateTable := strings.TrimSpace(sql[12:])

	// Extract table name (first word)
	// Find first opening parenthesis
	openParenPos := strings.Index(afterCreateTable, "(")
	if openParenPos < 0 {
		return "", nil, fmt.Errorf("missing opening parenthesis in column definitions")
	}

	tableName := strings.TrimSpace(afterCreateTable[:openParenPos])

	// Check if it's a temporary table
	if !IsTempTable(tableName) {
		return "", nil, fmt.Errorf("not a temporary table (missing # prefix)")
	}

	// Normalize table name
	normalizedName := NormalizeTableName(tableName)

	// Extract column definitions (between parentheses)
	columnDefs := strings.TrimSpace(afterCreateTable[openParenPos+1:])

	// Find closing parenthesis
	closeParenPos := strings.LastIndex(columnDefs, ")")
	if closeParenPos < 0 {
		return "", nil, fmt.Errorf("missing closing parenthesis in column definitions")
	}

	columnDefs = columnDefs[:closeParenPos]

	// Parse columns
	columns, err = parseColumnDefinitions(columnDefs)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse column definitions: %w", err)
	}

	return normalizedName, columns, nil
}

// parseColumnDefinitions parses column definitions from CREATE TABLE
func parseColumnDefinitions(defs string) ([]Column, error) {
	var columns []Column

	// Split by comma, but handle nested parentheses
	current := ""
	inParentheses := 0

	for _, ch := range defs {
		switch ch {
		case '(':
			inParentheses++
			current += string(ch)
		case ')':
			inParentheses--
			current += string(ch)
		case ',':
			if inParentheses == 0 {
				col, err := parseColumnDefinition(strings.TrimSpace(current))
				if err != nil {
					return nil, err
				}
				columns = append(columns, col)
				current = ""
			} else {
				current += string(ch)
			}
		default:
			current += string(ch)
		}
	}

	// Add last column
	if current != "" {
		col, err := parseColumnDefinition(strings.TrimSpace(current))
		if err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return columns, nil
}

// parseColumnDefinition parses a single column definition
// Format: name type [NULL|NOT NULL]
func parseColumnDefinition(def string) (Column, error) {
	// Split by whitespace
	parts := strings.Fields(def)

	if len(parts) < 2 {
		return Column{}, fmt.Errorf("invalid column definition: %s", def)
	}

	col := Column{
		Name: parts[0],
		Type: parts[1],
	}

	// Check for NULL/NOT NULL
	for i := 2; i < len(parts); i++ {
		if strings.ToUpper(parts[i]) == "NOT" && i+1 < len(parts) && strings.ToUpper(parts[i+1]) == "NULL" {
			col.Nullable = false
			break
		}
	}

	return col, nil
}

// DetectTempTableReference detects if SQL contains references to temporary tables
func DetectTempTableReference(sql string) []string {
	// Find all #table references
	re := regexp.MustCompile(`#[\w#]+`)
	matches := re.FindAllString(sql, -1)

	// Deduplicate
	seen := make(map[string]bool)
	var tables []string

	for _, match := range matches {
		if !seen[match] {
			seen[match] = true
			tables = append(tables, match)
		}
	}

	return tables
}

// ReplaceTempTableNames replaces #table names with internal names for a session
func ReplaceTempTableNames(sql, sessionID string) string {
	// Find all #table references
	re := regexp.MustCompile(`#[\w#]+`)

	return re.ReplaceAllStringFunc(sql, func(match string) string {
		// Return session-prefixed name
		return sessionID + "_" + NormalizeTableName(match)
	})
}
