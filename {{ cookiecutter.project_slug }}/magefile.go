//go:build mage
// +build mage

package main

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	. "github.com/go-git/go-git/v5/_examples"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

var (
	version     = ""
	date        = ""
	goVersion   = ""
	shortGitSHA = ""
	fullGitSHA  = ""
)

var Default = Iterate // default mage target

const (
	ldFlagsPrefix = "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_slug }}/version"
	buildTarget   = "{{ cookiecutter.project_slug }}"
)

func init() {
	var err error
	date = time.Now().UTC().Format(time.RFC3339)

	r, err := git.PlainOpen(".")
	CheckIfError(err)

	ref, err := r.Head()
	CheckIfError(err)

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err)

	commit, err := cIter.Next()
	CheckIfError(err)

	version = commit.Hash.String()[:8]
	if commit.NumParents() > 0 {
		version += "-dirty"
	}

	goVersion, err = sh.Output("go", "version")
	if err != nil {
		fmt.Printf("Error getting Go version: %v\n", err)
	}

	shortGitSHA = commit.Hash.String()[:7]
	fullGitSHA = commit.Hash.String()
}

func Iterate() {
	mg.Deps(Check, Build)
}

func Lint() error {
	return sh.Run("golangci-lint", "run")
}

func Fmt() error {
	return sh.Run("gofumpt", "-w", ".")
}

func Vet() error {
	return sh.Run("go", "vet", "./...")
}

func Check() {
	mg.Deps(Lint, Fmt, Vet)
}

func Tidy() error {
	return sh.Run("go", "mod", "tidy")
}

func Build() error {
	mg.Deps(Tidy)

	if rebuild, err := target.Glob(buildTarget, "*.go", "go.mod", "go.sum"); err != nil || !rebuild {
		if err != nil {
			return err
		}
		return nil
	}

	ldflags := fmt.Sprintf(`-s -w -X '%s.Version=%s' -X '%s.Date=%s' -X '%s.GoVersion=%s' -X '%s.ShortGitSHA=%s' -X '%s.FullGitSHA=%s'`,
		ldFlagsPrefix, version,
		ldFlagsPrefix, date,
		ldFlagsPrefix, goVersion,
		ldFlagsPrefix, shortGitSHA,
		ldFlagsPrefix, fullGitSHA,
	)

	return sh.Run("go", "build", "-ldflags", ldflags, "-o", buildTarget, "cmd/main.go")
}

func Install() error {
	mg.Deps(Build)
	return sh.Run("go", "install")
}

func Clean() error {
	return sh.Rm(buildTarget)
}
