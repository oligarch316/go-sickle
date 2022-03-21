package transform

import (
	blade "github.com/oligarch316/go-sickle-blade"
	"github.com/oligarch316/go-sickle/meta"
)

var (
	RequireCollectionItem = requireItemType(blade.ItemTypeCollection)
	RequireMediaItem      = requireItemType(blade.ItemTypeMedia)
)

type (
	PredicateError      func(error) (accept bool)
	PredicateClassified func(blade.ClassifiedItem) (accept bool)
	PredicateCollection func(blade.CollectionItem) (accept bool)
	PredicateMedia      func(blade.MediaItem) (accept bool)
)

func (pe PredicateError) Apply(m *meta.Meta)      { m.Transform.ErrorPredicates.Append(pe) }
func (pc PredicateClassified) Apply(m *meta.Meta) { m.Transform.ClassifiedPredicates.Append(pc) }
func (pc PredicateCollection) Apply(m *meta.Meta) { m.Transform.CollectionPredicates.Append(pc) }
func (pm PredicateMedia) Apply(m *meta.Meta)      { m.Transform.MediaPredicates.Append(pm) }

func requireItemType(itemType blade.ItemType) PredicateClassified {
	return func(item blade.ClassifiedItem) bool { return item.Type() == itemType }
}
