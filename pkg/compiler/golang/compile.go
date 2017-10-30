package golang

import (
	"fmt"
	"os"
	"strings"

	"github.com/upfluence/ubuild/pkg/config"
	"github.com/upfluence/ubuild/pkg/context"
	"github.com/upfluence/ubuild/pkg/sh"
)

const handlerPkg = "github.com/upfluence/pkg/thrift/handler"

func Compile(ctx *context.Context, cfg *config.Configuration) error {
	if err := os.MkdirAll(
		cfg.GetCompiler().GetDist(),
		os.ModeDir|os.ModePerm,
	); err != nil {
		return err
	}

	if ctx.Dist {
		for k, v := range map[string]string{
			"CGO_ENABLED": "0",
			"GOOS":        "linux",
			"GOARCH":      "amd64",
		} {
			os.Setenv(k, v)

			defer os.Unsetenv(k)
		}
	}

	for _, binary := range cfg.GetCompiler().Binaries {
		cmd := buildCommand(ctx, binary, cfg)

		if err := sh.RunCommand(cfg.GetVerbose(), cmd[0], cmd[1:]...); err != nil {
			return err
		}
	}

	return nil
}

func buildCommand(ctx *context.Context, binary config.Binary, cfg *config.Configuration) []string {
	ldFlagFunc := func(k, v string) string {
		return fmt.Sprintf("-X %s/vendor/%s.%s=%s", cfg.GetRepo(), handlerPkg, k, v)
	}

	ldFlags := []string{
		"-s",
		ldFlagFunc("GitCommit", ctx.Version.Commit),
		ldFlagFunc("GitBranch", ctx.Version.Branch),
		ldFlagFunc("GitRemote", ctx.Version.Remote),
		ldFlagFunc("Version", ctx.Version.Semver),
	}

	return []string{
		"go",
		"build",
		"-installsuffix",
		"netgo",
		"-installsuffix",
		"cgo",
		"-ldflags",
		strings.Join(ldFlags, " "),
		"-o",
		fmt.Sprintf("%s/%s", cfg.GetCompiler().GetDist(), binary.GetName()),
		fmt.Sprintf("%s/%s", cfg.GetRepo(), binary.GetPath()),
	}
}
