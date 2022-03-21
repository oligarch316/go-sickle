package consume

import "github.com/oligarch316/go-sickle/meta"

type Predicate func(error) (accept bool)

func (p Predicate) Apply(m *meta.Meta) { m.Consume.ErrorPredicates.Append(p) }
