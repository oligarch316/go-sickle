package main

import (
	"github.com/oligarch316/go-sickle/pkg/command/session"
	"github.com/oligarch316/go-sickle/pkg/embedded"
)

var defaultConfig = embedded.MustLoadConfig()

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
