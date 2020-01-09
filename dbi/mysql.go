package dbi

import (
	"database/sql"

	gouio "github.com/bingoohuang/gou/io"
	"github.com/bingoohuang/gou/str"

	"strings"

	"github.com/bingoohuang/pump/ds"
	"github.com/bingoohuang/showmetable/model"
	"github.com/bingoohuang/sqlmore"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

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
func (c MyTableColumn) GetDataType() string { return c.DataType }

// GetName ...
func (c MyTableColumn) GetName() string { return c.Name }

// GetComment ...
func (c MyTableColumn) GetComment() string { return c.Comment }

// GetDefault ...
func (c MyTableColumn) GetDefault() string {
	if c.Default.Valid {
		return c.Default.String
	}

	return ""
}

// MySQLSchema ...
type MySQLSchema struct {
	dbFn         func() (*gorm.DB, error)
	compatibleDs string

	verbose       bool
	showAllTables bool
	showTables    []string
}

// Tables get all the tables
func (s *MySQLSchema) Tables() ([]model.Table, error) {
	panic("implement me")
}

// TableColumns gets the all columns by table name
func (s *MySQLSchema) TableColumns(table string) ([]model.TableColumn, error) {
	db, err := s.dbFn()
	if err != nil {
		return nil, err
	}

	defer gouio.Close(db)

	columns := make([]MyTableColumn, 0)
	schema, tableName := ParseTable(table)

	if schema != "" {
		const s = `SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ` +
			`AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION`

		db.Raw(s, schema, tableName).Find(&columns)
	} else {
		const s = `SELECT * FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = database() ` +
			`AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION`

		db.Raw(s, tableName).Find(&columns)
	}

	ts := make([]model.TableColumn, len(columns))

	for i, t := range columns {
		t.Comment = strings.TrimSpace(t.Comment)
		ts[i] = t
	}

	return ts, db.Error
}

// ParseTable parses the schema and table name from table which may be like db1.t1
func ParseTable(table string) (schemaName, tableName string) {
	if strings.Contains(table, ".") {
		return str.Split2(table, ".", true, true)
	}

	return "", table
}

// CompatibleDs returns the dataSourceName from various the compatible format.
func (s *MySQLSchema) CompatibleDs() string { return s.compatibleDs }

var _ model.DbSchema = (*MySQLSchema)(nil)

// SetVerbose sets verbose
func (s *MySQLSchema) SetVerbose(verbose bool) {
	s.verbose = verbose
}

// SetShowTables set which tables to show
func (s *MySQLSchema) SetShowTables(tables []string) {
	if len(tables) == 0 {
		s.showAllTables = true
		return
	}

	s.showTables = tables
}

// CreateMySQLSchema ...
func CreateMySQLSchema(dataSourceName string) (*MySQLSchema, error) {
	compatibleDs := ds.CompatibleMySQLDs(dataSourceName)
	logrus.Infof("dataSourceName:%s", compatibleDs)

	return &MySQLSchema{
		dbFn:         func() (*gorm.DB, error) { return sqlmore.NewSQLMore("mysql", compatibleDs).GormOpen() },
		compatibleDs: compatibleDs,
	}, nil
}
