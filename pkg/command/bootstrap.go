package command

import (
	"log"
	"os"

	"github.com/oligarch316/go-sickle/pkg/config"
	"github.com/oligarch316/go-sickle/pkg/config/data"
	"github.com/oligarch316/go-sickle/pkg/config/flag"
	"github.com/oligarch316/go-sickle/pkg/observ"
	"github.com/spf13/pflag"
)

const configPathNone = "none"

type BootstrapParams struct {
	configFilePath string

	*data.Config
	flag.DeferredList
}

func (bp *BootstrapParams) SetFlags(fs *pflag.FlagSet) {
	// config file path flag
	fs.StringVar(&bp.configFilePath, "config", "", "TODO-usage")

	// observ flags
	bp.Var(fs, &bp.Observ.Log.Encoding, "log-encoding", "TODO-usage")
	bp.Var(fs, &bp.Observ.Log.Level, "log-level", "TODO-usage")
	bp.Var(fs, &bp.Observ.Log.EnableCaller, "log-caller", "TODO-usage")
	bp.Var(fs, &bp.Observ.Log.EnableStacktrace, "log-stacktrace", "TODO-usage")
}

// TODO: logger.Sync() on close

type Bootstrap struct {
	Logger *observ.Logger
	Config data.Config
}

func bootstrapConfigPath(params BootstrapParams) (string, bool) {
	if params.configFilePath != "" {
		return params.configFilePath, params.configFilePath != configPathNone
	}

	defaultPath := config.DefaultFilePath()
	if defaultPath == "" {
		return "", false
	}

	if info, err := os.Stat(defaultPath); err == nil && !info.IsDir() {
		return defaultPath, true
	}

	return "", false
}

func NewBootstrap(params BootstrapParams) (*Bootstrap, int) {
	bootLogger := log.New(os.Stderr, "[bootstrap] ", log.LstdFlags)

	// Load file data into config if available
	if fpath, ok := bootstrapConfigPath(params); ok {
		if err := data.LoadConfigFile(fpath, params.Config); err != nil {
			bootLogger.Printf("failed to load config file '%s': %s\n", fpath, err)
			return nil, 1
		}
	}

	// Load flag data into config
	if err := params.DeferredList.Apply(); err != nil {
		bootLogger.Printf("failed to load config flags: %s\n", err)
		return nil, 1
	}

	// Build logger from config
	logger, err := data.BuildObserv(params.Config.Observ)
	if err != nil {
		bootLogger.Printf("failed to build logger: %s\n", err)
		return nil, 1
	}

	res := &Bootstrap{
		Logger: logger,
		Config: *params.Config,
	}

	return res, 0
}
