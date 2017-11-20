# sql-gen-doc
[![GoDoc](https://godoc.org/github.com/ntindall/sql-gen-doc?status.svg)](https://godoc.org/github.com/ntindall/sql-gen-doc) [![CircleCI](https://circleci.com/gh/ntindall/sql-gen-doc.svg?style=svg)](https://circleci.com/gh/ntindall/sql-gen-doc)

A tool to automatically generate sql documentation.

## Installation
```sh
  go get -u github.com/ntindall/sql-gen-doc
```

Not yet ready for production use :)

## Example use case

```sh
./bin/sql-gen-doc -dsn 'user:password@tcp(localhost:3306)/database_to_generate'
```
