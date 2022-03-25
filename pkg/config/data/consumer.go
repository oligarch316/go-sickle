package data

import (
	"fmt"

	blade "github.com/oligarch316/go-sickle-blade"
	"github.com/oligarch316/go-sickle/pkg/config/value"
	"github.com/oligarch316/go-sickle/pkg/consume"
	"github.com/oligarch316/go-sickle/pkg/observ"
	"github.com/oligarch316/go-sickle/pkg/plugin"
)

type ConsumerPluginsConfig struct {
	Any        value.Set[value.String, *value.String] `dhall:"any"`
	Parsed     value.Set[value.String, *value.String] `dhall:"parsed"`
	Classified value.Set[value.String, *value.String] `dhall:"classified"`
	Collection value.Set[value.String, *value.String] `dhall:"collection"`
	Media      value.Set[value.String, *value.String] `dhall:"media"`
}

type ConsumerConfig struct {
	Plugins ConsumerPluginsConfig `dhall:"plugins"`
}

func BuildConsumer(data ConsumerConfig, registry *plugin.Registry, logger *observ.Logger) (*consume.Dispatcher, error) {
	res := consume.NewDispatcher(logger)

	buildItems := []struct {
		itemType     string
		names        []value.String
		registerFunc func(blade.Consumer) bool
	}{
		{"any", data.Plugins.Any, res.Register},
		{"parsed", data.Plugins.Parsed, res.Parsed.Register},
		{"classified", data.Plugins.Classified, res.Classified.Register},
		{"collection", data.Plugins.Collection, res.Collection.Register},
		{"media", data.Plugins.Media, res.Media.Register},
	}

	for _, item := range buildItems {
		for _, name := range item.names {
			consumer, err := registry.Consumers.Lookup(string(name))
			if err != nil {
				return nil, err
			}

			if !item.registerFunc(consumer) {
				return nil, fmt.Errorf("plugin %s cannot consume %s item", name, item.itemType)
			}
		}
	}

	return res, nil
}
