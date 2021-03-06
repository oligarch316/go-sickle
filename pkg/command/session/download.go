package session

import (
	"os"

	blade "github.com/oligarch316/go-sickle-blade"
	"github.com/oligarch316/go-sickle/pkg/command"
	"github.com/oligarch316/go-sickle/pkg/config/data"
	"github.com/oligarch316/go-sickle/pkg/meta"
	"github.com/oligarch316/go-sickle/pkg/transform"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type downloadParams struct {
	command.SessionParams
	targets []string
}

func NewDownload(defaultConfig data.Config) *cobra.Command {
	var params downloadParams

	res := &cobra.Command{
		Use:   "download",
		Short: "Download target urls",
		Long:  "TODO",
		Run: func(_ *cobra.Command, args []string) {
			params.targets = args
			os.Exit(runDownload(params))
		},
	}

	params.Config = &defaultConfig
	params.SetFlags(res.Flags())

	return res
}

func runDownload(params downloadParams) int {
	session, status := command.NewSession(params.SessionParams)
	if status != 0 {
		return status
	}

	ctx := meta.With(session.Context, transform.RequireMediaItem)

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
		if actual := classified.Class(); actual != blade.ItemClassMedia {
			session.Logger.Error(
				"internal error, unexpected classified item class",
				zap.Stringer("expected", blade.ItemClassMedia),
				zap.Stringer("actual", actual),
			)
			return 2
		}

		if err = session.Consumer.ConsumeClassified(ctx, classified); err != nil {
			session.Logger.Error("failed to consume classified item", zap.Error(err))
			return 1
		}

		normalURL, err := classified.NormalURL()
		if err != nil {
			session.Logger.Error("failed to read url from classified item", zap.Error(err))
			return 1
		}

		media, err := session.Transformer.Download(ctx, normalURL)
		if err != nil {
			session.Logger.Error("failed to download classified item", zap.Error(err))
			return 1
		}

		if err = session.Consumer.ConsumeMedia(ctx, media); err != nil {
			session.Logger.Error("failed to consume media", zap.Error(err))
			return 1
		}
	}

	return 0
}
