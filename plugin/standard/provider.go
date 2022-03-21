package standard

import blade "github.com/oligarch316/go-sickle-blade"

var info = blade.ProviderInfo{
	Name:  "std",
	Short: "Standard sickle consumers and transformers",
	Long:  "Standard sickle consumers and transformers",
}

type Provider struct{}

func (Provider) Info() blade.ProviderInfo { return info }

func (Provider) Consumers(data blade.ProviderData) ([]blade.Consumer, error) {
	res := []blade.Consumer{
		LogConsumer{logger: data.Logger.Named("log")},
		// TODO
	}

	return res, nil
}

func (Provider) Transformers(data blade.ProviderData) ([]blade.Transformer, error) {
	res := []blade.Transformer{
		// TODO
	}

	return res, nil
}
