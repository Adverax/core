package generic

import "sort"

type Set[T Ordered] map[T]struct{}

func NewSet[T Ordered](values ...T) Set[T] {
	set := make(map[T]struct{})
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func (set Set[T]) Append(values ...T) {
	for _, value := range values {
		set[value] = struct{}{}
	}
}

func (set Set[T]) Len() int {
	return len(set)
}

func (set Set[T]) Add(value T) {
	set[value] = struct{}{}
}

func (set Set[T]) Remove(value T) {
	delete(set, value)
}

func (set Set[T]) Contains(value T) bool {
	_, ok := set[value]
	return ok
}

func (set Set[T]) Values() List[T] {
	values := make(List[T], 0, len(set))
	for value := range set {
		values = append(values, value)
	}
	sort.Sort(values)
	return values
}

func Union[T Ordered](lists ...[]T) List[T] {
	set := make(Set[T])
	for _, list := range lists {
		set.Append(list...)
	}
	return set.Values()
}
