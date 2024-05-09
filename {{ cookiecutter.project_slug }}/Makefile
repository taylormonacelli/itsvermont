BIN := {{ cookiecutter.project_slug }}

SRC := $(wildcard *.go **/*.go)

DATE := $(shell date +"%Y-%m-%dT%H:%M:%SZ")
GOVERSION := $(shell go version)
VERSION := $(shell git describe --tags --abbrev=8 --dirty --always --long)
SHORT_SHA := $(shell git rev-parse --short HEAD)
FULL_SHA := $(shell git rev-parse HEAD)

export GOVERSION

PREFIX := github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_slug }}/version
LDFLAGS = -s -w
LDFLAGS += -X '$(PREFIX).Version=$(VERSION)'
LDFLAGS += -X '$(PREFIX).Date=$(DATE)'
LDFLAGS += -X '$(PREFIX).GoVersion=$(GOVERSION)'
LDFLAGS += -X '$(PREFIX).ShortGitSHA=$(SHORT_SHA)'
LDFLAGS += -X '$(PREFIX).FullGitSHA=$(FULL_SHA)'

.DEFAULT_GOAL := iterate

all: check $(BIN) install

.PHONY: iterate # lint and rebuild
iterate: check $(BIN)

.PHONY: check # lint and vet
check: tidy fmt lint vet

.PHONY: build # build
build: $(BIN)

$(BIN): $(SRC) # $(SRC)
	go build -ldflags "$(LDFLAGS)" -o $@ cmd/main.go

.PHONY: goreleaser
goreleaser:
	goreleaser --clean

.PHONY: tidy # go tidy
tidy:
	go mod tidy

.PHONY: fmt # go fmt
fmt:
	gofumpt -w .

.PHONY: lint # lint
lint:
	golangci-lint run

.PHONY: vet # go vet
vet:
	go vet ./...

.PHONY: install # go install
install: $(BIN)
	go install

.PHONY: help # show makefile rules
help:
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/\1	\2/' | expand -t20

.PHONY: clean # clean bin
clean:
	$(RM) $(BIN)
