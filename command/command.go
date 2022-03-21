package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type paramsSickle struct {
	paramsSession
	targets []string
}

func NewSickle() *cobra.Command {
	var params paramsSickle

	res := &cobra.Command{
		Use:   "sickle",
		Short: "TODO",
		Long:  "TODO",
		Run: func(_ *cobra.Command, args []string) {
			params.targets = args
			os.Exit(runSickle(params))
		},
	}

	params.SetFlags(res.Flags())

	res.Flags().SortFlags = false

	res.AddCommand(
		NewVersion(),
		NewConfig(),
		NewClassify(),
		NewCollect(),
		NewDownload(),
		NewSow(),
		NewReap(),
	)

	return res
}

func runSickle(params paramsSickle) int {
	fmt.Println("sickle: not yet implemented")
	return 1
}
