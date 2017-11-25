package format

import (
	"bytes"
	"database/sql"
	"html/template"
	"strings"
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
	FieldLen   int
	TypeLen    int
	NullLen    int
	KeyLen     int
	DefaultLen int
	ExtraLen   int
}

func monospaceIfNotEmpty(
	key string,
) string {
	return "{{if " + key + "}}`{{" + key + "}}`{{end}}"
}

func padRemainingWidth(
	s string,
	width int,
) string {
	val := strings.Repeat(" ", width-len(s))

	// When the string is empty, we pad an _additional_ two characters to
	// accomodate the missing backticks
	if s == "" {
		return val + "  "
	}

	return val
}

func (c ColumnDescription) format(f formatSpec) string {
	fieldTpl := monospaceIfNotEmpty(".Field")
	typeTpl := monospaceIfNotEmpty(".Type")
	nullTpl := monospaceIfNotEmpty(".Null")
	keyTpl := monospaceIfNotEmpty(".Key")
	defaultTpl := "{{if .Default.Valid}}`{{.Default.String}}`{{else}}`NULL`{{end}}"
	extraTpl := monospaceIfNotEmpty(".Extra")

	defaultStr := "NULL"
	if c.Default.Valid {
		defaultStr = c.Default.String
	}

	tplString := "| " + fieldTpl + padRemainingWidth(c.Field, f.FieldLen) +
		" | " + typeTpl + padRemainingWidth(c.Type, f.TypeLen) +
		" | " + nullTpl + padRemainingWidth(c.Null, f.NullLen) +
		" | " + keyTpl + padRemainingWidth(c.Key, f.KeyLen) +
		" | " + defaultTpl + padRemainingWidth(defaultStr, f.DefaultLen) +
		" | " + extraTpl + padRemainingWidth(c.Extra, f.ExtraLen) +
		" |\n"

	t := template.Must(template.New("column_description").Parse(tplString))
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, c); err != nil {
		panic(err)
	}

	return tpl.String()

}
