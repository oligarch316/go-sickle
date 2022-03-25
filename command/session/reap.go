package session

import (
	"fmt"
	"os"

	"github.com/oligarch316/go-sickle/command"
	"github.com/oligarch316/go-sickle/config/data"
	"github.com/spf13/cobra"
)

// TODO: Command to "drill down" to media items and download

type reapParams struct {
	command.SessionParams
	targets []string
}

func NewReap(defaultConfig data.Config) *cobra.Command {
	var params reapParams

	res := &cobra.Command{
		Use:   "reap",
		Short: "TODO",
		Long:  "TODO",
		Run: func(_ *cobra.Command, args []string) {
			params.targets = args
			os.Exit(runReap(params))
		},
	}

	params.Config = &defaultConfig
	params.SetFlags(res.Flags())

	return res
}

func runReap(params reapParams) int {
	fmt.Println("reap: not yet implemented")
	return 1
}
