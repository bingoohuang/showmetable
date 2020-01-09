package model

import "database/sql"

// Table abstract a table information
type Table interface {
	// GetName gets the table name
	GetName() string
	// GetComment gets the table comment
	GetComment() string
}

// TableColumn describes a column in a table
type TableColumn interface {
	// GetName get the column's name
	GetName() string
	// IsNullable tells if the column is nullable or not
	IsNullable() bool
	// GetComment get the column's comment
	GetComment() string
	// GetDefault get the columns's default value
	GetDefault() string
	// GetDataType get the columns's data type
	GetDataType() string
	// GetMaxSize get the max size of the column
	GetMaxSize() sql.NullInt64
}

// DbSchema ...
type DbSchema interface {
	// Tables get all the tables
	Tables() ([]Table, error)
	// TableColumns gets the all columns by table name
	TableColumns(table string) ([]TableColumn, error)
	// CompatibleDs returns the dataSourceName from various the compatible format.
	CompatibleDs() string
}
