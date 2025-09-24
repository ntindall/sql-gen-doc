# sql-gen-doc
[![GoDoc](https://godoc.org/github.com/ntindall/sql-gen-doc?status.svg)](https://godoc.org/github.com/ntindall/sql-gen-doc) [![CircleCI](https://circleci.com/gh/ntindall/sql-gen-doc.svg?style=svg)](https://circleci.com/gh/ntindall/sql-gen-doc)

A command-line tool to automatically generate comprehensive SQL database documentation in Markdown format. Connect to any MySQL database and generate detailed documentation including table schemas, column descriptions, indexes, and foreign key relationships.

## Installation

### Using Go
```sh
go get -u github.com/ntindall/sql-gen-doc
```

### Building from Source
```sh
git clone https://github.com/ntindall/sql-gen-doc.git
cd sql-gen-doc
go build -o sql-gen-doc ./cmd
```

## Usage

`sql-gen-doc` connects to a MySQL database and generates comprehensive Markdown documentation reflecting the current state of your database schema. This includes:

- **Table descriptions**: Complete table schemas with column details
- **Index information**: Primary keys, unique constraints, and regular indexes  
- **Foreign key relationships**: Table relationships and constraints
- **Column metadata**: Data types, constraints, and descriptions

This tool is particularly useful for databases undergoing frequent migrations. Set up your CI/CD pipeline to run this tool whenever new migrations are added to keep your documentation up-to-date automatically.

### Command Line Options

```sh
$ ./sql-gen-doc --help
Usage of sql-gen-doc:
  -dsn string
      a data source name for the database, e.g. user:password@tcp(mysql:3306)/database_name
  -o string
      the output file to write the documentation to (optional, writes to stdout if not specified)
  --sort-tables
      outputs tables in alphabetical order

$ ./sql-gen-doc --banana
Banana

$ ./sql-gen-doc -dsn 'user:password@tcp(localhost:3306)/database_name' -o documentation.md --sort-tables
```

### Updating Existing Documentation Files

`sql-gen-doc` can intelligently update existing Markdown files by inserting generated documentation between special comment markers. This is perfect for maintaining documentation files that contain additional content.

Simply add these comment tags to your existing Markdown file:

```markdown
# Database Documentation

<!-- sql-gen-doc BEGIN -->
The database documentation will be automatically generated and inserted here.
<!-- sql-gen-doc END -->

## Additional Notes
More documentation content can be placed below or above the generated section.
```

When you run `sql-gen-doc` with this file as the output destination (`-o mydocs.md`), it will replace only the content between the comment markers while preserving everything else.

## Example Output

The generated documentation includes comprehensive table schemas formatted as Markdown tables. Check out the example outputs in:

- [fixtures/expected1.md](fixtures/expected1.md) - Basic table documentation
- [fixtures/expected2.md](fixtures/expected2.md) - Documentation with indexes and foreign keys

Each table documentation includes:

- **Column information**: Name, data type, constraints, and description
- **Indexes**: Details about primary keys and indexes with columns
- **Foreign keys**: Relationships between tables with referenced tables and columns

## Development

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose (for testing)

### Setup

1. Clone the repository:

  ```sh
  git clone https://github.com/ntindall/sql-gen-doc.git
  cd sql-gen-doc
  ```

2. Install dependencies and setup the project:

  ```sh
  make setup
  ```

3. Run tests:

  ```sh
  make docker-test
  ```

### Running the Application

Build and run the application:

```sh
go build -o sql-gen-doc ./cmd
./sql-gen-doc -dsn "user:password@tcp(localhost:3306)/database_name" -o output.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.