package sqlparser

// StatementType represents the type of SQL statement
type StatementType int

const (
	StatementTypeUnknown StatementType = iota
	StatementTypeSelect
	StatementTypeInsert
	StatementTypeUpdate
	StatementTypeDelete
	StatementTypeCreateTable
	StatementTypeDropTable
	StatementTypeCreateView
	StatementTypeDropView
	StatementTypeBeginTransaction
	StatementTypeCommit
	StatementTypeRollback
)

// String returns the string representation of StatementType
func (st StatementType) String() string {
	switch st {
	case StatementTypeSelect:
		return "SELECT"
	case StatementTypeInsert:
		return "INSERT"
	case StatementTypeUpdate:
		return "UPDATE"
	case StatementTypeDelete:
		return "DELETE"
	case StatementTypeCreateTable:
		return "CREATE TABLE"
	case StatementTypeDropTable:
		return "DROP TABLE"
	case StatementTypeCreateView:
		return "CREATE VIEW"
	case StatementTypeDropView:
		return "DROP VIEW"
	case StatementTypeBeginTransaction:
		return "BEGIN TRANSACTION"
	case StatementTypeCommit:
		return "COMMIT"
	case StatementTypeRollback:
		return "ROLLBACK"
	default:
		return "UNKNOWN"
	}
}

// SelectStatement represents a SELECT statement
type SelectStatement struct {
	Columns           []string
	Table             string
	Joins             []JoinClause
	WhereClause        string
	Distinct          bool
	OrderBy           []OrderByClause
	Aggregates        []AggregateFunction
	IsAggregateQuery  bool
	GroupBy           []GroupByClause
	HavingClause      string
	HasSubqueries     bool // True if query contains subqueries
}

// OrderByClause represents an ORDER BY clause
type OrderByClause struct {
	Column    string
	Direction string // "ASC" or "DESC"
}

// GroupByClause represents a GROUP BY clause
type GroupByClause struct {
	Column string
}

// JoinClause represents a JOIN clause
type JoinClause struct {
	Type      string // "INNER", "LEFT", "RIGHT", "FULL"
	Table     string
	OnClause  string
	Alias     string // Optional table alias
}

// AggregateFunction represents an aggregate function
type AggregateFunction struct {
	Type   string // "COUNT", "SUM", "AVG", "MIN", "MAX"
	Column string
	Alias  string // Optional AS alias
}

// InsertStatement represents an INSERT statement
type InsertStatement struct {
	Table    string
	Columns  []string
	Values   [][]interface{}
}

// UpdateStatement represents an UPDATE statement
type UpdateStatement struct {
	Table       string
	SetClause   string
	WhereClause string
}

// DeleteStatement represents a DELETE statement
type DeleteStatement struct {
	Table       string
	WhereClause string
}

// CreateTableStatement represents a CREATE TABLE statement
type CreateTableStatement struct {
	TableName string
	Columns   []ColumnDefinition
}

// DropTableStatement represents a DROP TABLE statement
type DropTableStatement struct {
	TableName string
}

// CreateViewStatement represents a CREATE VIEW statement
type CreateViewStatement struct {
	ViewName  string
	SelectQuery string // The SELECT query that defines the view
}

// DropViewStatement represents a DROP VIEW statement
type DropViewStatement struct {
	ViewName string
}

// BeginTransactionStatement represents a BEGIN TRANSACTION statement
type BeginTransactionStatement struct {
	Name string // Optional transaction name (for named transactions)
}

// CommitStatement represents a COMMIT statement
type CommitStatement struct {
	Name string // Optional transaction name (for named transactions)
}

// RollbackStatement represents a ROLLBACK statement
type RollbackStatement struct {
	Name         string // Optional transaction name (for named transactions)
	SavepointName string // Optional savepoint name for ROLLBACK TO SAVEPOINT
}

// ColumnDefinition represents a column definition in CREATE TABLE
type ColumnDefinition struct {
	Name string
	Type string
}

// Statement represents a parsed SQL statement
type Statement struct {
	Type StatementType

	Select                 *SelectStatement
	Insert                 *InsertStatement
	Update                 *UpdateStatement
	Delete                 *DeleteStatement
	CreateTable            *CreateTableStatement
	DropTable              *DropTableStatement
	CreateView             *CreateViewStatement
	DropView               *DropViewStatement
	BeginTransaction       *BeginTransactionStatement
	Commit                *CommitStatement
	Rollback              *RollbackStatement
	RawQuery               string
}
