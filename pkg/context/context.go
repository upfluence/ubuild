package context

import (
	"strings"

	"github.com/upfluence/pkg/log"
	"gopkg.in/src-d/go-git.v4"
)

type Context struct {
	Dist    bool
	Version *Version
}

type Version struct {
	Semver, Commit, Branch, Remote string
}

func buildVersion(path, ver string) *Version {
	r, err := git.PlainOpen(path)

	if err != nil {
		log.Fatalf("git: %s", err.Error())
	}

	ref, err := r.Head()

	if err != nil {
		log.Fatalf("git: %s", err.Error())
	}

	remote, err := r.Remote("origin")

	if err != nil {
		log.Fatalf("git: remote origin: %s", err.Error())
	}

	refName := strings.Split(string(ref.Name()), "/")

	return &Version{
		Semver: ver,
		Commit: ref.Hash().String()[:7],
		Branch: refName[len(refName)-1],
		Remote: remote.Config().URLs[0],
	}
}

func BuildLocalContext(path string) *Context {
	return &Context{Dist: false, Version: buildVersion(path, "v0.0.0-dirtry")}
}

func BuildDistContext(path, ver string) *Context {
	return &Context{Dist: true, Version: buildVersion(path, ver)}
}
