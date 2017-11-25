GO_SRC_FILES = $(shell find . -type f -name '*.go' | sed /vendor/d )
GO_SRC_PACKAGES =$(shell go list ./... | sed /vendor/d )
GOLINT_SRC = ./vendor/github.com/golang/lint/golint
GOOSE_SRC = ./vendor/github.com/pressly/goose/cmd/goose

# vanity
GREEN = \033[0;32m
MAGENTA = \033[0;35m
RESET = \033[0;0m

# setup
.PHONY: setup
setup: install-dep vendor bin/golint bin/goose

.PHONY: clean
clean:
	rm -rf logs/*

.PHONY: install-dep
install-dep:
	@./scripts/install-dep.sh

.PHONY: update
update:
	@echo "$(GREEN)updating vendored dependencies...$(RESET)"
	@dep ensure -v

vendor: Gopkg.toml Gopkg.lock
	@echo "$(GREEN)installing vendored dependencies...$(RESET)"
	@# use the vendor-only flag to prevent us from removing dependencies before
	@# they are added to the docker container
	@dep ensure -v --vendor-only

bin/golint: vendor
	@echo "$(MAGENTA)building $(@)...$(RESET)"
	@go build -o $(@) $(GOLINT_SRC)

bin/goose: vendor
	@echo "$(MAGENTA)building $(@)...$(RESET)"
	@go build -o $(@) $(GOOSE_SRC)

# build

.PHONY: build
build: bin/sql-gen-doc

bin/sql-gen-doc: $(GO_SRC_FILES)
	@echo "$(MAGENTA)building $(@)...$(RESET)"
	go build -o bin/sql-gen-doc ./cmd

# images
.PHONY: images
VERSION = 0.0.1
images: .circleci/images/primary/Dockerfile
	# TODO: is there a way to get the last tag following semver?
	@echo "$(MAGENTA)current version is $(VERSION), did you bump this value?\n\
press (enter to continue)?$(RESET)"
	@read
	@echo "$(MAGENTA)building a new base image with tag $(VERSION)...$(RESET)"
	docker build -t pumpkinobsessed/sql-gen-doc:$(VERSION) $(^:%/Dockerfile=%)
	docker login
	docker tag pumpkinobsessed/sql-gen-doc:$(VERSION) pumpkinobsessed/sql-gen-doc:latest
	docker push pumpkinobsessed/sql-gen-doc:$(VERSION)
	docker push pumpkinobsessed/sql-gen-doc:latest

# testing / linting

.PHONY: docker-test
docker-test:
	@docker-compose -f docker-compose.yml build test_container
	@docker-compose -f docker-compose.yml run test_container make test integrate

.PHONY: integrate
integrate:
	./scripts/integrate.sh

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

# releasing
.PHONY: release
release:
	./scripts/release.sh

# migrations
.PHONY: migrate-up migrate-down migrate-reset
migrate-up migrate-down migrate-reset:
	@./bin/goose -dir goose mysql "$(MYSQL_USER):$(MYSQL_PASSWORD)@(mysql:3306)/example" $(@:migrate-%=%)
