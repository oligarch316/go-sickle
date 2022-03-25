package flag

import "github.com/spf13/pflag"

type deferredItem struct {
	list *DeferredList
	pflag.Value
}

func (di deferredItem) Set(val string) error {
	fn := func() error { return di.Value.Set(val) }
	*di.list = append(*di.list, fn)
	return nil
}

type DeferredList []func() error

func (dl DeferredList) Apply() error {
	for _, fn := range dl {
		// TODO: Look at pflag/cobra for error formatting inspiration, may need to store flag names
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func (dl *DeferredList) Var(fs *pflag.FlagSet, value pflag.Value, name, usage string, opts ...Option) {
	dl.VarP(fs, value, name, "", usage, opts...)
}

func (dl *DeferredList) VarP(fs *pflag.FlagSet, value pflag.Value, name, shorthand, usage string, opts ...Option) {
	if valOpt, ok := value.(Option); ok {
		opts = append([]Option{valOpt}, opts...)
	}

	item := deferredItem{list: dl, Value: value}
	flag := fs.VarPF(item, name, shorthand, usage)

	for _, opt := range opts {
		opt.FlagOption(flag)
	}
}
