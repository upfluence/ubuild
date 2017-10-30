package builder

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/upfluence/ubuild/pkg/config"
	"github.com/upfluence/ubuild/pkg/context"
	"github.com/upfluence/ubuild/pkg/sh"
)

func Build(ctx *context.Context, cfg *config.Configuration) error {
	buf, err := ioutil.ReadFile(cfg.GetBuilder().GetDockerfile())

	if err != nil {
		return err
	}

	for _, l := range strings.Split(string(buf), "\n") {
		if splittedLine := strings.Split(
			l,
			" ",
		); len(splittedLine) == 2 && splittedLine[0] == "FROM" {
			if err := sh.RunCommand(
				cfg.GetVerbose(),
				"docker",
				"pull",
				splittedLine[1],
			); err != nil {
				return err
			}

			break
		}
	}

	args := append(
		[]string{"build"},
		buildArgs(ctx)...,
	)

	args = append(
		args,
		"--no-cache",
		"-t",
		cfg.GetBuilder().GetImage()+":"+ctx.Version.Commit,
		cfg.GetPath(),
	)

	if err := sh.RunCommand(
		cfg.GetVerbose(),
		"docker",
		args...,
	); err != nil {
		return err
	}

	if err := tagImage(cfg, ctx.Version.Commit, ctx.Version.Semver); err != nil {
		return err
	}

	if t, ok := cfg.GetBuilder().Tags[ctx.Version.Branch]; ok {
		if err := tagImage(cfg, ctx.Version.Commit, t); err != nil {
			return err
		}
	}

	return sh.RunCommand(
		cfg.GetVerbose(),
		"docker",
		"push",
		cfg.GetBuilder().GetImage(),
	)
}

func buildArgs(ctx *context.Context) []string {
	var res []string

	for k, v := range map[string]string{
		"GIT_BRANCH":     ctx.Version.Branch,
		"GIT_COMMIT":     ctx.Version.Commit,
		"GIT_REMOTE":     ctx.Version.Remote,
		"SEMBER_VERSION": ctx.Version.Semver,
	} {
		res = append(res, "--build-arg", fmt.Sprintf("%s=%s", k, v))
	}

	return res
}

func tagImage(cfg *config.Configuration, from, to string) error {
	return sh.RunCommand(
		cfg.GetVerbose(),
		"docker",
		"tag",
		cfg.GetBuilder().GetImage()+":"+from,
		cfg.GetBuilder().GetImage()+":"+to,
	)
}
