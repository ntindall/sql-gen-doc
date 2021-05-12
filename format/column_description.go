package format

import (
	"database/sql"
)

// ColumnDescription contains all the data rendered about a sql column by the
// DESCRIBE command.
type ColumnDescription struct {
	Field      string         `db:"Field"`
	Type       string         `db:"Type"`
	Null       string         `db:"Null"`
	Key        string         `db:"Key"`
	Default    sql.NullString `db:"Default"`
	Extra      string         `db:"Extra"`
	Comment    string         `db:"Comment"`
	Collation  sql.NullString `db:"Collation"`
	Privileges string         `db:"Privileges"`
}
