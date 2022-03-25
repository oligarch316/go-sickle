package session

import (
	"fmt"
	"os"

	"github.com/oligarch316/go-sickle/command"
	"github.com/oligarch316/go-sickle/config/data"
	"github.com/spf13/cobra"
)

type sickleParams struct {
	command.SessionParams
	targets []string
}

func NewSickle(defaultConfig data.Config) *cobra.Command {
	var params sickleParams

	res := &cobra.Command{
		Use:   "sickle",
		Short: "TODO",
		Long:  "TODO",
		Run: func(_ *cobra.Command, args []string) {
			params.targets = args
			os.Exit(runSickle(params))
		},
	}

	params.Config = &defaultConfig
	params.SetFlags(res.Flags())

	// TODO: Drop all this, should be done in cmd/main.go
	// res.AddCommand(
	// 	NewVersion(),
	// 	NewConfig(),
	// 	NewClassify(),
	// 	NewCollect(),
	// 	NewDownload(),
	// 	NewSow(),
	// 	NewReap(),
	// )

	return res
}

func runSickle(params sickleParams) int {
	fmt.Println("sickle: not yet implemented")
	return 1
}
