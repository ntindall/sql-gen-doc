package format

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/jmoiron/sqlx"
)

const (
	defaultBufferSize = 65536
)

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

func DescribeTable(
	ctx context.Context,
	db *sqlx.DB,
	tableName string,
) ([]columnDescription, error) {
	result := []columnDescription{}

	if err := db.SelectContext(ctx, &result, "DESCRIBE "+tableName); err != nil {
		return nil, err
	}

	return result, nil
}

// TODO TEST
func WriteToFile(
	filename string,
	markdown string,
) error {
	file := ""
	out, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	defer out.Close()

	if err != nil {
		return fmt.Errorf("couldn't open %s for writing. reason: %s", filename, err)
	}

	n := defaultBufferSize
	for {
		buffer := make([]byte, n)

		n, err = out.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error while reading. reason: %s", err)
		}

		file += string(buffer)

		if n == 0 || err == io.EOF {
			break
		}
	}

	// If the existing file is annotated with the requisite comments, we insert
	// between them.
	processedMarkdown, err := insertBetweenTags(file, markdown)
	if err != nil {
		return fmt.Errorf("couldn't insert markdown into file. reason: %s", err)
	}

	// Reset the fd before writing
	out.Seek(0, 0)
	out.Truncate(int64(len(processedMarkdown)))

	remainingIdx := 0
	for {
		written, err := out.WriteString(processedMarkdown[remainingIdx:])
		if err != nil {
			return fmt.Errorf("error while writing. reason: %s", err)
		}

		if written == len(processedMarkdown) {
			break
		}
		remainingIdx += written
	}

	return nil
}
