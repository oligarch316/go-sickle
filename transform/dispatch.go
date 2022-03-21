package transform

import (
	"context"
	"errors"
	"net/url"

	blade "github.com/oligarch316/go-sickle-blade"
	"github.com/oligarch316/go-sickle/observ"
)

// TODO: Add transformer set (mimicing consumerSet)
// During classify/collect/download, try hostMux 1st, then iterate through
// transformerSet items the same way as consume.Dispatcher
// Use case is to facilitate stuff like file:// schemes, etc.

type Dispatcher struct {
	logger *observ.Logger

	hostMux
}

func NewDispatcher(logger *observ.Logger) *Dispatcher {
	return &Dispatcher{
		logger:  logger,
		hostMux: newHostMux(logger.Named("host")),
	}
}

func (d *Dispatcher) registerTransformer(handler blade.Transformer) error {
	return errors.New("naked transformers not yet implemented")
}

func (d *Dispatcher) registerHostHandler(handler blade.HostHandler) error {
	return d.hostMux.add(handler)
}

func (d *Dispatcher) registerPathHandler(handler blade.PathHandler) error {
	pathMux := newPathMux(handler, d.logger.Named("path"))
	for _, target := range handler.Paths() {
		if err := pathMux.add(target); err != nil {
			return err
		}
	}

	return d.registerHostHandler(pathMux)
}

func (d *Dispatcher) Register(transformer blade.Transformer) error {
	switch handler := transformer.(type) {
	case blade.PathHandler:
		return d.registerPathHandler(handler)
	case blade.HostHandler:
		return d.registerHostHandler(handler)
	default:
		return d.registerTransformer(handler)
	}
}

// TODO: So far this should just be a standalone function on the package
func (d *Dispatcher) Parse(_ context.Context, raw string) (*url.URL, error) {
	return url.Parse(raw)
}
