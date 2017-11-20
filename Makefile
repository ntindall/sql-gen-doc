GO_SRC_FILES = $(shell find . -type f -name '*.go')

.PHONY: build
build: bin/sql-gen-doc

.PHONY: test
test:
	@go test -v $(shell go list ./...)

.PHONY: lint
lint:
	@golint -set_exit_status ./...

bin/sql-gen-doc: $(GO_SRC_FILES)
	go build -o bin/sql-gen-doc ./

