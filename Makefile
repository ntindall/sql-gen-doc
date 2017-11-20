GO_SRC_FILES = $(shell find . -type f -name '*.go')

# build

.PHONY: build
build: bin/sql-gen-doc

bin/sql-gen-doc: $(GO_SRC_FILES)
	go build -o bin/sql-gen-doc ./

# testing / linting

.PHONY: test
test: go-test go-lint

.PHONY: go-test
go-test:
	@go test -v $(shell go list ./...)

.PHONY: go-lint
go-lint:
	@golint -set_exit_status ./...
