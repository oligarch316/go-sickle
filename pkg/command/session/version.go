package session

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type versionParams struct{}

func (pv *versionParams) SetFlags(fs *pflag.FlagSet) {}

func NewVersion() *cobra.Command {
	var params versionParams

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

func runVersion(params versionParams) int {
	fmt.Println("version: not yet implemented")
	return 1
}
