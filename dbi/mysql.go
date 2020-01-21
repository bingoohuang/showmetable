package dbi

import (
	"database/sql"

	"github.com/bingoohuang/gou/str"

	"strings"

	"github.com/bingoohuang/pump/ds"
	"github.com/bingoohuang/showmetable/model"
	"github.com/bingoohuang/sqlmore"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// MySQLTable ...
type MySQLTable struct {
	Schema  string `gorm:"column:TABLE_SCHEMA"`
	Name    string `gorm:"column:TABLE_NAME"`
	Comment string `gorm:"column:TABLE_COMMENT"`
}

var _ model.Table = (*MySQLTable)(nil)

// GetScheme gets the table scheme
func (m MySQLTable) GetScheme() string { return m.Schema }

// GetName ...
func (m MySQLTable) GetName() string { return m.Name }

// GetComment  ...
func (m MySQLTable) GetComment() string { return m.Comment }

// MyTableColumn ...
type MyTableColumn struct {
	Name      string         `gorm:"column:COLUMN_NAME"`
	Type      string         `gorm:"column:COLUMN_TYPE"`
	Extra     string         `gorm:"column:EXTRA"` // auto_increment
	Comment   string         `gorm:"column:COLUMN_COMMENT"`
	DataType  string         `gorm:"column:DATA_TYPE"`
	MaxLength sql.NullInt64  `gorm:"column:CHARACTER_MAXIMUM_LENGTH"`
	Nullable  string         `gorm:"column:IS_NULLABLE"`
	Default   sql.NullString `gorm:"column:COLUMN_DEFAULT"`

	NumericPrecision sql.NullInt64 `gorm:"column:NUMERIC_PRECISION"`
	NumericScale     sql.NullInt64 `gorm:"column:NUMERIC_SCALE"`
}

var _ model.TableColumn = (*MyTableColumn)(nil)

// IsNullable ...
func (c MyTableColumn) IsNullable() bool { return c.Nullable == "YES" }

// GetMaxSize ...
func (c MyTableColumn) GetMaxSize() sql.NullInt64 { return c.MaxLength }

// GetDataType ...
func (c MyTableColumn) GetDataType() string { return c.Type }

// GetName ...
func (c MyTableColumn) GetName() string { return c.Name }

// GetComment ...
func (c MyTableColumn) GetComment() string { return c.Comment }

// GetDefault ...
func (c MyTableColumn) GetDefault() string {
	d := ""

	if c.Default.Valid {
		d = c.Default.String
	}

	if c.DataType == "double" {
		d = TrimNumber(d)
	}

	return d
}

// TrimNumber trim zeros from s.
func TrimNumber(s string) string {
	for {
		if !strings.HasSuffix(s, "0") {
			break
		}

		s = s[:len(s)-1]
	}

	if strings.HasSuffix(s, ".") {
		s = s[:len(s)-1]
	}

	return s
}

// MySQLSchema ...
type MySQLSchema struct {
	dbFn    func() (*gorm.DB, error)
	verbose bool
	db      *gorm.DB
}

// Tables get all the tables
func (s *MySQLSchema) Tables() ([]model.Table, error) {
	var err error
	s.db, err = s.dbFn()

	if err != nil {
		return nil, err
	}

	//defer gouio.Close(db)

	var tables []MySQLTable

	const sq = `SELECT * FROM information_schema.TABLES
		WHERE TABLE_SCHEMA NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')`

	s.db.Raw(sq).Find(&tables)

	ts := make([]model.Table, len(tables))

	for i, t := range tables {
		//t.Comment = strings.TrimSpace(t.Comment)
		ts[i] = t
	}

	return ts, s.db.Error
}

// TableColumns gets the all columns by table name
func (s *MySQLSchema) TableColumns(table string) ([]model.TableColumn, error) {
	columns := make([]MyTableColumn, 0)
	schema, tableName := ParseTable(table)

	if schema != "" {
		const q = `SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ?
			AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION`

		s.db.Raw(q, schema, tableName).Find(&columns)
	} else {
		const q = `SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = database()
			AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION`

		s.db.Raw(q, tableName).Find(&columns)
	}

	ts := make([]model.TableColumn, len(columns))

	for i, t := range columns {
		t.Comment = strings.TrimSpace(t.Comment)
		ts[i] = t
	}

	return ts, s.db.Error
}

// ParseTable parses the schema and table name from table which may be like db1.t1
func ParseTable(table string) (schemaName, tableName string) {
	if strings.Contains(table, ".") {
		return str.Split2(table, ".", true, true)
	}

	return "", table
}

var _ model.DbSchema = (*MySQLSchema)(nil)

// SetVerbose sets verbose
func (s *MySQLSchema) SetVerbose(verbose bool) { s.verbose = verbose }

// CreateMySQLSchema ...
func CreateMySQLSchema(dataSourceName string) (*MySQLSchema, error) {
	compatibleDs := ds.CompatibleMySQLDs(dataSourceName)
	logrus.Infof("dataSourceName:%s", compatibleDs)

	return &MySQLSchema{
		dbFn: func() (*gorm.DB, error) { return sqlmore.NewSQLMore("mysql", compatibleDs).GormOpen() },
	}, nil
}
