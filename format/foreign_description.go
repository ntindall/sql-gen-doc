package format

// ForeignDescription is generated from INFORMATION_SCHEMA.KEY_COLUMN_USAGE
type ForeignDescription struct {
	TableName            string `db:"table_name"`
	ColumnName           string `db:"column_name"`
	ConstraintName       string `db:"constraint_name"`
	ReferencedTableName  string `db:"referenced_table_name"`
	ReferencedColumnName string `db:"referenced_column_name"`
}

// ForeignDescriptions is a set of foreign key descriptions
type ForeignDescriptions []ForeignDescription
