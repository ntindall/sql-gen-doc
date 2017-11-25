package format

import (
	"fmt"
	"strings"
)

func makeHeader(f formatSpec) string {
	// TODO: a better abstraction for styling column description fields.
	// we add two to the expected length to accomodate the missing backticks.
	header := "| " + strings.Join([]string{
		"Field" + padRemainingWidth("Field", f.FieldLen+2),
		"Type" + padRemainingWidth("Type", f.TypeLen+2),
		"Null" + padRemainingWidth("Null", f.NullLen+2),
		"Key" + padRemainingWidth("Key", f.KeyLen+2),
		"Default" + padRemainingWidth("Default", f.DefaultLen+2),
		"Extra" + padRemainingWidth("Extra", f.ExtraLen+2),
	}, " | ") + " |\n"

	dashes := "|" + strings.Join([]string{
		// we add two for each backtick present around the column name, plus an
		// additional two for the spacing on either side.
		strings.Repeat("-", f.FieldLen+4),
		strings.Repeat("-", f.TypeLen+4),
		strings.Repeat("-", f.NullLen+4),
		strings.Repeat("-", f.KeyLen+4),
		strings.Repeat("-", f.DefaultLen+4),
		strings.Repeat("-", f.ExtraLen+4),
	}, "|") + "|" + "\n"

	return header + dashes
}

func getFormatSpec(columns []ColumnDescription) formatSpec {
	spec := formatSpec{
		FieldLen:   len("Field"),
		TypeLen:    len("Type"),
		NullLen:    len("Null"),
		KeyLen:     len("Key"),
		DefaultLen: len("Default"),
		ExtraLen:   len("Extra"),
	}

	// Iterate over each column
	for _, c := range columns {
		if len(c.Field) > spec.FieldLen {
			spec.FieldLen = len(c.Field)
		}

		if len(c.Type) > spec.TypeLen {
			spec.TypeLen = len(c.Type)
		}

		if len(c.Null) > spec.NullLen {
			spec.NullLen = len(c.Null)
		}

		if len(c.Key) > spec.KeyLen {
			spec.KeyLen = len(c.Key)
		}

		if c.Default.Valid && len(c.Default.String) > spec.DefaultLen {
			spec.DefaultLen = len(c.Default.String)
		}

		if len(c.Extra) > spec.ExtraLen {
			spec.ExtraLen = len(c.Extra)
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
