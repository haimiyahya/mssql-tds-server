package variable

import (
	"fmt"
	"regexp"
	"strings"
)

// Type represents the data type of a variable
type Type string

const (
	TypeInt     Type = "INT"
	TypeBigInt  Type = "BIGINT"
	TypeSmallInt Type = "SMALLINT"
	TypeTinyInt Type = "TINYINT"
	TypeVarchar Type = "VARCHAR"
	TypeNVarchar Type = "NVARCHAR"
	TypeChar    Type = "CHAR"
	TypeNChar   Type = "NCHAR"
	TypeBit     Type = "BIT"
	TypeDecimal Type = "DECIMAL"
	TypeFloat   Type = "FLOAT"
	TypeReal    Type = "REAL"
	TypeNumeric Type = "NUMERIC"
)

// Variable represents a declared variable
type Variable struct {
	Name       string
	Type       Type
	Length     int
	Precision  int
	Scale      int
	Value      interface{}
	IsNullable bool
	IsNull     bool
}

// Context represents a variable context for procedure execution
type Context struct {
	variables map[string]*Variable
}

// NewContext creates a new variable context
func NewContext() *Context {
	return &Context{
		variables: make(map[string]*Variable),
	}
}

// Declare adds a new variable declaration to the context
func (c *Context) Declare(name string, varType Type, length int) (*Variable, error) {
	// Normalize name (remove @ prefix)
	if strings.HasPrefix(name, "@") {
		name = strings.ToUpper(name[1:])
	} else {
		name = strings.ToUpper(name)
	}

	// Check if variable already exists
	if _, exists := c.variables[name]; exists {
		return nil, fmt.Errorf("variable '@%s' already declared", name)
	}

	// Create variable
	variable := &Variable{
		Name:       name,
		Type:       varType,
		Length:     length,
		Value:      nil,
		IsNull:     true,
		IsNullable: true,
	}

	c.variables[name] = variable
	return variable, nil
}

// Get retrieves a variable from the context
func (c *Context) Get(name string) (*Variable, bool) {
	// Normalize name (remove @ prefix)
	if strings.HasPrefix(name, "@") {
		name = strings.ToUpper(name[1:])
	} else {
		name = strings.ToUpper(name)
	}

	variable, exists := c.variables[name]
	return variable, exists
}

// Set sets a variable value
func (c *Context) Set(name string, value interface{}) error {
	// Normalize name
	if strings.HasPrefix(name, "@") {
		name = strings.ToUpper(name[1:])
	} else {
		name = strings.ToUpper(name)
	}

	variable, exists := c.variables[name]
	if !exists {
		return fmt.Errorf("variable '@%s' not declared", name)
	}

	// Set value
	if value == nil {
		variable.IsNull = true
		variable.Value = nil
	} else {
		variable.IsNull = false
		variable.Value = value
	}

	return nil
}

// Clear clears all variables from the context
func (c *Context) Clear() {
	c.variables = make(map[string]*Variable)
}

// GetAll returns all variables
func (c *Context) GetAll() map[string]*Variable {
	return c.variables
}

// ParseDeclaration parses a DECLARE statement
// Format: DECLARE @var TYPE [LENGTH]
func ParseDeclaration(sql string) (*Variable, error) {
	// Normalize SQL
	sql = strings.TrimSpace(sql)
	sql = strings.ToUpper(sql)

	// Check if it starts with DECLARE
	if !strings.HasPrefix(sql, "DECLARE") {
		return nil, fmt.Errorf("not a DECLARE statement")
	}

	// Remove DECLARE keyword
	sql = strings.TrimSpace(sql[7:])

	// Parse variable name and type
	// Format: @var TYPE [LENGTH]
	re := regexp.MustCompile(`@(\w+)\s+(\w+)(?:\((\d+)(?:,(\d+))?\))?`)
	matches := re.FindStringSubmatch(sql)

	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid DECLARE syntax: %s", sql)
	}

	name := matches[1]
	varType := Type(matches[2])

	// Parse length and precision
	length := 0
	precision := 0
	scale := 0

	if len(matches) > 3 && matches[3] != "" {
		fmt.Sscanf(matches[3], "%d", &length)
	}

	if len(matches) > 4 && matches[4] != "" {
		fmt.Sscanf(matches[4], "%d", &scale)
	}

	variable := &Variable{
		Name:       name,
		Type:       varType,
		Length:     length,
		Precision:  precision,
		Scale:      scale,
		Value:      nil,
		IsNull:     true,
		IsNullable: true,
	}

	return variable, nil
}

// ParseSetAssignment parses a SET assignment
// Format: SET @var = value
func ParseSetAssignment(sql string) (string, interface{}, error) {
	// Normalize SQL
	sql = strings.TrimSpace(sql)
	sql = strings.ToUpper(sql)

	// Check if it starts with SET
	if !strings.HasPrefix(sql, "SET") {
		return "", nil, fmt.Errorf("not a SET statement")
	}

	// Remove SET keyword
	sql = strings.TrimSpace(sql[3:])

	// Parse variable name and value
	// Format: @var = value
	re := regexp.MustCompile(`@(\w+)\s*=\s*(.+)`)
	matches := re.FindStringSubmatch(sql)

	if len(matches) < 3 {
		return "", nil, fmt.Errorf("invalid SET syntax: %s", sql)
	}

	varName := matches[1]
	valueStr := strings.TrimSpace(matches[2])

	// Try to parse value
	value, err := parseValue(valueStr)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse value: %w", err)
	}

	return varName, value, nil
}

// ParseSelectAssignment parses a SELECT variable assignment
// Format: SELECT @var = column FROM table WHERE ...
func ParseSelectAssignment(sql string) (string, string, error) {
	// Normalize SQL
	sql = strings.TrimSpace(sql)

	// Check if it's a SELECT with variable assignment
	// Format: SELECT @var = expression
	re := regexp.MustCompile(`SELECT\s+@(\w+)\s*=\s+(.+)`)
	matches := re.FindStringSubmatch(sql)

	if len(matches) < 3 {
		return "", "", fmt.Errorf("not a SELECT assignment statement")
	}

	varName := matches[1]
	expression := strings.TrimSpace(matches[2])

	return varName, expression, nil
}

// parseValue parses a value string into appropriate Go type
func parseValue(s string) (interface{}, error) {
	s = strings.TrimSpace(s)

	// String literal (in quotes)
	if strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'") {
		// Remove quotes and handle escaped quotes
		return strings.ReplaceAll(s[1:len(s)-1], "''", "'"), nil
	}

	// Integer
	if matches, _ := regexp.MatchString(`^-?\d+$`, s); matches {
		var value int
		_, err := fmt.Sscanf(s, "%d", &value)
		return value, err
	}

	// Float
	if matches, _ := regexp.MatchString(`^-?\d+\.\d+$`, s); matches {
		var value float64
		_, err := fmt.Sscanf(s, "%f", &value)
		return value, err
	}

	// NULL
	if strings.ToUpper(s) == "NULL" {
		return nil, nil
	}

	// Default to string
	return s, nil
}

// FormatValue formats a value for SQL
func FormatValue(value interface{}, varType Type) string {
	if value == nil {
		return "NULL"
	}

	switch v := value.(type) {
	case string:
		// Escape single quotes
		escaped := strings.ReplaceAll(v, "'", "''")
		return fmt.Sprintf("'%s'", escaped)
	case int, int32, int64, uint, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	default:
		// Default to string format
		escaped := strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''")
		return fmt.Sprintf("'%s'", escaped)
	}
}
