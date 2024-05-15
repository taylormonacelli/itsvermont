//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
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
	sources, err := getSourceFiles(".")
	if err != nil {
		return err
	}

	lintTime := time.Time{}
	if _, err := os.Stat(".lint.time"); err == nil {
		lintTime, err = target.NewestModTime(".lint.time")
		if err != nil {
			return err
		}
	}

	if rebuild, err := target.DirNewer(lintTime, sources...); err != nil || !rebuild {
		if err != nil {
			return err
		}
		return nil
	}

	if err := sh.Run("golangci-lint", "run"); err != nil {
		return err
	}

	return touch(".lint.time")
}

func Fmt() error {
	sources, err := getSourceFiles(".")
	if err != nil {
		return err
	}

	fmtTime := time.Time{}
	if _, err := os.Stat(".fmt.time"); err == nil {
		fmtTime, err = target.NewestModTime(".fmt.time")
		if err != nil {
			return err
		}
	}

	if rebuild, err := target.DirNewer(fmtTime, sources...); err != nil || !rebuild {
		if err != nil {
			return err
		}
		return nil
	}

	if err := sh.Run("gofumpt", "-w", "."); err != nil {
		return err
	}

	return touch(".fmt.time")
}

func Vet() error {
	sources, err := getSourceFiles(".")
	if err != nil {
		return err
	}

	vetTime := time.Time{}
	if _, err := os.Stat(".vet.time"); err == nil {
		vetTime, err = target.NewestModTime(".vet.time")
		if err != nil {
			return err
		}
	}

	if rebuild, err := target.DirNewer(vetTime, sources...); err != nil || !rebuild {
		if err != nil {
			return err
		}
		return nil
	}

	if err := sh.Run("go", "vet", "./..."); err != nil {
		return err
	}

	return touch(".vet.time")
}

func Check() {
	mg.Deps(Lint, Fmt, Vet)
}

func Tidy() error {
	tidyTime := time.Time{}
	if _, err := os.Stat(".tidy.time"); err == nil {
		tidyTime, err = target.NewestModTime(".tidy.time")
		if err != nil {
			return err
		}
	}

	if rebuild, err := target.DirNewer(tidyTime, "go.mod", "go.sum"); err != nil || !rebuild {
		if err != nil {
			return err
		}
		return nil
	}

	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return err
	}

	return touch(".tidy.time")
}

func Build() error {
	mg.Deps(Tidy)
	buildTime := time.Time{}
	if _, err := os.Stat(buildTarget); err == nil {
		buildTime, err = target.NewestModTime(buildTarget)
		if err != nil {
			return err
		}
	}
	sources, err := getSourceFiles(".")
	if err != nil {
		return err
	}
	sources = append(sources, "go.mod", "go.sum")
	if rebuild, err := target.DirNewer(buildTime, sources...); err != nil || !rebuild {
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
	return sh.Run("go", "build", "-ldflags", ldflags, "-o", buildTarget)
}

func Install() error {
	mg.Deps(Build)
	return sh.Run("go", "install")
}

func Clean() error {
	return sh.Rm(buildTarget)
}

func getSourceFiles(dir string, exts ...string) ([]string, error) {
	if len(exts) == 0 {
		exts = []string{
			".go",
			".cue",
		}
	}
	var sources []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			for _, ext := range exts {
				if filepath.Ext(path) == ext {
					sources = append(sources, path)
					break
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return sources, nil
}

func touch(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	return f.Close()
}
