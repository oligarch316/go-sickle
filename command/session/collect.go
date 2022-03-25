package session

import (
	"os"

	blade "github.com/oligarch316/go-sickle-blade"
	"github.com/oligarch316/go-sickle/command"
	"github.com/oligarch316/go-sickle/config/data"
	"github.com/oligarch316/go-sickle/meta"
	"github.com/oligarch316/go-sickle/transform"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type collectParams struct {
	command.SessionParams
	targets []string
}

func NewCollect(defaultConfig data.Config) *cobra.Command {
	var params collectParams

	res := &cobra.Command{
		Use:   "collect",
		Short: "Collect target urls",
		Long:  "TODO",
		Run: func(_ *cobra.Command, args []string) {
			params.targets = args
			os.Exit(runCollect(params))
		},
	}

	params.Config = &defaultConfig
	params.SetFlags(res.Flags())

	return res
}

func runCollect(params collectParams) int {
	session, status := command.NewSession(params.SessionParams)
	if status != 0 {
		return status
	}

	ctx := meta.With(session.Context, transform.RequireCollectionItem)

	for _, target := range params.targets {
		parsedURL, err := session.Transformer.Parse(ctx, target)
		if err != nil {
			session.Logger.Error("failed to parse target", zap.Error(err))
			return 1
		}

		if err = session.Consumer.ConsumeParsed(ctx, parsedURL); err != nil {
			session.Logger.Error("failed to consume parsed item", zap.Error(err))
			return 1
		}

		classified, err := session.Transformer.Classify(ctx, parsedURL)
		if err != nil {
			session.Logger.Error("failed to classify parsed item", zap.Error(err))
			return 1
		}

		// TODO: ugly, type system should really be involved in this assertion
		if classified.Type() != blade.ItemTypeCollection {
			session.Logger.Error(
				"internal error, unexpected classified item class",
				zap.Stringer("expected", blade.ItemTypeCollection),
				zap.Stringer("actual", classified.Type()),
			)
			return 2
		}

		if err = session.Consumer.ConsumeClassified(ctx, classified); err != nil {
			session.Logger.Error("failed to consume classified item", zap.Error(err))
			return 1
		}

		classifiedURL, err := classified.URL()
		if err != nil {
			session.Logger.Error("failed to read url from classified item", zap.Error(err))
			return 1
		}

		collection, err := session.Transformer.Collect(ctx, classifiedURL)
		if err != nil {
			session.Logger.Error("failed to collect classified item", zap.Error(err))
			return 1
		}

		if err = session.Consumer.ConsumeCollection(ctx, collection); err != nil {
			session.Logger.Error("failed to consume collection", zap.Error(err))
			return 1
		}
	}

	return 0
}
