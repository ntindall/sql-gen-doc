package format

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func makeTitle(s string) string {
	return "### " + s + "\n"
}

func wrapBackTicks(s string) string {
	if s == "" {
		return s
	}
	return "`" + s + "`"
}

// CreateTableMarkdown takes the name of a table in a database and a list of
// ColumnDescription and returns a formatted markdown table with the
// corresponding data.
func CreateTableMarkdown(
	tableName string,
	columns []ColumnDescription,
	indexes []LogicalIndex,
) string {
	tableMarkdown := bytes.NewBufferString(`### ` + tableName + "\n")

	tableMarkdown.WriteString("#### SCHEMA\n")
	columnsTable := tablewriter.NewWriter(tableMarkdown)
	columnsTable.SetHeader([]string{"FIELD", "TYPE", "NULL", "KEY", "DEFAULT", "EXTRA"})
	columnsTable.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: true})
	columnsTable.SetCenterSeparator("|")

	for _, col := range columns {
		columnsTable.Append([]string{wrapBackTicks(col.Field), wrapBackTicks(col.Type), wrapBackTicks(col.Null), wrapBackTicks(col.Key), wrapBackTicks(col.Default.String), wrapBackTicks(col.Extra)})
	}

	// write the columns table to the buf
	columnsTable.Render()
	tableMarkdown.WriteString("\n")

	// format the indexes
	tableMarkdown.WriteString("#### INDEXES\n")
	indexesTable := tablewriter.NewWriter(tableMarkdown)
	indexesTable.SetHeader([]string{"KEY NAME", "NON UNIQUE", "COLUMNS"})
	indexesTable.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: true})
	indexesTable.SetCenterSeparator("|")

	for _, idx := range indexes {
		indexesTable.Append([]string{
			wrapBackTicks(idx.KeyName),
			wrapBackTicks(fmt.Sprintf("%t", idx.NonUnique)),
			wrapBackTicks(fmt.Sprintf(`(%s)`, strings.Join(idx.IndexedColumnNamesOrdered, ", "))),
		})
	}

	// write the columns table to the buf
	indexesTable.Render()
	tableMarkdown.WriteString("\n")

	return tableMarkdown.String()
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
