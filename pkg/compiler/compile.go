package compiler

import (
	"errors"

	"github.com/upfluence/ubuild/pkg/compiler/golang"
	"github.com/upfluence/ubuild/pkg/config"
	"github.com/upfluence/ubuild/pkg/context"
)

var errNotImplemented = errors.New("compile: Compiler not implemented")

func Compile(ctx *context.Context, cfg *config.Configuration) error {
	switch cfg.Type {
	case config.Go:
		return golang.Compile(ctx, cfg)
	case config.Ruby, config.Frontend, config.Python:
		return nil
	}

	return errNotImplemented
}
