package command

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type paramsClassify struct {
	paramsSession
	targets []string
}

func NewClassify() *cobra.Command {
	var params paramsClassify

	res := &cobra.Command{
		Use:   "classify",
		Short: "Classify target urls",
		Long:  "TODO",
		Run: func(_ *cobra.Command, args []string) {
			params.targets = args
			os.Exit(runClassify(params))
		},
	}

	params.SetFlags(res.Flags())
	return res
}

func runClassify(params paramsClassify) int {
	session, status := newSession(params.paramsSession)
	if status != 0 {
		return status
	}

	for _, target := range params.targets {
		parsed, err := session.transformer.Parse(session.ctx, target)
		if err != nil {
			session.logger.Error("failed to parse target", zap.Error(err))
			return 1
		}

		if err = session.consumer.ConsumeParsed(session.ctx, parsed); err != nil {
			session.logger.Error("failed to consume parsed item", zap.Error(err))
			return 1
		}

		classified, err := session.transformer.Classify(session.ctx, parsed)
		if err != nil {
			session.logger.Error("failed to classify parsed item", zap.Error(err))
			return 1
		}

		if err = session.consumer.ConsumeClassified(session.ctx, classified); err != nil {
			session.logger.Error("failed to consume classified item", zap.Error(err))
			return 1
		}
	}

	return 0
}
