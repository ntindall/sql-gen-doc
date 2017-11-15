package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	// We must import the mysql driver as it is required by sqlx.
	_ "github.com/go-sql-driver/mysql"
)

var (
	flagDSN     *string
	flagOutfile *string
	logger      *log.Logger
)

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
