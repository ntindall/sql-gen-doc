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

// ShowTables queries the database and returns a list of the tables that
// are present in the database.
func ShowTables(
	ctx context.Context,
	db *sqlx.DB,
) ([]string, error) {
	result := []string{}

	if err := db.SelectContext(ctx, &result, "SHOW TABLES"); err != nil {
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

	if err := db.SelectContext(ctx, &result, "DESCRIBE "+tableName); err != nil {
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
