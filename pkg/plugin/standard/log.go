package standard

import (
	"context"
	"fmt"
	"net/url"

	blade "github.com/oligarch316/go-sickle-blade"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConsumable interface{ StdLogFields() []zap.Field }

type LogConsumer struct{ logger blade.Logger }

func (LogConsumer) Namespace() string { return "log" }

func (LogConsumer) detailsMarshaler(item blade.Item) zapcore.ObjectMarshaler {
	switch tItem := item.(type) {
	case LogConsumable:
		omf := func(enc zapcore.ObjectEncoder) error {
			for _, field := range tItem.StdLogFields() {
				field.AddTo(enc)
			}
			return nil
		}
		return zapcore.ObjectMarshalerFunc(omf)
	case zapcore.ObjectMarshaler:
		return tItem
	}

	noop := func(_ zapcore.ObjectEncoder) error { return nil }
	return zapcore.ObjectMarshalerFunc(noop)
}

func (LogConsumer) marshalItem(item blade.Item, enc zapcore.ObjectEncoder) error {
	fingerprint, err := item.Fingerprint()
	if err != nil {
		return err
	}

	enc.AddString("type", fmt.Sprintf("%T", item))
	enc.AddBinary("fingerprint", fingerprint)
	return nil
}

func (lc LogConsumer) ConsumeParsed(_ context.Context, url *url.URL) error {
	lc.logger.Info("received parsed item", zap.Stringer("url", url))
	return nil
}

func (lc LogConsumer) ConsumeClassified(_ context.Context, item blade.ClassifiedItem) error {
	var itemMarshaler zapcore.ObjectMarshalerFunc = func(enc zapcore.ObjectEncoder) error {
		normalURL, err := item.NormalURL()
		if err != nil {
			return err
		}

		enc.AddString("class", item.Class().String())
		enc.AddString("normalURL", normalURL.String())
		return lc.marshalItem(item, enc)
	}

	lc.logger.Info(
		"received classified item",
		zap.Object("details", lc.detailsMarshaler(item)),
		zap.Object("item", itemMarshaler),
	)

	return nil
}

func (lc LogConsumer) ConsumeCollection(ctx context.Context, item blade.CollectionItem) error {
	var itemMarshaler zapcore.ObjectMarshalerFunc = func(enc zapcore.ObjectEncoder) error {
		childURLs, err := item.ChildURLs(ctx)
		if err != nil {
			return err
		}

		childURLStrs := make([]string, len(childURLs))
		for i, childURL := range childURLs {
			childURLStrs[i] = childURL.String()
		}

		zap.Strings("childURLs", childURLStrs).AddTo(enc)
		return lc.marshalItem(item, enc)
	}

	lc.logger.Info(
		"received collection item",
		zap.Object("details", lc.detailsMarshaler(item)),
		zap.Object("item", itemMarshaler),
	)

	return nil
}

func (lc LogConsumer) ConsumeMedia(_ context.Context, item blade.MediaItem) error {
	var itemMarshaler zapcore.ObjectMarshalerFunc = func(enc zapcore.ObjectEncoder) error {
		return lc.marshalItem(item, enc)
	}

	lc.logger.Info(
		"received media item",
		zap.String("type", fmt.Sprintf("%T", item)),
		zap.Object("details", lc.detailsMarshaler(item)),
		zap.Object("item", itemMarshaler),
	)

	return nil
}
