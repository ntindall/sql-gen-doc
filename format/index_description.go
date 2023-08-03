package format

import (
	"database/sql"
	"fmt"
)

// IndexDescription contains all the data known about a specific index.
// Note that some indexes may be related (e.g. in cases of clustered indexes).
type IndexDescription struct {
	Table      string         `db:"Table"`
	NonUnique  bool           `db:"Non_unique"`
	KeyName    string         `db:"Key_name"`
	SeqInIndex int            `db:"Seq_in_index"`
	ColumnName string         `db:"Column_name"`
	Comment    sql.NullString `db:"Comment"`

	// Not used (yet)
	Collation    sql.NullString `db:"Collation"`
	Cardinality  sql.NullString `db:"Cardinality"`
	SubPart      sql.NullString `db:"Sub_part"`
	Packed       sql.NullString `db:"Packed"`
	Null         sql.NullString `db:"Null"`
	IndexType    sql.NullString `db:"Index_type"`
	IndexComment sql.NullString `db:"Index_comment"`
	Visible      sql.NullString `db:"Visible"`
	Expression   sql.NullString `db:"Expression"`
}

// IndexDescriptions is a set of index descriptions
type IndexDescriptions []IndexDescription

// ConvertToLogicalIndexes converts raw index descriptions into a aggregated
// type.
func (descs IndexDescriptions) ConvertToLogicalIndexes() ([]LogicalIndex, error) {
	indices := map[string]LogicalIndex{}

	for _, description := range descs {
		var li LogicalIndex
		var ok bool
		if li, ok = indices[description.KeyName]; !ok {
			li = LogicalIndex{
				Table:     description.Table,
				NonUnique: description.NonUnique,
				KeyName:   description.KeyName,
				Comment:   description.Comment.String,
			}
		}

		li.IndexedColumnNamesOrdered = append(li.IndexedColumnNamesOrdered, description.ColumnName)

		if len(li.IndexedColumnNamesOrdered) != description.SeqInIndex {
			return nil, fmt.Errorf("internal logic error: expecting for indexed columns to always be returned in sequence")
		}

		indices[description.KeyName] = li
	}

	// return things in the order they were received -- required for the out
	// put to be deterministic.
	result := []LogicalIndex{}
	var addedKeys = map[string]struct{}{}
	for _, description := range descs {

		// only add keys once
		if _, ok := addedKeys[description.KeyName]; !ok {
			result = append(result, indices[description.KeyName])
			addedKeys[description.KeyName] = struct{}{}
		}
	}

	return result, nil
}

// LogicalIndex defines a "logical index -- e.g. what will make the most sense
// to humans"
type LogicalIndex struct {
	Table                     string
	NonUnique                 bool
	KeyName                   string
	IndexedColumnNamesOrdered []string
	Comment                   string
}
