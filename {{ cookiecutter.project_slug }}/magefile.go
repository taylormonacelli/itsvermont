//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	version     = ""
	date        = ""
	goVersion   = ""
	shortGitSHA = ""
	fullGitSHA  = ""
)


const (
	ldFlagsPrefix = "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_slug }}/version"
	buildTarget   = "{{ cookiecutter.project_slug }}"
)


func init() {
	var err error
	version, err = sh.Output("git", "describe", "--tags", "--abbrev=8", "--dirty", "--always", "--long")
	if err != nil {
		fmt.Printf("Error getting version: %v\n", err)
	}
	date, err = sh.Output("date", "+%Y-%m-%dT%H:%M:%SZ")
	if err != nil {
		fmt.Printf("Error getting date: %v\n", err)
	}
	goVersion, err = sh.Output("go", "version")
	if err != nil {
		fmt.Printf("Error getting Go version: %v\n", err)
	}
	shortGitSHA, err = sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		fmt.Printf("Error getting short Git SHA: %v\n", err)
	}
	fullGitSHA, err = sh.Output("git", "rev-parse", "HEAD")
	if err != nil {
		fmt.Printf("Error getting full Git SHA: %v\n", err)
	}
}

// Deps ensures that the required dependencies are installed.
func Deps() error {
	if err := sh.Run("go", "get", "github.com/magefile/mage/mg"); err != nil {
		return err
	}
	if err := sh.Run("go", "get", "github.com/magefile/mage/sh"); err != nil {
		return err
	}
	return nil
}

func Lint() error {
	mg.Deps(Deps)
	return sh.Run("golangci-lint", "run")
}

func Fmt() error {
	mg.Deps(Deps)
	return sh.Run("gofumpt", "-w", ".")
}

func Vet() error {
	mg.Deps(Deps)
	return sh.Run("go", "vet", "./...")
}

func Check() {
	mg.Deps(Lint, Fmt, Vet)
}

func Tidy() error {
	mg.Deps(Deps)
	return sh.Run("go", "mod", "tidy")
}

func Build() error {
	mg.Deps(Tidy)

	ldflags := fmt.Sprintf(`-s -w -X '%s.Version=%s' -X '%s.Date=%s' -X '%s.GoVersion=%s' -X '%s.ShortGitSHA=%s' -X '%s.FullGitSHA=%s'`,
		ldFlagsPrefix, version,
		ldFlagsPrefix, date,
		ldFlagsPrefix, goVersion,
		ldFlagsPrefix, shortGitSHA,
		ldFlagsPrefix, fullGitSHA,
	)

	return sh.Run("go", "build", "-ldflags", ldflags, "-o", buildTarget, "main.go")
}

func Install() error {
	mg.Deps(Build)
	return sh.Run("go", "install")
}

func Clean() error {
	mg.Deps(Deps)
	return sh.Rm(buildTarget)
}