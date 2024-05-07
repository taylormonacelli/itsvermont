BIN := {{ cookiecutter.project_slug }}

SRC := $(wildcard *.go **/*.go)

.DEFAULT_GOAL := all

all: check $(BIN)

.PHONY: check
check: tidy fmt lint vet

.PHONY: build
build: $(BIN)

LDFLAGS = -s -w
LDFLAGS += -X "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_slug }}/version.Version={% raw %}{{.Version}}{% endraw %}"
LDFLAGS += -X "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_slug }}/version.Date={% raw %}{{.CommitDate}}{% endraw %}"
LDFLAGS += -X "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_slug }}/version.GoVersion={% raw %}{{.Env.GOVERSION}}{% endraw %}"
LDFLAGS += -X "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_slug }}/version.ShortGitSHA={% raw %}{{.ShortCommit}}{% endraw %}"
LDFLAGS += -X "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_slug }}/version.FullGitSHA={% raw %}{{.FullCommit}}{% endraw %}"

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
