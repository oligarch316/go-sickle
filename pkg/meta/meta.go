package meta

import "context"

type metaContextKey struct{}

var metaKey = metaContextKey{}

type Option interface{ Apply(*Meta) }

type Meta struct {
	Consume   ConsumeMeta
	Transform TransformMeta
}

func (m Meta) clone() Meta {
	return Meta{
		Consume:   m.Consume.clone(),
		Transform: m.Transform.clone(),
	}
}

func Get(ctx context.Context) Meta {
	if res, ok := ctx.Value(metaKey).(Meta); ok {
		return res
	}
	return Meta{}
}

func With(ctx context.Context, opts ...Option) context.Context {
	m := Get(ctx).clone()

	for _, opt := range opts {
		opt.Apply(&m)
	}

	return context.WithValue(ctx, metaKey, m)
}
