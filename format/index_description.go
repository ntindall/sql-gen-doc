package format

// IndexDescription contains all the data known about a specific index.
// Note that some indexes may be related (e.g. in cases of clustered indexes).
type IndexDescription struct {
	Table      string `db:"Table"`
	NonUnique  bool   `db:"NonUnique"`
	KeyName    string `db:"Key_name"`
	SeqInIndex string `db:"Seq_in_index"`
	ColumnName string `db:"Column_name"`
}
