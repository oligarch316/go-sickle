package config

import (
	"io/ioutil"

	"github.com/philandstuff/dhall-golang"
)

type Data struct {
	Consumer    ConsumerData    `dhall:"consumer"`
	Observ      ObservData      `dhall:"observ"`
	Plugin      PluginData      `dhall:"plugin"`
	Transformer TransformerData `dhall:"transformer"`
}

func MergeData(base, priority Data) Data {
	return Data{
		Consumer:    MergeConsumerData(base.Consumer, priority.Consumer),
		Observ:      MergeObservData(base.Observ, priority.Observ),
		Plugin:      MergePluginData(base.Plugin, priority.Plugin),
		Transformer: MergeTransformerData(base.Transformer, priority.Transformer),
	}
}

func LoadData(fpath string) (Data, error) {
	var res Data

	fbytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		return res, err
	}

	return res, dhall.Unmarshal(fbytes, &res)
}
