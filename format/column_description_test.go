package format

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		f           formatSpec
		expectation string
	}{
		{
			desc:        "pads all fields to specified length (works with empty strings, writes NULL for Default)",
			cd:          ColumnDescription{},
			f:           testingFormatSpec,
			expectation: "|         |        |        |       | `NULL`    |         |\n",
		},
		{
			desc: "pads all fields to specified length",
			cd: ColumnDescription{
				Field: "Field", // same length
				Extra: "s",     //shorter
			},
			f:           testingFormatSpec,
			expectation: "| `Field` |        |        |       | `NULL`    | `s`     |\n",
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
			f:           testingFormatSpec,
			expectation: "| `b`     | `a`    | `n`    | `a`   | `n`       | `a`     |\n",
		},
		{
			desc: "pads appropriately with a custom format spec",
			cd: ColumnDescription{
				Field:   "b", // same length
				Type:    "a",
				Null:    "n",
				Key:     "a",
				Default: sql.NullString{String: "n", Valid: true},
				Extra:   "a",
			},
			f: func() formatSpec {
				// don't modify the global
				cpy := testingFormatSpec
				cpy.FieldLen = 5
				cpy.ExtraLen = 6

				return cpy
			}(),
			expectation: "| `b`     | `a`    | `n`    | `a`   | `n`       | `a`      |\n",
		},
	}

	for i, tc := range testcases {
		t.Logf("test case %d: %s", i, tc.desc)
		actual := tc.cd.format(tc.f)
		assert.Equal(t, tc.expectation, actual)
	}
}
