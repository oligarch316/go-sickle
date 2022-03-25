package value

import (
	"strconv"

	"github.com/spf13/pflag"
)

type String string

func (s *String) Set(val string) error {
	*s = String(val)
	return nil
}

func (s String) String() string { return string(s) }
func (String) Type() string     { return "string" }

type Bool bool

func (b *Bool) Set(val string) error {
	parsed, err := strconv.ParseBool(val)
	*b = Bool(parsed)
	return err
}

func (b *Bool) FlagOption(flag *pflag.Flag) { flag.NoOptDefVal = "true" }
func (b Bool) String() string               { return strconv.FormatBool(bool(b)) }
func (Bool) Type() string                   { return "bool" }
