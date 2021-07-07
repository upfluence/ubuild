package builder

import (
	"errors"

	"github.com/upfluence/ubuild/pkg/builder/docker"
	"github.com/upfluence/ubuild/pkg/config"
	"github.com/upfluence/ubuild/pkg/context"
)

var errNotImplemented = errors.New("builder: Builder not implemented")

func Build(ctx *context.Context, cfg *config.Configuration) error {
	switch cfg.Type {
	case config.Go, config.Ruby, config.Python, config.Node:
		return docker.Build(ctx, cfg)
	case config.Frontend:
		return nil
	}

	return errNotImplemented
}
