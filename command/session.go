package command

import (
	"context"

	"github.com/oligarch316/go-sickle/config"
	"github.com/oligarch316/go-sickle/consume"
	"github.com/oligarch316/go-sickle/observ"
	"github.com/oligarch316/go-sickle/transform"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type paramsSession struct{ paramsBootstrap }

func (ps *paramsSession) SetFlags(fs *pflag.FlagSet) {
	// bootstrap flags
	ps.paramsBootstrap.SetFlags(fs)

	// plugin flags
	fs.StringArrayVar(&ps.flagData.Plugin.Files, "plugin-file", nil, "TODO-usage")
	fs.StringArrayVar(&ps.flagData.Plugin.Directories, "plugin-directory", nil, "TODO-usage")
	fs.StringArrayVar(&ps.flagData.Plugin.Trees, "plugin-tree", nil, "TODO-usage")

	// consumer flags
	fs.StringArrayVar(&ps.flagData.Consumer.Plugins.Any, "out", nil, "TODO-usage")
	fs.StringArrayVar(&ps.flagData.Consumer.Plugins.Parsed, "out-parsed", nil, "TODO-usage")
	fs.StringArrayVar(&ps.flagData.Consumer.Plugins.Classified, "out-classified", nil, "TODO-usage")
	fs.StringArrayVar(&ps.flagData.Consumer.Plugins.Collection, "out-collection", nil, "TODO-usage")
	fs.StringArrayVar(&ps.flagData.Consumer.Plugins.Media, "out-media", nil, "TODO-usage")

	// transformer flags
	fs.StringArrayVar(&ps.flagData.Transformer.Plugins, "use", []string{config.TransformerPluginNameAll}, "TODO-usage")
}

// TODO: interrupt stuffs
// TODO: (bootstrap) logger.Sync() on close

type sessionData struct {
	ctx         context.Context
	logger      *observ.Logger
	consumer    *consume.Dispatcher
	transformer *transform.Dispatcher
}

func newSession(params paramsSession) (*sessionData, int) {
	bootstrap, status := newBootstrap(params.paramsBootstrap)
	if status != 0 {
		return nil, status
	}

	rootLogger := bootstrap.logger.Named("root")

	registry, err := config.BuildRegistry(
		bootstrap.configData.Plugin,
		bootstrap.logger.Named("plugin"),
	)

	if err != nil {
		rootLogger.Error("failed to build plugin registry", zap.Error(err))
		return nil, 1
	}

	consumerDispatch, err := config.BuildConsumer(
		bootstrap.configData.Consumer,
		registry,
		bootstrap.logger.Named("consumerDispatch"),
	)

	if err != nil {
		rootLogger.Error("failed to build consumer dispatch", zap.Error(err))
		return nil, 1
	}

	transformerDispatch, err := config.BuildTransformer(
		bootstrap.configData.Transformer,
		registry,
		bootstrap.logger.Named("transformerDispatch"),
	)

	if err != nil {
		rootLogger.Error("failed to build transformer dispatch", zap.Error(err))
		return nil, 1
	}

	res := &sessionData{
		ctx:         context.Background(),
		logger:      rootLogger,
		consumer:    consumerDispatch,
		transformer: transformerDispatch,
	}

	return res, 0
}
