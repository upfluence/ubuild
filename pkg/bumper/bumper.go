package bumper

import (
	"github.com/upfluence/pkg/log"

	"github.com/upfluence/ubuild/pkg/config"
	"github.com/upfluence/ubuild/pkg/context"
	"github.com/upfluence/ubuild/pkg/githubutil"
	"github.com/upfluence/ubuild/pkg/version"
)

var (
	branchBumpFns = map[string]func(*version.Version){
		"master": func(v *version.Version) { v.IncPatch() },
	}

	defaultBumpFn = func(v *version.Version) { v.IncRC() }
)

func Bump(ctx *context.Context, cfg *config.Configuration) (*version.Version, error) {
	v, err := githubutil.GetLastVersion(cfg.GetRepo())

	if err != nil {
		return nil, err
	}

	messages, err := githubutil.FetchCommits(
		cfg.GetRepo(),
		v.String(),
		ctx.Version.Commit,
	)

	if err != nil {
		return nil, err
	}

	if !version.IncrementVersionFromCommits(v, messages) {
		log.Notice(ctx.Version.Branch)
		if fn, ok := branchBumpFns[ctx.Version.Branch]; ok {
			fn(v)
		} else {
			defaultBumpFn(v)
		}
	}

	log.Notice(v)

	return v, githubutil.CreateRelease(cfg.GetRepo(), ctx.Version.Commit, v)
}
