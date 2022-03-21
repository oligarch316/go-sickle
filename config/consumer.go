package config

import (
	"fmt"

	blade "github.com/oligarch316/go-sickle-blade"
	"github.com/oligarch316/go-sickle/consume"
	"github.com/oligarch316/go-sickle/observ"
	"github.com/oligarch316/go-sickle/plugin"
)

type ConsumerDataPlugins struct {
	Any        []string `dhall:"any"`
	Parsed     []string `dhall:"parsed"`
	Classified []string `dhall:"classified"`
	Collection []string `dhall:"collection"`
	Media      []string `dhall:"media"`
}

type ConsumerData struct {
	Plugins ConsumerDataPlugins
}

func MergeConsumerData(base, priority ConsumerData) ConsumerData {
	return ConsumerData{
		Plugins: ConsumerDataPlugins{
			Any:        append(base.Plugins.Any, priority.Plugins.Any...),
			Parsed:     append(base.Plugins.Parsed, priority.Plugins.Parsed...),
			Classified: append(base.Plugins.Classified, priority.Plugins.Classified...),
			Collection: append(base.Plugins.Collection, priority.Plugins.Collection...),
			Media:      append(base.Plugins.Media, priority.Plugins.Media...),
		},
	}
}

func BuildConsumer(data ConsumerData, registry *plugin.Registry, logger *observ.Logger) (*consume.Dispatcher, error) {
	res := consume.NewDispatcher(logger)

	buildItems := []struct {
		itemType     string
		names        []string
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
			consumer, err := registry.Consumers.Lookup(name)
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
