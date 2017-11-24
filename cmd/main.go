package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	// We must import the mysql driver as it is required by sqlx.
	_ "github.com/go-sql-driver/mysql"
	"github.com/ntindall/sql-gen-doc/format"
)

var (
	flagDSN     *string
	flagOutfile *string
	logger      *log.Logger
)

func init() {
	// Parse flags
	flagDSN = flag.String("dsn", "", "a data source name for the database, e.g. user:password@tcp(mysql:3306)/database_name")
	flagOutfile = flag.String("o", "", "the outfile to write the documentation to, if no outfile is specified, the output is written to stdout")
	flag.Parse()

	// Setup logging
	logger = log.New(os.Stderr, "[sql-gen-doc] ", log.Lshortfile)

	// Validate flags
	if *flagDSN == "" {
		logger.Fatalln("the -dsn flag must be provided")
	}
}

func main() {
	ctx := context.Background()
	db, err := format.CreateDatabaseConnection(ctx, *flagDSN)
	if err != nil {
		logger.Fatalf("couldn't create database connection. reason: %s", err)
	}

	tables, err := format.ShowTables(ctx, db)
	if err != nil {
		logger.Fatalf("couldn't query database for tables. reason: %s", err)
	}

	markdown := ""
	for _, tableName := range tables {
		columns, err := format.DescribeTable(ctx, db, tableName)
		if err != nil {
			logger.Fatalf("couldn't query database to describe table %s. reason: %s", tableName, err)
		}

		markdown += format.CreateTableMarkdown(tableName, columns)
		markdown += "\n"
	}

	if *flagOutfile == "" {
		fmt.Fprintf(os.Stdout, markdown)
		return
	}
	if err := format.WriteToFile(*flagOutfile, markdown); err != nil {
		logger.Fatalln(err)
	}
}
