BIN := {{ cookiecutter.project_slug }}

SRC := $(wildcard *.go **/*.go)

.DEFAULT_GOAL := $(BIN)

all: check $(BIN)

.PHONY: check
check: tidy fmt lint vet

$(BIN): $(SRC)
	go build -o $@ cmd/main.go

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
