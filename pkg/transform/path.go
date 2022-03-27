package transform

import (
	"context"

	blade "github.com/oligarch316/go-sickle-blade"
	"github.com/oligarch316/go-sickle/pkg/meta"
	"github.com/oligarch316/go-sickle/pkg/observ"
	"github.com/oligarch316/go-urlrouter/component"
	"github.com/oligarch316/go-urlrouter/graph"
	"go.uber.org/zap"
)

type pathMux struct {
	blade.HostHandler

	logger *observ.Logger
	router *component.Router[blade.PathTarget]
}

func newPathMux(host blade.HostHandler, logger *observ.Logger) pathMux {
	return pathMux{
		HostHandler: host,
		logger:      logger,
		router:      component.NewPathRouter[blade.PathTarget](),
	}
}

func (pm *pathMux) add(handler blade.PathTarget) error {
	return pm.router.Add(handler.PathPattern(), handler)
}

func (pm pathMux) Classify(ctx context.Context, hostData blade.HostData) (blade.ClassifiedItem, error) {
	var (
		query   = hostData.URL.Path
		meta    = meta.Get(ctx).Transform
		vResult = visitResult[blade.ClassifiedItem]{
			err:           errClassifierNotFound,
			errPredicate:  meta.ErrorPredicates.Accept,
			itemPredicate: meta.ClassifiedPredicates.Accept,
		}
	)

	visit := func(result *graph.SearchResult[blade.PathTarget]) bool {
		visitLogger := pm.logger.With(
			zap.String("transformation", "classify"),
			zap.String("query", query),
			zap.String("namespace", pm.Namespace()),
			zap.String("hostPattern", pm.HostPattern()),
			zap.String("pathPattern", result.Value.PathPattern()),
		)

		visitLogger.Debug("pattern match success")

		handler, ok := result.Value.(blade.PathClassifier)
		if !ok {
			visitLogger.Debug("handler match failure")
			return false
		}

		visitLogger.Debug("handler match success")

		data := blade.PathData{
			URL:  hostData.URL,
			Host: hostData.Host,
			Path: blade.PatternData{
				Parameters: result.Parameters,
				Tail:       result.Tail,
			},
		}

		item, err := handler.Classify(ctx, data)
		return vResult.accept(item, err, visitLogger)
	}

	if err := pm.router.SearchFunc(visit, query); err != nil {
		return nil, err
	}

	return vResult.item, vResult.err
}

func (pm pathMux) Collect(ctx context.Context, hostData blade.HostData) (blade.CollectionItem, error) {
	var (
		query   = hostData.URL.Path
		meta    = meta.Get(ctx).Transform
		vResult = visitResult[blade.CollectionItem]{
			err:           errCollectorNotFound,
			errPredicate:  meta.ErrorPredicates.Accept,
			itemPredicate: meta.CollectionPredicates.Accept,
		}
	)

	visit := func(result *graph.SearchResult[blade.PathTarget]) bool {
		visitLogger := pm.logger.With(
			zap.String("transformation", "collect"),
			zap.String("query", query),
			zap.String("namespace", pm.Namespace()),
			zap.String("hostPattern", pm.HostPattern()),
			zap.String("pathPattern", result.Value.PathPattern()),
		)

		visitLogger.Debug("pattern match success")

		handler, ok := result.Value.(blade.PathCollector)
		if !ok {
			visitLogger.Debug("handler match failure")
			return false
		}

		visitLogger.Debug("handler match success")

		data := blade.PathData{
			URL:  hostData.URL,
			Host: hostData.Host,
			Path: blade.PatternData{
				Parameters: result.Parameters,
				Tail:       result.Tail,
			},
		}

		item, err := handler.Collect(ctx, data)
		return vResult.accept(item, err, visitLogger)
	}

	if err := pm.router.SearchFunc(visit, query); err != nil {
		return nil, err
	}

	return vResult.item, vResult.err
}

func (pm pathMux) Download(ctx context.Context, hostData blade.HostData) (blade.MediaItem, error) {
	var (
		query   = hostData.URL.Path
		meta    = meta.Get(ctx).Transform
		vResult = visitResult[blade.MediaItem]{
			err:           errDownloaderNotFound,
			errPredicate:  meta.ErrorPredicates.Accept,
			itemPredicate: meta.MediaPredicates.Accept,
		}
	)

	visit := func(result *graph.SearchResult[blade.PathTarget]) bool {
		visitLogger := pm.logger.With(
			zap.String("transformation", "download"),
			zap.String("query", query),
			zap.String("namespace", pm.Namespace()),
			zap.String("hostPattern", pm.HostPattern()),
			zap.String("pathPattern", result.Value.PathPattern()),
		)

		visitLogger.Debug("pattern match success")

		handler, ok := result.Value.(blade.PathDownloader)
		if !ok {
			visitLogger.Debug("handler match failure")
			return false
		}

		visitLogger.Debug("handler match success")

		data := blade.PathData{
			URL:  hostData.URL,
			Host: hostData.Host,
			Path: blade.PatternData{
				Parameters: result.Parameters,
				Tail:       result.Tail,
			},
		}

		item, err := handler.Download(ctx, data)
		return vResult.accept(item, err, visitLogger)
	}

	if err := pm.router.SearchFunc(visit, query); err != nil {
		return nil, err
	}

	return vResult.item, vResult.err
}
