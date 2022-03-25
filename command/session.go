package command

import (
	"context"

	"github.com/oligarch316/go-sickle/config/data"
	"github.com/oligarch316/go-sickle/consume"
	"github.com/oligarch316/go-sickle/observ"
	"github.com/oligarch316/go-sickle/transform"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type SessionParams struct{ BootstrapParams }

func (sp *SessionParams) SetFlags(fs *pflag.FlagSet) {
	// bootstrap flags
	sp.BootstrapParams.SetFlags(fs)

	// plugin flags
	sp.Var(fs, &sp.Plugin.Files, "plugin-file", "TODO-usage")
	sp.Var(fs, &sp.Plugin.Directories, "plugin-directory", "TODO-usage")
	sp.Var(fs, &sp.Plugin.Trees, "plugin-tree", "TODO-usage")

	// consumer flags
	sp.Var(fs, &sp.Consumer.Plugins.Any, "out", "TODO-usage")
	sp.Var(fs, &sp.Consumer.Plugins.Parsed, "out-parsed", "TODO-usage")
	sp.Var(fs, &sp.Consumer.Plugins.Classified, "out-classified", "TODO-usage")
	sp.Var(fs, &sp.Consumer.Plugins.Collection, "out-collection", "TODO-usage")
	sp.Var(fs, &sp.Consumer.Plugins.Media, "out-media", "TODO-usage")

	// transformer flags
	sp.Var(fs, &sp.Transformer.Plugins, "use", "TODO-usage")
}

// TODO: interrupt stuffs
// TODO: (bootstrap) logger.Sync() on close

type Session struct {
	Context     context.Context
	Logger      *observ.Logger
	Consumer    *consume.Dispatcher
	Transformer *transform.Dispatcher
}

func NewSession(params SessionParams) (*Session, int) {
	bootstrap, status := NewBootstrap(params.BootstrapParams)
	if status != 0 {
		return nil, status
	}

	rootLogger := bootstrap.Logger.Named("root")

	registry, err := data.BuildRegistry(
		bootstrap.Config.Plugin,
		bootstrap.Logger.Named("plugin"),
	)

	if err != nil {
		rootLogger.Error("failed to build plugin registry", zap.Error(err))
		return nil, 1
	}

	consumer, err := data.BuildConsumer(
		bootstrap.Config.Consumer,
		registry,
		bootstrap.Logger.Named("consumerDispatch"),
	)

	if err != nil {
		rootLogger.Error("failed to build consumer dispatch", zap.Error(err))
		return nil, 1
	}

	transformer, err := data.BuildTransformer(
		bootstrap.Config.Transformer,
		registry,
		bootstrap.Logger.Named("transformerDispatch"),
	)

	if err != nil {
		rootLogger.Error("failed to build transformer dispatch", zap.Error(err))
		return nil, 1
	}

	res := &Session{
		Context:     context.Background(),
		Logger:      rootLogger,
		Consumer:    consumer,
		Transformer: transformer,
	}

	return res, 0
}
