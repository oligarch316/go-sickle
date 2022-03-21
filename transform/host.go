package transform

import (
	"context"
	"net/url"

	blade "github.com/oligarch316/go-sickle-blade"
	"github.com/oligarch316/go-sickle/meta"
	"github.com/oligarch316/go-sickle/observ"
	"github.com/oligarch316/go-urlrouter/component"
	"github.com/oligarch316/go-urlrouter/graph"
	"go.uber.org/zap"
)

type hostMux struct {
	logger *observ.Logger
	router *component.Router[blade.HostHandler]
}

func newHostMux(logger *observ.Logger) hostMux {
	return hostMux{
		logger: logger,
		router: component.NewHostRouter[blade.HostHandler](),
	}
}

func (hm *hostMux) add(handler blade.HostHandler) error {
	return hm.router.Add(handler.HostPattern(), handler)
}

func (hm hostMux) Classify(ctx context.Context, url *url.URL) (blade.ClassifiedItem, error) {
	var (
		query   = url.Hostname()
		meta    = meta.Get(ctx).Transform
		vResult = visitResult[blade.ClassifiedItem]{
			err:           errClassifierNotFound,
			errPredicate:  meta.ErrorPredicates.Accept,
			itemPredicate: meta.ClassifiedPredicates.Accept,
		}
	)

	visit := func(result *graph.SearchResult[blade.HostHandler]) bool {
		visitLogger := hm.logger.With(
			zap.String("transformation", "classify"),
			zap.String("query", query),
			zap.String("namespace", result.Value.Namespace()),
			zap.String("hostPattern", result.Value.HostPattern()),
		)

		visitLogger.Debug("pattern match success")

		handler, ok := result.Value.(blade.HostClassifier)
		if !ok {
			visitLogger.Debug("handler match failure")
			return false
		}

		visitLogger.Debug("handler match success")

		data := blade.HostData{
			URL: url,
			Host: blade.PatternData{
				Parameters: result.Parameters,
				Tail:       result.Tail,
			},
		}

		item, err := handler.Classify(ctx, data)
		return vResult.accept(item, err, visitLogger)
	}

	if err := hm.router.SearchFunc(visit, query); err != nil {
		return nil, err
	}

	return vResult.item, vResult.err
}

func (hm hostMux) Collect(ctx context.Context, url *url.URL) (blade.CollectionItem, error) {
	var (
		query   = url.Hostname()
		meta    = meta.Get(ctx).Transform
		vResult = visitResult[blade.CollectionItem]{
			err:           errClassifierNotFound,
			errPredicate:  meta.ErrorPredicates.Accept,
			itemPredicate: meta.CollectionPredicates.Accept,
		}
	)

	visit := func(result *graph.SearchResult[blade.HostHandler]) bool {
		visitLogger := hm.logger.With(
			zap.String("transformation", "collect"),
			zap.String("query", query),
			zap.String("namespace", result.Value.Namespace()),
			zap.String("hostPattern", result.Value.HostPattern()),
		)

		visitLogger.Debug("pattern match success")

		handler, ok := result.Value.(blade.HostCollector)
		if !ok {
			visitLogger.Debug("handler match failure")
			return false
		}

		visitLogger.Debug("handler match success")

		data := blade.HostData{
			URL: url,
			Host: blade.PatternData{
				Parameters: result.Parameters,
				Tail:       result.Tail,
			},
		}

		item, err := handler.Collect(ctx, data)
		return vResult.accept(item, err, visitLogger)
	}

	if err := hm.router.SearchFunc(visit, query); err != nil {
		return nil, err
	}

	return vResult.item, vResult.err
}

func (hm hostMux) Download(ctx context.Context, url *url.URL) (blade.MediaItem, error) {
	var (
		query   = url.Hostname()
		meta    = meta.Get(ctx).Transform
		vResult = visitResult[blade.MediaItem]{
			err:           errClassifierNotFound,
			errPredicate:  meta.ErrorPredicates.Accept,
			itemPredicate: meta.MediaPredicates.Accept,
		}
	)

	visit := func(result *graph.SearchResult[blade.HostHandler]) bool {
		visitLogger := hm.logger.With(
			zap.String("transformation", "download"),
			zap.String("query", query),
			zap.String("namespace", result.Value.Namespace()),
			zap.String("hostPattern", result.Value.HostPattern()),
		)

		visitLogger.Debug("pattern match success")

		handler, ok := result.Value.(blade.HostDownloader)
		if !ok {
			visitLogger.Debug("handler match failure")
			return false
		}

		visitLogger.Debug("handler match success")

		data := blade.HostData{
			URL: url,
			Host: blade.PatternData{
				Parameters: result.Parameters,
				Tail:       result.Tail,
			},
		}

		item, err := handler.Download(ctx, data)
		return vResult.accept(item, err, visitLogger)
	}

	if err := hm.router.SearchFunc(visit, query); err != nil {
		return nil, err
	}

	return vResult.item, vResult.err
}
