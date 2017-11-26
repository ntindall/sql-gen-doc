# sql-gen-doc
[![GoDoc](https://godoc.org/github.com/ntindall/sql-gen-doc?status.svg)](https://godoc.org/github.com/ntindall/sql-gen-doc) [![CircleCI](https://circleci.com/gh/ntindall/sql-gen-doc.svg?style=svg)](https://circleci.com/gh/ntindall/sql-gen-doc)

A tool to automatically generate sql documentation.

## Installation
```sh
  go get -u github.com/ntindall/sql-gen-doc
```

## Usage

`sql-gen-doc` will connect to a database and generate a markdown table
corresponding to the current state of each table in the database. This is useful
for databases that undergo frequent migrations. You can set up your CI to run
this tool whenever a new migration is added.

```sh
$ ./bin/sql-gen-doc --help
Usage of ./bin/sql-gen-doc:
  -dsn string
      a data source name for the database, e.g. user:password@tcp(mysql:3306)/database_name
  -o string
      the outfile to write the documentation to, if no outfile is specified, the output is written to stdout

$ ./bin/sql-gen-doc -dsn 'user:password@tcp(localhost:3306)/database_to_generate' -out outfile.md
```

Additionally, the markdown file can be annotated with comments in order to have
`sql-gen-doc` insert the table into a specific location in an existing file. For this
to work, just add these comments to your markdown file and then specify it as the
outfile via the command line flag.

```markdown
# fake markdown

<!-- sql-gen-doc BEGIN -->
database documentation will go here!
<!-- sql-gen-doc END -->

more documentation!
```

## Output

See [fixtures/expected1.md](fixtures/expected1.md) and [fixtures/expected2.md](fixtures/expected2.md) for examples.

## Development

1. This project uses `docker` and `docker-compose` for testing, see [here](https://docs.docker.com/compose/install/)
   for the setup instructions for your operating system.

2. Run the following to clone and setup the project.

  ```sh
    git clone git@github.com:ntindall/sql-gen-doc.git $GOPATH/src/github.com/ntindall/sql-gen-doc
    cd $GOPATH/src/github.com/ntindall/sql-gen-doc
    make setup
  ```

3. Run the test suite

  ```sh
    make docker-test
  ```