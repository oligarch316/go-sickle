package command

import (
	"os"

	blade "github.com/oligarch316/go-sickle-blade"
	"github.com/oligarch316/go-sickle/meta"
	"github.com/oligarch316/go-sickle/transform"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type paramsDownload struct {
	paramsSession
	targets []string
}

func NewDownload() *cobra.Command {
	var params paramsDownload

	res := &cobra.Command{
		Use:   "download",
		Short: "Download target urls",
		Long:  "TODO",
		Run: func(_ *cobra.Command, args []string) {
			params.targets = args
			os.Exit(runDownload(params))
		},
	}

	params.SetFlags(res.Flags())
	return res
}

func runDownload(params paramsDownload) int {
	session, status := newSession(params.paramsSession)
	if status != 0 {
		return status
	}

	ctx := meta.With(session.ctx, transform.RequireMediaItem)

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
		if classified.Type() != blade.ItemTypeMedia {
			session.logger.Error(
				"internal error, unexpected classified item class",
				zap.Stringer("expected", blade.ItemTypeMedia),
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

		media, err := session.transformer.Download(ctx, classifiedURL)
		if err != nil {
			session.logger.Error("failed to download classified item", zap.Error(err))
			return 1
		}

		if err = session.consumer.ConsumeMedia(ctx, media); err != nil {
			session.logger.Error("failed to consume media", zap.Error(err))
			return 1
		}
	}

	return 0
}
