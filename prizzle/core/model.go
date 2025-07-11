package prizzle

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
// Custom Data Types ---------------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------------------------------

type DateTime struct {
	time.Time
}

func (dt *DateTime) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = str[1 : len(str)-1]

	parsedTime, err := time.Parse("2006-01-02T15:04:05.999999", str)
	if err != nil {
		return err
	}

	*dt = DateTime{parsedTime}

	return nil
}

func (dt *DateTime) MarshalJSON() ([]byte, error) {
	return []byte(dt.Time.Format("2006-01-02T15:04:05.999999")), nil
}

func (dt *DateTime) Scan(value interface{}) error {
	if value == nil {
		*dt = DateTime{Time: time.Time{}}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*dt = DateTime{Time: v}
		return nil
	default:
		return fmt.Errorf("unsupported Scan type for DateTime: %T", value)
	}
}

func (dt *DateTime) Value() (driver.Value, error) {
	if dt.Time.IsZero() {
		return nil, nil
	}
	return dt.Time, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// Model ---------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------------------------------

type SqlQuery struct {
	Client       DatabaseClient
	WithStr      string
	BaseStr      string
	FromStr      string
	WhereStr     string
	ReturningStr string
	GroupByStr   string
	OrderStr     string
	LimitStr     string
	OffsetStr    string
	Args         *[]interface{}
}

type PreparedSqlQuery struct {
	Client      DatabaseClient
	QueryString string
	Args        []interface{}
}

type SqlQueryString string

type SqlCondition string

type SqlOperation string

type SqlOrderDirection string

type SqlTable struct {
	Name string
}

type SqlName string

type SqlValues = map[SqlName]interface{}

type SqlValue struct {
	Prefix string
	Value  interface{}
	Suffix string
}

type SqlJson = map[SqlName]SqlName

type EmbedsSqlTable interface {
	GetSqlTable() SqlTable
	As(alias string) EmbedsSqlTable
}

type Cluster struct {
	WriteNodes    map[string]*DB `json:"writeNodes"`
	ReadNodes     map[string]*DB `json:"readNodes"`
	writeNodeKeys []string
	readNodeKeys  []string
}

type ClusterConfig struct {
	WriteNodesConfig map[string]ClusterNodeConfig `json:"writeNodesConfig"`
	ReadNodesConfig  map[string]ClusterNodeConfig `json:"readNodesConfig"`
}

type ClusterNodeConfig struct {
	Url                string         `json:"url"`
	MaxOpenConnections *int           `json:"maxOpenConnections"`
	MaxIdleConnections *int           `json:"maxIdleConnections"`
	MaxIdleTime        *time.Duration `json:"maxIdleTime"`
	MaxLifetime        *time.Duration `json:"maxLifetime"`
}

type DB struct {
	*sql.DB
}

type Tx struct {
	*sql.Tx
}

type Stmt struct {
	*sql.Stmt
	Args []interface{}
}

type Conn struct {
	*sql.Conn
}

// ---------------------------------------------------------------------------------------------------------------------
// Behavior ------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------------------------------

/* To String -------------------------------------------------------------------------------------------------------- */

func (s SqlQueryString) String() string {
	return string(s)
}

func (s SqlCondition) String() string {
	return string(s)
}

func (s SqlOperation) String() string {
	return string(s)
}

func (s SqlOrderDirection) String() string {
	return string(s)
}

func (s SqlTable) String() string {
	return s.Name
}

func (s SqlName) String() string {
	return string(s)
}

/* Alias Builders --------------------------------------------------------------------------------------------------- */

func (s SqlTable) As(alias string) SqlTable {
	return SqlTable{Name: s.String() + " AS " + alias}
}

func (s SqlTable) Aliased(alias string) SqlTable {
	return SqlTable{Name: s.String() + " " + alias}
}

func (s SqlTable) Namespaced(namespace string) SqlTable {
	return SqlTable{Name: namespace + "." + s.String()}
}

func (s SqlName) As(alias SqlName) SqlName {
	return SqlName(s.String() + " AS " + alias.String())
}

func (s SqlName) Aliased(alias string) SqlName {
	return SqlName(s.String() + " " + alias)
}

func (s SqlName) Append(after SqlName, with ...string) SqlName {
	var connector = " "

	if len(with) > 0 {
		connector = with[0]
	}

	return SqlName(s.String() + connector + after.String())
}

func (s SqlName) Prepend(before SqlName, with ...string) SqlName {
	var connector = " "

	if len(with) > 0 {
		connector = with[0]
	}

	return SqlName(before.String() + connector + s.String())
}

/* Namespace Builders ----------------------------------------------------------------------------------------------- */

func (s SqlTable) NamespacedColumn(column SqlName) SqlName {
	return SqlName(s.String() + "." + column.String())
}

func (s SqlTable) NamespacedColumns(columns ...SqlName) []SqlName {
	var _columns []SqlName

	for _, column := range columns {
		_columns = append(_columns, s.NamespacedColumn(column))
	}

	return _columns
}

func (s SqlName) NamespacedWith(namespace EmbedsSqlTable) SqlName {
	if namespace == nil || namespace.GetSqlTable().String() == "" {
		return s
	}

	return SqlName(namespace.GetSqlTable().String() + "." + s.String())
}

func (s SqlName) NamespacedWithCustom(namespace SqlName) SqlName {
	if namespace == "" {
		return s
	}

	return SqlName(namespace.String() + "." + s.String())
}
