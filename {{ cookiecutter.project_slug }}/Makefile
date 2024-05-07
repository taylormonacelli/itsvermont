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
LDFLAGS += -X "$(PREFIX).Version={% raw %}{{.Version}}{% endraw %}"
LDFLAGS += -X "$(PREFIX).Date={% raw %}{{.CommitDate}}{% endraw %}"
LDFLAGS += -X "$(PREFIX).GoVersion={% raw %}{{.Env.GOVERSION}}{% endraw %}"
LDFLAGS += -X "$(PREFIX).ShortGitSHA={% raw %}{{.ShortCommit}}{% endraw %}"
LDFLAGS += -X "$(PREFIX).FullGitSHA={% raw %}{{.FullCommit}}{% endraw %}"

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
