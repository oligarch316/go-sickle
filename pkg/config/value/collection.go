package value

import (
	"strings"
)

// TODO:
// Pretty constraint type inference like in the generics spec docs doesn't seem
// to apply when using Set[...] as a struct field type in the data package

// Currently no logic to wholesale replace set, which means added flags
// do not replace values from config, only append

type valStatic interface {
	String() string
	Type() string
}

type valMutable[T valStatic] interface {
	Set(string) error
	*T
}

type Set[V valStatic, VM valMutable[V]] []V

func (s *Set[V, VM]) Set(val string) error {
	var item VM = new(V)

	err := item.Set(val)
	*s = append(*s, *item)

	return err
}

func (s Set[V, VM]) String() string {
	if len(s) == 0 {
		return "<empty>"
	}

	strs := make([]string, len(s))
	for i, item := range s {
		strs[i] = item.String()
	}
	return "[" + strings.Join(strs, ",") + "]"
}

func (Set[V, VM]) Type() string {
	var item V
	return "set<" + item.Type() + ">"
}
