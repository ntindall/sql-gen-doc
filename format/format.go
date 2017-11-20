package format

import (
	"fmt"
	"strings"
)

func padRemainingWidth(
	s string,
	width int,
) string {
	return strings.Repeat(" ", width-len(s))
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

func getFormatSpec(columns []ColumnDescription) formatSpec {
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

// CreateTableMarkdown takes the name of a table in a database and a list of
// ColumnDescription and returns a formatted markdown table with the
// corresponding data.
func CreateTableMarkdown(
	table string,
	columns []ColumnDescription,
) string {
	tableMarkdown := makeTitle(table)

	formatSpec := getFormatSpec(columns)
	tableMarkdown = tableMarkdown + makeHeader(formatSpec)

	for _, c := range columns {
		tableMarkdown += c.format(formatSpec)
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
