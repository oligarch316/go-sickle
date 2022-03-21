package config

import (
	"os"
	"path/filepath"
	goplugin "plugin"

	"github.com/oligarch316/go-sickle/observ"
	"github.com/oligarch316/go-sickle/plugin"
	"github.com/oligarch316/go-sickle/plugin/standard"
)

var stdProvider = standard.Provider{}

type PluginData struct {
	Files       []string `dhall:"files"`
	Directories []string `dhall:"directories"`
	Trees       []string `dhall:"trees"`
}

func MergePluginData(base, priority PluginData) PluginData {
	return PluginData{
		Files:       append(base.Files, priority.Files...),
		Directories: append(base.Directories, priority.Directories...),
		Trees:       append(base.Trees, priority.Trees...),
	}
}

func BuildRegistry(data PluginData, logger *observ.Logger) (*plugin.Registry, error) {
	res := plugin.NewRegistry(logger)

	if err := res.AddProvider(stdProvider); err != nil {
		// TODO: Mark as internal error
		return nil, err
	}

	items, err := pluginLoadData(data)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if err := res.AddPlugin(item); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func pluginLoadData(data PluginData) ([]*goplugin.Plugin, error) {
	var res []*goplugin.Plugin

	for _, fpath := range data.Files {
		p, err := pluginLoadFile(fpath)
		if err != nil {
			return nil, err
		}

		res = append(res, p)
	}

	for _, dpath := range data.Directories {
		ps, _, err := pluginLoadDirectory(dpath)
		if err != nil {
			return nil, err
		}

		res = append(res, ps...)
	}

	for _, tpath := range data.Trees {
		ps, err := pluginLoadTree(tpath)
		if err != nil {
			return nil, err
		}

		res = append(res, ps...)
	}

	return res, nil
}

func pluginLoadFile(fpath string) (*goplugin.Plugin, error) { return goplugin.Open(fpath) }

func pluginLoadDirectory(dpath string) ([]*goplugin.Plugin, []string, error) {
	entries, err := os.ReadDir(dpath)
	if err != nil {
		return nil, nil, err
	}

	var (
		plugins []*goplugin.Plugin
		subdirs []string
	)

	for _, entry := range entries {
		subpath := filepath.Join(dpath, entry.Name())

		if entry.IsDir() {
			subdirs = append(subdirs, subpath)
			continue
		}

		p, err := pluginLoadFile(subpath)
		if err != nil {
			return nil, nil, err
		}

		plugins = append(plugins, p)
	}

	return plugins, subdirs, nil
}

func pluginLoadTree(tpath string) ([]*goplugin.Plugin, error) {
	agg, rest, err := pluginLoadDirectory(tpath)
	if err != nil {
		return nil, err
	}

	pop := func() (next string, ok bool) {
		if ok = len(rest) > 0; !ok {
			return
		}

		next, rest = rest[0], rest[1:]
		return
	}

	for next, ok := pop(); ok; next, ok = pop() {
		plugins, more, err := pluginLoadDirectory(next)
		if err != nil {
			return nil, err
		}

		agg = append(agg, plugins...)
		rest = append(rest, more...)
	}

	return agg, nil
}
