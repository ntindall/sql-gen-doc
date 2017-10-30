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
	fieldLen:   len("Field"),
	typeLen:    len("Type"),
	nullLen:    len("Null"),
	keyLen:     len("Key"),
	defaultLen: len("Default"),
	extraLen:   len("Extra"),
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
		fieldLen:   len("Field"),
		typeLen:    len("Type"),
		nullLen:    len("Null"),
		keyLen:     len("Key"),
		defaultLen: len("Default"),
		extraLen:   len("Extra"),
	}

	testcases := []struct {
		desc        string
		cd          columnDescription
		expectation string
	}{
		{
			desc:        "pads all fields to specified length (works with empty strings)",
			cd:          columnDescription{},
			expectation: "|       |      |      |     | NULL    |       |\n",
		},
		{
			desc: "pads all fields to specified length",
			cd: columnDescription{
				Field: "Field", // same length
				Extra: "s",     //shorter
			},
			expectation: "| Field |      |      |     | NULL    | s     |\n",
		},
		{
			desc: "writes values to the mark down",
			cd: columnDescription{
				Field:   "b", // same length
				Type:    "a",
				Null:    "n",
				Key:     "a",
				Default: sql.NullString{String: "n", Valid: true},
				Extra:   "a",
			},
			expectation: "| b     | a    | n    | a   | n       | a     |\n",
		},
	}

	for i, tc := range testcases {
		t.Logf("test case %d: %s", i, tc.desc)
		actual := tc.cd.Format(testingFormatSpec)
		assert.Equal(t, tc.expectation, actual)
	}
}

func TestGetFormatSpec(t *testing.T) {
	testcases := []struct {
		desc        string
		cds         []columnDescription
		expectation formatSpec
	}{
		{
			desc:        "returns default for empty cds",
			cds:         []columnDescription{},
			expectation: defaultFormatSpec,
		},
		{
			desc: "returns default if all entries have field less than default",
			cds: []columnDescription{
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
			cds: []columnDescription{
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
				fieldLen:   22,
				typeLen:    20,
				nullLen:    18,
				keyLen:     16,
				defaultLen: 14,
				extraLen:   12,
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
func TestFormatTable(t *testing.T) {
	testcases := []struct {
		desc        string
		tablename   string
		cds         []columnDescription
		expectation string
	}{
		{
			desc:      "works with just one column",
			tablename: "simple_table",
			cds: []columnDescription{
				{
					Field: "id",
					Type:  "bigint(20) unsigned",
					Null:  "NO",
					Key:   "PRI",
					Extra: "PRIMARY KEY",
				},
			},
			expectation: `### simple_table
| Field | Type                | Null | Key | Default | Extra       |
|------------------------------------------------------------------|
| id    | bigint(20) unsigned | NO   | PRI | NULL    | PRIMARY KEY |
`,
		},
		{
			desc:      "works with more complicated table",
			tablename: "complex_table",
			cds: []columnDescription{
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
			expectation: `### complex_table
| Field          | Type                | Null | Key | Default              | Extra       |
|----------------------------------------------------------------------------------------|
| id             | bigint(20) unsigned | NO   | PRI | NULL                 | PRIMARY KEY |
| created        | timestamp(6)        | NO   |     | CURRENT_TIMESTAMP(6) |             |
| indexed_column | bigint(20) unsigned | NO   | MUL | NULL                 |             |
| request_id     | varchar(255)        | YES  |     | NULL                 |             |
`,
		},
	}

	for i, tc := range testcases {
		t.Logf("test case %d: %s", i, tc.desc)
		actual := FormatTable(tc.tablename, tc.cds)
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
