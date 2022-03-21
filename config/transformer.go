package config

import (
	blade "github.com/oligarch316/go-sickle-blade"
	"github.com/oligarch316/go-sickle/observ"
	"github.com/oligarch316/go-sickle/plugin"
	"github.com/oligarch316/go-sickle/transform"
)

const TransformerPluginNameAll = "all"

type TransformerData struct {
	Plugins []string `dhall:"plugins"`
}

func MergeTransformerData(base, priority TransformerData) TransformerData {
	return TransformerData{
		Plugins: append(base.Plugins, priority.Plugins...),
	}
}

func BuildTransformer(data TransformerData, registry *plugin.Registry, logger *observ.Logger) (*transform.Dispatcher, error) {
	var (
		res   = transform.NewDispatcher(logger)
		items []blade.Transformer
	)

	for _, name := range data.Plugins {
		if name == TransformerPluginNameAll {
			items = registry.Transformers.List()
			break
		}

		item, err := registry.Transformers.Lookup(name)
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
