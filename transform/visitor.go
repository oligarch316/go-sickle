package transform

import (
	"errors"

	"github.com/oligarch316/go-sickle/observ"
	"go.uber.org/zap"
)

// TODO: Handle ctx.Done() during visit functions

var (
	errClassifierNotFound = errors.New("classifier not found")
	errCollectorNotFound  = errors.New("collector not found")
	errDownloaderNotFound = errors.New("downloader not found")
)

type visitResult[Item any] struct {
	item Item
	err  error

	itemPredicate func(Item) bool
	errPredicate  func(error) bool
}

func (vr *visitResult[Item]) accept(item Item, err error, logger *observ.Logger) bool {
	if err != nil {
		if vr.errPredicate(err) {
			logger.Debug("error predicate match success", zap.Error(err))
			vr.err = err
			return true
		}

		logger.Debug("error predicate match failure", zap.Error(err))
		return false
	}

	if vr.itemPredicate(item) {
		logger.Debug("item predicate match success", zap.String("item", "TODO"))
		vr.item = item
		vr.err = nil
		return true
	}

	logger.Debug("item predicate match failure", zap.String("item", "TODO"))
	return false
}
