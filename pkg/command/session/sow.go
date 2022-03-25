package session

import (
	"fmt"
	"os"

	"github.com/oligarch316/go-sickle/pkg/command"
	"github.com/oligarch316/go-sickle/pkg/config/data"
	"github.com/spf13/cobra"
)

// TODO: Command to "drill up" to collection items and collect

type sowParams struct {
	command.SessionParams
	targets []string
}

func NewSow(defaultConfig data.Config) *cobra.Command {
	var params sowParams

	res := &cobra.Command{
		Use:   "sow",
		Short: "TODO",
		Long:  "TODO",
		Run: func(_ *cobra.Command, args []string) {
			params.targets = args
			os.Exit(runSow(params))
		},
	}

	params.Config = &defaultConfig
	params.SetFlags(res.Flags())

	return res
}

func runSow(params sowParams) int {
	fmt.Println("sow: not yet implemented")
	return 1
}
