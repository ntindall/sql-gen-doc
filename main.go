package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	// We must import the mysql driver as it is required by sqlx.
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	defaultBufferSize = 65536
)

var (
	flagDSN     *string
	flagOutfile *string
	logger      *log.Logger
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

func formatDescription(
	table string,
	columns []columnDescription,
) string {
	tableMarkdown := makeTitle(table)

	formatSpec := getFormatSpec(columns)
	tableMarkdown = tableMarkdown + makeHeader(formatSpec)

	for _, c := range columns {
		tableMarkdown += c.Format(formatSpec)
	}

	tableMarkdown += "\n"

	return tableMarkdown
}

func init() {
	flagDSN = flag.String("dsn", "", "a data source name for the database, e.g. paramuser:password@tcp(mysql:3306)/database_name")
	flagOutfile = flag.String("o", "", "the outfile to write the documentation to, if no outfile is specified, the output is written to stdout")

	// Setup logging
	logger = log.New(os.Stderr, "[sql-gen-doc] ", log.Lshortfile)
}

func main() {
	flag.Parse()

	if *flagDSN == "" {
		logger.Fatalln("the -dsn flag must be provided")
	}

	ctx := context.Background()
	db, err := createDatabaseConnection(ctx, *flagDSN)
	if err != nil {
		logger.Fatalf("couldn't create database connection. reason: %s", err)
	}

	tables, err := showTables(ctx, db)
	if err != nil {
		logger.Fatalf("couldn't query database for tables. reason: %s", err)
	}

	markdown := ""
	for _, tableName := range tables {
		description, err := describeTable(ctx, db, tableName)
		if err != nil {
			logger.Fatalf("couldn't query database to describe table %s. reason: %s", tableName, err)
		}

		markdown += formatDescription(tableName, description)
	}

	if *flagOutfile == "" {
		fmt.Fprintf(os.Stdout, markdown)
	} else {
		writeToFile(*flagOutfile, markdown)
	}
}

func insertBetweenTags(
	file string,
	markdown string,
) string {
	startTag := "<!-- sql-gen-doc BEGIN -->"
	endTag := "<!-- sql-gen-doc END -->"

	// r := strings.NewReplacer(" ", "", "\t", "")
	// stripped := r.Replace(file)

	startIdx := strings.Index(file, startTag)
	endIdx := strings.Index(file, endTag)
	logger.Print(startIdx, endIdx)

	if startIdx == -1 || endIdx == -1 {
		logger.Printf("returning markdown")
		return markdown
	}

	startIdx += len(startTag)
	endIdx += len(endTag)
	return file[:startIdx] + "\n" + markdown + "\n"
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
