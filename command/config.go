package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type paramsConfig struct{ paramsSession }

func NewConfig() *cobra.Command {
	var params paramsConfig

	res := &cobra.Command{
		Use:   "config",
		Short: "Display loaded configuration",
		Long:  "Display loaded configuration",
		Run: func(_ *cobra.Command, _ []string) {
			os.Exit(runConfig(params))
		},
	}

	params.SetFlags(res.Flags())
	return res
}

func runConfig(params paramsConfig) int {
	bootstrap, status := newBootstrap(params.paramsBootstrap)
	if status != 0 {
		return status
	}

	// TODO: Fancier printing
	fmt.Printf("%+v\n", bootstrap.configData)
	return 0
}
