package format

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Fixture used for testing the format spec logic.
var defaultFormatSpec = formatSpec{
	FieldLen:   len("Field"),
	TypeLen:    len("Type"),
	NullLen:    len("Null"),
	KeyLen:     len("Key"),
	DefaultLen: len("Default"),
	ExtraLen:   len("Extra"),
}

func TestPadRemainingWidth(t *testing.T) {
	testcases := []struct {
		desc        string
		inputString string
		inputWidth  int
		expectation string
		shouldPanic bool
	}{
		{
			desc:        "should pad width - len(s)",
			inputString: "hello world", // 11 characters
			inputWidth:  13,
			expectation: "  ",
		},
		{
			desc:        "returns nothing if width == len(s)",
			inputString: "hello world", // 11 characters
			inputWidth:  11,
			expectation: "",
		},
		{
			desc:        "panicks if width < len(s)",
			inputString: "hello world", // 11 characters
			inputWidth:  10,
			shouldPanic: true,
		},
	}

	for i, tc := range testcases {
		func() {
			t.Logf("test case %d: %s", i, tc.desc)

			defer func() {
				if p := recover(); p != nil {
					require.True(t, tc.shouldPanic, fmt.Sprintf("unexpected recovered from panic: %v", p))
					assert.Equal(t, p, "strings: negative Repeat count")
				}
			}()

			actual := padRemainingWidth(tc.inputString, tc.inputWidth)
			assert.Equal(t, tc.expectation, actual)
		}()
	}
}

func TestColumnDescriptionFormat(t *testing.T) {
	testingFormatSpec := formatSpec{
		FieldLen:   len("Field"),
		TypeLen:    len("Type"),
		NullLen:    len("Null"),
		KeyLen:     len("Key"),
		DefaultLen: len("Default"),
		ExtraLen:   len("Extra"),
	}

	testcases := []struct {
		desc        string
		cd          ColumnDescription
		expectation string
	}{
		{
			desc:        "pads all fields to specified length (works with empty strings)",
			cd:          ColumnDescription{},
			expectation: "|       |      |      |     | `NULL`    |       |\n",
		},
		{
			desc: "pads all fields to specified length",
			cd: ColumnDescription{
				Field: "Field", // same length
				Extra: "s",     //shorter
			},
			expectation: "| `Field` |      |      |     | `NULL`    | `s`     |\n",
		},
		{
			desc: "writes values to the mark down",
			cd: ColumnDescription{
				Field:   "b", // same length
				Type:    "a",
				Null:    "n",
				Key:     "a",
				Default: sql.NullString{String: "n", Valid: true},
				Extra:   "a",
			},
			expectation: "| `b`     | `a`    | `n`    | `a`   | `n`       | `a`     |\n",
		},
	}

	for i, tc := range testcases {
		t.Logf("test case %d: %s", i, tc.desc)
		actual := tc.cd.format(testingFormatSpec)
		assert.Equal(t, tc.expectation, actual)
	}
}

func TestGetFormatSpec(t *testing.T) {
	testcases := []struct {
		desc        string
		cds         []ColumnDescription
		expectation formatSpec
	}{
		{
			desc:        "returns default for empty cds",
			cds:         []ColumnDescription{},
			expectation: defaultFormatSpec,
		},
		{
			desc: "returns default if all entries have field less than default",
			cds: []ColumnDescription{
				{
					Field:   "aa",
					Type:    "aa",
					Null:    "aa",
					Key:     "aa",
					Default: sql.NullString{"aa", true},
					Extra:   "aa",
				},
				{
					Field:   "bb",
					Type:    "bb",
					Null:    "bb",
					Key:     "bb",
					Default: sql.NullString{"bb", true},
					Extra:   "bb",
				},
			},
			expectation: defaultFormatSpec,
		},
		{
			desc: "returns longest value for each field name",
			cds: []ColumnDescription{
				{
					Field: "the longest field name", // 22
				},
				{
					Type: "the longest field na", // 20
				},
				{
					Null: "the longest field ", // 18
				},
				{
					Key: "the longest fiel", // 16
				},
				{
					Default: sql.NullString{"the longest fi", true}, // 14
				},
				{
					Extra: "the longest ", // 12
				},
			},
			expectation: formatSpec{
				// we add two to the length of the string to accomodate the backtick.
				FieldLen:   22,
				TypeLen:    20,
				NullLen:    18,
				KeyLen:     16,
				DefaultLen: 14,
				ExtraLen:   12,
			},
		},
	}

	for i, tc := range testcases {
		t.Logf("test case %d: %s", i, tc.desc)
		actual := getFormatSpec(tc.cds)
		assert.Equal(t, tc.expectation, actual)
	}
}

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
				"| `Field` | `Type`                | `Null` | `Key` | `Default` | `Extra`       |\n" +
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
				"| `Field`            | `Type`                  | `Null`   | `Key`   | `Default`              | `Extra`       |\n" +
				"|--------------------|-------------------------|----------|---------|------------------------|---------------|\n" +
				"| `id`               | `bigint(20) unsigned`   | `NO`     | `PRI`   | `NULL`                 | `PRIMARY KEY` |\n" +
				"| `created`          | `timestamp(6)`          | `NO`     |         | `CURRENT_TIMESTAMP(6)` |               |\n" +
				"| `indexed_column`   | `bigint(20) unsigned`   | `NO`     | `MUL`   | `NULL`                 |               |\n" +
				"| `request_id`       | `varchar(255)`          | `YES`    |         | `NULL`                 |               |\n",
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
