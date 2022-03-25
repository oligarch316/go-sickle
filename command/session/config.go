package session

import (
	"fmt"
	"os"

	"github.com/oligarch316/go-sickle/command"
	"github.com/oligarch316/go-sickle/config/data"
	"github.com/spf13/cobra"
)

type configParams struct{ command.SessionParams }

func NewConfig(defaultConfig data.Config) *cobra.Command {
	var params configParams

	res := &cobra.Command{
		Use:   "config",
		Short: "Display loaded configuration",
		Long:  "Display loaded configuration",
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

	// TODO: Fancier printing
	fmt.Printf("%+v\n", bootstrap.Config)
	return 0
}
