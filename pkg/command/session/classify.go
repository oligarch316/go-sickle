package session

import (
	"os"

	"github.com/oligarch316/go-sickle/pkg/command"
	"github.com/oligarch316/go-sickle/pkg/config/data"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type classifyParams struct {
	command.SessionParams
	targets []string
}

func NewClassify(defaultConfig data.Config) *cobra.Command {
	var params classifyParams

	res := &cobra.Command{
		Use:   "classify",
		Short: "Classify target urls",
		Long:  "TODO",
		Run: func(_ *cobra.Command, args []string) {
			params.targets = args
			os.Exit(runClassify(params))
		},
	}

	params.Config = &defaultConfig
	params.SetFlags(res.Flags())

	return res
}

func runClassify(params classifyParams) int {
	session, status := command.NewSession(params.SessionParams)
	if status != 0 {
		return status
	}

	for _, target := range params.targets {
		parsed, err := session.Transformer.Parse(session.Context, target)
		if err != nil {
			session.Logger.Error("failed to parse target", zap.Error(err))
			return 1
		}

		if err = session.Consumer.ConsumeParsed(session.Context, parsed); err != nil {
			session.Logger.Error("failed to consume parsed item", zap.Error(err))
			return 1
		}

		classified, err := session.Transformer.Classify(session.Context, parsed)
		if err != nil {
			session.Logger.Error("failed to classify parsed item", zap.Error(err))
			return 1
		}

		if err = session.Consumer.ConsumeClassified(session.Context, classified); err != nil {
			session.Logger.Error("failed to consume classified item", zap.Error(err))
			return 1
		}
	}

	return 0
}
