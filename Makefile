

.PHONY: build
build: bin/sql-gen-doc

bin/sql-gen-doc: *.go
	go build -o bin/sql-gen-doc ./
