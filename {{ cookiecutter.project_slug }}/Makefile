BIN := {{ cookiecutter.project_slug }}

SRC := $(wildcard *.go **/*.go)

.PHONY: all tidy fmt lint install clean

all: tidy fmt lint vet install

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

$(BIN): $(SRC)
	go build -o $@ cmd/main.go

.PHONY: install
install: $(BIN)
	go install

.PHONY: clean
clean:
	$(RM) $(BIN)
