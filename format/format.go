package format

import (
	"database/sql"
	"fmt"
	"strings"
)

type columnDescription struct {
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

func padRemainingWidth(
	s string,
	width int,
) string {
	return strings.Repeat(" ", width-len(s))
}

func (c columnDescription) Format(f formatSpec) string {
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

func makeHeader(f formatSpec) string {
	numberOfDashes :=
		f.defaultLen + 2 + // padding for the surrounding spaces
			f.extraLen + 2 +
			f.fieldLen + 2 +
			f.keyLen + 2 +
			f.nullLen + 2 +
			f.typeLen + 2 +
			5 // one for each interior column separator

	return "| Field" + padRemainingWidth("Field", f.fieldLen) +
		" | Type" + padRemainingWidth("Type", f.typeLen) +
		" | Null" + padRemainingWidth("Null", f.nullLen) +
		" | Key" + padRemainingWidth("Key", f.keyLen) +
		" | Default" + padRemainingWidth("Default", f.defaultLen) +
		" | Extra" + padRemainingWidth("Extra", f.extraLen) +
		" |\n" +
		"|" + strings.Repeat("-", numberOfDashes) + "|" +
		"\n"
}

func getFormatSpec(columns []columnDescription) formatSpec {
	spec := formatSpec{
		fieldLen:   len("Field"),
		typeLen:    len("Type"),
		nullLen:    len("Null"),
		keyLen:     len("Key"),
		defaultLen: len("Default"),
		extraLen:   len("Extra"),
	}

	for _, c := range columns {
		if len(c.Field) > spec.fieldLen {
			spec.fieldLen = len(c.Field)
		}

		if len(c.Type) > spec.typeLen {
			spec.typeLen = len(c.Type)
		}

		if len(c.Null) > spec.nullLen {
			spec.nullLen = len(c.Null)
		}

		if len(c.Key) > spec.keyLen {
			spec.keyLen = len(c.Key)
		}

		if c.Default.Valid && len(c.Default.String) > spec.defaultLen {
			spec.defaultLen = len(c.Default.String)
		}

		if len(c.Extra) > spec.extraLen {
			spec.extraLen = len(c.Extra)
		}
	}

	return spec
}

func makeTitle(s string) string {
	return "### " + s + "\n"
}

func FormatTable(
	table string,
	columns []columnDescription,
) string {
	tableMarkdown := makeTitle(table)

	formatSpec := getFormatSpec(columns)
	tableMarkdown = tableMarkdown + makeHeader(formatSpec)

	for _, c := range columns {
		tableMarkdown += c.Format(formatSpec)
	}

	return tableMarkdown
}

func insertBetweenTags(
	file string,
	markdown string,
) (string, error) {
	startTag := "<!-- sql-gen-doc BEGIN -->"
	endTag := "<!-- sql-gen-doc END -->"

	startIdx := strings.Index(file, startTag)
	endIdx := strings.Index(file, endTag)

	if startIdx == -1 && endIdx != -1 {
		return "", fmt.Errorf("missing start tag <!-- sql-gen-doc BEGIN -->")
	} else if startIdx != -1 && endIdx == -1 {
		return "", fmt.Errorf("missing end tag <!-- sql-gen-doc END -->")
	} else if startIdx == -1 && endIdx == -1 {
		return markdown, nil
	} else if startIdx > endIdx {
		return "", fmt.Errorf("tags out of order! <!-- sql-gen-doc BEGIN --> was after <!-- sql-gen-doc END -->")
	}

	// all is well, insert between the tags!
	startIdx += len(startTag)
	return file[:startIdx] + "\n" + markdown + "\n" + file[endIdx:], nil
}
