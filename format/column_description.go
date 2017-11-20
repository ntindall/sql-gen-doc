package format

import (
	"database/sql"
)

// ColumnDescription contains all the data rendered about a sql column by the
// DESCRIBE command.
type ColumnDescription struct {
	Field   string         `db:"Field"`
	Type    string         `db:"Type"`
	Null    string         `db:"Null"`
	Key     string         `db:"Key"`
	Default sql.NullString `db:"Default"`
	Extra   string         `db:"Extra"`
}

type formatSpec struct {
	fieldLen   int
	typeLen    int
	nullLen    int
	keyLen     int
	defaultLen int
	extraLen   int
}

func (c ColumnDescription) format(f formatSpec) string {
	defaultStr := "NULL"
	if c.Default.Valid {
		defaultStr = c.Default.String
	}

	return "| " + c.Field + padRemainingWidth(c.Field, f.fieldLen) +
		" | " + c.Type + padRemainingWidth(c.Type, f.typeLen) +
		" | " + c.Null + padRemainingWidth(c.Null, f.nullLen) +
		" | " + c.Key + padRemainingWidth(c.Key, f.keyLen) +
		" | " + defaultStr + padRemainingWidth(defaultStr, f.defaultLen) +
		" | " + c.Extra + padRemainingWidth(c.Extra, f.extraLen) +
		" |\n"
}
