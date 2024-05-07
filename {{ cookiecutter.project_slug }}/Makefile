BIN := {{ cookiecutter.project_slug }}

SRC := $(wildcard *.go **/*.go)

all: tidy fmt lint vet

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
