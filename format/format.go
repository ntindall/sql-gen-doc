package format

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func wrapBackTicks(s string) string {
	if s == "" {
		return s
	}
	return "`" + s + "`"
}

// CreateTableMarkdown takes the name of a table in a database and a list of
// ColumnDescription and returns a formatted markdown table with the
// corresponding data.
func CreateTableMarkdown(tableName string, comment string, columns []ColumnDescription, indexes []LogicalIndex, foreignKeys ForeignDescriptions) string {
	tableMarkdown := bytes.NewBufferString(`## ` + tableName + "\n")

	if comment != "" {
		tableMarkdown.WriteString("\n" + comment + "\n\n")
	}

	tableMarkdown.WriteString("#### SCHEMA\n")
	columnsTable := tablewriter.NewWriter(tableMarkdown)
	formatTable(columnsTable, []string{"FIELD", "TYPE", "NULL", "KEY", "DEFAULT", "EXTRA", "COMMENT"})

	for _, col := range columns {
		columnsTable.Append([]string{
			wrapBackTicks(col.Field),
			wrapBackTicks(col.Type),
			wrapBackTicks(col.Null),
			wrapBackTicks(col.Key),
			wrapBackTicks(col.Default.String),
			wrapBackTicks(col.Extra),
			wrapBackTicks(col.Comment),
		})
	}

	// write the columns table to the buf
	columnsTable.Render()

	// format the indexes
	tableMarkdown.WriteString("#### INDEXES\n")
	indexesTable := tablewriter.NewWriter(tableMarkdown)

	// EXPRESSION is a new column type introduced in MySQL 8.0.
	// Only include this header if one of the indexes has an expression.
	hasExpression := false
	for _, idx := range indexes {
		if idx.Expression != "" {
			hasExpression = true
			break
		}
	}
	header := []string{"KEY NAME", "UNIQUE", "COLUMNS", "COMMENT"}
	if hasExpression {
		header = []string{"KEY NAME", "UNIQUE", "COLUMNS", "COMMENT", "EXPRESSION"}
	}

	formatTable(indexesTable, header)

	for _, idx := range indexes {
		indexRow := []string{
			wrapBackTicks(idx.KeyName),
			wrapBackTicks(fmt.Sprintf("%t", !idx.NonUnique)),
			wrapBackTicks(fmt.Sprintf(`(%s)`, strings.Join(idx.IndexedColumnNamesOrdered, ", "))),
			wrapBackTicks(idx.Comment),
		}
		// Only include this if one of the indexes in the table has an expression.
		if hasExpression {
			indexRow = append(indexRow, wrapBackTicks(idx.Expression))
		}
		indexesTable.Append(indexRow)
	}

	// write the indexes table to the buf
	indexesTable.Render()

	if len(foreignKeys) != 0 {
		tableMarkdown.WriteString("#### Foreign Key\n")
		foreignKeyTable := tablewriter.NewWriter(tableMarkdown)
		formatTable(foreignKeyTable, []string{"KEY NAME", "TABLE NAME", "COLUMN NAME", "REFERENCES"})

		for _, fk := range foreignKeys {
			foreignKeyTable.Append([]string{
				wrapBackTicks(fk.ConstraintName),
				wrapBackTicks(fk.TableName),
				wrapBackTicks(fk.ColumnName),
				wrapBackTicks(fmt.Sprintf("%s.%s", fk.ReferencedTableName, fk.ReferencedColumnName)),
			})
		}
		foreignKeyTable.Render()
	}

	return tableMarkdown.String()
}

func formatTable(table *tablewriter.Table, header []string) {
	table.SetHeader(header)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetAutoWrapText(false)
	table.SetCenterSeparator("|")
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
