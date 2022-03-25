package plugin

import (
	"errors"
	"fmt"
	"plugin"

	blade "github.com/oligarch316/go-sickle-blade"
	"github.com/oligarch316/go-sickle/pkg/observ"
	"go.uber.org/zap"
)

const entrySymbol = "Entry"

var (
	errMalformedPlugin     = errors.New("malformed plugin")
	errInvalidPluginName   = errors.New("invalid plugin name")
	errDuplicatePluginName = errors.New("duplicate plugin name")
	errUnknownPluginName   = errors.New("unknown plugin name")
)

type pluginLogger struct{ *observ.Logger }

func (pl pluginLogger) Named(name string) blade.Logger {
	return pluginLogger{pl.Logger.Named(name)}
}

func (pl pluginLogger) With(fields ...zap.Field) blade.Logger {
	return pluginLogger{pl.Logger.With(fields...)}
}

type registrySet[T any] map[string]T

func (rs registrySet[T]) add(name string, value T) error {
	if existing, exists := rs[name]; exists {
		return fmt.Errorf(
			"%w, existing: %T, duplicate: %T",
			errDuplicatePluginName, existing, value,
		)
	}

	rs[name] = value
	return nil
}

func (rs registrySet[T]) List() (res []T) {
	for _, entry := range rs {
		res = append(res, entry)
	}
	return
}

func (rs registrySet[T]) Lookup(name string) (T, error) {
	res, ok := rs[name]
	if !ok {
		return res, fmt.Errorf("%w: %s", errUnknownPluginName, name)
	}
	return res, nil
}

type Registry struct {
	logger *observ.Logger

	Consumers    registrySet[blade.Consumer]
	Transformers registrySet[blade.Transformer]
}

func validateName(name string) error {
	if name == "" {
		return fmt.Errorf("%w '%s': cannot be empty", errInvalidPluginName, name)
	}

	// TODO: Anything else nefarious?
	return nil
}

func NewRegistry(logger *observ.Logger) *Registry {
	return &Registry{
		logger: logger,

		Consumers:    make(registrySet[blade.Consumer]),
		Transformers: make(registrySet[blade.Transformer]),
	}
}

func (r *Registry) addConsumers(consumers []blade.Consumer) error {
	for _, consumer := range consumers {
		cName := consumer.Namespace()
		if err := validateName(cName); err != nil {
			return err
		}

		if err := r.Consumers.add(cName, consumer); err != nil {
			return err
		}
	}

	return nil
}

func (r *Registry) addTransformers(transformers []blade.Transformer) error {
	for _, transformer := range transformers {
		tName := transformer.Namespace()
		if err := validateName(tName); err != nil {
			return err
		}

		if err := r.Transformers.add(tName, transformer); err != nil {
			return err
		}
	}

	return nil
}

func (r *Registry) AddProvider(provider blade.Provider) error {
	var (
		pName = provider.Info().Name
		pData = blade.ProviderData{
			Logger: pluginLogger{r.logger.Named(pName)},
		}
		matched bool
	)

	if err := validateName(pName); err != nil {
		return err
	}

	if consumerP, ok := provider.(blade.ConsumerProvider); ok {
		matched = true

		consumers, err := consumerP.Consumers(pData)
		if err != nil {
			return err
		}

		if err = r.addConsumers(consumers); err != nil {
			return err
		}
	}

	if transformerP, ok := provider.(blade.TransformerProvider); ok {
		matched = true

		transformers, err := transformerP.Transformers(pData)
		if err != nil {
			return err
		}

		if err = r.addTransformers(transformers); err != nil {
			return err
		}
	}

	if !matched {
		r.logger.Warn(
			"no valid constructors detected for provider",
			zap.String("name", provider.Info().Name),
		)
	}

	return nil
}

func (r *Registry) AddPlugin(p *plugin.Plugin) error {
	sym, err := p.Lookup(entrySymbol)
	if err != nil {
		return fmt.Errorf("%w, missing entry symbol: %s", errMalformedPlugin, entrySymbol)
	}

	provider, ok := sym.(blade.Provider)
	if !ok {
		return fmt.Errorf("%w, invalid entry symbol type: %T", errMalformedPlugin, sym)
	}

	return r.AddProvider(provider)
}
