package standard

import (
	"context"
	"net/url"

	blade "github.com/oligarch316/go-sickle-blade"
	"go.uber.org/zap"
)

// TODOS:
// - Two logs, one info and one debug with differnt info density ?
// - Check for and use zapcore.ObjectMarshaler ?
// - Check for and use fmt.Stringer ?
// - Decide on whether this consumer should ever return error (or just log errors)

type LogConsumer struct{ logger blade.Logger }

func (LogConsumer) Namespace() string { return "log" }

func (lc LogConsumer) ConsumeParsed(_ context.Context, url *url.URL) error {
	lc.logger.Info("received parsed item", zap.Stringer("url", url))
	return nil
}

func (lc LogConsumer) ConsumeClassified(_ context.Context, item blade.ClassifiedItem) error {
	url, err := item.URL()
	if err != nil {
		return err
	}

	fingerprint, err := item.Fingerprint()
	if err != nil {
		return err
	}

	lc.logger.Info(
		"received classified item",
		zap.Stringer("type", item.Type()),
		zap.Stringer("url", url),
		zap.Binary("fingerprint", fingerprint),
	)
	return nil
}

func (lc LogConsumer) ConsumeCollection(_ context.Context, item blade.CollectionItem) error {
	lc.logger.Info("received collection item", zap.String("TODO", "TODO"))
	return nil
}

func (lc LogConsumer) ConsumeMedia(_ context.Context, item blade.MediaItem) error {
	lc.logger.Info("received media item", zap.String("TODO", "TODO"))
	return nil
}
