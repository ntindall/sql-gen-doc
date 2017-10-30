package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"

	// We must import the mysql driver as it is required by sqlx.
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

	var out io.Writer = os.Stdout
	if *flagOutfile != "" {
		out, err = os.Create(*flagOutfile)
		if err != nil {
			logger.Fatalf("couldn't open %s for writing. reason: %s", *flagOutfile, err)
		}
	}

	out.Write([]byte(markdown))
}
