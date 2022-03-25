package data

import (
	blade "github.com/oligarch316/go-sickle-blade"
	"github.com/oligarch316/go-sickle/pkg/config/value"
	"github.com/oligarch316/go-sickle/pkg/observ"
	"github.com/oligarch316/go-sickle/pkg/plugin"
	"github.com/oligarch316/go-sickle/pkg/transform"
)

const TransformerPluginNameAll = "all"

type TransformerConfig struct {
	Plugins value.Set[value.String, *value.String] `dhall:"plugins"`
}

func BuildTransformer(data TransformerConfig, registry *plugin.Registry, logger *observ.Logger) (*transform.Dispatcher, error) {
	var (
		res   = transform.NewDispatcher(logger)
		items []blade.Transformer
	)

	for _, name := range data.Plugins {
		if name == TransformerPluginNameAll {
			items = registry.Transformers.List()
			break
		}

		item, err := registry.Transformers.Lookup(string(name))
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	for _, item := range items {
		if err := res.Register(item); err != nil {
			return nil, err
		}
	}

	return res, nil
}
