package command

import (
	"log"
	"os"

	"github.com/oligarch316/go-sickle/config"
	"github.com/oligarch316/go-sickle/observ"
	"github.com/spf13/pflag"
)

type paramsBootstrap struct {
	flagConfigFilepath string
	flagData           config.Data
}

func (pb *paramsBootstrap) SetFlags(fs *pflag.FlagSet) {
	// config file path flag
	fs.StringVar(&pb.flagConfigFilepath, "config", "", "TODO-usage")

	// observ flags
	fs.StringVar(&pb.flagData.Observ.LogLevel, "log-level", "info", "TODO-usage")
	fs.StringVar(&pb.flagData.Observ.LogEncoding, "log-encoding", "console", "TODO-usage")
	fs.BoolVar(&pb.flagData.Observ.LogCaller, "log-caller", false, "TODO-usage")
	fs.BoolVar(&pb.flagData.Observ.LogStacktrace, "log-stacktrace", false, "TODO-usage")
}

// TODO: logger.Sync() on close

type bootstrapData struct {
	logger     *observ.Logger
	configData config.Data
}

func newBootstrap(params paramsBootstrap) (*bootstrapData, int) {
	var (
		bootLogger = log.New(os.Stderr, "[bootstrap] ", log.LstdFlags)
		configData = params.flagData
	)

	if params.flagConfigFilepath != "" {
		fileData, err := config.LoadData(params.flagConfigFilepath)
		if err != nil {
			bootLogger.Printf("failed to load config file '%s': %s\n", params.flagConfigFilepath, err)
			return nil, 1
		}

		configData = config.MergeData(fileData, configData)
	}

	logger, err := config.BuildObserv(configData.Observ)
	if err != nil {
		bootLogger.Printf("failed to build logger: %s\n", err)
		return nil, 1
	}

	res := &bootstrapData{
		logger:     logger,
		configData: configData,
	}

	return res, 0
}
