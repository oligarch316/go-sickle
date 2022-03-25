package session

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oligarch316/go-sickle/pkg/command"
	"github.com/oligarch316/go-sickle/pkg/config/data"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type configParams struct{ command.SessionParams }

func NewConfig(defaultConfig data.Config) *cobra.Command {
	var params configParams

	res := &cobra.Command{
		Use:   "config",
		Short: "Display loaded configuration",
		Long:  "TODO",
		Run: func(_ *cobra.Command, _ []string) {
			os.Exit(runConfig(params))
		},
	}

	params.Config = &defaultConfig
	params.SetFlags(res.Flags())

	return res
}

func runConfig(params configParams) int {
	bootstrap, status := command.NewBootstrap(params.BootstrapParams)
	if status != 0 {
		return status
	}

	b, err := json.MarshalIndent(bootstrap.Config, "", "  ")
	if err != nil {
		bootstrap.Logger.Error("failed to marshal config", zap.Error(err))
		return 1
	}

	fmt.Println(string(b))
	return 0
}
