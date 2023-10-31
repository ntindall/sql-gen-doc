package format

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/stretchr/testify/assert"
)

func TestInsertBetweenTags(t *testing.T) {
	testcases := []struct {
		desc        string
		file        string
		markdown    string
		expectation string
		expectedErr string
	}{
		{
			desc:        "works with empty file",
			file:        "",
			markdown:    "markdown",
			expectation: "markdown",
		},
		{
			desc:        "errors when closing tag is missing",
			file:        "<!-- sql-gen-doc BEGIN -->",
			markdown:    "markdown",
			expectation: "markdown",
			expectedErr: "missing end tag <!-- sql-gen-doc END -->",
		},
		{
			desc:        "errors when begin tag is missing",
			file:        "<!-- sql-gen-doc END -->",
			markdown:    "markdown",
			expectedErr: "missing start tag <!-- sql-gen-doc BEGIN -->",
		},
		{
			desc:        "errors when begin tag is after end tag",
			file:        "<!-- sql-gen-doc END --><!-- sql-gen-doc BEGIN -->",
			markdown:    "markdown",
			expectedErr: "tags out of order! <!-- sql-gen-doc BEGIN --> was after <!-- sql-gen-doc END -->",
		},
		{
			desc: "inserts between tags with valid BEGIN and END tags",
			file: `
## hello world

<!-- sql-gen-doc BEGIN -->
some old stuff
<!-- sql-gen-doc END -->

# more stuff to follow!
			`,
			markdown: "markdown",
			expectation: `
## hello world

<!-- sql-gen-doc BEGIN -->
markdown
<!-- sql-gen-doc END -->

# more stuff to follow!
			`,
		},
	}

	for i, tc := range testcases {
		t.Logf("test case %d: %s", i, tc.desc)
		actual, err := insertBetweenTags(tc.file, tc.markdown)
		if tc.expectedErr != "" {
			assert.EqualError(t, err, tc.expectedErr)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expectation, actual)
		}
	}
}

func TestRenderForeignKeys(t *testing.T) {

	testForeignKeys := []ForeignDescription{
		{
			ConstraintName:       "fk1",
			TableName:            "table1",
			ColumnName:           "col1",
			ReferencedTableName:  "reftable",
			ReferencedColumnName: "refcol",
		},
	}

	expected := `
| KEY NAME | TABLE NAME | COLUMN NAME |    REFERENCES     |
|----------|------------|-------------|-------------------|
| ` + "`fk1`   " + ` | ` + "`table1`  " + ` | ` + "`col1`     " + ` | ` + "`reftable.refcol`" + ` |
`

	buf := &bytes.Buffer{}

	foreignKeyTable := tablewriter.NewWriter(buf)
	formatTable(foreignKeyTable, []string{"KEY NAME", "TABLE NAME", "COLUMN NAME", "REFERENCES"})

	for _, fk := range testForeignKeys {
		foreignKeyTable.Append([]string{
			wrapBackTicks(fk.ConstraintName),
			wrapBackTicks(fk.TableName),
			wrapBackTicks(fk.ColumnName),
			wrapBackTicks(fmt.Sprintf("%s.%s", fk.ReferencedTableName, fk.ReferencedColumnName)),
		})
	}

	foreignKeyTable.Render()

	actual := buf.String()

	assert.Equal(t, strings.TrimLeft(expected, "\n"), strings.TrimLeft(actual, "\n"))

}
