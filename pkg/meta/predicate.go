package meta

type predicateList[T any] []func(value T) (accpet bool)

func (pl predicateList[T]) clone() (res predicateList[T]) {
	if pl != nil {
		res = make(predicateList[T], len(pl))
		copy(res, pl)
	}
	return
}

func (pl predicateList[T]) Accept(value T) bool {
	for _, p := range pl {
		if !p(value) {
			return false
		}
	}
	return true
}

func (pl *predicateList[T]) Append(predicate func(T) bool) {
	*pl = append(*pl, predicate)
}
