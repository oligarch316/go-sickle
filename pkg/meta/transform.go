package meta

import blade "github.com/oligarch316/go-sickle-blade"

type TransformMeta struct {
	ErrorPredicates      predicateList[error]
	ClassifiedPredicates predicateList[blade.ClassifiedItem]
	CollectionPredicates predicateList[blade.CollectionItem]
	MediaPredicates      predicateList[blade.MediaItem]
}

func (tm TransformMeta) clone() TransformMeta {
	return TransformMeta{
		ErrorPredicates:      tm.ErrorPredicates.clone(),
		ClassifiedPredicates: tm.ClassifiedPredicates.clone(),
		CollectionPredicates: tm.CollectionPredicates.clone(),
		MediaPredicates:      tm.MediaPredicates.clone(),
	}
}
