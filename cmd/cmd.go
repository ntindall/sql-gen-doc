package cmd

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
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

// Execute is the primary driver for the sql-gen-doc functionality.
func Execute() {
	ctx := context.Background()
	db, err := format.CreateDatabaseConnection(ctx, *flagDSN)
	if err != nil {
		logger.Fatalf("couldn't create database connection. reason: %s", err)
	}

	cfg, err := mysql.ParseDSN(*flagDSN)
	if err != nil {
		logger.Fatalf("failed to parse DSN: %v", err)
	}

	tables, err := format.GetTables(ctx, db, cfg.DBName)
	if err != nil {
		logger.Fatalf("couldn't query database for tables. reason: %s", err)
	}

	markdown := &bytes.Buffer{}
	for idx, table := range tables {
		tableName := table.Name
		columns, err := format.DescribeTable(ctx, db, tableName)
		if err != nil {
			logger.Fatalf("couldn't query database to describe table %s. reason: %s", tableName, err)
		}

		indexData, err := format.GetIndexDescriptions(ctx, db, tableName)
		if err != nil {
			logger.Fatalf("couldn't query database to fetch index data: table %s. reason: %s", tableName, err)
		}

		logicalIndexes, err := indexData.ConvertToLogicalIndexes()
		if err != nil {
			logger.Fatalf("couldn't convert data to logical indexes: table %s. reason: %s", tableName, err)
		}

		_, err = markdown.WriteString(format.CreateTableMarkdown(tableName, table.Comment, columns, logicalIndexes))
		if err != nil {
			logger.Fatalf("error writing to buffer: table %s. reason: %s", tableName, err)
		}

		if idx != len(tables)-1 {
			if _, err = markdown.WriteString("\n"); err != nil {
				logger.Fatalf("error writing to buffer: table %s. reason: %s", tableName, err)
			}
		}

		if *flagOutfile == "" {
			fmt.Fprint(os.Stdout, markdown.String())
			return
		}
		if err := format.WriteToFile(*flagOutfile, markdown.String()); err != nil {
			logger.Fatalln(err)
		}
	}
}
