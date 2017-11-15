package main

import (
	"context"
	"io"
	"os"

	"github.com/jmoiron/sqlx"
)

const (
	defaultBufferSize = 65536
)

func createDatabaseConnection(
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

func showTables(
	ctx context.Context,
	db *sqlx.DB,
) ([]string, error) {
	result := []string{}

	if err := db.SelectContext(ctx, &result, "SHOW TABLES"); err != nil {
		return nil, err
	}

	return result, nil
}

func describeTable(
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

func writeToFile(
	filename string,
	markdown string,
) {
	file := ""
	out, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	defer out.Close()

	if err != nil {
		logger.Fatalf("couldn't open %s for writing. reason: %s", *flagOutfile, err)
	}

	n := defaultBufferSize
	for {
		buffer := make([]byte, n)

		n, err = out.Read(buffer)
		if err != nil && err != io.EOF {
			logger.Fatalf("error while reading. reason: %s", err)
		}

		file += string(buffer)

		if n == 0 || err == io.EOF {
			break
		}
	}

	// If the existing file is annotated with the requisite comments, we insert
	// between them.
	processedMarkdown := insertBetweenTags(file, markdown)

	// Reset the fd before writing
	out.Seek(0, 0)
	out.Truncate(int64(len(processedMarkdown)))

	remainingIdx := 0
	for {
		logger.Print("iterating")
		written, err := out.WriteString(processedMarkdown[remainingIdx:])
		if err != nil {
			logger.Fatalf("error while writing. reason: %s", err)
		}

		if written == len(processedMarkdown) {
			break
		}
		remainingIdx += written
	}
}
