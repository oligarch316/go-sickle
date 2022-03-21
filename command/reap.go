package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// TODO: Command to "drill down" to media items and download

type paramsReap struct {
	paramsSession
	targets []string
}

func NewReap() *cobra.Command {
	var params paramsReap

	res := &cobra.Command{
		Use:   "reap",
		Short: "TODO",
		Long:  "TODO",
		Run: func(_ *cobra.Command, args []string) {
			params.targets = args
			os.Exit(runReap(params))
		},
	}

	params.SetFlags(res.Flags())
	return res
}

func runReap(params paramsReap) int {
	fmt.Println("reap: not yet implemented")
	return 1
}
