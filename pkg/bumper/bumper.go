package bumper

import (
	"errors"

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

	bumpStrategies = map[string]func(*version.Version){
		"bump_patch": func(v *version.Version) { v.IncPatch() },
		"bump_rc":    func(v *version.Version) { v.IncRC() },
	}

	errBumpStrategiesNotFound = errors.New("bump strategies not found")
)

func Bump(ctx *context.Context, cfg *config.Configuration) (*version.Version, error) {
	v, err := bumpVersion(ctx, cfg)

	if err != nil {
		return nil, err
	}

	log.Notice(v)

	return v, githubutil.CreateRelease(cfg.GetRepo(), ctx.Version.Commit, v)
}

func bumpVersion(ctx *context.Context, cfg *config.Configuration) (*version.Version, error) {
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

	if version.IncrementVersionFromCommits(v, messages) {
		return v, nil
	}

	if st, ok := cfg.CustomBumpStrategies[ctx.Version.Branch]; ok {
		if fn, ok := bumpStrategies[st]; ok {
			fn(v)
			return v, nil
		} else {
			return nil, errBumpStrategiesNotFound
		}
	}

	log.Notice(ctx.Version.Branch)
	if fn, ok := branchBumpFns[ctx.Version.Branch]; ok {
		fn(v)
	} else {
		defaultBumpFn(v)
	}

	return v, nil
}
