package flag

import "github.com/spf13/pflag"

type Option interface{ FlagOption(*pflag.Flag) }

type (
	OptionFunc    func(*pflag.Flag)
	OptDeprecated string
	OptHidden     bool
)

func (of OptionFunc) FlagOption(f *pflag.Flag)    { of(f) }
func (od OptDeprecated) FlagOption(f *pflag.Flag) { f.Deprecated = string(od) }
func (oh OptHidden) FlagOption(f *pflag.Flag)     { f.Hidden = bool(oh) }
