package command

import (
	"os"

	blade "github.com/oligarch316/go-sickle-blade"
	"github.com/oligarch316/go-sickle/meta"
	"github.com/oligarch316/go-sickle/transform"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type paramsCollect struct {
	paramsSession
	targets []string
}

func NewCollect() *cobra.Command {
	var params paramsCollect

	res := &cobra.Command{
		Use:   "collect",
		Short: "Collect target urls",
		Long:  "TODO",
		Run: func(_ *cobra.Command, args []string) {
			params.targets = args
			os.Exit(runCollect(params))
		},
	}

	params.SetFlags(res.Flags())
	return res
}

func runCollect(params paramsCollect) int {
	session, status := newSession(params.paramsSession)
	if status != 0 {
		return status
	}

	ctx := meta.With(session.ctx, transform.RequireCollectionItem)

	for _, target := range params.targets {
		parsedURL, err := session.transformer.Parse(ctx, target)
		if err != nil {
			session.logger.Error("failed to parse target", zap.Error(err))
			return 1
		}

		if err = session.consumer.ConsumeParsed(ctx, parsedURL); err != nil {
			session.logger.Error("failed to consume parsed item", zap.Error(err))
			return 1
		}

		classified, err := session.transformer.Classify(ctx, parsedURL)
		if err != nil {
			session.logger.Error("failed to classify parsed item", zap.Error(err))
			return 1
		}

		// TODO: ugly, type system should really be involved in this assertion
		if classified.Type() != blade.ItemTypeCollection {
			session.logger.Error(
				"internal error, unexpected classified item class",
				zap.Stringer("expected", blade.ItemTypeCollection),
				zap.Stringer("actual", classified.Type()),
			)
			return 2
		}

		if err = session.consumer.ConsumeClassified(ctx, classified); err != nil {
			session.logger.Error("failed to consume classified item", zap.Error(err))
			return 1
		}

		classifiedURL, err := classified.URL()
		if err != nil {
			session.logger.Error("failed to read url from classified item", zap.Error(err))
			return 1
		}

		collection, err := session.transformer.Collect(ctx, classifiedURL)
		if err != nil {
			session.logger.Error("failed to collect classified item", zap.Error(err))
			return 1
		}

		if err = session.consumer.ConsumeCollection(ctx, collection); err != nil {
			session.logger.Error("failed to consume collection", zap.Error(err))
			return 1
		}
	}

	return 0
}
