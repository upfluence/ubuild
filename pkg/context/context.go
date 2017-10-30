package context

import (
	"regexp"
	"strings"

	"github.com/upfluence/pkg/log"
	"gopkg.in/src-d/go-git.v4"
)

var versionRegexp = regexp.MustCompile("^v\\d.\\d.\\d(-.+)?$")

type Context struct {
	Dist    bool
	Version *Version
}

type Version struct {
	Semver, Commit, Branch, Remote string
}

func buildVersion(path string) *Version {
	var r, err = git.PlainOpen(path)

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
		Commit: ref.Hash().String(),
		Branch: refName[len(refName)-1],
		Remote: remote.Config().URLs[0],
	}
}

func BuildContext(path string, dist bool) *Context {
	return &Context{Dist: dist, Version: buildVersion(path)}
}
