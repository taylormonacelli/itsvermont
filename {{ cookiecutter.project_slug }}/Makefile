BIN := {{ cookiecutter.project_slug }}

SRC := $(wildcard *.go **/*.go)

.DEFAULT_GOAL := iterate

all: check $(BIN) install

.PHONY: iterate
iterate: check $(BIN)

.PHONY: check
check: tidy fmt lint vet

.PHONY: build
build: $(BIN)

PREFIX := github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_slug }}/version
LDFLAGS = -s -w

DATE := $(shell date)
GOVERSION := $(shell go version)
VERSION := $(shell git describe --tags --abbrev=8 --dirty --always --long)
SHORT_SHA := $(shell git rev-parse --short HEAD)
FULL_SHA := $(shell git rev-parse HEAD)

LDFLAGS += -X '$(PREFIX).Version=$(VERSION)'
LDFLAGS += -X '$(PREFIX).Date=$(DATE)'
LDFLAGS += -X '$(PREFIX).GoVersion=$(GOVERSION)'
LDFLAGS += -X '$(PREFIX).ShortGitSHA=$(SHORT_SHA)'
LDFLAGS += -X '$(PREFIX).FullGitSHA=$(FULL_SHA)'

$(BIN): $(SRC)
	go build -ldflags "$(LDFLAGS)" -o $@ cmd/main.go

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	gofumpt -w .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: vet
vet:
	go vet ./...

.PHONY: install
install: $(BIN)
	go install

.PHONY: clean
clean:
	$(RM) $(BIN)
