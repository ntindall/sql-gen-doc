package format

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

// table abstraction -> markdown conversion
func TestCreateTableMarkdown(t *testing.T) {
	testcases := []struct {
		desc        string
		tablename   string
		cds         []ColumnDescription
		expectation string
	}{
		{
			desc:      "works with just one column",
			tablename: "simple_table",
			cds: []ColumnDescription{
				{
					Field: "id",
					Type:  "bigint(20) unsigned",
					Null:  "NO",
					Key:   "PRI",
					Extra: "PRIMARY KEY",
				},
			},
			expectation: "### simple_table\n" +
				"| Field   | Type                  | Null   | Key   | Default   | Extra         |\n" +
				"|---------|-----------------------|--------|-------|-----------|---------------|\n" +
				"| `id`    | `bigint(20) unsigned` | `NO`   | `PRI` | `NULL`    | `PRIMARY KEY` |\n",
		},
		{
			desc:      "works with more complicated table",
			tablename: "complex_table",
			cds: []ColumnDescription{
				{
					Field: "id",
					Type:  "bigint(20) unsigned",
					Null:  "NO",
					Key:   "PRI",
					Extra: "PRIMARY KEY",
				},
				{
					Field:   "created",
					Type:    "timestamp(6)",
					Null:    "NO",
					Default: sql.NullString{"CURRENT_TIMESTAMP(6)", true},
				},
				{
					Field: "indexed_column",
					Type:  "bigint(20) unsigned",
					Null:  "NO",
					Key:   "MUL",
				},
				{
					Field: "request_id",
					Type:  "varchar(255)",
					Null:  "YES",
				},
			},
			expectation: "### complex_table\n" +
				"| Field            | Type                  | Null   | Key   | Default                | Extra         |\n" +
				"|------------------|-----------------------|--------|-------|------------------------|---------------|\n" +
				"| `id`             | `bigint(20) unsigned` | `NO`   | `PRI` | `NULL`                 | `PRIMARY KEY` |\n" +
				"| `created`        | `timestamp(6)`        | `NO`   |       | `CURRENT_TIMESTAMP(6)` |               |\n" +
				"| `indexed_column` | `bigint(20) unsigned` | `NO`   | `MUL` | `NULL`                 |               |\n" +
				"| `request_id`     | `varchar(255)`        | `YES`  |       | `NULL`                 |               |\n",
		},
	}

	for i, tc := range testcases {
		t.Logf("test case %d: %s", i, tc.desc)
		actual := CreateTableMarkdown(tc.tablename, tc.cds)
		assert.Equal(t, tc.expectation, actual)
	}
}

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

func TestMakeHeader(t *testing.T) {
	testcases := []struct {
		desc        string
		f           formatSpec
		expectation string
	}{
		{
			desc: "works with default format spec",
			f:    defaultFormatSpec,
			expectation: "" +
				"| Field   | Type   | Null   | Key   | Default   | Extra   |\n" +
				"|---------|--------|--------|-------|-----------|---------|\n",
		},
		{
			desc: "pads column names and dashes to correct length",
			f: formatSpec{
				FieldLen:   10,
				TypeLen:    10,
				NullLen:    20,
				KeyLen:     10,
				DefaultLen: 10,
				ExtraLen:   15,
			},
			expectation: "" +
				"| Field        | Type         | Null                   | Key          | Default      | Extra             |\n" +
				"|--------------|--------------|------------------------|--------------|--------------|-------------------|\n",
		},
	}

	for i, tc := range testcases {
		t.Logf("test case %d: %s", i, tc.desc)
		actual := makeHeader(tc.f)
		assert.Equal(t, tc.expectation, actual)
	}
}
