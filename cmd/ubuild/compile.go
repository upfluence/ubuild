package main

import (
	"flag"
	"fmt"
	"os"

	upfcfg "github.com/upfluence/pkg/cfg"
	"github.com/upfluence/pkg/log"

	"github.com/upfluence/ubuild/pkg/builder"
	"github.com/upfluence/ubuild/pkg/bumper"
	"github.com/upfluence/ubuild/pkg/compiler"
	"github.com/upfluence/ubuild/pkg/config"
	"github.com/upfluence/ubuild/pkg/context"
)

const configDefaultName = "ubuild.yml"

var configPath = flag.String("config-file", "", "config file path")

func mustParseConfiguration() (string, *config.Configuration) {
	cfg := *configPath
	dir, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	if cfg == "" {
		cfg = fmt.Sprintf("%s/%s", dir, configDefaultName)
	}

	res, err := config.ParseFile(cfg)

	if err != nil {
		log.Fatal(err)
	}

	return dir, res
}

func main() {
	flag.Parse()

	dir, cfg := mustParseConfiguration()
	ctx := context.BuildContext(dir, upfcfg.FetchBool("RELEASE", false))

	if ctx.Dist {
		v, err := bumper.Bump(ctx, cfg)

		if err != nil {
			log.Fatal(err)
		}

		ctx.Version.Semver = v.String()
	} else {
		ctx.Version.Semver = "v0.0.0-dirty"
	}

	if err := compiler.Compile(ctx, cfg); err != nil {
		log.Fatal(err)
	}

	if ctx.Dist {
		if err := builder.Build(ctx, cfg); err != nil {
			log.Fatal(err)
		}
	}

	log.Notice("Success")
}
