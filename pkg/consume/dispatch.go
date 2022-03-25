package consume

import (
	"context"
	"net/url"

	blade "github.com/oligarch316/go-sickle-blade"
	"github.com/oligarch316/go-sickle/pkg/meta"
	"github.com/oligarch316/go-sickle/pkg/observ"
	"go.uber.org/zap"
)

// TODO: True set duplicate errors during register?
type consumerSet[T blade.Consumer] []T

func (cs *consumerSet[T]) Register(consumer blade.Consumer) bool {
	item, ok := consumer.(T)
	if ok {
		*cs = append(*cs, item)
	}
	return ok
}

type Dispatcher struct {
	logger *observ.Logger

	Parsed     consumerSet[blade.ParsedConsumer]
	Classified consumerSet[blade.ClassifiedConsumer]
	Collection consumerSet[blade.CollectionConsumer]
	Media      consumerSet[blade.MediaConsumer]
}

func NewDispatcher(logger *observ.Logger) *Dispatcher {
	return &Dispatcher{
		logger: logger,
	}
}

func (d *Dispatcher) Register(consumer blade.Consumer) bool {
	matched := d.Parsed.Register(consumer)
	matched = d.Classified.Register(consumer) || matched
	matched = d.Collection.Register(consumer) || matched
	matched = d.Media.Register(consumer) || matched

	return matched
}

func (d *Dispatcher) ConsumeParsed(ctx context.Context, url *url.URL) error {
	meta := meta.Get(ctx).Consume

	for _, consumer := range d.Parsed {
		if err := consumer.ConsumeParsed(ctx, url); err != nil {
			if meta.ErrorPredicates.Accept(err) {
				return err
			}

			d.logger.Warn(
				"ignoring consume error",
				zap.Error(err),
				zap.String("namesapce", consumer.Namespace()),
				zap.String("itemType", "parsed"),
				zap.Stringer("url", url),
			)
		}
	}

	return nil
}

func (d *Dispatcher) ConsumeClassified(ctx context.Context, item blade.ClassifiedItem) error {
	meta := meta.Get(ctx).Consume

	for _, consumer := range d.Classified {
		if err := consumer.ConsumeClassified(ctx, item); err != nil {
			if meta.ErrorPredicates.Accept(err) {
				return err
			}

			d.logger.Warn(
				"ignoring consume error",
				zap.Error(err),
				zap.String("namespace", consumer.Namespace()),
				zap.String("itemType", "classified"),
				zap.String("item", "TODO"),
			)
		}
	}

	return nil
}

func (d *Dispatcher) ConsumeCollection(ctx context.Context, item blade.CollectionItem) error {
	meta := meta.Get(ctx).Consume

	for _, consumer := range d.Collection {
		if err := consumer.ConsumeCollection(ctx, item); err != nil {
			if meta.ErrorPredicates.Accept(err) {
				return err
			}

			d.logger.Warn(
				"ignoring consume error",
				zap.Error(err),
				zap.String("namespace", consumer.Namespace()),
				zap.String("itemType", "collection"),
				zap.String("item", "TODO"),
			)
		}
	}

	return nil
}

func (d *Dispatcher) ConsumeMedia(ctx context.Context, item blade.MediaItem) error {
	meta := meta.Get(ctx).Consume

	for _, consumer := range d.Media {
		if err := consumer.ConsumeMedia(ctx, item); err != nil {
			if meta.ErrorPredicates.Accept(err) {
				return err
			}

			d.logger.Warn(
				"ignoring consume error",
				zap.Error(err),
				zap.String("namespace", consumer.Namespace()),
				zap.String("itemType", "collection"),
				zap.String("item", "TODO"),
			)
		}
	}

	return nil
}
