GO_SRC_FILES = $(shell find . -type f -name '*.go' | sed /vendor/d )
GO_SRC_PACKAGES =$(shell go list ./... | sed /vendor/d )
GOLINT_SRC = ./vendor/github.com/golang/lint/golint
# vanity
GREEN = \033[0;32m
MAGENTA = \033[0;35m
RESET = \033[0;0m

# setup
.PHONY: setup
setup: install-dep vendor bin/golint

.PHONY: install-dep
install-dep:
	@scripts/install-dep.sh

vendor: Gopkg.toml Gopkg.lock
	@echo "$(GREEN)installing vendored dependencies...$(RESET)"
	@dep ensure -v

bin/golint: vendor
	@echo "$(MAGENTA)building $(@)...$(RESET)"
	@go build -o $(@) ./vendor/github.com/golang/lint/golint

# build

.PHONY: build
build: bin/sql-gen-doc

bin/sql-gen-doc: $(GO_SRC_FILES)
	@echo "$(MAGENTA)building $(@)...$(RESET)"
	go build -o bin/sql-gen-doc ./cmd

# testing / linting

.PHONY: docker-test
docker-test:
	@docker-compose -f docker-compose.yml build test_container
	@docker-compose -f docker-compose.yml run test_container make test

.PHONY: test
test: go-test go-lint build

.PHONY: go-test
go-test:
	@echo "$(MAGENTA)running go tests...$(RESET)"
	@go test -v $(GO_SRC_PACKAGES)

.PHONY: go-lint
go-lint: bin/golint
	@echo "$(MAGENTA)linting $(GO_SRC_PACKAGES)$(RESET)"
	@bin/golint -set_exit_status $(GO_SRC_PACKAGES)
