package embedded

import (
	"bytes"
	_ "embed"

	"github.com/oligarch316/go-sickle/pkg/config/data"
	"github.com/philandstuff/dhall-golang/v6"
	"github.com/philandstuff/dhall-golang/v6/binary"
	"github.com/philandstuff/dhall-golang/v6/core"
)

//go:embed assets/config.bytes
var configBytes []byte

func MustLoadConfig() data.Config {
	res, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	return res
}

func LoadConfig() (data.Config, error) {
	var res data.Config

	term, err := binary.DecodeAsCbor(bytes.NewReader(configBytes))
	if err != nil {
		return res, err
	}

	return res, dhall.Decode(core.Eval(term), &res)
}
