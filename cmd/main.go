package main

import (
	"github.com/oligarch316/go-sickle/command/session"
	"github.com/oligarch316/go-sickle/config/data"
)

// TODO: Generate and embed from package.dhall
var defaultConfig = data.Config{
	Observ: data.ObservConfig{
		Log: data.ObservLogConfig{
			Encoding: "console",
			Level:    "info",
		},
	},
}

func main() {
	root := session.NewSickle(defaultConfig)

	root.AddCommand(
		session.NewVersion(),
		session.NewConfig(defaultConfig),
		session.NewClassify(defaultConfig),
		session.NewCollect(defaultConfig),
		session.NewDownload(defaultConfig),
		session.NewSow(defaultConfig),
		session.NewReap(defaultConfig),
	)

	root.Execute()
}
