.PHONY: build
build: bin/sql-gen-doc

.PHONY: test
test:
	@go test -v

bin/sql-gen-doc: *.go
	go build -o bin/sql-gen-doc ./

