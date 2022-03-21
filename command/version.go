package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type paramsVersion struct{}

func (pv *paramsVersion) SetFlags(fs *pflag.FlagSet) {}

func NewVersion() *cobra.Command {
	var params paramsVersion

	res := &cobra.Command{
		Use:   "version",
		Short: "TODO",
		Long:  "TODO",
		Run: func(_ *cobra.Command, _ []string) {
			os.Exit(runVersion(params))
		},
	}

	params.SetFlags(res.Flags())
	return res
}

func runVersion(params paramsVersion) int {
	fmt.Println("version: not yet implemented")
	return 1
}
