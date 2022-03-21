package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// TODO: Command to "drill up" to collection items and collect

type paramsSow struct {
	paramsSession
	targets []string
}

func NewSow() *cobra.Command {
	var params paramsSow

	res := &cobra.Command{
		Use:   "sow",
		Short: "TODO",
		Long:  "TODO",
		Run: func(_ *cobra.Command, args []string) {
			params.targets = args
			os.Exit(runSow(params))
		},
	}

	params.SetFlags(res.Flags())
	return res
}

func runSow(params paramsSow) int {
	fmt.Println("sow: not yet implemented")
	return 1
}
