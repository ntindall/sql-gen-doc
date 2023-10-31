package format

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jmoiron/sqlx"
)

// CreateDatabaseConnection creates a connection to the database. The connection
// is long lived and should only be created once per process.
func CreateDatabaseConnection(
	ctx context.Context,
	dsn string,
) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	return db, nil
}

// GetTablesRow see GetTablesRow
type GetTablesRow struct {
	Name    string
	Comment string
}

// GetTables queries the database and returns a list of the tables that
// are present in the database.
func GetTables(ctx context.Context, db *sqlx.DB, dbName string) ([]GetTablesRow, error) {
	result := []GetTablesRow{}

	if err := db.SelectContext(ctx, &result,
		`SELECT table_name AS name, table_comment AS comment 
		FROM information_schema.tables 
		WHERE table_schema = ?`,
		dbName,
	); err != nil {
		return nil, err
	}

	return result, nil
}

// DescribeTable queries the database for information about the specified
// table. The result is scanned into a ColumnDescription struct.
func DescribeTable(
	ctx context.Context,
	db *sqlx.DB,
	tableName string,
) ([]ColumnDescription, error) {
	result := []ColumnDescription{}

	// use show full columns rather than describe to access comments
	if err := db.SelectContext(ctx, &result, "SHOW FULL COLUMNS FROM "+tableName); err != nil {
		return nil, err
	}

	return result, nil
}

// GetIndexDescriptions queries the database for information about the specified
// table. The result is scanned into a IndexDescription struct.
func GetIndexDescriptions(
	ctx context.Context,
	db *sqlx.DB,
	tableName string,
) (IndexDescriptions, error) {
	result := []IndexDescription{}

	if err := db.SelectContext(ctx, &result, "SHOW INDEXES FROM "+tableName); err != nil {
		return nil, err
	}

	return result, nil
}

// GetForeignKeyDescriptions queries INFORMATION_SCHEMA.KEY_COLUMN_USAGE table about references information
func GetForeignKeyDescriptions(
	ctx context.Context,
	db *sqlx.DB,
	tableName string,
) (ForeignDescriptions, error) {
	var result []ForeignDescription

	if err := db.SelectContext(ctx, &result,
		`SELECT table_name, column_name, constraint_name, referenced_table_name, referenced_column_name
		FROM information_schema.key_column_usage
		WHERE table_name = ? AND referenced_table_name IS NOT NULL
		ORDER BY 1,2`,
		tableName,
	); err != nil {
		return nil, err
	}

	return result, nil
}

// WriteToFile takes a filename and a markdown string and writes the markdown
// to the file. If the file is annotated with markdown comments, the markdown
// will be inserted in between the comments. e.g.
//
// # fake markdown
//
// <!-- sql-gen-doc BEGIN -->
// markdown will go here!
// <!-- sql-gen-doc END -->"
//
// An error is returned if the file cannot be written.
func WriteToFile(
	filename string,
	markdown string,
) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("couldn't open %s for reading. reason: %s", filename, err)
	}

	// If the existing file is annotated with the requisite comments, we insert
	// between them.
	processedMarkdown, err := insertBetweenTags(string(file), markdown)
	if err != nil {
		return fmt.Errorf("couldn't insert markdown into file. reason: %s", err)
	}

	return ioutil.WriteFile(filename, []byte(processedMarkdown), 0644)
}
