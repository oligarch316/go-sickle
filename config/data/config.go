package data

import (
	"io/ioutil"

	"github.com/philandstuff/dhall-golang"
)

type Config struct {
	Observ      ObservConfig      `dhall:"observ"`
	Plugin      PluginConfig      `dhall:"plugin"`
	Consumer    ConsumerConfig    `dhall:"consumer"`
	Transformer TransformerConfig `dhall:"transformer"`
}

func LoadConfigFile(fpath string, c *Config) error {
	fbytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	return dhall.Unmarshal(fbytes, c)
}
