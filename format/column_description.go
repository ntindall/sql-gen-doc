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
	lenAccessor string,
) string {
	return "{{if " + key + "}}`{{" + key + "}}`{{end}}" +
		"{{pad_remaining_width " + key + " " + lenAccessor + "}}"
}

func padRemainingWidth(
	s string,
	width int,
) string {
	if width-len(s) < 0 {
		return ""
	}
	return strings.Repeat(" ", width-len(s))
}

func (c ColumnDescription) format(f formatSpec) string {
	fieldTpl := monospaceIfNotEmpty(".Field", ".FieldLen")
	typeTpl := monospaceIfNotEmpty(".Type", ".TypeLen")
	nullTpl := monospaceIfNotEmpty(".Null", ".NullLen")
	keyTpl := monospaceIfNotEmpty(".Key", ".KeyLen")
	defaultTpl :=
		"{{if .Default.Valid}}`{{.Default.String}}`{{pad_remaining_width `{{.Default.String}}` .DefaultLen}}" +
			"{{else}}`NULL`{{pad_remaining_width `NULL` .DefaultLen}}{{end}}"
	extraTpl := monospaceIfNotEmpty(".Extra", ".ExtraLen")

	tplString := "| " + strings.Join([]string{
		fieldTpl,
		typeTpl,
		nullTpl,
		keyTpl,
		defaultTpl,
		extraTpl,
	}, " | ") + " |\n"

	t := template.Must(template.New("column_description").Funcs(map[string]interface{}{
		"pad_remaining_width": padRemainingWidth,
	}).Parse(tplString))

	merged := struct {
		ColumnDescription
		formatSpec
	}{
		c,
		f,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, &merged); err != nil {
		panic(err)
	}

	return tpl.String()

}
